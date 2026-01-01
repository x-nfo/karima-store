package config

import (
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	// Database
	DBHost     string
	DBPort     string
	DBUser     string
	DBPassword string
	DBName     string
	DBSSLMode  string

	// Redis
	RedisHost     string
	RedisPort     string
	RedisPassword string

	// Application
	AppPort string
	AppEnv  string

	// JWT
	JWTSecret     string
	JWTExpiration string

	// File Storage
	FileStorage       string
	FileUploadMaxSize string

	// Cloudflare R2
	R2AccountID       string
	R2AccessKeyID     string
	R2SecretAccessKey string
	R2BucketName      string
	R2PublicURL       string
	R2Region          string

	// Email
	EmailHost     string
	EmailPort     string
	EmailUser     string
	EmailPassword string

	// Logging
	LogLevel string
	LogFile  string

	// Cache
	CacheType     string
	CacheDuration string

	// Rate Limiting
	RateLimitWindow string
	RateLimitLimit  string

	// CORS
	CORSOrigin string

	// API Version
	APIVersion string

	// Go Environment
	GoEnv     string
	GoVersion string

	// Migrations
	MigrationSource string

	// RajaOngkir
	RajaOngkirAPIKey  string
	RajaOngkirBaseURL string

	// Komerce
	KomerceAPIKey  string
	KomerceBaseURL string

	// Fonnte
	FonnteToken string
	FonnteURL   string
}

func Load() *Config {
	// Load .env file
	if err := godotenv.Load(); err != nil {
		log.Println("Warning: .env file not found, using environment variables")
	}

	return &Config{
		// Database
		DBHost:     getEnv("DB_HOST", "localhost"),
		DBPort:     getEnv("DB_PORT", "5432"),
		DBUser:     getEnv("DB_USER", "postgres"),
		DBPassword: getEnv("DB_PASSWORD", "secret"),
		DBName:     getEnv("DB_NAME", "karima_db"),
		DBSSLMode:  getEnv("DB_SSL_MODE", "disable"),

		// Redis
		RedisHost:     getEnv("REDIS_HOST", "localhost"),
		RedisPort:     getEnv("REDIS_PORT", "6380"),
		RedisPassword: getEnv("REDIS_PASSWORD", ""),

		// Application
		AppPort: getEnv("APP_PORT", "8080"),
		AppEnv:  getEnv("APP_ENV", "development"),

		// JWT
		JWTSecret:     getEnv("JWT_SECRET", "super_secret_key"),
		JWTExpiration: getEnv("JWT_EXPIRATION", "24h"),

		// File Storage
		FileStorage:       getEnv("FILE_STORAGE", "local"),
		FileUploadMaxSize: getEnv("FILE_UPLOAD_MAX_SIZE", "10MB"),

		// Cloudflare R2
		R2AccountID:       getEnv("R2_ACCOUNT_ID", ""),
		R2AccessKeyID:     getEnv("R2_ACCESS_KEY_ID", ""),
		R2SecretAccessKey: getEnv("R2_SECRET_ACCESS_KEY", ""),
		R2BucketName:      getEnv("R2_BUCKET_NAME", ""),
		R2PublicURL:       getEnv("R2_PUBLIC_URL", ""),
		R2Region:          getEnv("R2_REGION", "auto"),

		// Email
		EmailHost:     getEnv("EMAIL_HOST", "smtp.gmail.com"),
		EmailPort:     getEnv("EMAIL_PORT", "587"),
		EmailUser:     getEnv("EMAIL_USER", ""),
		EmailPassword: getEnv("EMAIL_PASSWORD", ""),

		// Logging
		LogLevel: getEnv("LOG_LEVEL", "info"),
		LogFile:  getEnv("LOG_FILE", "logs/app.log"),

		// Cache
		CacheType:     getEnv("CACHE_TYPE", "redis"),
		CacheDuration: getEnv("CACHE_DURATION", "1h"),

		// Rate Limiting
		RateLimitWindow: getEnv("RATE_LIMIT_WINDOW", "1m"),
		RateLimitLimit:  getEnv("RATE_LIMIT_LIMIT", "100"),

		// CORS
		CORSOrigin: getEnv("CORS_ORIGIN", "http://localhost:3000"),

		// API Version
		APIVersion: getEnv("API_VERSION", "v1"),

		// Go Environment
		GoEnv:     getEnv("GO_ENV", "development"),
		GoVersion: getEnv("GO_VERSION", "1.21"),

		// Migrations
		MigrationSource: getEnv("MIGRATION_SOURCE", "migrations"),

		// RajaOngkir by Komerce
		RajaOngkirAPIKey:  getEnv("RAJAONKIR_API_KEY_SHIPPING_DELIVERY", ""),
		RajaOngkirBaseURL: getEnv("RAJAONGKIR_BASE_URL", "https://api-sandbox.collaborator.komerce.id"),

		// Komerce
		KomerceAPIKey:  getEnv("KOMERCE_API_KEY", ""),
		KomerceBaseURL: getEnv("KOMERCE_BASE_URL", "https://api-sandbox.collaborator.komerce.id"),

		// Fonnte
		FonnteToken: getEnv("FONNTE_TOKEN", ""),
		FonnteURL:   getEnv("FONNTE_URL", "https://api.fonnte.com/send"),
	}
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func getEnvAsInt(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		if intVal, err := strconv.Atoi(value); err == nil {
			return intVal
		}
	}
	return defaultValue
}

func getEnvAsBool(key string, defaultValue bool) bool {
	if value := os.Getenv(key); value != "" {
		if boolVal, err := strconv.ParseBool(value); err == nil {
			return boolVal
		}
	}
	return defaultValue
}
