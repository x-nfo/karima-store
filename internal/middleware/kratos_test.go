package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/limiter"
	"github.com/gofiber/storage/redis/v3"
	"github.com/karima-store/internal/config"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockHTTPClient for testing
type MockHTTPClient struct {
	mock.Mock
}

func (m *MockHTTPClient) Do(req *http.Request) (*http.Response, error) {
	args := m.Called(req)
	return args.Get(0).(*http.Response), args.Error(1)
}

func TestKratosMiddleware_Authenticate(t *testing.T) {
	// Setup test app
	app := fiber.New()
	cfg := &config.Config{
		KratosPublicURL: "http://kratos-public",
		KratosAdminURL:  "http://kratos-admin",
	}

	// Create middleware
	kratosMiddleware := NewKratosMiddleware(cfg.KratosPublicURL, cfg.KratosAdminURL)

	// Test route with authentication
	app.Get("/protected", kratosMiddleware.Authenticate(), func(c *fiber.Ctx) error {
		return c.SendStatus(fiber.StatusOK)
	})

	// Test case 1: No session cookie
	req1 := httptest.NewRequest("GET", "/protected", nil)
	resp1 := httptest.NewRecorder()
	app.Handler()(resp1, req1)
	assert.Equal(t, fiber.StatusUnauthorized, resp1.Code)
	assert.Contains(t, resp1.Body.String(), "No session cookie found")

	// Test case 2: Invalid session cookie
	req2 := httptest.NewRequest("GET", "/protected", nil)
	req2.AddCookie(&http.Cookie{
		Name:  "ory_kratos_session",
		Value: "invalid-session-token",
	})
	resp2 := httptest.NewRecorder()
	app.Handler()(resp2, req2)
	assert.Equal(t, fiber.StatusUnauthorized, resp2.Code)
	assert.Contains(t, resp2.Body.String(), "Invalid or expired session")

	// Test case 3: Expired session
	// This would require mocking the HTTP client to return 401
	// For simplicity, we'll test the basic flow
}

func TestKratosMiddleware_RequireRole(t *testing.T) {
	app := fiber.New()
	cfg := &config.Config{
		KratosPublicURL: "http://kratos-public",
		KratosAdminURL:  "http://kratos-admin",
	}

	kratosMiddleware := NewKratosMiddleware(cfg.KratosPublicURL, cfg.KratosAdminURL)

	// Test route requiring admin role
	app.Get("/admin", kratosMiddleware.RequireRole("admin"), func(c *fiber.Ctx) error {
		return c.SendStatus(fiber.StatusOK)
	})

	// Test route requiring user role
	app.Get("/user", kratosMiddleware.RequireRole("user"), func(c *fiber.Ctx) error {
		return c.SendStatus(fiber.StatusOK)
	})

	// Test case 1: No authentication
	req1 := httptest.NewRequest("GET", "/admin", nil)
	resp1 := httptest.NewRecorder()
	app.Handler()(resp1, req1)
	assert.Equal(t, fiber.StatusUnauthorized, resp1.Code)

	// Test case 2: User role accessing admin route
	req2 := httptest.NewRequest("GET", "/admin", nil)
	req2.AddCookie(&http.Cookie{
		Name:  "ory_kratos_session",
		Value: "valid-session-token",
	})
	resp2 := httptest.NewRecorder()
	app.Handler()(resp2, req2)
	assert.Equal(t, fiber.StatusForbidden, resp2.Code)
	assert.Contains(t, resp2.Body.String(), "Insufficient permissions")

	// Test case 3: Admin role accessing admin route
	// This would require mocking the session validation
}

func TestKratosMiddleware_ValidateToken(t *testing.T) {
	app := fiber.New()
	cfg := &config.Config{
		KratosPublicURL: "http://kratos-public",
		KratosAdminURL:  "http://kratos-admin",
	}

	kratosMiddleware := NewKratosMiddleware(cfg.KratosPublicURL, cfg.KratosAdminURL)

	// Test route with token authentication
	app.Get("/api", kratosMiddleware.ValidateToken(), func(c *fiber.Ctx) error {
		return c.SendStatus(fiber.StatusOK)
	})

	// Test case 1: No authorization header
	req1 := httptest.NewRequest("GET", "/api", nil)
	resp1 := httptest.NewRecorder()
	app.Handler()(resp1, req1)
	assert.Equal(t, fiber.StatusUnauthorized, resp1.Code)

	// Test case 2: Invalid authorization header format
	req2 := httptest.NewRequest("GET", "/api", nil)
	req2.Header.Set("Authorization", "InvalidToken")
	resp2 := httptest.NewRecorder()
	app.Handler()(resp2, req2)
	assert.Equal(t, fiber.StatusUnauthorized, resp2.Code)

	// Test case 3: Valid Bearer token format
	req3 := httptest.NewRequest("GET", "/api", nil)
	req3.Header.Set("Authorization", "Bearer valid-session-token")
	resp3 := httptest.NewRecorder()
	app.Handler()(resp3, req3)
	// This would require mocking the session validation
	assert.Equal(t, fiber.StatusOK, resp3.Code)
}

