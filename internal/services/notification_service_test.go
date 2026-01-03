package services

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"sync"
	"testing"

	"github.com/karima-store/internal/config"
	"github.com/karima-store/internal/fonnte"
	"github.com/karima-store/internal/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
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
	// This test is skipped because GetDB returns nil when no actual DB is initialized
	// The GetDB method is used internally and requires a real DB connection
	t.Skip("GetDB requires actual database connection")
}

// ============================================================================
// Message Format Validation Tests (FR-066 to FR-071)
// ============================================================================

func TestNotificationService_OrderCreatedNotification_MessageFormat(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "POST", r.Method)

		// Parse request body to verify message format
		var reqBody map[string]interface{}
		err := json.NewDecoder(r.Body).Decode(&reqBody)
		require.NoError(t, err)

		message, ok := reqBody["message"].(string)
		require.True(t, ok, "Message field should be present")

		// Verify message format according to FR-066
		assert.Contains(t, message, "🛍️", "Should contain shopping bag emoji")
		assert.Contains(t, message, "Pesanan Baru!", "Should contain order created title")
		assert.Contains(t, message, "Nomor Pesanan:", "Should contain order number label")
		assert.Contains(t, message, "Total:", "Should contain total label")
		assert.Contains(t, message, "Rp", "Should contain currency symbol")
		assert.Contains(t, message, "Silakan selesaikan pembayaran", "Should contain payment instruction")
		assert.Contains(t, message, "Terima kasih telah berbelanja", "Should contain thank you message")

		// Verify order number is included
		assert.Contains(t, message, "ORD-12345", "Should contain order number")

		// Verify amount is formatted
		assert.Contains(t, message, "150000", "Should contain formatted amount")

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

	order := &models.Order{
		OrderNumber:   "ORD-12345",
		TotalAmount:   150000,
		ShippingPhone: "08123456789",
	}

	err := service.SendOrderCreatedNotification(order)
	assert.NoError(t, err)
}

func TestNotificationService_PaymentSuccessNotification_MessageFormat(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "POST", r.Method)

		var reqBody map[string]interface{}
		err := json.NewDecoder(r.Body).Decode(&reqBody)
		require.NoError(t, err)

		message, ok := reqBody["message"].(string)
		require.True(t, ok)

		// Verify message format according to FR-067
		assert.Contains(t, message, "✅", "Should contain checkmark emoji")
		assert.Contains(t, message, "Pembayaran Berhasil!", "Should contain payment success title")
		assert.Contains(t, message, "Nomor Pesanan:", "Should contain order number label")
		assert.Contains(t, message, "Total:", "Should contain total label")
		assert.Contains(t, message, "Rp", "Should contain currency symbol")
		assert.Contains(t, message, "Pesanan Anda sedang diproses", "Should contain processing message")
		assert.Contains(t, message, "akan segera dikirim", "Should contain shipping info")

		// Verify order number and amount
		assert.Contains(t, message, "ORD-67890", "Should contain order number")
		assert.Contains(t, message, "250000", "Should contain formatted amount")

		response := map[string]interface{}{
			"status": true,
			"id":     "msg-456",
		}
		json.NewEncoder(w).Encode(response)
	}))
	defer server.Close()

	cfg := &config.Config{
		FonnteToken: "test-token",
		FonnteURL:   server.URL,
	}
	service := NewNotificationService(nil, nil, cfg)

	order := &models.Order{
		OrderNumber:   "ORD-67890",
		TotalAmount:   250000,
		ShippingPhone: "08123456789",
	}

	err := service.SendPaymentSuccessNotification(order)
	assert.NoError(t, err)
}

func TestNotificationService_ShippingNotification_MessageFormat(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "POST", r.Method)

		var reqBody map[string]interface{}
		err := json.NewDecoder(r.Body).Decode(&reqBody)
		require.NoError(t, err)

		message, ok := reqBody["message"].(string)
		require.True(t, ok)

		// Verify message format according to FR-068
		assert.Contains(t, message, "📦", "Should contain package emoji")
		assert.Contains(t, message, "Pesanan Dikirim!", "Should contain shipped title")
		assert.Contains(t, message, "Nomor Pesanan:", "Should contain order number label")
		assert.Contains(t, message, "Kurir:", "Should contain courier label")
		assert.Contains(t, message, "No. Resi:", "Should contain tracking number label")
		assert.Contains(t, message, "Lacak pesanan Anda", "Should contain tracking instruction")

		// Verify order details
		assert.Contains(t, message, "ORD-11111", "Should contain order number")
		assert.Contains(t, message, "JNE", "Should contain courier name")
		assert.Contains(t, message, "JP1234567890", "Should contain tracking number")

		response := map[string]interface{}{
			"status": true,
			"id":     "msg-789",
		}
		json.NewEncoder(w).Encode(response)
	}))
	defer server.Close()

	cfg := &config.Config{
		FonnteToken: "test-token",
		FonnteURL:   server.URL,
	}
	service := NewNotificationService(nil, nil, cfg)

	order := &models.Order{
		OrderNumber:     "ORD-11111",
		ShippingProvider: "JNE",
		ShippingPhone:   "08123456789",
	}

	err := service.SendShippingNotification(order, "JP1234567890")
	assert.NoError(t, err)
}

