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
			app := fiber.New()
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

func TestSendSuccess_NilData(t *testing.T) {
	app := fiber.New()

	app.Get("/nil-data", func(c *fiber.Ctx) error {
		return SendSuccess(c, nil, "Success with nil data")
	})

	req := httptest.NewRequest("GET", "/nil-data", nil)
	resp, err := app.Test(req)
	assert.NoError(t, err)

	assert.Equal(t, fiber.StatusOK, resp.StatusCode)
	bodyBytes, _ := io.ReadAll(resp.Body)
	assert.Contains(t, string(bodyBytes), "success")
	assert.Contains(t, string(bodyBytes), "Success with nil data")
	// When data is nil, the field is omitted due to omitempty tag
	assert.NotContains(t, string(bodyBytes), `"data"`)
}

func TestSendSuccess_EmptyDataStructures(t *testing.T) {
	tests := []struct {
		name     string
		data     interface{}
		expected string
	}{
		{
			name:     "Empty map",
			data:     map[string]interface{}{},
			expected: `"data":{}`,
		},
		{
			name:     "Empty slice",
			data:     []interface{}{},
			expected: `"data":[]`,
		},
		{
			name:     "Empty string",
			data:     "",
			expected: `"data":""`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			app := fiber.New()
			app.Get("/test", func(c *fiber.Ctx) error {
				return SendSuccess(c, tt.data, "Empty data test")
			})

			req := httptest.NewRequest("GET", "/test", nil)
			resp, err := app.Test(req)
			assert.NoError(t, err)

			assert.Equal(t, fiber.StatusOK, resp.StatusCode)
			bodyBytes, _ := io.ReadAll(resp.Body)
			assert.Contains(t, string(bodyBytes), tt.expected)
		})
	}
}

func TestSendSuccess_NestedData(t *testing.T) {
	app := fiber.New()

	app.Get("/nested", func(c *fiber.Ctx) error {
		data := map[string]interface{}{
			"user": map[string]interface{}{
				"id":   1,
				"name": "John",
				"address": map[string]string{
					"street": "123 Main St",
					"city":   "Jakarta",
				},
			},
			"orders": []map[string]interface{}{
				{
					"id":     101,
					"total":  100.50,
					"items":  []string{"item1", "item2"},
				},
			},
		}
		return SendSuccess(c, data, "Nested data test")
	})

	req := httptest.NewRequest("GET", "/nested", nil)
	resp, err := app.Test(req)
	assert.NoError(t, err)

	assert.Equal(t, fiber.StatusOK, resp.StatusCode)
	bodyBytes, _ := io.ReadAll(resp.Body)
	assert.Contains(t, string(bodyBytes), "success")
	assert.Contains(t, string(bodyBytes), "John")
	assert.Contains(t, string(bodyBytes), "Jakarta")
	assert.Contains(t, string(bodyBytes), "101")
}

func TestSendSuccess_UnicodeAndSpecialCharacters(t *testing.T) {
	app := fiber.New()

	app.Get("/unicode", func(c *fiber.Ctx) error {
		data := map[string]string{
			"arabic":    "مرحبا",
			"chinese":   "你好",
			"emoji":     "🎉🚀💻",
			"special":   "Test <script>alert('xss')</script>",
			"newline":   "Line1\nLine2\nLine3",
			"quotes":    `"quoted" and 'single'`,
		}
		return SendSuccess(c, data, "Unicode and special characters test")
	})

	req := httptest.NewRequest("GET", "/unicode", nil)
	resp, err := app.Test(req)
	assert.NoError(t, err)

	assert.Equal(t, fiber.StatusOK, resp.StatusCode)
	bodyBytes, _ := io.ReadAll(resp.Body)
	bodyStr := string(bodyBytes)

	// Check that data is properly JSON encoded
	assert.Contains(t, bodyStr, "مرحبا")
	assert.Contains(t, bodyStr, "你好")
	assert.Contains(t, bodyStr, "🎉")
	assert.Contains(t, bodyStr, "\\u003cscript\\u003ealert") // Script tags should be escaped (lowercase hex)
	assert.Contains(t, bodyStr, "Line1\\nLine2")
	assert.Contains(t, bodyStr, "\\\"quoted\\\"")
}

