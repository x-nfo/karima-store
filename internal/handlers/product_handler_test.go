package handlers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http/httptest"
	"testing"

	apperrors "karima_store/internal/errors"
	"karima_store/internal/models"
	"karima_store/internal/repository"
	"karima_store/internal/services"
	"karima_store/internal/test_setup"

	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func setupProductHandlerTest(t *testing.T) (*fiber.App, func()) {
	// Setup test database
	db, cleanupDB := test_setup.SetupTestDB(t)

	// Clean up any existing data
	db.Exec("DELETE FROM products")

	// Setup test Redis
	redisClient := test_setup.SetupTestRedis(t)

	// Create repositories
	productRepo := repository.NewProductRepository(db)
	variantRepo := repository.NewVariantRepository(db)

	// Create service
	productService := services.NewProductService(productRepo, variantRepo, redisClient)

	// Create handler
	productHandler := NewProductHandler(productService)

	// Create Fiber app
	app := fiber.New(fiber.Config{
		ErrorHandler: func(c *fiber.Ctx, err error) error {
			if appErr := apperrors.GetAppError(err); appErr != nil {
				return c.Status(appErr.StatusCode).JSON(fiber.Map{
					"code":    string(appErr.Code),
					"message": appErr.Message,
				})
			}
			return c.Status(500).JSON(fiber.Map{
				"code":    "INTERNAL_ERROR",
				"message": "Internal server error",
			})
		},
	})

	// Setup routes
	app.Get("/products", productHandler.GetProducts)
	app.Get("/products/:id", productHandler.GetProduct)
	app.Post("/products", productHandler.CreateProduct)
	app.Put("/products/:id", productHandler.UpdateProduct)
	app.Delete("/products/:id", productHandler.DeleteProduct)

	cleanup := func() {
		cleanupDB()
		redisClient.Close()
	}

	return app, cleanup
}

func TestProductHandler_GetProducts(t *testing.T) {
	app, cleanup := setupProductHandlerTest(t)
	defer cleanup()

	// Create test products
	for i := 1; i <= 5; i++ {
		product := &models.Product{
			Name:        fmt.Sprintf("Test Product %d", i),
			Description: "Test Description",
			Price:       float64(i * 100),
			Category:    models.CategoryTops,
			Stock:       10,
			Status:      models.StatusAvailable,
			Slug:        fmt.Sprintf("test-product-%d", i),
		}
		// Note: In a real test, you would use the repository to create products
	}

	// Test GET /products
	req := httptest.NewRequest("GET", "/products", nil)
	resp, err := app.Test(req)
	require.NoError(t, err)
	assert.Equal(t, fiber.StatusOK, resp.StatusCode)
}

func TestProductHandler_GetProduct(t *testing.T) {
	app, cleanup := setupProductHandlerTest(t)
	defer cleanup()

	// Test GET /products/:id
	req := httptest.NewRequest("GET", "/products/1", nil)
	resp, err := app.Test(req)
	require.NoError(t, err)

	// Should return 404 for non-existent product
	assert.Equal(t, fiber.StatusNotFound, resp.StatusCode)
}

