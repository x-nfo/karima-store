package services

import (
	"errors"
	"testing"

	"github.com/karima-store/internal/models"
	"github.com/karima-store/internal/repository"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

// MockCategoryRepository
type MockCategoryRepository struct {
	mock.Mock
}

func (m *MockCategoryRepository) GetAllCategories() []models.ProductCategory {
	args := m.Called()
	return args.Get(0).([]models.ProductCategory)
}

func (m *MockCategoryRepository) GetCategoryStats() ([]repository.CategoryStats, error) {
	args := m.Called()
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]repository.CategoryStats), args.Error(1)
}

// ============================================
// SERVICE INITIALIZATION TESTS
// ============================================

func TestNewCategoryService(t *testing.T) {
	mockRepo := new(MockCategoryRepository)
	service := NewCategoryService(mockRepo)

	assert.NotNil(t, service, "Service should not be nil")
}

func TestCategoryService_ImplementsInterface(t *testing.T) {
	mockRepo := new(MockCategoryRepository)
	service := NewCategoryService(mockRepo)

	// Verify service implements the interface
	var _ CategoryService = service
}

// ============================================
// GET ALL CATEGORIES TESTS
// ============================================

func TestCategoryService_GetAllCategories_Success(t *testing.T) {
	mockRepo := new(MockCategoryRepository)
	service := NewCategoryService(mockRepo)

	expected := []models.ProductCategory{
		models.CategoryTops,
		models.CategoryBottoms,
		models.CategoryDresses,
		models.CategoryOuterwear,
		models.CategoryFootwear,
		models.CategoryAccessories,
	}
	mockRepo.On("GetAllCategories").Return(expected)

	result := service.GetAllCategories()

	assert.Equal(t, expected, result, "Should return all categories")
	mockRepo.AssertExpectations(t)
}

func TestCategoryService_GetAllCategories_EmptyList(t *testing.T) {
	mockRepo := new(MockCategoryRepository)
	service := NewCategoryService(mockRepo)

	expected := []models.ProductCategory{}
	mockRepo.On("GetAllCategories").Return(expected)

	result := service.GetAllCategories()

	assert.Equal(t, expected, result, "Should return empty list when no categories exist")
	assert.Len(t, result, 0, "Result should be empty")
	mockRepo.AssertExpectations(t)
}

func TestCategoryService_GetAllCategories_PartialList(t *testing.T) {
	mockRepo := new(MockCategoryRepository)
	service := NewCategoryService(mockRepo)

	expected := []models.ProductCategory{
		models.CategoryTops,
		models.CategoryBottoms,
	}
	mockRepo.On("GetAllCategories").Return(expected)

	result := service.GetAllCategories()

	assert.Equal(t, expected, result, "Should return partial list of categories")
	assert.Len(t, result, 2, "Should have 2 categories")
	mockRepo.AssertExpectations(t)
}

// ============================================
// GET CATEGORY STATS TESTS
// ============================================

func TestCategoryService_GetCategoryStats_Success(t *testing.T) {
	mockRepo := new(MockCategoryRepository)
	service := NewCategoryService(mockRepo)

	expected := []repository.CategoryStats{
		{Category: models.CategoryTops, ProductCount: 10},
		{Category: models.CategoryBottoms, ProductCount: 15},
		{Category: models.CategoryDresses, ProductCount: 5},
		{Category: models.CategoryOuterwear, ProductCount: 8},
		{Category: models.CategoryFootwear, ProductCount: 12},
		{Category: models.CategoryAccessories, ProductCount: 20},
	}
	mockRepo.On("GetCategoryStats").Return(expected, nil)

	result, err := service.GetCategoryStats()

	require.NoError(t, err, "Should not return error")
	assert.Equal(t, expected, result, "Should return category stats")
	assert.Len(t, result, 6, "Should have stats for all 6 categories")
	mockRepo.AssertExpectations(t)
}

func TestCategoryService_GetCategoryStats_EmptyStats(t *testing.T) {
	mockRepo := new(MockCategoryRepository)
	service := NewCategoryService(mockRepo)

	expected := []repository.CategoryStats{}
	mockRepo.On("GetCategoryStats").Return(expected, nil)

	result, err := service.GetCategoryStats()

	require.NoError(t, err, "Should not return error")
	assert.Equal(t, expected, result, "Should return empty stats")
	assert.Len(t, result, 0, "Should have no stats")
	mockRepo.AssertExpectations(t)
}

func TestCategoryService_GetCategoryStats_RepositoryError(t *testing.T) {
	mockRepo := new(MockCategoryRepository)
	service := NewCategoryService(mockRepo)

	expectedErr := errors.New("database connection failed")
	mockRepo.On("GetCategoryStats").Return(nil, expectedErr)

	result, err := service.GetCategoryStats()

	assert.Error(t, err, "Should return error")
	assert.Nil(t, result, "Result should be nil")
	assert.Contains(t, err.Error(), "database connection failed", "Error should contain repository error message")
	mockRepo.AssertExpectations(t)
}

