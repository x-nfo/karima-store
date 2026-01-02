package middleware

import (
	"bytes"
	"context"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"image"
	"image/jpeg"
	"image/png"
	"io"
	"karima_store/internal/errors"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/dutchcoders/go-clamd"
	"github.com/gofiber/fiber/v2"
)

// SecureFileUploadConfig holds configuration for secure file upload
type SecureFileUploadConfig struct {
	MaxFileSize       int64          // Maximum file size in bytes
	AllowedMimeTypes  []string       // Allowed MIME types
	AllowedExtensions []string       // Allowed file extensions
	MaxImageWidth     int            // Maximum image width (0 = no limit)
	MaxImageHeight    int            // Maximum image height (0 = no limit)
	MinImageWidth     int            // Minimum image width (0 = no limit)
	MinImageHeight    int            // Minimum image height (0 = no limit)
	ScanForMalware    bool           // Enable malware scanning
	SanitizeFilename  bool           // Sanitize filename
	Required          bool           // Is file required
	MalwareScanner    MalwareScanner // Malware scanner implementation
	FailOpen          bool           // If true, allow file upload when scanner fails (fail-open); if false, block upload (fail-closed)
}

// MalwareScanner defines the interface for malware scanning services
type MalwareScanner interface {
	ScanFile(ctx context.Context, file io.Reader, filename string) (*ScanResult, error)
	GetScanResult(ctx context.Context, scanID string) (*ScanResult, error)
	QuarantineFile(ctx context.Context, scanID string) error
}

// ScanResult represents the result of a malware scan
type ScanResult struct {
	ScanID       string        `json:"scan_id"`
	Filename     string        `json:"filename"`
	FileHash     string        `json:"file_hash"`
	IsClean      bool          `json:"is_clean"`
	Threats      []string      `json:"threats,omitempty"`
	ScannedAt    time.Time     `json:"scanned_at"`
	ScanDuration time.Duration `json:"scan_duration"`
	ScannerName  string        `json:"scanner_name"`
}

// ClamAVScanner implements MalwareScanner using ClamAV
type ClamAVScanner struct {
	endpoint    string
	timeout     time.Duration
	maxFileSize int64
}

// NewClamAVScanner creates a new ClamAV scanner instance
func NewClamAVScanner(endpoint string, timeout time.Duration, maxFileSize int64) *ClamAVScanner {
	return &ClamAVScanner{
		endpoint:    endpoint,
		timeout:     timeout,
		maxFileSize: maxFileSize,
	}
}

// ScanFile scans a file using ClamAV
func (s *ClamAVScanner) ScanFile(ctx context.Context, file io.Reader, filename string) (*ScanResult, error) {
	startTime := time.Now()

	// Read file content for hashing
	content, err := io.ReadAll(file)
	if err != nil {
		return nil, fmt.Errorf("failed to read file: %w", err)
	}

	// Calculate file hash
	hash := sha256.Sum256(content)
	fileHash := hex.EncodeToString(hash[:])

	// Create scan result
	result := &ScanResult{
		ScanID:      generateScanID(),
		Filename:    filename,
		FileHash:    fileHash,
		ScannedAt:   startTime,
		ScannerName: "ClamAV",
	}

	// In production, this would connect to ClamAV daemon
	// For now, we'll implement basic checks
	if err := s.scanWithClamAV(ctx, content); err != nil {
		result.IsClean = false
		result.Threats = []string{err.Error()}
		result.ScanDuration = time.Since(startTime)
		return result, nil
	}

	result.IsClean = true
	result.ScanDuration = time.Since(startTime)
	return result, nil
}

// scanWithClamAV performs actual ClamAV scan
func (s *ClamAVScanner) scanWithClamAV(ctx context.Context, content []byte) error {
	// Create ClamAV client
	client := clamd.NewClamd(s.endpoint)

	// Create a channel to receive scan results
	resultChan := client.ScanStream(ctx, bytes.NewReader(content), nil)

	// Wait for scan result with timeout
	select {
	case result := <-resultChan:
		if result == nil {
			return fmt.Errorf("clamAV scan returned no result")
		}

		// Check scan status
		if result.Status == clamd.RES_FOUND {
			return fmt.Errorf("malware detected: %s", result.Description)
		} else if result.Status == clamd.RES_ERROR {
			return fmt.Errorf("clamAV scan error: %s", result.Description)
		} else if result.Status == clamd.RES_PARSE_ERROR {
			return fmt.Errorf("clamAV parse error: %s", result.Description)
		}
		// RES_OK means file is clean
		return nil

	case <-ctx.Done():
		return fmt.Errorf("clamAV scan timeout: %v", ctx.Err())
	}
}

