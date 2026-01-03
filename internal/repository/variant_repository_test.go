package repository

import (
	"fmt"
	"testing"

	"github.com/karima-store/internal/models"
	"github.com/karima-store/internal/test_setup"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"
)

func setupVariantTest(t *testing.T) (*gorm.DB, *models.Product, func()) {
	db, cleanup := test_setup.SetupTestDB(t)

	// Clean up existing data
	db.Exec("DELETE FROM product_variants")
	db.Exec("DELETE FROM products")

	// Create a test product
	product := &models.Product{
		Name:     "Test Product",
		Price:    100.00,
		Category: models.CategoryTops,
		Stock:    100,
		Status:   models.StatusAvailable,
		Slug:     "test-product",
		SKU:      "test-product-sku",
		Weight:   0.5,
	}
	db.Create(product)

	return db, product, cleanup
}

func createTestVariant(productID uint, name, sku string) *models.ProductVariant {
	return &models.ProductVariant{
		ProductID: productID,
		Name:      name,
		Size:      "M",
		Color:     "Red",
		Price:     110.00,
		Stock:     10,
		SKU:       sku,
	}
}

func TestVariantRepository_NewVariantRepository(t *testing.T) {
	db, _, cleanup := setupVariantTest(t)
	defer cleanup()

	repo := NewVariantRepository(db)
	assert.NotNil(t, repo)
}

func TestVariantRepository_Create(t *testing.T) {
	db, product, cleanup := setupVariantTest(t)
	defer cleanup()

	repo := NewVariantRepository(db)

	variant := createTestVariant(product.ID, "Medium Red", "VAR-001")
	err := repo.Create(variant)
	require.NoError(t, err)
	assert.NotZero(t, variant.ID)
}

func TestVariantRepository_Create_DuplicateSKU(t *testing.T) {
	db, product, cleanup := setupVariantTest(t)
	defer cleanup()

	repo := NewVariantRepository(db)

	// Create first variant
	variant1 := createTestVariant(product.ID, "Medium Red", "VAR-DUPE")
	err := repo.Create(variant1)
	require.NoError(t, err)

	// Try to create second variant with same SKU
	variant2 := createTestVariant(product.ID, "Large Blue", "VAR-DUPE")
	err = repo.Create(variant2)
	assert.Error(t, err) // Should fail due to unique constraint
}

func TestVariantRepository_GetByID(t *testing.T) {
	db, product, cleanup := setupVariantTest(t)
	defer cleanup()

	repo := NewVariantRepository(db)

	// Create variant
	variant := createTestVariant(product.ID, "Medium Red", "VAR-001")
	err := repo.Create(variant)
	require.NoError(t, err)

	// Get by ID
	fetched, err := repo.GetByID(variant.ID)
	require.NoError(t, err)
	assert.Equal(t, variant.ID, fetched.ID)
	assert.Equal(t, "Medium Red", fetched.Name)
}

func TestVariantRepository_GetByID_NotFound(t *testing.T) {
	db, _, cleanup := setupVariantTest(t)
	defer cleanup()

	repo := NewVariantRepository(db)

	_, err := repo.GetByID(99999)
	assert.Error(t, err)
	assert.Equal(t, gorm.ErrRecordNotFound, err)
}

func TestVariantRepository_GetBySKU(t *testing.T) {
	db, product, cleanup := setupVariantTest(t)
	defer cleanup()

	repo := NewVariantRepository(db)

	// Create variant
	variant := createTestVariant(product.ID, "Medium Red", "SKU-UNIQUE-123")
	err := repo.Create(variant)
	require.NoError(t, err)

	// Get by SKU
	fetched, err := repo.GetBySKU("SKU-UNIQUE-123")
	require.NoError(t, err)
	assert.Equal(t, variant.ID, fetched.ID)
	assert.Equal(t, "SKU-UNIQUE-123", fetched.SKU)
}

func TestVariantRepository_GetBySKU_NotFound(t *testing.T) {
	db, _, cleanup := setupVariantTest(t)
	defer cleanup()

	repo := NewVariantRepository(db)

	_, err := repo.GetBySKU("NON-EXISTENT-SKU")
	assert.Error(t, err)
	assert.Equal(t, gorm.ErrRecordNotFound, err)
}

func TestVariantRepository_GetByProductID(t *testing.T) {
	db, product, cleanup := setupVariantTest(t)
	defer cleanup()

	repo := NewVariantRepository(db)
	productRepo := NewProductRepository(db)

	// Create another product
	product2 := &models.Product{
		Name:     "Test Product 2",
		Price:    200.00,
		Category: models.CategoryBottoms,
		Stock:    50,
		Status:   models.StatusAvailable,
		Slug:     "test-product-2",
		SKU:      "test-product-2-sku",
		Weight:   0.5,
	}
	err := productRepo.Create(product2)
	require.NoError(t, err)

	// Create variants for product 1
	for i := 1; i <= 3; i++ {
		variant := createTestVariant(product.ID, fmt.Sprintf("Variant %d", i), fmt.Sprintf("P1-VAR-%d", i))
		err := repo.Create(variant)
		require.NoError(t, err)
	}

	// Create variants for product 2
	for i := 1; i <= 2; i++ {
		variant := createTestVariant(product2.ID, fmt.Sprintf("Variant %d", i), fmt.Sprintf("P2-VAR-%d", i))
		err := repo.Create(variant)
		require.NoError(t, err)
	}

	// Get variants for product 1
	variants, err := repo.GetByProductID(product.ID)
	require.NoError(t, err)
	assert.Len(t, variants, 3)

	// All should belong to product 1
	for _, v := range variants {
		assert.Equal(t, product.ID, v.ProductID)
	}
}

