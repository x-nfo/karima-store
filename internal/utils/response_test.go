package utils

import (
	"net/http"
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
	resp := httptest.NewRecorder()
	app.Handler()(resp, req)

	assert.Equal(t, fiber.StatusOK, resp.Code)
	assert.Contains(t, resp.Body.String(), "success")
	assert.Contains(t, resp.Body.String(), "Success message")
	assert.Contains(t, resp.Body.String(), "Operation successful")
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
	resp := httptest.NewRecorder()
	app.Handler()(resp, req)

	assert.Equal(t, fiber.StatusCreated, resp.Code)
	assert.Contains(t, resp.Body.String(), "success")
	assert.Contains(t, resp.Body.String(), "Resource created")
	assert.Contains(t, resp.Body.String(), "123")
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
	resp := httptest.NewRecorder()
	app.Handler()(resp, req)

	assert.Equal(t, fiber.StatusBadRequest, resp.Code)
	assert.Contains(t, resp.Body.String(), "error")
	assert.Contains(t, resp.Body.String(), "Bad request")
	assert.Contains(t, resp.Body.String(), "Invalid value")
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
	resp := httptest.NewRecorder()
	app.Handler()(resp, req)

	assert.Equal(t, fiber.StatusBadRequest, resp.Code)
	assert.Contains(t, resp.Body.String(), "error")
	assert.Contains(t, resp.Body.String(), "Validation failed")
	assert.Contains(t, resp.Body.String(), "Invalid email format")
	assert.Contains(t, resp.Body.String(), "Password too short")
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
	resp := httptest.NewRecorder()
	app.Handler()(resp, req)

	assert.Equal(t, fiber.StatusCreated, resp.Code)
	assert.Contains(t, resp.Body.String(), "success")
	assert.Contains(t, resp.Body.String(), "Product created successfully")
	assert.Contains(t, resp.Body.String(), "123")
	assert.Contains(t, resp.Body.String(), "Test Product")
}

func TestErrorHandling_GenericError(t *testing.T) {
	app := fiber.New()

	app.Get("/generic-error", func(c *fiber.Ctx) error {
		return SendError(c, fiber.StatusInternalServerError, "Internal server error", nil)
	})

	req := httptest.NewRequest("GET", "/generic-error", nil)
	resp := httptest.NewRecorder()
	app.Handler()(resp, req)

	assert.Equal(t, fiber.StatusInternalServerError, resp.Code)
	assert.Contains(t, resp.Body.String(), "error")
	assert.Contains(t, resp.Body.String(), "Internal server error")
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
	resp := httptest.NewRecorder()
	app.Handler()(resp, req)

	assert.Equal(t, fiber.StatusBadRequest, resp.Code)
	assert.Contains(t, resp.Body.String(), "error")
	assert.Contains(t, resp.Body.String(), "Validation failed")
	assert.Contains(t, resp.Body.String(), "Name is required")
	assert.Contains(t, resp.Body.String(), "Price must be greater than 0")
}

func TestErrorHandling_AuthenticationError(t *testing.T) {
	app := fiber.New()

	app.Get("/auth-error", func(c *fiber.Ctx) error {
		return SendError(c, fiber.StatusUnauthorized, "Unauthorized access", nil)
	})

	req := httptest.NewRequest("GET", "/auth-error", nil)
	resp := httptest.NewRecorder()
	app.Handler()(resp, req)

	assert.Equal(t, fiber.StatusUnauthorized, resp.Code)
	assert.Contains(t, resp.Body.String(), "error")
	assert.Contains(t, resp.Body.String(), "Unauthorized access")
}

func TestErrorHandling_AuthorizationError(t *testing.T) {
	app := fiber.New()

	app.Get("/forbidden", func(c *fiber.Ctx) error {
		return SendError(c, fiber.StatusForbidden, "Access forbidden", nil)
	})

	req := httptest.NewRequest("GET", "/forbidden", nil)
	resp := httptest.NewRecorder()
	app.Handler()(resp, req)

	assert.Equal(t, fiber.StatusForbidden, resp.Code)
	assert.Contains(t, resp.Body.String(), "error")
	assert.Contains(t, resp.Body.String(), "Access forbidden")
}

func TestErrorHandling_NotFoundError(t *testing.T) {
	app := fiber.New()

	app.Get("/not-found", func(c *fiber.Ctx) error {
		return SendError(c, fiber.StatusNotFound, "Resource not found", nil)
	})

	req := httptest.NewRequest("GET", "/not-found", nil)
	resp := httptest.NewRecorder()
	app.Handler()(resp, req)

	assert.Equal(t, fiber.StatusNotFound, resp.Code)
	assert.Contains(t, resp.Body.String(), "error")
	assert.Contains(t, resp.Body.String(), "Resource not found")
}

