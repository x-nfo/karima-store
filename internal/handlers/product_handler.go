package handlers

import (
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/karima-store/internal/models"
	"github.com/karima-store/internal/services"
	"github.com/karima-store/internal/utils"
)

type ProductHandler struct {
	productService services.ProductService
	mediaService   *services.MediaService
}

func NewProductHandler(productService services.ProductService, mediaService *services.MediaService) *ProductHandler {
	return &ProductHandler{
		productService: productService,
		mediaService:   mediaService,
	}
}

// CreateProduct creates a new product
// @Summary Create a new product
// @Description Create a new product with the provided details
// @Tags products
// @Accept json
// @Produce json
// @Param product body models.Product true "Product object"
// @Success 201 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /api/v1/products [post]
func (h *ProductHandler) CreateProduct(c *fiber.Ctx) error {
	var product models.Product

	if err := utils.ParseAndValidate(c, &product); err != nil {
		return err
	}

	if err := h.productService.CreateProduct(&product); err != nil {
		return utils.SendError(c, fiber.StatusBadRequest, err.Error(), nil)
	}

	return utils.SendCreated(c, product, "Product created successfully")
}

// GetProductByID retrieves a product by ID
// @Summary Get product by ID
// @Description Retrieve a single product by its ID
// @Tags products
// @Produce json
// @Param id path int true "Product ID"
// @Success 200 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /api/v1/products/{id} [get]
func (h *ProductHandler) GetProductByID(c *fiber.Ctx) error {
	id, err := strconv.ParseUint(c.Params("id"), 10, 32)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "error",
			"message": "Invalid product ID",
			"code":    400,
		})
	}

	product, err := h.productService.GetProductByID(uint(id))
	if err != nil {
		return utils.SendError(c, fiber.StatusNotFound, err.Error(), nil)
	}

	return utils.SendSuccess(c, product, "Product retrieved successfully")
}

// GetProductBySlug retrieves a product by slug
// @Summary Get product by slug
// @Description Retrieve a single product by its slug
// @Tags products
// @Produce json
// @Param slug path string true "Product slug"
// @Success 200 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /api/v1/products/slug/{slug} [get]
func (h *ProductHandler) GetProductBySlug(c *fiber.Ctx) error {
	slug := c.Params("slug")

	product, err := h.productService.GetProductBySlug(slug)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"status":  "error",
			"message": err.Error(),
			"code":    404,
		})
	}

	return c.JSON(fiber.Map{
		"status": "success",
		"data":   product,
	})
}

// GetProducts retrieves a list of products with optional filters
// @Summary Get products
// @Description Retrieve a list of products with pagination and optional filters
// @Tags products
// @Produce json
// @Param page query int false "Page number" default(1)
// @Param limit query int false "Items per page" default(20)
// @Param category query string false "Filter by category"
// @Param status query string false "Filter by status"
// @Param min_price query number false "Minimum price"
// @Param max_price query number false "Maximum price"
// @Param brand query string false "Filter by brand"
// @Success 200 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /api/v1/products [get]
func (h *ProductHandler) GetProducts(c *fiber.Ctx) error {
	// Parse pagination parameters
	page, _ := strconv.Atoi(c.Query("page", "1"))
	limit, _ := strconv.Atoi(c.Query("limit", "20"))
	offset := (page - 1) * limit

	// Build filters
	filters := make(map[string]interface{})

	if category := c.Query("category"); category != "" {
		filters["category"] = category
	}

	if status := c.Query("status"); status != "" {
		filters["status"] = status
	}

	if minPrice := c.Query("min_price"); minPrice != "" {
		if price, err := strconv.ParseFloat(minPrice, 64); err == nil {
			filters["min_price"] = price
		}
	}

	if maxPrice := c.Query("max_price"); maxPrice != "" {
		if price, err := strconv.ParseFloat(maxPrice, 64); err == nil {
			filters["max_price"] = price
		}
	}

	if brand := c.Query("brand"); brand != "" {
		filters["brand"] = brand
	}

	products, total, err := h.productService.GetProducts(limit, offset, filters)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "error",
			"message": "Failed to retrieve products",
			"code":    500,
		})
	}

	return c.JSON(fiber.Map{
		"status": "success",
		"data": fiber.Map{
			"products": products,
			"pagination": fiber.Map{
				"page":        page,
				"limit":       limit,
				"total_items": total,
				"total_pages": (total + int64(limit) - 1) / int64(limit),
			},
		},
	})
}

