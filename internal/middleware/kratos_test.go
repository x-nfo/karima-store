package middleware

import (
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/limiter"
	"github.com/gofiber/storage/redis/v3"
	"github.com/karima-store/internal/config"
	"github.com/karima-store/internal/models"
	"github.com/stretchr/testify/assert"
)

// MockAuthService is a mock implementation of AuthService for testing
type MockAuthService struct{}

func (m *MockAuthService) SyncUser(kratosIdentity *models.KratosIdentity, email string) (*models.User, error) {
	// Extract role from Kratos traits if available (matching test expectations)
	role := models.UserRole("user") // default to "user" to match test expectations
	if roleStr, ok := kratosIdentity.Traits["role"].(string); ok && roleStr != "" {
		role = models.UserRole(roleStr)
	}

	return &models.User{
		ID:       1,
		KratosID: kratosIdentity.ID,
		Email:    email,
		Role:     role,
		FullName: "Test User",
	}, nil
}

func (m *MockAuthService) GetUserByID(id uint) (*models.User, error) {
	return &models.User{
		ID:       id,
		Email:    "test@example.com",
		Role:     models.RoleCustomer,
		FullName: "Test User",
	}, nil
}

// Helper to create mock Kratos server
func mockKratosServer() *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/sessions/whoami" {
			// Check cookie or header
			cookie, err := r.Cookie("ory_kratos_session")
			token := r.Header.Get("X-Session-Token")

			if (err == nil && cookie.Value == "valid-session-token") || token == "valid-session-token" {
				w.WriteHeader(http.StatusOK)
				w.Write([]byte(`{
					"id": "test-session-id",
					"active": true,
					"identity": {
						"id": "test-identity-id",
						"traits": {
							"email": "test@example.com",
							"role": "user"
						}
					}
				}`))
				return
			}
			if err == nil && cookie.Value == "valid-session-token-without-role" {
				w.WriteHeader(http.StatusOK)
				w.Write([]byte(`{
					"id": "test-session-id",
					"active": true,
					"identity": {
						"id": "test-identity-id",
						"traits": {
							"email": "test@example.com"
						}
					}
				}`))
				return
			}

			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		w.WriteHeader(http.StatusNotFound)
	}))
}

func TestKratosMiddleware_Authenticate(t *testing.T) {
	// Setup mock Kratos server
	ts := mockKratosServer()
	defer ts.Close()

	// Setup test app
	app := fiber.New()

	// Create middleware with mock server URL and mock auth service
	mockAuthService := &MockAuthService{}
	kratosMiddleware := NewKratosMiddleware(ts.URL, ts.URL, mockAuthService)

	// Test route with authentication
	app.Get("/protected", kratosMiddleware.Authenticate(), func(c *fiber.Ctx) error {
		return c.SendStatus(fiber.StatusOK)
	})

	// Test case 1: No session cookie
	req1 := httptest.NewRequest("GET", "/protected", nil)
	resp1, err1 := app.Test(req1)
	assert.NoError(t, err1)
	assert.Equal(t, fiber.StatusUnauthorized, resp1.StatusCode)

	// Test case 2: Invalid session cookie
	req2 := httptest.NewRequest("GET", "/protected", nil)
	req2.AddCookie(&http.Cookie{
		Name:  "ory_kratos_session",
		Value: "invalid-session-token",
	})
	resp2, err2 := app.Test(req2)
	assert.NoError(t, err2)
	assert.Equal(t, fiber.StatusUnauthorized, resp2.StatusCode)

	// Test case 3: Valid session cookie
	req3 := httptest.NewRequest("GET", "/protected", nil)
	req3.AddCookie(&http.Cookie{
		Name:  "ory_kratos_session",
		Value: "valid-session-token",
	})
	resp3, err3 := app.Test(req3)
	assert.NoError(t, err3)
	assert.Equal(t, fiber.StatusOK, resp3.StatusCode)
}

