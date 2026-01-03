package services

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"regexp"
	"strings"
	"time"

	"github.com/karima-store/internal/database"
	"github.com/karima-store/internal/models"
	"github.com/karima-store/internal/repository"
	"gorm.io/gorm"
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
	redis       database.RedisClient
}

func NewProductService(productRepo repository.ProductRepository, variantRepo repository.VariantRepository, redis database.RedisClient) ProductService {
	return &productService{
		productRepo: productRepo,
		variantRepo: variantRepo,
		redis:       redis,
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

	if err := s.productRepo.Create(product); err != nil {
		return err
	}

	// Invalidate list caches
	ctx := context.Background()
	_ = s.redis.DeleteByPattern(ctx, "products:*")

	return nil
}

func (s *productService) GetProductByID(id uint) (*models.Product, error) {
	ctx := context.Background()
	cacheKey := fmt.Sprintf("product:id:%d", id)

	// Try to get from cache
	var cachedProduct models.Product
	if err := s.redis.GetJSON(ctx, cacheKey, &cachedProduct); err == nil {
		return &cachedProduct, nil
	}

	product, err := s.productRepo.GetByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("product not found")
		}
		return nil, err
	}

	// Store in cache for 1 hour
	if err := s.redis.SetJSON(ctx, cacheKey, product, 1*time.Hour); err != nil {
		fmt.Printf("Failed to cache product ID %d: %v\n", id, err)
	}

	// Increment view count
	go s.productRepo.IncrementViewCount(id)

	return product, nil
}

func (s *productService) GetProductBySlug(slug string) (*models.Product, error) {
	ctx := context.Background()
	cacheKey := fmt.Sprintf("product:slug:%s", slug)

	// Try to get from cache
	var cachedProduct models.Product
	if err := s.redis.GetJSON(ctx, cacheKey, &cachedProduct); err == nil {
		return &cachedProduct, nil
	}

	product, err := s.productRepo.GetBySlug(slug)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("product not found")
		}
		return nil, err
	}

	// Store in cache for 1 hour
	if err := s.redis.SetJSON(ctx, cacheKey, product, 1*time.Hour); err != nil {
		fmt.Printf("Failed to cache product slug %s: %v\n", slug, err)
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

	// Cache key
	filterBytes, _ := json.Marshal(filters)
	cacheKey := fmt.Sprintf("products:list:limit:%d:offset:%d:filters:%s", limit, offset, string(filterBytes))

	// Try to get from cache
	type CachedResult struct {
		Products []models.Product
		Total    int64
	}
	var cachedResult CachedResult
	ctx := context.Background()
	if err := s.redis.GetJSON(ctx, cacheKey, &cachedResult); err == nil {
		return cachedResult.Products, cachedResult.Total, nil
	}

	products, total, err := s.productRepo.GetAll(limit, offset, filters)
	if err != nil {
		return nil, 0, err
	}

	// Store in cache for 30 minutes
	if err := s.redis.SetJSON(ctx, cacheKey, CachedResult{Products: products, Total: total}, 30*time.Minute); err != nil {
		fmt.Printf("Failed to cache product list: %v\n", err)
	}

	return products, total, nil
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

	if err := s.productRepo.Update(product); err != nil {
		return err
	}

	// Invalidate cache
	ctx := context.Background()
	_ = s.redis.Delete(ctx, fmt.Sprintf("product:id:%d", id))
	_ = s.redis.Delete(ctx, fmt.Sprintf("product:slug:%s", existingProduct.Slug)) // Old slug
	if product.Slug != "" && product.Slug != existingProduct.Slug {
		_ = s.redis.Delete(ctx, fmt.Sprintf("product:slug:%s", product.Slug)) // New slug (just in case)
	}

	// Invalidate list caches
	_ = s.redis.DeleteByPattern(ctx, "products:*")

	return nil
}

