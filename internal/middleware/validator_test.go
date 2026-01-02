package middleware

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
)

// Test data structures
type ProductInput struct {
	Name        string  `json:"name" validate:"required,min=3,max=100"`
	Price       float64 `json:"price" validate:"required,gt=0"`
	Description string  `json:"description" validate:"max=500"`
	Category    string  `json:"category" validate:"required"`
	Stock       int     `json:"stock" validate:"gte=0"`
}

func TestInputValidation_ValidInput(t *testing.T) {
	app := fiber.New()

	// Test route with validation
	app.Post("/products", func(c *fiber.Ctx) error {
		var input ProductInput
		if err := c.BodyParser(&input); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Invalid request body",
			})
		}

		// Validate input
		if input.Name == "" {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Product name is required",
			})
		}
		if input.Price <= 0 {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Price must be greater than 0",
			})
		}
		if input.Category == "" {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Category is required",
			})
		}

		return c.Status(fiber.StatusCreated).JSON(fiber.Map{
			"message": "Product created successfully",
		})
	})

	// Test valid input
	validInput := ProductInput{
		Name:        "Test Product",
		Price:       99.99,
		Description: "A test product",
		Category:    "tops",
		Stock:       10,
	}

	jsonInput, _ := json.Marshal(validInput)
	req := httptest.NewRequest("POST", "/products", bytes.NewReader(jsonInput))
	req.Header.Set("Content-Type", "application/json")
	resp := httptest.NewRecorder()
	app.Handler()(resp, req)

	assert.Equal(t, fiber.StatusCreated, resp.Code)
	assert.Contains(t, resp.Body.String(), "Product created successfully")
}

func TestInputValidation_MissingRequiredFields(t *testing.T) {
	app := fiber.New()

	app.Post("/products", func(c *fiber.Ctx) error {
		var input ProductInput
		if err := c.BodyParser(&input); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Invalid request body",
			})
		}

		if input.Name == "" {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Product name is required",
			})
		}
		if input.Price <= 0 {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Price must be greater than 0",
			})
		}
		if input.Category == "" {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Category is required",
			})
		}

		return c.SendStatus(fiber.StatusCreated)
	})

	// Test missing name
	input1 := map[string]interface{}{
		"price":    99.99,
		"category": "tops",
	}
	jsonInput1, _ := json.Marshal(input1)
	req1 := httptest.NewRequest("POST", "/products", bytes.NewReader(jsonInput1))
	req1.Header.Set("Content-Type", "application/json")
	resp1 := httptest.NewRecorder()
	app.Handler()(resp1, req1)
	assert.Equal(t, fiber.StatusBadRequest, resp1.Code)
	assert.Contains(t, resp1.Body.String(), "Product name is required")

	// Test missing price
	input2 := map[string]interface{}{
		"name":     "Test Product",
		"category": "tops",
	}
	jsonInput2, _ := json.Marshal(input2)
	req2 := httptest.NewRequest("POST", "/products", bytes.NewReader(jsonInput2))
	req2.Header.Set("Content-Type", "application/json")
	resp2 := httptest.NewRecorder()
	app.Handler()(resp2, req2)
	assert.Equal(t, fiber.StatusBadRequest, resp2.Code)
	assert.Contains(t, resp2.Body.String(), "Price must be greater than 0")

	// Test missing category
	input3 := map[string]interface{}{
		"name":  "Test Product",
		"price": 99.99,
	}
	jsonInput3, _ := json.Marshal(input3)
	req3 := httptest.NewRequest("POST", "/products", bytes.NewReader(jsonInput3))
	req3.Header.Set("Content-Type", "application/json")
	resp3 := httptest.NewRecorder()
	app.Handler()(resp3, req3)
	assert.Equal(t, fiber.StatusBadRequest, resp3.Code)
	assert.Contains(t, resp3.Body.String(), "Category is required")
}

