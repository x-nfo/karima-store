package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/storage/redis/v3"
	"github.com/karima-store/internal/config"
	"github.com/stretchr/testify/assert"
)

func TestRateLimiter_NewRateLimiter(t *testing.T) {
	// Test rate limiter creation with different environments
	tests := []struct {
		name          string
		env           string
		expectedLimit int
	}{
		{
			name:          "Production environment",
			env:           "production",
			expectedLimit: 120,
		},
		{
			name:          "Development environment",
			env:           "development",
			expectedLimit: 2400,
		},
		{
			name:          "Test environment",
			env:           "test",
			expectedLimit: 2400,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfg := &config.Config{
				AppEnv:       tt.env,
				RedisHost:    "localhost",
				RedisPort:    "6379",
				RateLimitLimit: "",
				RateLimitWindow: "",
			}

			// Create rate limiter
			app := fiber.New()
			app.Use(NewRateLimiter(cfg))

			// Test route
			app.Get("/test", func(c *fiber.Ctx) error {
				return c.SendStatus(fiber.StatusOK)
			})

			// Test first request (should succeed)
			req1 := httptest.NewRequest("GET", "/test", nil)
			resp1 := httptest.NewRecorder()
			app.Handler()(resp1, req1)
			assert.Equal(t, fiber.StatusOK, resp1.Code)

			// Test second request (should succeed)
			req2 := httptest.NewRequest("GET", "/test", nil)
			resp2 := httptest.NewRecorder()
			app.Handler()(resp2, req2)
			assert.Equal(t, fiber.StatusOK, resp2.Code)
		})
	}
}

func TestRateLimiter_CustomConfiguration(t *testing.T) {
	// Test custom rate limit configuration
	cfg := &config.Config{
		AppEnv:          "development",
		RedisHost:       "localhost",
		RedisPort:       "6379",
		RateLimitLimit:  "10",
		RateLimitWindow: "30s",
	}

	// Create rate limiter with custom config
	app := fiber.New()
	app.Use(NewRateLimiter(cfg))

	// Test route
	app.Get("/custom-rate", func(c *fiber.Ctx) error {
		return c.SendStatus(fiber.StatusOK)
	})

	// Test requests within limit
	for i := 0; i < 10; i++ {
		req := httptest.NewRequest("GET", "/custom-rate", nil)
		resp := httptest.NewRecorder()
		app.Handler()(resp, req)
		assert.Equal(t, fiber.StatusOK, resp.Code)
	}

	// 11th request should be rate limited
	req11 := httptest.NewRequest("GET", "/custom-rate", nil)
	resp11 := httptest.NewRecorder()
	app.Handler()(resp11, req11)
	assert.Equal(t, fiber.StatusTooManyRequests, resp11.Code)
}

func TestRateLimiter_IPBasedLimiting(t *testing.T) {
	// Test rate limiting per IP address
	cfg := &config.Config{
		AppEnv:       "development",
		RedisHost:    "localhost",
		RedisPort:    "6379",
		RateLimitLimit: "3",
		RateLimitWindow: "1m",
	}

	app := fiber.New()
	app.Use(NewRateLimiter(cfg))

	// Test route
	app.Get("/ip-limited", func(c *fiber.Ctx) error {
		return c.SendStatus(fiber.StatusOK)
	})

	// Test requests from same IP
	for i := 0; i < 3; i++ {
		req := httptest.NewRequest("GET", "/ip-limited", nil)
		req.RemoteAddr = "192.168.1.1:12345"
		resp := httptest.NewRecorder()
		app.Handler()(resp, req)
		assert.Equal(t, fiber.StatusOK, resp.Code)
	}

	// 4th request from same IP should be rate limited
	req4 := httptest.NewRequest("GET", "/ip-limited", nil)
	req4.RemoteAddr = "192.168.1.1:12345"
	resp4 := httptest.NewRecorder()
	app.Handler()(resp4, req4)
	assert.Equal(t, fiber.StatusTooManyRequests, resp4.Code)

	// Test different IP (should succeed)
	req5 := httptest.NewRequest("GET", "/ip-limited", nil)
	req5.RemoteAddr = "192.168.1.2:12345"
	resp5 := httptest.NewRecorder()
	app.Handler()(resp5, req5)
	assert.Equal(t, fiber.StatusOK, resp5.Code)
}

