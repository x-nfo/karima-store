package middleware

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDefaultCSRFConfig(t *testing.T) {
	config := DefaultCSRFConfig()

	assert.Equal(t, "header:X-CSRF-Token", config.KeyLookup)
	assert.Equal(t, "csrf_token", config.CookieName)
	assert.True(t, config.CookieSecure)
	assert.False(t, config.CookieHTTPOnly)
	assert.Equal(t, "Strict", config.CookieSameSite)
	assert.Equal(t, 24*60*60, config.Expiration) // 24 hours in seconds
	assert.Equal(t, "token", config.ContextKey)
}

func TestCSRFMiddleware_GET(t *testing.T) {
	app := fiber.New()
	config := DefaultCSRFConfig()
	app.Use(CSRF(config))

	app.Get("/test", func(c *fiber.Ctx) error {
		token := c.Locals(config.ContextKey)
		return c.JSON(fiber.Map{"token": token})
	})

	req := getCSRFTestRequest("GET", "/test", nil)
	resp, err := app.Test(req)
	require.NoError(t, err)
	assert.Equal(t, 200, resp.StatusCode)

	// Check for CSRF cookie
	cookies := resp.Cookies()
	var csrfCookie *http.Cookie
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

	req := getCSRFTestRequest("POST", "/test", nil)
	resp, err := app.Test(req)
	require.NoError(t, err)
	assert.Equal(t, 403, resp.StatusCode) // Forbidden - CSRF token required
}

func TestCSRFMiddleware_POST_WithInvalidToken(t *testing.T) {
	app := fiber.New()
	config := DefaultCSRFConfig()
	app.Use(CSRF(config))

	app.Post("/test", func(c *fiber.Ctx) error {
		return c.SendStatus(200)
	})

	req := getCSRFTestRequest("POST", "/test", nil)
	req.Header.Set("X-CSRF-Token", "invalid-token")
	resp, err := app.Test(req)
	require.NoError(t, err)
	assert.Equal(t, 403, resp.StatusCode) // Forbidden - Invalid CSRF token
}

func TestCSRFMiddleware_POST_WithValidToken(t *testing.T) {
	app := fiber.New()
	config := DefaultCSRFConfig()
	app.Use(CSRF(config))

	app.Get("/test", func(c *fiber.Ctx) error {
		token := c.Locals(config.ContextKey)
		return c.JSON(fiber.Map{"token": token})
	})

	app.Post("/test", func(c *fiber.Ctx) error {
		return c.SendStatus(200)
	})

	// First, get a token via GET request
	getReq := getCSRFTestRequest("GET", "/test", nil)
	getResp, err := app.Test(getReq)
	require.NoError(t, err)
	assert.Equal(t, 200, getResp.StatusCode)

	// Extract CSRF token from cookie
	var csrfToken string
	for _, cookie := range getResp.Cookies() {
		if cookie.Name == config.CookieName {
			csrfToken = cookie.Value
			break
		}
	}
	require.NotEmpty(t, csrfToken)

	// Now make POST request with valid token
	postReq := getCSRFTestRequest("POST", "/test", nil)
	postReq.Header.Set("X-CSRF-Token", csrfToken)
	postReq.AddCookie(&http.Cookie{
		Name:  config.CookieName,
		Value: csrfToken,
	})
	postResp, err := app.Test(postReq)
	require.NoError(t, err)
	assert.Equal(t, 200, postResp.StatusCode)
}

func TestCSRFMiddleware_NextFunction(t *testing.T) {
	app := fiber.New()
	config := DefaultCSRFConfig()
	config.Next = func(c *fiber.Ctx) bool {
		return c.Path() == "/skip-csrf"
	}
	app.Use(CSRF(config))

	app.Post("/skip-csrf", func(c *fiber.Ctx) error {
		return c.SendStatus(200)
	})

	app.Post("/require-csrf", func(c *fiber.Ctx) error {
		return c.SendStatus(200)
	})

	// Should succeed without CSRF token
	req1 := getCSRFTestRequest("POST", "/skip-csrf", nil)
	resp1, err := app.Test(req1)
	require.NoError(t, err)
	assert.Equal(t, 200, resp1.StatusCode)

	// Should fail without CSRF token
	req2 := getCSRFTestRequest("POST", "/require-csrf", nil)
	resp2, err := app.Test(req2)
	require.NoError(t, err)
	assert.Equal(t, 403, resp2.StatusCode)
}

func TestCSRFMiddleware_CustomConfig(t *testing.T) {
	app := fiber.New()
	config := CSRFConfig{
		KeyLookup:      "form:csrf_token",
		CookieName:     "custom_csrf",
		CookieSecure:   false,
		CookieHTTPOnly: true,
		CookieSameSite: "Lax",
		Expiration:     3600, // 1 hour
		ContextKey:     "custom_token",
	}
	app.Use(CSRF(config))

	app.Get("/test", func(c *fiber.Ctx) error {
		token := c.Locals(config.ContextKey)
		return c.JSON(fiber.Map{"token": token})
	})

	req := getCSRFTestRequest("GET", "/test", nil)
	resp, err := app.Test(req)
	require.NoError(t, err)
	assert.Equal(t, 200, resp.StatusCode)

	// Check for custom CSRF cookie
	cookies := resp.Cookies()
	var csrfCookie *http.Cookie
	for _, cookie := range cookies {
		if cookie.Name == config.CookieName {
			csrfCookie = cookie
			break
		}
	}
	assert.NotNil(t, csrfCookie)
	assert.NotEmpty(t, csrfCookie.Value)
	assert.True(t, csrfCookie.HttpOnly)
}

func TestCSRFMiddleware(t *testing.T) {
	app := fiber.New()
	app.Use(CSRFMiddleware())

	app.Get("/test", func(c *fiber.Ctx) error {
		return c.SendStatus(200)
	})

	req := getCSRFTestRequest("GET", "/test", nil)
	resp, err := app.Test(req)
	require.NoError(t, err)
	assert.Equal(t, 200, resp.StatusCode)
}

func TestGetSameSite(t *testing.T) {
	tests := []struct {
		sameSite string
		expected string
	}{
		{"Strict", "Strict"},
		{"Lax", "Lax"},
		{"None", "None"},
		{"invalid", "Strict"},
	}

	for _, tt := range tests {
		t.Run(tt.sameSite, func(t *testing.T) {
			result := getSameSite(tt.sameSite)
			assert.Equal(t, tt.expected, result)
		})
	}
}

// Helper function to create test requests
func getCSRFTestRequest(method, path string, body []byte) *http.Request {
	req := httptest.NewRequest(method, path, bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	return req
}