func TestInputValidation_SQLInjection(t *testing.T) {
	app := fiber.New()

	app.Post("/products", func(c *fiber.Ctx) error {
		var input ProductInput
		if err := c.BodyParser(&input); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Invalid request body",
			})
		}

		// Sanitize input (basic example)
		input.Name = strings.ReplaceAll(input.Name, "'", "''")
		input.Category = strings.ReplaceAll(input.Category, "'", "''")

		return c.JSON(fiber.Map{
			"name":     input.Name,
			"category": input.Category,
		})
	})

	// Test SQL injection attempt
	sqlInjectionInput := map[string]interface{}{
		"name":     "'; DROP TABLE products; --",
		"price":    99.99,
		"category": "tops'; DELETE FROM products; --",
	}

	jsonInput, _ := json.Marshal(sqlInjectionInput)
	req := httptest.NewRequest("POST", "/products", bytes.NewReader(jsonInput))
	req.Header.Set("Content-Type", "application/json")
	resp := httptest.NewRecorder()
	app.Handler()(resp, req)

	assert.Equal(t, fiber.StatusOK, resp.Code)
	// Verify that the input was sanitized
	assert.NotContains(t, resp.Body.String(), "'; DROP TABLE")
	assert.NotContains(t, resp.Body.String(), "DELETE FROM")
}

func TestInputValidation_XSSAttack(t *testing.T) {
	app := fiber.New()

	app.Post("/products", func(c *fiber.Ctx) error {
		var input ProductInput
		if err := c.BodyParser(&input); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Invalid request body",
			})
		}

		// Sanitize input (basic XSS prevention)
		input.Name = strings.ReplaceAll(input.Name, "<script>", "")
		input.Name = strings.ReplaceAll(input.Name, "</script>", "")
		input.Description = strings.ReplaceAll(input.Description, "<script>", "")
		input.Description = strings.ReplaceAll(input.Description, "</script>", "")

		return c.JSON(fiber.Map{
			"name":        input.Name,
			"description": input.Description,
		})
	})

	// Test XSS attack attempt
	xssInput := map[string]interface{}{
		"name":        "<script>alert('XSS')</script>Product",
		"price":       99.99,
		"category":    "tops",
		"description": "<script>document.cookie</script>Description",
	}

	jsonInput, _ := json.Marshal(xssInput)
	req := httptest.NewRequest("POST", "/products", bytes.NewReader(jsonInput))
	req.Header.Set("Content-Type", "application/json")
	resp := httptest.NewRecorder()
	app.Handler()(resp, req)

	assert.Equal(t, fiber.StatusOK, resp.Code)
	// Verify that the script tags were removed
	assert.NotContains(t, resp.Body.String(), "<script>")
	assert.NotContains(t, resp.Body.String(), "</script>")
}

func TestInputValidation_CommandInjection(t *testing.T) {
	app := fiber.New()

	app.Post("/products", func(c *fiber.Ctx) error {
		var input ProductInput
		if err := c.BodyParser(&input); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Invalid request body",
			})
		}

		// Sanitize input (basic command injection prevention)
		input.Name = strings.ReplaceAll(input.Name, ";", "")
		input.Name = strings.ReplaceAll(input.Name, "&", "")
		input.Name = strings.ReplaceAll(input.Name, "|", "")
		input.Name = strings.ReplaceAll(input.Name, "`", "")

		return c.JSON(fiber.Map{
			"name": input.Name,
		})
	})

	// Test command injection attempt
	cmdInjectionInput := map[string]interface{}{
		"name":     "product; rm -rf /",
		"price":    99.99,
		"category": "tops",
	}

	jsonInput, _ := json.Marshal(cmdInjectionInput)
	req := httptest.NewRequest("POST", "/products", bytes.NewReader(jsonInput))
	req.Header.Set("Content-Type", "application/json")
	resp := httptest.NewRecorder()
	app.Handler()(resp, req)

	assert.Equal(t, fiber.StatusOK, resp.Code)
	// Verify that the dangerous characters were removed
	assert.NotContains(t, resp.Body.String(), ";")
	assert.NotContains(t, resp.Body.String(), "rm -rf")
}

func TestInputValidation_PathTraversal(t *testing.T) {
	app := fiber.New()

	app.Post("/upload", func(c *fiber.Ctx) error {
		var input struct {
			Filename string `json:"filename"`
		}
		if err := c.BodyParser(&input); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Invalid request body",
			})
		}

		// Prevent path traversal
		if strings.Contains(input.Filename, "..") {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Invalid filename",
			})
		}

		return c.JSON(fiber.Map{
			"filename": input.Filename,
		})
	})

	// Test path traversal attempt
	pathTraversalInput := map[string]interface{}{
		"filename": "../../../etc/passwd",
	}

	jsonInput, _ := json.Marshal(pathTraversalInput)
	req := httptest.NewRequest("POST", "/upload", bytes.NewReader(jsonInput))
	req.Header.Set("Content-Type", "application/json")
	resp := httptest.NewRecorder()
	app.Handler()(resp, req)

	assert.Equal(t, fiber.StatusBadRequest, resp.Code)
	assert.Contains(t, resp.Body.String(), "Invalid filename")
}

