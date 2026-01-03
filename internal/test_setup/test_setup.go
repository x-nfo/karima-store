package test_setup

import (
	"log"
	"os"
	"testing"

	"github.com/karima-store/internal/config"
	"github.com/karima-store/internal/database"
	"github.com/karima-store/internal/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// TestDB is global test database instance
var TestDB *gorm.DB

// TestConfig is global test configuration
var TestCfg *config.Config

// TestMain is main entry point for tests
func TestMain(m *testing.M) {
	// Setup test environment
	setupTestEnvironment()

	// Run tests
	code := m.Run()

	// Cleanup
	cleanupTestEnvironment()

	// Exit with the test exit code
	os.Exit(code)
}

// setupTestEnvironment sets up test environment
func setupTestEnvironment() {
	// Load test configuration
	TestCfg = config.TestConfig()

	// Set up test database
	setupTestDatabase()

	// Set up test Redis (if needed)
	// setupTestRedis()

	log.Println("Test environment setup complete")
}

// cleanupTestEnvironment cleans up test environment
func cleanupTestEnvironment() {
	// Close database connection
	if TestDB != nil {
		sqlDB, err := TestDB.DB()
		if err == nil {
			sqlDB.Close()
		}
	}

	log.Println("Test environment cleanup complete")
}

// setupTestDatabase sets up test database
func setupTestDatabase() {
	// Build DSN for test database
	dsn := buildTestDSN()

	// Open database connection
	var err error
	TestDB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent), // Disable SQL logs in tests
	})
	if err != nil {
		log.Fatalf("Failed to connect to test database: %v", err)
	}

	// Run migrations for test database
	// Run migrations for test database
	if err := RunMigrations(TestDB); err != nil {
		log.Fatalf("Failed to run migrations: %v", err)
	}

	log.Println("Test database setup complete")
}

// buildTestDSN builds DSN for test database
func buildTestDSN() string {
	return "host=" + TestCfg.DBHost +
		" port=" + TestCfg.DBPort +
		" user=" + TestCfg.DBUser +
		" password=" + TestCfg.DBPassword +
		" dbname=" + TestCfg.DBName +
		" sslmode=disable"
}

// CleanupTestData cleans up test data after each test
// CleanupTestData cleans up test data after each test
func CleanupTestData(t *testing.T, db *gorm.DB) {
	if db == nil {
		return
	}

	// Delete all test data
	// This is a simple cleanup - you might want to be more selective
	// based on your test requirements
	// Delete all test data in correct order to handle foreign keys
	tables := []string{
		"coupon_usages",
		"flash_sale_products",
		"order_items",
		"cart_items",
		"stock_logs",
		"media",
		"reviews",
		"wishlists",
		"product_variants",
		"carts",
		"orders",
		"products",
		"coupons",
		"flash_sales",
		"users",
		"shipping_zones",
		"taxes",
	}

	for _, table := range tables {
		err := db.Exec("DELETE FROM " + table).Error
		if err != nil {
			// Don't fail test cleanup, just log it (some tables might not exist yet)
			t.Logf("Note: Failed to clean up table %s: %v", table, err)
		}
	}
}

// GetTestDB returns test database instance
func GetTestDB() *gorm.DB {
	return TestDB
}

// GetTestConfig returns test configuration
func GetTestConfig() *config.Config {
	return TestCfg
}

// SetupTestDB sets up a test database for individual tests
func SetupTestDB(t *testing.T) (*gorm.DB, func()) {
	// Initialize test config if not already done
	if TestCfg == nil {
		TestCfg = config.TestConfig()
	}

	// Create a new database connection for this test
	dsn := buildTestDSN()

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	if err != nil {
		t.Fatalf("Failed to connect to test database: %v", err)
	}

	// Run migrations for this connection to ensure schema exists
	if err := RunMigrations(db); err != nil {
		t.Fatalf("Failed to run migrations: %v", err)
	}

	// Clean up data before test to ensure clean state
	CleanupTestData(t, db)

	// Clean up function
	cleanup := func() {
		// Clean up data after test
		CleanupTestData(t, db)

		sqlDB, err := db.DB()
		if err == nil {
			sqlDB.Close()
		}
	}

	return db, cleanup
}

// SetupTestRedis sets up a test Redis client for individual tests
func SetupTestRedis(t *testing.T) database.RedisClient {
	// Use test config to create Redis instance
	// Use test config to create Redis instance
	testCfg := config.TestConfigWithRedis()

	// Create Redis instance
	redisInstance, err := database.NewRedis(testCfg)
	if err != nil {
		t.Logf("Warning: Redis not available for testing: %v", err)
		// Return a mock or nil - for now we'll return nil and handle in tests
		// In production tests, you'd want to skip tests that require Redis
		return nil
	}

	return redisInstance
}

// RunMigrations runs database migrations
func RunMigrations(db *gorm.DB) error {
	return db.AutoMigrate(
		&models.User{},
		&models.Product{},
		&models.ProductVariant{},
		&models.Media{},
		&models.Order{},
		&models.OrderItem{},
		&models.Cart{},
		&models.CartItem{},
		&models.Review{},
		&models.Wishlist{},
		&models.Coupon{},
		&models.CouponUsage{},
		&models.FlashSale{},
		&models.FlashSaleProduct{},
		&models.ShippingZone{},
		&models.Tax{},
		&models.StockLog{},
	)
}
