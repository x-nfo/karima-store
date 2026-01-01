package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/helmet"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/karima-store/internal/config"
	"github.com/karima-store/internal/database"
	"github.com/karima-store/internal/handlers"
	"github.com/karima-store/internal/komerce"
	"github.com/karima-store/internal/middleware"
	"github.com/karima-store/internal/repository"
	"github.com/karima-store/internal/routes"
	"github.com/karima-store/internal/services"
	"github.com/karima-store/internal/utils"
)

func main() {
	// Load configuration
	cfg := config.Load()

	// Validate critical configuration (production only, warnings only in dev)
	validateConfiguration(cfg)

	// Initialize database connections
	db, err := database.NewPostgreSQL(cfg)
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}
	defer db.Close()

	redis, err := database.NewRedis(cfg)
	if err != nil {
		log.Fatalf("Failed to initialize Redis: %v", err)
	}
	defer redis.Close()

	// Run database migrations using golang-migrate
	if err := db.Migrate("migrations"); err != nil {
		log.Fatalf("Failed to migrate database: %v", err)
	}

	// Initialize Fiber app
	app := fiber.New(fiber.Config{
		BodyLimit: 2 * 1024 * 1024, // 2MB
		ErrorHandler: func(c *fiber.Ctx, err error) error {
			code := fiber.StatusInternalServerError
			if e, ok := err.(*fiber.Error); ok {
				code = e.Code
			}
			return utils.SendError(c, code, err.Error(), nil)
		},
	})

	// Global middleware
	app.Use(logger.New(logger.Config{
		Format:     "[${time}] ${status} - ${method} ${path} (${latency})\\n",
		TimeFormat: "2006-01-02 15:04:05",
		TimeZone:   "Asia/Jakarta",
	}))
	app.Use(recover.New())
	app.Use(helmet.New()) // Security Headers
	app.Use(middleware.CORS(cfg.CORSOrigin))
	app.Use(middleware.NewRateLimiter(cfg))

	// Initialize Ory Kratos middleware for authentication
	authMiddleware := middleware.NewKratosMiddleware(cfg.KratosPublicURL, cfg.KratosAdminURL)

	// Initialize repositories
	productRepo := repository.NewProductRepository(db.DB())
	variantRepo := repository.NewVariantRepository(db.DB())
	categoryRepo := repository.NewCategoryRepository(db.DB())
	flashSaleRepo := repository.NewFlashSaleRepository(db.DB())
	couponRepo := repository.NewCouponRepository(db.DB())
	shippingZoneRepo := repository.NewShippingZoneRepository(db.DB())
	mediaRepo := repository.NewMediaRepository(db.DB())
	orderRepo := repository.NewOrderRepository(db.DB())
	stockLogRepo := repository.NewStockLogRepository(db.DB())

	// Initialize services
	productService := services.NewProductService(productRepo, variantRepo, redis)
	orderService := services.NewOrderService(orderRepo) // Added OrderService
	variantService := services.NewVariantService(variantRepo, productRepo)
	categoryService := services.NewCategoryService(categoryRepo)
	pricingService := services.NewPricingService(productRepo, variantRepo, flashSaleRepo, couponRepo, shippingZoneRepo)
	mediaService := services.NewMediaService(mediaRepo, productRepo, cfg)
	notificationService := services.NewNotificationService(db, redis, cfg)

	// Initialize Komerce client
	komerceClient := komerce.NewClient(cfg.KomerceAPIKey, cfg.KomerceBaseURL)
	komerceService := services.NewKomerceService(komerceClient)

	// Midtrans configuration
	midtransConfig := &services.MidtransConfig{
		ServerKey:    getEnv("MIDTRANS_SERVER_KEY", ""),
		ClientKey:    getEnv("MIDTRANS_CLIENT_KEY", ""),
		APIBaseURL:   getEnv("MIDTRANS_API_BASE_URL", "https://app.sandbox.midtrans.com/snap/v1"),
		IsProduction: getEnvAsBool("MIDTRANS_IS_PRODUCTION", false),
	}

	// Initialize checkout service with all dependencies
	checkoutService := services.NewCheckoutService(
		db,
		orderRepo,
		productRepo,
		variantRepo,
		stockLogRepo,
		pricingService,
		notificationService,
		midtransConfig,
	)

	// Initialize handlers
	productHandler := handlers.NewProductHandler(productService, mediaService)
	variantHandler := handlers.NewVariantHandler(variantService)
	categoryHandler := handlers.NewCategoryHandler(categoryService)
	pricingHandler := handlers.NewPricingHandler(pricingService, redis)
	mediaHandler := handlers.NewMediaHandler(mediaService)
	checkoutHandler := handlers.NewCheckoutHandler(checkoutService)
	komerceHandler := handlers.NewKomerceHandler(komerceService)
	orderHandler := handlers.NewOrderHandler(orderService) // Added OrderHandler
	whatsappHandler := handlers.NewWhatsAppHandler(notificationService)
	swaggerHandler := handlers.NewSwaggerHandler()

	// Register routes
	routes.RegisterRoutes(
		app,
		authMiddleware,
		productHandler,
		variantHandler,
		categoryHandler,
		pricingHandler,
		mediaHandler,
		checkoutHandler,
		komerceHandler,
		orderHandler,
		whatsappHandler,
		swaggerHandler,
	)

	// Health check endpoint
	app.Get("/health", func(c *fiber.Ctx) error {
		// Check database
		if err := checkDatabase(db); err != nil {
			return c.Status(fiber.StatusServiceUnavailable).JSON(fiber.Map{
				"status":   "error",
				"database": err.Error(),
			})
		}

		// Check Redis
		if err := checkRedis(redis); err != nil {
			return c.Status(fiber.StatusServiceUnavailable).JSON(fiber.Map{
				"status": "error",
				"redis":  err.Error(),
			})
		}

		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"status":      "ok",
			"database":    "connected",
			"redis":       "connected",
			"environment": cfg.AppEnv,
			"timestamp":   time.Now().Format(time.RFC3339),
		})
	})

	// Start server
	port := cfg.AppPort
	if port == "" {
		port = "8080"
	}

	log.Printf("Starting Karima Store API on port %s (Environment: %s)", port, cfg.AppEnv)

	// Graceful shutdown setup (production) or direct start (development)
	if cfg.AppEnv == "production" {
		// Production: Enable graceful shutdown for clean resource cleanup
		log.Println("Production mode: Graceful shutdown enabled")
		startServerWithGracefulShutdown(app, port, cfg, db, redis)
	} else {
		// Development: Simple start (Ctrl+C works immediately, faster iteration)
		log.Println("Development mode: Press Ctrl+C to stop server (immediate)")
		if err := app.Listen(fmt.Sprintf(":%s", port)); err != nil {
			log.Fatalf("Failed to start server: %v", err)
		}
	}
}

