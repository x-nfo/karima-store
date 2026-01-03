package telemetry

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestRecordHTTPRequest(t *testing.T) {
	// Test recording HTTP request metrics
	RecordHTTPRequest("GET", "/api/products", 200, 100*time.Millisecond)
	RecordHTTPRequest("POST", "/api/products", 201, 150*time.Millisecond)
	RecordHTTPRequest("GET", "/api/products/1", 404, 50*time.Millisecond)

	// If no panic, the test passes
	assert.True(t, true, "RecordHTTPRequest executed without error")
}

func TestIncrementDecrementInProgress(t *testing.T) {
	// Test in-progress counter
	IncrementInProgress("GET", "/api/test")
	DecrementInProgress("GET", "/api/test")

	// If no panic, the test passes
	assert.True(t, true, "Increment/Decrement in-progress executed without error")
}

func TestRecordOperationPrometheus(t *testing.T) {
	// Test recording operation metrics
	RecordOperationPrometheus("database_query", 25*time.Millisecond, nil)
	RecordOperationPrometheus("cache_get", 5*time.Millisecond, nil)
	RecordOperationPrometheus("external_api_call", 500*time.Millisecond, assert.AnError)

	// If no panic, the test passes
	assert.True(t, true, "RecordOperationPrometheus executed without error")
}

func TestGetRegistry(t *testing.T) {
	// Test getting the registry
	registry := GetRegistry()
	assert.NotNil(t, registry, "Registry should not be nil")
}

func TestMetricsIntegration(t *testing.T) {
	// Test a complete metrics workflow
	method := "GET"
	path := "/api/health"

	// Simulate a request
	IncrementInProgress(method, path)
	time.Sleep(10 * time.Millisecond)
	RecordHTTPRequest(method, path, 200, 10*time.Millisecond)
	DecrementInProgress(method, path)

	// Update runtime metrics
	UpdateRuntimeMetrics()

	// Record an operation
	RecordOperationPrometheus("health_check", 5*time.Millisecond, nil)

	// If no panic, the test passes
	assert.True(t, true, "Metrics integration test passed")
}