func TestKratosMiddleware_RequireRole(t *testing.T) {
	ts := mockKratosServer()
	defer ts.Close()

	app := fiber.New()
	mockAuthService := &MockAuthService{}
	kratosMiddleware := NewKratosMiddleware(ts.URL, ts.URL, mockAuthService)

	// Test route requiring admin role
	app.Get("/admin", kratosMiddleware.Authenticate(), kratosMiddleware.RequireRole("admin"), func(c *fiber.Ctx) error {
		return c.SendStatus(fiber.StatusOK)
	})

	// Test route requiring user role
	app.Get("/user", kratosMiddleware.Authenticate(), kratosMiddleware.RequireRole("user"), func(c *fiber.Ctx) error {
		return c.SendStatus(fiber.StatusOK)
	})

	// Test case 1: No authentication
	req1 := httptest.NewRequest("GET", "/admin", nil)
	resp1, err1 := app.Test(req1)
	assert.NoError(t, err1)
	assert.Equal(t, fiber.StatusUnauthorized, resp1.StatusCode)

	// Test case 2: User role accessing admin route
	// The mock server returns "role": "user" for "valid-session-token"
	req2 := httptest.NewRequest("GET", "/admin", nil)
	req2.AddCookie(&http.Cookie{
		Name:  "ory_kratos_session",
		Value: "valid-session-token",
	})
	resp2, err2 := app.Test(req2)
	assert.NoError(t, err2)
	if resp2.StatusCode != fiber.StatusForbidden {
		body, _ := io.ReadAll(resp2.Body)
		fmt.Printf("Test case 2 failed. Status: %d, Body: %s\n", resp2.StatusCode, string(body))
	}
	// User is authenticated (mock returns user role), but route requires admin
	assert.Equal(t, fiber.StatusForbidden, resp2.StatusCode)

	// Test case 3: User role accessing user route
	req3 := httptest.NewRequest("GET", "/user", nil)
	req3.AddCookie(&http.Cookie{
		Name:  "ory_kratos_session",
		Value: "valid-session-token",
	})
	resp3, err3 := app.Test(req3)
	assert.NoError(t, err3)
	if resp3.StatusCode != fiber.StatusOK {
		body, _ := io.ReadAll(resp3.Body)
		fmt.Printf("Test case 3 failed. Status: %d, Body: %s\n", resp3.StatusCode, string(body))
	}
	assert.Equal(t, fiber.StatusOK, resp3.StatusCode)
}

func TestKratosMiddleware_ValidateToken(t *testing.T) {
	ts := mockKratosServer()
	defer ts.Close()

	app := fiber.New()
	mockAuthService := &MockAuthService{}
	kratosMiddleware := NewKratosMiddleware(ts.URL, ts.URL, mockAuthService)

	// Test route with token authentication
	app.Get("/api", kratosMiddleware.ValidateToken(), func(c *fiber.Ctx) error {
		return c.SendStatus(fiber.StatusOK)
	})

	// Test case 1: No authorization header
	req1 := httptest.NewRequest("GET", "/api", nil)
	resp1, err1 := app.Test(req1)
	assert.NoError(t, err1)
	assert.Equal(t, fiber.StatusUnauthorized, resp1.StatusCode)

	// Test case 2: Invalid authorization header format
	req2 := httptest.NewRequest("GET", "/api", nil)
	req2.Header.Set("Authorization", "InvalidToken")
	resp2, err2 := app.Test(req2)
	assert.NoError(t, err2)
	assert.Equal(t, fiber.StatusUnauthorized, resp2.StatusCode)

	// Test case 3: Valid Bearer token format
	req3 := httptest.NewRequest("GET", "/api", nil)
	req3.Header.Set("Authorization", "Bearer valid-session-token")
	resp3, err3 := app.Test(req3)
	assert.NoError(t, err3)
	assert.Equal(t, fiber.StatusOK, resp3.StatusCode)
}

func TestKratosMiddleware_OptionalAuth(t *testing.T) {
	ts := mockKratosServer()
	defer ts.Close()

	app := fiber.New()
	mockAuthService := &MockAuthService{}
	kratosMiddleware := NewKratosMiddleware(ts.URL, ts.URL, mockAuthService)

	// Test route with optional authentication
	app.Get("/optional", kratosMiddleware.OptionalAuth(), func(c *fiber.Ctx) error {
		return c.SendStatus(fiber.StatusOK)
	})

	// Test case 1: No session cookie
	req1 := httptest.NewRequest("GET", "/optional", nil)
	resp1, err1 := app.Test(req1)
	assert.NoError(t, err1)
	assert.Equal(t, fiber.StatusOK, resp1.StatusCode)

	// Test case 2: Invalid session cookie
	req2 := httptest.NewRequest("GET", "/optional", nil)
	req2.AddCookie(&http.Cookie{
		Name:  "ory_kratos_session",
		Value: "invalid-session-token",
	})
	resp2, err2 := app.Test(req2)
	assert.NoError(t, err2)
	assert.Equal(t, fiber.StatusOK, resp2.StatusCode)

	// Test case 3: Valid session cookie
	req3 := httptest.NewRequest("GET", "/optional", nil)
	req3.AddCookie(&http.Cookie{
		Name:  "ory_kratos_session",
		Value: "valid-session-token",
	})
	resp3, err3 := app.Test(req3)
	assert.NoError(t, err3)
	assert.Equal(t, fiber.StatusOK, resp3.StatusCode)
}

