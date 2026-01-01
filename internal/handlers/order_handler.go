package handlers

import (
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/karima-store/internal/services"
)

type OrderHandler struct {
	orderService services.OrderService
}

func NewOrderHandler(orderService services.OrderService) *OrderHandler {
	return &OrderHandler{orderService: orderService}
}

// GetOrders godoc
// @Summary Get user orders
// @Description Get list of orders for the authenticated user
// @Tags orders
// @Accept json
// @Produce json
// @Security Bearer
// @Param page query int false "Page number"
// @Param limit query int false "Items per page"
// @Success 200 {object} map[string]interface{}
// @Router /api/v1/orders [get]
func (h *OrderHandler) GetOrders(c *fiber.Ctx) error {
	// Assuming auth middleware sets "user_id"
	// We need to make sure the type casting is safe or handle nil
	uID := c.Locals("user_id")
	if uID == nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Unauthorized"})
	}

	// Handle both float64 (typical from unknown JSON) or uint/int if set directly
	var userID uint
	switch v := uID.(type) {
	case uint:
		userID = v
	case int:
		userID = uint(v)
	case float64:
		userID = uint(v)
	default:
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Invalid user ID type"})
	}

	page, _ := strconv.Atoi(c.Query("page", "1"))
	limit, _ := strconv.Atoi(c.Query("limit", "10"))
	offset := (page - 1) * limit

	orders, total, err := h.orderService.GetOrders(userID, limit, offset)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to fetch orders",
		})
	}

	return c.JSON(fiber.Map{
		"data":  orders,
		"total": total,
		"page":  page,
		"limit": limit,
	})
}

// GetOrder godoc
// @Summary Get order details
// @Description Get details of a specific order
// @Tags orders
// @Accept json
// @Produce json
// @Security Bearer
// @Param id path int true "Order ID"
// @Success 200 {object} map[string]interface{} // Using map instead of models.Order for simplicity in swagger for now
// @Router /api/v1/orders/{id} [get]
func (h *OrderHandler) GetOrder(c *fiber.Ctx) error {
	uID := c.Locals("user_id")
	if uID == nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Unauthorized"})
	}

	var userID uint
	switch v := uID.(type) {
	case uint:
		userID = v
	case int:
		userID = uint(v)
	case float64:
		userID = uint(v)
	}

	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid order ID"})
	}

	order, err := h.orderService.GetOrder(uint(id), userID)
	if err != nil {
		if err.Error() == "unauthorized" {
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"error": "Unauthorized"})
		}
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Order not found"})
	}

	return c.JSON(order)
}

// TrackOrder godoc
// @Summary Track order
// @Description Get order status by order number (public)
// @Tags orders
// @Accept json
// @Produce json
// @Param order_number query string true "Order Number"
// @Success 200 {object} map[string]interface{}
// @Router /api/v1/orders/track [get]
func (h *OrderHandler) TrackOrder(c *fiber.Ctx) error {
	orderNumber := c.Query("order_number")
	if orderNumber == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Order number is required"})
	}

	order, err := h.orderService.GetOrderByNumber(orderNumber)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Order not found"})
	}

	return c.JSON(order)
}
