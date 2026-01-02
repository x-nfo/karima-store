package utils

import (
	"io"
	"net/http/httptest"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
)

func TestSendSuccess(t *testing.T) {
	app := fiber.New()

	app.Get("/success", func(c *fiber.Ctx) error {
		data := map[string]string{
			"message": "Operation successful",
		}
		return SendSuccess(c, data, "Success message")
	})

	req := httptest.NewRequest("GET", "/success", nil)
	resp, err := app.Test(req)
	assert.NoError(t, err)

	assert.Equal(t, fiber.StatusOK, resp.StatusCode)
	bodyBytes, _ := io.ReadAll(resp.Body)
	assert.Contains(t, string(bodyBytes), "success")
	assert.Contains(t, string(bodyBytes), "Success message")
	assert.Contains(t, string(bodyBytes), "Operation successful")
}

func TestSendSuccess_CustomStatus(t *testing.T) {
	app := fiber.New()

	app.Post("/created", func(c *fiber.Ctx) error {
		data := map[string]string{
			"id": "123",
		}
		return SendSuccess(c, data, "Resource created", fiber.StatusCreated)
	})

	req := httptest.NewRequest("POST", "/created", nil)
	resp, err := app.Test(req)
	assert.NoError(t, err)

	assert.Equal(t, fiber.StatusCreated, resp.StatusCode)
	bodyBytes, _ := io.ReadAll(resp.Body)
	assert.Contains(t, string(bodyBytes), "success")
	assert.Contains(t, string(bodyBytes), "Resource created")
	assert.Contains(t, string(bodyBytes), "123")
}

func TestSendError(t *testing.T) {
	app := fiber.New()

	app.Get("/error", func(c *fiber.Ctx) error {
		errors := map[string]string{
			"field": "Invalid value",
		}
		return SendError(c, fiber.StatusBadRequest, "Bad request", errors)
	})

	req := httptest.NewRequest("GET", "/error", nil)
	resp, err := app.Test(req)
	assert.NoError(t, err)

	assert.Equal(t, fiber.StatusBadRequest, resp.StatusCode)
	bodyBytes, _ := io.ReadAll(resp.Body)
	assert.Contains(t, string(bodyBytes), "error")
	assert.Contains(t, string(bodyBytes), "Bad request")
	assert.Contains(t, string(bodyBytes), "Invalid value")
}

func TestSendValidationError(t *testing.T) {
	app := fiber.New()

	app.Post("/validate", func(c *fiber.Ctx) error {
		errors := map[string]string{
			"email":    "Invalid email format",
			"password": "Password too short",
		}
		return SendValidationError(c, errors)
	})

	req := httptest.NewRequest("POST", "/validate", nil)
	resp, err := app.Test(req)
	assert.NoError(t, err)

	assert.Equal(t, fiber.StatusBadRequest, resp.StatusCode)
	bodyBytes, _ := io.ReadAll(resp.Body)
	assert.Contains(t, string(bodyBytes), "error")
	assert.Contains(t, string(bodyBytes), "Validation failed")
	assert.Contains(t, string(bodyBytes), "Invalid email format")
	assert.Contains(t, string(bodyBytes), "Password too short")
}

func TestSendCreated(t *testing.T) {
	app := fiber.New()

	app.Post("/create", func(c *fiber.Ctx) error {
		data := map[string]interface{}{
			"id":      123,
			"name":    "Test Product",
			"created": true,
		}
		return SendCreated(c, data, "Product created successfully")
	})

	req := httptest.NewRequest("POST", "/create", nil)
	resp, err := app.Test(req)
	assert.NoError(t, err)

	assert.Equal(t, fiber.StatusCreated, resp.StatusCode)
	bodyBytes, _ := io.ReadAll(resp.Body)
	assert.Contains(t, string(bodyBytes), "success")
	assert.Contains(t, string(bodyBytes), "Product created successfully")
	assert.Contains(t, string(bodyBytes), "123")
	assert.Contains(t, string(bodyBytes), "Test Product")
}

func TestErrorHandling_GenericError(t *testing.T) {
	app := fiber.New()

	app.Get("/generic-error", func(c *fiber.Ctx) error {
		return SendError(c, fiber.StatusInternalServerError, "Internal server error", nil)
	})

	req := httptest.NewRequest("GET", "/generic-error", nil)
	resp, err := app.Test(req)
	assert.NoError(t, err)

	assert.Equal(t, fiber.StatusInternalServerError, resp.StatusCode)
	bodyBytes, _ := io.ReadAll(resp.Body)
	assert.Contains(t, string(bodyBytes), "error")
	assert.Contains(t, string(bodyBytes), "Internal server error")
}

