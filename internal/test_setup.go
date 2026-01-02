package internal

import (
	"log"
	"os"
	"testing"

	"github.com/karima-store/internal/config"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// TestDB is the global test database instance
var TestDB *gorm.DB

// TestConfig is the global test configuration
var TestCfg *config.Config

// TestMain is the main entry point for tests
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

// setupTestEnvironment sets up the test environment
func setupTestEnvironment() {
	// Load test configuration
	TestCfg = config.TestConfig()

	// Set up test database
	setupTestDatabase()

	// Set up test Redis (if needed)
	// setupTestRedis()

	log.Println("Test environment setup complete")
}

// cleanupTestEnvironment cleans up the test environment
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

// setupTestDatabase sets up the test database
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
	err = TestDB.AutoMigrate()
	if err != nil {
		log.Fatalf("Failed to run migrations: %v", err)
	}

	log.Println("Test database setup complete")
}

// buildTestDSN builds the DSN for the test database
func buildTestDSN() string {
	return "host=" + TestCfg.DBHost +
		" port=" + TestCfg.DBPort +
		" user=" + TestCfg.DBUser +
		" password=" + TestCfg.DBPassword +
		" dbname=" + TestCfg.DBName +
		" sslmode=disable"
}

// CleanupTestData cleans up test data after each test
func CleanupTestData(t *testing.T) {
	if TestDB == nil {
		return
	}

	// Delete all test data
	// This is a simple cleanup - you might want to be more selective
	// based on your test requirements
	tables := []string{
		"media",
		"variants",
		"products",
		"categories",
		"orders",
		"coupons",
		"flash_sales",
		"reviews",
		"shipping_zones",
		"stock_logs",
		"taxes",
		"users",
		"wishlists",
		"carts",
	}

	for _, table := range tables {
		err := TestDB.Exec("DELETE FROM " + table).Error
		if err != nil {
			t.Logf("Failed to clean up table %s: %v", table, err)
		}
	}
}

// GetTestDB returns the test database instance
func GetTestDB() *gorm.DB {
	return TestDB
}

// GetTestConfig returns the test configuration
func GetTestConfig() *config.Config {
	return TestCfg
}