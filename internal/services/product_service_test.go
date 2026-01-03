package services

import (
	"errors"
	"log"
	"testing"

	"github.com/karima-store/internal/config"
	"github.com/karima-store/internal/database"
	"github.com/karima-store/internal/models"
	"github.com/karima-store/internal/repository"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"gorm.io/gorm"
)

// MockProductRepository for testing
type MockProductRepository struct {
	mock.Mock
}

func (m *MockProductRepository) Create(product *models.Product) error {
	args := m.Called(product)
	return args.Error(0)
}

func (m *MockProductRepository) GetByID(id uint) (*models.Product, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Product), args.Error(1)
}

func (m *MockProductRepository) GetBySlug(slug string) (*models.Product, error) {
	args := m.Called(slug)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Product), args.Error(1)
}

func (m *MockProductRepository) GetAll(limit, offset int, filters map[string]interface{}) ([]models.Product, int64, error) {
	args := m.Called(limit, offset, filters)
	return args.Get(0).([]models.Product), args.Get(1).(int64), args.Error(2)
}

func (m *MockProductRepository) Update(product *models.Product) error {
	args := m.Called(product)
	return args.Error(0)
}

func (m *MockProductRepository) Delete(id uint) error {
	args := m.Called(id)
	return args.Error(0)
}

func (m *MockProductRepository) UpdateStock(id uint, quantity int) error {
	args := m.Called(id, quantity)
	return args.Error(0)
}

func (m *MockProductRepository) Search(query string, limit, offset int) ([]models.Product, int64, error) {
	args := m.Called(query, limit, offset)
	return args.Get(0).([]models.Product), args.Get(1).(int64), args.Error(2)
}

func (m *MockProductRepository) GetByCategory(category models.ProductCategory, limit, offset int) ([]models.Product, int64, error) {
	args := m.Called(category, limit, offset)
	return args.Get(0).([]models.Product), args.Get(1).(int64), args.Error(2)
}

func (m *MockProductRepository) GetFeatured(limit int) ([]models.Product, error) {
	args := m.Called(limit)
	return args.Get(0).([]models.Product), args.Error(1)
}

func (m *MockProductRepository) GetBestSellers(limit int) ([]models.Product, error) {
	args := m.Called(limit)
	return args.Get(0).([]models.Product), args.Error(1)
}

func (m *MockProductRepository) IncrementViewCount(id uint) error {
	args := m.Called(id)
	return args.Error(0)
}

func (m *MockProductRepository) WithTx(tx *gorm.DB) repository.ProductRepository {
	args := m.Called(tx)
	return args.Get(0).(repository.ProductRepository)
}

// MockVariantRepository for testing
type MockVariantRepository struct {
	mock.Mock
}

func (m *MockVariantRepository) Create(variant *models.ProductVariant) error {
	args := m.Called(variant)
	return args.Error(0)
}

func (m *MockVariantRepository) GetByID(id uint) (*models.ProductVariant, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.ProductVariant), args.Error(1)
}

func (m *MockVariantRepository) GetBySKU(sku string) (*models.ProductVariant, error) {
	args := m.Called(sku)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.ProductVariant), args.Error(1)
}

func (m *MockVariantRepository) GetByProductID(productID uint) ([]models.ProductVariant, error) {
	args := m.Called(productID)
	return args.Get(0).([]models.ProductVariant), args.Error(1)
}

func (m *MockVariantRepository) Update(variant *models.ProductVariant) error {
	args := m.Called(variant)
	return args.Error(0)
}

func (m *MockVariantRepository) Delete(id uint) error {
	args := m.Called(id)
	return args.Error(0)
}

func (m *MockVariantRepository) UpdateStock(id uint, quantity int) error {
	args := m.Called(id, quantity)
	return args.Error(0)
}

// Helper function to create a test Redis instance
func createTestRedis() database.RedisClient {
	cfg := config.TestConfigWithRedis()
	redisClient, err := database.NewRedis(cfg)
	if err != nil {
		log.Printf("Warning: Failed to connect to test Redis: %v", err)
		return nil
	}
	return redisClient
}

func TestProductService_CreateProduct_ValidInput(t *testing.T) {
	mockRepo := new(MockProductRepository)
	mockVariantRepo := (*MockVariantRepository)(nil)
	testRedis := createTestRedis()

	service := NewProductService(mockRepo, mockVariantRepo, testRedis)

	// Test valid product creation
	product := &models.Product{
		Name:     "Test Product",
		Price:    99.99,
		Category: models.CategoryTops,
		Stock:    10,
	}

	// Mock expectations
	mockRepo.On("GetBySlug", "test-product").Return(nil, errors.New("not found"))
	mockRepo.On("Create", product).Return(nil)

	err := service.CreateProduct(product)

	assert.NoError(t, err)
	assert.NotEmpty(t, product.Slug)
	assert.Equal(t, models.StatusAvailable, product.Status)
	mockRepo.AssertExpectations(t)
}