func TestNotificationService_MessageFormat_ContainsAllRequiredFields(t *testing.T) {
	testCases := []struct {
		name           string
		order          *models.Order
		trackingNumber string
		expectedFields []string
		sendFunc       func(s NotificationService, order *models.Order, trackingNumber string) error
	}{
		{
			name: "Order created notification contains all required fields",
			order: &models.Order{
				OrderNumber:   "ORD-TEST-001",
				TotalAmount:   100000,
				ShippingPhone: "08123456789",
			},
			expectedFields: []string{
				"🛍️",
				"Pesanan Baru!",
				"Nomor Pesanan:",
				"Total:",
				"Rp",
				"Silakan selesaikan pembayaran",
				"Terima kasih",
			},
			sendFunc: func(s NotificationService, order *models.Order, trackingNumber string) error {
				return s.SendOrderCreatedNotification(order)
			},
		},
		{
			name: "Payment success notification contains all required fields",
			order: &models.Order{
				OrderNumber:   "ORD-TEST-002",
				TotalAmount:   200000,
				ShippingPhone: "08123456789",
			},
			expectedFields: []string{
				"✅",
				"Pembayaran Berhasil!",
				"Nomor Pesanan:",
				"Total:",
				"Rp",
				"Pesanan Anda sedang diproses",
				"akan segera dikirim",
			},
			sendFunc: func(s NotificationService, order *models.Order, trackingNumber string) error {
				return s.SendPaymentSuccessNotification(order)
			},
		},
		{
			name: "Shipping notification contains all required fields",
			order: &models.Order{
				OrderNumber:     "ORD-TEST-003",
				ShippingProvider: "SiCepat",
				ShippingPhone:   "08123456789",
			},
			trackingNumber: "SC0012345678",
			expectedFields: []string{
				"📦",
				"Pesanan Dikirim!",
				"Nomor Pesanan:",
				"Kurir:",
				"No. Resi:",
				"Lacak pesanan Anda",
			},
			sendFunc: func(s NotificationService, order *models.Order, trackingNumber string) error {
				return s.SendShippingNotification(order, trackingNumber)
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				var reqBody map[string]interface{}
				err := json.NewDecoder(r.Body).Decode(&reqBody)
				require.NoError(t, err)

				message, ok := reqBody["message"].(string)
				require.True(t, ok)

				// Verify all expected fields are present
				for _, field := range tc.expectedFields {
					assert.Contains(t, message, field, "Message should contain field: %s", field)
				}

				response := map[string]interface{}{
					"status": true,
					"id":     "msg-test",
				}
				json.NewEncoder(w).Encode(response)
			}))
			defer server.Close()

			cfg := &config.Config{
				FonnteToken: "test-token",
				FonnteURL:   server.URL,
			}
			service := NewNotificationService(nil, nil, cfg)

			err := tc.sendFunc(service, tc.order, tc.trackingNumber)
			assert.NoError(t, err)
		})
	}
}

