package handlers

import (
	"errors"
	"net/http/httptest"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/karima-store/internal/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockOrderService struct {
	mock.Mock
}

func (m *MockOrderService) GetOrders(userID uint, limit, offset int) ([]models.Order, int64, error) {
	args := m.Called(userID, limit, offset)
	return args.Get(0).([]models.Order), args.Get(1).(int64), args.Error(2)
}

func (m *MockOrderService) GetOrder(id uint, userID uint) (*models.Order, error) {
	args := m.Called(id, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Order), args.Error(1)
}

func (m *MockOrderService) GetOrderByNumber(orderNumber string) (*models.Order, error) {
	args := m.Called(orderNumber)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Order), args.Error(1)
}

func setupOrderHandlerTest(t *testing.T) (*fiber.App, *OrderHandler, *MockOrderService) {
	app := fiber.New()
	mockService := new(MockOrderService)
	handler := NewOrderHandler(mockService)
	return app, handler, mockService
}

func TestOrderHandler_GetOrders(t *testing.T) {
	app, handler, mockService := setupOrderHandlerTest(t)

	// Middleware to simulate authentication
	app.Use(func(c *fiber.Ctx) error {
		c.Locals("user_id", uint(1))
		return c.Next()
	})
	app.Get("/orders", handler.GetOrders)

	// Test Success
	orders := []models.Order{{ID: 1}}
	mockService.On("GetOrders", uint(1), 10, 0).Return(orders, int64(1), nil)

	req := httptest.NewRequest("GET", "/orders?page=1&limit=10", nil)
	resp, err := app.Test(req)
	assert.NoError(t, err)
	assert.Equal(t, 200, resp.StatusCode)

	// Test Unauthorized (if middleware fails/missing locals - wait, our mock middleware sets it)
	// We can test Unauthorized by simulating middleware not setting user_id
	app2 := fiber.New()
	app2.Get("/orders", handler.GetOrders)
	req2 := httptest.NewRequest("GET", "/orders", nil)
	resp2, err := app2.Test(req2)
	assert.NoError(t, err)
	assert.Equal(t, 401, resp2.StatusCode)
}

func TestOrderHandler_GetOrder(t *testing.T) {
	app, handler, mockService := setupOrderHandlerTest(t)

	app.Use(func(c *fiber.Ctx) error {
		c.Locals("user_id", uint(1))
		return c.Next()
	})
	app.Get("/orders/:id", handler.GetOrder)

	// Success
	order := &models.Order{ID: 1}
	mockService.On("GetOrder", uint(1), uint(1)).Return(order, nil)

	req := httptest.NewRequest("GET", "/orders/1", nil)
	resp, err := app.Test(req)
	assert.NoError(t, err)
	assert.Equal(t, 200, resp.StatusCode)

	// Not Found
	mockService.On("GetOrder", uint(999), uint(1)).Return(nil, errors.New("not found"))
	reqNotFound := httptest.NewRequest("GET", "/orders/999", nil)
	respNotFound, err := app.Test(reqNotFound)
	assert.NoError(t, err)
	assert.Equal(t, 404, respNotFound.StatusCode)

	// Forbidden
	mockService.On("GetOrder", uint(2), uint(1)).Return(nil, errors.New("unauthorized"))
	reqForbidden := httptest.NewRequest("GET", "/orders/2", nil)
	respForbidden, err := app.Test(reqForbidden)
	assert.NoError(t, err)
	assert.Equal(t, 403, respForbidden.StatusCode)
}

func TestOrderHandler_TrackOrder(t *testing.T) {
	app, handler, mockService := setupOrderHandlerTest(t)
	app.Get("/track", handler.TrackOrder)

	// Success
	order := &models.Order{OrderNumber: "ORD-123"}
	mockService.On("GetOrderByNumber", "ORD-123").Return(order, nil)

	req := httptest.NewRequest("GET", "/track?order_number=ORD-123", nil)
	resp, err := app.Test(req)
	assert.NoError(t, err)
	assert.Equal(t, 200, resp.StatusCode)

	// Missing order number
	reqMissing := httptest.NewRequest("GET", "/track", nil)
	respMissing, err := app.Test(reqMissing)
	assert.NoError(t, err)
	assert.Equal(t, 400, respMissing.StatusCode)

	// Not found
	mockService.On("GetOrderByNumber", "INVALID").Return(nil, errors.New("not found"))
	reqNotFound := httptest.NewRequest("GET", "/track?order_number=INVALID", nil)
	respNotFound, err := app.Test(reqNotFound)
	assert.NoError(t, err)
	assert.Equal(t, 404, respNotFound.StatusCode)
}