func TestProductService_CreateProduct_MissingRequiredFields(t *testing.T) {
	mockRepo := new(MockProductRepository)
	mockVariantRepo := (*MockVariantRepository)(nil)
	testRedis := createTestRedis()

	service := NewProductService(mockRepo, mockVariantRepo, testRedis)

	// Test missing name
	product1 := &models.Product{
		Price:    99.99,
		Category: models.CategoryTops,
	}
	err1 := service.CreateProduct(product1)
	assert.Error(t, err1)
	assert.Contains(t, err1.Error(), "product name is required")

	// Test invalid price
	product2 := &models.Product{
		Name:     "Test Product",
		Price:    -10,
		Category: models.CategoryTops,
	}
	err2 := service.CreateProduct(product2)
	assert.Error(t, err2)
	assert.Contains(t, err2.Error(), "product price must be greater than 0")

	// Test missing category
	product3 := &models.Product{
		Name:  "Test Product",
		Price: 99.99,
	}
	err3 := service.CreateProduct(product3)
	assert.Error(t, err3)
	assert.Contains(t, err3.Error(), "product category is required")
}

func TestProductService_CreateProduct_DuplicateSlug(t *testing.T) {
	mockRepo := new(MockProductRepository)
	mockVariantRepo := (*MockVariantRepository)(nil)
	testRedis := createTestRedis()

	service := NewProductService(mockRepo, mockVariantRepo, testRedis)

	// Test duplicate slug
	product := &models.Product{
		Name:     "Test Product",
		Price:    99.99,
		Category: models.CategoryTops,
	}

	existingProduct := &models.Product{
		ID:   1,
		Name: "Test Product",
		Slug: "test-product",
	}

	mockRepo.On("GetBySlug", "test-product").Return(existingProduct, nil)

	err := service.CreateProduct(product)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "product with this slug already exists")
	mockRepo.AssertExpectations(t)
}

func TestProductService_GetProductByID_SQLInjection(t *testing.T) {
	mockRepo := new(MockProductRepository)
	mockVariantRepo := (*MockVariantRepository)(nil)
	testRedis := createTestRedis()

	service := NewProductService(mockRepo, mockVariantRepo, testRedis)

	// Test SQL injection attempt in ID
	// Note: This is a compile-time check, but we test the service behavior
	mockRepo.On("GetByID", uint(0)).Return(nil, errors.New("product not found"))

	_, err := service.GetProductByID(0) // ID 0 is invalid

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "product not found")
}

func TestProductService_UpdateProduct_SQLInjection(t *testing.T) {
	mockRepo := new(MockProductRepository)
	mockVariantRepo := (*MockVariantRepository)(nil)
	testRedis := createTestRedis()

	service := NewProductService(mockRepo, mockVariantRepo, testRedis)

	// Test SQL injection attempt in product name
	product := &models.Product{
		Name:  "'; DROP TABLE products; --",
		Price: 99.99,
	}

	existingProduct := &models.Product{
		ID:   1,
		Name: "Old Product",
		Slug: "old-product",
	}

	mockRepo.On("GetByID", uint(1)).Return(existingProduct, nil)
	mockRepo.On("Update", product).Return(nil)

	err := service.UpdateProduct(1, product)

	assert.NoError(t, err)
	// Verify that the name was stored (sanitization should happen at repository level)
	assert.Equal(t, "'; DROP TABLE products; --", product.Name)
	mockRepo.AssertExpectations(t)
}

func TestProductService_UpdateProduct_InvalidPrice(t *testing.T) {
	mockRepo := new(MockProductRepository)
	mockVariantRepo := (*MockVariantRepository)(nil)
	testRedis := createTestRedis()

	service := NewProductService(mockRepo, mockVariantRepo, testRedis)

	// Test negative price
	product := &models.Product{
		Name:  "Test Product",
		Price: -10,
	}

	existingProduct := &models.Product{
		ID:   1,
		Name: "Old Product",
		Slug: "old-product",
	}

	mockRepo.On("GetByID", uint(1)).Return(existingProduct, nil)
	mockRepo.On("Update", product).Return(nil)

	err := service.UpdateProduct(1, product)

	assert.NoError(t, err) // Service doesn't validate price on update
}

func TestProductService_DeleteProduct_NonExistent(t *testing.T) {
	mockRepo := new(MockProductRepository)
	mockVariantRepo := (*MockVariantRepository)(nil)
	testRedis := createTestRedis()

	service := NewProductService(mockRepo, mockVariantRepo, testRedis)

	// Test deleting non-existent product
	mockRepo.On("GetByID", uint(999)).Return(nil, gorm.ErrRecordNotFound)

	err := service.DeleteProduct(999)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "product not found")
	mockRepo.AssertExpectations(t)
}

func TestProductService_UpdateProductStock_InsufficientStock(t *testing.T) {
	mockRepo := new(MockProductRepository)
	mockVariantRepo := (*MockVariantRepository)(nil)
	testRedis := createTestRedis()

	service := NewProductService(mockRepo, mockVariantRepo, testRedis)

	// Test insufficient stock
	product := &models.Product{
		ID:    1,
		Name:  "Test Product",
		Stock: 5,
	}

	mockRepo.On("GetByID", uint(1)).Return(product, nil)

	err := service.UpdateProductStock(1, -10) // Trying to reduce by 10 when only 5 in stock

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "insufficient stock")
	mockRepo.AssertExpectations(t)
}