func TestInputValidation_EmailValidation(t *testing.T) {
	app := fiber.New()

	app.Post("/users", func(c *fiber.Ctx) error {
		var input struct {
			Email string `json:"email"`
		}
		if err := c.BodyParser(&input); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Invalid request body",
			})
		}

		// Basic email validation
		if !strings.Contains(input.Email, "@") || !strings.Contains(input.Email, ".") {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Invalid email address",
			})
		}

		return c.SendStatus(fiber.StatusCreated)
	})

	// Test valid email
	validEmail := map[string]interface{}{
		"email": "test@example.com",
	}
	jsonInput1, _ := json.Marshal(validEmail)
	req1 := httptest.NewRequest("POST", "/users", bytes.NewReader(jsonInput1))
	req1.Header.Set("Content-Type", "application/json")
	resp1 := httptest.NewRecorder()
	app.Handler()(resp1, req1)
	assert.Equal(t, fiber.StatusCreated, resp1.Code)

	// Test invalid email
	invalidEmail := map[string]interface{}{
		"email": "invalid-email",
	}
	jsonInput2, _ := json.Marshal(invalidEmail)
	req2 := httptest.NewRequest("POST", "/users", bytes.NewReader(jsonInput2))
	req2.Header.Set("Content-Type", "application/json")
	resp2 := httptest.NewRecorder()
	app.Handler()(resp2, req2)
	assert.Equal(t, fiber.StatusBadRequest, resp2.Code)
	assert.Contains(t, resp2.Body.String(), "Invalid email address")
}

func TestInputValidation_NumericRangeValidation(t *testing.T) {
	app := fiber.New()

	app.Post("/products", func(c *fiber.Ctx) error {
		var input ProductInput
		if err := c.BodyParser(&input); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Invalid request body",
			})
		}

		// Validate price range
		if input.Price < 0 || input.Price > 1000000 {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Price must be between 0 and 1,000,000",
			})
		}

		// Validate stock range
		if input.Stock < 0 || input.Stock > 10000 {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Stock must be between 0 and 10,000",
			})
		}

		return c.SendStatus(fiber.StatusCreated)
	})

	// Test price too high
	highPriceInput := map[string]interface{}{
		"name":     "Test Product",
		"price":    9999999,
		"category": "tops",
		"stock":    10,
	}
	jsonInput1, _ := json.Marshal(highPriceInput)
	req1 := httptest.NewRequest("POST", "/products", bytes.NewReader(jsonInput1))
	req1.Header.Set("Content-Type", "application/json")
	resp1 := httptest.NewRecorder()
	app.Handler()(resp1, req1)
	assert.Equal(t, fiber.StatusBadRequest, resp1.Code)

	// Test negative price
	negativePriceInput := map[string]interface{}{
		"name":     "Test Product",
		"price":    -10,
		"category": "tops",
		"stock":    10,
	}
	jsonInput2, _ := json.Marshal(negativePriceInput)
	req2 := httptest.NewRequest("POST", "/products", bytes.NewReader(jsonInput2))
	req2.Header.Set("Content-Type", "application/json")
	resp2 := httptest.NewRecorder()
	app.Handler()(resp2, req2)
	assert.Equal(t, fiber.StatusBadRequest, resp2.Code)

	// Test valid price
	validPriceInput := map[string]interface{}{
		"name":     "Test Product",
		"price":    99.99,
		"category": "tops",
		"stock":    10,
	}
	jsonInput3, _ := json.Marshal(validPriceInput)
	req3 := httptest.NewRequest("POST", "/products", bytes.NewReader(jsonInput3))
	req3.Header.Set("Content-Type", "application/json")
	resp3 := httptest.NewRecorder()
	app.Handler()(resp3, req3)
	assert.Equal(t, fiber.StatusCreated, resp3.Code)
}

