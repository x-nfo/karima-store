package middleware

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDefaultCSRFConfig(t *testing.T) {
	config := DefaultCSRFConfig()

	assert.Equal(t, 32, config.TokenLength)
	assert.Equal(t, 24*time.Hour, config.TokenExpiration)
	assert.Equal(t, "X-CSRF-Token", config.TokenHeader)
	assert.Equal(t, "csrf_token", config.TokenFormField)
	assert.Equal(t, "csrf_token", config.CookieName)
	assert.True(t, config.CookieSecure)
	assert.False(t, config.CookieHTTPOnly)
	assert.Equal(t, "Strict", config.CookieSameSite)
	assert.NotEmpty(t, config.ExcludedPaths)
}

func TestNewCSRFManager(t *testing.T) {
	config := DefaultCSRFConfig()
	manager := NewCSRFManager(config)

	require.NotNil(t, manager)
	assert.Equal(t, config, manager.config)
	assert.NotNil(t, manager.tokens)
	assert.NotNil(t, manager.done)

	// Cleanup
	manager.Stop()
}

func TestCSRFManager_GenerateToken(t *testing.T) {
	config := DefaultCSRFConfig()
	manager := NewCSRFManager(config)
	defer manager.Stop()

	sessionID := "test-session-123"
	token, err := manager.GenerateToken(sessionID)

	require.NoError(t, err)
	assert.NotEmpty(t, token)

	// Verify token is stored
	tokenInfo := manager.GetTokenInfo(sessionID)
	require.NotNil(t, tokenInfo)
	assert.Equal(t, token, tokenInfo.value)
	assert.True(t, tokenInfo.expiresAt.After(time.Now()))
}

func TestCSRFManager_ValidateToken(t *testing.T) {
	config := DefaultCSRFConfig()
	manager := NewCSRFManager(config)
	defer manager.Stop()

	sessionID := "test-session-123"
	token, err := manager.GenerateToken(sessionID)
	require.NoError(t, err)

	// Valid token
	assert.True(t, manager.ValidateToken(sessionID, token))

	// Invalid token
	assert.False(t, manager.ValidateToken(sessionID, "invalid-token"))

	// Wrong session
	assert.False(t, manager.ValidateToken("wrong-session", token))
}

func TestCSRFManager_RotateToken(t *testing.T) {
	config := DefaultCSRFConfig()
	manager := NewCSRFManager(config)
	defer manager.Stop()

	sessionID := "test-session-123"
	oldToken, err := manager.GenerateToken(sessionID)
	require.NoError(t, err)

	newToken, err := manager.RotateToken(sessionID)
	require.NoError(t, err)

	assert.NotEqual(t, oldToken, newToken)

	// Old token should no longer be valid
	assert.False(t, manager.ValidateToken(sessionID, oldToken))

	// New token should be valid
	assert.True(t, manager.ValidateToken(sessionID, newToken))
}

func TestCSRFManager_RevokeToken(t *testing.T) {
	config := DefaultCSRFConfig()
	manager := NewCSRFManager(config)
	defer manager.Stop()

	sessionID := "test-session-123"
	token, err := manager.GenerateToken(sessionID)
	require.NoError(t, err)

	// Token should be valid
	assert.True(t, manager.ValidateToken(sessionID, token))

	// Revoke token
	manager.RevokeToken(sessionID)

	// Token should no longer be valid
	assert.False(t, manager.ValidateToken(sessionID, token))

	// Token info should be removed
	assert.Nil(t, manager.GetTokenInfo(sessionID))
}

func TestCSRFManager_TokenExpiration(t *testing.T) {
	config := CSRFConfig{
		TokenLength:     32,
		TokenExpiration: 100 * time.Millisecond, // Short expiration for testing
	}
	manager := NewCSRFManager(config)
	defer manager.Stop()

	sessionID := "test-session-123"
	token, err := manager.GenerateToken(sessionID)
	require.NoError(t, err)

	// Token should be valid immediately
	assert.True(t, manager.ValidateToken(sessionID, token))

	// Wait for expiration
	time.Sleep(150 * time.Millisecond)

	// Token should be expired
	assert.False(t, manager.ValidateToken(sessionID, token))
}

func TestCSRFManager_CleanupExpiredTokens(t *testing.T) {
	config := CSRFConfig{
		TokenLength:     32,
		TokenExpiration: 100 * time.Millisecond,
	}
	manager := NewCSRFManager(config)
	defer manager.Stop()

	// Generate multiple tokens
	for i := 0; i < 5; i++ {
		sessionID := "test-session-" + string(rune(i))
		_, err := manager.GenerateToken(sessionID)
		require.NoError(t, err)
	}

	// Wait for expiration
	time.Sleep(150 * time.Millisecond)

	// Trigger cleanup
	manager.cleanupExpiredTokens()

	// All tokens should be removed
	for i := 0; i < 5; i++ {
		sessionID := "test-session-" + string(rune(i))
		assert.Nil(t, manager.GetTokenInfo(sessionID))
	}
}

