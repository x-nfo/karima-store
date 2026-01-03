package services

import (
	"testing"

	"github.com/karima-store/internal/models"
	"github.com/karima-store/internal/repository"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
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

func TestCategoryService_GetAllCategories(t *testing.T) {
	mockRepo := new(MockCategoryRepository)
	service := NewCategoryService(mockRepo)

	expected := []models.ProductCategory{models.CategoryTops, models.CategoryBottoms}
	mockRepo.On("GetAllCategories").Return(expected)

	result := service.GetAllCategories()

	assert.Equal(t, expected, result)
	mockRepo.AssertExpectations(t)
}

func TestCategoryService_GetCategoryStats(t *testing.T) {
	mockRepo := new(MockCategoryRepository)
	service := NewCategoryService(mockRepo)

	expected := []repository.CategoryStats{
		{Category: models.CategoryTops, ProductCount: 10},
	}
	mockRepo.On("GetCategoryStats").Return(expected, nil)

	result, err := service.GetCategoryStats()

	assert.NoError(t, err)
	assert.Equal(t, expected, result)
	mockRepo.AssertExpectations(t)
}

func TestCategoryService_GetCategoryName(t *testing.T) {
	service := NewCategoryService(nil) // Repo not used here

	tests := []struct {
		category models.ProductCategory
		expected string
	}{
		{models.CategoryTops, "Tops"},
		{models.CategoryBottoms, "Bottoms"},
		{models.CategoryDresses, "Dresses"},
		{models.CategoryOuterwear, "Outerwear"},
		{models.CategoryFootwear, "Footwear"},
		{models.CategoryAccessories, "Accessories"},
		{models.ProductCategory("Unknown"), "Unknown"},
	}

	for _, tt := range tests {
		t.Run(string(tt.category), func(t *testing.T) {
			result := service.GetCategoryName(tt.category)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestCategoryService_IsValidCategory(t *testing.T) {
	service := NewCategoryService(nil) // Repo not used here

	tests := []struct {
		category models.ProductCategory
		expected bool
	}{
		{models.CategoryTops, true},
		{models.CategoryBottoms, true},
		{models.ProductCategory("Invalid"), false},
	}

	for _, tt := range tests {
		t.Run(string(tt.category), func(t *testing.T) {
			result := service.IsValidCategory(tt.category)
			assert.Equal(t, tt.expected, result)
		})
	}
}
