package middleware

import (
	"log"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/limiter"
	"github.com/gofiber/storage/redis/v3"
	"github.com/karima-store/internal/config"
)

// NewRateLimiter creates a new rate limiting middleware backed by Redis
func NewRateLimiter(cfg *config.Config) fiber.Handler {
	port, _ := strconv.Atoi(cfg.RedisPort)
	if port == 0 {
		port = 6379 // Default fallback
	}

	// Initialize Redis storage
	// We use a separate connection pool for the rate limiter to avoid contention with main app logic
	store := redis.New(redis.Config{
		Host:     cfg.RedisHost,
		Port:     port,
		Password: cfg.RedisPassword,
		Database: 0, // Use default DB or maybe separated one? Keep 0 for simplicity now
		Reset:    false,
	})

	// Default limits
	max := 60
	expiration := 1 * time.Minute

	// Adjust based on environment
	if cfg.AppEnv == "production" {
		max = 120 // 120 req/min (2 req/sec per IP)
	} else {
		max = 2400 // 2400 req/min for development (high enough not to be annoying)
	}

	// Override from config if present
	if cfg.RateLimitLimit != "" {
		if val, err := strconv.Atoi(cfg.RateLimitLimit); err == nil && val > 0 {
			max = val
		}
	}

	// Try parsing window if set (e.g. "1m", "1h")
	if cfg.RateLimitWindow != "" {
		if val, err := time.ParseDuration(cfg.RateLimitWindow); err == nil {
			expiration = val
		}
	}

	log.Printf("🛡️  Rate Limiter initialized: %d req / %s (Redis backend)", max, expiration)

	return limiter.New(limiter.Config{
		Max:        max,
		Expiration: expiration,
		Storage:    store,
		KeyGenerator: func(c *fiber.Ctx) string {
			// Allow overriding key for testing purposes in non-production environments
			if cfg.AppEnv != "production" {
				if testKey := c.Get("X-Test-Key"); testKey != "" {
					return testKey
				}
			}
			return c.IP()
		},
		LimitReached: func(c *fiber.Ctx) error {
			return c.Status(fiber.StatusTooManyRequests).JSON(fiber.Map{
				"status":  "error",
				"message": "Too many requests, please try again later.",
			})
		},
		// Use sliding window for better accuracy (optional, depends on need, fixed window is cheaper)
		// LimiterMiddleware: limiter.SlidingWindow{},
	})
}