func TestCategoryService_GetCategoryStats_SingleCategory(t *testing.T) {
	mockRepo := new(MockCategoryRepository)
	service := NewCategoryService(mockRepo)

	expected := []repository.CategoryStats{
		{Category: models.CategoryTops, ProductCount: 100},
	}
	mockRepo.On("GetCategoryStats").Return(expected, nil)

	result, err := service.GetCategoryStats()

	require.NoError(t, err, "Should not return error")
	assert.Equal(t, expected, result, "Should return single category stats")
	assert.Len(t, result, 1, "Should have stats for 1 category")
	assert.Equal(t, int64(100), result[0].ProductCount, "Product count should be 100")
	mockRepo.AssertExpectations(t)
}

func TestCategoryService_GetCategoryStats_ZeroProductCount(t *testing.T) {
	mockRepo := new(MockCategoryRepository)
	service := NewCategoryService(mockRepo)

	expected := []repository.CategoryStats{
		{Category: models.CategoryTops, ProductCount: 0},
		{Category: models.CategoryBottoms, ProductCount: 0},
	}
	mockRepo.On("GetCategoryStats").Return(expected, nil)

	result, err := service.GetCategoryStats()

	require.NoError(t, err, "Should not return error")
	assert.Len(t, result, 2, "Should have stats for 2 categories")
	for _, stat := range result {
		assert.Equal(t, int64(0), stat.ProductCount, "Product count should be 0")
	}
	mockRepo.AssertExpectations(t)
}

// ============================================
// GET CATEGORY NAME TESTS
// ============================================

func TestCategoryService_GetCategoryName_AllValidCategories(t *testing.T) {
	service := NewCategoryService(nil) // Repo not used here

	tests := []struct {
		name     string
		category models.ProductCategory
		expected string
	}{
		{"Tops", models.CategoryTops, "Tops"},
		{"Bottoms", models.CategoryBottoms, "Bottoms"},
		{"Dresses", models.CategoryDresses, "Dresses"},
		{"Outerwear", models.CategoryOuterwear, "Outerwear"},
		{"Footwear", models.CategoryFootwear, "Footwear"},
		{"Accessories", models.CategoryAccessories, "Accessories"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := service.GetCategoryName(tt.category)
			assert.Equal(t, tt.expected, result, "Should return correct display name for %s", tt.name)
		})
	}
}

func TestCategoryService_GetCategoryName_InvalidCategory(t *testing.T) {
	service := NewCategoryService(nil)

	tests := []struct {
		name     string
		category models.ProductCategory
		expected string
	}{
		{"Unknown", models.ProductCategory("Unknown"), "Unknown"},
		{"Invalid", models.ProductCategory("Invalid"), "Invalid"},
		{"Empty", models.ProductCategory(""), ""},
		{"Random", models.ProductCategory("RandomCategory"), "RandomCategory"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := service.GetCategoryName(tt.category)
			assert.Equal(t, tt.expected, result, "Should return the category string as-is for invalid category")
		})
	}
}

func TestCategoryService_GetCategoryName_CaseSensitivity(t *testing.T) {
	service := NewCategoryService(nil)

	// Test that category matching is case-sensitive
	// The constant CategoryTops equals "tops", so it returns display name "Tops"
	result := service.GetCategoryName(models.ProductCategory("tops"))
	assert.Equal(t, "Tops", result, "Should return display name for exact match to constant")

	result = service.GetCategoryName(models.ProductCategory("TOPS"))
	assert.Equal(t, "TOPS", result, "Should return raw string for uppercase (non-matching)")

	result = service.GetCategoryName(models.ProductCategory("Tops"))
	assert.Equal(t, "Tops", result, "Should return raw string for mixed case (non-matching)")

	// Only exact match to constant should return display name
	result = service.GetCategoryName(models.CategoryTops)
	assert.Equal(t, "Tops", result, "Should return display name for exact constant match")
}

func TestCategoryService_GetCategoryName_SpecialCharacters(t *testing.T) {
	service := NewCategoryService(nil)

	tests := []struct {
		name     string
		category models.ProductCategory
		expected string
	}{
		{"With Space", models.ProductCategory("With Space"), "With Space"},
		{"With-Hyphen", models.ProductCategory("With-Hyphen"), "With-Hyphen"},
		{"With_Underscore", models.ProductCategory("With_Underscore"), "With_Underscore"},
		{"With@Special", models.ProductCategory("With@Special"), "With@Special"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := service.GetCategoryName(tt.category)
			assert.Equal(t, tt.expected, result, "Should handle special characters in category name")
		})
	}
}

