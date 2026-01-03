# Prometheus Metrics Implementation

## Overview

This document describes the migration from custom in-memory metrics to Prometheus metrics for better server health monitoring and observability.

## Changes Made

### 1. Added Prometheus Client Library

The Prometheus client library has been added to the project dependencies:

```bash
go get github.com/prometheus/client_golang/prometheus
```

### 2. New Prometheus Metrics Implementation

Created [`internal/telemetry/prometheus_metrics.go`](../internal/telemetry/prometheus_metrics.go) with the following metrics:

#### HTTP Request Metrics
- **`http_requests_total`**: Counter for total HTTP requests labeled by method, path, and status
- **`http_request_duration_seconds`**: Histogram tracking request duration in seconds
- **`http_requests_in_progress`**: Gauge for currently in-progress requests

#### Operation Metrics
- **`operations_total`**: Counter for operations performed labeled by operation type and status
- **`operation_duration_seconds`**: Histogram tracking operation duration

#### Go Runtime Metrics
Prometheus automatically collects Go runtime metrics including:
- Goroutines count
- Memory allocation
- GC statistics

### 3. Updated Metrics Middleware

The [`internal/middleware/metrics.go`](../internal/middleware/metrics.go) has been updated to use Prometheus:

- Records HTTP request metrics with method, path, and status code
- Tracks request duration using histograms
- Monitors in-progress requests
- Updates runtime metrics automatically

### 4. Updated Metrics Handler

The metrics endpoint now exposes metrics in Prometheus format:

```go
// GET /metrics endpoint returns Prometheus-formatted metrics
```

## Usage

### Accessing Metrics

Metrics are available at the `/metrics` endpoint:

```bash
curl http://localhost:8080/metrics
```

Example output:

```
# HELP http_requests_total Total number of HTTP requests
# TYPE http_requests_total counter
http_requests_total{method="GET",path="/api/products",status="200"} 42
http_requests_total{method="POST",path="/api/products",status="201"} 5

# HELP http_request_duration_seconds HTTP request duration in seconds
# TYPE http_request_duration_seconds histogram
http_request_duration_seconds_bucket{method="GET",path="/api/products",le="0.005"} 10
http_request_duration_seconds_bucket{method="GET",path="/api/products",le="0.01"} 25
http_request_duration_seconds_bucket{method="GET",path="/api/products",le="+Inf"} 42
http_request_duration_seconds_sum{method="GET",path="/api/products"} 0.523
http_request_duration_seconds_count{method="GET",path="/api/products"} 42

# HELP http_requests_in_progress Number of HTTP requests currently in progress
# TYPE http_requests_in_progress gauge
http_requests_in_progress{method="GET",path="/api/products"} 0

# HELP operations_total Total number of operations performed
# TYPE operations_total counter
operations_total{operation="database_query",status="success"} 150
operations_total{operation="cache_get",status="success"} 300

# HELP operation_duration_seconds Operation duration in seconds
# TYPE operation_duration_seconds histogram
operation_duration_seconds_bucket{operation="database_query",le="0.005"} 50
operation_duration_seconds_bucket{operation="database_query",le="0.01"} 100
operation_duration_seconds_bucket{operation="database_query",le="+Inf"} 150
operation_duration_seconds_sum{operation="database_query"} 1.234
operation_duration_seconds_count{operation="database_query"} 150
```

### Integrating with Prometheus

Add the following to your `prometheus.yml` configuration:

```yaml
scrape_configs:
  - job_name: 'karima_store'
    scrape_interval: 15s
    static_configs:
      - targets: ['localhost:8080']
```

### Example Grafana Queries

#### Request Rate
```
rate(http_requests_total[5m])
```

#### Error Rate
```
rate(http_requests_total{status=~"5.."}[5m])
```

#### Average Response Time
```
rate(http_request_duration_seconds_sum[5m]) / rate(http_request_duration_seconds_count[5m])
```

#### P95 Response Time
```
histogram_quantile(0.95, rate(http_request_duration_seconds_bucket[5m]))
```

#### Current In-Progress Requests
```
http_requests_in_progress
```

## API Reference

### Metrics Middleware

