package middleware

import (
	"io"
	"net/http/httptest"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
)

func TestSecurityHeaders(t *testing.T) {
	app := fiber.New()
	app.Use(SecurityHeaders())

	app.Get("/test", func(c *fiber.Ctx) error {
		return c.SendString("test")
	})

	req := httptest.NewRequest("GET", "/test", nil)
	resp, err := app.Test(req)
	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusOK, resp.StatusCode)

	// Check security headers
	headers := map[string]string{
		"Content-Security-Policy":      "default-src 'self'; script-src 'self' 'unsafe-inline' 'unsafe-eval'; style-src 'self' 'unsafe-inline'; img-src 'self' data: https:; font-src 'self' data:; connect-src 'self'; frame-ancestors 'none'; base-uri 'self'; form-action 'self'",
		"X-Content-Type-Options":       "nosniff",
		"X-Frame-Options":              "DENY",
		"X-XSS-Protection":             "1; mode=block",
		"Strict-Transport-Security":    "max-age=31536000; includeSubDomains; preload",
		"Referrer-Policy":              "strict-origin-when-cross-origin",
		"X-DNS-Prefetch-Control":       "off",
		"Cross-Origin-Embedder-Policy": "require-corp",
		"Cross-Origin-Opener-Policy":   "same-origin",
		"Cross-Origin-Resource-Policy": "same-origin",
	}

	for header, expectedValue := range headers {
		t.Run(header, func(t *testing.T) {
			actualValue := resp.Header.Get(header)
			assert.Equal(t, expectedValue, actualValue, "Header %s should match", header)
		})
	}
}

func TestSecurityHeadersDevelopment(t *testing.T) {
	app := fiber.New()
	app.Use(SecurityHeadersDevelopment())

	app.Get("/test", func(c *fiber.Ctx) error {
		return c.SendString("test")
	})

	req := httptest.NewRequest("GET", "/test", nil)
	resp, err := app.Test(req)
	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusOK, resp.StatusCode)

	// Check security headers
	assert.Equal(t, "nosniff", resp.Header.Get("X-Content-Type-Options"))
	assert.Equal(t, "SAMEORIGIN", resp.Header.Get("X-Frame-Options"))
	assert.Equal(t, "1; mode=block", resp.Header.Get("X-XSS-Protection"))

	// HSTS should not be set in development
	assert.Empty(t, resp.Header.Get("Strict-Transport-Security"))
}

func TestSecurityHeadersChain(t *testing.T) {
	app := fiber.New()
	app.Use(SecurityHeaders())
	app.Use(func(c *fiber.Ctx) error {
		return c.SendString("test response")
	})

	req := httptest.NewRequest("GET", "/test", nil)
	resp, err := app.Test(req)
	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusOK, resp.StatusCode)

	// Read response body
	body, err := io.ReadAll(resp.Body)
	assert.NoError(t, err)
	assert.Equal(t, "test response", string(body))
}
