package middleware

import (
	"karima_store/internal/errors"
	"mime/multipart"
	"net/http"
	"regexp"
	"strings"

	"github.com/gofiber/fiber/v2"
)

// ValidationConfig holds configuration for validation middleware
type ValidationConfig struct {
	MaxBodySize       int64    // Maximum request body size in bytes
	AllowedMethods    []string // Allowed HTTP methods
	AllowedMimeTypes  []string // Allowed content types
	RequireAuth       bool     // Require authentication
	EnableXSS         bool     // Enable XSS protection
	EnableSQLInjection bool    // Enable SQL injection protection
}

// DefaultValidationConfig returns default validation configuration
func DefaultValidationConfig() ValidationConfig {
	return ValidationConfig{
		MaxBodySize:       10 * 1024 * 1024, // 10MB
		AllowedMethods:    []string{"GET", "POST", "PUT", "PATCH", "DELETE"},
		AllowedMimeTypes:  []string{"application/json", "multipart/form-data", "application/x-www-form-urlencoded"},
		RequireAuth:       false,
		EnableXSS:         true,
		EnableSQLInjection: true,
	}
}

// ValidationMiddleware creates a comprehensive validation middleware
func ValidationMiddleware(config ValidationConfig) fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Validate HTTP method
		if !isAllowedMethod(c.Method(), config.AllowedMethods) {
			return errors.NewInvalidInputError("Method not allowed")
		}

		// Validate content type for POST, PUT, PATCH requests
		if c.Method() == "POST" || c.Method() == "PUT" || c.Method() == "PATCH" {
			contentType := c.Get("Content-Type")
			if contentType != "" && !isAllowedContentType(contentType, config.AllowedMimeTypes) {
				return errors.NewInvalidInputError("Unsupported content type")
			}
		}

		// Validate body size
		if c.Method() == "POST" || c.Method() == "PUT" || c.Method() == "PATCH" {
			contentLength := c.Get("Content-Length")
			if contentLength != "" {
				length := c.Context().Request.Header.ContentLength()
				if length > config.MaxBodySize {
					return errors.NewInvalidInputError("Request body too large")
				}
			}
		}

		// XSS protection for query parameters and form data
		if config.EnableXSS {
			if err := sanitizeInput(c); err != nil {
				return err
			}
		}

		// SQL injection protection
		if config.EnableSQLInjection {
			if err := checkSQLInjection(c); err != nil {
				return err
			}
		}

		return c.Next()
	}
}

// isAllowedMethod checks if the HTTP method is allowed
func isAllowedMethod(method string, allowedMethods []string) bool {
	for _, allowed := range allowedMethods {
		if method == allowed {
			return true
		}
	}
	return false
}

// isAllowedContentType checks if the content type is allowed
func isAllowedContentType(contentType string, allowedTypes []string) bool {
	// Extract the main content type (ignore charset, boundary, etc.)
	mainType := strings.Split(contentType, ";")[0]
	mainType = strings.TrimSpace(mainType)

	for _, allowed := range allowedTypes {
		if mainType == allowed {
			return true
		}
	}
	return false
}

// sanitizeInput sanitizes input to prevent XSS attacks
func sanitizeInput(c *fiber.Ctx) error {
	// Sanitize query parameters
	for key, values := range c.Queries() {
		for i, value := range values {
			sanitized := sanitizeXSS(value)
			if sanitized != value {
				c.Query(key, sanitized)
			}
		}
	}

	// Sanitize route parameters
	for key, value := range c.Route().Params {
		sanitized := sanitizeXSS(value)
		if sanitized != value {
			c.Params(key, sanitized)
		}
	}

	// Sanitize headers (specific headers only)
	sensitiveHeaders := []string{"User-Agent", "Referer", "Origin"}
	for _, header := range sensitiveHeaders {
		value := c.Get(header)
		if value != "" {
			sanitized := sanitizeXSS(value)
			if sanitized != value {
				c.Set(header, sanitized)
			}
		}
	}

	return nil
}

// sanitizeXSS removes potential XSS patterns from input
func sanitizeXSS(input string) string {
	// Remove common XSS patterns
	xssPatterns := []string{
		"<script.*?>.*?</script>",
		"javascript:",
		"on\\w+\\s*=",
		"eval\\s*\\(",
		"expression\\s*\\(",
		"vbscript:",
		"fromCharCode",
		"&#x",
		"&#",
	}

	result := input
	for _, pattern := range xssPatterns {
		re := regexp.MustCompile("(?i)" + pattern)
		result = re.ReplaceAllString(result, "")
	}

	return result
}

// checkSQLInjection checks for potential SQL injection patterns
func checkSQLInjection(c *fiber.Ctx) error {
	// SQL injection patterns to detect
	sqlPatterns := []string{
		"'\\s*or\\s*'.*'",
		"'\\s*and\\s*'.*'",
		"'\\s*;\\s*",
		"'\\s*--",
		"'\\s*#",
		"\\bunion\\s+select\\b",
		"\\bdrop\\s+table\\b",
		"\\bdelete\\s+from\\b",
		"\\binsert\\s+into\\b",
		"\\bupdate\\s+\\w+\\s+set\\b",
		"\\bexec\\s*\\(",
		"\\bexecute\\s*\\(",
		"\\bsp_executesql\\b",
		"\\bxp_cmdshell\\b",
	}

	// Check query parameters
	for _, values := range c.Queries() {
		for _, value := range values {
			if containsSQLPattern(value, sqlPatterns) {
				return errors.NewInvalidInputError("Invalid input detected")
			}
		}
	}

	// Check route parameters
	for _, value := range c.Route().Params {
		if containsSQLPattern(value, sqlPatterns) {
			return errors.NewInvalidInputError("Invalid input detected")
		}
	}

	return nil
}

