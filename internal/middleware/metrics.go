package middleware

import (
	"fmt"
	"time"

	"github.com/gofiber/adaptor/v2"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
	"github.com/karima-store/internal/telemetry"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

// MetricsMiddleware collects metrics for each request using Prometheus
func MetricsMiddleware() fiber.Handler {
	return func(c *fiber.Ctx) error {
		startTime := time.Now()
		method := c.Method()

		// Use c.Route().Path to get the route pattern (e.g., /product/:id) instead of
		// the actual path (e.g., /product/123) to avoid high cardinality issues
		path := c.Route().Path
		if path == "" {
			// Fallback to c.Path() if route path is not available
			path = c.Path()
		}

		// Increment in-progress counter
		telemetry.IncrementInProgress(method, path)

		// Process request
		err := c.Next()

		// Calculate duration
		duration := time.Since(startTime)

		// Decrement in-progress counter
		telemetry.DecrementInProgress(method, path)

		// Record Prometheus metrics
		statusCode := c.Response().StatusCode()
		telemetry.RecordHTTPRequest(method, path, statusCode, duration)

		// Update runtime metrics
		telemetry.UpdateRuntimeMetrics()

		// Check performance thresholds (optional, for alerting)
		pm := NewPerformanceMonitor()
		pm.CheckPerformance(duration, path)

		return err
	}
}

// MetricsHandler returns a handler that exposes Prometheus metrics using the adaptor pattern
func MetricsHandler() fiber.Handler {
	// Get the custom Prometheus registry
	registry := telemetry.GetRegistry()

	// Create the Prometheus HTTP handler with proper options
	promHandler := promhttp.HandlerFor(registry, promhttp.HandlerOpts{
		EnableOpenMetrics: false,
	})

	// Use adaptor to convert the standard net/http handler to a Fiber handler
	// This is more efficient and cleaner than using httptest hacks
	return adaptor.HTTPHandler(promHandler)
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
