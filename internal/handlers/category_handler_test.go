package handlers

import (
	"encoding/json"
	"errors"
	"net/http/httptest"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/karima-store/internal/models"
	"github.com/karima-store/internal/repository"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockCategoryService for testing
type MockCategoryService struct {
	mock.Mock
}

func (m *MockCategoryService) GetAllCategories() []models.ProductCategory {
	args := m.Called()
	return args.Get(0).([]models.ProductCategory)
}

func (m *MockCategoryService) GetCategoryStats() ([]repository.CategoryStats, error) {
	args := m.Called()
	return args.Get(0).([]repository.CategoryStats), args.Error(1)
}

func (m *MockCategoryService) GetCategoryName(category models.ProductCategory) string {
	args := m.Called(category)
	return args.String(0)
}

func (m *MockCategoryService) IsValidCategory(category models.ProductCategory) bool {
	args := m.Called(category)
	return args.Bool(0)
}

func TestCategoryHandler_GetAllCategories(t *testing.T) {
	mockService := new(MockCategoryService)
	handler := NewCategoryHandler(mockService)
	app := fiber.New()
	app.Get("/api/v1/categories", handler.GetAllCategories)

	mockService.On("GetAllCategories").Return([]models.ProductCategory{models.CategoryTops, models.CategoryBottoms})
	mockService.On("GetCategoryName", models.CategoryTops).Return("Tops")
	mockService.On("GetCategoryName", models.CategoryBottoms).Return("Bottoms")

	req := httptest.NewRequest("GET", "/api/v1/categories", nil)
	resp, err := app.Test(req)

	assert.NoError(t, err)
	assert.Equal(t, 200, resp.StatusCode)

	var response map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&response)

	assert.Equal(t, "success", response["status"])
	data := response["data"].([]interface{})
	assert.Len(t, data, 2)

	mockService.AssertExpectations(t)
}

func TestCategoryHandler_GetCategoryStats(t *testing.T) {
	mockService := new(MockCategoryService)
	handler := NewCategoryHandler(mockService)
	app := fiber.New()
	app.Get("/api/v1/categories/stats", handler.GetCategoryStats)

	stats := []repository.CategoryStats{
		{Category: models.CategoryTops, ProductCount: 10},
	}

	mockService.On("GetCategoryStats").Return(stats, nil)
	mockService.On("GetCategoryName", models.CategoryTops).Return("Tops")

	req := httptest.NewRequest("GET", "/api/v1/categories/stats", nil)
	resp, err := app.Test(req)

	assert.NoError(t, err)
	assert.Equal(t, 200, resp.StatusCode)

	mockService.AssertExpectations(t)
}

func TestCategoryHandler_GetCategoryStats_Error(t *testing.T) {
	mockService := new(MockCategoryService)
	handler := NewCategoryHandler(mockService)
	app := fiber.New()
	app.Get("/api/v1/categories/stats", handler.GetCategoryStats)

	mockService.On("GetCategoryStats").Return([]repository.CategoryStats{}, errors.New("db error"))

	req := httptest.NewRequest("GET", "/api/v1/categories/stats", nil)
	resp, err := app.Test(req)

	assert.NoError(t, err)
	assert.Equal(t, 500, resp.StatusCode)

	mockService.AssertExpectations(t)
}

func TestCategoryHandler_GetCategoryProducts(t *testing.T) {
	mockService := new(MockCategoryService)
	handler := NewCategoryHandler(mockService)
	app := fiber.New()
	app.Get("/api/v1/categories/:category/products", handler.GetCategoryProducts)

	mockService.On("IsValidCategory", models.CategoryTops).Return(true)
	mockService.On("GetCategoryName", models.CategoryTops).Return("Tops")

	req := httptest.NewRequest("GET", "/api/v1/categories/tops/products", nil)
	resp, err := app.Test(req)

	assert.NoError(t, err)
	assert.Equal(t, 200, resp.StatusCode)

	mockService.AssertExpectations(t)
}

func TestCategoryHandler_GetCategoryProducts_Invalid(t *testing.T) {
	mockService := new(MockCategoryService)
	handler := NewCategoryHandler(mockService)
	app := fiber.New()
	app.Get("/api/v1/categories/:category/products", handler.GetCategoryProducts)

	mockService.On("IsValidCategory", models.ProductCategory("invalid")).Return(false)

	req := httptest.NewRequest("GET", "/api/v1/categories/invalid/products", nil)
	resp, err := app.Test(req)

	assert.NoError(t, err)
	assert.Equal(t, 400, resp.StatusCode)

	mockService.AssertExpectations(t)
}
