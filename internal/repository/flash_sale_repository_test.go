package repository

import (
	"testing"
	"time"

	"github.com/karima-store/internal/models"
	"github.com/karima-store/internal/test_setup"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"
)

func setupFlashSaleTest(t *testing.T) (*gorm.DB, *models.Product, func()) {
	db, cleanup := test_setup.SetupTestDB(t)

	// Migrate flash sale tables
	db.AutoMigrate(&models.FlashSale{}, &models.FlashSaleProduct{})

	// Clean up any existing data
	db.Exec("DELETE FROM flash_sale_products")
	db.Exec("DELETE FROM flash_sales")
	db.Exec("DELETE FROM products")

	// Create a test product
	product := &models.Product{
		Name:     "Flash Sale Product",
		Price:    100.00,
		Category: models.CategoryTops,
		Stock:    10,
		Status:   models.StatusAvailable,
		Slug:     "flash-sale-product",
		SKU:      "flash-sku",
		Weight:   0.5,
	}
	db.Create(product)

	return db, product, cleanup
}

func createTestFlashSale(name string, status models.FlashSaleStatus) *models.FlashSale {
	now := time.Now()
	return &models.FlashSale{
		Name:               name,
		Description:        "Test flash sale",
		Status:             status,
		StartTime:          now.Add(-1 * time.Hour),
		EndTime:            now.Add(24 * time.Hour),
		DiscountPercentage: 20.0,
		MaxQuantityPerUser: 5,
		TotalStockLimit:    100,
	}
}

func TestFlashSaleRepository_NewFlashSaleRepository(t *testing.T) {
	db, _, cleanup := setupFlashSaleTest(t)
	defer cleanup()

	repo := NewFlashSaleRepository(db)
	assert.NotNil(t, repo)
}

func TestFlashSaleRepository_Create(t *testing.T) {
	db, _, cleanup := setupFlashSaleTest(t)
	defer cleanup()

	repo := NewFlashSaleRepository(db)

	flashSale := createTestFlashSale("Test Sale", models.FlashSaleActive)
	err := repo.Create(flashSale)
	require.NoError(t, err)
	assert.NotZero(t, flashSale.ID)
}

func TestFlashSaleRepository_GetByID(t *testing.T) {
	db, _, cleanup := setupFlashSaleTest(t)
	defer cleanup()

	repo := NewFlashSaleRepository(db)

	// Create flash sale
	flashSale := createTestFlashSale("Test Sale", models.FlashSaleActive)
	err := repo.Create(flashSale)
	require.NoError(t, err)

	// Get by ID
	fetched, err := repo.GetByID(flashSale.ID)
	require.NoError(t, err)
	assert.Equal(t, flashSale.ID, fetched.ID)
	assert.Equal(t, "Test Sale", fetched.Name)
}

func TestFlashSaleRepository_GetByID_NotFound(t *testing.T) {
	db, _, cleanup := setupFlashSaleTest(t)
	defer cleanup()

	repo := NewFlashSaleRepository(db)

	_, err := repo.GetByID(99999)
	assert.Error(t, err)
	assert.Equal(t, gorm.ErrRecordNotFound, err)
}

func TestFlashSaleRepository_GetAll(t *testing.T) {
	db, _, cleanup := setupFlashSaleTest(t)
	defer cleanup()

	repo := NewFlashSaleRepository(db)

	// Create multiple flash sales
	for i := 1; i <= 3; i++ {
		flashSale := createTestFlashSale("Sale "+string(rune('0'+i)), models.FlashSaleActive)
		err := repo.Create(flashSale)
		require.NoError(t, err)
	}

	// Get all
	flashSales, err := repo.GetAll()
	require.NoError(t, err)
	assert.Len(t, flashSales, 3)
}

func TestFlashSaleRepository_GetActiveFlashSales(t *testing.T) {
	db, _, cleanup := setupFlashSaleTest(t)
	defer cleanup()

	repo := NewFlashSaleRepository(db)

	// Create active flash sale (current time within range)
	now := time.Now()
	activeSale := &models.FlashSale{
		Name:               "Active Sale",
		Status:             models.FlashSaleActive,
		StartTime:          now.Add(-1 * time.Hour),
		EndTime:            now.Add(2 * time.Hour),
		DiscountPercentage: 20.0,
	}
	err := repo.Create(activeSale)
	require.NoError(t, err)

	// Create ended flash sale
	endedSale := &models.FlashSale{
		Name:               "Ended Sale",
		Status:             models.FlashSaleEnded,
		StartTime:          now.Add(-48 * time.Hour),
		EndTime:            now.Add(-24 * time.Hour),
		DiscountPercentage: 15.0,
	}
	err = repo.Create(endedSale)
	require.NoError(t, err)

	// Get active flash sales
	active, err := repo.GetActiveFlashSales()
	require.NoError(t, err)
	assert.Len(t, active, 1)
	assert.Equal(t, "Active Sale", active[0].Name)
}

func TestFlashSaleRepository_GetUpcomingFlashSales(t *testing.T) {
	db, _, cleanup := setupFlashSaleTest(t)
	defer cleanup()

	repo := NewFlashSaleRepository(db)

	now := time.Now()

	// Create upcoming flash sale
	upcomingSale := &models.FlashSale{
		Name:               "Upcoming Sale",
		Status:             models.FlashSaleUpcoming,
		StartTime:          now.Add(24 * time.Hour),
		EndTime:            now.Add(48 * time.Hour),
		DiscountPercentage: 25.0,
	}
	err := repo.Create(upcomingSale)
	require.NoError(t, err)

	// Create active flash sale
	activeSale := createTestFlashSale("Active Sale", models.FlashSaleActive)
	err = repo.Create(activeSale)
	require.NoError(t, err)

	// Get upcoming
	upcoming, err := repo.GetUpcomingFlashSales()
	require.NoError(t, err)
	assert.Len(t, upcoming, 1)
	assert.Equal(t, "Upcoming Sale", upcoming[0].Name)
}

