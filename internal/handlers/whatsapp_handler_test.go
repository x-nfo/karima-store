package handlers

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/karima-store/internal/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockNotificationService is a mock implementation of NotificationService
type MockNotificationService struct {
	mock.Mock
}

func (m *MockNotificationService) SendWhatsAppMessage(order *models.Order, message string, recipient string) error {
	args := m.Called(order, message, recipient)
	return args.Error(0)
}

func (m *MockNotificationService) SendOrderCreatedNotification(order *models.Order) error {
	args := m.Called(order)
	return args.Error(0)
}

func (m *MockNotificationService) SendPaymentSuccessNotification(order *models.Order) error {
	args := m.Called(order)
	return args.Error(0)
}

func (m *MockNotificationService) SendShippingNotification(order *models.Order, trackingNumber string) error {
	args := m.Called(order, trackingNumber)
	return args.Error(0)
}

func (m *MockNotificationService) GetWhatsAppStatus() (string, error) {
	args := m.Called()
	return args.String(0), args.Error(1)
}

func (m *MockNotificationService) SendTestWhatsAppMessage(phoneNumber string, message string) error {
	args := m.Called(phoneNumber, message)
	return args.Error(0)
}

func (m *MockNotificationService) ProcessWhatsAppWebhook(data map[string]interface{}) error {
	args := m.Called(data)
	return args.Error(0)
}

func (m *MockNotificationService) GetWhatsAppWebhookURL() string {
	args := m.Called()
	return args.String(0)
}

func (m *MockNotificationService) GetDB() interface{} {
	args := m.Called()
	return args.Get(0)
}

func TestWhatsAppHandler_SendTestMessage(t *testing.T) {
	tests := []struct {
		name           string
		requestBody    map[string]string
		setupMock      func(*MockNotificationService)
		expectedStatus int
	}{
		{
			name: "Success",
			requestBody: map[string]string{
				"phone":   "08123456789",
				"message": "Hello World",
			},
			setupMock: func(m *MockNotificationService) {
				m.On("SendTestWhatsAppMessage", "08123456789", "Hello World").Return(nil)
			},
			expectedStatus: http.StatusOK,
		},
		{
			name: "Missing Phone",
			requestBody: map[string]string{
				"message": "Hello World",
			},
			setupMock:      func(m *MockNotificationService) {},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name: "Service Error",
			requestBody: map[string]string{
				"phone":   "08123456789",
				"message": "Hello World",
			},
			setupMock: func(m *MockNotificationService) {
				m.On("SendTestWhatsAppMessage", "08123456789", "Hello World").Return(errors.New("service error"))
			},
			expectedStatus: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockService := new(MockNotificationService)
			tt.setupMock(mockService)

			handler := NewWhatsAppHandler(mockService)
			app := fiber.New()
			app.Post("/whatsapp/test", handler.SendTestWhatsAppMessage)

			// Send query params instead of body as per handler implementation
			params := url.Values{}
			if val, ok := tt.requestBody["phone"]; ok {
				params.Add("phone_number", val)
			}
			if val, ok := tt.requestBody["message"]; ok {
				params.Add("message", val)
			}
			req := httptest.NewRequest("POST", "/whatsapp/test?"+params.Encode(), nil)
			resp, _ := app.Test(req)

			assert.Equal(t, tt.expectedStatus, resp.StatusCode)
			mockService.AssertExpectations(t)
		})
	}
}

func TestWhatsAppHandler_GetStatus(t *testing.T) {
	tests := []struct {
		name           string
		setupMock      func(*MockNotificationService)
		expectedStatus int
	}{
		{
			name: "Connected",
			setupMock: func(m *MockNotificationService) {
				m.On("GetWhatsAppStatus").Return("connected", nil)
			},
			expectedStatus: http.StatusOK,
		},
		{
			name: "Error",
			setupMock: func(m *MockNotificationService) {
				m.On("GetWhatsAppStatus").Return("", errors.New("connection error"))
			},
			expectedStatus: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockService := new(MockNotificationService)
			tt.setupMock(mockService)

			handler := NewWhatsAppHandler(mockService)
			app := fiber.New()
			app.Get("/whatsapp/status", handler.GetWhatsAppStatus)

			req := httptest.NewRequest("GET", "/whatsapp/status", nil)
			resp, _ := app.Test(req)

			assert.Equal(t, tt.expectedStatus, resp.StatusCode)
			mockService.AssertExpectations(t)
		})
	}
}