// ============================================
// IS VALID CATEGORY TESTS
// ============================================

func TestCategoryService_IsValidCategory_AllValidCategories(t *testing.T) {
	service := NewCategoryService(nil)

	tests := []struct {
		name     string
		category models.ProductCategory
		expected bool
	}{
		{"Tops", models.CategoryTops, true},
		{"Bottoms", models.CategoryBottoms, true},
		{"Dresses", models.CategoryDresses, true},
		{"Outerwear", models.CategoryOuterwear, true},
		{"Footwear", models.CategoryFootwear, true},
		{"Accessories", models.CategoryAccessories, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := service.IsValidCategory(tt.category)
			assert.Equal(t, tt.expected, result, "Should validate %s as valid category", tt.name)
		})
	}
}

func TestCategoryService_IsValidCategory_InvalidCategories(t *testing.T) {
	service := NewCategoryService(nil)

	tests := []struct {
		name     string
		category models.ProductCategory
		expected bool
	}{
		{"Unknown", models.ProductCategory("Unknown"), false},
		{"Invalid", models.ProductCategory("Invalid"), false},
		{"Empty", models.ProductCategory(""), false},
		{"Random", models.ProductCategory("RandomCategory"), false},
		{"Numbers", models.ProductCategory("123"), false},
		{"SpecialChars", models.ProductCategory("@#$%"), false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := service.IsValidCategory(tt.category)
			assert.Equal(t, tt.expected, result, "Should validate %s as invalid category", tt.name)
		})
	}
}

func TestCategoryService_IsValidCategory_CaseSensitivity(t *testing.T) {
	service := NewCategoryService(nil)

	// Test that validation is case-sensitive
	// The constant CategoryTops equals "tops", so it's valid
	result := service.IsValidCategory(models.ProductCategory("tops"))
	assert.True(t, result, "Should be valid for exact match to constant 'tops'")

	result = service.IsValidCategory(models.ProductCategory("TOPS"))
	assert.False(t, result, "Should be invalid for uppercase 'TOPS'")

	result = service.IsValidCategory(models.ProductCategory("Tops"))
	assert.False(t, result, "Should be invalid for mixed case 'Tops'")

	// Only exact match to constant should be valid
	result = service.IsValidCategory(models.CategoryTops)
	assert.True(t, result, "Should be valid for exact constant match")
}

func TestCategoryService_IsValidCategory_SimilarNames(t *testing.T) {
	service := NewCategoryService(nil)

	tests := []struct {
		name     string
		category models.ProductCategory
		expected bool
	}{
		{"Top", models.ProductCategory("Top"), false},
		{"Bottom", models.ProductCategory("Bottom"), false},
		{"Dress", models.ProductCategory("Dress"), false},
		{"Outer", models.ProductCategory("Outer"), false},
		{"Foot", models.ProductCategory("Foot"), false},
		{"Accessory", models.ProductCategory("Accessory"), false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := service.IsValidCategory(tt.category)
			assert.Equal(t, tt.expected, result, "Should be invalid for similar but not exact category name")
		})
	}
}

// ============================================
// INTEGRATION TESTS
// ============================================

func TestCategoryService_CategoryNameAndValidationConsistency(t *testing.T) {
	service := NewCategoryService(nil)

	// Test that all valid categories have display names
	validCategories := []models.ProductCategory{
		models.CategoryTops,
		models.CategoryBottoms,
		models.CategoryDresses,
		models.CategoryOuterwear,
		models.CategoryFootwear,
		models.CategoryAccessories,
	}

	for _, category := range validCategories {
		t.Run(string(category), func(t *testing.T) {
			// Should be valid
			assert.True(t, service.IsValidCategory(category), "Category %s should be valid", category)

			// Should have a display name
			name := service.GetCategoryName(category)
			assert.NotEmpty(t, name, "Category %s should have a display name", category)

			// Display name should not equal the raw category value
			assert.NotEqual(t, string(category), name, "Display name should differ from raw value for %s", category)
		})
	}
}

func TestCategoryService_InvalidCategoryFallback(t *testing.T) {
	service := NewCategoryService(nil)

	invalidCategory := models.ProductCategory("InvalidCategory")

	// Should be invalid
	assert.False(t, service.IsValidCategory(invalidCategory), "InvalidCategory should be invalid")

	// Should return the raw string as fallback
	name := service.GetCategoryName(invalidCategory)
	assert.Equal(t, "InvalidCategory", name, "Should return raw string for invalid category")
}

// ============================================
// EDGE CASES AND BOUNDARY TESTS
// ============================================

func TestCategoryService_GetCategoryName_VeryLongCategory(t *testing.T) {
	service := NewCategoryService(nil)

	longCategory := models.ProductCategory("ThisIsAVeryLongCategoryNameThatExceedsNormalLength")
	result := service.GetCategoryName(longCategory)

	assert.Equal(t, string(longCategory), result, "Should handle very long category names")
}

