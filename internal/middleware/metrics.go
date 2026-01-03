package middleware

import (
	"fmt"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
	"github.com/karima-store/internal/telemetry"
)

// MetricsMiddleware collects metrics for each request
func MetricsMiddleware() fiber.Handler {
	return func(c *fiber.Ctx) error {
		startTime := time.Now()

		// Process request
		err := c.Next()

		// Calculate duration
		duration := time.Since(startTime)

		// Record metrics
		telemetry.RecordMetrics(c, duration, err)

		return err
	}
}

// MetricsHandler returns a handler that exposes metrics
func MetricsHandler() fiber.Handler {
	return func(c *fiber.Ctx) error {
		metrics := telemetry.GetMetrics()
		endpointMetrics := telemetry.GetAllEndpointMetrics()

		return c.JSON(fiber.Map{
			"application": metrics,
			"endpoints":   endpointMetrics,
		})
	}
}

// PerformanceMonitor monitors application performance
// Kept here as it was not moved to telemetry yet and seems specific
type PerformanceMonitor struct {
	thresholds map[string]time.Duration
	alerts     []string
}

// NewPerformanceMonitor creates a new performance monitor
func NewPerformanceMonitor() *PerformanceMonitor {
	return &PerformanceMonitor{
		thresholds: map[string]time.Duration{
			"slow_request":      1 * time.Second,
			"very_slow_request": 5 * time.Second,
		},
		alerts: make([]string, 0),
	}
}

// CheckPerformance checks if performance thresholds are exceeded
func (pm *PerformanceMonitor) CheckPerformance(duration time.Duration, path string) {
	if duration > pm.thresholds["very_slow_request"] {
		alert := fmt.Sprintf("Very slow request detected: %s took %v", path, duration)
		pm.alerts = append(pm.alerts, alert)
		log.Warn(alert)
	} else if duration > pm.thresholds["slow_request"] {
		log.Warnf("Slow request detected: %s took %v", path, duration)
	}
}

// GetAlerts returns all performance alerts
func (pm *PerformanceMonitor) GetAlerts() []string {
	return pm.alerts
}

// ClearAlerts clears all performance alerts
func (pm *PerformanceMonitor) ClearAlerts() {
	pm.alerts = make([]string, 0)
}