// containsSQLPattern checks if input contains SQL injection patterns
func containsSQLPattern(input string, patterns []string) bool {
	lowerInput := strings.ToLower(input)
	for _, pattern := range patterns {
		re := regexp.MustCompile("(?i)" + pattern)
		if re.MatchString(lowerInput) {
			return true
		}
	}
	return false
}

// ValidateJSON validates JSON request body
func ValidateJSON(c *fiber.Ctx, target interface{}) error {
	if err := c.BodyParser(target); err != nil {
		return errors.NewInvalidInputError("Invalid JSON format")
	}
	return nil
}

// ValidateRequiredFields checks if required fields are present
func ValidateRequiredFields(data map[string]interface{}, requiredFields []string) error {
	var missingFields []string

	for _, field := range requiredFields {
		if _, exists := data[field]; !exists {
			missingFields = append(missingFields, field)
		}
	}

	if len(missingFields) > 0 {
		return errors.NewValidationErrorWithDetails(
			"Required fields are missing",
			map[string]interface{}{
				"missing_fields": missingFields,
			},
		)
	}

	return nil
}

// ValidateEmail validates email format
func ValidateEmail(email string) bool {
	emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	return emailRegex.MatchString(email)
}

// ValidatePhoneNumber validates phone number format
func ValidatePhoneNumber(phone string) bool {
	// Allow various phone number formats
	phoneRegex := regexp.MustCompile(`^\+?[\d\s\-()]+$`)
	return phoneRegex.MatchString(phone) && len(strings.ReplaceAll(phone, " ", "")) >= 10
}

// ValidateURL validates URL format
func ValidateURL(url string) bool {
	return regexp.MustCompile(`^https?://`).MatchString(url)
}

// ValidateStringLength validates string length
func ValidateStringLength(value string, min, max int) bool {
	length := len(value)
	return length >= min && length <= max
}

// ValidateNumeric validates numeric value
func ValidateNumeric(value string) bool {
	return regexp.MustCompile(`^\d+$`).MatchString(value)
}

// ValidateDecimal validates decimal value
func ValidateDecimal(value string) bool {
	return regexp.MustCompile(`^\d+(\.\d+)?$`).MatchString(value)
}

// ValidateUUID validates UUID format
func ValidateUUID(uuid string) bool {
	uuidRegex := regexp.MustCompile(`^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$`)
	return uuidRegex.MatchString(uuid)
}

// FileValidationConfig holds configuration for file validation
type FileValidationConfig struct {
	MaxFileSize      int64    // Maximum file size in bytes
	AllowedMimeTypes []string // Allowed MIME types
	AllowedExtensions []string // Allowed file extensions
	Required         bool     // Is file required
}

// DefaultFileValidationConfig returns default file validation configuration
func DefaultFileValidationConfig() FileValidationConfig {
	return FileValidationConfig{
		MaxFileSize: 5 * 1024 * 1024, // 5MB
		AllowedMimeTypes: []string{
			"image/jpeg",
			"image/png",
			"image/gif",
			"image/webp",
			"application/pdf",
		},
		AllowedExtensions: []string{".jpg", ".jpeg", ".png", ".gif", ".webp", ".pdf"},
		Required: false,
	}
}

// ValidateFile validates uploaded file
func ValidateFile(fileHeader *multipart.FileHeader, config FileValidationConfig) error {
	// Check if file is required
	if config.Required && fileHeader == nil {
		return errors.NewValidationError("File is required")
	}

	// Check file size
	if fileHeader.Size > config.MaxFileSize {
		return errors.NewValidationErrorWithDetails(
			"File size exceeds limit",
			map[string]interface{}{
				"max_size": config.MaxFileSize,
				"file_size": fileHeader.Size,
			},
		)
	}

	// Check file extension
	ext := strings.ToLower(getFileExtension(fileHeader.Filename))
	if !isAllowedExtension(ext, config.AllowedExtensions) {
		return errors.NewValidationErrorWithDetails(
			"File extension not allowed",
			map[string]interface{}{
				"allowed_extensions": config.AllowedExtensions,
				"file_extension": ext,
			},
		)
	}

	// Open file to check MIME type
	file, err := fileHeader.Open()
	if err != nil {
		return errors.NewInternalError("Failed to open file")
	}
	defer file.Close()

	// Detect actual MIME type
	buffer := make([]byte, 512)
	_, err = file.Read(buffer)
	if err != nil {
		return errors.NewInternalError("Failed to read file")
	}

	mimeType := http.DetectContentType(buffer)
	if !isAllowedMimeType(mimeType, config.AllowedMimeTypes) {
		return errors.NewValidationErrorWithDetails(
			"File type not allowed",
			map[string]interface{}{
				"allowed_types": config.AllowedMimeTypes,
				"file_type": mimeType,
			},
		)
	}

	return nil
}

// getFileExtension extracts file extension from filename
func getFileExtension(filename string) string {
	parts := strings.Split(filename, ".")
	if len(parts) > 1 {
		return "." + parts[len(parts)-1]
	}
	return ""
}

// isAllowedExtension checks if extension is allowed
func isAllowedExtension(ext string, allowedExtensions []string) bool {
	for _, allowed := range allowedExtensions {
		if strings.EqualFold(ext, allowed) {
			return true
		}
	}
	return false
}

// isAllowedMimeType checks if MIME type is allowed
func isAllowedMimeType(mimeType string, allowedMimeTypes []string) bool {
	for _, allowed := range allowedMimeTypes {
		if strings.HasPrefix(mimeType, allowed) {
			return true
		}
	}
	return false
}
