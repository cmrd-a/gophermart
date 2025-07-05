package service

import (
	"context"
	"encoding/json"
	"errors"
	"log"
	"net/http"

	"github.com/shopspring/decimal"

	"github.com/jackc/pgx/v5"

	"github.com/cmrd-a/gophermart/internal/accrual"
	"github.com/cmrd-a/gophermart/internal/config"
	"github.com/cmrd-a/gophermart/internal/domain"
	"github.com/cmrd-a/gophermart/internal/repository"

	"database/sql"
	"fmt"

	_ "github.com/jackc/pgx/v5/stdlib"

	"go.dataddo.com/pgq"
)

var queueName = "order_queue"

type Service struct {
	ctx  context.Context
	repo repository.Repository
}

func NewService(ctx context.Context, repo repository.Repository) *Service {
	s := &Service{ctx: ctx, repo: repo}
	go s.consumerJob(ctx)
	return s
}

func (s *Service) AddUser(ctx context.Context, login string, password string) (userID int64, err error) {
	return s.repo.InsertUser(ctx, login, password) //TODO: хэш, соль, bcrypt
}

func (s *Service) CheckLoginPassword(ctx context.Context, login string, password string) (userID int64, err error) {
	userID, err = s.repo.GetUserID(ctx, login, password)
	if errors.Is(err, pgx.ErrNoRows) {
		return 0, nil
	}
	return userID, err
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

func (s *Service) GetUserBalance(ctx context.Context, userID int64) (domain.Balance, error) {
	return s.repo.GetUserBalance(ctx, userID)
}
func (s *Service) WithdrawUserBalance(ctx context.Context, orderNumber string, userID int64, withdraw decimal.Decimal) error {
	order, err := s.repo.GetOrder(ctx, orderNumber)
	if err != nil && !errors.Is(err, pgx.ErrNoRows) {
		return err
	}
	empty := domain.Order{}
	if order != empty {
		return errors.New("old order")
	}
	err = s.AddOrder(ctx, orderNumber, userID)
	if err != nil {
		return err
	}
	log.Println(order)
	balance, err := s.repo.GetUserBalance(ctx, userID)
	if err != nil {
		return err
	}
	if balance.Current.GreaterThanOrEqual(withdraw) {
		err = s.repo.AddWithdraw(ctx, userID, withdraw)
		return err
	}
	return errors.New("insufficient balance")
}

func (s *Service) Publish(orderNumber string) {
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
	if config.Config.DatabaseURI == "" {
		return
	}
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
	consumer, err := pgq.NewConsumer(db, queueName, h, pgq.WithMaxConsumeCount(3))
	if err != nil {
		panic(err.Error())
	}

	err = consumer.Run(ctx)
	if err != nil {
		panic(err.Error())
	}
}

type Handler struct {
	repo repository.Repository
}

func NewHandler(repo repository.Repository) *Handler {
	return &Handler{repo: repo}
}

func (h *Handler) HandleMessage(ctx context.Context, msg *pgq.MessageIncoming) (processed bool, err error) {
	fmt.Println("Message payload:", string(msg.Payload))
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
	switch order.Status {
	case string(domain.NEW):
		err = h.repo.UpdateOrderStatus(ctx, payload.OrderNumber, domain.PROCESSING)
		if err != nil {
			return false, err
		}
		return h.processRequest(ctx, payload.OrderNumber)
	case string(domain.PROCESSING):
		return h.processRequest(ctx, payload.OrderNumber)
	case string(domain.PROCESSED):
		return true, nil
	case string(domain.INVALID):
		return true, nil
	}

	return true, nil
}
func (h *Handler) processRequest(ctx context.Context, orderNumber string) (processed bool, err error) {
	client := accrual.NewClient()
	acc, statusCode, err := client.GetOrderInfo(orderNumber)
	if err != nil {
		return false, err
	}
	switch statusCode {
	case http.StatusOK:
		return h.processSuccessResponse(ctx, acc)
	case http.StatusNoContent:
		err = h.repo.UpdateOrderStatus(ctx, acc.Order, domain.PROCESSED)
		if err != nil {
			return false, err
		}
		return true, nil
	case http.StatusTooManyRequests:
		return false, nil //TODO:  Retry-After: 60
	case http.StatusInternalServerError:
		return false, nil
	}
	return false, nil
}

func (h *Handler) processSuccessResponse(ctx context.Context, acc accrual.OrderInfoResponse) (processed bool, err error) {
	switch acc.Status {
	case string(accrual.REGISTERED):
		return false, nil
	case string(accrual.INVALID):
		err = h.repo.UpdateOrderStatus(ctx, acc.Order, domain.INVALID)
		if err != nil {
			return false, err
		}
		return true, err
	case string(accrual.PROCESSING):
		return false, nil
	case string(accrual.PROCESSED):
		d := decimal.NewFromFloat(acc.Accrual)
		err = h.repo.UpdateOrderAccrualStatus(ctx, acc.Order, d, domain.PROCESSED)
		if err != nil {
			return false, err
		}
		return true, nil
	}
	return false, nil
}
