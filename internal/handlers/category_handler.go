package handlers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/karima-store/internal/models"
	"github.com/karima-store/internal/repository"
	"github.com/karima-store/internal/services"
)

type CategoryHandler struct {
	categoryService services.CategoryService
}

func NewCategoryHandler(categoryService services.CategoryService) *CategoryHandler {
	return &CategoryHandler{
		categoryService: categoryService,
	}
}

// GetAllCategories retrieves all available product categories
// @Summary Get all categories
// @Description Retrieve all available product categories
// @Tags categories
// @Produce json
// @Success 200 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /api/v1/categories [get]
func (h *CategoryHandler) GetAllCategories(c *fiber.Ctx) error {
	categories := h.categoryService.GetAllCategories()

	// Convert to response format with display names
	categoryList := make([]map[string]interface{}, 0, len(categories))
	for _, category := range categories {
		categoryList = append(categoryList, map[string]interface{}{
			"slug": string(category),
			"name": h.categoryService.GetCategoryName(category),
		})
	}

	return c.JSON(fiber.Map{
		"status": "success",
		"data":   categoryList,
	})
}

// GetCategoryStats retrieves statistics for each category
// @Summary Get category statistics
// @Description Retrieve product count statistics for each category
// @Tags categories
// @Produce json
// @Success 200 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /api/v1/categories/stats [get]
func (h *CategoryHandler) GetCategoryStats(c *fiber.Ctx) error {
	stats, err := h.categoryService.GetCategoryStats()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "error",
			"message": "Failed to retrieve category statistics",
			"code":    500,
		})
	}

	// Enhance stats with category names
	enhancedStats := make([]map[string]interface{}, 0, len(stats))
	for _, stat := range stats {
		enhancedStats = append(enhancedStats, map[string]interface{}{
			"category":      stat.Category,
			"category_name": h.categoryService.GetCategoryName(stat.Category),
			"product_count": stat.ProductCount,
		})
	}

	return c.JSON(fiber.Map{
		"status": "success",
		"data":   enhancedStats,
	})
}

// GetCategoryProducts retrieves products for a specific category
// @Summary Get products by category
// @Description Retrieve products filtered by category
// @Tags categories
// @Produce json
// @Param category path string true "Product category"
// @Param page query int false "Page number" default(1)
// @Param limit query int false "Items per page" default(20)
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /api/v1/categories/{category}/products [get]
func (h *CategoryHandler) GetCategoryProducts(c *fiber.Ctx) error {
	categoryStr := c.Params("category")
	category := models.ProductCategory(categoryStr)

	// Validate category
	if !h.categoryService.IsValidCategory(category) {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "error",
			"message": "Invalid category",
			"code":    400,
		})
	}

	// This endpoint would typically delegate to product handler
	// For now, return category info
	return c.JSON(fiber.Map{
		"status": "success",
		"data": fiber.Map{
			"category":      category,
			"category_name": h.categoryService.GetCategoryName(category),
			"message":       "Use /api/v1/products/category/{category} endpoint for products",
		},
	})
}

// CategoryResponse represents the category response structure
type CategoryResponse struct {
	Slug string `json:"slug"`
	Name string `json:"name"`
}

// CategoryStatsResponse represents the category statistics response
type CategoryStatsResponse struct {
	Category      models.ProductCategory `json:"category"`
	CategoryName  string                 `json:"category_name"`
	ProductCount  int64                  `json:"product_count"`
}

// Helper function to convert repository stats to response
func convertCategoryStats(stats []repository.CategoryStats, categoryService services.CategoryService) []CategoryStatsResponse {
	response := make([]CategoryStatsResponse, 0, len(stats))
	for _, stat := range stats {
		response = append(response, CategoryStatsResponse{
			Category:     stat.Category,
			CategoryName: categoryService.GetCategoryName(stat.Category),
			ProductCount: stat.ProductCount,
		})
	}
	return response
}
