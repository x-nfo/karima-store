package services

import (
	"errors"

	"github.com/karima-store/internal/models"
	"github.com/karima-store/internal/repository"
	"gorm.io/gorm"
)

type OrderService interface {
	GetOrders(userID uint, limit, offset int) ([]models.Order, int64, error)
	GetOrder(id uint, userID uint) (*models.Order, error)
	GetOrderByNumber(orderNumber string) (*models.Order, error)
}

type orderService struct {
	orderRepo repository.OrderRepository
}

func NewOrderService(orderRepo repository.OrderRepository) OrderService {
	return &orderService{
		orderRepo: orderRepo,
	}
}

func (s *orderService) GetOrders(userID uint, limit, offset int) ([]models.Order, int64, error) {
	if limit <= 0 || limit > 50 {
		limit = 10
	}
	if offset < 0 {
		offset = 0
	}
	return s.orderRepo.GetByUserID(userID, limit, offset)
}

func (s *orderService) GetOrder(id uint, userID uint) (*models.Order, error) {
	order, err := s.orderRepo.GetByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("order not found")
		}
		return nil, err
	}

	if order.UserID != userID {
		return nil, errors.New("unauthorized")
	}

	return order, nil
}

func (s *orderService) GetOrderByNumber(orderNumber string) (*models.Order, error) {
	order, err := s.orderRepo.GetByOrderNumber(orderNumber)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("order not found")
		}
		return nil, err
	}
	return order, nil
}