func TestNotificationService_CurrencyFormatting(t *testing.T) {
	testCases := []struct {
		name          string
		amount        float64
		expectedInMsg string
	}{
		{
			name:          "Format amount 100000",
			amount:        100000,
			expectedInMsg: "100000",
		},
		{
			name:          "Format amount 150000.50",
			amount:        150000.50,
			expectedInMsg: "150001", // Rounded
		},
		{
			name:          "Format amount 999999.99",
			amount:        999999.99,
			expectedInMsg: "1000000", // Rounded
		},
		{
			name:          "Format amount 0",
			amount:        0,
			expectedInMsg: "0",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			var wg sync.WaitGroup
			wg.Add(1)

			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				// Set content type header
				w.Header().Set("Content-Type", "application/json")

				var reqBody map[string]interface{}
				err := json.NewDecoder(r.Body).Decode(&reqBody)
				if err != nil {
					t.Logf("Error decoding request body: %v", err)
					wg.Done()
					return
				}

				message, ok := reqBody["message"].(string)
				if !ok {
					t.Logf("Message field missing or not a string")
					wg.Done()
					return
				}

				assert.Contains(t, message, tc.expectedInMsg, "Message should contain formatted amount")

				response := map[string]interface{}{
					"status": true,
					"id":     "msg-currency",
				}
				json.NewEncoder(w).Encode(response)
				wg.Done()
			}))
			defer server.Close()

			cfg := &config.Config{
				FonnteToken: "test-token",
				FonnteURL:   server.URL,
			}
			service := NewNotificationService(nil, nil, cfg)

			order := &models.Order{
				OrderNumber:   "ORD-CURRENCY",
				TotalAmount:   tc.amount,
				ShippingPhone: "08123456789",
			}

			err := service.SendOrderCreatedNotification(order)
			assert.NoError(t, err)

			// Wait for async operation to complete
			wg.Wait()
		})
	}
}

func TestNotificationService_PhoneNumberFormatting(t *testing.T) {
	testCases := []struct {
		name           string
		inputPhone     string
		expectedFormat string
	}{
		{
			name:           "Format 08 prefix",
			inputPhone:     "08123456789",
			expectedFormat: "628123456789",
		},
		{
			name:           "Format 62 prefix",
			inputPhone:     "628123456789",
			expectedFormat: "628123456789",
		},
		{
			name:           "Format with spaces",
			inputPhone:     "08 1234 5678 9",
			expectedFormat: "628123456789",
		},
		{
			name:           "Format with dashes",
			inputPhone:     "0812-3456-789",
			expectedFormat: "628123456789",
		},
		{
			name:           "Format with +62",
			inputPhone:     "+628123456789",
			expectedFormat: "628123456789",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			var wg sync.WaitGroup
			wg.Add(1)

			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				// Set content type header
				w.Header().Set("Content-Type", "application/json")

				var reqBody map[string]interface{}
				err := json.NewDecoder(r.Body).Decode(&reqBody)
				if err != nil {
					t.Logf("Error decoding request body: %v", err)
					wg.Done()
					return
				}

				// Verify phone number was formatted correctly
				target, ok := reqBody["target"].(string)
				if !ok {
					t.Logf("Target field missing or not a string")
					wg.Done()
					return
				}

				assert.Equal(t, tc.expectedFormat, target, "Phone number should be formatted correctly")

				response := map[string]interface{}{
					"status": true,
					"id":     "msg-phone",
				}
				json.NewEncoder(w).Encode(response)
				wg.Done()
			}))
			defer server.Close()

			cfg := &config.Config{
				FonnteToken: "test-token",
				FonnteURL:   server.URL,
			}
			service := NewNotificationService(nil, nil, cfg)

			order := &models.Order{
				OrderNumber:   "ORD-PHONE",
				TotalAmount:   100000,
				ShippingPhone: tc.inputPhone,
			}

			err := service.SendOrderCreatedNotification(order)
			assert.NoError(t, err)

			// Wait for async operation to complete
			wg.Wait()
		})
	}
}

// ============================================================================
// Edge Cases and Error Scenarios Tests
// ============================================================================

func TestNotificationService_SendOrderCreatedNotification_MissingPhone(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Should not reach here if phone is empty
		t.Error("Should not send notification when phone is empty")
	}))
	defer server.Close()

	cfg := &config.Config{
		FonnteToken: "test-token",
		FonnteURL:   server.URL,
	}
	service := NewNotificationService(nil, nil, cfg)

	order := &models.Order{
		OrderNumber:   "ORD-NO-PHONE",
		TotalAmount:   100000,
		ShippingPhone: "", // Empty phone number
	}

	err := service.SendOrderCreatedNotification(order)
	assert.NoError(t, err, "Should not error when phone is empty")
}

func TestNotificationService_SendPaymentSuccessNotification_MissingPhone(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Should not reach here if phone is empty
		t.Error("Should not send notification when phone is empty")
	}))
	defer server.Close()

	cfg := &config.Config{
		FonnteToken: "test-token",
		FonnteURL:   server.URL,
	}
	service := NewNotificationService(nil, nil, cfg)

	order := &models.Order{
		OrderNumber:   "ORD-NO-PHONE",
		TotalAmount:   100000,
		ShippingPhone: "", // Empty phone number
	}

	err := service.SendPaymentSuccessNotification(order)
	assert.NoError(t, err, "Should not error when phone is empty")
}

