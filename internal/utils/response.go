package utils

import (
	"github.com/gofiber/fiber/v2"
)

// APIResponse Standard API response structure
type APIResponse struct {
	Status  string      `json:"status"`
	Message string      `json:"message,omitempty"`
	Data    interface{} `json:"data,omitempty"`
	Errors  interface{} `json:"errors,omitempty"`
}

// SendSuccess sends a success response with 200 OK (or custom status)
func SendSuccess(c *fiber.Ctx, data interface{}, message string, status ...int) error {
	code := fiber.StatusOK
	if len(status) > 0 {
		code = status[0]
	}

	return c.Status(code).JSON(APIResponse{
		Status:  "success",
		Message: message,
		Data:    data,
	})
}

// SendError sends an error response with appropriate status code
func SendError(c *fiber.Ctx, code int, message string, errors interface{}) error {
	return c.Status(code).JSON(APIResponse{
		Status:  "error",
		Message: message,
		Errors:  errors,
	})
}

// SendValidationError is a shortcut for validation errors (400 Bad Request)
func SendValidationError(c *fiber.Ctx, errors interface{}) error {
	return SendError(c, fiber.StatusBadRequest, "Validation failed", errors)
}

// SendCreated sends a created response (201 Created)
func SendCreated(c *fiber.Ctx, data interface{}, message string) error {
	return SendSuccess(c, data, message, fiber.StatusCreated)
}