func TestSendError_NilErrors(t *testing.T) {
	app := fiber.New()

	app.Get("/error-nil", func(c *fiber.Ctx) error {
		return SendError(c, fiber.StatusBadRequest, "Error with nil errors", nil)
	})

	req := httptest.NewRequest("GET", "/error-nil", nil)
	resp, err := app.Test(req)
	assert.NoError(t, err)

	assert.Equal(t, fiber.StatusBadRequest, resp.StatusCode)
	bodyBytes, _ := io.ReadAll(resp.Body)
	assert.Contains(t, string(bodyBytes), "error")
	assert.Contains(t, string(bodyBytes), "Error with nil errors")
	// When errors is nil, the field is omitted due to omitempty tag
	assert.NotContains(t, string(bodyBytes), `"errors"`)
}

func TestSendError_EmptyErrors(t *testing.T) {
	tests := []struct {
		name     string
		errors   interface{}
		expected string
	}{
		{
			name:     "Empty map",
			errors:   map[string]string{},
			expected: `"errors":{}`,
		},
		{
			name:     "Empty slice",
			errors:   []string{},
			expected: `"errors":[]`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			app := fiber.New()
			app.Get("/test", func(c *fiber.Ctx) error {
				return SendError(c, fiber.StatusBadRequest, "Error with empty errors", tt.errors)
			})

			req := httptest.NewRequest("GET", "/test", nil)
			resp, err := app.Test(req)
			assert.NoError(t, err)

			assert.Equal(t, fiber.StatusBadRequest, resp.StatusCode)
			bodyBytes, _ := io.ReadAll(resp.Body)
			assert.Contains(t, string(bodyBytes), tt.expected)
		})
	}
}

func TestSendError_NestedErrors(t *testing.T) {
	app := fiber.New()

	app.Get("/nested-error", func(c *fiber.Ctx) error {
		errors := map[string]interface{}{
			"field1": map[string]string{
				"rule":  "required",
				"value": "missing",
			},
			"field2": []string{
				"Error 1",
				"Error 2",
			},
		}
		return SendError(c, fiber.StatusBadRequest, "Nested errors", errors)
	})

	req := httptest.NewRequest("GET", "/nested-error", nil)
	resp, err := app.Test(req)
	assert.NoError(t, err)

	assert.Equal(t, fiber.StatusBadRequest, resp.StatusCode)
	bodyBytes, _ := io.ReadAll(resp.Body)
	assert.Contains(t, string(bodyBytes), "error")
	assert.Contains(t, string(bodyBytes), "Nested errors")
	assert.Contains(t, string(bodyBytes), "required")
	assert.Contains(t, string(bodyBytes), "Error 1")
}

func TestSendError_UnicodeInErrors(t *testing.T) {
	app := fiber.New()

	app.Get("/unicode-error", func(c *fiber.Ctx) error {
		errors := map[string]string{
			"arabic":   "خطأ في الإدخال",
			"chinese":  "输入错误",
			"special":  "Special <>&\"' chars",
		}
		return SendError(c, fiber.StatusBadRequest, "Unicode errors", errors)
	})

	req := httptest.NewRequest("GET", "/unicode-error", nil)
	resp, err := app.Test(req)
	assert.NoError(t, err)

	assert.Equal(t, fiber.StatusBadRequest, resp.StatusCode)
	bodyBytes, _ := io.ReadAll(resp.Body)
	bodyStr := string(bodyBytes)

	assert.Contains(t, bodyStr, "خطأ في الإدخال")
	assert.Contains(t, bodyStr, "输入错误")
	assert.Contains(t, bodyStr, "\\u003c\\u003e\\u0026\\\"'") // Special chars should be escaped (lowercase hex)
}

func TestSendSuccess_LargeDataPayload(t *testing.T) {
	app := fiber.New()

	app.Get("/large-data", func(c *fiber.Ctx) error {
		// Create a large data structure
		items := make([]map[string]interface{}, 1000)
		for i := 0; i < 1000; i++ {
			items[i] = map[string]interface{}{
				"id":          i,
				"name":        "Item " + string(rune('a'+(i%26))),
				"description": "This is a long description for item " + string(rune(i)),
				"metadata": map[string]interface{}{
					"created_at": "2024-01-01T00:00:00Z",
					"updated_at": "2024-01-01T00:00:00Z",
					"tags":       []string{"tag1", "tag2", "tag3"},
				},
			}
		}
		data := map[string]interface{}{
			"total":  1000,
			"items":  items,
			"page":   1,
			"limit":  1000,
		}
		return SendSuccess(c, data, "Large data payload")
	})

	req := httptest.NewRequest("GET", "/large-data", nil)
	resp, err := app.Test(req)
	assert.NoError(t, err)

	assert.Equal(t, fiber.StatusOK, resp.StatusCode)
	bodyBytes, _ := io.ReadAll(resp.Body)
	bodyStr := string(bodyBytes)

	assert.Contains(t, bodyStr, "success")
	assert.Contains(t, bodyStr, "Large data payload")
	assert.Contains(t, bodyStr, `"total":1000`)
	assert.Contains(t, bodyStr, `"page":1`)
	assert.Contains(t, bodyStr, `"limit":1000`)
	// Verify we have multiple items
	assert.Contains(t, bodyStr, `"items":[`)
}

