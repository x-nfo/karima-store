package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/karima-store/internal/database"
	"github.com/karima-store/internal/services"
)

type PricingHandler struct {
	pricingService *services.PricingService
	redisClient *database.Redis
}

func NewPricingHandler(pricingService *services.PricingService, redis *database.Redis) *PricingHandler {
	return &PricingHandler{
		pricingService: pricingService,
		redisClient: redis,
	}
}

// CalculatePrice calculates the price for a product
// @Summary Calculate product price
// @Description Calculate the final price based on customer type, quantity, and active flash sales
// @Tags pricing
// @Accept json
// @Produce json
// @Param request body services.PriceCalculationRequest true "Price calculation request"
// @Success 200 {object} map[string]interface{} "Price calculation response"
// @Failure 400 {object} map[string]interface{} "Bad request"
// @Failure 404 {object} map[string]interface{} "Product not found"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /api/v1/pricing/calculate [post]
func (h *PricingHandler) CalculatePrice(c *fiber.Ctx) error {
	var req services.PriceCalculationRequest

	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "error",
			"message": "Invalid request body",
			"code":    400,
		})
	}

	// Validate customer type
	if req.CustomerType != services.CustomerRetail && req.CustomerType != services.CustomerReseller {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "error",
			"message": "Invalid customer type. Must be 'retail' or 'reseller'",
			"code":    400,
		})
	}

	response, err := h.pricingService.CalculatePrice(req)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"status":  "error",
			"message": err.Error(),
			"code":    404,
		})
	}

	return c.JSON(fiber.Map{
		"status": "success",
		"data":   response,
	})
}

// CalculatePriceByParams calculates price using query parameters
// @Summary Calculate product price by parameters
// @Description Calculate the final price using query parameters instead of JSON body
// @Tags pricing
// @Accept json
// @Produce json
// @Param product_id query int true "Product ID"
// @Param variant_id query int false "Variant ID"
// @Param quantity query int true "Quantity"
// @Param customer_type query string true "Customer type (retail or reseller)"
// @Success 200 {object} map[string]interface{} "Price calculation response"
// @Failure 400 {object} map[string]interface{} "Bad request"
// @Failure 404 {object} map[string]interface{} "Product not found"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /api/v1/pricing/calculate [get]
func (h *PricingHandler) CalculatePriceByParams(c *fiber.Ctx) error {
	// Parse product_id
	productIDStr := c.Query("product_id")
	if productIDStr == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "error",
			"message": "product_id is required",
			"code":    400,
		})
	}

	productID, err := strconv.ParseUint(productIDStr, 10, 32)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "error",
			"message": "Invalid product_id",
			"code":    400,
		})
	}

	// Parse variant_id (optional)
	var variantID *uint
	variantIDStr := c.Query("variant_id")
	if variantIDStr != "" {
		vid, err := strconv.ParseUint(variantIDStr, 10, 32)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"status":  "error",
				"message": "Invalid variant_id",
				"code":    400,
			})
		}
		vidUint := uint(vid)
		variantID = &vidUint
	}

	// Parse quantity
	quantityStr := c.Query("quantity")
	if quantityStr == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "error",
			"message": "quantity is required",
			"code":    400,
		})
	}

	quantity, err := strconv.Atoi(quantityStr)
	if err != nil || quantity <= 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "error",
			"message": "Invalid quantity. Must be a positive integer",
			"code":    400,
		})
	}

	// Parse customer_type
	customerTypeStr := c.Query("customer_type")
	if customerTypeStr == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "error",
			"message": "customer_type is required",
			"code":    400,
		})
	}

	customerType := services.CustomerType(customerTypeStr)
	if customerType != services.CustomerRetail && customerType != services.CustomerReseller {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "error",
			"message": "Invalid customer_type. Must be 'retail' or 'reseller'",
			"code":    400,
		})
	}

	req := services.PriceCalculationRequest{
		ProductID:    uint(productID),
		VariantID:    variantID,
		Quantity:     quantity,
		CustomerType: customerType,
	}

	response, err := h.pricingService.CalculatePrice(req)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"status":  "error",
			"message": err.Error(),
			"code":    404,
		})
	}

	return c.JSON(fiber.Map{
		"status": "success",
		"data":   response,
	})
}