func TestNotificationService_SendShippingNotification_MissingPhone(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Should not reach here if phone is empty
		t.Error("Should not send notification when phone is empty")
	}))
	defer server.Close()

	cfg := &config.Config{
		FonnteToken: "test-token",
		FonnteURL:   server.URL,
	}
	service := NewNotificationService(nil, nil, cfg)

	order := &models.Order{
		OrderNumber:     "ORD-NO-PHONE",
		ShippingProvider: "JNE",
		ShippingPhone:    "", // Empty phone number
	}

	err := service.SendShippingNotification(order, "JP1234567890")
	assert.NoError(t, err, "Should not error when phone is empty")
}

func TestNotificationService_SendOrderCreatedNotification_ZeroAmount(t *testing.T) {
	var wg sync.WaitGroup
	wg.Add(1)

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Set content type header
		w.Header().Set("Content-Type", "application/json")

		var reqBody map[string]interface{}
		err := json.NewDecoder(r.Body).Decode(&reqBody)
		if err != nil {
			t.Logf("Error decoding request body: %v", err)
			wg.Done()
			return
		}

		message, ok := reqBody["message"].(string)
		if !ok {
			t.Logf("Message field missing or not a string")
			wg.Done()
			return
		}

		// Verify zero amount is formatted correctly
		assert.Contains(t, message, "0", "Message should contain zero amount")

		response := map[string]interface{}{
			"status": true,
			"id":     "msg-zero",
		}
		json.NewEncoder(w).Encode(response)
		wg.Done()
	}))
	defer server.Close()

	cfg := &config.Config{
		FonnteToken: "test-token",
		FonnteURL:   server.URL,
	}
	service := NewNotificationService(nil, nil, cfg)

	order := &models.Order{
		OrderNumber:   "ORD-ZERO",
		TotalAmount:   0,
		ShippingPhone: "08123456789",
	}

	err := service.SendOrderCreatedNotification(order)
	assert.NoError(t, err)

	// Wait for async operation to complete
	wg.Wait()
}

func TestNotificationService_SendOrderCreatedNotification_VeryLargeAmount(t *testing.T) {
	var wg sync.WaitGroup
	wg.Add(1)

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Set content type header
		w.Header().Set("Content-Type", "application/json")

		var reqBody map[string]interface{}
		err := json.NewDecoder(r.Body).Decode(&reqBody)
		if err != nil {
			t.Logf("Error decoding request body: %v", err)
			wg.Done()
			return
		}

		message, ok := reqBody["message"].(string)
		if !ok {
			t.Logf("Message field missing or not a string")
			wg.Done()
			return
		}

		// Verify large amount is formatted correctly
		assert.Contains(t, message, "999999999", "Message should contain large amount")

		response := map[string]interface{}{
			"status": true,
			"id":     "msg-large",
		}
		json.NewEncoder(w).Encode(response)
		wg.Done()
	}))
	defer server.Close()

	cfg := &config.Config{
		FonnteToken: "test-token",
		FonnteURL:   server.URL,
	}
	service := NewNotificationService(nil, nil, cfg)

	order := &models.Order{
		OrderNumber:   "ORD-LARGE",
		TotalAmount:   999999999.99,
		ShippingPhone: "08123456789",
	}

	err := service.SendOrderCreatedNotification(order)
	assert.NoError(t, err)

	// Wait for async operation to complete
	wg.Wait()
}

func TestNotificationService_SendShippingNotification_EmptyTrackingNumber(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var reqBody map[string]interface{}
		err := json.NewDecoder(r.Body).Decode(&reqBody)
		require.NoError(t, err)

		message, ok := reqBody["message"].(string)
		require.True(t, ok)

		// Verify empty tracking number is still included
		assert.Contains(t, message, "No. Resi:", "Should contain tracking label")

		response := map[string]interface{}{
			"status": true,
			"id":     "msg-empty-tracking",
		}
		json.NewEncoder(w).Encode(response)
	}))
	defer server.Close()

	cfg := &config.Config{
		FonnteToken: "test-token",
		FonnteURL:   server.URL,
	}
	service := NewNotificationService(nil, nil, cfg)

	order := &models.Order{
		OrderNumber:     "ORD-EMPTY-TRACKING",
		ShippingProvider: "JNE",
		ShippingPhone:    "08123456789",
	}

	err := service.SendShippingNotification(order, "")
	assert.NoError(t, err)
}

