package repository

import (
	"testing"

	"github.com/karima-store/internal/models"
	"github.com/karima-store/internal/test_setup"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"
)

func setupOrderTest(t *testing.T) (*gorm.DB, *models.User, *models.Product, func()) {
	db, cleanup := test_setup.SetupTestDB(t)

	// Clean up existing data
	db.Exec("DELETE FROM order_items")
	db.Exec("DELETE FROM orders")
	db.Exec("DELETE FROM products")
	db.Exec("DELETE FROM users")

	// Create test user
	user := &models.User{
		FullName: "Test User",
		Email:    "test@example.com",
		Password: "hashedpassword",
		Role:     models.RoleCustomer,
	}
	db.Create(user)

	// Create test product
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

	return db, user, product, cleanup
}

func createTestOrder(userID uint, orderNumber string) *models.Order {
	return &models.Order{
		UserID:             userID,
		OrderNumber:        orderNumber,
		Status:             models.StatusPending,
		PaymentStatus:      models.PaymentPending,
		Subtotal:           100.00,
		ShippingCost:       10.00,
		TotalAmount:        110.00,
		ShippingName:       "Test User",
		ShippingPhone:      "08123456789",
		ShippingAddress:    "123 Test Street",
		ShippingCity:       "Jakarta",
		ShippingProvince:   "DKI Jakarta",
		ShippingPostalCode: "12345",
	}
}

func TestOrderRepository_NewOrderRepository(t *testing.T) {
	db, _, _, cleanup := setupOrderTest(t)
	defer cleanup()

	repo := NewOrderRepository(db)
	assert.NotNil(t, repo)
}

func TestOrderRepository_Create(t *testing.T) {
	db, user, _, cleanup := setupOrderTest(t)
	defer cleanup()

	repo := NewOrderRepository(db)

	order := createTestOrder(user.ID, "ORD-001")
	err := repo.Create(order)
	require.NoError(t, err)
	assert.NotZero(t, order.ID)
}

func TestOrderRepository_GetByID(t *testing.T) {
	db, user, _, cleanup := setupOrderTest(t)
	defer cleanup()

	repo := NewOrderRepository(db)

	// Create order
	order := createTestOrder(user.ID, "ORD-001")
	err := repo.Create(order)
	require.NoError(t, err)

	// Get by ID
	fetched, err := repo.GetByID(order.ID)
	require.NoError(t, err)
	assert.Equal(t, order.ID, fetched.ID)
	assert.Equal(t, "ORD-001", fetched.OrderNumber)
}

func TestOrderRepository_GetByID_NotFound(t *testing.T) {
	db, _, _, cleanup := setupOrderTest(t)
	defer cleanup()

	repo := NewOrderRepository(db)

	_, err := repo.GetByID(99999)
	assert.Error(t, err)
	assert.Equal(t, gorm.ErrRecordNotFound, err)
}

func TestOrderRepository_GetByOrderNumber(t *testing.T) {
	db, user, _, cleanup := setupOrderTest(t)
	defer cleanup()

	repo := NewOrderRepository(db)

	// Create order
	order := createTestOrder(user.ID, "ORD-UNIQUE-123")
	err := repo.Create(order)
	require.NoError(t, err)

	// Get by order number
	fetched, err := repo.GetByOrderNumber("ORD-UNIQUE-123")
	require.NoError(t, err)
	assert.Equal(t, order.ID, fetched.ID)
	assert.Equal(t, "ORD-UNIQUE-123", fetched.OrderNumber)
}

func TestOrderRepository_GetByOrderNumber_NotFound(t *testing.T) {
	db, _, _, cleanup := setupOrderTest(t)
	defer cleanup()

	repo := NewOrderRepository(db)

	_, err := repo.GetByOrderNumber("NON-EXISTENT")
	assert.Error(t, err)
	assert.Equal(t, gorm.ErrRecordNotFound, err)
}

func TestOrderRepository_GetByUserID(t *testing.T) {
	db, user, _, cleanup := setupOrderTest(t)
	defer cleanup()

	repo := NewOrderRepository(db)

	// Create multiple orders for user
	for i := 1; i <= 5; i++ {
		order := createTestOrder(user.ID, "ORD-00"+string(rune('0'+i)))
		err := repo.Create(order)
		require.NoError(t, err)
	}

	// Get orders by user ID
	orders, total, err := repo.GetByUserID(user.ID, 10, 0)
	require.NoError(t, err)
	assert.Len(t, orders, 5)
	assert.Equal(t, int64(5), total)
}

