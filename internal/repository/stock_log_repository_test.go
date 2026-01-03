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

func setupStockLogTest(t *testing.T) (*gorm.DB, *models.Product, func()) {
	// SetupTestDB now calls CleanupTestData implicitly, but we need product first
	db, cleanup := test_setup.SetupTestDB(t)

	// Create test product
	product := &models.Product{
		Name:     "Stock Log Product",
		Price:    100.00,
		Category: models.CategoryTops,
		Stock:    50,
		Status:   models.StatusAvailable,
		Slug:     "stock-log-product",
		SKU:      "STOCK-LOG-SKU",
		Weight:   0.5,
	}
	db.Create(product)

	return db, product, cleanup
}

func TestStockLogRepository_NewStockLogRepository(t *testing.T) {
	db, _, cleanup := setupStockLogTest(t)
	defer cleanup()

	repo := NewStockLogRepository(db)
	assert.NotNil(t, repo)
}

func TestStockLogRepository_Create(t *testing.T) {
	db, product, cleanup := setupStockLogTest(t)
	defer cleanup()

	repo := NewStockLogRepository(db)

	log := &models.StockLog{
		ProductID:     product.ID,
		ChangeAmount:  10,
		PreviousStock: 50,
		NewStock:      60,
		Reason:        "Restock",
		ReferenceID:   "PO-001",
		CreatedAt:     time.Now(),
	}

	err := repo.Create(log)
	require.NoError(t, err)
	assert.NotZero(t, log.ID)
}

func TestStockLogRepository_WithTx(t *testing.T) {
	db, product, cleanup := setupStockLogTest(t)
	defer cleanup()

	repo := NewStockLogRepository(db)

	tx := db.Begin()
	defer tx.Rollback()

	txRepo := repo.WithTx(tx)

	log := &models.StockLog{
		ProductID:     product.ID,
		ChangeAmount:  -5,
		PreviousStock: 60,
		NewStock:      55,
		Reason:        "Order",
		ReferenceID:   "ORD-001",
		CreatedAt:     time.Now(),
	}

	err := txRepo.Create(log)
	require.NoError(t, err)

	tx.Commit()

	// Verify existence
	var count int64
	err = db.Model(&models.StockLog{}).Where("id = ?", log.ID).Count(&count).Error
	require.NoError(t, err)
	assert.Equal(t, int64(1), count)
}
