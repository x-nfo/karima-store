package repository

import (
	"fmt"
	"strings"

	"github.com/karima-store/internal/models"
	"gorm.io/gorm"
)

type ProductRepository interface {
	Create(product *models.Product) error
	GetByID(id uint) (*models.Product, error)
	GetBySlug(slug string) (*models.Product, error)
	GetAll(limit, offset int, filters map[string]interface{}) ([]models.Product, int64, error)
	Update(product *models.Product) error
	Delete(id uint) error
	UpdateStock(id uint, quantity int) error
	IncrementViewCount(id uint) error
	Search(query string, limit, offset int) ([]models.Product, int64, error)
	GetByCategory(category models.ProductCategory, limit, offset int) ([]models.Product, int64, error)
	GetFeatured(limit int) ([]models.Product, error)
	GetBestSellers(limit int) ([]models.Product, error)
	WithTx(tx *gorm.DB) ProductRepository
}

type productRepository struct {
	db *gorm.DB
}

func NewProductRepository(db *gorm.DB) ProductRepository {
	return &productRepository{db: db}
}

func (r *productRepository) WithTx(tx *gorm.DB) ProductRepository {
	return &productRepository{db: tx}
}

func (r *productRepository) Create(product *models.Product) error {
	return r.db.Create(product).Error
}

func (r *productRepository) GetByID(id uint) (*models.Product, error) {
	var product models.Product
	err := r.db.Preload("Media").Preload("Variants").First(&product, id).Error
	if err != nil {
		return nil, err
	}
	return &product, nil
}

func (r *productRepository) GetBySlug(slug string) (*models.Product, error) {
	var product models.Product
	err := r.db.Preload("Media").Preload("Variants").Where("slug = ?", slug).First(&product).Error
	if err != nil {
		return nil, err
	}
	return &product, nil
}

func (r *productRepository) GetAll(limit, offset int, filters map[string]interface{}) ([]models.Product, int64, error) {
	var products []models.Product
	var total int64

	query := r.db.Model(&models.Product{})

	// Apply filters
	for key, value := range filters {
		switch key {
		case "category":
			query = query.Where("category = ?", value)
		case "status":
			query = query.Where("status = ?", value)
		case "min_price":
			query = query.Where("price >= ?", value)
		case "max_price":
			query = query.Where("price <= ?", value)
		case "brand":
			query = query.Where("brand = ?", value)
		}
	}

	// Count total
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Get products with pagination
	err := query.Preload("Media").Preload("Variants").
		Order("created_at DESC").
		Limit(limit).
		Offset(offset).
		Find(&products).Error

	return products, total, err
}

func (r *productRepository) Update(product *models.Product) error {
	return r.db.Save(product).Error
}

func (r *productRepository) Delete(id uint) error {
	return r.db.Delete(&models.Product{}, id).Error
}

func (r *productRepository) UpdateStock(id uint, quantity int) error {
	if quantity < 0 {
		result := r.db.Model(&models.Product{}).
			Where("id = ? AND stock >= ?", id, -quantity).
			Update("stock", gorm.Expr("stock + ?", quantity))

		if result.Error != nil {
			return result.Error
		}
		if result.RowsAffected == 0 {
			return fmt.Errorf("insufficient stock")
		}
		return nil
	}

	return r.db.Model(&models.Product{}).
		Where("id = ?", id).
		Update("stock", gorm.Expr("stock + ?", quantity)).Error
}

func (r *productRepository) IncrementViewCount(id uint) error {
	return r.db.Model(&models.Product{}).
		Where("id = ?", id).
		UpdateColumn("view_count", gorm.Expr("view_count + 1")).Error
}

func (r *productRepository) Search(query string, limit, offset int) ([]models.Product, int64, error) {
	var products []models.Product
	var total int64

	searchQuery := fmt.Sprintf("%%%s%%", strings.ToLower(query))

	countQuery := r.db.Model(&models.Product{}).
		Where("LOWER(name) LIKE ? OR LOWER(description) LIKE ? OR LOWER(brand) LIKE ?",
			searchQuery, searchQuery, searchQuery)

	if err := countQuery.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	err := r.db.Preload("Media").
		Where("LOWER(name) LIKE ? OR LOWER(description) LIKE ? OR LOWER(brand) LIKE ?",
			searchQuery, searchQuery, searchQuery).
		Order("created_at DESC").
		Limit(limit).
		Offset(offset).
		Find(&products).Error

	return products, total, err
}

func (r *productRepository) GetByCategory(category models.ProductCategory, limit, offset int) ([]models.Product, int64, error) {
	var products []models.Product
	var total int64

	query := r.db.Model(&models.Product{}).Where("category = ?", category)

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	err := query.Preload("Media").
		Order("created_at DESC").
		Limit(limit).
		Offset(offset).
		Find(&products).Error

	return products, total, err
}

func (r *productRepository) GetFeatured(limit int) ([]models.Product, error) {
	var products []models.Product
	err := r.db.Preload("Media").
		Where("status = ? AND is_featured = ?", models.StatusAvailable, true).
		Order("view_count DESC, sold_count DESC").
		Limit(limit).
		Find(&products).Error
	return products, err
}

func (r *productRepository) GetBestSellers(limit int) ([]models.Product, error) {
	var products []models.Product
	err := r.db.Preload("Media").
		Where("status = ?", models.StatusAvailable).
		Order("sold_count DESC").
		Limit(limit).
		Find(&products).Error
	return products, err
}