// GetScanResult retrieves a scan result by ID
func (s *ClamAVScanner) GetScanResult(ctx context.Context, scanID string) (*ScanResult, error) {
	// In production, this would query a database or cache for scan results
	// For now, return not found
	return nil, fmt.Errorf("scan result not found")
}

// QuarantineFile moves a file to quarantine
func (s *ClamAVScanner) QuarantineFile(ctx context.Context, scanID string) error {
	// In production, this would move the file to a quarantine directory
	// For now, just log
	return nil
}

// VirusTotalScanner implements MalwareScanner using VirusTotal API
type VirusTotalScanner struct {
	apiKey      string
	endpoint    string
	timeout     time.Duration
	maxFileSize int64
}

// NewVirusTotalScanner creates a new VirusTotal scanner instance
func NewVirusTotalScanner(apiKey, endpoint string, timeout time.Duration, maxFileSize int64) *VirusTotalScanner {
	return &VirusTotalScanner{
		apiKey:      apiKey,
		endpoint:    endpoint,
		timeout:     timeout,
		maxFileSize: maxFileSize,
	}
}

// ScanFile scans a file using VirusTotal API
func (s *VirusTotalScanner) ScanFile(ctx context.Context, file io.Reader, filename string) (*ScanResult, error) {
	startTime := time.Now()

	// Read file content for hashing
	content, err := io.ReadAll(file)
	if err != nil {
		return nil, fmt.Errorf("failed to read file: %w", err)
	}

	// Calculate file hash
	hash := sha256.Sum256(content)
	fileHash := hex.EncodeToString(hash[:])

	// Create scan result
	result := &ScanResult{
		ScanID:      generateScanID(),
		Filename:    filename,
		FileHash:    fileHash,
		ScannedAt:   startTime,
		ScannerName: "VirusTotal",
	}

	// In production, this would call VirusTotal API
	// For now, we'll implement basic checks
	if err := s.scanWithVirusTotal(ctx, content, fileHash); err != nil {
		result.IsClean = false
		result.Threats = []string{err.Error()}
		result.ScanDuration = time.Since(startTime)
		return result, nil
	}

	result.IsClean = true
	result.ScanDuration = time.Since(startTime)
	return result, nil
}

// scanWithVirusTotal performs actual VirusTotal scan
func (s *VirusTotalScanner) scanWithVirusTotal(ctx context.Context, content []byte, fileHash string) error {
	// In production implementation:
	// 1. Upload file to VirusTotal API
	// 2. Get scan ID
	// 3. Poll for scan results
	// 4. Parse and return result

	// Placeholder implementation - always return clean
	// TODO: Implement actual VirusTotal API integration
	return nil
}

// GetScanResult retrieves a scan result by ID
func (s *VirusTotalScanner) GetScanResult(ctx context.Context, scanID string) (*ScanResult, error) {
	// In production, this would query VirusTotal API for scan results
	// For now, return not found
	return nil, fmt.Errorf("scan result not found")
}

// QuarantineFile moves a file to quarantine
func (s *VirusTotalScanner) QuarantineFile(ctx context.Context, scanID string) error {
	// In production, this would move the file to a quarantine directory
	// For now, just log
	return nil
}

// LocalScanner implements a basic local malware scanner
type LocalScanner struct {
	enabled        bool
	maxFileSize    int64
	quarantinePath string
	scanTimeout    time.Duration
}

// NewLocalScanner creates a new local scanner instance
func NewLocalScanner(enabled bool, maxFileSize int64, quarantinePath string, timeout time.Duration) *LocalScanner {
	return &LocalScanner{
		enabled:        enabled,
		maxFileSize:    maxFileSize,
		quarantinePath: quarantinePath,
		scanTimeout:    timeout,
	}
}

// ScanFile scans a file locally
func (s *LocalScanner) ScanFile(ctx context.Context, file io.Reader, filename string) (*ScanResult, error) {
	startTime := time.Now()

	if !s.enabled {
		return &ScanResult{
			ScanID:       generateScanID(),
			Filename:     filename,
			IsClean:      true,
			ScannedAt:    startTime,
			ScanDuration: time.Since(startTime),
			ScannerName:  "LocalScanner",
		}, nil
	}

	// Read file content for hashing
	content, err := io.ReadAll(file)
	if err != nil {
		return nil, fmt.Errorf("failed to read file: %w", err)
	}

	// Calculate file hash
	hash := sha256.Sum256(content)
	fileHash := hex.EncodeToString(hash[:])

	// Create scan result
	result := &ScanResult{
		ScanID:      generateScanID(),
		Filename:    filename,
		FileHash:    fileHash,
		ScannedAt:   startTime,
		ScannerName: "LocalScanner",
	}

	// Perform basic local checks
	threats := s.performLocalChecks(content, filename)
	if len(threats) > 0 {
		result.IsClean = false
		result.Threats = threats
		result.ScanDuration = time.Since(startTime)
		return result, nil
	}

	result.IsClean = true
	result.ScanDuration = time.Since(startTime)
	return result, nil
}

