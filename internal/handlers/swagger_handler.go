package handlers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/karima-store/docs/swagger"
	"github.com/swaggo/fiber-swagger"
	"io/fs"
	"net/http"
)

// SwaggerHandler handles Swagger documentation
type SwaggerHandler struct{}

// NewSwaggerHandler creates a new Swagger handler
func NewSwaggerHandler() *SwaggerHandler {
	return &SwaggerHandler{}
}

// GetSwagger godoc
// @Summary Show Swagger UI
// @Description Get Swagger UI for API documentation
// @Tags Documentation
// @Accept  json
// @Produce  json
// @Success 200 {string} string "Swagger UI"
// @Router /swagger/* [get]
func (h *SwaggerHandler) GetSwagger(c *fiber.Ctx) error {
	// Serve Swagger UI files
	return c.SendFile("docs/swagger/index.html")
}

// GetSwaggerJSON godoc
// @Summary Get Swagger JSON
// @Description Get the OpenAPI specification in JSON format
// @Tags Documentation
// @Accept  json
// @Produce  json
// @Success 200 {object} swagger.KarimaStoreAPI
// @Router /swagger.json [get]
func (h *SwaggerHandler) GetSwaggerJSON(c *fiber.Ctx) error {
	return c.JSON(swagger.GetSwagger())
}

// ServeSwaggerAssets godoc
// @Summary Serve Swagger assets
// @Description Serve Swagger UI assets
// @Tags Documentation
// @Router /swagger/{path:*} [get]
func (h *SwaggerHandler) ServeSwaggerAssets(c *fiber.Ctx) error {
	// Serve Swagger UI assets
	return c.SendFile(c.Params("path"), true)
}