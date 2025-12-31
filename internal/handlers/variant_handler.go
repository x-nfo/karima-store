package handlers

import (
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/karima-store/internal/models"
	"github.com/karima-store/internal/services"
)

type VariantHandler struct {
	variantService services.VariantService
}

func NewVariantHandler(variantService services.VariantService) *VariantHandler {
	return &VariantHandler{
		variantService: variantService,
	}
}

// CreateVariant creates a new product variant
// @Summary Create a new variant
// @Description Create a new product variant with the provided details
// @Tags variants
// @Accept json
// @Produce json
// @Param variant body models.ProductVariant true "Variant object"
// @Success 201 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /api/v1/variants [post]
func (h *VariantHandler) CreateVariant(c *fiber.Ctx) error {
	var variant models.ProductVariant

	if err := c.BodyParser(&variant); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "error",
			"message": "Invalid request body",
			"code":    400,
		})
	}

	if err := h.variantService.CreateVariant(&variant); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "error",
			"message": err.Error(),
			"code":    400,
		})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"status": "success",
		"data":   variant,
	})
}

// GetVariantByID retrieves a variant by ID
// @Summary Get variant by ID
// @Description Retrieve a single variant by its ID
// @Tags variants
// @Produce json
// @Param id path int true "Variant ID"
// @Success 200 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /api/v1/variants/{id} [get]
func (h *VariantHandler) GetVariantByID(c *fiber.Ctx) error {
	id, err := strconv.ParseUint(c.Params("id"), 10, 32)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "error",
			"message": "Invalid variant ID",
			"code":    400,
		})
	}

	variant, err := h.variantService.GetVariantByID(uint(id))
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"status":  "error",
			"message": err.Error(),
			"code":    404,
		})
	}

	return c.JSON(fiber.Map{
		"status": "success",
		"data":   variant,
	})
}

// GetVariantBySKU retrieves a variant by SKU
// @Summary Get variant by SKU
// @Description Retrieve a single variant by its SKU
// @Tags variants
// @Produce json
// @Param sku path string true "Variant SKU"
// @Success 200 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /api/v1/variants/sku/{sku} [get]
func (h *VariantHandler) GetVariantBySKU(c *fiber.Ctx) error {
	sku := c.Params("sku")

	variant, err := h.variantService.GetVariantBySKU(sku)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"status":  "error",
			"message": err.Error(),
			"code":    404,
		})
	}

	return c.JSON(fiber.Map{
		"status": "success",
		"data":   variant,
	})
}

// GetVariantsByProductID retrieves all variants for a product
// @Summary Get variants by product ID
// @Description Retrieve all variants for a specific product
// @Tags variants
// @Produce json
// @Param product_id path int true "Product ID"
// @Success 200 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /api/v1/products/{product_id}/variants [get]
func (h *VariantHandler) GetVariantsByProductID(c *fiber.Ctx) error {
	productID, err := strconv.ParseUint(c.Params("product_id"), 10, 32)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "error",
			"message": "Invalid product ID",
			"code":    400,
		})
	}

	variants, err := h.variantService.GetVariantsByProductID(uint(productID))
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"status":  "error",
			"message": err.Error(),
			"code":    404,
		})
	}

	return c.JSON(fiber.Map{
		"status": "success",
		"data":   variants,
	})
}

// UpdateVariant updates an existing variant
// @Summary Update variant
// @Description Update an existing variant by ID
// @Tags variants
// @Accept json
// @Produce json
// @Param id path int true "Variant ID"
// @Param variant body models.ProductVariant true "Variant object"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /api/v1/variants/{id} [put]
func (h *VariantHandler) UpdateVariant(c *fiber.Ctx) error {
	id, err := strconv.ParseUint(c.Params("id"), 10, 32)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "error",
			"message": "Invalid variant ID",
			"code":    400,
		})
	}

	var variant models.ProductVariant
	if err := c.BodyParser(&variant); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "error",
			"message": "Invalid request body",
			"code":    400,
		})
	}

	if err := h.variantService.UpdateVariant(uint(id), &variant); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "error",
			"message": err.Error(),
			"code":    400,
		})
	}

	// Get updated variant
	updatedVariant, err := h.variantService.GetVariantByID(uint(id))
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "error",
			"message": "Failed to retrieve updated variant",
			"code":    500,
		})
	}

	return c.JSON(fiber.Map{
		"status": "success",
		"data":   updatedVariant,
	})
}

// DeleteVariant deletes a variant
// @Summary Delete variant
// @Description Delete a variant by ID
// @Tags variants
// @Produce json
// @Param id path int true "Variant ID"
// @Success 200 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /api/v1/variants/{id} [delete]
func (h *VariantHandler) DeleteVariant(c *fiber.Ctx) error {
	id, err := strconv.ParseUint(c.Params("id"), 10, 32)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "error",
			"message": "Invalid variant ID",
			"code":    400,
		})
	}

	if err := h.variantService.DeleteVariant(uint(id)); err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"status":  "error",
			"message": err.Error(),
			"code":    404,
		})
	}

	return c.JSON(fiber.Map{
		"status":  "success",
		"message": "Variant deleted successfully",
	})
}

// UpdateVariantStock updates variant stock
// @Summary Update variant stock
// @Description Update the stock quantity for a variant
// @Tags variants
// @Accept json
// @Produce json
// @Param id path int true "Variant ID"
// @Param body body map[string]int true "Stock update object"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Router /api/v1/variants/{id}/stock [patch]
func (h *VariantHandler) UpdateVariantStock(c *fiber.Ctx) error {
	id, err := strconv.ParseUint(c.Params("id"), 10, 32)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "error",
			"message": "Invalid variant ID",
			"code":    400,
		})
	}

	var req struct {
		Quantity int `json:"quantity"`
	}

	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "error",
			"message": "Invalid request body",
			"code":    400,
		})
	}

	if err := h.variantService.UpdateVariantStock(uint(id), req.Quantity); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "error",
			"message": err.Error(),
			"code":    400,
		})
	}

	return c.JSON(fiber.Map{
		"status":  "success",
		"message": "Stock updated successfully",
	})
}
