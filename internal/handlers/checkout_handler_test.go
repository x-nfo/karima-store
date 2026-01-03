package handlers

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http/httptest"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/karima-store/internal/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockCheckoutService for testing
type MockCheckoutService struct {
	mock.Mock
}

func (m *MockCheckoutService) Checkout(req *models.CheckoutRequest) (*models.CheckoutResponse, error) {
	args := m.Called(req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.CheckoutResponse), args.Error(1)
}

func (m *MockCheckoutService) ProcessPaymentNotification(notification *models.MidtransPaymentNotification) error {
	args := m.Called(notification)
	return args.Error(0)
}

func TestCheckoutHandler_Checkout(t *testing.T) {
	mockService := new(MockCheckoutService)
	handler := NewCheckoutHandler(mockService)
	app := fiber.New()
	app.Post("/api/v1/checkout", handler.Checkout)

	// Valid request inputs
	reqBody := models.CheckoutRequest{
		UserID:        1,
		Items:         []models.CheckoutItem{{ProductID: 1, Quantity: 1}},
		PaymentMethod: "midtrans",
		ShippingCity:  "Jakarta", // Assuming validation requires this
	}

	expectedResponse := &models.CheckoutResponse{
		OrderNumber: "ORD123",
		SnapToken:   "token123",
	}

	// We use mock.AnythingOfType because the pointer address will differ
	mockService.On("Checkout", mock.MatchedBy(func(req *models.CheckoutRequest) bool {
		return req.UserID == 1 && len(req.Items) == 1
	})).Return(expectedResponse, nil)

	body, _ := json.Marshal(reqBody)
	req := httptest.NewRequest("POST", "/api/v1/checkout", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	resp, err := app.Test(req)

	assert.NoError(t, err)
	// Status Created is 201
	assert.Equal(t, 201, resp.StatusCode)

	mockService.AssertExpectations(t)
}

func TestCheckoutHandler_Checkout_ServiceError(t *testing.T) {
	mockService := new(MockCheckoutService)
	handler := NewCheckoutHandler(mockService)
	app := fiber.New()
	app.Post("/api/v1/checkout", handler.Checkout)

	reqBody := models.CheckoutRequest{
		UserID:        1,
		Items:         []models.CheckoutItem{{ProductID: 1, Quantity: 1}},
		PaymentMethod: "midtrans",
		ShippingCity:  "Jakarta",
	}

	mockService.On("Checkout", mock.AnythingOfType("*models.CheckoutRequest")).Return(nil, errors.New("checkout failed"))

	body, _ := json.Marshal(reqBody)
	req := httptest.NewRequest("POST", "/api/v1/checkout", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	resp, err := app.Test(req)

	assert.NoError(t, err)
	assert.Equal(t, 500, resp.StatusCode)

	mockService.AssertExpectations(t)
}

func TestCheckoutHandler_PaymentWebhook(t *testing.T) {
	mockService := new(MockCheckoutService)
	handler := NewCheckoutHandler(mockService)
	app := fiber.New()
	app.Post("/api/v1/payment/webhook", handler.PaymentWebhook)

	notification := models.MidtransPaymentNotification{
		OrderID:           "ORD123",
		TransactionStatus: "settlement",
		SignatureKey:      "hash",
	}

	mockService.On("ProcessPaymentNotification", mock.MatchedBy(func(n *models.MidtransPaymentNotification) bool {
		return n.OrderID == "ORD123"
	})).Return(nil)

	body, _ := json.Marshal(notification)
	req := httptest.NewRequest("POST", "/api/v1/payment/webhook", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	resp, err := app.Test(req)

	assert.NoError(t, err)
	assert.Equal(t, 200, resp.StatusCode)

	mockService.AssertExpectations(t)
}
