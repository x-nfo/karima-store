package services

import (
	"testing"
	"time"

	"github.com/karima-store/internal/models"
	"github.com/karima-store/internal/repository"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"gorm.io/gorm"
)

// MockFlashSaleRepository
type MockFlashSaleRepository struct {
	mock.Mock
}

func (m *MockFlashSaleRepository) GetByID(id uint) (*models.FlashSale, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.FlashSale), args.Error(1)
}
func (m *MockFlashSaleRepository) GetAll() ([]models.FlashSale, error) {
	args := m.Called()
	return args.Get(0).([]models.FlashSale), args.Error(1)
}
func (m *MockFlashSaleRepository) GetActiveFlashSales() ([]models.FlashSale, error) {
	args := m.Called()
	return args.Get(0).([]models.FlashSale), args.Error(1)
}
func (m *MockFlashSaleRepository) GetUpcomingFlashSales() ([]models.FlashSale, error) {
	args := m.Called()
	return args.Get(0).([]models.FlashSale), args.Error(1)
}
func (m *MockFlashSaleRepository) Create(flashSale *models.FlashSale) error {
	args := m.Called(flashSale)
	return args.Error(0)
}
func (m *MockFlashSaleRepository) Update(flashSale *models.FlashSale) error {
	args := m.Called(flashSale)
	return args.Error(0)
}
func (m *MockFlashSaleRepository) Delete(id uint) error {
	args := m.Called(id)
	return args.Error(0)
}
func (m *MockFlashSaleRepository) AddProductToFlashSale(flashSaleProduct *models.FlashSaleProduct) error {
	args := m.Called(flashSaleProduct)
	return args.Error(0)
}
func (m *MockFlashSaleRepository) RemoveProductFromFlashSale(flashSaleID, productID uint) error {
	args := m.Called(flashSaleID, productID)
	return args.Error(0)
}
func (m *MockFlashSaleRepository) GetFlashSaleProducts(flashSaleID uint) ([]models.FlashSaleProduct, error) {
	args := m.Called(flashSaleID)
	return args.Get(0).([]models.FlashSaleProduct), args.Error(1)
}
func (m *MockFlashSaleRepository) UpdateFlashSaleProduct(flashSaleProduct *models.FlashSaleProduct) error {
	args := m.Called(flashSaleProduct)
	return args.Error(0)
}

// MockCouponRepository
type MockCouponRepository struct {
	mock.Mock
}

func (m *MockCouponRepository) Create(coupon *models.Coupon) error { return m.Called(coupon).Error(0) }
func (m *MockCouponRepository) GetByID(id uint) (*models.Coupon, error) {
	args := m.Called(id)
	return nil, args.Error(1)
} // Simplified
func (m *MockCouponRepository) GetByCode(code string) (*models.Coupon, error) {
	args := m.Called(code)
	return nil, args.Error(1)
}
func (m *MockCouponRepository) GetAll(limit, offset int) ([]models.Coupon, int64, error) {
	args := m.Called(limit, offset)
	return nil, 0, args.Error(2)
}
func (m *MockCouponRepository) GetActive() ([]models.Coupon, error) {
	args := m.Called()
	return nil, args.Error(1)
}
func (m *MockCouponRepository) Update(coupon *models.Coupon) error { return m.Called(coupon).Error(0) }
func (m *MockCouponRepository) Delete(id uint) error               { return m.Called(id).Error(0) }
func (m *MockCouponRepository) ValidateCoupon(code string, userID uint, purchaseAmount float64, customerType string) (*models.Coupon, error) {
	args := m.Called(code, userID, purchaseAmount, customerType)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Coupon), args.Error(1)
}
func (m *MockCouponRepository) RecordUsage(couponID, userID, orderID uint, discountAmount float64) error {
	return m.Called(couponID, userID, orderID, discountAmount).Error(0)
}
func (m *MockCouponRepository) GetUserUsageCount(couponID, userID uint) (int, error) {
	args := m.Called(couponID, userID)
	return args.Int(0), args.Error(1)
}

// MockShippingZoneRepository
type MockShippingZoneRepository struct {
	mock.Mock
}