// CalculateShippingCost calculates shipping cost
// @Summary Calculate shipping cost
// @Description Calculate shipping cost based on items, weight, and destination
// @Tags pricing
// @Accept json
// @Produce json
// @Param request body services.ShippingCalculationRequest true "Shipping calculation request"
// @Success 200 {object} map[string]interface{} "Shipping calculation response"
// @Failure 400 {object} map[string]interface{} "Bad request"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /api/v1/pricing/shipping [post]
func (h *PricingHandler) CalculateShippingCost(c *fiber.Ctx) error {
	var req services.ShippingCalculationRequest

	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "error",
			"message": "Invalid request body",
			"code":    400,
		})
	}

	// Validate shipping type
	validShippingTypes := []string{"jne", "tiki", "pos", "sicepat"}
	isValid := false
	for _, validType := range validShippingTypes {
		if req.ShippingType == validType {
			isValid = true
			break
		}
	}

	if !isValid {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "error",
			"message": "Invalid shipping type. Valid types: jne, tiki, pos, sicepat",
			"code":    400,
		})
	}

	response, err := h.pricingService.CalculateShippingCost(req)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "error",
			"message": err.Error(),
			"code":    400,
		})
	}

	return c.JSON(fiber.Map{
		"status": "success",
		"data":   response,
	})
}

// CalculateOrderSummary calculates complete order summary
// @Summary Calculate order summary
// @Description Calculate complete order summary including pricing and shipping
// @Tags pricing
// @Accept json
// @Produce json
// @Param request body OrderSummaryRequest true "Order summary request"
// @Success 200 {object} map[string]interface{} "Order summary response"
// @Failure 400 {object} map[string]interface{} "Bad request"
// @Failure 404 {object} map[string]interface{} "Product not found"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /api/v1/pricing/order-summary [post]
func (h *PricingHandler) CalculateOrderSummary(c *fiber.Ctx) error {
	var req OrderSummaryRequest

	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "error",
			"message": "Invalid request body",
			"code":    400,
		})
	}

	// Validate customer type
	if req.CustomerType != services.CustomerRetail && req.CustomerType != services.CustomerReseller {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "error",
			"message": "Invalid customer type. Must be 'retail' or 'reseller'",
			"code":    400,
		})
	}

	// Validate shipping type
	validShippingTypes := []string{"jne", "tiki", "pos", "sicepat"}
	isValid := false
	for _, validType := range validShippingTypes {
		if req.Shipping.ShippingType == validType {
			isValid = true
			break
		}
	}

	if !isValid {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "error",
			"message": "Invalid shipping type. Valid types: jne, tiki, pos, sicepat",
			"code":    400,
		})
	}

	response, err := h.pricingService.CalculateOrderSummary(req.Items, req.Shipping, req.CustomerType)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"status":  "error",
			"message": err.Error(),
			"code":    404,
		})
	}

	return c.JSON(fiber.Map{
		"status": "success",
		"data":   response,
	})
}

// OrderSummaryRequest represents the request body for order summary calculation
type OrderSummaryRequest struct {
	Items        []services.PriceCalculationRequest `json:"items"`
	Shipping     services.ShippingCalculationRequest `json:"shipping"`
	CustomerType services.CustomerType               `json:"customer_type"`
	CouponCode    string                           `json:"coupon_code,omitempty"`
	UserID       uint                             `json:"user_id,omitempty"`
}

