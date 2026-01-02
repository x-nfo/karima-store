package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
)

func TestCORSConfiguration(t *testing.T) {
	// Test CORS middleware with different configurations
	tests := []struct {
		name           string
		allowedOrigins string
		method         string
		expectedStatus int
	}{
		{
			name:           "Allowed origin - GET request",
			allowedOrigins: "https://example.com",
			method:        "GET",
			expectedStatus: fiber.StatusOK,
		},
		{
			name:           "Allowed origin - POST request",
			allowedOrigins: "https://example.com",
			method:        "POST",
			expectedStatus: fiber.StatusOK,
		},
		{
			name:           "Allowed origin - OPTIONS request (preflight)",
			allowedOrigins: "https://example.com",
			method:        "OPTIONS",
			expectedStatus: fiber.StatusNoContent,
		},
		{
			name:           "Disallowed origin",
			allowedOrigins: "https://example.com",
			method:        "GET",
			expectedStatus: fiber.StatusOK, // Should still work, just no CORS headers
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create app with CORS middleware
			app := fiber.New()
			app.Use(CORS(tt.allowedOrigins))

			// Test route
			app.Get("/test", func(c *fiber.Ctx) error {
				return c.SendStatus(fiber.StatusOK)
			})

			// Create request
			req := httptest.NewRequest(tt.method, "/test", nil)
			if tt.method == "OPTIONS" {
				req.Header.Set("Origin", "https://example.com")
				req.Header.Set("Access-Control-Request-Method", "GET")
				req.Header.Set("Access-Control-Request-Headers", "Content-Type")
			} else if tt.method != "GET" {
				req.Header.Set("Origin", "https://example.com")
			}

			resp := httptest.NewRecorder()
			app.Handler()(resp, req)

			// Check response
			assert.Equal(t, tt.expectedStatus, resp.Code)

			// Check CORS headers for non-OPTIONS requests
			if tt.method != "OPTIONS" && tt.allowedOrigins != "" {
				assert.Equal(t, tt.allowedOrigins, resp.Header().Get("Access-Control-Allow-Origin"))
				assert.Equal(t, "GET,POST,PUT,DELETE,OPTIONS", resp.Header().Get("Access-Control-Allow-Methods"))
				assert.Equal(t, "Origin,Content-Type,Accept,Authorization,X-Requested-With", resp.Header().Get("Access-Control-Allow-Headers"))
				assert.Equal(t, "true", resp.Header().Get("Access-Control-Allow-Credentials"))
				assert.Equal(t, "86400", resp.Header().Get("Access-Control-Max-Age"))
			}
		})
	}
}

func TestCORSSecurity(t *testing.T) {
	// Test CORS security aspects
	app := fiber.New()
	app.Use(CORS("https://example.com"))

	// Test route that requires authentication
	app.Get("/secure", func(c *fiber.Ctx) error {
		return c.SendStatus(fiber.StatusOK)
	})

	// Test case 1: Cross-origin request with credentials
	req1 := httptest.NewRequest("GET", "/secure", nil)
	req1.Header.Set("Origin", "https://malicious.com")
	req1.Header.Set("Cookie", "session=token")
	resp1 := httptest.NewRecorder()
	app.Handler()(resp1, req1)

	// Should not expose sensitive headers to malicious origin
	assert.Equal(t, fiber.StatusOK, resp1.Code)
	assert.Equal(t, "", resp1.Header().Get("Access-Control-Allow-Origin")) // Should not be set for disallowed origin

	// Test case 2: Same-origin request
	req2 := httptest.NewRequest("GET", "/secure", nil)
	req2.Header.Set("Origin", "https://example.com")
	resp2 := httptest.NewRecorder()
	app.Handler()(resp2, req2)

	// Should have CORS headers for allowed origin
	assert.Equal(t, fiber.StatusOK, resp2.Code)
	assert.Equal(t, "https://example.com", resp2.Header().Get("Access-Control-Allow-Origin"))

	// Test case 3: Preflight request with authentication headers
	req3 := httptest.NewRequest("OPTIONS", "/secure", nil)
	req3.Header.Set("Origin", "https://example.com")
	req3.Header.Set("Access-Control-Request-Method", "GET")
	req3.Header.Set("Access-Control-Request-Headers", "Authorization")
	resp3 := httptest.NewRecorder()
	app.Handler()(resp3, req3)

	// Should handle preflight correctly
	assert.Equal(t, fiber.StatusNoContent, resp3.Code)
	assert.Equal(t, "https://example.com", resp3.Header().Get("Access-Control-Allow-Origin"))
}

func TestCORSWithAuthentication(t *testing.T) {
	app := fiber.New()
	app.Use(CORS("https://example.com"))

	// Test route with authentication
	app.Get("/auth", func(c *fiber.Ctx) error {
		return c.SendStatus(fiber.StatusOK)
	})

	// Test case 1: Request with authentication token
	req1 := httptest.NewRequest("GET", "/auth", nil)
	req1.Header.Set("Origin", "https://example.com")
	req1.Header.Set("Authorization", "Bearer token")
	resp1 := httptest.NewRecorder()
	app.Handler()(resp1, req1)

	// Should handle authentication headers correctly
	assert.Equal(t, fiber.StatusOK, resp1.Code)
	assert.Equal(t, "https://example.com", resp1.Header().Get("Access-Control-Allow-Origin"))

	// Test case 2: Cross-origin request with authentication token
	req2 := httptest.NewRequest("GET", "/auth", nil)
	req2.Header.Set("Origin", "https://malicious.com")
	req2.Header.Set("Authorization", "Bearer token")
	resp2 := httptest.NewRecorder()
	app.Handler()(resp2, req2)

	// Should not expose sensitive information
	assert.Equal(t, fiber.StatusOK, resp2.Code)
	assert.Equal(t, "", resp2.Header().Get("Access-Control-Allow-Origin"))
}

