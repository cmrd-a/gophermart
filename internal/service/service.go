package service

import (
	"context"
	"encoding/json"

	"github.com/cmrd-a/gophermart/internal/accrual"
	"github.com/cmrd-a/gophermart/internal/config"
	"github.com/cmrd-a/gophermart/internal/domain"
	"github.com/cmrd-a/gophermart/internal/repository"

	"database/sql"
	"fmt"

	_ "github.com/jackc/pgx/v5/stdlib"

	"go.dataddo.com/pgq"
)

type Service struct {
	ctx  context.Context
	repo repository.Repository
}

func NewService(ctx context.Context, repo repository.Repository) *Service {
	s := &Service{ctx: ctx, repo: repo}
	go s.consumerJob(ctx)
	return s
}

func (s *Service) AddUser(ctx context.Context, login string, password string) (int64, error) {
	return s.repo.InsertUser(ctx, login, password)
}

func (s *Service) AddOrder(ctx context.Context, orderNumber string, userID int64) error {
	return s.repo.AddOrder(ctx, orderNumber, userID)
}

func (s *Service) GetOrder(ctx context.Context, orderNumber string) *domain.Order {
	var order domain.Order
	order, err := s.repo.GetOrder(ctx, orderNumber)
	if err != nil {
		return nil
	}
	return &order
}

func (s *Service) GetUserOrders(ctx context.Context, userID int64) ([]domain.Order, error) {
	return s.repo.GetUserOrders(ctx, userID)
}

func (s *Service) Publish(orderNumber string) {
	// create a new postgres connection
	db, err := sql.Open("pgx", config.Config.DatabaseURI)
	if err != nil {
		panic(err.Error())
	}
	defer func() {
		if closeErr := db.Close(); closeErr != nil {
			fmt.Printf("Warning: failed to close database connection: %v\n", closeErr)
		}
	}()

	// create the publisher which may be reused for multiple messages
	// you may pass the optional PublisherOptions when creating it
	publisher := pgq.NewPublisher(db)

	// publish the message to the queue
	// provide the payload which is the JSON object
	// and optional metadata which is the map[string]string
	message := fmt.Sprintf(`{"order_number":"%s"}`, orderNumber)
	msg := &pgq.MessageOutgoing{
		Payload: json.RawMessage(message),
	}
	msgID, err := publisher.Publish(context.Background(), queueName, msg)
	if err != nil {
		panic(err.Error())
	}

	fmt.Println("Message published with ID:", msgID)
}

func (s *Service) consumerJob(ctx context.Context) {
	db, err := sql.Open("pgx", config.Config.DatabaseURI)
	if err != nil {
		panic(err.Error())
	}
	defer func() {
		if closeErr := db.Close(); closeErr != nil {
			fmt.Printf("Warning: failed to close database connection: %v\n", closeErr)
		}
	}()

	// create the consumer which gets attached to handling function we defined above
	h := NewHandler(s.repo)
	consumer, err := pgq.NewConsumer(db, queueName, h)
	if err != nil {
		panic(err.Error())
	}

	err = consumer.Run(ctx)
	if err != nil {
		panic(err.Error())
	}
}

// we must specify the message handler, which implements simple interface
type handler struct {
	repo repository.Repository
}

func NewHandler(repo repository.Repository) *handler {
	return &handler{repo: repo}
}

func (h *handler) HandleMessage(ctx context.Context, msg *pgq.MessageIncoming) (processed bool, err error) {
	fmt.Println("Message payload:", string(msg.Payload))

	// Parse the JSON payload to extract order_id
	var payload struct {
		OrderNumber string `json:"order_number"`
	}

	err = json.Unmarshal(msg.Payload, &payload)
	if err != nil {
		return false, err
	}

	order, err := h.repo.GetOrder(ctx, payload.OrderNumber)
	if err != nil {
		return false, err
	}
	if order.Status == string(domain.NEW) {
		err := h.repo.UpdateOrderStatus(ctx, payload.OrderNumber, string(domain.PROCESSING))
		if err != nil {
			return false, err
		}
		as, acc, err := accrual.GetAccrual(payload.OrderNumber)
		if err != nil {
			return false, err
		}
		if as == accrual.REGISTERED {
			return false, nil
		}
		if acc > 0 {
			err = h.repo.UpdateOrderAccrualStatus(ctx, payload.OrderNumber, acc, string(domain.PROCESSED))
			if err != nil {
				return false, err
			}
		}
	}

	return true, nil
}
