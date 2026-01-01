package handlers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/karima-store/internal/models"
	"github.com/karima-store/internal/services"
)

type CheckoutHandler struct {
	checkoutService *services.CheckoutService
}

func NewCheckoutHandler(checkoutService *services.CheckoutService) *CheckoutHandler {
	return &CheckoutHandler{
		checkoutService: checkoutService,
	}
}

// Checkout initiates the checkout process
// @Summary Create Order and Generate Payment Token
// @Description Creates a new order with items, calculates pricing including shipping and tax, then generates Midtrans Snap payment token. Returns order details and payment URL.
// @Tags payment
// @Accept json
// @Produce json
// @Param checkout body models.CheckoutRequest true "Checkout request with items and shipping information"
// @Success 200 {object} map[string]interface{} "Success response with order number, snap token, and payment URL"
// @Failure 400 {object} map[string]interface{} "Invalid request body"
// @Failure 500 {object} map[string]interface{} "Server error during order creation or payment token generation"
// @Router /api/v1/checkout [post]
func (h *CheckoutHandler) Checkout(c *fiber.Ctx) error {
	var req models.CheckoutRequest

	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "error",
			"message": "Invalid request body",
			"code":    400,
		})
	}

	// Process checkout
	response, err := h.checkoutService.Checkout(&req)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "error",
			"message": err.Error(),
			"code":    500,
		})
	}

	return c.JSON(fiber.Map{
		"status": "success",
		"data":   response,
	})
}

// PaymentWebhook handles Midtrans payment notifications
// @Summary Process Midtrans Payment Webhook
// @Description Receives and processes payment status notifications from Midtrans. Updates order status and manages stock based on payment result. This endpoint is called by Midtrans automatically.
// @Tags payment
// @Accept json
// @Produce json
// @Param notification body models.MidtransPaymentNotification true "Payment notification from Midtrans"
// @Success 200 {string} string "Notification processed successfully"
// @Failure 400 {object} map[string]interface{} "Invalid notification format"
// @Failure 500 {object} map[string]interface{} "Error processing notification"
// @Router /api/v1/payment/webhook [post]
func (h *CheckoutHandler) PaymentWebhook(c *fiber.Ctx) error {
	var notification models.MidtransPaymentNotification

	if err := c.BodyParser(&notification); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "error",
			"message": "Invalid notification body",
			"code":    400,
		})
	}

	// Process payment notification
	if err := h.checkoutService.ProcessPaymentNotification(&notification); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "error",
			"message": err.Error(),
			"code":    500,
		})
	}

	return c.SendStatus(fiber.StatusOK)
}