```go
import "github.com/karima-store/internal/middleware"

app.Use(middleware.MetricsMiddleware())
```

### Metrics Handler

```go
import "github.com/karima-store/internal/middleware"

app.Get("/metrics", middleware.MetricsHandler())
```

### Recording Custom Operations

```go
import "github.com/karima-store/internal/telemetry"

startTime := time.Now()
err := performOperation()
duration := time.Since(startTime)

telemetry.RecordOperationPrometheus("custom_operation", duration, err)
```

## Benefits of Prometheus Metrics

1. **Standard Format**: Prometheus uses industry-standard metrics format
2. **Better Visualization**: Works seamlessly with Grafana and other monitoring tools
3. **Advanced Queries**: Supports PromQL for complex metric queries
4. **Alerting**: Built-in alerting capabilities
5. **Long-term Storage**: Efficient time-series storage
6. **Multi-dimensional**: Labels allow for flexible metric aggregation
7. **Automatic Collection**: Go runtime metrics collected automatically

## Removed Code

The old custom metrics implementation in [`internal/telemetry/metrics.go`](../internal/telemetry/metrics.go) has been **removed** to avoid confusion and maintain a single, consistent metrics implementation using Prometheus.

The following deprecated functions are no longer available:
- `GetMetrics()`: Returns application metrics in JSON format (removed)
- `GetEndpointMetrics()`: Returns metrics for specific endpoint (removed)
- `GetAllEndpointMetrics()`: Returns all endpoint metrics (removed)
- `ResetMetrics()`: Resets all metrics (removed)
- `RecordMetrics()`: Records HTTP request metrics (removed)
- `RecordOperation()`: Records operation metrics (removed)
- `GetOperationMetrics()`: Returns metrics for specific operation (removed)

All metrics are now handled exclusively through Prometheus.

## Testing

Run the Prometheus metrics tests:

```bash
go test ./internal/telemetry/... -v
```

## Migration Guide

### For Developers

The old custom metrics implementation has been removed. All metrics are now handled through Prometheus:

**Current Usage (Automatic):**
```go
// Middleware handles HTTP request metrics automatically
app.Use(middleware.MetricsMiddleware())

// Expose metrics endpoint
app.Get("/metrics", middleware.MetricsHandler())
```

**For Custom Operations:**
```go
// Record custom operation metrics
telemetry.RecordOperationPrometheus("operation_name", duration, err)
```

### For Operations

1. Update Prometheus configuration to scrape the `/metrics` endpoint
2. Update Grafana dashboards to use new metric names
3. Update alerting rules to use PromQL queries
4. Remove any references to old custom metrics endpoints

## Troubleshooting

### Metrics Not Appearing

1. Verify the `/metrics` endpoint is accessible
2. Check Prometheus configuration for correct scrape interval
3. Ensure the metrics middleware is registered in the Fiber app

### High Cardinality Warnings

Avoid using high-cardinality labels (e.g., user IDs, timestamps) as they can cause performance issues:

```go
// BAD - high cardinality
httpRequestsTotal.WithLabelValues(method, path, userID).Inc()

// GOOD - low cardinality
httpRequestsTotal.WithLabelValues(method, path, status).Inc()
```

## Performance Considerations

- Prometheus metrics are designed for high-performance scenarios
- Minimal overhead on request processing
- Histograms use configurable buckets (default: Prometheus default buckets)
- Consider customizing buckets for your specific use case:

```go
customBuckets := []float64{0.001, 0.005, 0.01, 0.025, 0.05, 0.1, 0.25, 0.5, 1, 2.5, 5, 10}
```

## Future Enhancements

Potential improvements:

1. Add custom business metrics (e.g., orders per minute, revenue tracking)
2. Implement metric labels for additional context (e.g., service version, environment)
3. Add metric aggregation and summary endpoints
4. Implement metric export for external monitoring systems
5. Add custom histogram buckets for specific endpoints

## References

- [Prometheus Documentation](https://prometheus.io/docs/)
- [Prometheus Go Client](https://github.com/prometheus/client_golang)
- [PromQL Query Language](https://prometheus.io/docs/prometheus/latest/querying/basics/)
- [Grafana Dashboards](https://grafana.com/grafana/dashboards/)
