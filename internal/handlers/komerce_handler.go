package handlers

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
	"github.com/karima-store/internal/models"
	"github.com/karima-store/internal/services"
)

// KomerceHandler handles Komerce API-related HTTP requests
type KomerceHandler struct {
	komerceService services.KomerceService
}

// NewKomerceHandler creates a new Komerce handler
func NewKomerceHandler(komerceService services.KomerceService) *KomerceHandler {
	return &KomerceHandler{
		komerceService: komerceService,
	}
}

// SearchDestination godoc
// @Summary Search destination
// @Description Search for destination by postal code, village, sub-district, or district
// @Tags komerce
// @Accept json
// @Produce json
// @Param keyword query string true "Search keyword"
// @Success 200 {object} map[string]interface{}
// @Router /komerce/destination/search [get]
func (h *KomerceHandler) SearchDestination(c *fiber.Ctx) error {
	keyword := c.Query("keyword")

	if keyword == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   "Keyword is required",
			"message": "Keyword query parameter is missing",
		})
	}

	destinations, err := h.komerceService.SearchDestination(keyword)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   "Failed to search destination",
			"message": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"success": true,
		"data":    destinations,
	})
}

// CalculateShippingCostRequest represents the request body for calculating shipping cost
type CalculateShippingCostRequest struct {
	ShipperDestinationID  string  `json:"shipper_destination_id"`
	ReceiverDestinationID string  `json:"receiver_destination_id"`
	Weight               float64 `json:"weight"`
	ItemValue            int     `json:"item_value"`
	COD                  string  `json:"cod"`
}

// CalculateShippingCost godoc
// @Summary Calculate shipping cost
// @Description Calculate shipping cost between shipper and receiver destination
// @Tags komerce
// @Accept json
// @Produce json
// @Param request body CalculateShippingCostRequest false "Shipping cost request (JSON body)"
// @Param shipper_destination_id query string false "Shipper destination ID"
// @Param receiver_destination_id query string false "Receiver destination ID"
// @Param weight query number false "Weight in kg"
// @Param item_value query number false "Item value"
// @Param cod query string false "COD option (yes/no)"
// @Success 200 {object} map[string]interface{}
// @Router /komerce/calculate [post]
func (h *KomerceHandler) CalculateShippingCost(c *fiber.Ctx) error {
	var req CalculateShippingCostRequest

	// Try to parse from JSON body first
	bodyErr := c.BodyParser(&req)

	// If body parsing fails or body is empty, try query parameters
	if bodyErr != nil || c.Body() == nil || len(c.Body()) == 0 {
		req.ShipperDestinationID = c.Query("shipper_destination_id")
		req.ReceiverDestinationID = c.Query("receiver_destination_id")

		// Parse weight from query parameter
		if weightStr := c.Query("weight"); weightStr != "" {
			var weight float64
			if _, err := fmt.Sscanf(weightStr, "%f", &weight); err == nil {
				req.Weight = weight
			}
		}

		// Parse item value from query parameter
		if itemValueStr := c.Query("item_value"); itemValueStr != "" {
			var itemValue int
			if _, err := fmt.Sscanf(itemValueStr, "%d", &itemValue); err == nil {
				req.ItemValue = itemValue
			}
		}

		req.COD = c.Query("cod", "no")
	}

	// Validate required fields
	if req.ShipperDestinationID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   "Shipper destination ID is required",
			"message": "shipper_destination_id field is missing",
		})
	}
	if req.ReceiverDestinationID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   "Receiver destination ID is required",
			"message": "receiver_destination_id field is missing",
		})
	}
	if req.Weight <= 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   "Invalid weight",
			"message": "Weight must be greater than 0",
		})
	}
	if req.ItemValue < 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   "Invalid item value",
			"message": "Item value must be greater than or equal to 0",
		})
	}

	// Set default COD value if not provided
	if req.COD == "" {
		req.COD = "no"
	}

	result, err := h.komerceService.CalculateShippingCost(
		req.ShipperDestinationID,
		req.ReceiverDestinationID,
		req.Weight,
		req.ItemValue,
		req.COD,
	)

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   "Failed to calculate shipping cost",
			"message": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"success": true,
		"data":    result,
	})
}

// CreateOrder godoc
// @Summary Create order
// @Description Create a new shipping order
// @Tags komerce
// @Accept json
// @Produce json
// @Param request body models.KomerceCreateOrderRequest true "Create order request"
// @Success 201 {object} map[string]interface{}
// @Router /komerce/orders [post]
func (h *KomerceHandler) CreateOrder(c *fiber.Ctx) error {
	var req models.KomerceCreateOrderRequest

	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   "Invalid request body",
			"message": err.Error(),
		})
	}

	result, err := h.komerceService.CreateOrder(req)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   "Failed to create order",
			"message": err.Error(),
		})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"success": true,
		"data":    result,
	})
}

