package middleware

import (
	"fmt"
	"karima_store/internal/errors"
	"karima_store/internal/utils"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
)

// ErrorHandler handles all errors in the application
func ErrorHandler() fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Process the request
		err := c.Next()

		// If no error, return
		if err == nil {
			return nil
		}

		// Log the error
		logError(c, err)

		// Handle the error based on its type
		return handleError(c, err)
	}
}

// handleError processes the error and returns appropriate response
func handleError(c *fiber.Ctx, err error) error {
	// Check if it's an AppError
	if appErr := errors.GetAppError(err); appErr != nil {
		// In production, don't expose stack traces
		if fiber.IsProduction() {
			appErr.Stack = ""
		}

		return utils.SendError(c, appErr.StatusCode, map[string]interface{}{
			"code":      string(appErr.Code),
			"message":   appErr.Message,
			"details":   appErr.Details,
			"timestamp": appErr.Timestamp,
		})
	}

	// Check if it's a Fiber error
	if fiberErr, ok := err.(*fiber.Error); ok {
		return utils.SendError(c, fiberErr.Code, map[string]interface{}{
			"code":    "FIBER_ERROR",
			"message": fiberErr.Message,
		})
	}

	// Handle unknown errors
	return utils.SendError(c, fiber.StatusInternalServerError, map[string]interface{}{
		"code":    "INTERNAL_ERROR",
		"message": "An unexpected error occurred",
	})
}

// logError logs error with context
func logError(c *fiber.Ctx, err error) {
	// Extract request information
	method := c.Method()
	path := c.Path()
	ip := c.IP()
	userAgent := c.Get("User-Agent")

	// Build error message
	errorMsg := fmt.Sprintf("Error occurred - Method: %s, Path: %s, IP: %s, UserAgent: %s, Error: %v",
		method, path, ip, userAgent, err)

	// Log based on error type
	if appErr := errors.GetAppError(err); appErr != nil {
		// Log with error code
		switch appErr.Code {
		case errors.ErrCodeValidation, errors.ErrCodeInvalidInput, errors.ErrCodeMissingField:
			log.Warnf("%s", errorMsg)
		case errors.ErrCodeUnauthorized, errors.ErrCodeForbidden:
			log.Warnf("%s", errorMsg)
		case errors.ErrCodeNotFound:
			log.Infof("%s", errorMsg)
		case errors.ErrCodeInternal, errors.ErrCodeDatabase, errors.ErrCodeExternalAPI:
			log.Errorf("%s\nStack: %s", errorMsg, appErr.Stack)
		default:
			log.Errorf("%s", errorMsg)
		}
	} else {
		// Log unknown errors
		log.Errorf("%s", errorMsg)
	}
}

// RecoverMiddleware recovers from panics
func RecoverMiddleware() fiber.Handler {
	return func(c *fiber.Ctx) error {
		defer func() {
			if r := recover(); r != nil {
				err := fmt.Errorf("panic recovered: %v", r)
				logError(c, err)
				
				// Return internal server error
				_ = utils.SendError(c, fiber.StatusInternalServerError, map[string]interface{}{
					"code":    "INTERNAL_ERROR",
					"message": "An unexpected error occurred",
				})
			}
		}()

		return c.Next()
	}
}

// ValidationError represents a validation error
type ValidationError struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

// ValidateAndHandle validates input and handles validation errors
func ValidateAndHandle(c *fiber.Ctx, validator interface{}) error {
	// This is a placeholder for validation logic
	// In a real implementation, you would use a validation library
	// like go-playground/validator
	
	// Example validation logic:
	/*
	if err := c.BodyParser(validator); err != nil {
		return errors.NewInvalidInputError("Invalid request body")
	}
	
	if err := validate.Struct(validator); err != nil {
		var validationErrors []ValidationError
		for _, err := range err.(validator.ValidationErrors) {
			validationErrors = append(validationErrors, ValidationError{
				Field:   strings.ToLower(err.Field()),
				Message: getValidationMessage(err),
			})
		}
		return errors.NewValidationErrorFromDetails(validationErrors)
	}
	*/
	
	return c.Next()
}

// getValidationMessage returns a user-friendly validation message
func getValidationMessage(field string, tag string) string {
	messages := map[string]string{
		"required": "This field is required",
		"email":    "Invalid email format",
		"min":      "Value is too short",
		"max":      "Value is too long",
		"gte":      "Value is too small",
		"lte":      "Value is too large",
	}
	
	if msg, ok := messages[tag]; ok {
		return msg
	}
	
	return fmt.Sprintf("Validation failed on %s", tag)
}

// SanitizeErrorMessage sanitizes error messages for client responses
func SanitizeErrorMessage(message string) string {
	// Remove sensitive information from error messages
	sensitivePatterns := []string{
		"password",
		"secret",
		"token",
		"key",
		"credential",
	}

	sanitized := message
	for _, pattern := range sensitivePatterns {
		if strings.Contains(strings.ToLower(sanitized), pattern) {
			sanitized = "Sensitive information redacted"
			break
		}
	}

	return sanitized
}