func TestErrorHandling_ValidationError(t *testing.T) {
	app := fiber.New()

	app.Post("/validation-error", func(c *fiber.Ctx) error {
		errors := []string{
			"Name is required",
			"Price must be greater than 0",
		}
		return SendError(c, fiber.StatusBadRequest, "Validation failed", errors)
	})

	req := httptest.NewRequest("POST", "/validation-error", nil)
	resp, err := app.Test(req)
	assert.NoError(t, err)

	assert.Equal(t, fiber.StatusBadRequest, resp.StatusCode)
	bodyBytes, _ := io.ReadAll(resp.Body)
	assert.Contains(t, string(bodyBytes), "error")
	assert.Contains(t, string(bodyBytes), "Validation failed")
	assert.Contains(t, string(bodyBytes), "Name is required")
	assert.Contains(t, string(bodyBytes), "Price must be greater than 0")
}

func TestErrorHandling_AuthenticationError(t *testing.T) {
	app := fiber.New()

	app.Get("/auth-error", func(c *fiber.Ctx) error {
		return SendError(c, fiber.StatusUnauthorized, "Unauthorized access", nil)
	})

	req := httptest.NewRequest("GET", "/auth-error", nil)
	resp, err := app.Test(req)
	assert.NoError(t, err)

	assert.Equal(t, fiber.StatusUnauthorized, resp.StatusCode)
	bodyBytes, _ := io.ReadAll(resp.Body)
	assert.Contains(t, string(bodyBytes), "error")
	assert.Contains(t, string(bodyBytes), "Unauthorized access")
}

func TestErrorHandling_AuthorizationError(t *testing.T) {
	app := fiber.New()

	app.Get("/forbidden", func(c *fiber.Ctx) error {
		return SendError(c, fiber.StatusForbidden, "Access forbidden", nil)
	})

	req := httptest.NewRequest("GET", "/forbidden", nil)
	resp, err := app.Test(req)
	assert.NoError(t, err)

	assert.Equal(t, fiber.StatusForbidden, resp.StatusCode)
	bodyBytes, _ := io.ReadAll(resp.Body)
	assert.Contains(t, string(bodyBytes), "error")
	assert.Contains(t, string(bodyBytes), "Access forbidden")
}

func TestErrorHandling_NotFoundError(t *testing.T) {
	app := fiber.New()

	app.Get("/not-found", func(c *fiber.Ctx) error {
		return SendError(c, fiber.StatusNotFound, "Resource not found", nil)
	})

	req := httptest.NewRequest("GET", "/not-found", nil)
	resp, err := app.Test(req)
	assert.NoError(t, err)

	assert.Equal(t, fiber.StatusNotFound, resp.StatusCode)
	bodyBytes, _ := io.ReadAll(resp.Body)
	assert.Contains(t, string(bodyBytes), "error")
	assert.Contains(t, string(bodyBytes), "Resource not found")
}

func TestSecurityErrorMessages_NoSensitiveInfo(t *testing.T) {
	app := fiber.New()

	app.Get("/secure-error", func(c *fiber.Ctx) error {
		// Ensure no sensitive information is leaked in error messages
		return SendError(c, fiber.StatusInternalServerError, "An error occurred", nil)
	})

	req := httptest.NewRequest("GET", "/secure-error", nil)
	resp, err := app.Test(req)
	assert.NoError(t, err)

	assert.Equal(t, fiber.StatusInternalServerError, resp.StatusCode)
	bodyBytes, _ := io.ReadAll(resp.Body)
	assert.Contains(t, string(bodyBytes), "error")
	assert.Contains(t, string(bodyBytes), "An error occurred")
	// Ensure no sensitive information like file paths, stack traces, etc.
	assert.NotContains(t, string(bodyBytes), "/")
	assert.NotContains(t, string(bodyBytes), "stack")
	assert.NotContains(t, string(bodyBytes), "trace")
}

func TestSecurityErrorMessages_ConsistentFormatting(t *testing.T) {
	app := fiber.New()

	// Test that all error responses have consistent format
	errorCodes := []int{
		fiber.StatusBadRequest,
		fiber.StatusUnauthorized,
		fiber.StatusForbidden,
		fiber.StatusNotFound,
		fiber.StatusInternalServerError,
	}

	for _, code := range errorCodes {
		app.Get("/error/"+string(rune(code)), func(c *fiber.Ctx) error {
			return SendError(c, code, "Test error", nil)
		})

		req := httptest.NewRequest("GET", "/error/"+string(rune(code)), nil)
		resp, err := app.Test(req)
		assert.NoError(t, err)

		assert.Equal(t, code, resp.StatusCode)
		bodyBytes, _ := io.ReadAll(resp.Body)
		assert.Contains(t, string(bodyBytes), "status")
		assert.Contains(t, string(bodyBytes), "message")
	}
}

