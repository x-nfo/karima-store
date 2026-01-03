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

func setupCouponTest(t *testing.T) (*gorm.DB, func()) {
	db, cleanup := test_setup.SetupTestDB(t)

	// Clean up any existing data
	db.Exec("DELETE FROM coupon_usages")
	db.Exec("DELETE FROM coupons")

	return db, cleanup
}

func createTestCoupon(code string) *models.Coupon {
	now := time.Now()
	validUntil := now.Add(24 * time.Hour)
	return &models.Coupon{
		Code:              code,
		Name:              "Test Coupon",
		Description:       "Test Description",
		Type:              models.CouponTypePercentage,
		Status:            models.CouponStatusActive,
		DiscountValue:     10.0,
		MaxDiscount:       50.0,
		MinPurchaseAmount: 100.0,
		MaxUsageCount:     100,
		MaxUsagePerUser:   1,
		ValidFrom:         &now,
		ValidUntil:        &validUntil,
		ForRetail:         true,
		ForReseller:       true,
	}
}

func TestCouponRepository_NewCouponRepository(t *testing.T) {
	db, cleanup := setupCouponTest(t)
	defer cleanup()

	repo := NewCouponRepository(db)
	assert.NotNil(t, repo)
}

func TestCouponRepository_Create(t *testing.T) {
	db, cleanup := setupCouponTest(t)
	defer cleanup()

	repo := NewCouponRepository(db)

	coupon := createTestCoupon("TEST10")

	err := repo.Create(coupon)
	require.NoError(t, err)
	assert.NotZero(t, coupon.ID)
	assert.Equal(t, "TEST10", coupon.Code)
}

func TestCouponRepository_GetByID(t *testing.T) {
	db, cleanup := setupCouponTest(t)
	defer cleanup()

	repo := NewCouponRepository(db)

	// Create a coupon
	coupon := createTestCoupon("TEST10")
	err := repo.Create(coupon)
	require.NoError(t, err)

	// Get by ID
	fetched, err := repo.GetByID(coupon.ID)
	require.NoError(t, err)
	assert.Equal(t, coupon.ID, fetched.ID)
	assert.Equal(t, "TEST10", fetched.Code)
}

func TestCouponRepository_GetByID_NotFound(t *testing.T) {
	db, cleanup := setupCouponTest(t)
	defer cleanup()

	repo := NewCouponRepository(db)

	_, err := repo.GetByID(99999)
	assert.Error(t, err)
	assert.Equal(t, gorm.ErrRecordNotFound, err)
}

func TestCouponRepository_GetByCode(t *testing.T) {
	db, cleanup := setupCouponTest(t)
	defer cleanup()

	repo := NewCouponRepository(db)

	// Create a coupon
	coupon := createTestCoupon("SUMMER2026")
	err := repo.Create(coupon)
	require.NoError(t, err)

	// Get by code
	fetched, err := repo.GetByCode("SUMMER2026")
	require.NoError(t, err)
	assert.Equal(t, coupon.ID, fetched.ID)
	assert.Equal(t, "SUMMER2026", fetched.Code)
}

func TestCouponRepository_GetByCode_NotFound(t *testing.T) {
	db, cleanup := setupCouponTest(t)
	defer cleanup()

	repo := NewCouponRepository(db)

	_, err := repo.GetByCode("NONEXISTENT")
	assert.Error(t, err)
	assert.Equal(t, gorm.ErrRecordNotFound, err)
}

func TestCouponRepository_GetAll(t *testing.T) {
	db, cleanup := setupCouponTest(t)
	defer cleanup()

	repo := NewCouponRepository(db)

	// Create multiple coupons
	for i := 1; i <= 5; i++ {
		coupon := createTestCoupon("COUPON" + string(rune('0'+i)))
		err := repo.Create(coupon)
		require.NoError(t, err)
	}

	// Get all with pagination
	coupons, total, err := repo.GetAll(10, 0)
	require.NoError(t, err)
	assert.Len(t, coupons, 5)
	assert.Equal(t, int64(5), total)
}

