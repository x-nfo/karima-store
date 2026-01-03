package repository

import (
	"testing"

	"github.com/karima-store/internal/models"
	"github.com/karima-store/internal/test_setup"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"
)

func setupMediaTest(t *testing.T) (*gorm.DB, *models.Product, func()) {
	db, cleanup := test_setup.SetupTestDB(t)

	// Clean up any existing data
	db.Exec("DELETE FROM media")
	db.Exec("DELETE FROM products")

	// Create a test product for media
	product := &models.Product{
		Name:     "Test Product",
		Price:    100.00,
		Category: models.CategoryTops,
		Stock:    10,
		Status:   models.StatusAvailable,
		Slug:     "test-product",
		SKU:      "test-sku",
		Weight:   0.5,
	}
	db.Create(product)

	return db, product, cleanup
}

func TestMediaRepository_NewMediaRepository(t *testing.T) {
	db, _, cleanup := setupMediaTest(t)
	defer cleanup()

	repo := NewMediaRepository(db)
	assert.NotNil(t, repo)
}

func TestMediaRepository_Create(t *testing.T) {
	db, product, cleanup := setupMediaTest(t)
	defer cleanup()

	repo := NewMediaRepository(db)

	media := &models.Media{
		ProductID: product.ID,
		URL:       "https://example.com/image.jpg",
		Type:      "image",
		IsPrimary: true,
		Position:  0,
	}

	err := repo.Create(media)
	require.NoError(t, err)
	assert.NotZero(t, media.ID)
}

func TestMediaRepository_GetByID(t *testing.T) {
	db, product, cleanup := setupMediaTest(t)
	defer cleanup()

	repo := NewMediaRepository(db)

	// Create media
	media := &models.Media{
		ProductID: product.ID,
		URL:       "https://example.com/image.jpg",
		Type:      "image",
		IsPrimary: true,
		Position:  0,
	}
	err := repo.Create(media)
	require.NoError(t, err)

	// Get by ID
	fetched, err := repo.GetByID(media.ID)
	require.NoError(t, err)
	assert.Equal(t, media.ID, fetched.ID)
	assert.Equal(t, media.URL, fetched.URL)
}

func TestMediaRepository_GetByID_NotFound(t *testing.T) {
	db, _, cleanup := setupMediaTest(t)
	defer cleanup()

	repo := NewMediaRepository(db)

	_, err := repo.GetByID(99999)
	assert.Error(t, err)
	assert.Equal(t, gorm.ErrRecordNotFound, err)
}

func TestMediaRepository_GetAll(t *testing.T) {
	db, product, cleanup := setupMediaTest(t)
	defer cleanup()

	repo := NewMediaRepository(db)

	// Create multiple media items
	for i := 0; i < 3; i++ {
		media := &models.Media{
			ProductID: product.ID,
			URL:       "https://example.com/image" + string(rune('0'+i)) + ".jpg",
			Type:      "image",
			IsPrimary: i == 0,
			Position:  i,
		}
		err := repo.Create(media)
		require.NoError(t, err)
	}

	// Get all
	mediaList, err := repo.GetAll()
	require.NoError(t, err)
	assert.Len(t, mediaList, 3)

	// Verify ordered by position
	for i, m := range mediaList {
		assert.Equal(t, i, m.Position)
	}
}

func TestMediaRepository_GetByProductID(t *testing.T) {
	db, product, cleanup := setupMediaTest(t)
	defer cleanup()

	repo := NewMediaRepository(db)
	productRepo := NewProductRepository(db)

	// Create another product
	product2 := &models.Product{
		Name:     "Test Product 2",
		Price:    200.00,
		Category: models.CategoryBottoms,
		Stock:    10,
		Status:   models.StatusAvailable,
		Slug:     "test-product-2",
		SKU:      "test-sku-2",
		Weight:   0.5,
	}
	err := productRepo.Create(product2)
	require.NoError(t, err)

	// Create media for product 1
	for i := 0; i < 2; i++ {
		media := &models.Media{
			ProductID: product.ID,
			URL:       "https://example.com/p1-image" + string(rune('0'+i)) + ".jpg",
			Type:      "image",
			Position:  i,
		}
		err := repo.Create(media)
		require.NoError(t, err)
	}

	// Create media for product 2
	media := &models.Media{
		ProductID: product2.ID,
		URL:       "https://example.com/p2-image.jpg",
		Type:      "image",
		Position:  0,
	}
	err = repo.Create(media)
	require.NoError(t, err)

	// Get by product ID
	mediaList, err := repo.GetByProductID(product.ID)
	require.NoError(t, err)
	assert.Len(t, mediaList, 2)

	// All should belong to product 1
	for _, m := range mediaList {
		assert.Equal(t, product.ID, m.ProductID)
	}
}