func TestCategoryService_IsValidCategory_UnicodeCharacters(t *testing.T) {
	service := NewCategoryService(nil)

	tests := []struct {
		name     string
		category models.ProductCategory
		expected bool
	}{
		{"Chinese", models.ProductCategory("中文"), false},
		{"Japanese", models.ProductCategory("日本語"), false},
		{"Arabic", models.ProductCategory("العربية"), false},
		{"Emoji", models.ProductCategory("👕"), false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := service.IsValidCategory(tt.category)
			assert.Equal(t, tt.expected, result, "Should handle unicode characters")
		})
	}
}

func TestCategoryService_GetCategoryName_VeryShortCategory(t *testing.T) {
	service := NewCategoryService(nil)

	tests := []struct {
		name     string
		category models.ProductCategory
		expected string
	}{
		{"Single char", models.ProductCategory("A"), "A"},
		{"Two chars", models.ProductCategory("AB"), "AB"},
		{"Three chars", models.ProductCategory("ABC"), "ABC"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := service.GetCategoryName(tt.category)
			assert.Equal(t, tt.expected, result, "Should handle very short category names")
		})
	}
}

// ============================================
// PERFORMANCE AND STRESS TESTS
// ============================================

func TestCategoryService_GetCategoryName_Performance(t *testing.T) {
	service := NewCategoryService(nil)

	// Test performance with repeated calls
	iterations := 1000
	category := models.CategoryTops

	for i := 0; i < iterations; i++ {
		result := service.GetCategoryName(category)
		assert.Equal(t, "Tops", result)
	}
}

func TestCategoryService_IsValidCategory_Performance(t *testing.T) {
	service := NewCategoryService(nil)

	// Test performance with repeated calls
	iterations := 1000
	category := models.CategoryTops

	for i := 0; i < iterations; i++ {
		result := service.IsValidCategory(category)
		assert.True(t, result)
	}
}

// ============================================
// CATEGORY ENUMERATION TESTS
// ============================================

func TestCategoryService_AllCategoriesCovered(t *testing.T) {
	service := NewCategoryService(nil)

	// Verify all 6 predefined categories are covered
	allValidCategories := []models.ProductCategory{
		models.CategoryTops,
		models.CategoryBottoms,
		models.CategoryDresses,
		models.CategoryOuterwear,
		models.CategoryFootwear,
		models.CategoryAccessories,
	}

	for _, category := range allValidCategories {
		t.Run(string(category), func(t *testing.T) {
			// Should be valid
			assert.True(t, service.IsValidCategory(category))

			// Should have a non-empty display name
			name := service.GetCategoryName(category)
			assert.NotEmpty(t, name)

			// Display name should be different from the enum value
			assert.NotEqual(t, string(category), name)
		})
	}
}

func TestCategoryService_CategoryConstants(t *testing.T) {
	_ = NewCategoryService(nil)

	// Verify category constants have expected values
	assert.Equal(t, models.ProductCategory("tops"), models.CategoryTops)
	assert.Equal(t, models.ProductCategory("bottoms"), models.CategoryBottoms)
	assert.Equal(t, models.ProductCategory("dresses"), models.CategoryDresses)
	assert.Equal(t, models.ProductCategory("outerwear"), models.CategoryOuterwear)
	assert.Equal(t, models.ProductCategory("footwear"), models.CategoryFootwear)
	assert.Equal(t, models.ProductCategory("accessories"), models.CategoryAccessories)
}

// ============================================
// MOCK VERIFICATION TESTS
// ============================================

func TestCategoryService_RepositoryCalls(t *testing.T) {
	mockRepo := new(MockCategoryRepository)
	service := NewCategoryService(mockRepo)

	// Test GetCategoryName doesn't call repository (it uses internal map)
	service.GetCategoryName(models.CategoryTops)
	mockRepo.AssertNotCalled(t, "GetAllCategories")
	mockRepo.AssertNotCalled(t, "GetCategoryStats")

	// Test IsValidCategory doesn't call repository (it uses internal map)
	service.IsValidCategory(models.CategoryTops)
	mockRepo.AssertNotCalled(t, "GetAllCategories")
	mockRepo.AssertNotCalled(t, "GetCategoryStats")

	// Test GetAllCategories calls repository
	mockRepo.On("GetAllCategories").Return([]models.ProductCategory{models.CategoryTops}).Once()
	service.GetAllCategories()
	mockRepo.AssertExpectations(t)

	// Test GetCategoryStats calls repository
	mockRepo.On("GetCategoryStats").Return([]repository.CategoryStats{}, nil).Once()
	service.GetCategoryStats()
	mockRepo.AssertExpectations(t)
}