func TestCouponRepository_GetAll_Pagination(t *testing.T) {
	db, cleanup := setupCouponTest(t)
	defer cleanup()

	repo := NewCouponRepository(db)

	// Create 10 coupons
	for i := 1; i <= 10; i++ {
		coupon := createTestCoupon("COUPON" + string(rune('A'+i-1)))
		err := repo.Create(coupon)
		require.NoError(t, err)
	}

	// Get first page
	coupons, total, err := repo.GetAll(5, 0)
	require.NoError(t, err)
	assert.Len(t, coupons, 5)
	assert.Equal(t, int64(10), total)

	// Get second page
	coupons2, total2, err := repo.GetAll(5, 5)
	require.NoError(t, err)
	assert.Len(t, coupons2, 5)
	assert.Equal(t, int64(10), total2)
}

func TestCouponRepository_GetActive(t *testing.T) {
	db, cleanup := setupCouponTest(t)
	defer cleanup()

	repo := NewCouponRepository(db)

	// Create active coupon
	activeCoupon := createTestCoupon("ACTIVE")
	err := repo.Create(activeCoupon)
	require.NoError(t, err)

	// Create inactive coupon
	inactiveCoupon := createTestCoupon("INACTIVE")
	inactiveCoupon.Status = models.CouponStatusInactive
	err = repo.Create(inactiveCoupon)
	require.NoError(t, err)

	// Get active coupons
	active, err := repo.GetActive()
	require.NoError(t, err)
	assert.Len(t, active, 1)
	assert.Equal(t, "ACTIVE", active[0].Code)
}

func TestCouponRepository_Update(t *testing.T) {
	db, cleanup := setupCouponTest(t)
	defer cleanup()

	repo := NewCouponRepository(db)

	// Create a coupon
	coupon := createTestCoupon("UPDATE")
	err := repo.Create(coupon)
	require.NoError(t, err)

	// Update the coupon
	coupon.DiscountValue = 20.0
	coupon.Name = "Updated Coupon"
	err = repo.Update(coupon)
	require.NoError(t, err)

	// Verify update
	fetched, err := repo.GetByID(coupon.ID)
	require.NoError(t, err)
	assert.Equal(t, 20.0, fetched.DiscountValue)
	assert.Equal(t, "Updated Coupon", fetched.Name)
}

func TestCouponRepository_Delete(t *testing.T) {
	db, cleanup := setupCouponTest(t)
	defer cleanup()

	repo := NewCouponRepository(db)

	// Create a coupon
	coupon := createTestCoupon("DELETE")
	err := repo.Create(coupon)
	require.NoError(t, err)

	// Delete the coupon
	err = repo.Delete(coupon.ID)
	require.NoError(t, err)

	// Verify deletion (soft delete)
	_, err = repo.GetByID(coupon.ID)
	assert.Error(t, err)
}

func TestCouponRepository_ValidateCoupon_Success(t *testing.T) {
	db, cleanup := setupCouponTest(t)
	defer cleanup()

	repo := NewCouponRepository(db)

	// Create a valid coupon
	coupon := createTestCoupon("VALID")
	err := repo.Create(coupon)
	require.NoError(t, err)

	// Validate coupon
	validated, err := repo.ValidateCoupon("VALID", 1, 150.0, "retail")
	require.NoError(t, err)
	assert.Equal(t, coupon.ID, validated.ID)
}

func TestCouponRepository_ValidateCoupon_InactiveCoupon(t *testing.T) {
	db, cleanup := setupCouponTest(t)
	defer cleanup()

	repo := NewCouponRepository(db)

	// Create an inactive coupon
	coupon := createTestCoupon("INACTIVE")
	coupon.Status = models.CouponStatusInactive
	err := repo.Create(coupon)
	require.NoError(t, err)

	// Validate coupon - should fail
	_, err = repo.ValidateCoupon("INACTIVE", 1, 150.0, "retail")
	assert.Error(t, err)
}

func TestCouponRepository_ValidateCoupon_MinPurchaseNotMet(t *testing.T) {
	db, cleanup := setupCouponTest(t)
	defer cleanup()

	repo := NewCouponRepository(db)

	// Create a coupon with min purchase
	coupon := createTestCoupon("MINPURCHASE")
	coupon.MinPurchaseAmount = 200.0
	err := repo.Create(coupon)
	require.NoError(t, err)

	// Validate coupon with low purchase amount - should fail
	_, err = repo.ValidateCoupon("MINPURCHASE", 1, 100.0, "retail")
	assert.Error(t, err)
}