func TestOrderRepository_GetByUserID_Pagination(t *testing.T) {
	db, user, _, cleanup := setupOrderTest(t)
	defer cleanup()

	repo := NewOrderRepository(db)

	// Create 10 orders
	for i := 1; i <= 10; i++ {
		order := createTestOrder(user.ID, "ORD-"+string(rune('A'+i-1)))
		err := repo.Create(order)
		require.NoError(t, err)
	}

	// Get first page
	orders, total, err := repo.GetByUserID(user.ID, 5, 0)
	require.NoError(t, err)
	assert.Len(t, orders, 5)
	assert.Equal(t, int64(10), total)

	// Get second page
	orders2, total2, err := repo.GetByUserID(user.ID, 5, 5)
	require.NoError(t, err)
	assert.Len(t, orders2, 5)
	assert.Equal(t, int64(10), total2)
}

func TestOrderRepository_Update(t *testing.T) {
	db, user, _, cleanup := setupOrderTest(t)
	defer cleanup()

	repo := NewOrderRepository(db)

	// Create order
	order := createTestOrder(user.ID, "ORD-UPDATE")
	err := repo.Create(order)
	require.NoError(t, err)

	// Update order
	order.Status = models.StatusProcessing
	order.CustomerNotes = "Updated notes"
	err = repo.Update(order)
	require.NoError(t, err)

	// Verify update
	fetched, err := repo.GetByID(order.ID)
	require.NoError(t, err)
	assert.Equal(t, models.StatusProcessing, fetched.Status)
	assert.Equal(t, "Updated notes", fetched.CustomerNotes)
}

func TestOrderRepository_UpdateStatus(t *testing.T) {
	db, user, _, cleanup := setupOrderTest(t)
	defer cleanup()

	repo := NewOrderRepository(db)

	// Create order
	order := createTestOrder(user.ID, "ORD-STATUS")
	err := repo.Create(order)
	require.NoError(t, err)

	// Update status
	err = repo.UpdateStatus(order.ID, models.StatusShipped)
	require.NoError(t, err)

	// Verify
	fetched, err := repo.GetByID(order.ID)
	require.NoError(t, err)
	assert.Equal(t, models.StatusShipped, fetched.Status)
}

func TestOrderRepository_UpdatePaymentStatus(t *testing.T) {
	db, user, _, cleanup := setupOrderTest(t)
	defer cleanup()

	repo := NewOrderRepository(db)

	// Create order
	order := createTestOrder(user.ID, "ORD-PAYMENT")
	err := repo.Create(order)
	require.NoError(t, err)

	// Update payment status
	err = repo.UpdatePaymentStatus(order.ID, models.PaymentPaid)
	require.NoError(t, err)

	// Verify
	fetched, err := repo.GetByID(order.ID)
	require.NoError(t, err)
	assert.Equal(t, models.PaymentPaid, fetched.PaymentStatus)
}

func TestOrderRepository_Delete(t *testing.T) {
	db, user, _, cleanup := setupOrderTest(t)
	defer cleanup()

	repo := NewOrderRepository(db)

	// Create order
	order := createTestOrder(user.ID, "ORD-DELETE")
	err := repo.Create(order)
	require.NoError(t, err)

	// Delete order
	err = repo.Delete(order.ID)
	require.NoError(t, err)

	// Verify deletion
	_, err = repo.GetByID(order.ID)
	assert.Error(t, err)
}

func TestOrderRepository_WithTx(t *testing.T) {
	db, user, _, cleanup := setupOrderTest(t)
	defer cleanup()

	repo := NewOrderRepository(db)

	// Start transaction
	tx := db.Begin()
	defer tx.Rollback()

	// Get repository with transaction
	txRepo := repo.WithTx(tx)
	assert.NotNil(t, txRepo)

	// Create order within transaction
	order := createTestOrder(user.ID, "ORD-TX")
	err := txRepo.Create(order)
	require.NoError(t, err)

	// Commit transaction
	tx.Commit()

	// Verify order exists
	fetched, err := repo.GetByID(order.ID)
	require.NoError(t, err)
	assert.Equal(t, "ORD-TX", fetched.OrderNumber)
}