func TestFlashSaleRepository_Update(t *testing.T) {
	db, _, cleanup := setupFlashSaleTest(t)
	defer cleanup()

	repo := NewFlashSaleRepository(db)

	// Create flash sale
	flashSale := createTestFlashSale("Original Name", models.FlashSaleActive)
	err := repo.Create(flashSale)
	require.NoError(t, err)

	// Update
	flashSale.Name = "Updated Name"
	flashSale.DiscountPercentage = 30.0
	err = repo.Update(flashSale)
	require.NoError(t, err)

	// Verify
	fetched, err := repo.GetByID(flashSale.ID)
	require.NoError(t, err)
	assert.Equal(t, "Updated Name", fetched.Name)
	assert.Equal(t, 30.0, fetched.DiscountPercentage)
}

func TestFlashSaleRepository_Delete(t *testing.T) {
	db, _, cleanup := setupFlashSaleTest(t)
	defer cleanup()

	repo := NewFlashSaleRepository(db)

	// Create flash sale
	flashSale := createTestFlashSale("To Delete", models.FlashSaleActive)
	err := repo.Create(flashSale)
	require.NoError(t, err)

	// Delete
	err = repo.Delete(flashSale.ID)
	require.NoError(t, err)

	// Verify deletion (soft delete)
	_, err = repo.GetByID(flashSale.ID)
	assert.Error(t, err)
}

func TestFlashSaleRepository_AddProductToFlashSale(t *testing.T) {
	db, product, cleanup := setupFlashSaleTest(t)
	defer cleanup()

	repo := NewFlashSaleRepository(db)

	// Create flash sale
	flashSale := createTestFlashSale("Product Sale", models.FlashSaleActive)
	err := repo.Create(flashSale)
	require.NoError(t, err)

	// Add product to flash sale
	flashSaleProduct := &models.FlashSaleProduct{
		FlashSaleID:    flashSale.ID,
		ProductID:      product.ID,
		FlashSalePrice: 80.0,
		FlashSaleStock: 50,
	}
	err = repo.AddProductToFlashSale(flashSaleProduct)
	require.NoError(t, err)
	assert.NotZero(t, flashSaleProduct.ID)
}

func TestFlashSaleRepository_GetFlashSaleProducts(t *testing.T) {
	db, product, cleanup := setupFlashSaleTest(t)
	defer cleanup()

	repo := NewFlashSaleRepository(db)

	// Create flash sale
	flashSale := createTestFlashSale("Product Sale", models.FlashSaleActive)
	err := repo.Create(flashSale)
	require.NoError(t, err)

	// Add product to flash sale
	flashSaleProduct := &models.FlashSaleProduct{
		FlashSaleID:    flashSale.ID,
		ProductID:      product.ID,
		FlashSalePrice: 80.0,
		FlashSaleStock: 50,
	}
	err = repo.AddProductToFlashSale(flashSaleProduct)
	require.NoError(t, err)

	// Get flash sale products
	products, err := repo.GetFlashSaleProducts(flashSale.ID)
	require.NoError(t, err)
	assert.Len(t, products, 1)
	assert.Equal(t, 80.0, products[0].FlashSalePrice)
}

func TestFlashSaleRepository_RemoveProductFromFlashSale(t *testing.T) {
	db, product, cleanup := setupFlashSaleTest(t)
	defer cleanup()

	repo := NewFlashSaleRepository(db)

	// Create flash sale
	flashSale := createTestFlashSale("Product Sale", models.FlashSaleActive)
	err := repo.Create(flashSale)
	require.NoError(t, err)

	// Add product
	flashSaleProduct := &models.FlashSaleProduct{
		FlashSaleID:    flashSale.ID,
		ProductID:      product.ID,
		FlashSalePrice: 80.0,
		FlashSaleStock: 50,
	}
	err = repo.AddProductToFlashSale(flashSaleProduct)
	require.NoError(t, err)

	// Remove product
	err = repo.RemoveProductFromFlashSale(flashSale.ID, product.ID)
	require.NoError(t, err)

	// Verify removal
	products, err := repo.GetFlashSaleProducts(flashSale.ID)
	require.NoError(t, err)
	assert.Len(t, products, 0)
}

func TestFlashSaleRepository_UpdateFlashSaleProduct(t *testing.T) {
	db, product, cleanup := setupFlashSaleTest(t)
	defer cleanup()

	repo := NewFlashSaleRepository(db)

	// Create flash sale
	flashSale := createTestFlashSale("Product Sale", models.FlashSaleActive)
	err := repo.Create(flashSale)
	require.NoError(t, err)

	// Add product
	flashSaleProduct := &models.FlashSaleProduct{
		FlashSaleID:    flashSale.ID,
		ProductID:      product.ID,
		FlashSalePrice: 80.0,
		FlashSaleStock: 50,
	}
	err = repo.AddProductToFlashSale(flashSaleProduct)
	require.NoError(t, err)

	// Update flash sale product
	flashSaleProduct.FlashSalePrice = 70.0
	flashSaleProduct.SoldCount = 10
	err = repo.UpdateFlashSaleProduct(flashSaleProduct)
	require.NoError(t, err)

	// Verify update
	products, err := repo.GetFlashSaleProducts(flashSale.ID)
	require.NoError(t, err)
	assert.Len(t, products, 1)
	assert.Equal(t, 70.0, products[0].FlashSalePrice)
	assert.Equal(t, 10, products[0].SoldCount)
}
