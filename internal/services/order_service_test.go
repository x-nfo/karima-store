package services

import (
	"testing"

	"github.com/karima-store/internal/models"
	"github.com/karima-store/internal/repository"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"gorm.io/gorm"
)

// MockOrderRepository
type MockOrderRepository struct {
	mock.Mock
}

func (m *MockOrderRepository) Create(order *models.Order) error {
	args := m.Called(order)
	return args.Error(0)
}

func (m *MockOrderRepository) GetByID(id uint) (*models.Order, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Order), args.Error(1)
}

func (m *MockOrderRepository) GetByOrderNumber(orderNumber string) (*models.Order, error) {
	args := m.Called(orderNumber)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Order), args.Error(1)
}

func (m *MockOrderRepository) GetByUserID(userID uint, limit, offset int) ([]models.Order, int64, error) {
	args := m.Called(userID, limit, offset)
	return args.Get(0).([]models.Order), args.Get(1).(int64), args.Error(2)
}

func (m *MockOrderRepository) Update(order *models.Order) error {
	args := m.Called(order)
	return args.Error(0)
}

func (m *MockOrderRepository) UpdateStatus(id uint, status models.OrderStatus) error {
	args := m.Called(id, status)
	return args.Error(0)
}

func (m *MockOrderRepository) UpdatePaymentStatus(id uint, status models.PaymentStatus) error {
	args := m.Called(id, status)
	return args.Error(0)
}

func (m *MockOrderRepository) Delete(id uint) error {
	args := m.Called(id)
	return args.Error(0)
}

func (m *MockOrderRepository) WithTx(tx *gorm.DB) repository.OrderRepository {
	args := m.Called(tx)
	return args.Get(0).(repository.OrderRepository)
}

func TestOrderService_GetOrders(t *testing.T) {
	mockRepo := new(MockOrderRepository)
	service := NewOrderService(mockRepo)

	userID := uint(1)
	orders := []models.Order{{ID: 1}}

	// Test default limits
	mockRepo.On("GetByUserID", userID, 10, 0).Return(orders, int64(1), nil)

	result, total, err := service.GetOrders(userID, 0, 0)

	assert.NoError(t, err)
	assert.Equal(t, int64(1), total)
	assert.Equal(t, orders, result)
	mockRepo.AssertExpectations(t)
}

func TestOrderService_GetOrder(t *testing.T) {
	mockRepo := new(MockOrderRepository)
	service := NewOrderService(mockRepo)

	// Success
	order := &models.Order{ID: 1, UserID: 1}
	mockRepo.On("GetByID", uint(1)).Return(order, nil)

	result, err := service.GetOrder(1, 1)
	assert.NoError(t, err)
	assert.Equal(t, order, result)

	// Not found
	mockRepo.On("GetByID", uint(999)).Return(nil, gorm.ErrRecordNotFound)
	_, err = service.GetOrder(999, 1)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "order not found")

	// Unauthorized (Order UserID != request UserID)
	orderUnauth := &models.Order{ID: 2, UserID: 2}
	mockRepo.On("GetByID", uint(2)).Return(orderUnauth, nil)
	_, err = service.GetOrder(2, 1)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "unauthorized")
}

func TestOrderService_GetOrderByNumber(t *testing.T) {
	mockRepo := new(MockOrderRepository)
	service := NewOrderService(mockRepo)

	order := &models.Order{ID: 1, OrderNumber: "ORD-123"}
	mockRepo.On("GetByOrderNumber", "ORD-123").Return(order, nil)

	result, err := service.GetOrderByNumber("ORD-123")
	assert.NoError(t, err)
	assert.Equal(t, order, result)

	// Not found
	mockRepo.On("GetByOrderNumber", "INVALID").Return(nil, gorm.ErrRecordNotFound)
	_, err = service.GetOrderByNumber("INVALID")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "order not found")
}