// performLocalChecks performs basic local security checks
func (s *LocalScanner) performLocalChecks(content []byte, filename string) []string {
	var threats []string

	// Check for suspicious patterns
	suspiciousPatterns := [][]byte{
		[]byte("<script"),
		[]byte("javascript:"),
		[]byte("eval("),
		[]byte("document.cookie"),
	}

	for _, pattern := range suspiciousPatterns {
		if bytes.Contains(content, pattern) {
			threats = append(threats, fmt.Sprintf("Suspicious pattern detected: %s", string(pattern)))
		}
	}

	// Check for double extensions (e.g., file.jpg.exe)
	ext := filepath.Ext(filename)
	baseName := strings.TrimSuffix(filename, ext)
	if filepath.Ext(baseName) != "" {
		threats = append(threats, "Double extension detected - potential malware")
	}

	return threats
}

// GetScanResult retrieves a scan result by ID
func (s *LocalScanner) GetScanResult(ctx context.Context, scanID string) (*ScanResult, error) {
	// Local scanner doesn't store results
	return nil, fmt.Errorf("scan result not found")
}

// QuarantineFile moves a file to quarantine
func (s *LocalScanner) QuarantineFile(ctx context.Context, scanID string) error {
	// In production, this would move the file to quarantinePath
	return nil
}

// generateScanID generates a unique scan ID
func generateScanID() string {
	return fmt.Sprintf("scan_%d", time.Now().UnixNano())
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
		MalwareScanner:    NewLocalScanner(false, 10*1024*1024, "/tmp/quarantine", 30*time.Second),
		FailOpen:          false, // Default to fail-closed (block upload if scanner fails)
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
				"file_size": fileHeader.Size,
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

	// Scan for malware if enabled and scanner is configured
	if config.ScanForMalware && config.MalwareScanner != nil {
		if err := scanForMalware(fileHeader, config.MalwareScanner, config.FailOpen); err != nil {
			return err
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

// scanForMalware scans file for malware using configured scanner
func scanForMalware(fileHeader *multipart.FileHeader, scanner MalwareScanner, failOpen bool) error {
	// Open file for scanning
	file, err := fileHeader.Open()
	if err != nil {
		return errors.NewInternalError("Failed to open file for scanning")
	}
	defer file.Close()

	// Create context with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Scan file
	result, err := scanner.ScanFile(ctx, file, fileHeader.Filename)
	if err != nil {
		// Handle scanner failure based on fail_open configuration
		if failOpen {
			// Fail-open: Log warning but allow the file upload
			fmt.Printf("WARNING: Malware scan failed but allowing upload (fail-open mode): %v\n", err)
			return nil
		}
		// Fail-closed: Block the file upload
		return errors.NewInternalErrorWithDetails(
			"Malware scan failed - upload blocked (fail-closed mode)",
			map[string]interface{}{
				"error": err.Error(),
			},
		)
	}

	// Check if file is clean
	if !result.IsClean {
		// Delete temporary file if it exists
		if tempFile, ok := file.(*os.File); ok {
			tempPath := tempFile.Name()
			if err := os.Remove(tempPath); err != nil {
				fmt.Printf("Failed to delete temporary file %s: %v\n", tempPath, err)
			}
		}

		// Quarantine the file
		if err := scanner.QuarantineFile(ctx, result.ScanID); err != nil {
			// Log quarantine failure but still reject the file
			fmt.Printf("Failed to quarantine file: %v\n", err)
		}

		return errors.NewValidationErrorWithDetails(
			"File contains malware and has been rejected",
			map[string]interface{}{
				"scan_id":       result.ScanID,
				"threats":       result.Threats,
				"file_hash":     result.FileHash,
				"scanner":       result.ScannerName,
				"scan_duration": result.ScanDuration.String(),
			},
		)
	}

	// Log successful scan
	fmt.Printf("File scan completed successfully: %s (ID: %s, Duration: %s)\n",
		fileHeader.Filename, result.ScanID, result.ScanDuration)

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
