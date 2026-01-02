package handlers

import (
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/karima-store/internal/models"
	"github.com/karima-store/internal/repository"
	"github.com/karima-store/internal/services"
	"gorm.io/gorm"
)

// WhatsAppHandler handles WhatsApp API endpoints
type WhatsAppHandler struct {
	notificationService *services.NotificationService
}

// NewWhatsAppHandler creates a new WhatsApp handler instance
func NewWhatsAppHandler(notificationService *services.NotificationService) *WhatsAppHandler {
	return &WhatsAppHandler{
		notificationService: notificationService,
	}
}

// SendWhatsAppMessage sends a message via WhatsApp
// @Summary Send WhatsApp message (Admin only)
// @Description Send a message to a WhatsApp number. **Admin only**: Requires authentication with admin role.
// @Tags whatsapp
// @Accept json
// @Produce json
// @Security KratosSession []
// @Security KratosSessionCookie []
// @Param phoneNumber query string true "Phone number in E.164 format"
// @Param message body string true "Message content"
// @Success 200 {object} map[string]interface{} "Success response"
// @Failure 400 {object} map[string]interface{} "Bad request"
// @Failure 401 {object} map[string]interface{} "Unauthorized: No valid session or session expired"
// @Failure 403 {object} map[string]interface{} "Forbidden: Insufficient permissions (admin role required)"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /api/v1/whatsapp/send [post]
func (h *WhatsAppHandler) SendWhatsAppMessage(c *fiber.Ctx) error {
	var req struct {
		PhoneNumber string `json:"phone_number"`
		Message     string `json:"message"`
	}

	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "error",
			"message": "Invalid request body",
			"code":    400,
		})
	}

	// Validate phone number format
	if len(req.PhoneNumber) < 10 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "error",
			"message": "Invalid phone number format",
			"code":    400,
		})
	}

	// Send message
	err := h.notificationService.SendWhatsAppMessage(&models.Order{}, req.Message, req.PhoneNumber)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "error",
			"message": err.Error(),
			"code":    500,
		})
	}

	return c.JSON(fiber.Map{
		"status": "success",
		"data": fiber.Map{
			"message": "Message sent successfully",
		},
	})
}

// SendOrderCreatedNotification sends order created notification
// @Summary Send order created notification (Admin only)
// @Description Send notification when order is created. **Admin only**: Requires authentication with admin role.
// @Tags whatsapp
// @Accept json
// @Produce json
// @Security KratosSession []
// @Security KratosSessionCookie []
// @Param order_id path uint true "Order ID"
// @Success 200 {object} map[string]interface{} "Success response"
// @Failure 400 {object} map[string]interface{} "Bad request"
// @Failure 401 {object} map[string]interface{} "Unauthorized: No valid session or session expired"
// @Failure 403 {object} map[string]interface{} "Forbidden: Insufficient permissions (admin role required)"
// @Failure 404 {object} map[string]interface{} "Order not found"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /api/v1/whatsapp/order-created/{order_id} [get]
func (h *WhatsAppHandler) SendOrderCreatedNotification(c *fiber.Ctx) error {
	orderIDStr := c.Params("order_id")
	if orderIDStr == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "error",
			"message": "order_id is required",
			"code":    400,
		})
	}

	orderID, err := strconv.ParseUint(orderIDStr, 10, 64)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "error",
			"message": "Invalid order_id",
			"code":    400,
		})
	}

	// Get order from database
	db, ok := h.notificationService.GetDB().(*gorm.DB)
	if !ok {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "error",
			"message": "Database connection error",
			"code":    500,
		})
	}

	order, err := repository.NewOrderRepository(db).GetByID(uint(orderID))
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"status":  "error",
			"message": "Order not found",
			"code":    404,
		})
	}

	// Send notification
	err = h.notificationService.SendOrderCreatedNotification(order)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "error",
			"message": err.Error(),
			"code":    500,
		})
	}

	return c.JSON(fiber.Map{
		"status": "success",
		"data": fiber.Map{
			"message": "Order created notification sent successfully",
		},
	})
}

// SendPaymentSuccessNotification sends payment success notification
// @Summary Send payment success notification (Admin only)
// @Description Send notification when payment is successful. **Admin only**: Requires authentication with admin role.
// @Tags whatsapp
// @Accept json
// @Produce json
// @Security KratosSession []
// @Security KratosSessionCookie []
// @Param order_id path uint true "Order ID"
// @Success 200 {object} map[string]interface{} "Success response"
// @Failure 400 {object} map[string]interface{} "Bad request"
// @Failure 401 {object} map[string]interface{} "Unauthorized: No valid session or session expired"
// @Failure 403 {object} map[string]interface{} "Forbidden: Insufficient permissions (admin role required)"
// @Failure 404 {object} map[string]interface{} "Order not found"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /api/v1/whatsapp/payment-success/{order_id} [get]
func (h *WhatsAppHandler) SendPaymentSuccessNotification(c *fiber.Ctx) error {
	orderIDStr := c.Params("order_id")
	if orderIDStr == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "error",
			"message": "order_id is required",
			"code":    400,
		})
	}

	orderID, err := strconv.ParseUint(orderIDStr, 10, 64)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "error",
			"message": "Invalid order_id",
			"code":    400,
		})
	}

	// Get order from database
	db, ok := h.notificationService.GetDB().(*gorm.DB)
	if !ok {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "error",
			"message": "Database connection error",
			"code":    500,
		})
	}

	order, err := repository.NewOrderRepository(db).GetByID(uint(orderID))
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"status":  "error",
			"message": "Order not found",
			"code":    404,
		})
	}

	// Send notification
	err = h.notificationService.SendPaymentSuccessNotification(order)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "error",
			"message": err.Error(),
			"code":    500,
		})
	}

	return c.JSON(fiber.Map{
		"status": "success",
		"data": fiber.Map{
			"message": "Payment success notification sent successfully",
		},
	})
}

