package services

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/karima-store/internal/config"
	"github.com/karima-store/internal/database"
	"github.com/karima-store/internal/models"
	"github.com/karima-store/internal/repository"
	"github.com/karima-store/internal/utils"
	"github.com/whatsapigo/whatsapigo"
)

// NotificationService handles all notification-related operations
type NotificationService struct {
	db *database.DB
	redis *database.Redis
}

// NewNotificationService creates a new notification service instance
func NewNotificationService(db *database.DB, redis *database.Redis) *NotificationService {
	return &NotificationService{
		db:    db,
		redis: redis,
	}
}

// SendWhatsAppMessage sends a message via WhatsApp Gateway API
func (s *NotificationService) SendWhatsAppMessage(c *models.Order, message string) error {
	// Get WhatsApp configuration from environment
	config := config.GetConfig()

	// Create WhatsApp client
	client := whatsapigo.NewClient(
		whatsapigo.WithAPIKey(config.WhatsAppAPIKey),
		whatsapigo.WithBaseURL(config.WhatsAppBaseURL),
	)

	// Prepare message
	messageData := map[string]interface{}{
		"to":       c.UserID,
		"message":  message,
		"reference": c.OrderNumber,
	}

	// Send message
	resp, err := client.SendMessage(context.Background(), messageData)
	if err != nil {
		log.Printf("Failed to send WhatsApp message: %v", err)
		return fmt.Errorf("failed to send WhatsApp message: %w", err)
	}

	// Save notification to database
	notification := &models.Notification{
		OrderID:     c.ID,
		OrderNumber: c.OrderNumber,
		RecipientID: c.UserID,
		Message:     message,
		Status:      "sent",
		Type:        "whatsapp",
		CreatedAt:   time.Now(),
	}

	if err := repository.NewNotificationRepository(s.db.DB()).Create(notification); err != nil {
		log.Printf("Failed to save notification record: %v", err)
		return fmt.Errorf("failed to save notification record: %w", err)
	}

	return nil
}

// SendOrderCreatedNotification sends notification when order is created
func (s *NotificationService) SendOrderCreatedNotification(order *models.Order) error {
	// Get user from database
	user, err := repository.NewUserRepository(s.db.DB()).GetByID(order.UserID)
	if err != nil {
		log.Printf("Failed to get user for order %d: %v", order.ID, err)
		return fmt.Errorf("failed to get user: %w", err)
	}

	// Prepare message
	message := fmt.Sprintf("Your order #%s has been successfully created. Thank you for shopping with us!", order.OrderNumber)

	// Send notification asynchronously
	go func() {
		if err := s.SendWhatsAppMessage(order, message); err != nil {
			log.Printf("Failed to send order created notification: %v", err)
		}
	}()

	return nil
}

// SendPaymentSuccessNotification sends notification when payment is successful
func (s *NotificationService) SendPaymentSuccessNotification(order *models.Order) error {
	// Get user from database
	user, err := repository.NewUserRepository(s.db.DB()).GetByID(order.UserID)
	if err != nil {
		log.Printf("Failed to get user for order %d: %v", order.ID, err)
		return fmt.Errorf("failed to get user: %w", err)
	}

	// Prepare message
	message := fmt.Sprintf("Payment for order #%s was successful. Your order is now confirmed.", order.OrderNumber)

	// Send notification asynchronously
	go func() {
		if err := s.SendWhatsAppMessage(order, message); err != nil {
			log.Printf("Failed to send payment success notification: %v", err)
		}
	}()

	return nil
}

// SendWhatsAppWebhook handles WhatsApp webhook events
func (s *NotificationService) SendWhatsAppWebhook(c *models.Order, event string) error {
	// Get user from database
	user, err := repository.NewUserRepository(s.db.DB()).GetByID(c.UserID)
	if err != nil {
		log.Printf("Failed to get user for order %d: %v", c.ID, err)
		return fmt.Errorf("failed to get user: %w", err)
	}

	// Prepare message based on event
	var message string
	switch event {
	case "payment_success":
		message = fmt.Sprintf("Payment for order #%s was successful. Your order is now confirmed.", c.OrderNumber)
	case "order_created":
		message = fmt.Sprintf("Your order #%s has been successfully created. Thank you for shopping with us!", c.OrderNumber)
	default:
		message = fmt.Sprintf("Event %s occurred for order #%s", event, c.OrderNumber)
	}

	// Send notification asynchronously
	go func() {
		if err := s.SendWhatsAppMessage(c, message); err != nil {
			log.Printf("Failed to send WhatsApp webhook notification: %v", err)
		}
	}()

	return nil
}

