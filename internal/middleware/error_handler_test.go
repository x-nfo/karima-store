package middleware

import (
	"errors"
	"net/http/httptest"
	"testing"

	apperrors "karima_store/internal/errors"

	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
)

func TestErrorHandler(t *testing.T) {
	app := fiber.New()
	app.Use(ErrorHandler())

	app.Get("/test", func(c *fiber.Ctx) error {
		return apperrors.NewNotFoundError("resource")
	})

	req := httptest.NewRequest("GET", "/test", nil)
	resp, err := app.Test(req)
	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusNotFound, resp.StatusCode)
}

func TestErrorHandlerWithAppError(t *testing.T) {
	app := fiber.New()
	app.Use(ErrorHandler())

	app.Get("/validation", func(c *fiber.Ctx) error {
		return apperrors.NewValidationError("Invalid input")
	})

	req := httptest.NewRequest("GET", "/validation", nil)
	resp, err := app.Test(req)
	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusBadRequest, resp.StatusCode)
}

func TestErrorHandlerWithUnauthorized(t *testing.T) {
	app := fiber.New()
	app.Use(ErrorHandler())

	app.Get("/unauthorized", func(c *fiber.Ctx) error {
		return apperrors.NewUnauthorizedError("Access denied")
	})

	req := httptest.NewRequest("GET", "/unauthorized", nil)
	resp, err := app.Test(req)
	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusUnauthorized, resp.StatusCode)
}

func TestErrorHandlerWithInternalError(t *testing.T) {
	app := fiber.New()
	app.Use(ErrorHandler())

	app.Get("/internal", func(c *fiber.Ctx) error {
		return apperrors.NewInternalError("Something went wrong")
	})

	req := httptest.NewRequest("GET", "/internal", nil)
	resp, err := app.Test(req)
	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusInternalServerError, resp.StatusCode)
}

func TestErrorHandlerWithFiberError(t *testing.T) {
	app := fiber.New()
	app.Use(ErrorHandler())

	app.Get("/fiber-error", func(c *fiber.Ctx) error {
		return fiber.NewError(fiber.StatusTeapot, "I'm a teapot")
	})

	req := httptest.NewRequest("GET", "/fiber-error", nil)
	resp, err := app.Test(req)
	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusTeapot, resp.StatusCode)
}

func TestRecoverMiddleware(t *testing.T) {
	app := fiber.New()
	app.Use(RecoverMiddleware())

	app.Get("/panic", func(c *fiber.Ctx) error {
		panic("test panic")
	})

	req := httptest.NewRequest("GET", "/panic", nil)
	resp, err := app.Test(req)
	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusInternalServerError, resp.StatusCode)
}

func TestErrorHandlerWithDetails(t *testing.T) {
	app := fiber.New()
	app.Use(ErrorHandler())

	app.Get("/details", func(c *fiber.Ctx) error {
		details := []apperrors.ValidationErrorDetail{
			{Field: "email", Message: "Invalid email format"},
			{Field: "password", Message: "Password too short"},
		}
		return apperrors.NewValidationErrorFromDetails(details)
	})

	req := httptest.NewRequest("GET", "/details", nil)
	resp, err := app.Test(req)
	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusBadRequest, resp.StatusCode)
}

func TestSanitizeErrorMessage(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "Normal message",
			input:    "This is a normal error message",
			expected: "This is a normal error message",
		},
		{
			name:     "Password in message",
			input:    "Invalid password: secret123",
			expected: "Sensitive information redacted",
		},
		{
			name:     "Token in message",
			input:    "Invalid token: abc123def456",
			expected: "Sensitive information redacted",
		},
		{
			name:     "Secret key in message",
			input:    "API secret key: xyz789",
			expected: "Sensitive information redacted",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := SanitizeErrorMessage(tt.input)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestErrorHandlerWithWrappedError(t *testing.T) {
	app := fiber.New()
	app.Use(ErrorHandler())

	app.Get("/wrapped", func(c *fiber.Ctx) error {
		originalErr := errors.New("original error")
		return apperrors.WrapError(apperrors.ErrCodeDatabase, "Database operation failed", originalErr)
	})

	req := httptest.NewRequest("GET", "/wrapped", nil)
	resp, err := app.Test(req)
	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusInternalServerError, resp.StatusCode)
}

func TestErrorHandlerNoError(t *testing.T) {
	app := fiber.New()
	app.Use(ErrorHandler())

	app.Get("/success", func(c *fiber.Ctx) error {
		return c.SendString("success")
	})

	req := httptest.NewRequest("GET", "/success", nil)
	resp, err := app.Test(req)
	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusOK, resp.StatusCode)
}
