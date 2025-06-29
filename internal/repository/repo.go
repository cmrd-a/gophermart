package repository

import (
	"context"

	"github.com/cmrd-a/gophermart/internal/config"
	"github.com/cmrd-a/gophermart/internal/domain"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Repository struct {
	*pgxpool.Pool
}

func NewRepository() (*Repository, error) {
	pool, err := pgxpool.New(context.Background(), config.Config.DatabaseURI)
	if err != nil {
		return nil, err
	}
	_, err = pool.Exec(context.Background(), `
		CREATE TABLE IF NOT EXISTS users
		(
			id       BIGSERIAL PRIMARY KEY,
			login    text NOT NULL,
			password text NOT NULL
		)
	`)
	if err != nil {
		return nil, err
	}
	_, err = pool.Exec(context.Background(), `
		CREATE UNIQUE INDEX IF NOT EXISTS users_login_uindex
		ON users (login)
	`)
	if err != nil {
		return nil, err
	}
	_, err = pool.Exec(context.Background(), `
		CREATE TABLE IF NOT EXISTS orders
		(
			number   text PRIMARY KEY,
			status   text NOT NULL,
			accrual  bigint default 0,
			uploaded_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
			user_id  bigint NOT NULL
		)
	`)
	if err != nil {
		return nil, err
	}
	return &Repository{pool}, nil
}

func (r *Repository) InsertUser(ctx context.Context, login, password string) (int64, error) {
	var id int64
	err := r.QueryRow(ctx, "INSERT INTO users (login, password) VALUES ($1, $2) RETURNING id", login, password).Scan(&id)
	return id, err
}

func (r *Repository) AddOrder(ctx context.Context, orderNumber string, userID int64) error {
	_, err := r.Exec(ctx, "INSERT INTO orders (number, status, user_id) VALUES ($1, $2, $3)", orderNumber, "REGISTERED", userID)
	return err
}

func (r *Repository) GetOrder(ctx context.Context, orderNumber string) (domain.Order, error) {
	var order domain.Order
	err := r.QueryRow(ctx, "SELECT number, status, accrual, uploaded_at, user_id FROM orders WHERE number = $1", orderNumber).Scan(&order.Number, &order.Status, &order.Accrual, &order.UploadedAt, &order.UserID)
	return order, err
}

func (r *Repository) GetUserOrders(ctx context.Context, userID int64) ([]domain.Order, error) {
	rows, err := r.Query(ctx, "SELECT number, status, accrual, uploaded_at, user_id FROM orders WHERE user_id = $1", userID)
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

func (r *Repository) UpdateOrderStatus(ctx context.Context, orderNumber string, status string) error {
	_, err := r.Exec(ctx, "UPDATE orders SET status = $1 WHERE number = $2", status, orderNumber)
	return err
}

func (r *Repository) UpdateOrderAccrualStatus(ctx context.Context, orderNumber string, accrual int64, status string) error {
	_, err := r.Exec(ctx, "UPDATE orders SET accrual = $1, status = $2 WHERE number = $3", accrual, status, orderNumber)
	return err
}
