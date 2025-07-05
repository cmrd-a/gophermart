package repository

import (
	"context"

	"github.com/shopspring/decimal"

	"github.com/cmrd-a/gophermart/internal/config"
	"github.com/cmrd-a/gophermart/internal/domain"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Pgxer interface {
	Begin(context.Context) (pgx.Tx, error)
	Exec(context.Context, string, ...any) (pgconn.CommandTag, error)
	QueryRow(context.Context, string, ...any) pgx.Row
	Query(context.Context, string, ...any) (pgx.Rows, error)
	Ping(context.Context) error
	Close()
}

type Repository struct {
	Pgxer
}

func NewRepository() (*Repository, error) {
	pool, err := pgxpool.New(context.Background(), config.Config.DatabaseURI)
	if err != nil {
		return nil, err
	}
	return &Repository{pool}, nil
}

func (r *Repository) InsertUser(ctx context.Context, login, password string) (id int64, err error) {
	err = r.QueryRow(ctx, "INSERT INTO users (login, password) VALUES ($1, $2) RETURNING id", login, password).Scan(&id)
	return id, err
}
func (r *Repository) GetUserID(ctx context.Context, login, password string) (id int64, err error) {
	err = r.QueryRow(ctx, "SELECT id FROM users WHERE login=$1 AND password=$2", login, password).Scan(&id)
	return id, err
}

func (r *Repository) AddOrder(ctx context.Context, orderNumber string, userID int64) error {
	q := "INSERT INTO orders (number, status, user_id) VALUES ($1, $2, $3)"
	_, err := r.Exec(ctx, q, orderNumber, string(domain.NEW), userID)
	return err
}

func (r *Repository) GetOrder(ctx context.Context, orderNumber string) (order domain.Order, err error) {
	q := "SELECT number, status, accrual, uploaded_at, user_id FROM orders WHERE number = $1"
	err = r.QueryRow(ctx, q, orderNumber).Scan(&order.Number, &order.Status, &order.Accrual, &order.UploadedAt, &order.UserID)
	return order, err
}

func (r *Repository) GetUserOrder(ctx context.Context, userID int64, orderNumber string) (order domain.Order, err error) {
	q := "SELECT number, status, accrual, uploaded_at, user_id FROM orders WHERE user_id=$1 AND number = $2"
	err = r.QueryRow(ctx, q, userID, orderNumber).Scan(&order.Number, &order.Status, &order.Accrual, &order.UploadedAt, &order.UserID)
	return order, err
}
func (r *Repository) GetUserOrders(ctx context.Context, userID int64) ([]domain.Order, error) {
	q := "SELECT number, status, accrual, uploaded_at, user_id FROM orders WHERE user_id = $1"
	rows, err := r.Query(ctx, q, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var orders []domain.Order
	for rows.Next() {
		var order domain.Order
		err := rows.Scan(&order.Number, &order.Status, &order.Accrual, &order.UploadedAt, &order.UserID)
		if err != nil {
			return nil, err
		}
		orders = append(orders, order)
	}

	return orders, nil
}

func (r *Repository) UpdateOrderStatus(ctx context.Context, orderNumber string, status domain.Status) error {
	_, err := r.Exec(ctx, "UPDATE orders SET status = $1 WHERE number = $2", string(status), orderNumber)
	return err
}

func (r *Repository) UpdateOrderAccrualStatus(ctx context.Context, orderNumber string, accrual decimal.Decimal, status domain.Status) error {
	q := "UPDATE orders SET accrual = $1, status = $2 WHERE number = $3"
	_, err := r.Exec(ctx, q, accrual, string(status), orderNumber)
	return err
}

func (r *Repository) GetUserBalance(ctx context.Context, userID int64) (balance domain.Balance, err error) {
	q := `
	SELECT
		COALESCE(total_accrual.sum, 0) - COALESCE(total_withdrawals.sum, 0) AS current_accrual,
		COALESCE(total_withdrawals.sum, 0) AS total_withdrawals
	FROM
		(SELECT SUM(accrual) AS sum FROM orders WHERE user_id = $1) total_accrual,
		(SELECT SUM(withdraw) AS sum FROM withdrawals WHERE user_id = $1) total_withdrawals`
	err = r.QueryRow(ctx, q, userID).Scan(&balance.Current, &balance.Withdraw)
	return balance, err
}

func (r *Repository) AddWithdraw(ctx context.Context, userID int64, orderNumber string, withdraw decimal.Decimal) error {
	q := "INSERT INTO withdrawals (user_id, order_number, withdraw) VALUES ($1, $2, $3)"
	_, err := r.Exec(ctx, q, userID, orderNumber, withdraw)
	return err
}

func (r *Repository) GetUserWithdrawals(ctx context.Context, userID int64) ([]domain.Withdraw, error) {
	q := "SELECT order_number, withdraw, processed_at FROM withdrawals WHERE user_id = $1"
	rows, err := r.Query(ctx, q, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var withdrawals []domain.Withdraw
	for rows.Next() {
		var withdraw domain.Withdraw
		err := rows.Scan(&withdraw.OrderNumber, &withdraw.Withdraw, &withdraw.ProcessedAt)
		if err != nil {
			return nil, err
		}
		withdrawals = append(withdrawals, withdraw)
	}
	return withdrawals, nil
}
