package services

import (
	"context"
	"errors"
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/karima-store/internal/config"
	"github.com/karima-store/internal/models"
	"github.com/karima-store/internal/repository"
	"github.com/karima-store/internal/storage"
)

type MediaService struct {
	mediaRepo   repository.MediaRepository
	productRepo repository.ProductRepository
	cfg         *config.Config
	r2Storage   *storage.R2Storage
}

type UploadResponse struct {
	MediaID   uint   `json:"media_id"`
	URL       string `json:"url"`
	FileName  string `json:"file_name"`
	FileSize  int64  `json:"file_size"`
	Width     int    `json:"width"`
	Height    int    `json:"height"`
}

func NewMediaService(
	mediaRepo repository.MediaRepository,
	productRepo repository.ProductRepository,
	cfg *config.Config,
) *MediaService {
	service := &MediaService{
		mediaRepo:   mediaRepo,
		productRepo: productRepo,
		cfg:         cfg,
	}

	// Initialize R2 storage if configured
	if cfg.FileStorage == "r2" && cfg.R2AccountID != "" && cfg.R2AccessKeyID != "" && cfg.R2SecretAccessKey != "" && cfg.R2BucketName != "" {
		r2Config := &storage.R2Config{
			AccountID:       cfg.R2AccountID,
			AccessKeyID:     cfg.R2AccessKeyID,
			SecretAccessKey: cfg.R2SecretAccessKey,
			BucketName:      cfg.R2BucketName,
			PublicURL:       cfg.R2PublicURL,
			Region:          cfg.R2Region,
		}
		r2Storage, err := storage.NewR2Storage(r2Config)
		if err != nil {
			fmt.Printf("Warning: Failed to initialize R2 storage: %v. Falling back to local storage.\n", err)
		} else {
			service.r2Storage = r2Storage
		}
	}

	return service
}

// UploadImage uploads an image file and creates a media record
// Supports both local storage and Cloudflare R2
func (s *MediaService) UploadImage(fileHeader *multipart.FileHeader, productID uint, position int, isPrimary bool) (*UploadResponse, error) {
	// Validate file
	if fileHeader == nil {
		return nil, errors.New("no file provided")
	}

	// Check file size (max 10MB)
	const maxFileSize = 10 * 1024 * 1024
	if fileHeader.Size > maxFileSize {
		return nil, errors.New("file size exceeds 10MB limit")
	}

	// Open file
	src, err := fileHeader.Open()
	if err != nil {
		return nil, fmt.Errorf("failed to open file: %w", err)
	}
	defer src.Close()

	// Read file content
	fileBytes, err := io.ReadAll(src)
	if err != nil {
		return nil, fmt.Errorf("failed to read file: %w", err)
	}

	// Get file extension
	ext := filepath.Ext(fileHeader.Filename)
	if ext == "" {
		return nil, errors.New("invalid file extension")
	}

	// Generate unique filename
	uniqueFilename := fmt.Sprintf("%d%s", time.Now().UnixNano(), ext)

	// Get content type
	contentType := fileHeader.Header.Get("Content-Type")
	if contentType == "" {
		contentType = "image/jpeg" // default
	}

	var filePath string
	var storageProvider string
	var storagePath string
	var publicURL string

	// Upload to storage based on configuration
	if s.r2Storage != nil {
		// Upload to R2
		key := fmt.Sprintf("products/%d/%s", productID, uniqueFilename)
		publicURL, err = s.r2Storage.UploadFile(context.Background(), key, fileBytes, contentType)
		if err != nil {
			return nil, fmt.Errorf("failed to upload to R2: %w", err)
		}
		filePath = key
		storageProvider = "r2"
		storagePath = key
	} else {
		// Upload to local storage
		uploadDir := "uploads"
		if err := os.MkdirAll(uploadDir, 0755); err != nil {
			return nil, fmt.Errorf("failed to create upload directory: %w", err)
		}

		filePath = filepath.Join(uploadDir, uniqueFilename)
		if err := os.WriteFile(filePath, fileBytes, 0644); err != nil {
			return nil, fmt.Errorf("failed to save file: %w", err)
		}
		publicURL = "/" + filePath
		storageProvider = "local"
		storagePath = filePath
	}

	// Create media record
	media := &models.Media{
		Type:            models.MediaTypeImage,
		URL:             publicURL,
		AltText:         strings.TrimSuffix(fileHeader.Filename, ext),
		Status:          models.MediaStatusActive,
		Position:        position,
		IsPrimary:       isPrimary,
		FileName:        uniqueFilename,
		FileSize:        fileHeader.Size,
		ContentType:     contentType,
		StorageProvider: storageProvider,
		StoragePath:     storagePath,
		ProductID:       productID,
	}

	if err := s.mediaRepo.Create(media); err != nil {
		// Clean up file if database insert fails
		if storageProvider == "local" {
			os.Remove(filePath)
		} else if s.r2Storage != nil {
			s.r2Storage.DeleteFile(context.Background(), storagePath)
		}
		return nil, fmt.Errorf("failed to create media record: %w", err)
	}

	// Set as primary if requested
	if isPrimary {
		if err := s.mediaRepo.SetAsPrimary(media.ID); err != nil {
			return nil, fmt.Errorf("failed to set media as primary: %w", err)
		}
	}

	return &UploadResponse{
		MediaID:  media.ID,
		URL:       media.URL,
		FileName:  uniqueFilename,
		FileSize:  fileHeader.Size,
	}, nil
}

