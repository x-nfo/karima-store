package middleware

import (
	"fmt"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/karima-store/internal/telemetry"
)

// TracingMiddleware adds distributed tracing to requests
func TracingMiddleware() fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Generate or retrieve trace ID
		traceID := telemetry.GetOrCreateTraceID(c)

		// Generate span ID
		spanID := telemetry.GenerateSpanID()

		// Get parent span ID if exists
		parentSpanID := c.Get("X-Parent-Span-ID")

		// Create span
		span := &telemetry.Span{
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
		telemetry.StoreSpan(span)

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
		telemetry.UpdateSpan(span)

		// Log trace
		telemetry.LogTrace(span)

		return err
	}
}

// TraceHandler returns a handler that exposes trace information
func TraceHandler() fiber.Handler {
	return func(c *fiber.Ctx) error {
		traceID := c.Params("id")

		if traceID != "" {
			// Get specific trace
			traces, err := telemetry.GetTrace(telemetry.TraceID(traceID))
			if err != nil {
				return err
			}
			return c.JSON(traces)
		}

		// Get all traces
		traces := telemetry.GetAllTraces()
		return c.JSON(traces)
	}
}

// Re-export helper functions for backward compatibility/ease of use if needed,
// but better to use telemetry directly in other packages.
// Since we are refactoring, we will update consumers to use telemetry directly.