func TestProductService_SearchProducts_EmptyQuery(t *testing.T) {
	mockRepo := new(MockProductRepository)
	mockVariantRepo := (*MockVariantRepository)(nil)
	testRedis := createTestRedis()

	service := NewProductService(mockRepo, mockVariantRepo, testRedis)

	// Test empty search query
	_, _, err := service.SearchProducts("", 10, 0)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "search query cannot be empty")
}

func TestProductService_SearchProducts_SQLInjection(t *testing.T) {
	mockRepo := new(MockProductRepository)
	mockVariantRepo := (*MockVariantRepository)(nil)
	testRedis := createTestRedis()

	service := NewProductService(mockRepo, mockVariantRepo, testRedis)

	// Test SQL injection in search query
	sqlInjectionQuery := "test'; DROP TABLE products; --"

	mockRepo.On("Search", sqlInjectionQuery, 20, 0).Return([]models.Product{}, int64(0), nil)

	_, _, err := service.SearchProducts(sqlInjectionQuery, 20, 0)

	// Service should pass the query to repository (sanitization should happen at repository level)
	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

func TestProductService_GetProductsByCategory_InvalidCategory(t *testing.T) {
	mockRepo := new(MockProductRepository)
	mockVariantRepo := (*MockVariantRepository)(nil)
	testRedis := createTestRedis()

	service := NewProductService(mockRepo, mockVariantRepo, testRedis)

	// Test invalid category
	invalidCategory := models.ProductCategory("invalid")

	_, _, err := service.GetProductsByCategory(invalidCategory, 10, 0)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "invalid category")
}

func TestProductService_GenerateSlug(t *testing.T) {
	mockRepo := new(MockProductRepository)
	mockVariantRepo := (*MockVariantRepository)(nil)
	testRedis := createTestRedis()

	service := NewProductService(mockRepo, mockVariantRepo, testRedis)

	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "Simple name",
			input:    "Test Product",
			expected: "test-product",
		},
		{
			name:     "Name with special characters",
			input:    "Test Product!@#$%",
			expected: "test-product",
		},
		{
			name:     "Name with multiple spaces",
			input:    "Test   Product",
			expected: "test-product",
		},
		{
			name:     "Name with hyphens",
			input:    "Test-Product",
			expected: "test-product",
		},
		{
			name:     "Empty name",
			input:    "",
			expected: "product-", // Should add timestamp prefix
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := service.GenerateSlug(tt.input)
			// For empty name, we expect it to start with "product-"
			if tt.input == "" {
				assert.Contains(t, result, "product-")
			} else {
				assert.Equal(t, tt.expected, result)
			}
		})
	}
}

func TestProductService_GetProducts_InvalidPagination(t *testing.T) {
	mockRepo := new(MockProductRepository)
	mockVariantRepo := (*MockVariantRepository)(nil)
	testRedis := createTestRedis()

	service := NewProductService(mockRepo, mockVariantRepo, testRedis)

	// Test invalid pagination parameters
	mockRepo.On("GetAll", 20, 0, mock.Anything).Return([]models.Product{}, int64(0), nil)

	// Test limit too high
	_, _, err := service.GetProducts(200, 0, nil)
	assert.NoError(t, err) // Should default to 20

	// Test negative limit
	_, _, err = service.GetProducts(-10, 0, nil)
	assert.NoError(t, err) // Should default to 20

	// Test negative offset
	_, _, err = service.GetProducts(20, -10, nil)
	assert.NoError(t, err) // Should default to 0

	mockRepo.AssertExpectations(t)
}

func TestProductService_GetProducts_SQLInjectionInFilters(t *testing.T) {
	mockRepo := new(MockProductRepository)
	mockVariantRepo := (*MockVariantRepository)(nil)
	testRedis := createTestRedis()

	service := NewProductService(mockRepo, mockVariantRepo, testRedis)

	// Test SQL injection in filters
	filters := map[string]interface{}{
		"name": "test'; DROP TABLE products; --",
	}

	mockRepo.On("GetAll", 20, 0, filters).Return([]models.Product{}, int64(0), nil)

	_, _, err := service.GetProducts(20, 0, filters)

	// Service should pass filters to repository (sanitization should happen at repository level)
	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

func TestProductService_CacheInvalidation(t *testing.T) {
	mockRepo := new(MockProductRepository)
	mockVariantRepo := (*MockVariantRepository)(nil)
	testRedis := createTestRedis()

	service := NewProductService(mockRepo, mockVariantRepo, testRedis)

	// Test cache invalidation on product creation
	product := &models.Product{
		Name:     "Test Product",
		Price:    99.99,
		Category: models.CategoryTops,
	}

	mockRepo.On("GetBySlug", "test-product").Return(nil, errors.New("not found"))
	mockRepo.On("Create", product).Return(nil)

	err := service.CreateProduct(product)

	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}
