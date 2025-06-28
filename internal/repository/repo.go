package repository

import (
	"context"
	"log"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Repository struct {
	*pgxpool.Pool
}

func NewRepository() *Repository {
	uri := os.Getenv("DATABASE_URI")
	log.Printf("Connecting to database: %s", uri)
	pool, _ := pgxpool.New(context.Background(), uri)
	pool.Exec(context.Background(), `
		CREATE TABLE IF NOT EXISTS users
		(
			id       BIGSERIAL PRIMARY KEY,
			login    text NOT NULL,
			password text NOT NULL
		)
	`)
	pool.Exec(context.Background(), `
		CREATE UNIQUE INDEX IF NOT EXISTS users_login_uindex
		ON users (login)
	`)

	pool.Exec(context.Background(), `
		CREATE TABLE IF NOT EXISTS orders
		(
			number   text PRIMARY KEY,
			status   text NOT NULL,
			accural  bigint default 0,
			uploaded_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
			user_id  bigint NOT NULL
		)
	`)
	return &Repository{pool}
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

func (r *Repository) GetOrder(ctx context.Context, orderNumber string) (Order, error) {
	var order Order
	err := r.QueryRow(ctx, "SELECT number, status, accural, uploaded_at, user_id FROM orders WHERE number = $1", orderNumber).Scan(&order.Number, &order.Status, &order.Accural, &order.UploadedAt, &order.UserID)
	return order, err
}

func (r *Repository) GetUserOrders(ctx context.Context, userID int64) ([]Order, error) {
	rows, err := r.Query(ctx, "SELECT number, status, accural, uploaded_at, user_id FROM orders WHERE user_id = $1", userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var orders []Order
	for rows.Next() {
		var order Order
		err := rows.Scan(&order.Number, &order.Status, &order.Accural, &order.UploadedAt, &order.UserID)
		if err != nil {
			return nil, err
		}
		orders = append(orders, order)
	}

	return orders, nil
}
