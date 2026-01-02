package middleware

import (
	"crypto/rand"
	"fmt"
	"github.com/karima-store/internal/errors"
	"sync"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
)

// TraceID represents a unique trace identifier
type TraceID string

// Span represents a single operation in a trace
type Span struct {
	TraceID      TraceID     `json:"trace_id"`
	SpanID       string      `json:"span_id"`
	ParentSpanID string      `json:"parent_span_id,omitempty"`
	Operation    string      `json:"operation"`
	StartTime    time.Time   `json:"start_time"`
	EndTime      time.Time   `json:"end_time"`
	Duration     time.Duration `json:"duration"`
	Status       string      `json:"status"`
	Error        string      `json:"error,omitempty"`
	Tags         map[string]string `json:"tags,omitempty"`
	Metadata     map[string]interface{} `json:"metadata,omitempty"`
}

// TraceStore stores trace data
type TraceStore struct {
	mu     sync.RWMutex
	traces map[TraceID][]*Span
}

var (
	traceStore = &TraceStore{
		traces: make(map[TraceID][]*Span),
	}
	activeSpans = make(map[string]*Span)
	spanMu sync.RWMutex
)

// TracingMiddleware adds distributed tracing to requests
func TracingMiddleware() fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Generate or retrieve trace ID
		traceID := getOrCreateTraceID(c)

		// Generate span ID
		spanID := generateSpanID()

		// Get parent span ID if exists
		parentSpanID := c.Get("X-Parent-Span-ID")

		// Create span
		span := &Span{
			TraceID:      traceID,
			SpanID:       spanID,
			ParentSpanID: parentSpanID,
			Operation:    fmt.Sprintf("%s %s", c.Method(), c.Path()),
			StartTime:    time.Now(),
			Status:       "started",
			Tags: map[string]string{
				"http.method":     c.Method(),
				"http.url":        c.Path(),
				"http.host":       c.Hostname(),
				"http.user_agent": c.Get("User-Agent"),
				"http.ip":         c.IP(),
			},
			Metadata: map[string]interface{}{
				"query_params": c.Queries(),
			},
		}

		// Store span
		storeSpan(span)

		// Set trace headers
		c.Set("X-Trace-ID", string(traceID))
		c.Set("X-Span-ID", spanID)
		c.Locals("trace_id", traceID)
		c.Locals("span_id", spanID)

		// Process request
		err := c.Next()

		// Update span
		span.EndTime = time.Now()
		span.Duration = span.EndTime.Sub(span.StartTime)

		if err != nil {
			span.Status = "error"
			span.Error = err.Error()
			span.Tags["error"] = "true"
		} else {
			span.Status = "completed"
			span.Tags["http.status_code"] = fmt.Sprintf("%d", c.Response().StatusCode())
		}

		// Update span in store
		updateSpan(span)

		// Log trace
		logTrace(span)

		return err
	}
}

// getOrCreateTraceID retrieves existing trace ID or creates a new one
func getOrCreateTraceID(c *fiber.Ctx) TraceID {
	// Check if trace ID is in headers
	if traceID := c.Get("X-Trace-ID"); traceID != "" {
		return TraceID(traceID)
	}

	// Check if trace ID is in locals
	if traceID, ok := c.Locals("trace_id").(TraceID); ok {
		return traceID
	}

	// Generate new trace ID
	return generateTraceID()
}

// generateTraceID generates a new trace ID
func generateTraceID() TraceID {
	return TraceID(fmt.Sprintf("%s-%s-%s-%s-%s",
		randomHex(4),
		randomHex(4),
		randomHex(4),
		randomHex(4),
		randomHex(12),
	))
}

// generateSpanID generates a new span ID
func generateSpanID() string {
	return randomHex(16)
}

// randomHex generates a random hex string
func randomHex(length int) string {
	bytes := make([]byte, length/2)
	rand.Read(bytes)
	return fmt.Sprintf("%x", bytes)
}

// storeSpan stores a span in the trace store
func storeSpan(span *Span) {
	spanMu.Lock()
	defer spanMu.Unlock()

	activeSpans[span.SpanID] = span

	traceStore.mu.Lock()
	defer traceStore.mu.Unlock()

	traceStore.traces[span.TraceID] = append(traceStore.traces[span.TraceID], span)
}

// updateSpan updates an existing span
func updateSpan(span *Span) {
	spanMu.Lock()
	defer spanMu.Unlock()

	if existingSpan, exists := activeSpans[span.SpanID]; exists {
		existingSpan.EndTime = span.EndTime
		existingSpan.Duration = span.Duration
		existingSpan.Status = span.Status
		existingSpan.Error = span.Error
		existingSpan.Tags = span.Tags
		existingSpan.Metadata = span.Metadata
	}
}

// logTrace logs trace information
func logTrace(span *Span) {
	log.Infof("Trace: %s, Span: %s, Operation: %s, Duration: %v, Status: %s",
		span.TraceID,
		span.SpanID,
		span.Operation,
		span.Duration,
		span.Status,
	)
}

