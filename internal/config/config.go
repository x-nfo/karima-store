package config

import (
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
	"github.com/karima-store/internal/logger"
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

	// JWT
	JWTSecret string

	// Ory Kratos
	KratosPublicURL string
	KratosAdminURL  string
	KratosUIURL     string
}

func Load() *Config {
	// Load .env file
	if err := godotenv.Load(); err != nil {
		if logger.Log != nil {
			logger.Log.Warnw("Environment file not found, using environment variables")
		} else {
			log.Println("Warning: .env file not found, using environment variables")
		}
	}

	return &Config{
		// Database
		DBHost:     getEnv("DB_HOST", "localhost"),
		DBPort:     getEnv("DB_PORT", "5432"),
		DBUser:     getEnv("DB_USER", "postgres"),
		DBPassword: getEnv("DB_PASSWORD", ""),
		DBName:     getEnv("DB_NAME", "karima_db"),
		DBSSLMode:  getEnv("DB_SSL_MODE", "disable"),

		// Redis
		RedisHost:     getEnv("REDIS_HOST", "localhost"),
		RedisPort:     getEnv("REDIS_PORT", "6380"),
		RedisPassword: getEnv("REDIS_PASSWORD", ""),

		// Application
		AppPort: getEnv("APP_PORT", "8080"),
		AppEnv:  getEnv("APP_ENV", "development"),

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

		// JWT
		JWTSecret: getEnv("JWT_SECRET", ""),

		// Ory Kratos
		KratosPublicURL: getEnv("KRATOS_PUBLIC_URL", "http://127.0.0.1:4433"),
		KratosAdminURL:  getEnv("KRATOS_ADMIN_URL", "http://127.0.0.1:4434"),
		KratosUIURL:     getEnv("KRATOS_UI_URL", "http://127.0.0.1:4455"),
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

// Validate validates critical configuration and terminates the application if required variables are missing
func (c *Config) Validate() {
	var errors []string

	// Critical variables that must always be set
	if c.DBPassword == "" {
		errors = append(errors, "DB_PASSWORD is required but not set")
	}

	// Validate based on environment
	if c.AppEnv == "production" {
		// Production-specific strict validation
		if c.JWTSecret == "" {
			errors = append(errors, "JWT_SECRET is required in production")
		}

		// Validate CORS_ORIGIN is set to specific domains, not wildcard
		if c.CORSOrigin == "*" || c.CORSOrigin == "" {
			errors = append(errors, "CORS_ORIGIN must be set to specific domains in production, not wildcard")
		}

		// Validate Redis password in production
		if c.RedisPassword == "" {
			if logger.Log != nil {
				logger.Log.Warnw("Redis password not set (recommended for production)")
			} else {
				log.Println("⚠️  WARNING: REDIS_PASSWORD not set (recommended for production)")
			}
		}
	}

	// If there are any errors, log them and exit
	if len(errors) > 0 {
		if logger.Log != nil {
			logger.Log.Errorw("Configuration validation failed", "errors", errors)
			logger.Log.Fatal("Application cannot start due to missing required configuration")
		} else {
			log.Println("❌ Configuration validation failed:")
			for _, err := range errors {
				log.Printf("  - %s", err)
			}
			log.Fatal("Application cannot start due to missing required configuration")
		}
	}

	if logger.Log != nil {
		logger.Log.Info("Configuration validated successfully")
	} else {
		log.Println("✅ Configuration validated successfully")
	}
}

// InitLogger initializes the structured logger based on configuration
func (c *Config) InitLogger() error {
	cfg := &logger.Config{
		Level:      c.LogLevel,
		Output:     c.LogFile,
		Format:     "json",
		Env:        c.AppEnv,
		WithCaller: true,
	}

	// Use console format in development
	if c.AppEnv == "development" {
		cfg.Format = "console"
	}

	// Use stdout if log file is not specified
	if c.LogFile == "" {
		cfg.Output = "stdout"
	}

	return logger.Init(cfg)
}
