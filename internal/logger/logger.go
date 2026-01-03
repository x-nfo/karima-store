package logger

import (
	"fmt"
	"io"
	"os"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var (
	// Global logger instance
	Log *Logger
)

// Logger wraps zap logger with additional functionality
type Logger struct {
	*zap.SugaredLogger
}

// Config holds logger configuration
type Config struct {
	Level      string // debug, info, warn, error, fatal
	Output     string // stdout, stderr, or file path
	Format     string // json or console
	Env        string // development or production
	WithCaller bool   // include caller information
}

// NewLogger creates a new structured logger instance
func NewLogger(cfg *Config) (*Logger, error) {
	if cfg == nil {
		cfg = &Config{
			Level:      "info",
			Output:     "stdout",
			Format:     "json",
			Env:        "production",
			WithCaller: true,
		}
	}

	// Parse log level
	level, err := zapcore.ParseLevel(cfg.Level)
	if err != nil {
		return nil, fmt.Errorf("invalid log level: %w", err)
	}

	// Configure encoder
	encoderConfig := zapcore.EncoderConfig{
		TimeKey:        "timestamp",
		LevelKey:       "level",
		NameKey:        "logger",
		CallerKey:      "caller",
		FunctionKey:    zapcore.OmitKey,
		MessageKey:     "message",
		StacktraceKey:  "stacktrace",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.LowercaseLevelEncoder,
		EncodeTime:     zapcore.ISO8601TimeEncoder,
		EncodeDuration: zapcore.SecondsDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
	}

	// Use different encoding for development vs production
	if cfg.Env == "development" && cfg.Format == "console" {
		encoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
	}

	// Create encoder
	var encoder zapcore.Encoder
	if cfg.Format == "console" {
		encoder = zapcore.NewConsoleEncoder(encoderConfig)
	} else {
		encoder = zapcore.NewJSONEncoder(encoderConfig)
	}

	// Configure output
	var writer io.Writer
	switch cfg.Output {
	case "stderr":
		writer = os.Stderr
	case "stdout":
		writer = os.Stdout
	default:
		// File output
		file, err := os.OpenFile(cfg.Output, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			return nil, fmt.Errorf("failed to open log file: %w", err)
		}
		writer = file
	}

	// Create core
	core := zapcore.NewCore(
		encoder,
		zapcore.AddSync(writer),
		level,
	)

	// Create logger with options
	opts := []zap.Option{
		zap.AddCallerSkip(1),
		zap.AddStacktrace(zapcore.ErrorLevel),
	}

	if cfg.WithCaller {
		opts = append(opts, zap.AddCaller())
	}

	zapLogger := zap.New(core, opts...)
	sugaredLogger := zapLogger.Sugar()

	return &Logger{
		SugaredLogger: sugaredLogger,
	}, nil
}

// Init initializes the global logger
func Init(cfg *Config) error {
	logger, err := NewLogger(cfg)
	if err != nil {
		return err
	}
	Log = logger
	return nil
}

// WithFields creates a new logger with additional fields
func (l *Logger) WithFields(fields map[string]interface{}) *Logger {
	args := make([]interface{}, 0, len(fields)*2)
	for k, v := range fields {
		args = append(args, k, v)
	}
	return &Logger{
		SugaredLogger: l.SugaredLogger.With(args...),
	}
}

// WithError creates a new logger with error field
func (l *Logger) WithError(err error) *Logger {
	return &Logger{
		SugaredLogger: l.SugaredLogger.With("error", err.Error()),
	}
}

// WithRequest creates a new logger with request context
func (l *Logger) WithRequest(method, path, requestID string) *Logger {
	return &Logger{
		SugaredLogger: l.SugaredLogger.With(
			"method", method,
			"path", path,
			"request_id", requestID,
		),
	}
}

// WithUser creates a new logger with user context
func (l *Logger) WithUser(userID, email string) *Logger {
	return &Logger{
		SugaredLogger: l.SugaredLogger.With(
			"user_id", userID,
			"email", email,
		),
	}
}

// WithDuration creates a new logger with duration field
func (l *Logger) WithDuration(duration time.Duration) *Logger {
	return &Logger{
		SugaredLogger: l.SugaredLogger.With("duration_ms", duration.Milliseconds()),
	}
}

// HTTPRequest logs HTTP request information
func (l *Logger) HTTPRequest(method, path, requestID, ip string, statusCode int, duration time.Duration) {
	l.WithRequest(method, path, requestID).WithDuration(duration).Infow("HTTP request",
		"status_code", statusCode,
		"ip", ip,
	)
}

// DatabaseQuery logs database query information
func (l *Logger) DatabaseQuery(query string, duration time.Duration, rowsAffected int64) {
	l.WithDuration(duration).Debugw("Database query",
		"query", query,
		"rows_affected", rowsAffected,
	)
}

// Error logs error with stack trace
func (l *Logger) ErrorWithStack(message string, err error) {
	if err != nil {
		l.WithError(err).Errorw(message)
	} else {
		l.Error(message)
	}
}

// FatalWithStack logs fatal error with stack trace and exits
func (l *Logger) FatalWithStack(message string, err error) {
	if err != nil {
		l.WithError(err).Fatalw(message)
	} else {
		l.Fatal(message)
	}
}

// Sync flushes any buffered log entries
func (l *Logger) Sync() error {
	return l.SugaredLogger.Sync()
}

// DefaultConfig returns default logger configuration
func DefaultConfig() *Config {
	return &Config{
		Level:      "info",
		Output:     "stdout",
		Format:     "json",
		Env:        "production",
		WithCaller: true,
	}
}

// DevelopmentConfig returns development logger configuration
func DevelopmentConfig() *Config {
	return &Config{
		Level:      "debug",
		Output:     "stdout",
		Format:     "console",
		Env:        "development",
		WithCaller: true,
	}
}
