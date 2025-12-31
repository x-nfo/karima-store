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
// @Summary Checkout
// @Description Create an order and generate Midtrans Snap token for payment
// @Tags checkout
// @Accept json
// @Produce json
// @Param checkout body models.CheckoutRequest true "Checkout request"
// @Success 200 {object} map[string]interface{} "Checkout response with Snap token"
// @Failure 400 {object} map[string]interface{} "Bad request"
// @Failure 500 {object} map[string]interface{} "Internal server error"
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
// @Summary Payment Webhook
// @Description Process payment notifications from Midtrans
// @Tags checkout
// @Accept json
// @Produce json
// @Param notification body models.MidtransPaymentNotification true "Payment notification"
// @Success 200 {string} string "OK"
// @Failure 400 {object} map[string]interface{} "Bad request"
// @Failure 500 {object} map[string]interface{} "Internal server error"
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