func TestVariantRepository_Update(t *testing.T) {
	db, product, cleanup := setupVariantTest(t)
	defer cleanup()

	repo := NewVariantRepository(db)

	// Create variant
	variant := createTestVariant(product.ID, "Original Name", "VAR-UPDATE")
	err := repo.Create(variant)
	require.NoError(t, err)

	// Update variant
	variant.Name = "Updated Name"
	variant.Price = 150.00
	variant.Size = "L"
	variant.Color = "Blue"
	err = repo.Update(variant)
	require.NoError(t, err)

	// Verify update
	fetched, err := repo.GetByID(variant.ID)
	require.NoError(t, err)
	assert.Equal(t, "Updated Name", fetched.Name)
	assert.Equal(t, 150.00, fetched.Price)
	assert.Equal(t, "L", fetched.Size)
	assert.Equal(t, "Blue", fetched.Color)
}

func TestVariantRepository_Delete(t *testing.T) {
	db, product, cleanup := setupVariantTest(t)
	defer cleanup()

	repo := NewVariantRepository(db)

	// Create variant
	variant := createTestVariant(product.ID, "To Delete", "VAR-DELETE")
	err := repo.Create(variant)
	require.NoError(t, err)

	// Delete variant
	err = repo.Delete(variant.ID)
	require.NoError(t, err)

	// Verify deletion
	_, err = repo.GetByID(variant.ID)
	assert.Error(t, err)
}

func TestVariantRepository_UpdateStock(t *testing.T) {
	db, product, cleanup := setupVariantTest(t)
	defer cleanup()

	repo := NewVariantRepository(db)

	// Create variant with initial stock
	variant := createTestVariant(product.ID, "Stock Test", "VAR-STOCK")
	variant.Stock = 50
	err := repo.Create(variant)
	require.NoError(t, err)

	// Increase stock
	err = repo.UpdateStock(variant.ID, 10)
	require.NoError(t, err)

	// Verify
	fetched, err := repo.GetByID(variant.ID)
	require.NoError(t, err)
	assert.Equal(t, 60, fetched.Stock)
}

func TestVariantRepository_UpdateStock_Decrease(t *testing.T) {
	db, product, cleanup := setupVariantTest(t)
	defer cleanup()

	repo := NewVariantRepository(db)

	// Create variant with initial stock
	variant := createTestVariant(product.ID, "Stock Test", "VAR-STOCK-DEC")
	variant.Stock = 50
	err := repo.Create(variant)
	require.NoError(t, err)

	// Decrease stock
	err = repo.UpdateStock(variant.ID, -20)
	require.NoError(t, err)

	// Verify
	fetched, err := repo.GetByID(variant.ID)
	require.NoError(t, err)
	assert.Equal(t, 30, fetched.Stock)
}

func TestVariantRepository_MultipleVariants(t *testing.T) {
	db, product, cleanup := setupVariantTest(t)
	defer cleanup()

	repo := NewVariantRepository(db)

	// Create variants with different sizes and colors
	sizes := []string{"S", "M", "L", "XL"}
	colors := []string{"Red", "Blue"}

	for i, size := range sizes {
		for j, color := range colors {
			variant := &models.ProductVariant{
				ProductID: product.ID,
				Name:      fmt.Sprintf("%s - %s", size, color),
				Size:      size,
				Color:     color,
				Price:     100.00 + float64(i*10),
				Stock:     10 * (i + 1),
				SKU:       fmt.Sprintf("VAR-%s-%s", size, color[:1]),
			}
			err := repo.Create(variant)
			require.NoError(t, err, "Failed to create variant %s-%s (%d-%d)", size, color, i, j)
		}
	}

	// Get all variants for product
	variants, err := repo.GetByProductID(product.ID)
	require.NoError(t, err)
	assert.Len(t, variants, 8) // 4 sizes x 2 colors
}

func TestVariantRepository_UpdateVariantPrice(t *testing.T) {
	db, product, cleanup := setupVariantTest(t)
	defer cleanup()

	repo := NewVariantRepository(db)

	// Create variant
	variant := createTestVariant(product.ID, "Price Test", "VAR-PRICE")
	variant.Price = 100.00
	err := repo.Create(variant)
	require.NoError(t, err)

	// Update price
	variant.Price = 120.00
	err = repo.Update(variant)
	require.NoError(t, err)

	// Verify
	fetched, err := repo.GetByID(variant.ID)
	require.NoError(t, err)
	assert.Equal(t, 120.00, fetched.Price)
}
