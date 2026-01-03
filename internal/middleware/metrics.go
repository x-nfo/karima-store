package middleware

import (
	"fmt"
	"net/http/httptest"
	"time"

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
		path := c.Path()

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

// MetricsHandler returns a handler that exposes Prometheus metrics
func MetricsHandler() fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Set content type for Prometheus format
		c.Set("Content-Type", "text/plain; version=0.0.4; charset=utf-8")

		// Get Prometheus registry
		registry := telemetry.GetRegistry()

		// Create a test request to use with the Prometheus handler
		req := httptest.NewRequest("GET", "/metrics", nil)

		// Create a response recorder to capture the output
		recorder := httptest.NewRecorder()

		// Create the Prometheus HTTP handler
		handler := promhttp.HandlerFor(registry, promhttp.HandlerOpts{})

		// Serve the metrics to the recorder
		handler.ServeHTTP(recorder, req)

		// Write the response to Fiber context
		return c.Send(recorder.Body.Bytes())
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
