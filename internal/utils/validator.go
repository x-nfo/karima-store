package utils

import (
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

// Validator instance
var validate = validator.New()

// ValidationErrorResponse represents the structure of validation errors
type ValidationErrorResponse struct {
	Field string `json:"field"`
	Tag   string `json:"tag"`
	Value string `json:"value,omitempty"`
}

// ValidateStruct validates a struct and returns format-friendly errors
// Returns nil if no errors
func ValidateStruct(payload interface{}) []*ValidationErrorResponse {
	var errors []*ValidationErrorResponse
	err := validate.Struct(payload)
	if err != nil {
		for _, err := range err.(validator.ValidationErrors) {
			var element ValidationErrorResponse
			element.Field = err.Field()
			element.Tag = err.Tag()
			element.Value = err.Param()
			errors = append(errors, &element)
		}
	}
	return errors
}

// ParseAndValidate binds the request body to the struct and validates it
// Returns error (Fiber response) if validation fails, nil otherwise
func ParseAndValidate(c *fiber.Ctx, payload interface{}) error {
	// Parse body
	if err := c.BodyParser(payload); err != nil {
		return SendError(c, fiber.StatusBadRequest, "Invalid request body", err.Error())
	}

	// Validate
	if errors := ValidateStruct(payload); len(errors) > 0 {
		return SendValidationError(c, errors)
	}

	return nil
}
