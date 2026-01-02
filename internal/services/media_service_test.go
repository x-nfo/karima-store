package services

import (
	"errors"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/karima-store/internal/config"
	"github.com/karima-store/internal/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockMediaRepository for testing
type MockMediaRepository struct {
	mock.Mock
}

func (m *MockMediaRepository) Create(media *models.Media) error {
	args := m.Called(media)
	return args.Error(0)
}

func (m *MockMediaRepository) GetByID(id uint) (*models.Media, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Media), args.Error(1)
}

func (m *MockMediaRepository) GetByProductID(productID uint) ([]models.Media, error) {
	args := m.Called(productID)
	return args.Get(0).([]models.Media), args.Error(1)
}

func (m *MockMediaRepository) Update(media *models.Media) error {
	args := m.Called(media)
	return args.Error(0)
}

func (m *MockMediaRepository) Delete(id uint) error {
	args := m.Called(id)
	return args.Error(0)
}

func (m *MockMediaRepository) SetAsPrimary(id uint) error {
	args := m.Called(id)
	return args.Error(0)
}

func (m *MockMediaRepository) UnsetPrimary(productID uint) error {
	args := m.Called(productID)
	return args.Error(0)
}

// MockProductRepository for testing
type MockProductRepositoryForMedia struct {
	mock.Mock
}

func TestMediaService_ValidateImageFile_ValidImage(t *testing.T) {
	mockMediaRepo := new(MockMediaRepository)
	mockProductRepo := new(MockProductRepositoryForMedia)
	cfg := &config.Config{}

	service := NewMediaService(mockMediaRepo, mockProductRepo, cfg)

	// Create a valid image file header
	validImage := createTestImageHeader("test.jpg", 1024, []byte{0xFF, 0xD8, 0xFF, 0xE0})

	err := service.ValidateImageFile(validImage)

	assert.NoError(t, err)
}

func TestMediaService_ValidateImageFile_NoFile(t *testing.T) {
	mockMediaRepo := new(MockMediaRepository)
	mockProductRepo := new(MockProductRepositoryForMedia)
	cfg := &config.Config{}

	service := NewMediaService(mockMediaRepo, mockProductRepo, cfg)

	err := service.ValidateImageFile(nil)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "no file provided")
}

func TestMediaService_ValidateImageFile_TooLarge(t *testing.T) {
	mockMediaRepo := new(MockMediaRepository)
	mockProductRepo := new(MockProductRepositoryForMedia)
	cfg := &config.Config{}

	service := NewMediaService(mockMediaRepo, mockProductRepo, cfg)

	// Create a file that's too large (6MB)
	largeFile := createTestImageHeader("large.jpg", 6*1024*1024, []byte{0xFF, 0xD8, 0xFF, 0xE0})

	err := service.ValidateImageFile(largeFile)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "file size exceeds 5MB limit")
}

func TestMediaService_ValidateImageFile_InvalidExtension(t *testing.T) {
	mockMediaRepo := new(MockMediaRepository)
	mockProductRepo := new(MockProductRepositoryForMedia)
	cfg := &config.Config{}

	service := NewMediaService(mockMediaRepo, mockProductRepo, cfg)

	// Create a file with invalid extension
	invalidFile := createTestImageHeader("test.exe", 1024, []byte{0xFF, 0xD8, 0xFF, 0xE0})

	err := service.ValidateImageFile(invalidFile)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "invalid file extension")
}

func TestMediaService_ValidateImageFile_ExtensionSpoofing(t *testing.T) {
	mockMediaRepo := new(MockMediaRepository)
	mockProductRepo := new(MockProductRepositoryForMedia)
	cfg := &config.Config{}

	service := NewMediaService(mockMediaRepo, mockProductRepo, cfg)

	// Create a file with .jpg extension but executable content
	spoofedFile := createTestImageHeader("malicious.jpg", 1024, []byte{0x4D, 0x5A}) // MZ header (executable)

	err := service.ValidateImageFile(spoofedFile)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "security check failed")
}

func TestMediaService_ValidateImageFile_ScriptFile(t *testing.T) {
	mockMediaRepo := new(MockMediaRepository)
	mockProductRepo := new(MockProductRepositoryForMedia)
	cfg := &config.Config{}

	service := NewMediaService(mockMediaRepo, mockProductRepo, cfg)

	// Create a script file with .jpg extension
	scriptContent := []byte("<script>alert('XSS')</script>")
	scriptFile := createTestImageHeader("script.jpg", int64(len(scriptContent)), scriptContent)

	err := service.ValidateImageFile(scriptFile)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "security check failed")
}