func TestSecurityErrorMessages_NoSensitiveInfo(t *testing.T) {
	app := fiber.New()

	app.Get("/secure-error", func(c *fiber.Ctx) error {
		// Ensure no sensitive information is leaked in error messages
		return SendError(c, fiber.StatusInternalServerError, "An error occurred", nil)
	})

	req := httptest.NewRequest("GET", "/secure-error", nil)
	resp := httptest.NewRecorder()
	app.Handler()(resp, req)

	assert.Equal(t, fiber.StatusInternalServerError, resp.Code)
	assert.Contains(t, resp.Body.String(), "error")
	assert.Contains(t, resp.Body.String(), "An error occurred")
	// Ensure no sensitive information like file paths, stack traces, etc.
	assert.NotContains(t, resp.Body.String(), "/")
	assert.NotContains(t, resp.Body.String(), "stack")
	assert.NotContains(t, resp.Body.String(), "trace")
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
		resp := httptest.NewRecorder()
		app.Handler()(resp, req)

		assert.Equal(t, code, resp.Code)
		assert.Contains(t, resp.Body.String(), "status")
		assert.Contains(t, resp.Body.String(), "message")
	}
}

func TestSecurityErrorMessages_DatabaseError(t *testing.T) {
	app := fiber.New()

	app.Get("/db-error", func(c *fiber.Ctx) error {
		// Simulate database error without exposing sensitive details
		return SendError(c, fiber.StatusInternalServerError, "Database operation failed", nil)
	})

	req := httptest.NewRequest("GET", "/db-error", nil)
	resp := httptest.NewRecorder()
	app.Handler()(resp, req)

	assert.Equal(t, fiber.StatusInternalServerError, resp.Code)
	assert.Contains(t, resp.Body.String(), "error")
	assert.Contains(t, resp.Body.String(), "Database operation failed")
	// Ensure no database details are exposed
	assert.NotContains(t, resp.Body.String(), "sql")
	assert.NotContains(t, resp.Body.String(), "table")
	assert.NotContains(t, resp.Body.String(), "column")
}

func TestSecurityErrorMessages_AuthenticationError(t *testing.T) {
	app := fiber.New()

	app.Get("/auth-error", func(c *fiber.Ctx) error {
		// Generic authentication error message
		return SendError(c, fiber.StatusUnauthorized, "Authentication required", nil)
	})

	req := httptest.NewRequest("GET", "/auth-error", nil)
	resp := httptest.NewRecorder()
	app.Handler()(resp, req)

	assert.Equal(t, fiber.StatusUnauthorized, resp.Code)
	assert.Contains(t, resp.Body.String(), "error")
	assert.Contains(t, resp.Body.String(), "Authentication required")
	// Ensure no authentication details are exposed
	assert.NotContains(t, resp.Body.String(), "token")
	assert.NotContains(t, resp.Body.String(), "password")
	assert.NotContains(t, resp.Body.String(), "session")
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
			resp := httptest.NewRecorder()
			app.Handler()(resp, req)

			assert.Equal(t, tt.expectedStatus, resp.Code)
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
	resp1 := httptest.NewRecorder()
	app.Handler()(resp1, req1)

	assert.Equal(t, fiber.StatusBadRequest, resp1.Code)
	assert.Contains(t, resp1.Body.String(), "Validation failed")
	assert.Contains(t, resp1.Body.String(), "Invalid value")

	// In production, less details should be shown
	app.Get("/error-prod", func(c *fiber.Ctx) error {
		return SendError(c, fiber.StatusBadRequest, "Validation failed", nil)
	})

	req2 := httptest.NewRequest("GET", "/error-prod", nil)
	resp2 := httptest.NewRecorder()
	app.Handler()(resp2, req2)

	assert.Equal(t, fiber.StatusBadRequest, resp2.Code)
	assert.Contains(t, resp2.Body.String(), "Validation failed")
}

func TestAPIResponse_Structure(t *testing.T) {
	app := fiber.New()

	app.Get("/test", func(c *fiber.Ctx) error {
		return SendSuccess(c, map[string]string{"key": "value"}, "Success")
	})

	req := httptest.NewRequest("GET", "/test", nil)
	resp := httptest.NewRecorder()
	app.Handler()(resp, req)

	// Verify response structure
	assert.Contains(t, resp.Body.String(), `"status"`)
	assert.Contains(t, resp.Body.String(), `"message"`)
	assert.Contains(t, resp.Body.String(), `"data"`)
}