// ValidateCoupon validates a coupon code
// @Summary Validate coupon code
// @Description Validate if a coupon code is valid and applicable
// @Tags pricing
// @Accept json
// @Produce json
// @Param request body CouponValidationRequest true "Coupon validation request"
// @Success 200 {object} map[string]interface{} "Coupon validation response"
// @Failure 400 {object} map[string]interface{} "Bad request"
// @Failure 404 {object} map[string]interface{} "Invalid coupon"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /api/v1/pricing/coupons/validate [post]
func (h *PricingHandler) ValidateCoupon(c *fiber.Ctx) error {
	var req CouponValidationRequest

	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "error",
			"message": "Invalid request body",
			"code":    400,
		})
	}

	couponReq := services.CouponCalculationRequest{
		Code:          req.Code,
		UserID:        req.UserID,
		PurchaseAmount: req.PurchaseAmount,
		CustomerType:  req.CustomerType,
	}

	discount, couponName, err := h.pricingService.CalculateCouponDiscount(couponReq)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"status":  "error",
			"message": err.Error(),
			"code":    404,
		})
	}

	return c.JSON(fiber.Map{
		"status": "success",
		"data": fiber.Map{
			"valid":          true,
			"coupon_name":     couponName,
			"discount_amount": discount,
		},
	})
}

// CouponValidationRequest represents the request body for coupon validation
type CouponValidationRequest struct {
	Code           string                 `json:"code"`
	UserID         uint                   `json:"user_id"`
	PurchaseAmount float64                `json:"purchase_amount"`
	CustomerType  services.CustomerType    `json:"customer_type"`
}

// GetPricingInfo returns pricing information for a product
// @Summary Get pricing information
// @Description Get pricing information including available discounts and tiering
// @Tags pricing
// @Accept json
// @Produce json
// @Param product_id path int true "Product ID"
// @Success 200 {object} map[string]interface{} "Pricing information"
// @Failure 400 {object} map[string]interface{} "Bad request"
// @Failure 404 {object} map[string]interface{} "Product not found"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /api/v1/pricing/products/{product_id} [get]
func (h *PricingHandler) GetPricingInfo(c *fiber.Ctx) error {
	productIDStr := c.Params("product_id")
	if productIDStr == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "error",
			"message": "product_id is required",
			"code":    400,
		})
	}

	productID, err := strconv.ParseUint(productIDStr, 10, 32)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "error",
			"message": "Invalid product_id",
			"code":    400,
		})
	}

	// Create cache key
	cacheKey := fmt.Sprintf("pricing:%d", productID)

	// Check cache first
	val, err := h.redisClient.Get(c.Context(), cacheKey)
	if err == nil {
		// Cache hit - parse and return
		var cachedData map[string]interface{}
		if err := json.Unmarshal([]byte(val), &cachedData); err == nil {
			return c.JSON(fiber.Map{
				"status": "success",
				"data":   cachedData,
			})
		}
	}

	// Cache miss - calculate pricing
	// Get pricing info for both retail and reseller
	retailReq := services.PriceCalculationRequest{
		ProductID:    uint(productID),
		Quantity:     1,
		CustomerType: services.CustomerRetail,
	}

	resellerReq := services.PriceCalculationRequest{
		ProductID:    uint(productID),
		Quantity:     1,
		CustomerType: services.CustomerReseller,
	}

	retailPrice, err := h.pricingService.CalculatePrice(retailReq)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"status":  "error",
			"message": err.Error(),
			"code":    404,
		})
	}

	resellerPrice, err := h.pricingService.CalculatePrice(resellerReq)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"status":  "error",
			"message": err.Error(),
			"code":    404,
		})
	}

	// Create response
	response := fiber.Map{
		"retail":   retailPrice,
		"reseller": resellerPrice,
	}

	// Store in cache with 1 hour expiration
	if err := h.redisClient.Set(c.Context(), cacheKey, response, 1*time.Hour); err != nil {
		log.Printf("Failed to cache pricing data: %v", err)
	}

	return c.JSON(fiber.Map{
		"status": "success",
		"data":   response,
	})
}