func TestSendSuccess_DeeplyNestedStructures(t *testing.T) {
	app := fiber.New()

	app.Get("/deep-nested", func(c *fiber.Ctx) error {
		// Create a deeply nested structure (5 levels deep)
		data := map[string]interface{}{
			"level1": map[string]interface{}{
				"level2": map[string]interface{}{
					"level3": map[string]interface{}{
						"level4": map[string]interface{}{
							"level5": map[string]string{
								"deep": "value",
							},
						},
					},
				},
			},
		}
		return SendSuccess(c, data, "Deeply nested structure")
	})

	req := httptest.NewRequest("GET", "/deep-nested", nil)
	resp, err := app.Test(req)
	assert.NoError(t, err)

	assert.Equal(t, fiber.StatusOK, resp.StatusCode)
	bodyBytes, _ := io.ReadAll(resp.Body)
	bodyStr := string(bodyBytes)

	assert.Contains(t, bodyStr, "success")
	assert.Contains(t, bodyStr, "Deeply nested structure")
	assert.Contains(t, bodyStr, "level1")
	assert.Contains(t, bodyStr, "level2")
	assert.Contains(t, bodyStr, "level3")
	assert.Contains(t, bodyStr, "level4")
	assert.Contains(t, bodyStr, "level5")
	assert.Contains(t, bodyStr, "deep")
	assert.Contains(t, bodyStr, "value")
}

func TestSendSuccess_MixedDataTypes(t *testing.T) {
	app := fiber.New()

	app.Get("/mixed-types", func(c *fiber.Ctx) error {
		data := map[string]interface{}{
			"string":    "text",
			"integer":   42,
			"float":     3.14159,
			"boolean":   true,
			"null":      nil,
			"array":     []interface{}{1, "two", 3.0, true, nil},
			"object":    map[string]interface{}{"key": "value"},
			"zero_int":  0,
			"zero_float": 0.0,
			"empty_str": "",
		}
		return SendSuccess(c, data, "Mixed data types")
	})

	req := httptest.NewRequest("GET", "/mixed-types", nil)
	resp, err := app.Test(req)
	assert.NoError(t, err)

	assert.Equal(t, fiber.StatusOK, resp.StatusCode)
	bodyBytes, _ := io.ReadAll(resp.Body)
	bodyStr := string(bodyBytes)

	assert.Contains(t, bodyStr, "success")
	assert.Contains(t, bodyStr, "Mixed data types")
	assert.Contains(t, bodyStr, `"string":"text"`)
	assert.Contains(t, bodyStr, `"integer":42`)
	assert.Contains(t, bodyStr, `"float":3.14159`)
	assert.Contains(t, bodyStr, `"boolean":true`)
	assert.Contains(t, bodyStr, `"null":null`)
	assert.Contains(t, bodyStr, `"array":[1,"two",3,true,null]`)
	assert.Contains(t, bodyStr, `"object":{"key":"value"}`)
	assert.Contains(t, bodyStr, `"zero_int":0`)
	assert.Contains(t, bodyStr, `"zero_float":0`)
	assert.Contains(t, bodyStr, `"empty_str":""`)
}

func TestSendError_MixedErrorTypes(t *testing.T) {
	app := fiber.New()

	app.Get("/mixed-errors", func(c *fiber.Ctx) error {
		errors := map[string]interface{}{
			"string_error": "This is a string error",
			"map_errors": map[string]string{
				"email":    "Invalid email",
				"password": "Too short",
			},
			"array_errors": []string{
				"Error 1",
				"Error 2",
				"Error 3",
			},
			"null_error": nil,
		}
		return SendError(c, fiber.StatusBadRequest, "Mixed error types", errors)
	})

	req := httptest.NewRequest("GET", "/mixed-errors", nil)
	resp, err := app.Test(req)
	assert.NoError(t, err)

	assert.Equal(t, fiber.StatusBadRequest, resp.StatusCode)
	bodyBytes, _ := io.ReadAll(resp.Body)
	bodyStr := string(bodyBytes)

	assert.Contains(t, bodyStr, "error")
	assert.Contains(t, bodyStr, "Mixed error types")
	assert.Contains(t, bodyStr, "This is a string error")
	assert.Contains(t, bodyStr, "Invalid email")
	assert.Contains(t, bodyStr, "Too short")
	assert.Contains(t, bodyStr, "Error 1")
	assert.Contains(t, bodyStr, "Error 2")
	assert.Contains(t, bodyStr, "Error 3")
	assert.Contains(t, bodyStr, `"null_error":null`)
}