// UpdateProduct updates an existing product
// @Summary Update product
// @Description Update an existing product by ID
// @Tags products
// @Accept json
// @Produce json
// @Param id path int true "Product ID"
// @Param product body models.Product true "Product object"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /api/v1/products/{id} [put]
func (h *ProductHandler) UpdateProduct(c *fiber.Ctx) error {
	id, err := strconv.ParseUint(c.Params("id"), 10, 32)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "error",
			"message": "Invalid product ID",
			"code":    400,
		})
	}

	var product models.Product
	if err := utils.ParseAndValidate(c, &product); err != nil {
		return err
	}

	if err := h.productService.UpdateProduct(uint(id), &product); err != nil {
		return utils.SendError(c, fiber.StatusBadRequest, err.Error(), nil)
	}

	// Get updated product
	updatedProduct, err := h.productService.GetProductByID(uint(id))
	if err != nil {
		return utils.SendError(c, fiber.StatusInternalServerError, "Failed to retrieve updated product", nil)
	}

	return utils.SendSuccess(c, updatedProduct, "Product updated successfully")
}

// DeleteProduct deletes a product
// @Summary Delete product
// @Description Delete a product by ID
// @Tags products
// @Produce json
// @Param id path int true "Product ID"
// @Success 200 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /api/v1/products/{id} [delete]
func (h *ProductHandler) DeleteProduct(c *fiber.Ctx) error {
	id, err := strconv.ParseUint(c.Params("id"), 10, 32)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "error",
			"message": "Invalid product ID",
			"code":    400,
		})
	}

	if err := h.productService.DeleteProduct(uint(id)); err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"status":  "error",
			"message": err.Error(),
			"code":    404,
		})
	}

	return c.JSON(fiber.Map{
		"status":  "success",
		"message": "Product deleted successfully",
	})
}

// UpdateProductStock updates product stock
// @Summary Update product stock
// @Description Update the stock quantity for a product
// @Tags products
// @Accept json
// @Produce json
// @Param id path int true "Product ID"
// @Param body body map[string]int true "Stock update object"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Router /api/v1/products/{id}/stock [patch]
func (h *ProductHandler) UpdateProductStock(c *fiber.Ctx) error {
	id, err := strconv.ParseUint(c.Params("id"), 10, 32)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "error",
			"message": "Invalid product ID",
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

	if err := h.productService.UpdateProductStock(uint(id), req.Quantity); err != nil {
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

// SearchProducts searches for products
// @Summary Search products
// @Description Search for products by name, description, or brand
// @Tags products
// @Produce json
// @Param q query string true "Search query"
// @Param page query int false "Page number" default(1)
// @Param limit query int false "Items per page" default(20)
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /api/v1/products/search [get]
func (h *ProductHandler) SearchProducts(c *fiber.Ctx) error {
	query := c.Query("q")
	if query == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "error",
			"message": "Search query is required",
			"code":    400,
		})
	}

	page, _ := strconv.Atoi(c.Query("page", "1"))
	limit, _ := strconv.Atoi(c.Query("limit", "20"))
	offset := (page - 1) * limit

	products, total, err := h.productService.SearchProducts(query, limit, offset)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "error",
			"message": err.Error(),
			"code":    400,
		})
	}

	return c.JSON(fiber.Map{
		"status": "success",
		"data": fiber.Map{
			"products": products,
			"pagination": fiber.Map{
				"page":        page,
				"limit":       limit,
				"total_items": total,
				"total_pages": (total + int64(limit) - 1) / int64(limit),
			},
		},
	})
}