// DeleteMedia deletes a media record and its file
func (s *MediaService) DeleteMedia(mediaID uint) error {
	// Get media record
	media, err := s.mediaRepo.GetByID(mediaID)
	if err != nil {
		return err
	}

	// Delete file from storage
	if media.StorageProvider == "local" && media.StoragePath != "" {
		if err := os.Remove(media.StoragePath); err != nil {
			// Log error but continue to delete database record
			fmt.Printf("Warning: failed to delete file %s: %v", media.StoragePath, err)
		}
	} else if media.StorageProvider == "r2" && s.r2Storage != nil && media.StoragePath != "" {
		if err := s.r2Storage.DeleteFile(context.Background(), media.StoragePath); err != nil {
			// Log error but continue to delete database record
			fmt.Printf("Warning: failed to delete file from R2 %s: %v", media.StoragePath, err)
		}
	}

	// Delete database record
	return s.mediaRepo.Delete(mediaID)
}

// UpdateMedia updates media information
func (s *MediaService) UpdateMedia(media *models.Media) error {
	return s.mediaRepo.Update(media)
}

// GetMediaByProduct retrieves all media for a product
func (s *MediaService) GetMediaByProduct(productID uint) ([]models.Media, error) {
	return s.mediaRepo.GetByProductID(productID)
}

// SetPrimaryMedia sets a media item as primary for a product
func (s *MediaService) SetPrimaryMedia(mediaID, productID uint) error {
	// Get media record
	media, err := s.mediaRepo.GetByID(mediaID)
	if err != nil {
		return err
	}

	// Verify media belongs to product
	if media.ProductID != productID {
		return errors.New("media does not belong to the specified product")
	}

	// Unset primary for all media of this product
	if err := s.mediaRepo.UnsetPrimary(productID); err != nil {
		return err
	}

	// Set primary for this media
	return s.mediaRepo.SetAsPrimary(mediaID)
}

// ValidateImageFile validates an uploaded image file
func (s *MediaService) ValidateImageFile(fileHeader *multipart.FileHeader) error {
	if fileHeader == nil {
		return errors.New("no file provided")
	}

	// Check file size (max 10MB)
	const maxFileSize = 10 * 1024 * 1024
	if fileHeader.Size > maxFileSize {
		return errors.New("file size exceeds 10MB limit")
	}

	// Check file extension
	ext := strings.ToLower(filepath.Ext(fileHeader.Filename))
	allowedExtensions := map[string]bool{
		".jpg":  true,
		".jpeg": true,
		".png":  true,
		".gif":  true,
		".webp": true,
	}

	if !allowedExtensions[ext] {
		return errors.New("invalid file type. Only JPG, PNG, GIF, and WebP are allowed")
	}

	// Check content type
	contentType := fileHeader.Header.Get("Content-Type")
	allowedContentTypes := map[string]bool{
		"image/jpeg":      true,
		"image/jpg":       true,
		"image/png":       true,
		"image/gif":       true,
		"image/webp":      true,
	}

	if !allowedContentTypes[contentType] {
		return errors.New("invalid content type")
	}

	return nil
}
