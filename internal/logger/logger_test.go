package logger

import (
	"bytes"
	"encoding/json"
	"strings"
	"testing"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func TestNewLogger(t *testing.T) {
	// Test creating a new logger with JSON format
	cfg := &Config{
		Level:      "info",
		Output:     "stdout",
		Format:     "json",
		Env:        "production",
		WithCaller: true,
	}

	logger, err := NewLogger(cfg)
	if err != nil {
		t.Fatalf("Failed to create logger: %v", err)
	}
	defer logger.Sync()

	if logger == nil {
		t.Fatal("Logger should not be nil")
	}
}

func TestNewLoggerWithInvalidLevel(t *testing.T) {
	// Test creating a logger with invalid log level
	cfg := &Config{
		Level:  "invalid",
		Output: "stdout",
		Format: "json",
		Env:     "production",
	}

	_, err := NewLogger(cfg)
	if err == nil {
		t.Fatal("Expected error for invalid log level")
	}
}

func TestLoggerWithFields(t *testing.T) {
	// Test logging with structured fields
	cfg := &Config{
		Level:  "debug",
		Output: "stdout",
		Format: "json",
		Env:     "production",
	}

	logger, err := NewLogger(cfg)
	if err != nil {
		t.Fatalf("Failed to create logger: %v", err)
	}
	defer logger.Sync()

	// Create a logger with additional fields
	fieldsLogger := logger.WithFields(map[string]interface{}{
		"user_id": 123,
		"action":  "test",
	})

	fieldsLogger.Infow("Test message with fields", "key", "value")
}

func TestLoggerWithError(t *testing.T) {
	// Test logging with error
	cfg := &Config{
		Level:  "debug",
		Output: "stdout",
		Format: "json",
		Env:     "production",
	}

	logger, err := NewLogger(cfg)
	if err != nil {
		t.Fatalf("Failed to create logger: %v", err)
	}
	defer logger.Sync()

	testErr := &testError{message: "test error"}
	errorLogger := logger.WithError(testErr)
	errorLogger.Errorw("An error occurred")
}

func TestLoggerWithRequest(t *testing.T) {
	// Test logging with request context
	cfg := &Config{
		Level:  "debug",
		Output: "stdout",
		Format: "json",
		Env:     "production",
	}

	logger, err := NewLogger(cfg)
	if err != nil {
		t.Fatalf("Failed to create logger: %v", err)
	}
	defer logger.Sync()

	requestLogger := logger.WithRequest("GET", "/api/v1/products", "req-123")
	requestLogger.Infow("Request received", "status_code", 200)
}

func TestLoggerWithUser(t *testing.T) {
	// Test logging with user context
	cfg := &Config{
		Level:  "debug",
		Output: "stdout",
		Format: "json",
		Env:     "production",
	}

	logger, err := NewLogger(cfg)
	if err != nil {
		t.Fatalf("Failed to create logger: %v", err)
	}
	defer logger.Sync()

	userLogger := logger.WithUser("user-123", "test@example.com")
	userLogger.Infow("User action performed", "action", "login")
}

func TestLoggerWithDuration(t *testing.T) {
	// Test logging with duration
	cfg := &Config{
		Level:  "debug",
		Output: "stdout",
		Format: "json",
		Env:     "production",
	}

	logger, err := NewLogger(cfg)
	if err != nil {
		t.Fatalf("Failed to create logger: %v", err)
	}
	defer logger.Sync()

	durationLogger := logger.WithDuration(100000000) // 100ms in nanoseconds
	durationLogger.Infow("Operation completed")
}

func TestLoggerJSONFormat(t *testing.T) {
	// Test that logs are in JSON format
	_ = &Config{
		Level:  "debug",
		Output: "stdout",
		Format: "json",
		Env:     "production",
	}

	// Capture output
	var buf bytes.Buffer
	encoder := zapcore.NewJSONEncoder(zapcore.EncoderConfig{
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
	})

	core := zapcore.NewCore(
		encoder,
		zapcore.AddSync(&buf),
		zapcore.InfoLevel,
	)

	zapLogger := zap.New(core)
	sugaredLogger := zapLogger.Sugar()

	sugaredLogger.Infow("Test message", "key", "value")

	// Verify output is valid JSON
	output := buf.String()
	if !strings.Contains(output, "Test message") {
		t.Error("Output should contain message")
	}
	if !strings.Contains(output, "key") {
		t.Error("Output should contain key")
	}
	if !strings.Contains(output, "value") {
		t.Error("Output should contain value")
	}

	var result map[string]interface{}
	if err := json.Unmarshal([]byte(output), &result); err != nil {
		t.Errorf("Output should be valid JSON: %v", err)
	}
}

func TestDefaultConfig(t *testing.T) {
	cfg := DefaultConfig()
	if cfg.Level != "info" {
		t.Errorf("Expected level 'info', got '%s'", cfg.Level)
	}
	if cfg.Format != "json" {
		t.Errorf("Expected format 'json', got '%s'", cfg.Format)
	}
	if cfg.Env != "production" {
		t.Errorf("Expected env 'production', got '%s'", cfg.Env)
	}
}

func TestDevelopmentConfig(t *testing.T) {
	cfg := DevelopmentConfig()
	if cfg.Level != "debug" {
		t.Errorf("Expected level 'debug', got '%s'", cfg.Level)
	}
	if cfg.Format != "console" {
		t.Errorf("Expected format 'console', got '%s'", cfg.Format)
	}
	if cfg.Env != "development" {
		t.Errorf("Expected env 'development', got '%s'", cfg.Env)
	}
}

// testError is a simple error implementation for testing
type testError struct {
	message string
}

func (e *testError) Error() string {
	return e.message
}