func TestProductHandler_CreateProduct(t *testing.T) {
	app, cleanup := setupProductHandlerTest(t)
	defer cleanup()

	// Test data
	newProduct := map[string]interface{}{
		"name":        "New Product",
		"description": "Product Description",
		"price":       100.00,
		"category":    "tops",
		"stock":       10,
		"status":      "available",
	}

	body, _ := json.Marshal(newProduct)
	req := httptest.NewRequest("POST", "/products", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	resp, err := app.Test(req)
	require.NoError(t, err)

	// Should return 201 Created
	assert.Equal(t, fiber.StatusCreated, resp.StatusCode)
}

func TestProductHandler_CreateProduct_ValidationError(t *testing.T) {
	app, cleanup := setupProductHandlerTest(t)
	defer cleanup()

	// Test with missing required fields
	newProduct := map[string]interface{}{
		"description": "Product Description",
		// Missing name, price, category
	}

	body, _ := json.Marshal(newProduct)
	req := httptest.NewRequest("POST", "/products", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	resp, err := app.Test(req)
	require.NoError(t, err)

	// Should return 400 Bad Request
	assert.Equal(t, fiber.StatusBadRequest, resp.StatusCode)
}

func TestProductHandler_CreateProduct_InvalidPrice(t *testing.T) {
	app, cleanup := setupProductHandlerTest(t)
	defer cleanup()

	// Test with invalid price
	newProduct := map[string]interface{}{
		"name":        "New Product",
		"description": "Product Description",
		"price":       -10.00, // Invalid price
		"category":    "tops",
		"stock":       10,
		"status":      "available",
	}

	body, _ := json.Marshal(newProduct)
	req := httptest.NewRequest("POST", "/products", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	resp, err := app.Test(req)
	require.NoError(t, err)

	// Should return 400 Bad Request
	assert.Equal(t, fiber.StatusBadRequest, resp.StatusCode)
}

func TestProductHandler_UpdateProduct(t *testing.T) {
	app, cleanup := setupProductHandlerTest(t)
	defer cleanup()

	// Test data
	updateProduct := map[string]interface{}{
		"name":  "Updated Product",
		"price": 150.00,
	}

	body, _ := json.Marshal(updateProduct)
	req := httptest.NewRequest("PUT", "/products/1", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	resp, err := app.Test(req)
	require.NoError(t, err)

	// Should return 404 for non-existent product
	assert.Equal(t, fiber.StatusNotFound, resp.StatusCode)
}

func TestProductHandler_DeleteProduct(t *testing.T) {
	app, cleanup := setupProductHandlerTest(t)
	defer cleanup()

	// Test DELETE /products/:id
	req := httptest.NewRequest("DELETE", "/products/1", nil)
	resp, err := app.Test(req)
	require.NoError(t, err)

	// Should return 404 for non-existent product
	assert.Equal(t, fiber.StatusNotFound, resp.StatusCode)
}

func TestProductHandler_SearchProducts(t *testing.T) {
	app, cleanup := setupProductHandlerTest(t)
	defer cleanup()

	// Test GET /products/search?query=test
	req := httptest.NewRequest("GET", "/products/search?query=test", nil)
	resp, err := app.Test(req)
	require.NoError(t, err)

	// Should return 200 OK
	assert.Equal(t, fiber.StatusOK, resp.StatusCode)
}

func TestProductHandler_GetProductsByCategory(t *testing.T) {
	app, cleanup := setupProductHandlerTest(t)
	defer cleanup()

	// Test GET /products/category/tops
	req := httptest.NewRequest("GET", "/products/category/tops", nil)
	resp, err := app.Test(req)
	require.NoError(t, err)

	// Should return 200 OK
	assert.Equal(t, fiber.StatusOK, resp.StatusCode)
}

func TestProductHandler_GetFeaturedProducts(t *testing.T) {
	app, cleanup := setupProductHandlerTest(t)
	defer cleanup()

	// Test GET /products/featured
	req := httptest.NewRequest("GET", "/products/featured", nil)
	resp, err := app.Test(req)
	require.NoError(t, err)

	// Should return 200 OK
	assert.Equal(t, fiber.StatusOK, resp.StatusCode)
}

func TestProductHandler_GetBestSellers(t *testing.T) {
	app, cleanup := setupProductHandlerTest(t)
	defer cleanup()

	// Test GET /products/bestsellers
	req := httptest.NewRequest("GET", "/products/bestsellers", nil)
	resp, err := app.Test(req)
	require.NoError(t, err)

	// Should return 200 OK
	assert.Equal(t, fiber.StatusOK, resp.StatusCode)
}

func TestProductHandler_Pagination(t *testing.T) {
	app, cleanup := setupProductHandlerTest(t)
	defer cleanup()

	// Test GET /products?limit=10&offset=0
	req := httptest.NewRequest("GET", "/products?limit=10&offset=0", nil)
	resp, err := app.Test(req)
	require.NoError(t, err)

	// Should return 200 OK
	assert.Equal(t, fiber.StatusOK, resp.StatusCode)
}

func TestProductHandler_InvalidPagination(t *testing.T) {
	app, cleanup := setupProductHandlerTest(t)
	defer cleanup()

	// Test with invalid limit
	req := httptest.NewRequest("GET", "/products?limit=1000", nil)
	resp, err := app.Test(req)
	require.NoError(t, err)

	// Should return 200 OK (limit should be capped)
	assert.Equal(t, fiber.StatusOK, resp.StatusCode)
}

func TestProductHandler_InvalidJSON(t *testing.T) {
	app, cleanup := setupProductHandlerTest(t)
	defer cleanup()

	// Test with invalid JSON
	req := httptest.NewRequest("POST", "/products", bytes.NewReader([]byte("{invalid json")))
	req.Header.Set("Content-Type", "application/json")

	resp, err := app.Test(req)
	require.NoError(t, err)

	// Should return 400 Bad Request
	assert.Equal(t, fiber.StatusBadRequest, resp.StatusCode)
}

func TestProductHandler_MissingContentType(t *testing.T) {
	app, cleanup := setupProductHandlerTest(t)
	defer cleanup()

	// Test without Content-Type header
	newProduct := map[string]interface{}{
		"name":     "New Product",
		"price":    100.00,
		"category": "tops",
	}

	body, _ := json.Marshal(newProduct)
	req := httptest.NewRequest("POST", "/products", bytes.NewReader(body))

	resp, err := app.Test(req)
	require.NoError(t, err)

	// Should return 400 Bad Request
	assert.Equal(t, fiber.StatusBadRequest, resp.StatusCode)
}

func TestProductHandler_MethodNotAllowed(t *testing.T) {
	app, cleanup := setupProductHandlerTest(t)
	defer cleanup()

	// Test unsupported method
	req := httptest.NewRequest("PATCH", "/products/1", nil)
	resp, err := app.Test(req)
	require.NoError(t, err)

	// Should return 405 Method Not Allowed
	assert.Equal(t, fiber.StatusMethodNotAllowed, resp.StatusCode)
}
