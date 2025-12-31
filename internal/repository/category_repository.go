package repository

import (
	"gorm.io/gorm"
	"github.com/karima-store/internal/models"
)

type CategoryRepository interface {
	GetAllCategories() []models.ProductCategory
	GetCategoryStats() ([]CategoryStats, error)
}

type CategoryStats struct {
	Category     models.ProductCategory `json:"category"`
	ProductCount int64                  `json:"product_count"`
}

type categoryRepository struct {
	db *gorm.DB
}

func NewCategoryRepository(db *gorm.DB) CategoryRepository {
	return &categoryRepository{db: db}
}

func (r *categoryRepository) GetAllCategories() []models.ProductCategory {
	return []models.ProductCategory{
		models.CategoryTops,
		models.CategoryBottoms,
		models.CategoryDresses,
		models.CategoryOuterwear,
		models.CategoryFootwear,
		models.CategoryAccessories,
	}
}

func (r *categoryRepository) GetCategoryStats() ([]CategoryStats, error) {
	var stats []CategoryStats

	err := r.db.Model(&models.Product{}).
		Select("category, COUNT(*) as product_count").
		Group("category").
		Find(&stats).Error

	return stats, err
}
