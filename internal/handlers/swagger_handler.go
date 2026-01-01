package handlers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/swagger"
	_ "github.com/karima-store/docs" // Import for side-effects
)

// SwaggerHandler handles Swagger documentation
type SwaggerHandler struct{}

// NewSwaggerHandler creates a new Swagger handler
func NewSwaggerHandler() *SwaggerHandler {
	return &SwaggerHandler{}
}

// ServeSwagger serves the Swagger UI
// @Summary Show Swagger UI
// @Description Get Swagger UI for API documentation
// @Tags Documentation
// @Success 200 {string} string "Swagger UI"
// @Router /swagger/index.html [get]
func (h *SwaggerHandler) ServeSwagger(c *fiber.Ctx) error {
	return swagger.HandlerDefault(c)
}
