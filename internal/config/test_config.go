package config

// TestConfig returns a test configuration for unit tests
func TestConfig() *Config {
	return &Config{
		AppEnv:            getEnv("TEST_APP_ENV", "test"),
		AppPort:           getEnv("TEST_APP_PORT", "8080"),
		DBHost:            getEnv("TEST_DB_HOST", "localhost"),
		DBPort:            getEnv("TEST_DB_PORT", "5432"),
		DBUser:            getEnv("TEST_DB_USER", "test_user"),
		DBPassword:        getEnv("TEST_DB_PASSWORD", "test_password"),
		DBName:            getEnv("TEST_DB_NAME", "karima_store_test"),
		RedisHost:         getEnv("TEST_REDIS_HOST", "localhost"),
		RedisPort:         getEnv("TEST_REDIS_PORT", "6379"),
		RedisPassword:     getEnv("TEST_REDIS_PASSWORD", ""),
		KratosPublicURL:   "http://localhost:4433",
		KratosAdminURL:    "http://localhost:4434",
		FileStorage:       "local",
		RateLimitLimit:    "100",
		RateLimitWindow:   "1m",
		R2AccountID:       "",
		R2AccessKeyID:     "",
		R2SecretAccessKey: "",
		R2BucketName:      "",
		R2PublicURL:       "",
		R2Region:          "",
		JWTSecret:         "test-secret-key-for-testing-only",
	}
}

// TestConfigWithRedis returns a test configuration with Redis enabled
func TestConfigWithRedis() *Config {
	cfg := TestConfig()
	cfg.RedisHost = "localhost"
	cfg.RedisPort = "6380"
	return cfg
}

// TestConfigWithR2 returns a test configuration with R2 storage enabled
func TestConfigWithR2() *Config {
	cfg := TestConfig()
	cfg.FileStorage = "r2"
	cfg.R2AccountID = "test-account-id"
	cfg.R2AccessKeyID = "test-access-key-id"
	cfg.R2SecretAccessKey = "test-secret-access-key"
	cfg.R2BucketName = "test-bucket"
	cfg.R2PublicURL = "https://test.r2.dev"
	cfg.R2Region = "auto"
	return cfg
}

// TestConfigProduction returns a test configuration simulating production environment
func TestConfigProduction() *Config {
	cfg := TestConfig()
	cfg.AppEnv = "production"
	cfg.RateLimitLimit = "120"
	cfg.RateLimitWindow = "1m"
	return cfg
}
