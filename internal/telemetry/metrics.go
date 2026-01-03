package telemetry

import (
	"runtime"
	"sync"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/karima-store/internal/errors"
)

// Metrics holds application metrics
type Metrics struct {
	RequestCount        int64         `json:"request_count"`
	ResponseTime        time.Duration `json:"response_time"`
	ErrorRate           float64       `json:"error_rate"`
	ActiveGoroutines    int           `json:"active_goroutines"`
	MemoryUsage         uint64        `json:"memory_usage"`
	TotalErrors         int64         `json:"total_errors"`
	SuccessCount        int64         `json:"success_count"`
	AverageResponseTime time.Duration `json:"average_response_time"`
}

// MetricsStore stores metrics data
type MetricsStore struct {
	mu                sync.RWMutex
	requests          map[string]*RequestMetrics
	totalRequests     int64
	totalErrors       int64
	totalSuccess      int64
	totalResponseTime time.Duration
}

// RequestMetrics holds metrics for a specific endpoint
type RequestMetrics struct {
	Path         string
	Method       string
	Count        int64
	ErrorCount   int64
	SuccessCount int64
	TotalTime    time.Duration
	MinTime      time.Duration
	MaxTime      time.Duration
}

var (
	metricsStore = &MetricsStore{
		requests: make(map[string]*RequestMetrics),
	}
)

// RecordMetrics records metrics for a request
func RecordMetrics(c *fiber.Ctx, duration time.Duration, err error) {
	metricsStore.mu.Lock()
	defer metricsStore.mu.Unlock()

	// Generate key for the endpoint
	key := c.Method() + ":" + c.Path()

	// Get or create request metrics
	metrics, exists := metricsStore.requests[key]
	if !exists {
		metrics = &RequestMetrics{
			Path:    c.Path(),
			Method:  c.Method(),
			MinTime: duration,
			MaxTime: duration,
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
		RequestCount:        metricsStore.totalRequests,
		ResponseTime:        avgResponseTime,
		ErrorRate:           errorRate,
		ActiveGoroutines:    runtime.NumGoroutine(),
		MemoryUsage:         m.Alloc,
		TotalErrors:         metricsStore.totalErrors,
		SuccessCount:        metricsStore.totalSuccess,
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