func TestRateLimiter_ConcurrentRequests(t *testing.T) {
	// Test rate limiting with concurrent requests
	cfg := &config.Config{
		AppEnv:       "development",
		RedisHost:    "localhost",
		RedisPort:    "6379",
		RateLimitLimit: "2",
		RateLimitWindow: "1s",
	}

	app := fiber.New()
	app.Use(NewRateLimiter(cfg))

	// Test route
	app.Get("/concurrent", func(c *fiber.Ctx) error {
		time.Sleep(100 * time.Millisecond) // Simulate some processing time
		return c.SendStatus(fiber.StatusOK)
	})

	// Send 3 concurrent requests
	done := make(chan bool, 3)
	responses := make([]int, 3)

	for i := 0; i < 3; i++ {
		go func(index int) {
			req := httptest.NewRequest("GET", "/concurrent", nil)
			resp := httptest.NewRecorder()
			app.Handler()(resp, req)
			responses[index] = resp.Code
			done <- true
		}(i)
	}

	// Wait for all requests to complete
	for i := 0; i < 3; i++ {
		<-done
	}

	// First 2 requests should succeed, 3rd should be rate limited
	assert.Equal(t, fiber.StatusOK, responses[0])
	assert.Equal(t, fiber.StatusOK, responses[1])
	assert.Equal(t, fiber.StatusTooManyRequests, responses[2])
}

func TestRateLimiter_ErrorHandling(t *testing.T) {
	// Test rate limiter with Redis connection issues
	cfg := &config.Config{
		AppEnv:       "development",
		RedisHost:    "invalid-host", // Invalid Redis host
		RedisPort:    "6379",
		RateLimitLimit: "10",
		RateLimitWindow: "1m",
	}

	// Create app with rate limiter
	app := fiber.New()
	app.Use(NewRateLimiter(cfg))

	// Test route
	app.Get("/error-test", func(c *fiber.Ctx) error {
		return c.SendStatus(fiber.StatusOK)
	})

	// Request should still work even with Redis issues (fallback behavior)
	req := httptest.NewRequest("GET", "/error-test", nil)
	resp := httptest.NewRecorder()
	app.Handler()(resp, req)
	assert.Equal(t, fiber.StatusOK, resp.Code)
}

func TestRateLimiter_ProductionVsDevelopment(t *testing.T) {
	// Test different rate limits for production vs development
	tests := []struct {
		env           string
		expectedLimit int
	}{
		{"production", 120},
		{"development", 2400},
		{"staging", 2400}, // Default to development limits
		{"test", 2400},
	}

	for _, tt := range tests {
		t.Run(tt.env, func(t *testing.T) {
			cfg := &config.Config{
				AppEnv:          tt.env,
				RedisHost:       "localhost",
				RedisPort:       "6379",
				RateLimitLimit:  "",
				RateLimitWindow: "",
			}

			app := fiber.New()
			app.Use(NewRateLimiter(cfg))

			// Test route
			app.Get("/env-test", func(c *fiber.Ctx) error {
				return c.SendStatus(fiber.StatusOK)
			})

			// Test first request
			req := httptest.NewRequest("GET", "/env-test", nil)
			resp := httptest.NewRecorder()
			app.Handler()(resp, req)
			assert.Equal(t, fiber.StatusOK, resp.Code)
		})
	}
}

func TestRateLimiter_RateLimitReset(t *testing.T) {
	// Test rate limit reset functionality
	cfg := &config.Config{
		AppEnv:          "development",
		RedisHost:       "localhost",
		RedisPort:       "6379",
		RateLimitLimit:  "2",
		RateLimitWindow: "2s", // Short window for testing
	}

	app := fiber.New()
	app.Use(NewRateLimiter(cfg))

	// Test route
	app.Get("/reset-test", func(c *fiber.Ctx) error {
		return c.SendStatus(fiber.StatusOK)
	})

	// First 2 requests should succeed
	for i := 0; i < 2; i++ {
		req := httptest.NewRequest("GET", "/reset-test", nil)
		resp := httptest.NewRecorder()
		app.Handler()(resp, req)
		assert.Equal(t, fiber.StatusOK, resp.Code)
	}

	// 3rd request should be rate limited
	req3 := httptest.NewRequest("GET", "/reset-test", nil)
	resp3 := httptest.NewRecorder()
	app.Handler()(resp3, req3)
	assert.Equal(t, fiber.StatusTooManyRequests, resp3.Code)

	// Wait for window to reset
	time.Sleep(3 * time.Second)

	// Request should succeed after reset
	req4 := httptest.NewRequest("GET", "/reset-test", nil)
	resp4 := httptest.NewRecorder()
	app.Handler()(resp4, req4)
	assert.Equal(t, fiber.StatusOK, resp4.Code)
}