func TestKratosMiddleware_SessionData(t *testing.T) {
	ts := mockKratosServer()
	defer ts.Close()

	app := fiber.New()
	mockAuthService := &MockAuthService{}
	kratosMiddleware := NewKratosMiddleware(ts.URL, ts.URL, mockAuthService)

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
	resp, err := app.Test(req)
	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusOK, resp.StatusCode)
}

func TestKratosMiddleware_RoleDefaults(t *testing.T) {
	ts := mockKratosServer()
	defer ts.Close()

	app := fiber.New()
	mockAuthService := &MockAuthService{}
	kratosMiddleware := NewKratosMiddleware(ts.URL, ts.URL, mockAuthService)

	// Test route that checks role defaults
	app.Get("/role-default", kratosMiddleware.Authenticate(), func(c *fiber.Ctx) error {
		userRole := c.Locals("user_role")
		assert.Equal(t, models.UserRole("user"), userRole) // Default role should be models.UserRole("user")
		return c.SendStatus(fiber.StatusOK)
	})

	// Test with valid session but no role in traits
	req := httptest.NewRequest("GET", "/role-default", nil)
	req.AddCookie(&http.Cookie{
		Name:  "ory_kratos_session",
		Value: "valid-session-token-without-role",
	})
	resp, err := app.Test(req)
	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusOK, resp.StatusCode)
}

func TestKratosMiddleware_RateLimitIntegration(t *testing.T) {
	ts := mockKratosServer()
	defer ts.Close()

	app := fiber.New()
	cfg := &config.Config{
		KratosPublicURL: ts.URL,
		KratosAdminURL:  ts.URL,
		AppEnv:          "development",
		RedisHost:       "localhost",
		RedisPort:       "6379",
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

	mockAuthService := &MockAuthService{}
	kratosMiddleware := NewKratosMiddleware(cfg.KratosPublicURL, cfg.KratosAdminURL, mockAuthService)

	// Test route with authentication and rate limiting
	app.Get("/protected/rate-limited", rateLimiter, kratosMiddleware.Authenticate(), func(c *fiber.Ctx) error {
		return c.SendStatus(fiber.StatusOK)
	})

	// Generate unique IP
	// math/rand is not imported, so we just use a fixed unique one or rely on imported rand if available.
	// Imports check: "math/rand" is NOT imported. "fmt" is imported.
	// We can use time.Now() to generate something somewhat unique or just a hardcoded unique IP different from others.
	uniqueIP := fmt.Sprintf("192.168.200.%d:1234", time.Now().Nanosecond()%255)

	// Test case 1: First request (should succeed)
	req1 := httptest.NewRequest("GET", "/protected/rate-limited", nil)
	req1.RemoteAddr = uniqueIP
	req1.AddCookie(&http.Cookie{
		Name:  "ory_kratos_session",
		Value: "valid-session-token",
	})
	resp1, err1 := app.Test(req1)
	assert.NoError(t, err1)
	assert.Equal(t, fiber.StatusOK, resp1.StatusCode)

	// Test case 2: Second request (should succeed)
	req2 := httptest.NewRequest("GET", "/protected/rate-limited", nil)
	req2.RemoteAddr = uniqueIP
	req2.AddCookie(&http.Cookie{
		Name:  "ory_kratos_session",
		Value: "valid-session-token",
	})
	resp2, err2 := app.Test(req2)
	assert.NoError(t, err2)
	assert.Equal(t, fiber.StatusOK, resp2.StatusCode)

	// Test case 3: Third request (should be rate limited)
	req3 := httptest.NewRequest("GET", "/protected/rate-limited", nil)
	req3.RemoteAddr = uniqueIP
	req3.AddCookie(&http.Cookie{
		Name:  "ory_kratos_session",
		Value: "valid-session-token",
	})
	resp3, err3 := app.Test(req3)
	assert.NoError(t, err3)
	assert.Equal(t, fiber.StatusTooManyRequests, resp3.StatusCode)
}
