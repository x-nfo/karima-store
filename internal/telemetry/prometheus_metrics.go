package telemetry

import (
	"strconv"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	// HTTP request metrics
	httpRequestsTotal = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "http_requests_total",
			Help: "Total number of HTTP requests",
		},
		[]string{"method", "path", "status"},
	)

	httpRequestDuration = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "http_request_duration_seconds",
			Help:    "HTTP request duration in seconds",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"method", "path"},
	)

	httpRequestsInProgress = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "http_requests_in_progress",
			Help: "Number of HTTP requests currently in progress",
		},
		[]string{"method", "path"},
	)

	// Operation metrics
	operationsTotal = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "operations_total",
			Help: "Total number of operations performed",
		},
		[]string{"operation", "status"},
	)

	operationDuration = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "operation_duration_seconds",
			Help:    "Operation duration in seconds",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"operation"},
	)
)

// RecordHTTPRequest records HTTP request metrics
func RecordHTTPRequest(method, path string, statusCode int, duration time.Duration) {
	status := strconv.Itoa(statusCode)

	httpRequestsTotal.WithLabelValues(method, path, status).Inc()
	httpRequestDuration.WithLabelValues(method, path).Observe(duration.Seconds())
}

// IncrementInProgress increments in-progress request counter
func IncrementInProgress(method, path string) {
	httpRequestsInProgress.WithLabelValues(method, path).Inc()
}

// DecrementInProgress decrements in-progress request counter
func DecrementInProgress(method, path string) {
	httpRequestsInProgress.WithLabelValues(method, path).Dec()
}

// UpdateRuntimeMetrics updates Go runtime metrics
// Note: Go runtime metrics are automatically collected by Prometheus
func UpdateRuntimeMetrics() {
	// Prometheus automatically collects Go runtime metrics
	// This function is kept for backward compatibility
}

// GetRegistry returns default Prometheus registry
func GetRegistry() *prometheus.Registry {
	return prometheus.DefaultRegisterer.(*prometheus.Registry)
}