func TestMediaRepository_Update(t *testing.T) {
	db, product, cleanup := setupMediaTest(t)
	defer cleanup()

	repo := NewMediaRepository(db)

	// Create media
	media := &models.Media{
		ProductID: product.ID,
		URL:       "https://example.com/image.jpg",
		Type:      "image",
		IsPrimary: false,
		Position:  0,
	}
	err := repo.Create(media)
	require.NoError(t, err)

	// Update media
	media.URL = "https://example.com/updated-image.jpg"
	media.IsPrimary = true
	err = repo.Update(media)
	require.NoError(t, err)

	// Verify update
	fetched, err := repo.GetByID(media.ID)
	require.NoError(t, err)
	assert.Equal(t, "https://example.com/updated-image.jpg", fetched.URL)
	assert.True(t, fetched.IsPrimary)
}

func TestMediaRepository_Delete(t *testing.T) {
	db, product, cleanup := setupMediaTest(t)
	defer cleanup()

	repo := NewMediaRepository(db)

	// Create media
	media := &models.Media{
		ProductID: product.ID,
		URL:       "https://example.com/image.jpg",
		Type:      "image",
		Position:  0,
	}
	err := repo.Create(media)
	require.NoError(t, err)

	// Delete media
	err = repo.Delete(media.ID)
	require.NoError(t, err)

	// Verify deletion
	_, err = repo.GetByID(media.ID)
	assert.Error(t, err)
}

func TestMediaRepository_SetAsPrimary(t *testing.T) {
	db, product, cleanup := setupMediaTest(t)
	defer cleanup()

	repo := NewMediaRepository(db)

	// Create multiple media items
	media1 := &models.Media{
		ProductID: product.ID,
		URL:       "https://example.com/image1.jpg",
		Type:      "image",
		IsPrimary: true,
		Position:  0,
	}
	err := repo.Create(media1)
	require.NoError(t, err)

	media2 := &models.Media{
		ProductID: product.ID,
		URL:       "https://example.com/image2.jpg",
		Type:      "image",
		IsPrimary: false,
		Position:  1,
	}
	err = repo.Create(media2)
	require.NoError(t, err)

	// Set media2 as primary
	err = repo.SetAsPrimary(media2.ID)
	require.NoError(t, err)

	// Verify media2 is now primary
	fetched2, err := repo.GetByID(media2.ID)
	require.NoError(t, err)
	assert.True(t, fetched2.IsPrimary)

	// Verify media1 is no longer primary
	fetched1, err := repo.GetByID(media1.ID)
	require.NoError(t, err)
	assert.False(t, fetched1.IsPrimary)
}

func TestMediaRepository_UnsetPrimary(t *testing.T) {
	db, product, cleanup := setupMediaTest(t)
	defer cleanup()

	repo := NewMediaRepository(db)

	// Create a primary media item
	media := &models.Media{
		ProductID: product.ID,
		URL:       "https://example.com/image.jpg",
		Type:      "image",
		IsPrimary: true,
		Position:  0,
	}
	err := repo.Create(media)
	require.NoError(t, err)

	// Unset primary for product
	err = repo.UnsetPrimary(product.ID)
	require.NoError(t, err)

	// Verify media is no longer primary
	fetched, err := repo.GetByID(media.ID)
	require.NoError(t, err)
	assert.False(t, fetched.IsPrimary)
}