// GetProductsByCategory retrieves products by category
// @Summary Get products by category
// @Description Retrieve products filtered by category
// @Tags products
// @Produce json
// @Param category path string true "Product category"
// @Param page query int false "Page number" default(1)
// @Param limit query int false "Items per page" default(20)
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /api/v1/products/category/{category} [get]
func (h *ProductHandler) GetProductsByCategory(c *fiber.Ctx) error {
	category := models.ProductCategory(c.Params("category"))

	page, _ := strconv.Atoi(c.Query("page", "1"))
	limit, _ := strconv.Atoi(c.Query("limit", "20"))
	offset := (page - 1) * limit

	products, total, err := h.productService.GetProductsByCategory(category, limit, offset)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "error",
			"message": err.Error(),
			"code":    400,
		})
	}

	return c.JSON(fiber.Map{
		"status": "success",
		"data": fiber.Map{
			"products": products,
			"pagination": fiber.Map{
				"page":        page,
				"limit":       limit,
				"total_items": total,
				"total_pages": (total + int64(limit) - 1) / int64(limit),
			},
		},
	})
}

// GetFeaturedProducts retrieves featured products
// @Summary Get featured products
// @Description Retrieve featured products based on views and sales
// @Tags products
// @Produce json
// @Param limit query int false "Number of products to return" default(10)
// @Success 200 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /api/v1/products/featured [get]
func (h *ProductHandler) GetFeaturedProducts(c *fiber.Ctx) error {
	limit, _ := strconv.Atoi(c.Query("limit", "10"))

	products, err := h.productService.GetFeaturedProducts(limit)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "error",
			"message": "Failed to retrieve featured products",
			"code":    500,
		})
	}

	return c.JSON(fiber.Map{
		"status": "success",
		"data":   products,
	})
}

// UploadProductMedia uploads media for a product
// @Summary Upload product media
// @Description Upload an image file for a product
// @Tags products
// @Accept multipart/form-data
// @Produce json
// @Param product_id formData int true "Product ID"
// @Param file formData file true "Image file"
// @Param position formData int false "Position"
// @Param is_primary formData bool false "Is primary"
// @Success 200 {object} map[string]interface{} "Upload response"
// @Failure 400 {object} map[string]interface{} "Bad request"
// @Failure 404 {object} map[string]interface{} "Product not found"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /api/v1/products/{id}/media [post]
func (h *ProductHandler) UploadProductMedia(c *fiber.Ctx) error {
	productIDStr := c.Params("id")
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

	// Parse position (optional)
	positionStr := c.FormValue("position")
	position := 0
	if positionStr != "" {
		position, err = strconv.Atoi(positionStr)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"status":  "error",
				"message": "Invalid position",
				"code":    400,
			})
		}
	}

	// Parse is_primary (optional)
	isPrimaryStr := c.FormValue("is_primary")
	isPrimary := false
	if isPrimaryStr == "true" {
		isPrimary = true
	}

	// Get file
	file, err := c.FormFile("file")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "error",
			"message": "No file provided",
			"code":    400,
		})
	}

	// Validate file
	if err := h.mediaService.ValidateImageFile(file); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "error",
			"message": err.Error(),
			"code":    400,
		})
	}

	// Upload file
	response, err := h.mediaService.UploadImage(file, uint(productID), position, isPrimary)
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

// GetProductMedia retrieves all media for a product
// @Summary Get product media
// @Description Retrieve all media files for a specific product
// @Tags products
// @Accept json
// @Produce json
// @Param product_id path int true "Product ID"
// @Success 200 {object} map[string]interface{} "Media list"
// @Failure 400 {object} map[string]interface{} "Bad request"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /api/v1/products/{id}/media [get]
func (h *ProductHandler) GetProductMedia(c *fiber.Ctx) error {
	productIDStr := c.Params("id")
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

	mediaList, err := h.mediaService.GetMediaByProduct(uint(productID))
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "error",
			"message": err.Error(),
			"code":    500,
		})
	}

	return c.JSON(fiber.Map{
		"status": "success",
		"data":   mediaList,
	})
}

// GetBestSellers retrieves best-selling products
// @Summary Get best sellers
// @Description Retrieve best-selling products based on sales count
// @Tags products
// @Produce json
// @Param limit query int false "Number of products to return" default(10)
// @Success 200 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /api/v1/products/bestsellers [get]
func (h *ProductHandler) GetBestSellers(c *fiber.Ctx) error {
	limit, _ := strconv.Atoi(c.Query("limit", "10"))

	products, err := h.productService.GetBestSellers(limit)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "error",
			"message": "Failed to retrieve best sellers",
			"code":    500,
		})
	}

	return c.JSON(fiber.Map{
		"status": "success",
		"data":   products,
	})
}