func TestSendSuccess_ArrayOfStructs(t *testing.T) {
	app := fiber.New()

	app.Get("/array-structs", func(c *fiber.Ctx) error {
		type User struct {
			ID    int    `json:"id"`
			Name  string `json:"name"`
			Email string `json:"email"`
		}

		users := []User{
			{ID: 1, Name: "Alice", Email: "alice@example.com"},
			{ID: 2, Name: "Bob", Email: "bob@example.com"},
			{ID: 3, Name: "Charlie", Email: "charlie@example.com"},
		}

		return SendSuccess(c, users, "Array of structs")
	})

	req := httptest.NewRequest("GET", "/array-structs", nil)
	resp, err := app.Test(req)
	assert.NoError(t, err)

	assert.Equal(t, fiber.StatusOK, resp.StatusCode)
	bodyBytes, _ := io.ReadAll(resp.Body)
	bodyStr := string(bodyBytes)

	assert.Contains(t, bodyStr, "success")
	assert.Contains(t, bodyStr, "Array of structs")
	assert.Contains(t, bodyStr, `"id":1`)
	assert.Contains(t, bodyStr, `"name":"Alice"`)
	assert.Contains(t, bodyStr, `"email":"alice@example.com"`)
	assert.Contains(t, bodyStr, `"id":2`)
	assert.Contains(t, bodyStr, `"name":"Bob"`)
	assert.Contains(t, bodyStr, `"id":3`)
	assert.Contains(t, bodyStr, `"name":"Charlie"`)
}

func TestSendCreated_NilData(t *testing.T) {
	app := fiber.New()

	app.Post("/create-nil", func(c *fiber.Ctx) error {
		return SendCreated(c, nil, "Resource created with nil data")
	})

	req := httptest.NewRequest("POST", "/create-nil", nil)
	resp, err := app.Test(req)
	assert.NoError(t, err)

	assert.Equal(t, fiber.StatusCreated, resp.StatusCode)
	bodyBytes, _ := io.ReadAll(resp.Body)
	assert.Contains(t, string(bodyBytes), "success")
	assert.Contains(t, string(bodyBytes), "Resource created with nil data")
	// When data is nil, the field is omitted due to omitempty tag
	assert.NotContains(t, string(bodyBytes), `"data"`)
}

func TestSendValidationError_NilErrors(t *testing.T) {
	app := fiber.New()

	app.Post("/validate-nil", func(c *fiber.Ctx) error {
		return SendValidationError(c, nil)
	})

	req := httptest.NewRequest("POST", "/validate-nil", nil)
	resp, err := app.Test(req)
	assert.NoError(t, err)

	assert.Equal(t, fiber.StatusBadRequest, resp.StatusCode)
	bodyBytes, _ := io.ReadAll(resp.Body)
	assert.Contains(t, string(bodyBytes), "error")
	assert.Contains(t, string(bodyBytes), "Validation failed")
	// When errors is nil, the field is omitted due to omitempty tag
	assert.NotContains(t, string(bodyBytes), `"errors"`)
}

func TestAPIResponse_JSONEncodingConsistency(t *testing.T) {
	tests := []struct {
		name           string
		data           interface{}
		expectedFields []string
	}{
		{
			name: "Simple map",
			data: map[string]string{"key": "value"},
			expectedFields: []string{`"status":"success"`, `"key":"value"`},
		},
		{
			name: "Nil data",
			data: nil,
			expectedFields: []string{`"status":"success"`},
		},
		{
			name: "Empty map",
			data: map[string]interface{}{},
			expectedFields: []string{`"status":"success"`, `"data":{}`},
		},
		{
			name: "Array",
			data: []int{1, 2, 3},
			expectedFields: []string{`"status":"success"`, `"data":[1,2,3]`},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			app := fiber.New()
			app.Get("/test", func(c *fiber.Ctx) error {
				return SendSuccess(c, tt.data, "Test")
			})

			req := httptest.NewRequest("GET", "/test", nil)
			resp, err := app.Test(req)
			assert.NoError(t, err)

			bodyBytes, _ := io.ReadAll(resp.Body)
			bodyStr := string(bodyBytes)

			for _, field := range tt.expectedFields {
				assert.Contains(t, bodyStr, field)
			}
		})
	}
}