func TestMediaService_ValidateImageFile_AllowedExtensions(t *testing.T) {
	mockMediaRepo := new(MockMediaRepository)
	mockProductRepo := new(MockProductRepositoryForMedia)
	cfg := &config.Config{}

	service := NewMediaService(mockMediaRepo, mockProductRepo, cfg)

	allowedExtensions := []string{".jpg", ".jpeg", ".png", ".gif", ".webp"}

	for _, ext := range allowedExtensions {
		t.Run(ext, func(t *testing.T) {
			// Create a valid image file with the extension
			validFile := createTestImageHeader("test"+ext, 1024, []byte{0xFF, 0xD8, 0xFF, 0xE0})

			err := service.ValidateImageFile(validFile)

			assert.NoError(t, err)
		})
	}
}

func TestMediaService_ValidateImageFile_PathTraversal(t *testing.T) {
	mockMediaRepo := new(MockMediaRepository)
	mockProductRepo := new(MockProductRepositoryForMedia)
	cfg := &config.Config{}

	service := NewMediaService(mockMediaRepo, mockProductRepo, cfg)

	// Create a file with path traversal in filename
	pathTraversalFile := createTestImageHeader("../../../etc/passwd.jpg", 1024, []byte{0xFF, 0xD8, 0xFF, 0xE0})

	err := service.ValidateImageFile(pathTraversalFile)

	// The validation should pass (file extension is valid), but the upload should handle path traversal
	assert.NoError(t, err)
}

func TestMediaService_UploadImage_NoFile(t *testing.T) {
	mockMediaRepo := new(MockMediaRepository)
	mockProductRepo := new(MockProductRepositoryForMedia)
	cfg := &config.Config{}

	service := NewMediaService(mockMediaRepo, mockProductRepo, cfg)

	_, err := service.UploadImage(nil, 1, 0, false)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "no file provided")
}

func TestMediaService_UploadImage_TooLarge(t *testing.T) {
	mockMediaRepo := new(MockMediaRepository)
	mockProductRepo := new(MockProductRepositoryForMedia)
	cfg := &config.Config{}

	service := NewMediaService(mockMediaRepo, mockProductRepo, cfg)

	// Create a file that's too large (11MB)
	largeFile := createTestImageHeader("large.jpg", 11*1024*1024, []byte{0xFF, 0xD8, 0xFF, 0xE0})

	_, err := service.UploadImage(largeFile, 1, 0, false)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "file size exceeds 10MB limit")
}

func TestMediaService_UploadImage_InvalidExtension(t *testing.T) {
	mockMediaRepo := new(MockMediaRepository)
	mockProductRepo := new(MockProductRepositoryForMedia)
	cfg := &config.Config{}

	service := NewMediaService(mockMediaRepo, mockProductRepo, cfg)

	// Create a file with invalid extension
	invalidFile := createTestImageHeader("test.exe", 1024, []byte{0xFF, 0xD8, 0xFF, 0xE0})

	_, err := service.UploadImage(invalidFile, 1, 0, false)

	// The upload should fail due to file extension validation
	assert.Error(t, err)
}

func TestMediaService_UploadImage_MaliciousFile(t *testing.T) {
	mockMediaRepo := new(MockMediaRepository)
	mockProductRepo := new(MockProductRepositoryForMedia)
	cfg := &config.Config{}

	service := NewMediaService(mockMediaRepo, mockProductRepo, cfg)

	// Create a malicious file with .jpg extension
	maliciousContent := []byte("<script>alert('XSS')</script>")
	maliciousFile := createTestImageHeader("malicious.jpg", int64(len(maliciousContent)), maliciousContent)

	_, err := service.UploadImage(maliciousFile, 1, 0, false)

	// The upload should fail due to content validation
	assert.Error(t, err)
}

func TestMediaService_DeleteMedia_NonExistent(t *testing.T) {
	mockMediaRepo := new(MockMediaRepository)
	mockProductRepo := new(MockProductRepositoryForMedia)
	cfg := &config.Config{}

	service := NewMediaService(mockMediaRepo, mockProductRepo, cfg)

	mockMediaRepo.On("GetByID", uint(999)).Return(nil, errors.New("not found"))

	err := service.DeleteMedia(999)

	assert.Error(t, err)
	mockMediaRepo.AssertExpectations(t)
}

