# Structured JSON Logging Implementation

## Overview

This document describes the structured JSON logging implementation for the Karima Store project. The logging system uses [Uber Zap](https://github.com/uber-go/zap), a fast, structured, leveled logging library for Go.

## Features

- **Structured JSON Format**: All logs are output in JSON format for easy parsing by log aggregation tools
- **Machine-Readable**: Logs are structured with consistent field names for automated analysis
- **Contextual Logging**: Support for adding contextual information (user, request, duration, etc.)
- **Multiple Log Levels**: Debug, Info, Warn, Error, Fatal
- **Environment-Aware**: Different configurations for development and production
- **Caller Information**: Optional caller information for debugging
- **Stack Traces**: Automatic stack traces for error-level logs

## Implementation

### Core Components

#### 1. Logger Package (`internal/logger/logger.go`)

The main logger package provides:

- `Logger` struct wrapping Zap's SugaredLogger
- `Config` struct for logger configuration
- Helper methods for contextual logging:
  - `WithFields(fields)` - Add custom fields
  - `WithError(err)` - Add error context
  - `WithRequest(method, path, requestID)` - Add HTTP request context
  - `WithUser(userID, email)` - Add user context
  - `WithDuration(duration)` - Add timing information
  - `HTTPRequest(method, path, requestID, ip, statusCode, duration)` - Log HTTP requests
  - `DatabaseQuery(query, duration, rowsAffected)` - Log database queries

#### 2. Configuration Integration (`internal/config/config.go`)

Added `InitLogger()` method to Config struct:
```go
func (c *Config) InitLogger() error {
    cfg := &logger.Config{
        Level:      c.LogLevel,
        Output:     c.LogFile,
        Format:     "json",
        Env:        c.AppEnv,
        WithCaller: true,
    }

    // Use console format in development
    if c.AppEnv == "development" {
        cfg.Format = "console"
    }

    // Use stdout if log file is not specified
    if c.LogFile == "" {
        cfg.Output = "stdout"
    }

    return logger.Init(cfg)
}
```

## Log Format

### Production (JSON Format)

```json
{
  "level": "info",
  "timestamp": "2026-01-03T10:30:45.123+0700",
  "caller": "internal/services/notification_service.go:74",
  "message": "WhatsApp message sent successfully",
  "phone": "628123456789",
  "message_id": "12345"
}
```

### Development (Console Format)

```
10:30:45.123 INFO  WhatsApp message sent successfully phone=628123456789 message_id=12345
```

## Usage Examples

### Basic Logging

```go
import "github.com/karima-store/internal/logger"

// Initialize logger (typically in main.go)
cfg := config.Load()
if err := cfg.InitLogger(); err != nil {
    log.Fatalf("Failed to initialize logger: %v", err)
}

// Basic logging
logger.Log.Info("Server started")
logger.Log.Errorw("Database connection failed", "error", err)
logger.Log.Debugw("Processing request", "request_id", reqID)
```

### Contextual Logging

```go
// With custom fields
logger.Log.WithFields(map[string]interface{}{
    "user_id": 123,
    "action": "create_order",
}).Infow("Order created")

// With error
logger.Log.WithError(err).Errorw("Failed to process payment")

// With request context
logger.Log.WithRequest("POST", "/api/v1/orders", "req-123").Infow("Order request received")

// With user context
logger.Log.WithUser("user-123", "user@example.com").Infow("User logged in")

// With duration
logger.Log.WithDuration(time.Duration(100 * time.Millisecond)).Infow("Query executed")
```

### Specialized Logging Methods

```go
// HTTP Request Logging
logger.Log.HTTPRequest(
    "GET",
    "/api/v1/products",
    "req-123",
    "192.168.1.1",
    200,
    time.Duration(150 * time.Millisecond),
)

// Database Query Logging
logger.Log.DatabaseQuery(
    "SELECT * FROM products",
    time.Duration(50 * time.Millisecond),
    10,
)

// Error with Stack Trace
logger.Log.ErrorWithStack("Critical error occurred", err)

// Fatal with Stack Trace
logger.Log.FatalWithStack("Cannot start server", err)
```

## Configuration

### Environment Variables

| Variable | Default | Description |
|----------|---------|-------------|
| `LOG_LEVEL` | `info` | Log level: debug, info, warn, error, fatal |
| `LOG_FILE` | `logs/app.log` | Path to log file (empty for stdout) |
| `APP_ENV` | `development` | Environment: development, production |

### Log Levels

- `debug`: Detailed information for debugging
- `info`: General informational messages
- `warn`: Warning messages for potentially harmful situations
- `error`: Error messages for error events
- `fatal`: Severe error messages followed by application exit

## Integration Points

The structured logger has been integrated into:

1. **Configuration** (`internal/config/config.go`)
   - Environment file loading
   - Configuration validation
   - Logger initialization

2. **Database** (`internal/database/`)
   - PostgreSQL connection
   - Redis connection
   - Migration operations
   - Health checks

3. **Services** (`internal/services/`)
   - Notification service (WhatsApp messages)
   - Checkout service (order processing)
   - Product service (test setup)

4. **Handlers** (`internal/handlers/`)
   - Pricing handler (cache operations)

5. **Middleware** (`internal/middleware/`)
   - Rate limiter initialization

6. **Utils** (`internal/utils/`)
   - Server startup and shutdown
   - Graceful shutdown procedures

7. **Telemetry** (`internal/telemetry/`)
   - Distributed tracing
   - Span logging

8. **CLI Tools** (`cmd/migrate/`)
   - Database migration operations

## Backward Compatibility

All log calls maintain backward compatibility by checking if the structured logger is initialized:

```go
if logger.Log != nil {
    logger.Log.Infow("Message", "key", "value")
} else {
    log.Printf("Message: %v", value)
}
```

This ensures that the application continues to work even if the logger is not initialized (e.g., during early startup).

## Testing

Comprehensive tests have been added in `internal/logger/logger_test.go`:

- Logger creation and configuration
- Invalid log level handling
- Field-based logging
- Error logging
- Request context logging
- User context logging
- Duration logging
- JSON format validation
- Default and development config validation

Run tests:
```bash
go test ./internal/logger/... -v
```

## Best Practices

1. **Use Structured Fields**: Always use `Infow()`, `Errorw()`, etc. with key-value pairs
2. **Add Context**: Use `WithFields()`, `WithError()`, `WithRequest()`, `WithUser()` for contextual information
3. **Appropriate Levels**: Choose the right log level for the situation
4. **Consistent Field Names**: Use consistent naming across the codebase
5. **Avoid Sensitive Data**: Don't log passwords, tokens, or PII
6. **Performance**: Use appropriate log levels to avoid excessive logging in production

## Log Aggregation

The JSON format is compatible with popular log aggregation tools:

- **ELK Stack** (Elasticsearch, Logstash, Kibana)
- **Grafana Loki**
- **Splunk**
- **CloudWatch Logs**
- **Datadog**

Example Logstash configuration:
```yaml
input {
  file {
    path => "/var/log/karima_store/*.log"
    codec => json
  }
}
```

## Migration from Standard Log

To migrate existing code from standard `log` package:

### Before:
```go
log.Printf("User %s performed action %s", userID, action)
log.Printf("Error occurred: %v", err)
```

### After:
```go
logger.Log.WithFields(map[string]interface{}{
    "user_id": userID,
    "action": action,
}).Infow("User performed action")

logger.Log.WithError(err).Errorw("Error occurred")
```

## Troubleshooting

### Logger Not Initialized

If logs are not appearing in JSON format:

1. Check if `cfg.InitLogger()` is called in `main.go`
2. Verify `APP_ENV` environment variable
3. Check `LOG_LEVEL` setting

### Missing Fields

If expected fields are missing from logs:

1. Verify using the `*w` methods (e.g., `Infow`, `Errorw`)
2. Check that fields are passed as key-value pairs
3. Ensure logger is not nil before use

## Future Enhancements

Potential improvements:

1. **Log Rotation**: Implement automatic log file rotation
2. **Sampling**: Add log sampling for high-volume scenarios
3. **Filtering**: Implement dynamic log filtering based on conditions
4. **Metrics Integration**: Add Prometheus metrics integration
5. **Correlation IDs**: Enhance distributed tracing with better correlation
6. **Sensitive Data Masking**: Automatic masking of sensitive fields

## References

- [Uber Zap Documentation](https://github.com/uber-go/zap)
- [Structured Logging Best Practices](https://www.youtube.com/watch?v=Nln7q4j0jU)
- [Go Logging Patterns](https://go.dev/blog/log)
