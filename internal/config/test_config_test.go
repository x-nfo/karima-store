package config

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTestConfig(t *testing.T) {
	// Unset TEST_DB_NAME to ensure we get the default value
	os.Unsetenv("TEST_DB_NAME")
	defer os.Unsetenv("TEST_DB_NAME")

	cfg := TestConfig()

	assert.Equal(t, "test", cfg.AppEnv)
	assert.Equal(t, "8080", cfg.AppPort)
	assert.Equal(t, "localhost", cfg.DBHost)
	assert.Equal(t, "5432", cfg.DBPort)
	assert.Equal(t, "karima_store", cfg.DBUser)
	assert.Equal(t, "lokal", cfg.DBPassword)
	assert.Equal(t, "karima_test_db", cfg.DBName)
	assert.Equal(t, "disable", cfg.DBSSLMode)
	assert.Equal(t, "localhost", cfg.RedisHost)
	assert.Equal(t, "6380", cfg.RedisPort)
	assert.Equal(t, "", cfg.RedisPassword)
	assert.Equal(t, "test-secret-key-for-testing-only", cfg.JWTSecret)
	assert.Equal(t, "test-google-key", cfg.GoogleKey)
	assert.Equal(t, "test-google-secret", cfg.GoogleSecret)
	assert.Equal(t, "local", cfg.FileStorage)
	assert.Equal(t, "100", cfg.RateLimitLimit)
	assert.Equal(t, "1m", cfg.RateLimitWindow)
}

func TestTestConfigWithEnvVars(t *testing.T) {
	// Set environment variables
	os.Setenv("TEST_DB_HOST", "test-host")
	os.Setenv("TEST_DB_PORT", "5433")
	os.Setenv("TEST_DB_USER", "test-user")
	os.Setenv("TEST_DB_PASSWORD", "test-pass")
	os.Setenv("TEST_DB_NAME", "test-db")
	os.Setenv("TEST_REDIS_HOST", "test-redis")
	os.Setenv("TEST_REDIS_PORT", "6379")
	defer func() {
		os.Unsetenv("TEST_DB_HOST")
		os.Unsetenv("TEST_DB_PORT")
		os.Unsetenv("TEST_DB_USER")
		os.Unsetenv("TEST_DB_PASSWORD")
		os.Unsetenv("TEST_DB_NAME")
		os.Unsetenv("TEST_REDIS_HOST")
		os.Unsetenv("TEST_REDIS_PORT")
	}()

	cfg := TestConfig()

	assert.Equal(t, "test-host", cfg.DBHost)
	assert.Equal(t, "5433", cfg.DBPort)
	assert.Equal(t, "test-user", cfg.DBUser)
	assert.Equal(t, "test-pass", cfg.DBPassword)
	assert.Equal(t, "test-db", cfg.DBName)
	assert.Equal(t, "test-redis", cfg.RedisHost)
	assert.Equal(t, "6379", cfg.RedisPort)
}

func TestTestConfigWithRedis(t *testing.T) {
	cfg := TestConfigWithRedis()

	assert.Equal(t, "test", cfg.AppEnv)
	assert.Equal(t, "localhost", cfg.RedisHost)
	assert.Equal(t, "6380", cfg.RedisPort)
}

func TestTestConfigWithR2(t *testing.T) {
	cfg := TestConfigWithR2()

	assert.Equal(t, "test", cfg.AppEnv)
	assert.Equal(t, "r2", cfg.FileStorage)
	assert.Equal(t, "test-account-id", cfg.R2AccountID)
	assert.Equal(t, "test-access-key-id", cfg.R2AccessKeyID)
	assert.Equal(t, "test-secret-access-key", cfg.R2SecretAccessKey)
	assert.Equal(t, "test-bucket", cfg.R2BucketName)
	assert.Equal(t, "https://test.r2.dev", cfg.R2PublicURL)
	assert.Equal(t, "auto", cfg.R2Region)
}

func TestTestConfigProduction(t *testing.T) {
	cfg := TestConfigProduction()

	assert.Equal(t, "production", cfg.AppEnv)
	assert.Equal(t, "120", cfg.RateLimitLimit)
	assert.Equal(t, "1m", cfg.RateLimitWindow)
}
