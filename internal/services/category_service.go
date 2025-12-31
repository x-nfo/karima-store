package services

import (
	"github.com/karima-store/internal/models"
	"github.com/karima-store/internal/repository"
)

type CategoryService interface {
	GetAllCategories() []models.ProductCategory
	GetCategoryStats() ([]repository.CategoryStats, error)
	GetCategoryName(category models.ProductCategory) string
	IsValidCategory(category models.ProductCategory) bool
}

type categoryService struct {
	categoryRepo repository.CategoryRepository
}

func NewCategoryService(categoryRepo repository.CategoryRepository) CategoryService {
	return &categoryService{
		categoryRepo: categoryRepo,
	}
}

func (s *categoryService) GetAllCategories() []models.ProductCategory {
	return s.categoryRepo.GetAllCategories()
}

func (s *categoryService) GetCategoryStats() ([]repository.CategoryStats, error) {
	return s.categoryRepo.GetCategoryStats()
}

func (s *categoryService) GetCategoryName(category models.ProductCategory) string {
	categoryNames := map[models.ProductCategory]string{
		models.CategoryTops:        "Tops",
		models.CategoryBottoms:     "Bottoms",
		models.CategoryDresses:     "Dresses",
		models.CategoryOuterwear:   "Outerwear",
		models.CategoryFootwear:    "Footwear",
		models.CategoryAccessories: "Accessories",
	}

	if name, exists := categoryNames[category]; exists {
		return name
	}

	return string(category)
}

func (s *categoryService) IsValidCategory(category models.ProductCategory) bool {
	validCategories := map[models.ProductCategory]bool{
		models.CategoryTops:        true,
		models.CategoryBottoms:     true,
		models.CategoryDresses:     true,
		models.CategoryOuterwear:   true,
		models.CategoryFootwear:    true,
		models.CategoryAccessories: true,
	}

	return validCategories[category]
}