func TestCSRFMiddleware_GET(t *testing.T) {
	app := fiber.New()
	config := DefaultCSRFConfig()
	app.Use(CSRF(config))

	app.Get("/test", func(c *fiber.Ctx) error {
		token := c.Locals(config.TokenContextKey)
		return c.JSON(fiber.Map{"token": token})
	})

	req := getTestRequest("GET", "/test", nil)
	resp, err := app.Test(req)
	require.NoError(t, err)
	assert.Equal(t, 200, resp.StatusCode)

	// Check for CSRF cookie
	cookies := resp.Cookies()
	var csrfCookie *fiber.Cookie
	for _, cookie := range cookies {
		if cookie.Name == config.CookieName {
			csrfCookie = cookie
			break
		}
	}
	assert.NotNil(t, csrfCookie)
	assert.NotEmpty(t, csrfCookie.Value)
}

func TestCSRFMiddleware_POST_WithoutToken(t *testing.T) {
	app := fiber.New()
	config := DefaultCSRFConfig()
	app.Use(CSRF(config))

	app.Post("/test", func(c *fiber.Ctx) error {
		return c.SendStatus(200)
	})

	req := getTestRequest("POST", "/test", nil)
	resp, err := app.Test(req)
	require.NoError(t, err)
	assert.Equal(t, 400, resp.StatusCode) // Bad Request - CSRF token required
}

func TestCSRFMiddleware_POST_WithInvalidToken(t *testing.T) {
	app := fiber.New()
	config := DefaultCSRFConfig()
	app.Use(CSRF(config))

	app.Post("/test", func(c *fiber.Ctx) error {
		return c.SendStatus(200)
	})

	req := getTestRequest("POST", "/test", nil)
	req.Header.Set(config.TokenHeader, "invalid-token")
	resp, err := app.Test(req)
	require.NoError(t, err)
	assert.Equal(t, 400, resp.StatusCode) // Bad Request - Invalid CSRF token
}

func TestCSRFMiddleware_ExcludedPaths(t *testing.T) {
	app := fiber.New()
	config := DefaultCSRFConfig()
	config.ExcludedPaths = []string{"/api/health"}
	app.Use(CSRF(config))

	app.Post("/api/health", func(c *fiber.Ctx) error {
		return c.SendStatus(200)
	})

	req := getTestRequest("POST", "/api/health", nil)
	resp, err := app.Test(req)
	require.NoError(t, err)
	assert.Equal(t, 200, resp.StatusCode) // Should succeed without CSRF token
}

func TestIsExcludedPath(t *testing.T) {
	tests := []struct {
		path          string
		excludedPaths []string
		expected      bool
	}{
		{"/api/health", []string{"/api/health"}, true},
		{"/api/health/status", []string{"/api/health"}, true},
		{"/api/users", []string{"/api/health"}, false},
		{"/api/metrics", []string{"/api/health", "/api/metrics"}, true},
	}

	for _, tt := range tests {
		t.Run(tt.path, func(t *testing.T) {
			result := isExcludedPath(tt.path, tt.excludedPaths)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestGetSameSite(t *testing.T) {
	tests := []struct {
		sameSite string
		expected fiber.CookieSameSite
	}{
		{"strict", fiber.CookieSameSiteStrictMode},
		{"Strict", fiber.CookieSameSiteStrictMode},
		{"lax", fiber.CookieSameSiteLaxMode},
		{"Lax", fiber.CookieSameSiteLaxMode},
		{"none", fiber.CookieSameSiteNoneMode},
		{"None", fiber.CookieSameSiteNoneMode},
		{"invalid", fiber.CookieSameSiteStrictMode},
	}

	for _, tt := range tests {
		t.Run(tt.sameSite, func(t *testing.T) {
			result := getSameSite(tt.sameSite)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestCSRFTokenHandler(t *testing.T) {
	app := fiber.New()
	config := DefaultCSRFConfig()
	manager := NewCSRFManager(config)
	defer manager.Stop()

	app.Get("/csrf-token", CSRFTokenHandler(manager, config))

	req := getTestRequest("GET", "/csrf-token", nil)
	resp, err := app.Test(req)
	require.NoError(t, err)
	assert.Equal(t, 200, resp.StatusCode)

	// Check for CSRF cookie
	cookies := resp.Cookies()
	var csrfCookie *fiber.Cookie
	for _, cookie := range cookies {
		if cookie.Name == config.CookieName {
			csrfCookie = cookie
			break
		}
	}
	assert.NotNil(t, csrfCookie)
	assert.NotEmpty(t, csrfCookie.Value)
}

// Helper function to create test requests
func getTestRequest(method, path string, body []byte) *http.Request {
	req := httptest.NewRequest(method, path, bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	return req
}
