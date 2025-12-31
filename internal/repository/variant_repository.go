package repository

import (
	"gorm.io/gorm"
	"github.com/karima-store/internal/models"
)

type VariantRepository interface {
	Create(variant *models.ProductVariant) error
	GetByID(id uint) (*models.ProductVariant, error)
	GetBySKU(sku string) (*models.ProductVariant, error)
	GetByProductID(productID uint) ([]models.ProductVariant, error)
	Update(variant *models.ProductVariant) error
	Delete(id uint) error
	UpdateStock(id uint, quantity int) error
}

type variantRepository struct {
	db *gorm.DB
}

func NewVariantRepository(db *gorm.DB) VariantRepository {
	return &variantRepository{db: db}
}

func (r *variantRepository) Create(variant *models.ProductVariant) error {
	return r.db.Create(variant).Error
}

func (r *variantRepository) GetByID(id uint) (*models.ProductVariant, error) {
	var variant models.ProductVariant
	err := r.db.First(&variant, id).Error
	if err != nil {
		return nil, err
	}
	return &variant, nil
}

func (r *variantRepository) GetBySKU(sku string) (*models.ProductVariant, error) {
	var variant models.ProductVariant
	err := r.db.Where("sku = ?", sku).First(&variant).Error
	if err != nil {
		return nil, err
	}
	return &variant, nil
}

func (r *variantRepository) GetByProductID(productID uint) ([]models.ProductVariant, error) {
	var variants []models.ProductVariant
	err := r.db.Where("product_id = ?", productID).Find(&variants).Error
	return variants, err
}

func (r *variantRepository) Update(variant *models.ProductVariant) error {
	return r.db.Save(variant).Error
}

func (r *variantRepository) Delete(id uint) error {
	return r.db.Delete(&models.ProductVariant{}, id).Error
}

func (r *variantRepository) UpdateStock(id uint, quantity int) error {
	return r.db.Model(&models.ProductVariant{}).
		Where("id = ?", id).
		Update("stock", gorm.Expr("stock + ?", quantity)).Error
}