func TestSecurityErrorMessages_DatabaseError(t *testing.T) {
	app := fiber.New()

	app.Get("/db-error", func(c *fiber.Ctx) error {
		// Simulate database error without exposing sensitive details
		return SendError(c, fiber.StatusInternalServerError, "Database operation failed", nil)
	})

	req := httptest.NewRequest("GET", "/db-error", nil)
	resp, err := app.Test(req)
	assert.NoError(t, err)

	assert.Equal(t, fiber.StatusInternalServerError, resp.StatusCode)
	bodyBytes, _ := io.ReadAll(resp.Body)
	assert.Contains(t, string(bodyBytes), "error")
	assert.Contains(t, string(bodyBytes), "Database operation failed")
	// Ensure no database details are exposed
	assert.NotContains(t, string(bodyBytes), "sql")
	assert.NotContains(t, string(bodyBytes), "table")
	assert.NotContains(t, string(bodyBytes), "column")
}

func TestSecurityErrorMessages_AuthenticationError(t *testing.T) {
	app := fiber.New()

	app.Get("/auth-error", func(c *fiber.Ctx) error {
		// Generic authentication error message
		return SendError(c, fiber.StatusUnauthorized, "Authentication required", nil)
	})

	req := httptest.NewRequest("GET", "/auth-error", nil)
	resp, err := app.Test(req)
	assert.NoError(t, err)

	assert.Equal(t, fiber.StatusUnauthorized, resp.StatusCode)
	bodyBytes, _ := io.ReadAll(resp.Body)
	assert.Contains(t, string(bodyBytes), "error")
	assert.Contains(t, string(bodyBytes), "Authentication required")
	// Ensure no authentication details are exposed
	assert.NotContains(t, string(bodyBytes), "token")
	assert.NotContains(t, string(bodyBytes), "password")
	assert.NotContains(t, string(bodyBytes), "session")
}

func TestErrorHandling_ErrorCodeConsistency(t *testing.T) {
	app := fiber.New()

	// Test that error codes are consistent with HTTP standards
	tests := []struct {
		name           string
		statusCode     int
		expectedStatus int
	}{
		{
			name:           "Bad request",
			statusCode:     fiber.StatusBadRequest,
			expectedStatus: 400,
		},
		{
			name:           "Unauthorized",
			statusCode:     fiber.StatusUnauthorized,
			expectedStatus: 401,
		},
		{
			name:           "Forbidden",
			statusCode:     fiber.StatusForbidden,
			expectedStatus: 403,
		},
		{
			name:           "Not found",
			statusCode:     fiber.StatusNotFound,
			expectedStatus: 404,
		},
		{
			name:           "Internal server error",
			statusCode:     fiber.StatusInternalServerError,
			expectedStatus: 500,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			app.Get("/test-error", func(c *fiber.Ctx) error {
				return SendError(c, tt.statusCode, "Test error", nil)
			})

			req := httptest.NewRequest("GET", "/test-error", nil)
			resp, err := app.Test(req)
			assert.NoError(t, err)

			assert.Equal(t, tt.expectedStatus, resp.StatusCode)
		})
	}
}

func TestErrorHandling_ErrorDetailLevels(t *testing.T) {
	app := fiber.New()

	// Test that error details are appropriate for different environments
	app.Get("/error-dev", func(c *fiber.Ctx) error {
		// In development, more details might be shown
		errors := map[string]string{
			"field": "Invalid value",
		}
		return SendError(c, fiber.StatusBadRequest, "Validation failed", errors)
	})

	req1 := httptest.NewRequest("GET", "/error-dev", nil)
	resp1, err1 := app.Test(req1)
	assert.NoError(t, err1)

	assert.Equal(t, fiber.StatusBadRequest, resp1.StatusCode)
	bodyBytes1, _ := io.ReadAll(resp1.Body)
	assert.Contains(t, string(bodyBytes1), "Validation failed")
	assert.Contains(t, string(bodyBytes1), "Invalid value")

	// In production, less details should be shown
	app.Get("/error-prod", func(c *fiber.Ctx) error {
		return SendError(c, fiber.StatusBadRequest, "Validation failed", nil)
	})

	req2 := httptest.NewRequest("GET", "/error-prod", nil)
	resp2, err2 := app.Test(req2)
	assert.NoError(t, err2)

	assert.Equal(t, fiber.StatusBadRequest, resp2.StatusCode)
	bodyBytes2, _ := io.ReadAll(resp2.Body)
	assert.Contains(t, string(bodyBytes2), "Validation failed")
}

func TestAPIResponse_Structure(t *testing.T) {
	app := fiber.New()

	app.Get("/test", func(c *fiber.Ctx) error {
		return SendSuccess(c, map[string]string{"key": "value"}, "Success")
	})

	req := httptest.NewRequest("GET", "/test", nil)
	resp, err := app.Test(req)
	assert.NoError(t, err)

	// Verify response structure
	bodyBytes, _ := io.ReadAll(resp.Body)
	assert.Contains(t, string(bodyBytes), `"status"`)
	assert.Contains(t, string(bodyBytes), `"message"`)
	assert.Contains(t, string(bodyBytes), `"data"`)
}