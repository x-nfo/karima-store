package repository

import (
	"github.com/karima-store/internal/models"
	"gorm.io/gorm"
)

type ShippingZoneRepository interface {
	Create(zone *models.ShippingZone) error
	GetByID(id uint) (*models.ShippingZone, error)
	GetAll() ([]models.ShippingZone, error)
	GetActive() ([]models.ShippingZone, error)
	Update(zone *models.ShippingZone) error
	Delete(id uint) error
	GetByRegion(regionCode string) (*models.ShippingZone, error)
}

type shippingZoneRepository struct {
	db *gorm.DB
}

func NewShippingZoneRepository(db *gorm.DB) ShippingZoneRepository {
	return &shippingZoneRepository{db: db}
}

func (r *shippingZoneRepository) Create(zone *models.ShippingZone) error {
	return r.db.Create(zone).Error
}

func (r *shippingZoneRepository) GetByID(id uint) (*models.ShippingZone, error) {
	var zone models.ShippingZone
	err := r.db.First(&zone, id).Error
	if err != nil {
		return nil, err
	}
	return &zone, nil
}

func (r *shippingZoneRepository) GetAll() ([]models.ShippingZone, error) {
	var zones []models.ShippingZone
	err := r.db.Order("created_at DESC").Find(&zones).Error
	return zones, err
}

func (r *shippingZoneRepository) GetActive() ([]models.ShippingZone, error) {
	var zones []models.ShippingZone
	err := r.db.Where("status = ?", models.ShippingZoneActive).
		Order("created_at DESC").
		Find(&zones).Error
	return zones, err
}

func (r *shippingZoneRepository) Update(zone *models.ShippingZone) error {
	return r.db.Save(zone).Error
}

func (r *shippingZoneRepository) Delete(id uint) error {
	return r.db.Delete(&models.ShippingZone{}, id).Error
}

// GetByRegion finds the shipping zone that applies to a specific region
func (r *shippingZoneRepository) GetByRegion(regionCode string) (*models.ShippingZone, error) {
	var zones []models.ShippingZone
	err := r.db.Where("status = ?", models.ShippingZoneActive).
		Order("created_at DESC").
		Find(&zones).Error
	if err != nil {
		return nil, err
	}

	// Find the first zone that includes this region and doesn't exclude it
	for _, zone := range zones {
		// Check if region is in zone's regions
		for _, region := range zone.Regions {
			if region == regionCode {
				// Check if region is not excluded
				excluded := false
				for _, excludeRegion := range zone.ExcludeRegions {
					if excludeRegion == regionCode {
						excluded = true
						break
					}
				}
				if !excluded {
					return &zone, nil
				}
			}
		}
	}

	return nil, gorm.ErrRecordNotFound
}