func TestKratosMiddleware_OptionalAuth(t *testing.T) {
	app := fiber.New()
	cfg := &config.Config{
		KratosPublicURL: "http://kratos-public",
		KratosAdminURL:  "http://kratos-admin",
	}

	kratosMiddleware := NewKratosMiddleware(cfg.KratosPublicURL, cfg.KratosAdminURL)

	// Test route with optional authentication
	app.Get("/optional", kratosMiddleware.OptionalAuth(), func(c *fiber.Ctx) error {
		return c.SendStatus(fiber.StatusOK)
	})

	// Test case 1: No session cookie
	req1 := httptest.NewRequest("GET", "/optional", nil)
	resp1 := httptest.NewRecorder()
	app.Handler()(resp1, req1)
	assert.Equal(t, fiber.StatusOK, resp1.Code)

	// Test case 2: Invalid session cookie
	req2 := httptest.NewRequest("GET", "/optional", nil)
	req2.AddCookie(&http.Cookie{
		Name:  "ory_kratos_session",
		Value: "invalid-session-token",
	})
	resp2 := httptest.NewRecorder()
	app.Handler()(resp2, req2)
	assert.Equal(t, fiber.StatusOK, resp2.Code)

	// Test case 3: Valid session cookie
	req3 := httptest.NewRequest("GET", "/optional", nil)
	req3.AddCookie(&http.Cookie{
		Name:  "ory_kratos_session",
		Value: "valid-session-token",
	})
	resp3 := httptest.NewRecorder()
	app.Handler()(resp3, req3)
	assert.Equal(t, fiber.StatusOK, resp3.Code)
}

func TestKratosMiddleware_SessionData(t *testing.T) {
	app := fiber.New()
	cfg := &config.Config{
		KratosPublicURL: "http://kratos-public",
		KratosAdminURL:  "http://kratos-admin",
	}

	kratosMiddleware := NewKratosMiddleware(cfg.KratosPublicURL, cfg.KratosAdminURL)

	// Test route that checks session data
	app.Get("/session-data", kratosMiddleware.Authenticate(), func(c *fiber.Ctx) error {
		identityID := c.Locals("identity_id")
		userEmail := c.Locals("user_email")
		userRole := c.Locals("user_role")

		assert.NotEmpty(t, identityID)
		assert.NotEmpty(t, userEmail)
		assert.NotEmpty(t, userRole)

		return c.JSON(fiber.Map{
			"identity_id": identityID,
			"user_email":  userEmail,
			"user_role":   userRole,
		})
	})

	// Test with valid session
	req := httptest.NewRequest("GET", "/session-data", nil)
	req.AddCookie(&http.Cookie{
		Name:  "ory_kratos_session",
		Value: "valid-session-token",
	})
	resp := httptest.NewRecorder()
	app.Handler()(resp, req)
	assert.Equal(t, fiber.StatusOK, resp.Code)
}

func TestKratosMiddleware_RoleDefaults(t *testing.T) {
	app := fiber.New()
	cfg := &config.Config{
		KratosPublicURL: "http://kratos-public",
		KratosAdminURL:  "http://kratos-admin",
	}

	kratosMiddleware := NewKratosMiddleware(cfg.KratosPublicURL, cfg.KratosAdminURL)

	// Test route that checks role defaults
	app.Get("/role-default", kratosMiddleware.Authenticate(), func(c *fiber.Ctx) error {
		userRole := c.Locals("user_role")
		assert.Equal(t, "user", userRole) // Default role should be "user"
		return c.SendStatus(fiber.StatusOK)
	})

	// Test with valid session but no role in traits
	req := httptest.NewRequest("GET", "/role-default", nil)
	req.AddCookie(&http.Cookie{
		Name:  "ory_kratos_session",
		Value: "valid-session-token-without-role",
	})
	resp := httptest.NewRecorder()
	app.Handler()(resp, req)
	assert.Equal(t, fiber.StatusOK, resp.Code)
}

func TestKratosMiddleware_RateLimitIntegration(t *testing.T) {
	app := fiber.New()
	cfg := &config.Config{
		KratosPublicURL: "http://kratos-public",
		KratosAdminURL:  "http://kratos-admin",
		AppEnv:         "development",
		RedisHost:      "localhost",
		RedisPort:      "6379",
	}

	// Create rate limiter
	store := redis.New(redis.Config{
		Host: "localhost",
		Port: 6379,
	})
	rateLimiter := limiter.New(limiter.Config{
		Max:        2,
		Expiration: 1 * time.Minute,
		Storage:    store,
		KeyGenerator: func(c *fiber.Ctx) string {
			return c.IP()
		},
	})

	kratosMiddleware := NewKratosMiddleware(cfg.KratosPublicURL, cfg.KratosAdminURL)

	// Test route with authentication and rate limiting
	app.Get("/protected/rate-limited", rateLimiter, kratosMiddleware.Authenticate(), func(c *fiber.Ctx) error {
		return c.SendStatus(fiber.StatusOK)
	})

	// Test case 1: First request (should succeed)
	req1 := httptest.NewRequest("GET", "/protected/rate-limited", nil)
	req1.AddCookie(&http.Cookie{
		Name:  "ory_kratos_session",
		Value: "valid-session-token",
	})
	resp1 := httptest.NewRecorder()
	app.Handler()(resp1, req1)
	assert.Equal(t, fiber.StatusOK, resp1.Code)

	// Test case 2: Second request (should succeed)
	req2 := httptest.NewRequest("GET", "/protected/rate-limited", nil)
	req2.AddCookie(&http.Cookie{
		Name:  "ory_kratos_session",
		Value: "valid-session-token",
	})
	resp2 := httptest.NewRecorder()
	app.Handler()(resp2, req2)
	assert.Equal(t, fiber.StatusOK, resp2.Code)

	// Test case 3: Third request (should be rate limited)
	req3 := httptest.NewRequest("GET", "/protected/rate-limited", nil)
	req3.AddCookie(&http.Cookie{
		Name:  "ory_kratos_session",
		Value: "valid-session-token",
	})
	resp3 := httptest.NewRecorder()
	app.Handler()(resp3, req3)
	assert.Equal(t, fiber.StatusTooManyRequests, resp3.Code)
}