package telemetry

import (
	"strconv"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

// Custom registry to avoid polluting the global registry
var registry = prometheus.NewRegistry()

// HTTP request metrics
var (
	HttpRequestsTotal = promauto.With(registry).NewCounterVec(
		prometheus.CounterOpts{
			Name: "http_requests_total",
			Help: "Total number of HTTP requests",
		},
		[]string{"method", "path", "status"},
	)

	HttpRequestDuration = promauto.With(registry).NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "http_request_duration_seconds",
			Help:    "HTTP request duration in seconds",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"method", "path"},
	)

	HttpRequestsInProgress = promauto.With(registry).NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "http_requests_in_progress",
			Help: "Number of HTTP requests currently in progress",
		},
		[]string{"method", "path"},
	)
)

// Operation metrics
var (
	OperationsTotal = promauto.With(registry).NewCounterVec(
		prometheus.CounterOpts{
			Name: "operations_total",
			Help: "Total number of operations performed",
		},
		[]string{"operation", "status"},
	)

	OperationDuration = promauto.With(registry).NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "operation_duration_seconds",
			Help:    "Operation duration in seconds",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"operation"},
	)
)

// GetRegistry returns the custom Prometheus registry
func GetRegistry() *prometheus.Registry {
	return registry
}

// IncrementInProgress increments the in-progress request counter
func IncrementInProgress(method, path string) {
	HttpRequestsInProgress.WithLabelValues(method, path).Inc()
}

// DecrementInProgress decrements the in-progress request counter
func DecrementInProgress(method, path string) {
	HttpRequestsInProgress.WithLabelValues(method, path).Dec()
}

// RecordHTTPRequest records HTTP request metrics
func RecordHTTPRequest(method, path string, statusCode int, duration time.Duration) {
	status := strconv.Itoa(statusCode)

	HttpRequestsTotal.WithLabelValues(method, path, status).Inc()
	HttpRequestDuration.WithLabelValues(method, path).Observe(duration.Seconds())
}

// RecordOperationPrometheus records an operation metric using Prometheus
func RecordOperationPrometheus(operation string, duration time.Duration, err error) {
	status := "success"
	if err != nil {
		status = "error"
	}
	OperationsTotal.WithLabelValues(operation, status).Inc()
	OperationDuration.WithLabelValues(operation).Observe(duration.Seconds())
}

// UpdateRuntimeMetrics updates Go runtime metrics
// Note: Go runtime metrics are automatically collected by Prometheus
func UpdateRuntimeMetrics() {
	// Prometheus automatically collects Go runtime metrics
	// This function is kept for backward compatibility
}
