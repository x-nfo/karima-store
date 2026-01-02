package middleware

import (
	"bytes"
	"image"
	"image/jpeg"
	"image/png"
	"io"
	"karima_store/internal/errors"
	"mime/multipart"
	"net/http"
	"path/filepath"
	"strings"

	"github.com/gofiber/fiber/v2"
)

// SecureFileUploadConfig holds configuration for secure file upload
type SecureFileUploadConfig struct {
	MaxFileSize         int64    // Maximum file size in bytes
	AllowedMimeTypes    []string // Allowed MIME types
	AllowedExtensions   []string // Allowed file extensions
	MaxImageWidth       int      // Maximum image width (0 = no limit)
	MaxImageHeight      int      // Maximum image height (0 = no limit)
	MinImageWidth       int      // Minimum image width (0 = no limit)
	MinImageHeight      int      // Minimum image height (0 = no limit)
	ScanForMalware      bool     // Enable malware scanning (placeholder)
	SanitizeFilename    bool     // Sanitize filename
	Required            bool     // Is file required
}

// DefaultSecureFileUploadConfig returns default secure file upload configuration
func DefaultSecureFileUploadConfig() SecureFileUploadConfig {
	return SecureFileUploadConfig{
		MaxFileSize: 5 * 1024 * 1024, // 5MB
		AllowedMimeTypes: []string{
			"image/jpeg",
			"image/png",
			"image/gif",
			"image/webp",
			"application/pdf",
		},
		AllowedExtensions: []string{".jpg", ".jpeg", ".png", ".gif", ".webp", ".pdf"},
		MaxImageWidth:     4096,
		MaxImageHeight:    4096,
		MinImageWidth:     100,
		MinImageHeight:    100,
		ScanForMalware:    false,
		SanitizeFilename:  true,
		Required:          false,
	}
}

// SecureFileUpload creates a secure file upload middleware
func SecureFileUpload(fieldName string, config SecureFileUploadConfig) fiber.Handler {
	return func(c *fiber.Ctx) error {
		file, err := c.FormFile(fieldName)
		if err != nil {
			if err.Error() == "no such file" {
				if config.Required {
					return errors.NewValidationError("File is required")
				}
				return c.Next()
			}
			return errors.NewInvalidInputError("Failed to retrieve file")
		}

		// Validate file
		if err := validateFile(file, config); err != nil {
			return err
		}

		// Sanitize filename if enabled
		if config.SanitizeFilename {
			sanitizedFilename := sanitizeFilename(file.Filename)
			c.Locals("sanitized_filename", sanitizedFilename)
		}

		return c.Next()
	}
}

// validateFile performs comprehensive file validation
func validateFile(fileHeader *multipart.FileHeader, config SecureFileUploadConfig) error {
	// Check file size
	if fileHeader.Size > config.MaxFileSize {
		return errors.NewValidationErrorWithDetails(
			"File size exceeds limit",
			map[string]interface{}{
				"max_size":  config.MaxFileSize,
				"file_size":  fileHeader.Size,
			},
		)
	}

	// Check file extension
	ext := strings.ToLower(filepath.Ext(fileHeader.Filename))
	if !isAllowedExtension(ext, config.AllowedExtensions) {
		return errors.NewValidationErrorWithDetails(
			"File extension not allowed",
			map[string]interface{}{
				"allowed_extensions": config.AllowedExtensions,
				"file_extension":     ext,
			},
		)
	}

	// Open file for validation
	file, err := fileHeader.Open()
	if err != nil {
		return errors.NewInternalError("Failed to open file")
	}
	defer file.Close()

	// Detect actual MIME type
	buffer := make([]byte, 512)
	_, err = file.Read(buffer)
	if err != nil && err != io.EOF {
		return errors.NewInternalError("Failed to read file")
	}

	mimeType := http.DetectContentType(buffer)
	if !isAllowedMimeType(mimeType, config.AllowedMimeTypes) {
		return errors.NewValidationErrorWithDetails(
			"File type not allowed",
			map[string]interface{}{
				"allowed_types": config.AllowedMimeTypes,
				"file_type":     mimeType,
			},
		)
	}

	// Validate image dimensions if it's an image
	if strings.HasPrefix(mimeType, "image/") {
		if err := validateImageDimensions(bytes.NewReader(buffer), config); err != nil {
			return err
		}
	}

	// Scan for malware (placeholder - would need actual malware scanning library)
	if config.ScanForMalware {
		if err := scanForMalware(fileHeader); err != nil {
			return errors.NewInternalError("File scan failed")
		}
	}

	return nil
}

// validateImageDimensions validates image dimensions
func validateImageDimensions(reader io.Reader, config SecureFileUploadConfig) error {
	// Reset reader
	if seeker, ok := reader.(io.Seeker); ok {
		seeker.Seek(0, io.SeekStart)
	}

	// Decode image to get dimensions
	img, format, err := image.DecodeConfig(reader)
	if err != nil {
		return errors.NewInvalidInputError("Failed to decode image")
	}

	// Validate width
	if config.MaxImageWidth > 0 && img.Width > config.MaxImageWidth {
		return errors.NewValidationErrorWithDetails(
			"Image width exceeds maximum",
			map[string]interface{}{
				"max_width": config.MaxImageWidth,
				"width":     img.Width,
			},
		)
	}

	if config.MinImageWidth > 0 && img.Width < config.MinImageWidth {
		return errors.NewValidationErrorWithDetails(
			"Image width below minimum",
			map[string]interface{}{
				"min_width": config.MinImageWidth,
				"width":     img.Width,
			},
		)
	}

	// Validate height
	if config.MaxImageHeight > 0 && img.Height > config.MaxImageHeight {
		return errors.NewValidationErrorWithDetails(
			"Image height exceeds maximum",
			map[string]interface{}{
				"max_height": config.MaxImageHeight,
				"height":     img.Height,
			},
		)
	}

	if config.MinImageHeight > 0 && img.Height < config.MinImageHeight {
		return errors.NewValidationErrorWithDetails(
			"Image height below minimum",
			map[string]interface{}{
				"min_height": config.MinImageHeight,
				"height":     img.Height,
			},
		)
	}

	return nil
}

