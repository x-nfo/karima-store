package repository

import (
	"github.com/karima-store/internal/models"
	"gorm.io/gorm"
)

type MediaRepository interface {
	Create(media *models.Media) error
	GetByID(id uint) (*models.Media, error)
	GetAll() ([]models.Media, error)
	GetByProductID(productID uint) ([]models.Media, error)
	Update(media *models.Media) error
	Delete(id uint) error
	SetAsPrimary(id uint) error
	UnsetPrimary(productID uint) error
}

type mediaRepository struct {
	db *gorm.DB
}

func NewMediaRepository(db *gorm.DB) MediaRepository {
	return &mediaRepository{db: db}
}

func (r *mediaRepository) Create(media *models.Media) error {
	return r.db.Create(media).Error
}

func (r *mediaRepository) GetByID(id uint) (*models.Media, error) {
	var media models.Media
	err := r.db.First(&media, id).Error
	if err != nil {
		return nil, err
	}
	return &media, nil
}

func (r *mediaRepository) GetAll() ([]models.Media, error) {
	var mediaList []models.Media
	err := r.db.Order("position ASC").Find(&mediaList).Error
	return mediaList, err
}

func (r *mediaRepository) GetByProductID(productID uint) ([]models.Media, error) {
	var mediaList []models.Media
	err := r.db.Where("product_id = ?", productID).
		Order("position ASC").
		Find(&mediaList).Error
	return mediaList, err
}

func (r *mediaRepository) Update(media *models.Media) error {
	return r.db.Save(media).Error
}

func (r *mediaRepository) Delete(id uint) error {
	return r.db.Delete(&models.Media{}, id).Error
}

// SetAsPrimary sets a media item as primary for its product
func (r *mediaRepository) SetAsPrimary(id uint) error {
	// First, unset primary flag for all media of the same product
	var media models.Media
	err := r.db.Where("id = ?", id).First(&media).Error
	if err != nil {
		return err
	}

	// Unset primary for all media of this product
	r.db.Model(&models.Media{}).
		Where("product_id = ? AND id != ?", media.ProductID, id).
		Update("is_primary", false)

	// Set primary for this media
	return r.db.Model(&models.Media{}).
		Where("id = ?", id).
		Update("is_primary", true).Error
}

// UnsetPrimary removes primary flag from all media of a product
func (r *mediaRepository) UnsetPrimary(productID uint) error {
	return r.db.Model(&models.Media{}).
		Where("product_id = ?", productID).
		Update("is_primary", false).Error
}
