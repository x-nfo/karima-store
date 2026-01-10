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
// Note: For tests, we use AutoMigrate to create tables based on GORM models
// This ensures test isolation and doesn't depend on SQL migration files
func RunMigrations(db *gorm.DB) error {
	// Drop existing tables to ensure clean state for tests
	// This is safe for test databases
	tables := []interface{}{
		&models.StockLog{},
		&models.Tax{},
		&models.ShippingZone{},
		&models.FlashSale{},
		&models.CouponUsage{},
		&models.Coupon{},
		&models.Wishlist{},
		&models.Review{},
		&models.CartItem{},
		&models.Cart{},
		&models.OrderItem{},
		&models.Order{},
		&models.Media{},
		&models.ProductVariant{},
		&models.Product{},
		&models.User{},
	}

	// Drop tables in reverse order of dependencies
	for i := len(tables) - 1; i >= 0; i-- {
		if err := db.Migrator().DropTable(tables[i]); err != nil {
			// Ignore errors if tables don't exist
			continue
		}
	}

	// AutoMigrate will create tables with proper schema
	// We need to disable foreign key constraints temporarily for AutoMigrate
	// to work correctly with many-to-many relationships
	db.Exec("SET CONSTRAINTS ALL DEFERRED")

	if err := db.AutoMigrate(tables...); err != nil {
		return err
	}

	// Re-enable foreign key constraints
	db.Exec("SET CONSTRAINTS ALL IMMEDIATE")

	// Manually create flash_sale_products table with correct schema
	// This is a workaround for GORM AutoMigrate not creating the table correctly
	// Drop it first if it exists, then recreate with all required columns
	db.Exec("DROP TABLE IF EXISTS flash_sale_products CASCADE")

	db.Exec(`
		CREATE TABLE flash_sale_products (
			id BIGSERIAL PRIMARY KEY,
			created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
			updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
			flash_sale_id BIGINT NOT NULL,
			product_id BIGINT NOT NULL,
			flash_sale_price DECIMAL(10, 2) NOT NULL,
			flash_sale_stock INTEGER NOT NULL,
			sold_count INTEGER NOT NULL DEFAULT 0,
			CONSTRAINT fk_flash_sale_products_flash_sale FOREIGN KEY (flash_sale_id) REFERENCES flash_sales(id) ON DELETE CASCADE,
			CONSTRAINT fk_flash_sale_products_product FOREIGN KEY (product_id) REFERENCES products(id) ON DELETE CASCADE,
			CONSTRAINT uq_flash_sale_products UNIQUE (flash_sale_id, product_id)
		)
	`)

	// Create indexes
	db.Exec("CREATE INDEX IF NOT EXISTS idx_flash_sale_products_flash_sale_id ON flash_sale_products(flash_sale_id)")
	db.Exec("CREATE INDEX IF NOT EXISTS idx_flash_sale_products_product_id ON flash_sale_products(product_id)")

	return nil
}