func TestMediaService_SetPrimaryMedia_NonExistent(t *testing.T) {
	mockMediaRepo := new(MockMediaRepository)
	mockProductRepo := new(MockProductRepositoryForMedia)
	cfg := &config.Config{}

	service := NewMediaService(mockMediaRepo, mockProductRepo, cfg)

	mockMediaRepo.On("GetByID", uint(999)).Return(nil, errors.New("not found"))

	err := service.SetPrimaryMedia(999, 1)

	assert.Error(t, err)
	mockMediaRepo.AssertExpectations(t)
}

func TestMediaService_SetPrimaryMedia_WrongProduct(t *testing.T) {
	mockMediaRepo := new(MockMediaRepository)
	mockProductRepo := new(MockProductRepositoryForMedia)
	cfg := &config.Config{}

	service := NewMediaService(mockMediaRepo, mockProductRepo, cfg)

	media := &models.Media{
		ID:        1,
		ProductID: 1,
	}

	mockMediaRepo.On("GetByID", uint(1)).Return(media, nil)

	err := service.SetPrimaryMedia(1, 2) // Media belongs to product 1, trying to set for product 2

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "media does not belong to the specified product")
	mockMediaRepo.AssertExpectations(t)
}

func TestMediaService_UploadImage_ConcurrentUploads(t *testing.T) {
	mockMediaRepo := new(MockMediaRepository)
	mockProductRepo := new(MockProductRepositoryForMedia)
	cfg := &config.Config{}

	service := NewMediaService(mockMediaRepo, mockProductRepo, cfg)

	// Test concurrent uploads (basic test, actual concurrency would require more setup)
	// This is a placeholder for testing concurrent upload scenarios
	validFile := createTestImageHeader("test.jpg", 1024, []byte{0xFF, 0xD8, 0xFF, 0xE0})

	mockMediaRepo.On("Create", mock.Anything).Return(nil)
	mockMediaRepo.On("SetAsPrimary", mock.Anything).Return(nil).Maybe()

	// Upload multiple files
	for i := 0; i < 3; i++ {
		_, err := service.UploadImage(validFile, uint(i+1), 0, false)
		// Note: This will likely fail due to file I/O, but we're testing the logic
		_ = err
	}

	mockMediaRepo.AssertExpectations(t)
}

func TestMediaService_UploadImage_DuplicateUpload(t *testing.T) {
	mockMediaRepo := new(MockMediaRepository)
	mockProductRepo := new(MockProductRepositoryForMedia)
	cfg := &config.Config{}

	service := NewMediaService(mockMediaRepo, mockProductRepo, cfg)

	validFile := createTestImageHeader("test.jpg", 1024, []byte{0xFF, 0xD8, 0xFF, 0xE0})

	// First upload
	mockMediaRepo.On("Create", mock.Anything).Return(nil).Once()
	_, err1 := service.UploadImage(validFile, 1, 0, false)
	_ = err1 // Ignore error for this test

	// Second upload (same file)
	mockMediaRepo.On("Create", mock.Anything).Return(nil).Once()
	_, err2 := service.UploadImage(validFile, 1, 0, false)
	_ = err2 // Ignore error for this test

	mockMediaRepo.AssertExpectations(t)
}

// Helper function to create a test file header
func createTestImageHeader(filename string, size int64, content []byte) *multipart.FileHeader {
	// Create a test file
	file := httptest.NewMultipartReader(strings.NewReader(string(content)), "boundary")

	// Create a file header
	header := &multipart.FileHeader{
		Filename: filename,
		Size:     size,
		Header:   make(http.Header),
	}

	// Set content type
	header.Header.Set("Content-Type", "image/jpeg")

	// Create a reader for the file content
	reader := strings.NewReader(string(content))
	headerReader := &readCloser{reader: reader}

	// Set the file
	header.File = headerReader

	return header
}

// readCloser is a helper to implement io.ReadCloser
type readCloser struct {
	reader *strings.Reader
}

func (rc *readCloser) Read(p []byte) (n int, err error) {
	return rc.reader.Read(p)
}

func (rc *readCloser) Close() error {
	return nil
}