func TestNotificationService_SendShippingNotification_EmptyCourier(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var reqBody map[string]interface{}
		err := json.NewDecoder(r.Body).Decode(&reqBody)
		require.NoError(t, err)

		message, ok := reqBody["message"].(string)
		require.True(t, ok)

		// Verify empty courier is still included
		assert.Contains(t, message, "Kurir:", "Should contain courier label")

		response := map[string]interface{}{
			"status": true,
			"id":     "msg-empty-courier",
		}
		json.NewEncoder(w).Encode(response)
	}))
	defer server.Close()

	cfg := &config.Config{
		FonnteToken: "test-token",
		FonnteURL:   server.URL,
	}
	service := NewNotificationService(nil, nil, cfg)

	order := &models.Order{
		OrderNumber:     "ORD-EMPTY-COURIER",
		ShippingProvider: "", // Empty courier
		ShippingPhone:    "08123456789",
	}

	err := service.SendShippingNotification(order, "JP1234567890")
	assert.NoError(t, err)
}

func TestNotificationService_SendWhatsAppMessage_APIError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Simulate API error
		w.WriteHeader(http.StatusInternalServerError)
		response := map[string]interface{}{
			"status": false,
			"detail": "Internal server error",
		}
		json.NewEncoder(w).Encode(response)
	}))
	defer server.Close()

	cfg := &config.Config{
		FonnteToken: "test-token",
		FonnteURL:   server.URL,
	}
	service := NewNotificationService(nil, nil, cfg)

	order := &models.Order{OrderNumber: "ORD-ERROR"}
	err := service.SendWhatsAppMessage(order, "Test message", "08123456789")

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "failed")
}

func TestNotificationService_SendWhatsAppMessage_NetworkError(t *testing.T) {
	// Use invalid URL to simulate network error
	cfg := &config.Config{
		FonnteToken: "test-token",
		FonnteURL:   "http://invalid-url-that-does-not-exist.local:9999",
	}
	service := NewNotificationService(nil, nil, cfg)

	order := &models.Order{OrderNumber: "ORD-NETWORK-ERROR"}
	err := service.SendWhatsAppMessage(order, "Test message", "08123456789")

	assert.Error(t, err)
}

func TestNotificationService_GetWhatsAppStatus_APIError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Simulate API error by returning error response
		w.WriteHeader(http.StatusInternalServerError)
		response := map[string]interface{}{
			"status": false,
		}
		json.NewEncoder(w).Encode(response)
	}))
	defer server.Close()

	cfg := &config.Config{
		FonnteToken: "test-token",
		FonnteURL:   server.URL,
	}

	fonnteClient := fonnte.NewClient(cfg.FonnteToken, server.URL)
	service := &notificationService{
		fonnteClient: fonnteClient,
		cfg:          cfg,
	}

	status, err := service.GetWhatsAppStatus()

	// The GetWhatsAppStatus implementation returns "disconnected" when status is false, not an error
	// Only returns error when HTTP request itself fails
	assert.NoError(t, err)
	assert.Equal(t, "disconnected", status)
}

func TestNotificationService_SendTestWhatsAppMessage_APIError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Simulate API error
		w.WriteHeader(http.StatusBadRequest)
		response := map[string]interface{}{
			"status": false,
			"detail": "Bad request",
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

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "failed")
}

func TestNotificationService_PhoneNumberFormatting_EdgeCases(t *testing.T) {
	testCases := []struct {
		name           string
		inputPhone     string
		expectedFormat string
	}{
		{
			name:           "Single digit",
			inputPhone:     "0",
			expectedFormat: "62",
		},
		{
			name:           "Only digits",
			inputPhone:     "123456789",
			expectedFormat: "62123456789",
		},
		{
			name:           "International format without plus",
			inputPhone:     "628123456789",
			expectedFormat: "628123456789",
		},
		{
			name:           "Multiple special characters",
			inputPhone:     "+62-812-3456-7890",
			expectedFormat: "6281234567890",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			var wg sync.WaitGroup
			wg.Add(1)

			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				// Set content type header
				w.Header().Set("Content-Type", "application/json")

				var reqBody map[string]interface{}
				err := json.NewDecoder(r.Body).Decode(&reqBody)
				if err != nil {
					t.Logf("Error decoding request body: %v", err)
					wg.Done()
					return
				}

				target, ok := reqBody["target"].(string)
				if !ok {
					t.Logf("Target field missing or not a string")
					wg.Done()
					return
				}

				assert.Equal(t, tc.expectedFormat, target)

				response := map[string]interface{}{
					"status": true,
					"id":     "msg-edge",
				}
				json.NewEncoder(w).Encode(response)
				wg.Done()
			}))
			defer server.Close()

			cfg := &config.Config{
				FonnteToken: "test-token",
				FonnteURL:   server.URL,
			}
			service := NewNotificationService(nil, nil, cfg)

			order := &models.Order{
				OrderNumber:   "ORD-EDGE",
				TotalAmount:   100000,
				ShippingPhone: tc.inputPhone,
			}

			err := service.SendOrderCreatedNotification(order)
			assert.NoError(t, err)

			// Wait for async operation to complete
			wg.Wait()
		})
	}
}