func TestCORSHeaderConfiguration(t *testing.T) {
	// Test different CORS header configurations
	tests := []struct {
		name           string
		allowedOrigins string
		expectedOrigin string
	}{
		{
			name:           "Single origin",
			allowedOrigins: "https://example.com",
			expectedOrigin: "https://example.com",
		},
		{
			name:           "Multiple origins (comma-separated)",
			allowedOrigins: "https://example.com,https://test.com",
			expectedOrigin: "https://example.com,https://test.com",
		},
		{
			name:           "Wildcard origin",
			allowedOrigins: "*",
			expectedOrigin: "*",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			app := fiber.New()
			app.Use(CORS(tt.allowedOrigins))

			// Test route
			app.Get("/headers", func(c *fiber.Ctx) error {
				return c.SendStatus(fiber.StatusOK)
			})

			// Create request
			req := httptest.NewRequest("GET", "/headers", nil)
			req.Header.Set("Origin", "https://example.com")
			resp := httptest.NewRecorder()
			app.Handler()(resp, req)

			// Check CORS headers
			assert.Equal(t, fiber.StatusOK, resp.Code)
			assert.Equal(t, tt.expectedOrigin, resp.Header().Get("Access-Control-Allow-Origin"))
			assert.Equal(t, "GET,POST,PUT,DELETE,OPTIONS", resp.Header().Get("Access-Control-Allow-Methods"))
			assert.Equal(t, "Origin,Content-Type,Accept,Authorization,X-Requested-With", resp.Header().Get("Access-Control-Allow-Headers"))
			assert.Equal(t, "true", resp.Header().Get("Access-Control-Allow-Credentials"))
			assert.Equal(t, "86400", resp.Header().Get("Access-Control-Max-Age"))
		})
	}
}

func TestCORS_MisconfigurationAttacks(t *testing.T) {
	app := fiber.New()
	app.Use(CORS("https://example.com"))

	// Test route
	app.Get("/vulnerable", func(c *fiber.Ctx) error {
		return c.SendStatus(fiber.StatusOK)
	})

	// Test case 1: Attempt to exploit CORS misconfiguration
	req1 := httptest.NewRequest("GET", "/vulnerable", nil)
	req1.Header.Set("Origin", "null") // Null origin attack
	resp1 := httptest.NewRecorder()
	app.Handler()(resp1, req1)

	// Should not be vulnerable to null origin attack
	assert.Equal(t, fiber.StatusOK, resp1.Code)
	assert.Equal(t, "", resp1.Header().Get("Access-Control-Allow-Origin"))

	// Test case 2: File protocol attack
	req2 := httptest.NewRequest("GET", "/vulnerable", nil)
	req2.Header.Set("Origin", "file://")
	resp2 := httptest.NewRecorder()
	app.Handler()(resp2, req2)

	// Should block file protocol
	assert.Equal(t, fiber.StatusOK, resp2.Code)
	assert.Equal(t, "", resp2.Header().Get("Access-Control-Allow-Origin"))

	// Test case 3: Data URL attack
	req3 := httptest.NewRequest("GET", "/vulnerable", nil)
	req3.Header.Set("Origin", "data://text/plain,hello")
	resp3 := httptest.NewRecorder()
	app.Handler()(resp3, req3)

	// Should block data URL
	assert.Equal(t, fiber.StatusOK, resp3.Code)
	assert.Equal(t, "", resp3.Header().Get("Access-Control-Allow-Origin"))
}

func TestCORS_PreflightRequest(t *testing.T) {
	app := fiber.New()
	app.Use(CORS("https://example.com"))

	// Test route
	app.Get("/preflight", func(c *fiber.Ctx) error {
		return c.SendStatus(fiber.StatusOK)
	})

	// Test preflight request
	req := httptest.NewRequest("OPTIONS", "/preflight", nil)
	req.Header.Set("Origin", "https://example.com")
	req.Header.Set("Access-Control-Request-Method", "GET")
	req.Header.Set("Access-Control-Request-Headers", "Content-Type,Authorization")
	resp := httptest.NewRecorder()
	app.Handler()(resp, req)

	// Check preflight response
	assert.Equal(t, fiber.StatusNoContent, resp.Code)
	assert.Equal(t, "https://example.com", resp.Header().Get("Access-Control-Allow-Origin"))
	assert.Equal(t, "GET,POST,PUT,DELETE,OPTIONS", resp.Header().Get("Access-Control-Allow-Methods"))
	assert.Equal(t, "Origin,Content-Type,Accept,Authorization,X-Requested-With", resp.Header().Get("Access-Control-Allow-Headers"))
	assert.Equal(t, "true", resp.Header().Get("Access-Control-Allow-Credentials"))
	assert.Equal(t, "86400", resp.Header().Get("Access-Control-Max-Age"))
}