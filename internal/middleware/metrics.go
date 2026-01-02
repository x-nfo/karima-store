package middleware

import (
	"fmt"
	"karima_store/internal/errors"
	"runtime"
	"sync"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
)

// Metrics holds application metrics
type Metrics struct {
	RequestCount      int64         `json:"request_count"`
	ResponseTime      time.Duration `json:"response_time"`
	ErrorRate         float64       `json:"error_rate"`
	ActiveGoroutines  int           `json:"active_goroutines"`
	MemoryUsage       uint64        `json:"memory_usage"`
	TotalErrors       int64         `json:"total_errors"`
	SuccessCount      int64         `json:"success_count"`
	AverageResponseTime time.Duration `json:"average_response_time"`
}

// MetricsStore stores metrics data
type MetricsStore struct {
	mu              sync.RWMutex
	requests        map[string]*RequestMetrics
	totalRequests   int64
	totalErrors     int64
	totalSuccess    int64
	totalResponseTime time.Duration
}

// RequestMetrics holds metrics for a specific endpoint
type RequestMetrics struct {
	Path        string
	Method      string
	Count       int64
	ErrorCount  int64
	SuccessCount int64
	TotalTime   time.Duration
	MinTime     time.Duration
	MaxTime     time.Duration
}

var (
	metricsStore = &MetricsStore{
		requests: make(map[string]*RequestMetrics),
	}
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
		recordMetrics(c, duration, err)

		return err
	}
}

// recordMetrics records metrics for a request
func recordMetrics(c *fiber.Ctx, duration time.Duration, err error) {
	metricsStore.mu.Lock()
	defer metricsStore.mu.Unlock()

	// Generate key for the endpoint
	key := c.Method() + ":" + c.Path()

	// Get or create request metrics
	metrics, exists := metricsStore.requests[key]
	if !exists {
		metrics = &RequestMetrics{
			Path:      c.Path(),
			Method:    c.Method(),
			MinTime:   duration,
			MaxTime:   duration,
		}
		metricsStore.requests[key] = metrics
	}

	// Update metrics
	metrics.Count++
	metrics.TotalTime += duration

	// Update min/max times
	if duration < metrics.MinTime {
		metrics.MinTime = duration
	}
	if duration > metrics.MaxTime {
		metrics.MaxTime = duration
	}

	// Track success/error
	statusCode := c.Response().StatusCode()
	if statusCode >= 400 {
		metrics.ErrorCount++
		metricsStore.totalErrors++
	} else {
		metrics.SuccessCount++
		metricsStore.totalSuccess++
	}

	metricsStore.totalRequests++
	metricsStore.totalResponseTime += duration
}

// GetMetrics returns current application metrics
func GetMetrics() Metrics {
	metricsStore.mu.RLock()
	defer metricsStore.mu.RUnlock()

	var m runtime.MemStats
	runtime.ReadMemStats(&m)

	var avgResponseTime time.Duration
	if metricsStore.totalRequests > 0 {
		avgResponseTime = metricsStore.totalResponseTime / time.Duration(metricsStore.totalRequests)
	}

	var errorRate float64
	if metricsStore.totalRequests > 0 {
		errorRate = float64(metricsStore.totalErrors) / float64(metricsStore.totalRequests) * 100
	}

	return Metrics{
		RequestCount:       metricsStore.totalRequests,
		ResponseTime:       avgResponseTime,
		ErrorRate:          errorRate,
		ActiveGoroutines:   runtime.NumGoroutine(),
		MemoryUsage:        m.Alloc,
		TotalErrors:        metricsStore.totalErrors,
		SuccessCount:       metricsStore.totalSuccess,
		AverageResponseTime: avgResponseTime,
	}
}

// GetEndpointMetrics returns metrics for a specific endpoint
func GetEndpointMetrics(path, method string) (*RequestMetrics, error) {
	metricsStore.mu.RLock()
	defer metricsStore.mu.RUnlock()

	key := method + ":" + path
	metrics, exists := metricsStore.requests[key]
	if !exists {
		return nil, errors.NewNotFoundError("Endpoint metrics not found")
	}

	// Return a copy to avoid race conditions
	copy := *metrics
	return &copy, nil
}

// GetAllEndpointMetrics returns metrics for all endpoints
func GetAllEndpointMetrics() map[string]*RequestMetrics {
	metricsStore.mu.RLock()
	defer metricsStore.mu.RUnlock()

	result := make(map[string]*RequestMetrics)
	for key, metrics := range metricsStore.requests {
		copy := *metrics
		result[key] = &copy
	}

	return result
}

// ResetMetrics resets all metrics
func ResetMetrics() {
	metricsStore.mu.Lock()
	defer metricsStore.mu.Unlock()

	metricsStore.requests = make(map[string]*RequestMetrics)
	metricsStore.totalRequests = 0
	metricsStore.totalErrors = 0
	metricsStore.totalSuccess = 0
	metricsStore.totalResponseTime = 0
}

// MetricsHandler returns a handler that exposes metrics
func MetricsHandler() fiber.Handler {
	return func(c *fiber.Ctx) error {
		metrics := GetMetrics()
		endpointMetrics := GetAllEndpointMetrics()

		return c.JSON(fiber.Map{
			"application": metrics,
			"endpoints":   endpointMetrics,
		})
	}
}

// PerformanceMonitor monitors application performance
type PerformanceMonitor struct {
	thresholds map[string]time.Duration
	alerts     []string
}

// NewPerformanceMonitor creates a new performance monitor
func NewPerformanceMonitor() *PerformanceMonitor {
	return &PerformanceMonitor{
		thresholds: map[string]time.Duration{
			"slow_request": 1 * time.Second,
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

// HealthMetrics holds health check metrics
type HealthMetrics struct {
	Status    string            `json:"status"`
	Timestamp time.Time         `json:"timestamp"`
	Metrics   map[string]Metric `json:"metrics"`
}

// Metric represents a single health metric
type Metric struct {
	Status  string `json:"status"`
	Message string `json:"message,omitempty"`
	Latency string `json:"latency,omitempty"`
}

// RecordOperation records an operation metric
func RecordOperation(operation string, duration time.Duration, err error) {
	metricsStore.mu.Lock()
	defer metricsStore.mu.Unlock()

	key := "operation:" + operation
	metrics, exists := metricsStore.requests[key]
	if !exists {
		metrics = &RequestMetrics{
			Path:    operation,
			Method:  "OPERATION",
			MinTime: duration,
			MaxTime: duration,
		}
		metricsStore.requests[key] = metrics
	}

	metrics.Count++
	metrics.TotalTime += duration

	if duration < metrics.MinTime {
		metrics.MinTime = duration
	}
	if duration > metrics.MaxTime {
		metrics.MaxTime = duration
	}

	if err != nil {
		metrics.ErrorCount++
		metricsStore.totalErrors++
	} else {
		metrics.SuccessCount++
		metricsStore.totalSuccess++
	}

	metricsStore.totalRequests++
	metricsStore.totalResponseTime += duration
}

// GetOperationMetrics returns metrics for a specific operation
func GetOperationMetrics(operation string) (*RequestMetrics, error) {
	metricsStore.mu.RLock()
	defer metricsStore.mu.RUnlock()

	key := "operation:" + operation
	metrics, exists := metricsStore.requests[key]
	if !exists {
		return nil, errors.NewNotFoundError("Operation metrics not found")
	}

	copy := *metrics
	return &copy, nil
}
