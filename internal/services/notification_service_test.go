package services

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/karima-store/internal/config"
	"github.com/karima-store/internal/database"
	"github.com/karima-store/internal/fonnte"
	"github.com/karima-store/internal/models"
	"github.com/stretchr/testify/assert"
)

func TestNotificationService_NewNotificationService(t *testing.T) {
	cfg := &config.Config{
		FonnteToken: "test-token",
		FonnteURL:   "https://api.fonnte.com/send",
	}

	service := NewNotificationService(nil, nil, cfg)
	assert.NotNil(t, service)
}

func TestNotificationService_SendWhatsAppMessage_NotConfigured(t *testing.T) {
	cfg := &config.Config{} // No Fonnte token
	service := NewNotificationService(nil, nil, cfg)

	order := &models.Order{OrderNumber: "ORD-001"}
	err := service.SendWhatsAppMessage(order, "Test", "08123456789")

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "not configured")
}

func TestNotificationService_SendWhatsAppMessage_Success(t *testing.T) {
	// Mock Fonnte API
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "POST", r.Method)

		response := map[string]interface{}{
			"status": true,
			"id":     "msg-123",
		}
		json.NewEncoder(w).Encode(response)
	}))
	defer server.Close()

	cfg := &config.Config{
		FonnteToken: "test-token",
		FonnteURL:   server.URL,
	}
	service := NewNotificationService(nil, nil, cfg)

	order := &models.Order{OrderNumber: "ORD-001"}
	err := service.SendWhatsAppMessage(order, "Test message", "08123456789")

	assert.NoError(t, err)
}

func TestNotificationService_SendWhatsAppMessage_Failed(t *testing.T) {
	// Mock Fonnte API with failure
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		response := map[string]interface{}{
			"status": false,
			"detail": "Invalid phone number",
		}
		json.NewEncoder(w).Encode(response)
	}))
	defer server.Close()

	cfg := &config.Config{
		FonnteToken: "test-token",
		FonnteURL:   server.URL,
	}
	service := NewNotificationService(nil, nil, cfg)

	order := &models.Order{OrderNumber: "ORD-001"}
	err := service.SendWhatsAppMessage(order, "Test", "invalid")

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "failed")
}

func TestNotificationService_GetWhatsAppStatus_NotConfigured(t *testing.T) {
	cfg := &config.Config{}
	service := NewNotificationService(nil, nil, cfg)

	status, err := service.GetWhatsAppStatus()

	assert.NoError(t, err)
	assert.Equal(t, "not_configured", status)
}

func TestNotificationService_GetWhatsAppStatus_Connected(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		response := map[string]interface{}{
			"status": true,
		}
		json.NewEncoder(w).Encode(response)
	}))
	defer server.Close()

	cfg := &config.Config{
		FonnteToken: "test-token",
		FonnteURL:   server.URL,
	}

	// Create fonnte client with modified URL for device endpoint
	fonnteClient := fonnte.NewClient(cfg.FonnteToken, server.URL)
	service := &notificationService{
		fonnteClient: fonnteClient,
		cfg:          cfg,
	}

	status, err := service.GetWhatsAppStatus()

	assert.NoError(t, err)
	assert.Equal(t, "connected", status)
}

func TestNotificationService_SendTestWhatsAppMessage(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		response := map[string]interface{}{
			"status": true,
			"id":     "test-123",
		}
		json.NewEncoder(w).Encode(response)
	}))
	defer server.Close()

	cfg := &config.Config{
		FonnteToken: "test-token",
		FonnteURL:   server.URL,
	}
	service := NewNotificationService(nil, nil, cfg)

	err := service.SendTestWhatsAppMessage("08123456789", "Test message")

	assert.NoError(t, err)
}

func TestNotificationService_SendOrderCreatedNotification_NotConfigured(t *testing.T) {
	cfg := &config.Config{}
	service := NewNotificationService(nil, nil, cfg)

	order := &models.Order{
		OrderNumber:   "ORD-001",
		TotalAmount:   100000,
		ShippingPhone: "08123456789",
	}

	err := service.SendOrderCreatedNotification(order)

	// Should not error when not configured, just skip
	assert.NoError(t, err)
}

func TestNotificationService_SendPaymentSuccessNotification_NotConfigured(t *testing.T) {
	cfg := &config.Config{}
	service := NewNotificationService(nil, nil, cfg)

	order := &models.Order{
		OrderNumber:   "ORD-001",
		TotalAmount:   100000,
		ShippingPhone: "08123456789",
	}

	err := service.SendPaymentSuccessNotification(order)

	assert.NoError(t, err)
}

func TestNotificationService_ProcessWhatsAppWebhook(t *testing.T) {
	cfg := &config.Config{}
	service := NewNotificationService(nil, nil, cfg)

	data := map[string]interface{}{
		"event": "message_received",
		"from":  "628123456789",
	}

	err := service.ProcessWhatsAppWebhook(data)

	assert.NoError(t, err)
}

func TestNotificationService_GetWhatsAppWebhookURL(t *testing.T) {
	cfg := &config.Config{}
	service := NewNotificationService(nil, nil, cfg)

	url := service.GetWhatsAppWebhookURL()

	assert.Contains(t, url, "webhook")
	assert.Contains(t, url, "whatsapp")
}

func TestNotificationService_GetDB(t *testing.T) {
	mockDB := &database.PostgreSQL{}
	cfg := &config.Config{}
	service := NewNotificationService(mockDB, nil, cfg)

	db := service.GetDB()

	assert.NotNil(t, db)
}