func (s *productService) DeleteProduct(id uint) error {
	// Check if product exists
	existingProduct, err := s.productRepo.GetByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("product not found")
		}
		return err
	}

	if err := s.productRepo.Delete(id); err != nil {
		return err
	}

	// Invalidate cache
	ctx := context.Background()
	_ = s.redis.Delete(ctx, fmt.Sprintf("product:id:%d", id))
	_ = s.redis.Delete(ctx, fmt.Sprintf("product:slug:%s", existingProduct.Slug))

	// Invalidate list caches
	_ = s.redis.DeleteByPattern(ctx, "products:*")

	return nil
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

	if err := s.productRepo.UpdateStock(id, quantity); err != nil {
		return err
	}

	// Invalidate cache
	ctx := context.Background()
	_ = s.redis.Delete(ctx, fmt.Sprintf("product:id:%d", id))
	_ = s.redis.Delete(ctx, fmt.Sprintf("product:slug:%s", product.Slug))

	// Invalidate list caches
	_ = s.redis.DeleteByPattern(ctx, "products:*")

	return nil
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

	// Cache key
	cacheKey := fmt.Sprintf("products:category:%s:limit:%d:offset:%d", category, limit, offset)

	// Try to get from cache
	type CachedResult struct {
		Products []models.Product
		Total    int64
	}
	var cachedResult CachedResult
	ctx := context.Background()
	if err := s.redis.GetJSON(ctx, cacheKey, &cachedResult); err == nil {
		return cachedResult.Products, cachedResult.Total, nil
	}

	products, total, err := s.productRepo.GetByCategory(category, limit, offset)
	if err != nil {
		return nil, 0, err
	}

	// Store in cache for 30 minutes
	if err := s.redis.SetJSON(ctx, cacheKey, CachedResult{Products: products, Total: total}, 30*time.Minute); err != nil {
		fmt.Printf("Failed to cache category products: %v\n", err)
	}

	return products, total, nil
}

func (s *productService) GetFeaturedProducts(limit int) ([]models.Product, error) {
	if limit <= 0 || limit > 50 {
		limit = 10
	}

	// Cache key
	cacheKey := fmt.Sprintf("products:featured:limit:%d", limit)

	// Try to get from cache
	var cachedProducts []models.Product
	ctx := context.Background()
	if err := s.redis.GetJSON(ctx, cacheKey, &cachedProducts); err == nil {
		return cachedProducts, nil
	}

	products, err := s.productRepo.GetFeatured(limit)
	if err != nil {
		return nil, err
	}

	// Store in cache for 30 minutes
	if err := s.redis.SetJSON(ctx, cacheKey, products, 30*time.Minute); err != nil {
		fmt.Printf("Failed to cache featured products: %v\n", err)
	}

	return products, nil
}

func (s *productService) GetBestSellers(limit int) ([]models.Product, error) {
	if limit <= 0 || limit > 50 {
		limit = 10
	}

	// Cache key
	cacheKey := fmt.Sprintf("products:bestsellers:limit:%d", limit)

	// Try to get from cache
	var cachedProducts []models.Product
	ctx := context.Background()
	if err := s.redis.GetJSON(ctx, cacheKey, &cachedProducts); err == nil {
		return cachedProducts, nil
	}

	products, err := s.productRepo.GetBestSellers(limit)
	if err != nil {
		return nil, err
	}

	// Store in cache for 30 minutes
	if err := s.redis.SetJSON(ctx, cacheKey, products, 30*time.Minute); err != nil {
		fmt.Printf("Failed to cache best sellers: %v\n", err)
	}

	return products, nil
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
	regex := regexp.MustCompile(`-+`)
	slug = regex.ReplaceAllString(slug, "-")

	// Trim hyphens from start and end
	slug = strings.Trim(slug, "-")

	// Add timestamp if slug is empty
	if slug == "" {
		slug = fmt.Sprintf("product-%d", time.Now().Unix())
	}

	return slug
}