// scanForMalware scans file for malware (placeholder implementation)
func scanForMalware(fileHeader *multipart.FileHeader) error {
	// In a real implementation, you would integrate with a malware scanning service
	// like ClamAV, VirusTotal API, or similar
	
	// Placeholder: Always return success
	// TODO: Implement actual malware scanning
	return nil
}

// sanitizeFilename sanitizes the filename to prevent path traversal and other attacks
func sanitizeFilename(filename string) string {
	// Remove directory path
	filename = filepath.Base(filename)

	// Remove dangerous characters
	dangerousChars := []string{
		"..",
		"/",
		"\\",
		":",
		"*",
		"?",
		"\"",
		"<",
		">",
		"|",
		"\x00",
	}

	sanitized := filename
	for _, char := range dangerousChars {
		sanitized = strings.ReplaceAll(sanitized, char, "")
	}

	// Remove leading and trailing whitespace and dots
	sanitized = strings.Trim(sanitized, " .")

	// If filename is empty after sanitization, use a default
	if sanitized == "" {
		sanitized = "upload"
	}

	return sanitized
}

// GenerateSecureFilename generates a secure filename with timestamp and random string
func GenerateSecureFilename(originalFilename string) string {
	ext := filepath.Ext(originalFilename)
	sanitized := sanitizeFilename(originalFilename[:len(originalFilename)-len(ext)])
	
	// In a real implementation, you would add timestamp and random string
	// For now, just return sanitized filename
	return sanitized + ext
}

// ValidateImageFile validates an image file specifically
func ValidateImageFile(fileHeader *multipart.FileHeader) error {
	config := SecureFileUploadConfig{
		MaxFileSize: 5 * 1024 * 1024,
		AllowedMimeTypes: []string{
			"image/jpeg",
			"image/png",
			"image/gif",
			"image/webp",
		},
		AllowedExtensions: []string{".jpg", ".jpeg", ".png", ".gif", ".webp"},
		MaxImageWidth:     4096,
		MaxImageHeight:    4096,
		MinImageWidth:     100,
		MinImageHeight:    100,
	}

	return validateFile(fileHeader, config)
}

// ValidatePDFFile validates a PDF file specifically
func ValidatePDFFile(fileHeader *multipart.FileHeader) error {
	config := SecureFileUploadConfig{
		MaxFileSize: 10 * 1024 * 1024, // 10MB for PDFs
		AllowedMimeTypes: []string{
			"application/pdf",
		},
		AllowedExtensions: []string{".pdf"},
	}

	return validateFile(fileHeader, config)
}

// GetFileMimeType detects the MIME type of a file
func GetFileMimeType(fileHeader *multipart.FileHeader) (string, error) {
	file, err := fileHeader.Open()
	if err != nil {
		return "", errors.NewInternalError("Failed to open file")
	}
	defer file.Close()

	buffer := make([]byte, 512)
	_, err = file.Read(buffer)
	if err != nil && err != io.EOF {
		return "", errors.NewInternalError("Failed to read file")
	}

	return http.DetectContentType(buffer), nil
}

// IsImageFile checks if a file is an image
func IsImageFile(mimeType string) bool {
	imageTypes := []string{
		"image/jpeg",
		"image/png",
		"image/gif",
		"image/webp",
		"image/bmp",
	}

	for _, imageType := range imageTypes {
		if strings.HasPrefix(mimeType, imageType) {
			return true
		}
	}
	return false
}

// OptimizeImage optimizes an image file (placeholder)
func OptimizeImage(fileHeader *multipart.FileHeader) ([]byte, error) {
	// In a real implementation, you would use an image optimization library
	// like imaging, imgo, or similar
	
	// Placeholder: Return original file content
	file, err := fileHeader.Open()
	if err != nil {
		return nil, errors.NewInternalError("Failed to open file")
	}
	defer file.Close()

	return io.ReadAll(file)
}

// ConvertToJPEG converts an image to JPEG format (placeholder)
func ConvertToJPEG(fileHeader *multipart.FileHeader) ([]byte, error) {
	// In a real implementation, you would use an image conversion library
	
	// Placeholder: Return original file content
	file, err := fileHeader.Open()
	if err != nil {
		return nil, errors.NewInternalError("Failed to open file")
	}
	defer file.Close()

	return io.ReadAll(file)
}

// CreateThumbnail creates a thumbnail of an image (placeholder)
func CreateThumbnail(fileHeader *multipart.FileHeader, width, height int) ([]byte, error) {
	// In a real implementation, you would use an image processing library
	// like imaging, resize, or similar
	
	// Placeholder: Return original file content
	file, err := fileHeader.Open()
	if err != nil {
		return nil, errors.NewInternalError("Failed to open file")
	}
	defer file.Close()

	return io.ReadAll(file)
}
