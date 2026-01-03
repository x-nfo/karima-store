package repository

import (
	"testing"

	"github.com/karima-store/internal/models"
	"github.com/karima-store/internal/test_setup"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"
)

func setupCategoryTest(t *testing.T) (*gorm.DB, func()) {
	db, cleanup := test_setup.SetupTestDB(t)

	// Clean up any existing data
	db.Exec("DELETE FROM products")

	return db, cleanup
}

func TestCategoryRepository_NewCategoryRepository(t *testing.T) {
	db, cleanup := setupCategoryTest(t)
	defer cleanup()

	repo := NewCategoryRepository(db)
	assert.NotNil(t, repo)
}

func TestCategoryRepository_GetAllCategories(t *testing.T) {
	db, cleanup := setupCategoryTest(t)
	defer cleanup()

	repo := NewCategoryRepository(db)

	categories := repo.GetAllCategories()

	// Check that all expected categories are returned
	assert.Len(t, categories, 6)
	assert.Contains(t, categories, models.CategoryTops)
	assert.Contains(t, categories, models.CategoryBottoms)
	assert.Contains(t, categories, models.CategoryDresses)
	assert.Contains(t, categories, models.CategoryOuterwear)
	assert.Contains(t, categories, models.CategoryFootwear)
	assert.Contains(t, categories, models.CategoryAccessories)
}

func TestCategoryRepository_GetCategoryStats(t *testing.T) {
	db, cleanup := setupCategoryTest(t)
	defer cleanup()

	repo := NewCategoryRepository(db)
	productRepo := NewProductRepository(db)

	// Create products in different categories
	products := []*models.Product{
		{
			Name:     "Top Product 1",
			Price:    100.00,
			Category: models.CategoryTops,
			Stock:    10,
			Status:   models.StatusAvailable,
			Slug:     "top-product-1",
			SKU:      "top-1",
			Weight:   0.5,
		},
		{
			Name:     "Top Product 2",
			Price:    150.00,
			Category: models.CategoryTops,
			Stock:    10,
			Status:   models.StatusAvailable,
			Slug:     "top-product-2",
			SKU:      "top-2",
			Weight:   0.5,
		},
		{
			Name:     "Bottom Product 1",
			Price:    200.00,
			Category: models.CategoryBottoms,
			Stock:    10,
			Status:   models.StatusAvailable,
			Slug:     "bottom-product-1",
			SKU:      "bottom-1",
			Weight:   0.5,
		},
	}

	for _, p := range products {
		err := productRepo.Create(p)
		require.NoError(t, err)
	}

	// Get category stats
	stats, err := repo.GetCategoryStats()
	require.NoError(t, err)

	// Verify stats
	assert.NotNil(t, stats)

	// Find tops category stats
	var topsCount, bottomsCount int64
	for _, stat := range stats {
		if stat.Category == models.CategoryTops {
			topsCount = stat.ProductCount
		}
		if stat.Category == models.CategoryBottoms {
			bottomsCount = stat.ProductCount
		}
	}

	assert.Equal(t, int64(2), topsCount)
	assert.Equal(t, int64(1), bottomsCount)
}

func TestCategoryRepository_GetCategoryStats_Empty(t *testing.T) {
	db, cleanup := setupCategoryTest(t)
	defer cleanup()

	repo := NewCategoryRepository(db)

	// Get category stats when no products exist
	stats, err := repo.GetCategoryStats()
	require.NoError(t, err)

	// Stats should be empty or contain zero counts
	assert.Empty(t, stats)
}
