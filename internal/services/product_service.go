package services

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"gorm.io/gorm"
	"github.com/karima-store/internal/models"
	"github.com/karima-store/internal/repository"
)

type ProductService interface {
	CreateProduct(product *models.Product) error
	GetProductByID(id uint) (*models.Product, error)
	GetProductBySlug(slug string) (*models.Product, error)
	GetProducts(limit, offset int, filters map[string]interface{}) ([]models.Product, int64, error)
	UpdateProduct(id uint, product *models.Product) error
	DeleteProduct(id uint) error
	UpdateProductStock(id uint, quantity int) error
	SearchProducts(query string, limit, offset int) ([]models.Product, int64, error)
	GetProductsByCategory(category models.ProductCategory, limit, offset int) ([]models.Product, int64, error)
	GetFeaturedProducts(limit int) ([]models.Product, error)
	GetBestSellers(limit int) ([]models.Product, error)
	GenerateSlug(name string) string
}

type productService struct {
	productRepo repository.ProductRepository
	variantRepo repository.VariantRepository
}

func NewProductService(productRepo repository.ProductRepository, variantRepo repository.VariantRepository) ProductService {
	return &productService{
		productRepo: productRepo,
		variantRepo: variantRepo,
	}
}

func (s *productService) CreateProduct(product *models.Product) error {
	// Generate slug if not provided
	if product.Slug == "" {
		product.Slug = s.GenerateSlug(product.Name)
	}

	// Validate required fields
	if product.Name == "" {
		return errors.New("product name is required")
	}
	if product.Price <= 0 {
		return errors.New("product price must be greater than 0")
	}
	if product.Category == "" {
		return errors.New("product category is required")
	}

	// Set default status if not provided
	if product.Status == "" {
		product.Status = models.StatusAvailable
	}

	// Check if slug already exists
	existingProduct, err := s.productRepo.GetBySlug(product.Slug)
	if err == nil && existingProduct != nil {
		return errors.New("product with this slug already exists")
	}

	return s.productRepo.Create(product)
}

func (s *productService) GetProductByID(id uint) (*models.Product, error) {
	product, err := s.productRepo.GetByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("product not found")
		}
		return nil, err
	}

	// Increment view count
	go s.productRepo.IncrementViewCount(id)

	return product, nil
}

func (s *productService) GetProductBySlug(slug string) (*models.Product, error) {
	product, err := s.productRepo.GetBySlug(slug)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("product not found")
		}
		return nil, err
	}

	// Increment view count
	go s.productRepo.IncrementViewCount(product.ID)

	return product, nil
}

func (s *productService) GetProducts(limit, offset int, filters map[string]interface{}) ([]models.Product, int64, error) {
	// Validate pagination parameters
	if limit <= 0 || limit > 100 {
		limit = 20 // default limit
	}
	if offset < 0 {
		offset = 0
	}

	return s.productRepo.GetAll(limit, offset, filters)
}

func (s *productService) UpdateProduct(id uint, product *models.Product) error {
	// Check if product exists
	existingProduct, err := s.productRepo.GetByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("product not found")
		}
		return err
	}

	// Update slug if name changed
	if product.Name != "" && product.Name != existingProduct.Name {
		product.Slug = s.GenerateSlug(product.Name)
	}

	// Ensure ID is set
	product.ID = id

	return s.productRepo.Update(product)
}

func (s *productService) DeleteProduct(id uint) error {
	// Check if product exists
	_, err := s.productRepo.GetByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("product not found")
		}
		return err
	}

	return s.productRepo.Delete(id)
}

func (s *productService) UpdateProductStock(id uint, quantity int) error {
	// Check if product exists
	product, err := s.productRepo.GetByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("product not found")
		}
		return err
	}

	// Validate stock update
	newStock := product.Stock + quantity
	if newStock < 0 {
		return errors.New("insufficient stock")
	}

	return s.productRepo.UpdateStock(id, quantity)
}

func (s *productService) SearchProducts(query string, limit, offset int) ([]models.Product, int64, error) {
	if strings.TrimSpace(query) == "" {
		return nil, 0, errors.New("search query cannot be empty")
	}

	// Validate pagination parameters
	if limit <= 0 || limit > 100 {
		limit = 20
	}
	if offset < 0 {
		offset = 0
	}

	return s.productRepo.Search(query, limit, offset)
}

func (s *productService) GetProductsByCategory(category models.ProductCategory, limit, offset int) ([]models.Product, int64, error) {
	// Validate category
	validCategories := map[models.ProductCategory]bool{
		models.CategoryTops:        true,
		models.CategoryBottoms:     true,
		models.CategoryDresses:     true,
		models.CategoryOuterwear:   true,
		models.CategoryFootwear:    true,
		models.CategoryAccessories: true,
	}

	if !validCategories[category] {
		return nil, 0, errors.New("invalid category")
	}

	// Validate pagination parameters
	if limit <= 0 || limit > 100 {
		limit = 20
	}
	if offset < 0 {
		offset = 0
	}

	return s.productRepo.GetByCategory(category, limit, offset)
}

func (s *productService) GetFeaturedProducts(limit int) ([]models.Product, error) {
	if limit <= 0 || limit > 50 {
		limit = 10
	}

	return s.productRepo.GetFeatured(limit)
}

func (s *productService) GetBestSellers(limit int) ([]models.Product, error) {
	if limit <= 0 || limit > 50 {
		limit = 10
	}

	return s.productRepo.GetBestSellers(limit)
}

func (s *productService) GenerateSlug(name string) string {
	// Convert to lowercase
	slug := strings.ToLower(name)

	// Replace spaces with hyphens
	slug = strings.ReplaceAll(slug, " ", "-")

	// Remove special characters
	slug = strings.Map(func(r rune) rune {
		if (r >= 'a' && r <= 'z') || (r >= '0' && r <= '9') || r == '-' {
			return r
		}
		return -1
	}, slug)

	// Remove multiple consecutive hyphens
	slug = strings.ReplaceAll(slug, "--", "-")

	// Trim hyphens from start and end
	slug = strings.Trim(slug, "-")

	// Add timestamp if slug is empty
	if slug == "" {
		slug = fmt.Sprintf("product-%d", time.Now().Unix())
	}

	return slug
}
