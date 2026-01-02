package test_setup

import (
	"context"
	"log"
	"os"
	"testing"

	"github.com/karima-store/internal/config"
	"github.com/redis/go-redis/v9"
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
	err = TestDB.AutoMigrate()
	if err != nil {
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
	// Create a new database connection for this test
	dsn := buildTestDSN()

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	if err != nil {
		t.Fatalf("Failed to connect to test database: %v", err)
	}

	// Clean up function
	cleanup := func() {
		sqlDB, err := db.DB()
		if err == nil {
			sqlDB.Close()
		}
	}

	return db, cleanup
}

// SetupTestRedis sets up a test Redis client for individual tests
func SetupTestRedis(t *testing.T) *redis.Client {
	client := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       1, // Use DB 1 for tests
	})

	// Test connection
	ctx := context.Background()
	_, err := client.Ping(ctx).Result()
	if err != nil {
		t.Logf("Warning: Redis not available for testing: %v", err)
	}

	return client
}