// ProcessWhatsAppWebhook handles WhatsApp webhook events
// @Summary Process WhatsApp webhook
// @Description Handle incoming webhook events from WhatsApp (public endpoint, no authentication required)
// @Tags whatsapp
// @Accept json
// @Produce json
// @Param body body map[string]interface{} true "Webhook data"
// @Success 200 {object} map[string]interface{} "Success response"
// @Failure 400 {object} map[string]interface{} "Bad request"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /api/v1/whatsapp/webhook [post]
func (h *WhatsAppHandler) ProcessWhatsAppWebhook(c *fiber.Ctx) error {
	var data map[string]interface{}
	if err := c.BodyParser(&data); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "error",
			"message": "Invalid request body",
			"code":    400,
		})
	}

	// Process webhook
	err := h.notificationService.ProcessWhatsAppWebhook(data)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "error",
			"message": err.Error(),
			"code":    500,
		})
	}

	return c.JSON(fiber.Map{
		"status": "success",
		"data": fiber.Map{
			"message": "Webhook processed successfully",
		},
	})
}

// GetWhatsAppStatus returns WhatsApp service status
// @Summary Get WhatsApp status
// @Description Get current status of WhatsApp service (public endpoint, no authentication required)
// @Tags whatsapp
// @Accept json
// @Produce json
// @Success 200 {object} map[string]interface{} "Status response"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /api/v1/whatsapp/status [get]
func (h *WhatsAppHandler) GetWhatsAppStatus(c *fiber.Ctx) error {
	status, err := h.notificationService.GetWhatsAppStatus()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "error",
			"message": err.Error(),
			"code":    500,
		})
	}

	return c.JSON(fiber.Map{
		"status": "success",
		"data": fiber.Map{
			"status": status,
		},
	})
}

// SendTestWhatsAppMessage sends a test message to verify WhatsApp integration
// @Summary Send test WhatsApp message (Admin only)
// @Description Send a test message to verify WhatsApp integration. **Admin only**: Requires authentication with admin role.
// @Tags whatsapp
// @Accept json
// @Produce json
// @Security KratosSession []
// @Security KratosSessionCookie []
// @Param phoneNumber query string true "Phone number in E.164 format"
// @Param message body string true "Message content"
// @Success 200 {object} map[string]interface{} "Success response"
// @Failure 400 {object} map[string]interface{} "Bad request"
// @Failure 401 {object} map[string]interface{} "Unauthorized: No valid session or session expired"
// @Failure 403 {object} map[string]interface{} "Forbidden: Insufficient permissions (admin role required)"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /api/v1/whatsapp/test [post]
func (h *WhatsAppHandler) SendTestWhatsAppMessage(c *fiber.Ctx) error {
	phoneNumber := c.Query("phone_number")
	message := c.Query("message")

	if phoneNumber == "" || message == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "error",
			"message": "phone_number and message are required",
			"code":    400,
		})
	}

	// Validate phone number format
	if len(phoneNumber) < 10 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "error",
			"message": "Invalid phone number format",
			"code":    400,
		})
	}

	// Send test message
	err := h.notificationService.SendTestWhatsAppMessage(phoneNumber, message)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "error",
			"message": err.Error(),
			"code":    500,
		})
	}

	return c.JSON(fiber.Map{
		"status": "success",
		"data": fiber.Map{
			"message": "Test message sent successfully",
		},
	})
}

// GetWhatsAppWebhookURL returns the webhook URL for WhatsApp
// @Summary Get WhatsApp webhook URL
// @Description Get the webhook URL for WhatsApp
// @Tags whatsapp
// @Accept json
// @Produce json
// @Success 200 {object} map[string]interface{} "Webhook URL response"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /api/v1/whatsapp/webhook-url [get]
func (h *WhatsAppHandler) GetWhatsAppWebhookURL(c *fiber.Ctx) error {
	url := h.notificationService.GetWhatsAppWebhookURL()
	if url == "" {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "error",
			"message": "Failed to get webhook URL",
			"code":    500,
		})
	}

	return c.JSON(fiber.Map{
		"status": "success",
		"data": fiber.Map{
			"webhook_url": url,
		},
	})
}