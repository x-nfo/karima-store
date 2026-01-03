package telemetry

import (
	"crypto/rand"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/karima-store/internal/errors"
)

// TraceID represents a unique trace identifier
type TraceID string

// Span represents a single operation in a trace
type Span struct {
	TraceID      TraceID                `json:"trace_id"`
	SpanID       string                 `json:"span_id"`
	ParentSpanID string                 `json:"parent_span_id,omitempty"`
	Operation    string                 `json:"operation"`
	StartTime    time.Time              `json:"start_time"`
	EndTime      time.Time              `json:"end_time"`
	Duration     time.Duration          `json:"duration"`
	Status       string                 `json:"status"`
	Error        string                 `json:"error,omitempty"`
	Tags         map[string]string      `json:"tags,omitempty"`
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
	spanMu      sync.RWMutex
)

// GetOrCreateTraceID retrieves existing trace ID or creates a new one
func GetOrCreateTraceID(c *fiber.Ctx) TraceID {
	if c == nil {
		return GenerateTraceID()
	}

	// Check if trace ID is in headers
	if traceID := c.Get("X-Trace-ID"); traceID != "" {
		return TraceID(traceID)
	}

	// Check if trace ID is in locals
	if traceID, ok := c.Locals("trace_id").(TraceID); ok {
		return traceID
	}

	// Generate new trace ID
	return GenerateTraceID()
}

// generateTraceID generates a new trace ID
func GenerateTraceID() TraceID {
	return TraceID(fmt.Sprintf("%s-%s-%s-%s-%s",
		randomHex(4),
		randomHex(4),
		randomHex(4),
		randomHex(4),
		randomHex(12),
	))
}

// generateSpanID generates a new span ID
func GenerateSpanID() string {
	return randomHex(16)
}

// randomHex generates a random hex string
func randomHex(length int) string {
	bytes := make([]byte, length/2)
	rand.Read(bytes)
	return fmt.Sprintf("%x", bytes)
}

// storeSpan stores a span in the trace store
func StoreSpan(span *Span) {
	spanMu.Lock()
	defer spanMu.Unlock()

	activeSpans[span.SpanID] = span

	traceStore.mu.Lock()
	defer traceStore.mu.Unlock()

	traceStore.traces[span.TraceID] = append(traceStore.traces[span.TraceID], span)
}

// updateSpan updates an existing span
func UpdateSpan(span *Span) {
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
func LogTrace(span *Span) {
	// Using standard log for now, can be replaced with structured logger
	log.Printf("Trace: %s, Span: %s, Operation: %s, Duration: %v, Status: %s",
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
		SpanID:       GenerateSpanID(),
		ParentSpanID: parentSpanID,
		Operation:    operation,
		StartTime:    time.Now(),
		Status:       "started",
		Tags:         make(map[string]string),
		Metadata:     make(map[string]interface{}),
	}

	StoreSpan(span)
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

	UpdateSpan(span)
	LogTrace(span)
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

// GetCurrentSpanID returns the current span ID from context
func GetCurrentSpanID(c *fiber.Ctx) string {
	if spanID, ok := c.Locals("span_id").(string); ok {
		return spanID
	}
	return ""
}

// GetCurrentTraceID returns the current trace ID from context
func GetCurrentTraceID(c *fiber.Ctx) TraceID {
	if c == nil {
		return ""
	}
	if traceID, ok := c.Locals("trace_id").(TraceID); ok {
		return traceID
	}
	return ""
}

// TraceOperation traces a database operation
func TraceOperation(traceID TraceID, parentSpanID, operation, query string, duration time.Duration, err error) {
	span := &Span{
		TraceID:      traceID,
		SpanID:       GenerateSpanID(),
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

	StoreSpan(span)
	LogTrace(span)
}

// TraceCacheOperation traces a cache operation
func TraceCacheOperation(traceID TraceID, parentSpanID, operation, key string, duration time.Duration, err error) {
	span := &Span{
		TraceID:      traceID,
		SpanID:       GenerateSpanID(),
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

	StoreSpan(span)
	LogTrace(span)
}
