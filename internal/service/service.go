package service

import (
	"context"

	"github.com/cmrd-a/gophermart/internal/domain"
	"github.com/cmrd-a/gophermart/internal/repository"
)

type Service struct {
	repo repository.Repository
}

func NewService(repo repository.Repository) *Service {
	return &Service{repo: repo}
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
