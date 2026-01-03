package services

import (
	"fmt"
	"log"
	"sync"

	"github.com/karima-store/internal/config"
	"github.com/karima-store/internal/database"
	"github.com/karima-store/internal/fonnte"
	"github.com/karima-store/internal/logger"
	"github.com/karima-store/internal/models"
)

// NotificationService handles all notification-related operations
type NotificationService interface {
	SendWhatsAppMessage(order *models.Order, message string, recipient string) error
	SendOrderCreatedNotification(order *models.Order) error
	SendPaymentSuccessNotification(order *models.Order) error
	SendShippingNotification(order *models.Order, trackingNumber string) error
	GetWhatsAppStatus() (string, error)
	SendTestWhatsAppMessage(phoneNumber string, message string) error
	ProcessWhatsAppWebhook(data map[string]interface{}) error
	GetWhatsAppWebhookURL() string
	GetDB() interface{}
}

type notificationService struct {
	db           *database.PostgreSQL
	redis        database.RedisClient
	fonnteClient *fonnte.Client
	cfg          *config.Config
}

// NewNotificationService creates a new notification service instance
func NewNotificationService(db *database.PostgreSQL, redis database.RedisClient, cfg *config.Config) NotificationService {
	var fonnteClient *fonnte.Client
	if cfg.FonnteToken != "" {
		fonnteClient = fonnte.NewClient(cfg.FonnteToken, cfg.FonnteURL)
	}

	return &notificationService{
		db:           db,
		redis:        redis,
		fonnteClient: fonnteClient,
		cfg:          cfg,
	}
}

func (s *notificationService) GetDB() interface{} {
	return s.db.DB()
}

// sendWhatsAppAsync sends WhatsApp message asynchronously using goroutine
func (s *notificationService) sendWhatsAppAsync(phoneNumber, message string, wg *sync.WaitGroup) {
	if wg != nil {
		defer wg.Done()
	}

	if s.fonnteClient == nil {
		if logger.Log != nil {
			logger.Log.Warnw("Fonnte client not configured, skipping WhatsApp message", "phone", phoneNumber)
		} else {
			log.Printf("[WhatsApp] Fonnte client not configured, skipping message to %s", phoneNumber)
		}
		return
	}

	// Format phone number (ensure 628xxx format)
	formattedPhone := formatPhoneNumber(phoneNumber)

	resp, err := s.fonnteClient.SendMessage(formattedPhone, message)
	if err != nil {
		if logger.Log != nil {
			logger.Log.Errorw("Failed to send WhatsApp message", "phone", formattedPhone, "error", err)
		} else {
			log.Printf("[WhatsApp] Failed to send message to %s: %v", formattedPhone, err)
		}
		return
	}

	if resp.Status {
		if logger.Log != nil {
			logger.Log.Infow("WhatsApp message sent successfully", "phone", formattedPhone, "message_id", resp.ID)
		} else {
			log.Printf("[WhatsApp] Message sent successfully to %s (ID: %v)", formattedPhone, resp.ID)
		}
	} else {
		if logger.Log != nil {
			logger.Log.Errorw("WhatsApp message failed", "phone", formattedPhone, "detail", resp.Detail)
		} else {
			log.Printf("[WhatsApp] Message failed to %s: %s", formattedPhone, resp.Detail)
		}
	}
}

// SendWhatsAppMessage sends a message via WhatsApp Gateway API (Fonnte)
func (s *notificationService) SendWhatsAppMessage(order *models.Order, message string, recipient string) error {
	if s.fonnteClient == nil {
		if logger.Log != nil {
			logger.Log.Warnw("Fonnte client not configured")
		} else {
			log.Printf("[WhatsApp] Fonnte client not configured")
		}
		return fmt.Errorf("fonnte client not configured")
	}

	// Format phone number
	formattedPhone := formatPhoneNumber(recipient)

	resp, err := s.fonnteClient.SendMessage(formattedPhone, message)
	if err != nil {
		return fmt.Errorf("failed to send WhatsApp message: %w", err)
	}

	if !resp.Status {
		return fmt.Errorf("WhatsApp message failed: %s", resp.Detail)
	}

	if logger.Log != nil {
		logger.Log.Infow("WhatsApp message sent for order", "phone", formattedPhone, "order_number", order.OrderNumber)
	} else {
		log.Printf("[WhatsApp] Message sent to %s for order %s", formattedPhone, order.OrderNumber)
	}
	return nil
}

// SendOrderCreatedNotification sends notification when order is created (ASYNC)
func (s *notificationService) SendOrderCreatedNotification(order *models.Order) error {
	if s.fonnteClient == nil {
		if logger.Log != nil {
			logger.Log.Warnw("Fonnte not configured, skipping order created notification", "order_number", order.OrderNumber)
		} else {
			log.Printf("[WhatsApp] Fonnte not configured, skipping order created notification for %s", order.OrderNumber)
		}
		return nil
	}

	// Build message
	message := fmt.Sprintf(
		"🛍️ *Pesanan Baru!*\n\n"+
			"Nomor Pesanan: *%s*\n"+
			"Total: *Rp %s*\n\n"+
			"Silakan selesaikan pembayaran Anda.\n\n"+
			"Terima kasih telah berbelanja di Karima Store! 🙏",
		order.OrderNumber,
		formatCurrency(order.TotalAmount),
	)

	// Get customer phone from order
	customerPhone := order.ShippingPhone
	if customerPhone == "" {
		if logger.Log != nil {
			logger.Log.Warnw("No phone number for order", "order_number", order.OrderNumber)
		} else {
			log.Printf("[WhatsApp] No phone number for order %s", order.OrderNumber)
		}
		return nil
	}

	// Send asynchronously using goroutine
	go s.sendWhatsAppAsync(customerPhone, message, nil)

	if logger.Log != nil {
		logger.Log.Infow("Order created notification queued", "order_number", order.OrderNumber)
	} else {
		log.Printf("[WhatsApp] Order created notification queued for order %s", order.OrderNumber)
	}
	return nil
}