func checkDatabase(db *database.PostgreSQL) error {
	dbClient := db.DB()
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	// Try to execute a simple query
	sqlDB, err := dbClient.DB()
	if err != nil {
		return fmt.Errorf("failed to get sql.DB: %w", err)
	}

	if err := sqlDB.PingContext(ctx); err != nil {
		return fmt.Errorf("database ping failed: %w", err)
	}

	return nil
}

func checkRedis(redis *database.Redis) error {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	// Try to set and get a value
	if err := redis.Set(ctx, "health-check", "ok", 10*time.Second); err != nil {
		return fmt.Errorf("redis set failed: %w", err)
	}

	val, err := redis.Get(ctx, "health-check")
	if err != nil {
		return fmt.Errorf("redis get failed: %w", err)
	}

	if val != "ok" {
		return fmt.Errorf("redis value mismatch")
	}

	return nil
}

// validateConfiguration validates critical configuration based on environment
func validateConfiguration(cfg *config.Config) {
	if cfg.AppEnv == "production" {
		// Strict validation for production
		var errors []string

		if cfg.DBPassword == "" || cfg.DBPassword == "secret" {
			errors = append(errors, "DB_PASSWORD must be set and not default")
		}
		if cfg.RedisPassword == "" {
			log.Println("⚠️  WARNING: REDIS_PASSWORD not set (recommended for production)")
		}
		// Remove JWT_SECRET validation since we're using Kratos only

		if len(errors) > 0 {
			log.Println("❌ Production configuration errors:")
			for _, err := range errors {
				log.Printf("  - %s", err)
			}
			log.Fatal("Fix configuration errors before running in production")
		}

		log.Println("✅ Production configuration validated")
	} else {
		// Development: Just warnings
		log.Println("ℹ️  Development mode: Using relaxed configuration validation")
		if cfg.DBPassword == "secret" || cfg.DBPassword == "lokal" {
			log.Println("⚠️  Using default database password (OK for development)")
		}
	}
}

// Helper functions for environment variables
func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}

func getEnvAsBool(key string, defaultValue bool) bool {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	boolValue, err := strconv.ParseBool(value)
	if err != nil {
		return defaultValue
	}
	return boolValue
}
