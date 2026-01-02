package errors

import (
	"errors"
	"fmt"
	"runtime/debug"
	"time"
)

// ErrorCode represents the type of error
type ErrorCode string

const (
	// Validation errors
	ErrCodeValidation    ErrorCode = "VALIDATION_ERROR"
	ErrCodeInvalidInput  ErrorCode = "INVALID_INPUT"
	ErrCodeMissingField  ErrorCode = "MISSING_FIELD"

	// Authentication errors
	ErrCodeUnauthorized    ErrorCode = "UNAUTHORIZED"
	ErrCodeForbidden       ErrorCode = "FORBIDDEN"
	ErrCodeInvalidToken    ErrorCode = "INVALID_TOKEN"
	ErrCodeExpiredToken    ErrorCode = "EXPIRED_TOKEN"

	// Resource errors
	ErrCodeNotFound      ErrorCode = "NOT_FOUND"
	ErrCodeAlreadyExists ErrorCode = "ALREADY_EXISTS"
	ErrCodeConflict      ErrorCode = "CONFLICT"

	// Business logic errors
	ErrCodeBusinessLogic ErrorCode = "BUSINESS_LOGIC_ERROR"
	ErrCodeOperationFailed ErrorCode = "OPERATION_FAILED"

	// System errors
	ErrCodeInternal      ErrorCode = "INTERNAL_ERROR"
	ErrCodeDatabase      ErrorCode = "DATABASE_ERROR"
	ErrCodeExternalAPI   ErrorCode = "EXTERNAL_API_ERROR"
	ErrCodeRateLimit     ErrorCode = "RATE_LIMIT_EXCEEDED"
	ErrCodeTimeout       ErrorCode = "TIMEOUT"
)

// AppError represents a structured application error
type AppError struct {
	Code       ErrorCode     `json:"code"`
	Message    string        `json:"message"`
	Details    interface{}   `json:"details,omitempty"`
	StatusCode int           `json:"-"`
	Stack      string        `json:"-"`
	Timestamp  time.Time     `json:"timestamp"`
	Err        error         `json:"-"`
}

// Error implements the error interface
func (e *AppError) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("[%s] %s: %v", e.Code, e.Message, e.Err)
	}
	return fmt.Sprintf("[%s] %s", e.Code, e.Message)
}

// Unwrap returns the underlying error
func (e *AppError) Unwrap() error {
	return e.Err
}

// NewAppError creates a new AppError
func NewAppError(code ErrorCode, message string) *AppError {
	return &AppError{
		Code:       code,
		Message:    message,
		StatusCode: getStatusCode(code),
		Timestamp:  time.Now(),
		Stack:      string(debug.Stack()),
	}
}

// NewAppErrorWithDetails creates a new AppError with details
func NewAppErrorWithDetails(code ErrorCode, message string, details interface{}) *AppError {
	return &AppError{
		Code:       code,
		Message:    message,
		Details:    details,
		StatusCode: getStatusCode(code),
		Timestamp:  time.Now(),
		Stack:      string(debug.Stack()),
	}
}

// WrapError wraps an existing error with AppError
func WrapError(code ErrorCode, message string, err error) *AppError {
	return &AppError{
		Code:       code,
		Message:    message,
		StatusCode: getStatusCode(code),
		Timestamp:  time.Now(),
		Stack:      string(debug.Stack()),
		Err:        err,
	}
}

// WrapErrorWithDetails wraps an existing error with AppError and details
func WrapErrorWithDetails(code ErrorCode, message string, err error, details interface{}) *AppError {
	return &AppError{
		Code:       code,
		Message:    message,
		Details:    details,
		StatusCode: getStatusCode(code),
		Timestamp:  time.Now(),
		Stack:      string(debug.Stack()),
		Err:        err,
	}
}

// getStatusCode returns the appropriate HTTP status code for an error code
func getStatusCode(code ErrorCode) int {
	switch code {
	case ErrCodeValidation, ErrCodeInvalidInput, ErrCodeMissingField:
		return 400
	case ErrCodeUnauthorized, ErrCodeInvalidToken, ErrCodeExpiredToken:
		return 401
	case ErrCodeForbidden:
		return 403
	case ErrCodeNotFound:
		return 404
	case ErrCodeAlreadyExists, ErrCodeConflict:
		return 409
	case ErrCodeRateLimit:
		return 429
	case ErrCodeInternal, ErrCodeDatabase, ErrCodeExternalAPI, ErrCodeTimeout, ErrCodeBusinessLogic, ErrCodeOperationFailed:
		return 500
	default:
		return 500
	}
}

// Predefined error constructors
func NewValidationError(message string) *AppError {
	return NewAppError(ErrCodeValidation, message)
}

func NewValidationErrorWithDetails(message string, details interface{}) *AppError {
	return NewAppErrorWithDetails(ErrCodeValidation, message, details)
}

func NewInvalidInputError(message string) *AppError {
	return NewAppError(ErrCodeInvalidInput, message)
}

func NewUnauthorizedError(message string) *AppError {
	return NewAppError(ErrCodeUnauthorized, message)
}

func NewForbiddenError(message string) *AppError {
	return NewAppError(ErrCodeForbidden, message)
}

func NewNotFoundError(resource string) *AppError {
	return NewAppError(ErrCodeNotFound, fmt.Sprintf("%s not found", resource))
}

func NewAlreadyExistsError(resource string) *AppError {
	return NewAppError(ErrCodeAlreadyExists, fmt.Sprintf("%s already exists", resource))
}

func NewConflictError(message string) *AppError {
	return NewAppError(ErrCodeConflict, message)
}

func NewBusinessLogicError(message string) *AppError {
	return NewAppError(ErrCodeBusinessLogic, message)
}

func NewInternalError(message string) *AppError {
	return NewAppError(ErrCodeInternal, message)
}

func NewDatabaseError(message string) *AppError {
	return NewAppError(ErrCodeDatabase, message)
}

func NewExternalAPIError(service, message string) *AppError {
	return NewAppError(ErrCodeExternalAPI, fmt.Sprintf("%s: %s", service, message))
}

func NewRateLimitError(message string) *AppError {
	return NewAppError(ErrCodeRateLimit, message)
}

func NewTimeoutError(message string) *AppError {
	return NewAppError(ErrCodeTimeout, message)
}

// IsAppError checks if an error is an AppError
func IsAppError(err error) bool {
	var appErr *AppError
	return errors.As(err, &appErr)
}

// GetAppError extracts AppError from an error
func GetAppError(err error) *AppError {
	var appErr *AppError
	if errors.As(err, &appErr) {
		return appErr
	}
	return nil
}

// ValidationErrorDetail represents a single validation error
type ValidationErrorDetail struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

// NewValidationErrorFromDetails creates a validation error from field errors
func NewValidationErrorFromDetails(details []ValidationErrorDetail) *AppError {
	return NewAppErrorWithDetails(
		ErrCodeValidation,
		"Validation failed",
		details,
	)
}