// ProcessWhatsAppWebhook processes incoming webhook events
func (s *NotificationService) ProcessWhatsAppWebhook(data map[string]interface{}) error {
	// Parse incoming webhook data
	var event string
	if val, ok := data["event"].(string); ok {
		event = val
	}

	// Get order from database
	orderNumber := data["reference"].(string)
	order, err := repository.NewOrderRepository(s.db.DB()).GetByOrderNumber(orderNumber)
	if err != nil {
		log.Printf("Failed to get order by number %s: %v", orderNumber, err)
		return fmt.Errorf("failed to get order: %w", err)
	}

	// Process based on event
	switch event {
	case "payment_success":
		return s.SendPaymentSuccessNotification(order)
	case "order_created":
		return s.SendOrderCreatedNotification(order)
	default:
		log.Printf("Unknown WhatsApp event: %s", event)
		return nil
	}
}

// GetWhatsAppConfig retrieves WhatsApp configuration
func (s *NotificationService) GetWhatsAppConfig() (config.Config, error) {
	return config.GetConfig()
}

// GetWhatsAppStatus checks WhatsApp service status
func (s *NotificationService) GetWhatsAppStatus() (string, error) {
	config, err := s.GetWhatsAppConfig()
	if err != nil {
		return "", fmt.Errorf("failed to get WhatsApp config: %w", err)
	}

	// Simple health check - check if API key is set
	if config.WhatsAppAPIKey == "" {
		return "unavailable", nil
	}

	return "available", nil
}

// SendTestWhatsAppMessage sends a test message to verify WhatsApp integration
func (s *NotificationService) SendTestWhatsAppMessage(phoneNumber string, message string) error {
	// Get WhatsApp configuration
	config, err := s.GetWhatsAppConfig()
	if err != nil {
		return fmt.Errorf("failed to get WhatsApp config: %w", err)
	}

	// Create WhatsApp client
	client := whatsapigo.NewClient(
		whatsapigo.WithAPIKey(config.WhatsAppAPIKey),
		whatsapigo.WithBaseURL(config.WhatsAppBaseURL),
	)

	// Prepare message
	messageData := map[string]interface{}{
		"to":       phoneNumber,
		"message":  message,
		"reference": "test",
	}

	// Send message
	resp, err := client.SendMessage(context.Background(), messageData)
	if err != nil {
		log.Printf("Failed to send test WhatsApp message: %v", err)
		return fmt.Errorf("failed to send test message: %w", err)
	}

	return nil
}

// GetWhatsAppWebhookURL returns the webhook URL for WhatsApp
func (s *NotificationService) GetWhatsAppWebhookURL() string {
	config, err := s.GetWhatsAppConfig()
	if err != nil {
		return ""
	}

	return fmt.Sprintf("%s/api/v1/whatsapp/webhook", config.APIBaseURL)
}

// GetWhatsAppStatus returns the current status of WhatsApp service
func (s *NotificationService) GetWhatsAppStatus() (string, error) {
	config, err := s.GetWhatsAppConfig()
	if err != nil {
		return "", fmt.Errorf("failed to get WhatsApp config: %w", err)
	}

	// Check if API key is set
	if config.WhatsAppAPIKey == "" {
		return "unavailable", nil
	}

	return "available", nil
}

// GetWhatsAppConfig returns WhatsApp configuration
func (s *NotificationService) GetWhatsAppConfig() (config.Config, error) {
	return config.GetConfig()
}

// GetWhatsAppWebhookURL returns the webhook URL for WhatsApp
func (s *NotificationService) GetWhatsAppWebhookURL() string {
	config, err := s.GetWhatsAppConfig()
	if err != nil {
		return ""
	}

	return fmt.Sprintf("%s/api/v1/whatsapp/webhook", config.APIBaseURL)
}