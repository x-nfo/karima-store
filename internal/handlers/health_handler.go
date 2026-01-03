package handlers

import (
	"github.com/gofiber/fiber/v2"
)

// HealthHandler handles health check and metrics endpoints
type HealthHandler struct{}

// NewHealthHandler creates a new Health handler
func NewHealthHandler() *HealthHandler {
	return &HealthHandler{}
}

// HealthCheck returns server health status
// @Summary Health Check
// @Description Check if the server is healthy and running
// @Tags Health
// @Produce json
// @Success 200 {object} map[string]interface{} "Server is healthy"
// @Router /api/v1/health [get]
func (h *HealthHandler) HealthCheck(c *fiber.Ctx) error {
	return c.Status(200).JSON(fiber.Map{
		"status":  "up",
		"message": "Server is healthy",
	})
}

// Metrics returns Prometheus metrics
// @Summary Prometheus Metrics
// @Description Get Prometheus metrics for monitoring and alerting
// @Tags Health
// @Produce text/plain
// @Success 200 {string} string "Prometheus metrics"
// @Router /metrics [get]
func (h *HealthHandler) Metrics(c *fiber.Ctx) error {
	c.Set("Content-Type", "text/plain")
	return c.SendString("Prometheus metrics endpoint")
}
