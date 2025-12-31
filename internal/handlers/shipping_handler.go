package handlers

import (
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/karima-store/internal/models"
	"github.com/karima-store/internal/services"
)

// ShippingHandler handles shipping-related HTTP requests
type ShippingHandler struct {
	shippingService services.ShippingService
}

// NewShippingHandler creates a new shipping handler
func NewShippingHandler(shippingService services.ShippingService) *ShippingHandler {
	return &ShippingHandler{
		shippingService: shippingService,
	}
}

// GetAllProvinces godoc
// @Summary Get all provinces
// @Description Get all provinces from RajaOngkir
// @Tags shipping
// @Accept json
// @Produce json
// @Success 200 {object} map[string]interface{}
// @Router /shipping/provinces [get]
func (h *ShippingHandler) GetAllProvinces(c *fiber.Ctx) error {
	provinces, err := h.shippingService.GetAllProvinces()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   "Failed to get provinces",
			"message": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"success": true,
		"data":    provinces,
	})
}

// GetProvinceByID godoc
// @Summary Get province by ID
// @Description Get a specific province by ID
// @Tags shipping
// @Accept json
// @Produce json
// @Param id path string true "Province ID"
// @Success 200 {object} map[string]interface{}
// @Router /shipping/provinces/{id} [get]
func (h *ShippingHandler) GetProvinceByID(c *fiber.Ctx) error {
	provinceID := c.Params("id")

	if provinceID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   "Province ID is required",
			"message": "Province ID parameter is missing",
		})
	}

	province, err := h.shippingService.GetProvinceByID(provinceID)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error":   "Province not found",
			"message": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"success": true,
		"data":    province,
	})
}

// GetAllCities godoc
// @Summary Get all cities
// @Description Get all cities from RajaOngkir
// @Tags shipping
// @Accept json
// @Produce json
// @Success 200 {object} map[string]interface{}
// @Router /shipping/cities [get]
func (h *ShippingHandler) GetAllCities(c *fiber.Ctx) error {
	cities, err := h.shippingService.GetAllCities()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   "Failed to get cities",
			"message": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"success": true,
		"data":    cities,
	})
}

// GetCitiesByProvince godoc
// @Summary Get cities by province
// @Description Get cities for a specific province
// @Tags shipping
// @Accept json
// @Produce json
// @Param province_id path string true "Province ID"
// @Success 200 {object} map[string]interface{}
// @Router /shipping/provinces/{province_id}/cities [get]
func (h *ShippingHandler) GetCitiesByProvince(c *fiber.Ctx) error {
	provinceID := c.Params("province_id")

	if provinceID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   "Province ID is required",
			"message": "Province ID parameter is missing",
		})
	}

	cities, err := h.shippingService.GetCitiesByProvince(provinceID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   "Failed to get cities",
			"message": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"success": true,
		"data":    cities,
	})
}

// GetCityByID godoc
// @Summary Get city by ID
// @Description Get a specific city by ID
// @Tags shipping
// @Accept json
// @Produce json
// @Param id path string true "City ID"
// @Success 200 {object} map[string]interface{}
// @Router /shipping/cities/{id} [get]
func (h *ShippingHandler) GetCityByID(c *fiber.Ctx) error {
	cityID := c.Params("id")

	if cityID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   "City ID is required",
			"message": "City ID parameter is missing",
		})
	}

	city, err := h.shippingService.GetCityByID(cityID)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error":   "City not found",
			"message": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"success": true,
		"data":    city,
	})
}

// GetSubdistricts godoc
// @Summary Get subdistricts by city
// @Description Get subdistricts for a specific city
// @Tags shipping
// @Accept json
// @Produce json
// @Param city_id path string true "City ID"
// @Success 200 {object} map[string]interface{}
// @Router /shipping/cities/{city_id}/subdistricts [get]
func (h *ShippingHandler) GetSubdistricts(c *fiber.Ctx) error {
	cityID := c.Params("city_id")

	if cityID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   "City ID is required",
			"message": "City ID parameter is missing",
		})
	}

	subdistricts, err := h.shippingService.GetSubdistricts(cityID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   "Failed to get subdistricts",
			"message": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"success": true,
		"data":    subdistricts,
	})
}

