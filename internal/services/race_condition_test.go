package services_test

import (
	"fmt"
	"sync"
	"testing"

	"github.com/karima-store/internal/database"
	"github.com/karima-store/internal/models"
	"github.com/karima-store/internal/repository"
	"github.com/karima-store/internal/services"
	"github.com/karima-store/internal/test_setup"
	"github.com/stretchr/testify/assert"
)

// MockPricingService satisfies services.PricingService interface
type MockPricingService struct{}

func (m *MockPricingService) CalculatePrice(req services.PriceCalculationRequest) (*services.PriceCalculationResponse, error) {
	return &services.PriceCalculationResponse{
		BasePrice:     10000,
		FinalPrice:    10000,
		OriginalPrice: 10000,
	}, nil
}

func (m *MockPricingService) CalculateOrderSummary(items []services.PriceCalculationRequest, shippingReq services.ShippingCalculationRequest, customerType services.CustomerType) (*services.OrderSummary, error) {
	return &services.OrderSummary{
		Subtotal:     10000,
		ShippingCost: 5000,
		Total:        15000,
	}, nil
}

func (m *MockPricingService) CalculateShippingCost(req services.ShippingCalculationRequest) (*services.ShippingCalculationResponse, error) {
	return &services.ShippingCalculationResponse{}, nil
}
func (m *MockPricingService) CheckFreeShipping(orderAmount float64, regionCode string) (bool, error) {
	return false, nil
}
func (m *MockPricingService) CalculateCouponDiscount(req services.CouponCalculationRequest) (float64, string, error) {
	return 0, "", nil
}
func (m *MockPricingService) ApplyCouponToPriceCalculation(resp *services.PriceCalculationResponse, couponReq services.CouponCalculationRequest) error {
	return nil
}

// MockNotificationService satisfies services.NotificationService interface
type MockNotificationService struct{}

func (m *MockNotificationService) SendWhatsAppMessage(order *models.Order, message string, recipient string) error {
	return nil
}
func (m *MockNotificationService) SendOrderCreatedNotification(order *models.Order) error { return nil }
func (m *MockNotificationService) SendPaymentSuccessNotification(order *models.Order) error {
	return nil
}
func (m *MockNotificationService) SendShippingNotification(order *models.Order, trackingNumber string) error {
	return nil
}
func (m *MockNotificationService) GetWhatsAppStatus() (string, error) { return "connected", nil }
func (m *MockNotificationService) SendTestWhatsAppMessage(phoneNumber string, message string) error {
	return nil
}
func (m *MockNotificationService) ProcessWhatsAppWebhook(data map[string]interface{}) error {
	return nil
}
func (m *MockNotificationService) GetWhatsAppWebhookURL() string { return "" }
func (m *MockNotificationService) GetDB() interface{}            { return nil }

func TestCheckoutConcurrency(t *testing.T) {
	// 1. Setup Data Access Layer (Real DB)
	db, cleanup := test_setup.SetupTestDB(t)
	defer cleanup()

	// Create Postgres wrapper using the test configuration
	// Note: We need to ensure this connects to the SAME database as SetupTestDB
	cfg := test_setup.GetTestConfig()
	pg, err := database.NewPostgreSQL(cfg)
	assert.NoError(t, err)
	defer pg.Close()

	// Repositories
	productRepo := repository.NewProductRepository(db)
	orderRepo := repository.NewOrderRepository(db)
	variantRepo := repository.NewVariantRepository(db)
	stockLogRepo := repository.NewStockLogRepository(db)

	// Services
	pricingService := &MockPricingService{}
	notificationService := &MockNotificationService{}
	midtransConfig := &services.MidtransConfig{
		ServerKey: "dummy_server_key",
		ClientKey: "dummy_client_key",
	}

	checkoutService := services.NewCheckoutService(
		pg,
		orderRepo,
		productRepo,
		variantRepo,
		stockLogRepo,
		pricingService,
		notificationService,
		midtransConfig,
	)

	// 2. Setup Initial Data
	// Create User (Required for Order FK)
	user := &models.User{
		FullName: "Test User",
		Email:    "test@example.com",
		Password: "password123", // required field
		Role:     models.RoleCustomer,
	}
	err = db.Create(user).Error
	assert.NoError(t, err)

	initialStock := 10
	product := &models.Product{
		Name:  "Test Product",
		Slug:  "test-product-concurrency",
		Price: 10000,
		Stock: initialStock,
	}
	err = productRepo.Create(product)
	assert.NoError(t, err)

	// 3. Concurrency Test
	// We will attempt to buy 1 item 20 times concurrently.
	// Only 10 should succeed.
	var wg sync.WaitGroup
	concurrentRequests := 20
	successCount := 0
	failCount := 0
	var mu sync.Mutex

	fmt.Printf("Starting concurrency test: %d requests for %d stock\n", concurrentRequests, initialStock)

	for i := 0; i < concurrentRequests; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			req := &models.CheckoutRequest{
				UserID:        user.ID,
				Items:         []models.CheckoutItem{{ProductID: product.ID, Quantity: 1}},
				PaymentMethod: "bank_transfer",
				ShippingCity:  "Jakarta",
				ShippingName:  "Test User",
			}

			_, err := checkoutService.Checkout(req)

			mu.Lock()
			if err == nil {
				successCount++
			} else {
				failCount++
				// Optional: Log failures to see if they are indeed stock issues
				// if err.Error() != "stock reservation failed: insufficient stock" {
				// 	t.Logf("Unexpected error: %v", err)
				// }
			}
			mu.Unlock()
		}(i)
	}

	wg.Wait()

	// 4. Verification

	// Check final stock
	updatedProduct, err := productRepo.GetByID(product.ID)
	assert.NoError(t, err)

	fmt.Printf("Results - Success: %d, Fail: %d, Final Stock: %d\n", successCount, failCount, updatedProduct.Stock)

	// Assertions
	assert.Equal(t, 0, updatedProduct.Stock, "Stock should be exactly 0")
	assert.Equal(t, initialStock, successCount, fmt.Sprintf("Should have exactly %d successful orders", initialStock))
	assert.Equal(t, concurrentRequests-initialStock, failCount, fmt.Sprintf("Should have exactly %d failed orders", concurrentRequests-initialStock))

	// Verify Stock Logs
	// Since we are mocking other things, we just check if logs were created for successful transactions
	// Note: StockLogRepository might track this.
}