// StartSpan starts a new span for an operation
func StartSpan(traceID TraceID, parentSpanID, operation string) *Span {
	span := &Span{
		TraceID:      traceID,
		SpanID:       generateSpanID(),
		ParentSpanID: parentSpanID,
		Operation:    operation,
		StartTime:    time.Now(),
		Status:       "started",
		Tags:         make(map[string]string),
		Metadata:     make(map[string]interface{}),
	}

	storeSpan(span)
	return span
}

// FinishSpan finishes a span
func FinishSpan(span *Span, err error) {
	span.EndTime = time.Now()
	span.Duration = span.EndTime.Sub(span.StartTime)

	if err != nil {
		span.Status = "error"
		span.Error = err.Error()
		span.Tags["error"] = "true"
	} else {
		span.Status = "completed"
	}

	updateSpan(span)
	logTrace(span)
}

// GetTrace retrieves a trace by ID
func GetTrace(traceID TraceID) ([]*Span, error) {
	traceStore.mu.RLock()
	defer traceStore.mu.RUnlock()

	traces, exists := traceStore.traces[traceID]
	if !exists {
		return nil, errors.NewNotFoundError("Trace not found")
	}

	return traces, nil
}

// GetAllTraces retrieves all traces
func GetAllTraces() map[TraceID][]*Span {
	traceStore.mu.RLock()
	defer traceStore.mu.RUnlock()

	result := make(map[TraceID][]*Span)
	for traceID, spans := range traceStore.traces {
		result[traceID] = append([]*Span{}, spans...)
	}

	return result
}

// ClearTraces clears all traces
func ClearTraces() {
	traceStore.mu.Lock()
	defer traceStore.mu.Unlock()

	traceStore.traces = make(map[TraceID][]*Span)

	spanMu.Lock()
	defer spanMu.Unlock()

	activeSpans = make(map[string]*Span)
}

// TraceHandler returns a handler that exposes trace information
func TraceHandler() fiber.Handler {
	return func(c *fiber.Ctx) error {
		traceID := c.Params("id")

		if traceID != "" {
			// Get specific trace
			traces, err := GetTrace(TraceID(traceID))
			if err != nil {
				return err
			}
			return c.JSON(traces)
		}

		// Get all traces
		traces := GetAllTraces()
		return c.JSON(traces)
	}
}

// WithSpan executes a function with tracing
func WithSpan(traceID TraceID, parentSpanID, operation string, fn func() error) error {
	span := StartSpan(traceID, parentSpanID, operation)
	err := fn()
	FinishSpan(span, err)
	return err
}

// AddTag adds a tag to the current span
func AddTag(spanID string, key, value string) {
	spanMu.RLock()
	defer spanMu.RUnlock()

	if span, exists := activeSpans[spanID]; exists {
		span.Tags[key] = value
	}
}

// AddMetadata adds metadata to the current span
func AddMetadata(spanID string, key string, value interface{}) {
	spanMu.RLock()
	defer spanMu.RUnlock()

	if span, exists := activeSpans[spanID]; exists {
		if span.Metadata == nil {
			span.Metadata = make(map[string]interface{})
		}
		span.Metadata[key] = value
	}
}

// GetCurrentSpanID returns the current span ID from context
func GetCurrentSpanID(c *fiber.Ctx) string {
	if spanID, ok := c.Locals("span_id").(string); ok {
		return spanID
	}
	return ""
}

// GetCurrentTraceID returns the current trace ID from context
func GetCurrentTraceID(c *fiber.Ctx) TraceID {
	if traceID, ok := c.Locals("trace_id").(TraceID); ok {
		return traceID
	}
	return ""
}

// TraceOperation traces a database operation
func TraceOperation(traceID TraceID, parentSpanID, operation, query string, duration time.Duration, err error) {
	span := &Span{
		TraceID:      traceID,
		SpanID:       generateSpanID(),
		ParentSpanID: parentSpanID,
		Operation:    operation,
		StartTime:    time.Now().Add(-duration),
		EndTime:      time.Now(),
		Duration:     duration,
		Status:       "completed",
		Tags: map[string]string{
			"db.operation": operation,
			"db.type":      "sql",
		},
		Metadata: map[string]interface{}{
			"db.query": query,
		},
	}

	if err != nil {
		span.Status = "error"
		span.Error = err.Error()
		span.Tags["error"] = "true"
	}

	storeSpan(span)
	logTrace(span)
}

// TraceCacheOperation traces a cache operation
func TraceCacheOperation(traceID TraceID, parentSpanID, operation, key string, duration time.Duration, err error) {
	span := &Span{
		TraceID:      traceID,
		SpanID:       generateSpanID(),
		ParentSpanID: parentSpanID,
		Operation:    operation,
		StartTime:    time.Now().Add(-duration),
		EndTime:      time.Now(),
		Duration:     duration,
		Status:       "completed",
		Tags: map[string]string{
			"cache.operation": operation,
			"cache.type":      "redis",
		},
		Metadata: map[string]interface{}{
			"cache.key": key,
		},
	}

	if err != nil {
		span.Status = "error"
		span.Error = err.Error()
		span.Tags["error"] = "true"
	}

	storeSpan(span)
	logTrace(span)
}