// GetSubdistrictByID godoc
// @Summary Get subdistrict by ID
// @Description Get a specific subdistrict by ID
// @Tags shipping
// @Accept json
// @Produce json
// @Param id path string true "Subdistrict ID"
// @Success 200 {object} map[string]interface{}
// @Router /shipping/subdistricts/{id} [get]
func (h *ShippingHandler) GetSubdistrictByID(c *fiber.Ctx) error {
	subdistrictID := c.Params("id")

	if subdistrictID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   "Subdistrict ID is required",
			"message": "Subdistrict ID parameter is missing",
		})
	}

	subdistrict, err := h.shippingService.GetSubdistrictByID(subdistrictID)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error":   "Subdistrict not found",
			"message": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"success": true,
		"data":    subdistrict,
	})
}

// RajaOngkirShippingCostRequest represents the request body for calculating shipping cost via RajaOngkir
type RajaOngkirShippingCostRequest struct {
	Origin          string `json:"origin"`
	OriginType      string `json:"originType"`
	Destination     string `json:"destination"`
	DestinationType string `json:"destinationType"`
	Weight          int    `json:"weight"`
	Courier         string `json:"courier"`
}

// CalculateShippingCost godoc
// @Summary Calculate shipping cost
// @Description Calculate shipping cost between origin and destination
// @Tags shipping
// @Accept json
// @Produce json
// @Param request body RajaOngkirShippingCostRequest true "Shipping cost request"
// @Success 200 {object} map[string]interface{}
// @Router /shipping/cost [post]
func (h *ShippingHandler) CalculateShippingCost(c *fiber.Ctx) error {
	var req RajaOngkirShippingCostRequest

	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   "Invalid request body",
			"message": err.Error(),
		})
	}

	// Validate required fields
	if req.Origin == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   "Origin is required",
			"message": "Origin field is missing",
		})
	}
	if req.Destination == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   "Destination is required",
			"message": "Destination field is missing",
		})
	}
	if req.Weight <= 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   "Invalid weight",
			"message": "Weight must be greater than 0",
		})
	}
	if req.Courier == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   "Courier is required",
			"message": "Courier field is missing",
		})
	}

	rajaongkirReq := models.RajaOngkirCostRequest{
		Origin:          req.Origin,
		OriginType:      req.OriginType,
		Destination:     req.Destination,
		DestinationType: req.DestinationType,
		Weight:          req.Weight,
		Courier:         req.Courier,
	}

	result, err := h.shippingService.CalculateShippingCost(rajaongkirReq)
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

// RajaOngkirMultipleShippingCostsRequest represents the request body for calculating shipping costs for multiple couriers
type RajaOngkirMultipleShippingCostsRequest struct {
	Origin          string   `json:"origin"`
	OriginType      string   `json:"originType"`
	Destination     string   `json:"destination"`
	DestinationType string   `json:"destinationType"`
	Weight          int      `json:"weight"`
	Couriers        []string `json:"couriers"`
}

// CalculateShippingCostsForMultipleCouriers godoc
// @Summary Calculate shipping costs for multiple couriers
// @Description Calculate shipping costs for multiple couriers between origin and destination
// @Tags shipping
// @Accept json
// @Produce json
// @Param request body RajaOngkirMultipleShippingCostsRequest true "Shipping costs request"
// @Success 200 {object} map[string]interface{}
// @Router /shipping/costs [post]
func (h *ShippingHandler) CalculateShippingCostsForMultipleCouriers(c *fiber.Ctx) error {
	var req RajaOngkirMultipleShippingCostsRequest

	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   "Invalid request body",
			"message": err.Error(),
		})
	}

	// Validate required fields
	if req.Origin == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   "Origin is required",
			"message": "Origin field is missing",
		})
	}
	if req.Destination == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   "Destination is required",
			"message": "Destination field is missing",
		})
	}
	if req.Weight <= 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   "Invalid weight",
			"message": "Weight must be greater than 0",
		})
	}
	if len(req.Couriers) == 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   "Couriers are required",
			"message": "At least one courier must be specified",
		})
	}

	result, err := h.shippingService.CalculateShippingCostsForMultipleCouriers(
		req.Origin,
		req.OriginType,
		req.Destination,
		req.DestinationType,
		req.Weight,
		req.Couriers,
	)

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   "Failed to calculate shipping costs",
			"message": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"success": true,
		"data":    result,
	})
}