// SendPaymentSuccessNotification sends notification when payment is successful (ASYNC)
func (s *notificationService) SendPaymentSuccessNotification(order *models.Order) error {
	if s.fonnteClient == nil {
		if logger.Log != nil {
			logger.Log.Warnw("Fonnte not configured, skipping payment success notification", "order_number", order.OrderNumber)
		} else {
			log.Printf("[WhatsApp] Fonnte not configured, skipping payment success notification for %s", order.OrderNumber)
		}
		return nil
	}

	// Build message
	message := fmt.Sprintf(
		"✅ *Pembayaran Berhasil!*\n\n"+
			"Nomor Pesanan: *%s*\n"+
			"Total: *Rp %s*\n\n"+
			"Pesanan Anda sedang diproses dan akan segera dikirim.\n\n"+
			"Terima kasih! 🙏",
		order.OrderNumber,
		formatCurrency(order.TotalAmount),
	)

	// Get customer phone
	customerPhone := order.ShippingPhone
	if customerPhone == "" {
		if logger.Log != nil {
			logger.Log.Warnw("No phone number for order", "order_number", order.OrderNumber)
		} else {
			log.Printf("[WhatsApp] No phone number for order %s", order.OrderNumber)
		}
		return nil
	}

	// Send asynchronously using goroutine
	go s.sendWhatsAppAsync(customerPhone, message, nil)

	if logger.Log != nil {
		logger.Log.Infow("Payment success notification queued", "order_number", order.OrderNumber)
	} else {
		log.Printf("[WhatsApp] Payment success notification queued for order %s", order.OrderNumber)
	}
	return nil
}

// SendShippingNotification sends notification when order is shipped (ASYNC)
func (s *notificationService) SendShippingNotification(order *models.Order, trackingNumber string) error {
	if s.fonnteClient == nil {
		return nil
	}

	message := fmt.Sprintf(
		"📦 *Pesanan Dikirim!*\n\n"+
			"Nomor Pesanan: *%s*\n"+
			"Kurir: *%s*\n"+
			"No. Resi: *%s*\n\n"+
			"Lacak pesanan Anda untuk melihat status pengiriman.\n\n"+
			"Terima kasih! 🙏",
		order.OrderNumber,
		order.ShippingProvider,
		trackingNumber,
	)

	customerPhone := order.ShippingPhone
	if customerPhone == "" {
		return nil
	}

	go s.sendWhatsAppAsync(customerPhone, message, nil)
	return nil
}

// GetWhatsAppStatus checks WhatsApp service status
func (s *notificationService) GetWhatsAppStatus() (string, error) {
	if s.fonnteClient == nil {
		return "not_configured", nil
	}

	resp, err := s.fonnteClient.GetDeviceStatus()
	if err != nil {
		return "error", err
	}

	if resp.Status {
		return "connected", nil
	}
	return "disconnected", nil
}

// SendTestWhatsAppMessage sends a test message to verify WhatsApp integration
func (s *notificationService) SendTestWhatsAppMessage(phoneNumber string, message string) error {
	if s.fonnteClient == nil {
		return fmt.Errorf("fonnte client not configured")
	}

	formattedPhone := formatPhoneNumber(phoneNumber)
	resp, err := s.fonnteClient.SendMessage(formattedPhone, message)
	if err != nil {
		return err
	}

	if !resp.Status {
		return fmt.Errorf("test message failed: %s", resp.Detail)
	}

	if logger.Log != nil {
		logger.Log.Infow("Test WhatsApp message sent", "phone", formattedPhone)
	} else {
		log.Printf("[WhatsApp] Test message sent to %s", formattedPhone)
	}
	return nil
}

// ProcessWhatsAppWebhook handles WhatsApp webhook events
func (s *notificationService) ProcessWhatsAppWebhook(data map[string]interface{}) error {
	if logger.Log != nil {
		logger.Log.Infow("WhatsApp webhook received", "data", data)
	} else {
		log.Printf("[WhatsApp Webhook] Received: %v", data)
	}
	// Process webhook data as needed
	return nil
}

// GetWhatsAppWebhookURL returns the webhook URL for WhatsApp
func (s *notificationService) GetWhatsAppWebhookURL() string {
	return fmt.Sprintf("https://api.karimastore.com/api/v1/whatsapp/webhook")
}

// Helper functions

// formatPhoneNumber formats phone number to 628xxx format
func formatPhoneNumber(phone string) string {
	// Remove spaces and dashes
	phone = removeNonDigits(phone)

	// Convert 08xxx to 628xxx
	if len(phone) > 0 && phone[0] == '0' {
		phone = "62" + phone[1:]
	}

	// Ensure starts with 62
	if len(phone) > 2 && phone[0:2] != "62" {
		phone = "62" + phone
	}

	return phone
}

// removeNonDigits removes non-digit characters from string
func removeNonDigits(s string) string {
	result := ""
	for _, r := range s {
		if r >= '0' && r <= '9' {
			result += string(r)
		}
	}
	return result
}

// formatCurrency formats number to Indonesian currency format
func formatCurrency(amount float64) string {
	// Simple formatting without external lib
	return fmt.Sprintf("%.0f", amount)
}