// GetOrderDetail godoc
// @Summary Get order detail
// @Description Get order detail by order number
// @Tags komerce
// @Accept json
// @Produce json
// @Param order_no path string true "Order number"
// @Success 200 {object} map[string]interface{}
// @Router /komerce/orders/:order_no [get]
func (h *KomerceHandler) GetOrderDetail(c *fiber.Ctx) error {
	orderNo := c.Params("order_no")

	if orderNo == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   "Order number is required",
			"message": "Order number parameter is missing",
		})
	}

	result, err := h.komerceService.GetOrderDetail(orderNo)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   "Failed to get order detail",
			"message": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"success": true,
		"data":    result,
	})
}

// CancelOrderRequest represents the request body for canceling order
type CancelOrderRequest struct {
	OrderNo string `json:"order_no"`
}

// CancelOrder godoc
// @Summary Cancel order
// @Description Cancel an order by order number
// @Tags komerce
// @Accept json
// @Produce json
// @Param request body CancelOrderRequest true "Cancel order request"
// @Success 200 {object} map[string]interface{}
// @Router /komerce/orders/cancel [put]
func (h *KomerceHandler) CancelOrder(c *fiber.Ctx) error {
	var req CancelOrderRequest

	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   "Invalid request body",
			"message": err.Error(),
		})
	}

	if req.OrderNo == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   "Order number is required",
			"message": "order_no field is missing",
		})
	}

	err := h.komerceService.CancelOrder(req.OrderNo)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   "Failed to cancel order",
			"message": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"success": true,
		"message": "Order cancelled successfully",
	})
}

// RequestPickupRequest represents the request body for requesting pickup
type RequestPickupRequest struct {
	PickupVehicle string   `json:"pickup_vehicle"`
	PickupTime    string   `json:"pickup_time"`
	PickupDate    string   `json:"pickup_date"`
	Orders        []string `json:"orders"`
}

// RequestPickup godoc
// @Summary Request pickup
// @Description Request pickup for orders
// @Tags komerce
// @Accept json
// @Produce json
// @Param request body RequestPickupRequest true "Pickup request"
// @Success 201 {object} map[string]interface{}
// @Router /komerce/pickup [post]
func (h *KomerceHandler) RequestPickup(c *fiber.Ctx) error {
	var req RequestPickupRequest

	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   "Invalid request body",
			"message": err.Error(),
		})
	}

	// Validate required fields
	if req.PickupVehicle == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   "Pickup vehicle is required",
			"message": "pickup_vehicle field is missing",
		})
	}
	if req.PickupTime == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   "Pickup time is required",
			"message": "pickup_time field is missing",
		})
	}
	if req.PickupDate == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   "Pickup date is required",
			"message": "pickup_date field is missing",
		})
	}
	if len(req.Orders) == 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   "Orders are required",
			"message": "At least one order must be specified",
		})
	}

	err := h.komerceService.RequestPickup(req.PickupVehicle, req.PickupTime, req.PickupDate, req.Orders)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   "Failed to request pickup",
			"message": err.Error(),
		})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"success": true,
		"message": "Pickup requested successfully",
	})
}

// PrintLabel godoc
// @Summary Print label
// @Description Generate print label for order
// @Tags komerce
// @Accept json
// @Produce json
// @Param order_no query string true "Order number"
// @Param page query string false "Page number"
// @Success 200 {object} map[string]interface{}
// @Router /komerce/orders/print-label [post]
func (h *KomerceHandler) PrintLabel(c *fiber.Ctx) error {
	orderNo := c.Query("order_no")
	page := c.Query("page")
	if page == "" {
		page = "page_1"
	}

	if orderNo == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   "Order number is required",
			"message": "order_no query parameter is missing",
		})
	}

	path, err := h.komerceService.PrintLabel(orderNo, page)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   "Failed to print label",
			"message": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"success": true,
		"data": fiber.Map{
			"path": path,
			"url":  "https://api-sandbox.collaborator.komerce.id" + path,
		},
	})
}

// TrackOrder godoc
// @Summary Track order
// @Description Track order by shipping provider and airway bill
// @Tags komerce
// @Accept json
// @Produce json
// @Param shipping query string true "Shipping provider"
// @Param airway_bill query string true "Airway bill number"
// @Success 200 {object} map[string]interface{}
// @Router /komerce/orders/track [get]
func (h *KomerceHandler) TrackOrder(c *fiber.Ctx) error {
	shipping := c.Query("shipping")
	airwayBill := c.Query("airway_bill")

	if shipping == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   "Shipping provider is required",
			"message": "shipping query parameter is missing",
		})
	}
	if airwayBill == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   "Airway bill is required",
			"message": "airway_bill query parameter is missing",
		})
	}

	result, err := h.komerceService.TrackOrder(shipping, airwayBill)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   "Failed to track order",
			"message": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"success": true,
		"data":    result,
	})
}