func TestInputValidation_StringLengthValidation(t *testing.T) {
	app := fiber.New()

	app.Post("/products", func(c *fiber.Ctx) error {
		var input ProductInput
		if err := c.BodyParser(&input); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Invalid request body",
			})
		}

		// Validate name length
		if len(input.Name) < 3 || len(input.Name) > 100 {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Product name must be between 3 and 100 characters",
			})
		}

		// Validate description length
		if len(input.Description) > 500 {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Description must not exceed 500 characters",
			})
		}

		return c.SendStatus(fiber.StatusCreated)
	})

	// Test name too short
	shortNameInput := map[string]interface{}{
		"name":        "AB",
		"price":       99.99,
		"category":    "tops",
		"description": "A test product",
	}
	jsonInput1, _ := json.Marshal(shortNameInput)
	req1 := httptest.NewRequest("POST", "/products", bytes.NewReader(jsonInput1))
	req1.Header.Set("Content-Type", "application/json")
	resp1 := httptest.NewRecorder()
	app.Handler()(resp1, req1)
	assert.Equal(t, fiber.StatusBadRequest, resp1.Code)

	// Test description too long
	longDescription := strings.Repeat("A", 501)
	longDescInput := map[string]interface{}{
		"name":        "Test Product",
		"price":       99.99,
		"category":    "tops",
		"description": longDescription,
	}
	jsonInput2, _ := json.Marshal(longDescInput)
	req2 := httptest.NewRequest("POST", "/products", bytes.NewReader(jsonInput2))
	req2.Header.Set("Content-Type", "application/json")
	resp2 := httptest.NewRecorder()
	app.Handler()(resp2, req2)
	assert.Equal(t, fiber.StatusBadRequest, resp2.Code)

	// Test valid input
	validInput := map[string]interface{}{
		"name":        "Test Product",
		"price":       99.99,
		"category":    "tops",
		"description": "A test product",
	}
	jsonInput3, _ := json.Marshal(validInput)
	req3 := httptest.NewRequest("POST", "/products", bytes.NewReader(jsonInput3))
	req3.Header.Set("Content-Type", "application/json")
	resp3 := httptest.NewRecorder()
	app.Handler()(resp3, req3)
	assert.Equal(t, fiber.StatusCreated, resp3.Code)
}

func TestRequestBodyParsing_MalformedJSON(t *testing.T) {
	app := fiber.New()

	app.Post("/products", func(c *fiber.Ctx) error {
		var input ProductInput
		if err := c.BodyParser(&input); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Invalid JSON format",
			})
		}
		return c.SendStatus(fiber.StatusCreated)
	})

	// Test malformed JSON
	req := httptest.NewRequest("POST", "/products", bytes.NewReader([]byte("{invalid json}")))
	req.Header.Set("Content-Type", "application/json")
	resp := httptest.NewRecorder()
	app.Handler()(resp, req)

	assert.Equal(t, fiber.StatusBadRequest, resp.Code)
	assert.Contains(t, resp.Body.String(), "Invalid JSON format")
}

func TestRequestBodyParsing_EmptyRequestBody(t *testing.T) {
	app := fiber.New()

	app.Post("/products", func(c *fiber.Ctx) error {
		var input ProductInput
		if err := c.BodyParser(&input); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Request body is required",
			})
		}
		return c.SendStatus(fiber.StatusCreated)
	})

	// Test empty request body
	req := httptest.NewRequest("POST", "/products", bytes.NewReader([]byte{}))
	req.Header.Set("Content-Type", "application/json")
	resp := httptest.NewRecorder()
	app.Handler()(resp, req)

	assert.Equal(t, fiber.StatusBadRequest, resp.Code)
}

func TestRequestBodyParsing_InvalidContentType(t *testing.T) {
	app := fiber.New()

	app.Post("/products", func(c *fiber.Ctx) error {
		if c.Get("Content-Type") != "application/json" {
			return c.Status(fiber.StatusUnsupportedMediaType).JSON(fiber.Map{
				"error": "Content-Type must be application/json",
			})
		}
		return c.SendStatus(fiber.StatusCreated)
	})

	// Test invalid content type
	req := httptest.NewRequest("POST", "/products", bytes.NewReader([]byte("{}")))
	req.Header.Set("Content-Type", "text/plain")
	resp := httptest.NewRecorder()
	app.Handler()(resp, req)

	assert.Equal(t, fiber.StatusUnsupportedMediaType, resp.Code)
	assert.Contains(t, resp.Body.String(), "Content-Type must be application/json")
}