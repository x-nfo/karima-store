package handlers

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http/httptest"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/karima-store/internal/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockVariantService for testing
type MockVariantService struct {
	mock.Mock
}

func (m *MockVariantService) CreateVariant(variant *models.ProductVariant) error {
	args := m.Called(variant)
	return args.Error(0)
}

func (m *MockVariantService) GetVariantByID(id uint) (*models.ProductVariant, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.ProductVariant), args.Error(1)
}

func (m *MockVariantService) GetVariantBySKU(sku string) (*models.ProductVariant, error) {
	args := m.Called(sku)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.ProductVariant), args.Error(1)
}

func (m *MockVariantService) GetVariantsByProductID(productID uint) ([]models.ProductVariant, error) {
	args := m.Called(productID)
	return args.Get(0).([]models.ProductVariant), args.Error(1)
}

func (m *MockVariantService) UpdateVariant(id uint, variant *models.ProductVariant) error {
	args := m.Called(id, variant)
	return args.Error(0)
}

func (m *MockVariantService) DeleteVariant(id uint) error {
	args := m.Called(id)
	return args.Error(0)
}

func (m *MockVariantService) UpdateVariantStock(id uint, quantity int) error {
	args := m.Called(id, quantity)
	return args.Error(0)
}

func (m *MockVariantService) GenerateSKU(productName, size, color string) string {
	args := m.Called(productName, size, color)
	return args.String(0)
}

func TestVariantHandler_CreateVariant(t *testing.T) {
	mockService := new(MockVariantService)
	handler := NewVariantHandler(mockService)
	app := fiber.New()
	app.Post("/api/v1/variants", handler.CreateVariant)

	variant := &models.ProductVariant{
		ProductID: 1,
		Name:      "Test Variant",
		Price:     100,
	}

	mockService.On("CreateVariant", mock.AnythingOfType("*models.ProductVariant")).Return(nil)

	body, _ := json.Marshal(variant)
	req := httptest.NewRequest("POST", "/api/v1/variants", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	resp, err := app.Test(req)

	assert.NoError(t, err)
	assert.Equal(t, 201, resp.StatusCode)

	mockService.AssertExpectations(t)
}

func TestVariantHandler_CreateVariant_Error(t *testing.T) {
	mockService := new(MockVariantService)
	handler := NewVariantHandler(mockService)
	app := fiber.New()
	app.Post("/api/v1/variants", handler.CreateVariant)

	variant := &models.ProductVariant{
		ProductID: 1,
		Name:      "Test Variant",
	}

	mockService.On("CreateVariant", mock.AnythingOfType("*models.ProductVariant")).Return(errors.New("creation failed"))

	body, _ := json.Marshal(variant)
	req := httptest.NewRequest("POST", "/api/v1/variants", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	resp, err := app.Test(req)

	assert.NoError(t, err)
	assert.Equal(t, 400, resp.StatusCode)

	mockService.AssertExpectations(t)
}

func TestVariantHandler_GetVariantByID(t *testing.T) {
	mockService := new(MockVariantService)
	handler := NewVariantHandler(mockService)
	app := fiber.New()
	app.Get("/api/v1/variants/:id", handler.GetVariantByID)

	variant := &models.ProductVariant{ID: 1, Name: "Test Variant"}
	mockService.On("GetVariantByID", uint(1)).Return(variant, nil)

	req := httptest.NewRequest("GET", "/api/v1/variants/1", nil)
	resp, err := app.Test(req)

	assert.NoError(t, err)
	assert.Equal(t, 200, resp.StatusCode)

	mockService.AssertExpectations(t)
}

func TestVariantHandler_GetVariantByID_NotFound(t *testing.T) {
	mockService := new(MockVariantService)
	handler := NewVariantHandler(mockService)
	app := fiber.New()
	app.Get("/api/v1/variants/:id", handler.GetVariantByID)

	mockService.On("GetVariantByID", uint(999)).Return(nil, errors.New("not found"))

	req := httptest.NewRequest("GET", "/api/v1/variants/999", nil)
	resp, err := app.Test(req)

	assert.NoError(t, err)
	assert.Equal(t, 404, resp.StatusCode)

	mockService.AssertExpectations(t)
}

// Additional tests for other methods could be added here following the same pattern
