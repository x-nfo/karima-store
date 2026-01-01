package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
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
)

func main() {
	// Load configuration
	cfg := config.Load()

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
			return c.Status(code).JSON(fiber.Map{
				"error": err.Error(),
			})
		},
	})

	// Global middleware
	app.Use(logger.New(logger.Config{
		Format:     "[${time}] ${status} - ${method} ${path} (${latency})\n",
		TimeFormat: "2006-01-02 15:04:05",
		TimeZone:   "Asia/Jakarta",
	}))
	app.Use(recover.New())
	app.Use(middleware.CORS(cfg.CORSOrigin))

	// Initialize auth middleware
	authMiddleware := middleware.NewAuthMiddleware(cfg.JWTSecret)

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
	checkoutService := services.NewCheckoutService(db, orderRepo, productRepo, variantRepo, stockLogRepo, pricingService, notificationService, midtransConfig)

	// Initialize handlers
	productHandler := handlers.NewProductHandler(productService, mediaService)
	variantHandler := handlers.NewVariantHandler(variantService)
	categoryHandler := handlers.NewCategoryHandler(categoryService)
	pricingHandler := handlers.NewPricingHandler(pricingService, redis)
	mediaHandler := handlers.NewMediaHandler(mediaService)
	checkoutHandler := handlers.NewCheckoutHandler(checkoutService)
	komerceHandler := handlers.NewKomerceHandler(komerceService)
	orderHandler := handlers.NewOrderHandler(orderService)
	whatsappHandler := handlers.NewWhatsAppHandler(notificationService)
	swaggerHandler := handlers.NewSwaggerHandler()

	// Setup routes
	routes.RegisterRoutes(app, authMiddleware.Authenticate(), productHandler, variantHandler, categoryHandler, pricingHandler, mediaHandler, checkoutHandler, komerceHandler, orderHandler, whatsappHandler, swaggerHandler)

	// Health check endpoint
	app.Get("/health", func(c *fiber.Ctx) error {
		// Check database connection
		if err := checkDatabase(db); err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"status":  "error",
				"message": "Database connection failed",
				"details": err.Error(),
			})
		}

		// Check Redis connection
		if err := checkRedis(redis); err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"status":  "error",
				"message": "Redis connection failed",
				"details": err.Error(),
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
	if err := app.Listen(fmt.Sprintf(":%s", port)); err != nil {
		log.Fatalf("Failed to start server: %v", err)
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

	return sqlDB.PingContext(ctx)
}

func checkRedis(redis *database.Redis) error {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	// Try to set and get a value
	err := redis.Set(ctx, "health-check", "ok", 10*time.Second)
	if err != nil {
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