// GetShippingOptions godoc
// @Summary Get all shipping options
// @Description Get all available shipping options for a given route
// @Tags shipping
// @Accept json
// @Produce json
// @Param origin query string true "Origin city ID"
// @Param destination query string true "Destination city ID"
// @Param weight query int true "Weight in grams"
// @Success 200 {object} map[string]interface{}
// @Router /shipping/options [get]
func (h *ShippingHandler) GetShippingOptions(c *fiber.Ctx) error {
	origin := c.Query("origin")
	destination := c.Query("destination")
	weightStr := c.Query("weight")

	if origin == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   "Origin is required",
			"message": "Origin query parameter is missing",
		})
	}

	if destination == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   "Destination is required",
			"message": "Destination query parameter is missing",
		})
	}

	if weightStr == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   "Weight is required",
			"message": "Weight query parameter is missing",
		})
	}

	weight, err := strconv.Atoi(weightStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   "Invalid weight",
			"message": "Weight must be a valid integer",
		})
	}

	result, err := h.shippingService.GetShippingOptions(origin, destination, weight)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   "Failed to get shipping options",
			"message": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"success": true,
		"data":    result,
	})
}

// GetCheapestShippingOption godoc
// @Summary Get cheapest shipping option
// @Description Get the cheapest shipping option for a given route
// @Tags shipping
// @Accept json
// @Produce json
// @Param origin query string true "Origin city ID"
// @Param destination query string true "Destination city ID"
// @Param weight query int true "Weight in grams"
// @Success 200 {object} map[string]interface{}
// @Router /shipping/options/cheapest [get]
func (h *ShippingHandler) GetCheapestShippingOption(c *fiber.Ctx) error {
	origin := c.Query("origin")
	destination := c.Query("destination")
	weightStr := c.Query("weight")

	if origin == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   "Origin is required",
			"message": "Origin query parameter is missing",
		})
	}

	if destination == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   "Destination is required",
			"message": "Destination query parameter is missing",
		})
	}

	if weightStr == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   "Weight is required",
			"message": "Weight query parameter is missing",
		})
	}

	weight, err := strconv.Atoi(weightStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   "Invalid weight",
			"message": "Weight must be a valid integer",
		})
	}

	result, err := h.shippingService.GetCheapestShippingOption(origin, destination, weight)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   "Failed to get cheapest shipping option",
			"message": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"success": true,
		"data":    result,
	})
}

// GetFastestShippingOption godoc
// @Summary Get fastest shipping option
// @Description Get the fastest shipping option for a given route
// @Tags shipping
// @Accept json
// @Produce json
// @Param origin query string true "Origin city ID"
// @Param destination query string true "Destination city ID"
// @Param weight query int true "Weight in grams"
// @Success 200 {object} map[string]interface{}
// @Router /shipping/options/fastest [get]
func (h *ShippingHandler) GetFastestShippingOption(c *fiber.Ctx) error {
	origin := c.Query("origin")
	destination := c.Query("destination")
	weightStr := c.Query("weight")

	if origin == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   "Origin is required",
			"message": "Origin query parameter is missing",
		})
	}

	if destination == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   "Destination is required",
			"message": "Destination query parameter is missing",
		})
	}

	if weightStr == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   "Weight is required",
			"message": "Weight query parameter is missing",
		})
	}

	weight, err := strconv.Atoi(weightStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   "Invalid weight",
			"message": "Weight must be a valid integer",
		})
	}

	result, err := h.shippingService.GetFastestShippingOption(origin, destination, weight)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   "Failed to get fastest shipping option",
			"message": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"success": true,
		"data":    result,
	})
}