func TestCouponRepository_ValidateCoupon_CustomerTypeRestriction(t *testing.T) {
	db, cleanup := setupCouponTest(t)
	defer cleanup()

	repo := NewCouponRepository(db)

	// Create a reseller-only coupon
	coupon := createTestCoupon("RESELLERONLY")
	coupon.ForRetail = false
	coupon.ForReseller = true
	err := repo.Create(coupon)
	require.NoError(t, err)

	// Validate for retail customer - should fail
	_, err = repo.ValidateCoupon("RESELLERONLY", 1, 150.0, "retail")
	assert.Error(t, err)

	// Validate for reseller - should succeed
	validated, err := repo.ValidateCoupon("RESELLERONLY", 1, 150.0, "reseller")
	require.NoError(t, err)
	assert.Equal(t, coupon.ID, validated.ID)
}

func TestCouponRepository_ValidateCoupon_MaxUsageReached(t *testing.T) {
	db, cleanup := setupCouponTest(t)
	defer cleanup()

	repo := NewCouponRepository(db)

	// Create a coupon that has reached max usage
	coupon := createTestCoupon("MAXUSED")
	coupon.MaxUsageCount = 5
	coupon.UsageCount = 5
	err := repo.Create(coupon)
	require.NoError(t, err)

	// Validate coupon - should fail
	_, err = repo.ValidateCoupon("MAXUSED", 1, 150.0, "retail")
	assert.Error(t, err)
}

func TestCouponRepository_ValidateCoupon_ExpiredCoupon(t *testing.T) {
	db, cleanup := setupCouponTest(t)
	defer cleanup()

	repo := NewCouponRepository(db)

	// Create an expired coupon
	pastTime := time.Now().Add(-48 * time.Hour)
	pastEndTime := time.Now().Add(-24 * time.Hour)
	coupon := createTestCoupon("EXPIRED")
	coupon.ValidFrom = &pastTime
	coupon.ValidUntil = &pastEndTime
	err := repo.Create(coupon)
	require.NoError(t, err)

	// Validate coupon - should fail
	_, err = repo.ValidateCoupon("EXPIRED", 1, 150.0, "retail")
	assert.Error(t, err)
}

func TestCouponRepository_RecordUsage(t *testing.T) {
	db, cleanup := setupCouponTest(t)
	defer cleanup()

	// Migrate CouponUsage table
	db.AutoMigrate(&models.CouponUsage{})

	repo := NewCouponRepository(db)

	// Create a coupon
	coupon := createTestCoupon("RECORDUSAGE")
	err := repo.Create(coupon)
	require.NoError(t, err)

	// Record usage
	err = repo.RecordUsage(coupon.ID, 1, 100, 10.0)
	require.NoError(t, err)

	// Verify usage count increased
	fetched, err := repo.GetByID(coupon.ID)
	require.NoError(t, err)
	assert.Equal(t, 1, fetched.UsageCount)
}

func TestCouponRepository_GetUserUsageCount(t *testing.T) {
	db, cleanup := setupCouponTest(t)
	defer cleanup()

	// Migrate CouponUsage table
	db.AutoMigrate(&models.CouponUsage{})

	repo := NewCouponRepository(db)

	// Create a coupon
	coupon := createTestCoupon("USERUSAGE")
	err := repo.Create(coupon)
	require.NoError(t, err)

	// Record multiple usages by same user
	err = repo.RecordUsage(coupon.ID, 1, 100, 10.0)
	require.NoError(t, err)
	err = repo.RecordUsage(coupon.ID, 1, 101, 15.0)
	require.NoError(t, err)
	err = repo.RecordUsage(coupon.ID, 2, 102, 20.0) // Different user
	require.NoError(t, err)

	// Get user usage count
	count, err := repo.GetUserUsageCount(coupon.ID, 1)
	require.NoError(t, err)
	assert.Equal(t, 2, count)

	// Get usage count for different user
	count2, err := repo.GetUserUsageCount(coupon.ID, 2)
	require.NoError(t, err)
	assert.Equal(t, 1, count2)
}