func TestNotificationService_MessageFormat_SpecialCharacters(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var reqBody map[string]interface{}
		err := json.NewDecoder(r.Body).Decode(&reqBody)
		require.NoError(t, err)

		message, ok := reqBody["message"].(string)
		require.True(t, ok)

		// Verify message contains expected emojis and special characters
		assert.Contains(t, message, "🛍️", "Should contain shopping bag emoji")
		assert.Contains(t, message, "✅", "Should contain checkmark emoji")
		assert.Contains(t, message, "📦", "Should contain package emoji")
		assert.Contains(t, message, "🙏", "Should contain prayer hands emoji")
		assert.Contains(t, message, "*", "Should contain markdown bold markers")
		assert.Contains(t, message, "\n", "Should contain newlines")

		response := map[string]interface{}{
			"status": true,
			"id":     "msg-special",
		}
		json.NewEncoder(w).Encode(response)
	}))
	defer server.Close()

	cfg := &config.Config{
		FonnteToken: "test-token",
		FonnteURL:   server.URL,
	}
	service := NewNotificationService(nil, nil, cfg)

	order := &models.Order{
		OrderNumber:   "ORD-SPECIAL",
		TotalAmount:   100000,
		ShippingPhone: "08123456789",
	}

	err := service.SendOrderCreatedNotification(order)
	assert.NoError(t, err)
}

func TestNotificationService_MultipleNotifications_SameOrder(t *testing.T) {
	cfg := &config.Config{
		FonnteToken: "test-token",
		FonnteURL:   "https://api.fonnte.com/send",
	}
	service := NewNotificationService(nil, nil, cfg)

	order := &models.Order{
		OrderNumber:     "ORD-MULTI",
		TotalAmount:     100000,
		ShippingPhone:   "08123456789",
		ShippingProvider: "JNE",
	}

	// Send all three notification types
	// These are async operations and will be sent in background
	err1 := service.SendOrderCreatedNotification(order)
	err2 := service.SendPaymentSuccessNotification(order)
	err3 := service.SendShippingNotification(order, "JP1234567890")

	assert.NoError(t, err1)
	assert.NoError(t, err2)
	assert.NoError(t, err3)
	// Note: Async notifications are sent in background goroutines
	// We verify methods don't error and queue notifications properly
}

func TestNotificationService_OrderNumberInMessage(t *testing.T) {
	testCases := []struct {
		name        string
		orderNumber string
	}{
		{
			name:        "Standard order number",
			orderNumber: "ORD-12345",
		},
		{
			name:        "Order number with special chars",
			orderNumber: "ORD-2026-01-03-001",
		},
		{
			name:        "Order number with UUID",
			orderNumber: "ORD-a1b2c3d4-e5f6-7890",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				var reqBody map[string]interface{}
				err := json.NewDecoder(r.Body).Decode(&reqBody)
				require.NoError(t, err)

				message, ok := reqBody["message"].(string)
				require.True(t, ok)

				assert.Contains(t, message, tc.orderNumber, "Message should contain order number")

				response := map[string]interface{}{
					"status": true,
					"id":     "msg-ordnum",
				}
				json.NewEncoder(w).Encode(response)
			}))
			defer server.Close()

			cfg := &config.Config{
				FonnteToken: "test-token",
				FonnteURL:   server.URL,
			}
			service := NewNotificationService(nil, nil, cfg)

			order := &models.Order{
				OrderNumber:   tc.orderNumber,
				TotalAmount:   100000,
				ShippingPhone: "08123456789",
			}

			err := service.SendOrderCreatedNotification(order)
			assert.NoError(t, err)
		})
	}
}