func (m *MockShippingZoneRepository) Create(zone *models.ShippingZone) error {
	return m.Called(zone).Error(0)
}
func (m *MockShippingZoneRepository) GetByID(id uint) (*models.ShippingZone, error) {
	args := m.Called(id)
	return nil, args.Error(1)
}
func (m *MockShippingZoneRepository) GetByRegion(regionCode string) (*models.ShippingZone, error) {
	args := m.Called(regionCode)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.ShippingZone), args.Error(1)
}
func (m *MockShippingZoneRepository) GetAll() ([]models.ShippingZone, error) {
	args := m.Called()
	return nil, args.Error(1)
}
func (m *MockShippingZoneRepository) GetActive() ([]models.ShippingZone, error) {
	args := m.Called()
	return nil, args.Error(1)
}
func (m *MockShippingZoneRepository) Update(zone *models.ShippingZone) error {
	return m.Called(zone).Error(0)
}
func (m *MockShippingZoneRepository) Delete(id uint) error { return m.Called(id).Error(0) }
func (m *MockShippingZoneRepository) WithTx(tx *gorm.DB) repository.ShippingZoneRepository {
	args := m.Called(tx)
	return args.Get(0).(repository.ShippingZoneRepository)
}

// MockProductRepo (Assume it's in the same package services from product_service_test.go)
// MockVariantRepo (Assume it's in the same package services from product_service_test.go)

func TestPricingService_CalculatePrice(t *testing.T) {
	mockProductRepo := new(MockProductRepository)
	mockVariantRepo := new(MockVariantRepository)
	mockFlashSaleRepo := new(MockFlashSaleRepository)
	mockCouponRepo := new(MockCouponRepository)
	mockZoneRepo := new(MockShippingZoneRepository)

	service := NewPricingService(mockProductRepo, mockVariantRepo, mockFlashSaleRepo, mockCouponRepo, mockZoneRepo)

	// Test 1: Basic Retail Price (No discount)
	t.Run("Basic Retail", func(t *testing.T) {
		req := PriceCalculationRequest{
			ProductID:    1,
			Quantity:     1,
			CustomerType: CustomerRetail,
		}
		product := &models.Product{ID: 1, Price: 100000}

		mockProductRepo.On("GetByID", uint(1)).Return(product, nil).Once()
		mockFlashSaleRepo.On("GetActiveFlashSales").Return([]models.FlashSale{}, nil).Once()

		resp, err := service.CalculatePrice(req)
		assert.NoError(t, err)
		assert.Equal(t, 100000.0, resp.FinalPrice)
		assert.Equal(t, 0.0, resp.Savings)
	})

	// Test 2: Reseller Tiering (Quantity 20 -> 20% discount)
	t.Run("Reseller Discount", func(t *testing.T) {
		req := PriceCalculationRequest{
			ProductID:    1,
			Quantity:     20,
			CustomerType: CustomerReseller,
		}
		product := &models.Product{ID: 1, Price: 100000}

		mockProductRepo.On("GetByID", uint(1)).Return(product, nil).Once()
		mockFlashSaleRepo.On("GetActiveFlashSales").Return([]models.FlashSale{}, nil).Once()

		resp, err := service.CalculatePrice(req)
		assert.NoError(t, err)

		// Base: 100,000 * 20 = 2,000,000
		// Discount 20% -> 0.20 * 100,000 = 20,000 per item
		// Final per item: 80,000
		// Total: 1,600,000

		expectedFinal := 1600000.0
		assert.Equal(t, expectedFinal, resp.FinalPrice)
		assert.Equal(t, "reseller", resp.DiscountType)
	})

	// Test 3: Flash Sale
	t.Run("Flash Sale", func(t *testing.T) {
		req := PriceCalculationRequest{ProductID: 2, Quantity: 1, CustomerType: CustomerRetail}
		product := &models.Product{ID: 2, Price: 100000}

		flashSale := models.FlashSale{
			ID:        1,
			Status:    models.FlashSaleActive,
			StartTime: time.Now().Add(-1 * time.Hour),
			EndTime:   time.Now().Add(1 * time.Hour),
			Products:  []models.Product{{ID: 2}},
		}
		fsProduct := models.FlashSaleProduct{ProductID: 2, FlashSalePrice: 50000}

		mockProductRepo.On("GetByID", uint(2)).Return(product, nil).Once()
		mockFlashSaleRepo.On("GetActiveFlashSales").Return([]models.FlashSale{flashSale}, nil).Once()
		mockFlashSaleRepo.On("GetFlashSaleProducts", uint(1)).Return([]models.FlashSaleProduct{fsProduct}, nil).Once()

		resp, err := service.CalculatePrice(req)
		assert.NoError(t, err)
		assert.Equal(t, 50000.0, resp.FinalPrice)
		assert.Equal(t, "flash_sale", resp.DiscountType)
		assert.True(t, resp.FlashSaleActive)
	})
}
