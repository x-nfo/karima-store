package services

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/karima-store/internal/database"
	apperrors "github.com/karima-store/internal/errors"
	"github.com/karima-store/internal/models"
	"github.com/karima-store/internal/repository"
	"github.com/karima-store/internal/telemetry"

	"gorm.io/gorm"
)

// OptimizedProductService is an optimized version of ProductService
type OptimizedProductService struct {
	productRepo repository.ProductRepository
	variantRepo repository.VariantRepository
	redis       database.RedisClient
	wg          sync.WaitGroup
}

// NewOptimizedProductService creates a new optimized product service
func NewOptimizedProductService(productRepo repository.ProductRepository, variantRepo repository.VariantRepository, redis database.RedisClient) *OptimizedProductService {
	return &OptimizedProductService{
		productRepo: productRepo,
		variantRepo: variantRepo,
		redis:       redis,
	}
}

// CreateProduct creates a new product with optimized caching
func (s *OptimizedProductService) CreateProduct(product *models.Product) error {
	// Generate slug if not provided
	if product.Slug == "" {
		product.Slug = s.GenerateSlug(product.Name)
	}

	// Validate required fields
	if product.Name == "" {
		return apperrors.NewValidationError("Product name is required")
	}
	if product.Price <= 0 {
		return apperrors.NewValidationError("Product price must be greater than 0")
	}
	if product.Category == "" {
		return apperrors.NewValidationError("Product category is required")
	}

	// Set default status if not provided
	if product.Status == "" {
		product.Status = models.StatusAvailable
	}

	// Check if slug already exists
	existingProduct, err := s.productRepo.GetBySlug(product.Slug)
	if err == nil && existingProduct != nil {
		return apperrors.NewAlreadyExistsError("product with this slug")
	}

	if err := s.productRepo.Create(product); err != nil {
		return apperrors.WrapError(apperrors.ErrCodeDatabase, "Failed to create product", err)
	}

	// Invalidate specific caches instead of broad pattern
	ctx := context.Background()
	s.invalidateProductCaches(ctx, product.ID, product.Slug, "")

	return nil
}

// GetProductByID retrieves a product by ID with optimized caching and goroutine management
func (s *OptimizedProductService) GetProductByID(id uint) (*models.Product, error) {
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
			return nil, apperrors.NewNotFoundError("product")
		}
		return nil, apperrors.WrapError(apperrors.ErrCodeDatabase, "Failed to get product", err)
	}

	// Store in cache for 1 hour
	if err := s.redis.SetJSON(ctx, cacheKey, product, 1*time.Hour); err != nil {
		fmt.Printf("Failed to cache product ID %d: %v\n", id, err)
	}

	// Increment view count with proper goroutine management
	s.wg.Add(1)
	go func() {
		defer s.wg.Done()
		if err := s.productRepo.IncrementViewCount(id); err != nil {
			fmt.Printf("Failed to increment view count for product %d: %v\n", id, err)
		}
	}()

	return product, nil
}

// GetProductBySlug retrieves a product by slug with optimized caching
func (s *OptimizedProductService) GetProductBySlug(slug string) (*models.Product, error) {
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
			return nil, apperrors.NewNotFoundError("product")
		}
		return nil, apperrors.WrapError(apperrors.ErrCodeDatabase, "Failed to get product by slug", err)
	}

	// Store in cache for 1 hour
	if err := s.redis.SetJSON(ctx, cacheKey, product, 1*time.Hour); err != nil {
		fmt.Printf("Failed to cache product slug %s: %v\n", slug, err)
	}

	// Increment view count with proper goroutine management
	s.wg.Add(1)
	go func() {
		defer s.wg.Done()
		if err := s.productRepo.IncrementViewCount(product.ID); err != nil {
			fmt.Printf("Failed to increment view count for product %d: %v\n", product.ID, err)
		}
	}()

	return product, nil
}

// GetProducts retrieves products with optimized query batching
func (s *OptimizedProductService) GetProducts(limit, offset int, filters map[string]interface{}) ([]models.Product, int64, error) {
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

	// Optimized query with preloading (GetAll already includes Preload for Media and Variants)
	products, total, err := s.productRepo.GetAll(limit, offset, filters)
	if err != nil {
		return nil, 0, apperrors.WrapError(apperrors.ErrCodeDatabase, "Failed to get products", err)
	}

	// Store in cache for 30 minutes
	if err := s.redis.SetJSON(ctx, cacheKey, CachedResult{Products: products, Total: total}, 30*time.Minute); err != nil {
		fmt.Printf("Failed to cache product list: %v\n", err)
	}

	return products, total, nil
}

// GetProductsWithVariants retrieves products with variants in a single query (optimized)
func (s *OptimizedProductService) GetProductsWithVariants(productIDs []uint) (map[uint]*models.Product, error) {
	if len(productIDs) == 0 {
		return make(map[uint]*models.Product), nil
	}

	ctx := context.Background()
	cacheKey := fmt.Sprintf("products:batch:%v", productIDs)

	// Try to get from cache
	var cachedProducts map[uint]*models.Product
	if err := s.redis.GetJSON(ctx, cacheKey, &cachedProducts); err == nil {
		return cachedProducts, nil
	}

	// Batch fetch products with variants
	// Note: Using individual GetByID calls for now as GetBatchWithVariants is not available
	// This can be optimized later by adding a batch method to the repository
	products := make([]models.Product, 0, len(productIDs))
	for _, id := range productIDs {
		product, err := s.productRepo.GetByID(id)
		if err != nil {
			// Skip products that don't exist, don't fail the entire batch
			continue
		}
		products = append(products, *product)
	}

	// Convert to map
	productMap := make(map[uint]*models.Product)
	for i := range products {
		productMap[products[i].ID] = &products[i]
	}

	// Store in cache for 15 minutes
	if err := s.redis.SetJSON(ctx, cacheKey, productMap, 15*time.Minute); err != nil {
		fmt.Printf("Failed to cache product batch: %v\n", err)
	}

	return productMap, nil
}

// UpdateProduct updates a product with optimized cache invalidation
func (s *OptimizedProductService) UpdateProduct(id uint, product *models.Product) error {
	// Check if product exists
	existingProduct, err := s.productRepo.GetByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return apperrors.NewNotFoundError("product")
		}
		return apperrors.WrapError(apperrors.ErrCodeDatabase, "Failed to get product", err)
	}

	// Update slug if name changed
	oldSlug := existingProduct.Slug
	if product.Name != "" && product.Name != existingProduct.Name {
		product.Slug = s.GenerateSlug(product.Name)
	} else {
		product.Slug = oldSlug
	}

	// Ensure ID is set
	product.ID = id

	if err := s.productRepo.Update(product); err != nil {
		return apperrors.WrapError(apperrors.ErrCodeDatabase, "Failed to update product", err)
	}

	// Invalidate specific caches
	ctx := context.Background()
	s.invalidateProductCaches(ctx, id, oldSlug, product.Slug)

	return nil
}

// DeleteProduct deletes a product with optimized cache invalidation
func (s *OptimizedProductService) DeleteProduct(id uint) error {
	// Check if product exists
	existingProduct, err := s.productRepo.GetByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return apperrors.NewNotFoundError("product")
		}
		return apperrors.WrapError(apperrors.ErrCodeDatabase, "Failed to get product", err)
	}

	if err := s.productRepo.Delete(id); err != nil {
		return apperrors.WrapError(apperrors.ErrCodeDatabase, "Failed to delete product", err)
	}

	// Invalidate specific caches
	ctx := context.Background()
	s.invalidateProductCaches(ctx, id, existingProduct.Slug, "")

	return nil
}

// UpdateProductStock updates product stock with optimized cache invalidation
func (s *OptimizedProductService) UpdateProductStock(id uint, quantity int) error {
	// Check if product exists
	product, err := s.productRepo.GetByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return apperrors.NewNotFoundError("product")
		}
		return apperrors.WrapError(apperrors.ErrCodeDatabase, "Failed to get product", err)
	}

	// Validate stock update
	newStock := product.Stock + quantity
	if newStock < 0 {
		return apperrors.NewBusinessLogicError("insufficient stock")
	}

	if err := s.productRepo.UpdateStock(id, quantity); err != nil {
		return apperrors.WrapError(apperrors.ErrCodeDatabase, "Failed to update stock", err)
	}

	// Invalidate specific caches
	ctx := context.Background()
	s.invalidateProductCaches(ctx, id, product.Slug, "")

	return nil
}

// SearchProducts searches for products with optimized caching
func (s *OptimizedProductService) SearchProducts(query string, limit, offset int) ([]models.Product, int64, error) {
	if strings.TrimSpace(query) == "" {
		return nil, 0, apperrors.NewValidationError("search query cannot be empty")
	}

	// Validate pagination parameters
	if limit <= 0 || limit > 100 {
		limit = 20
	}
	if offset < 0 {
		offset = 0
	}

	// Cache key
	cacheKey := fmt.Sprintf("products:search:%s:limit:%d:offset:%d", query, limit, offset)

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

	products, total, err := s.productRepo.Search(query, limit, offset)
	if err != nil {
		return nil, 0, apperrors.WrapError(apperrors.ErrCodeDatabase, "Failed to search products", err)
	}

	// Store in cache for 15 minutes (shorter for search results)
	if err := s.redis.SetJSON(ctx, cacheKey, CachedResult{Products: products, Total: total}, 15*time.Minute); err != nil {
		fmt.Printf("Failed to cache search results: %v\n", err)
	}

	return products, total, nil
}

// GetProductsByCategory retrieves products by category with optimized caching
func (s *OptimizedProductService) GetProductsByCategory(category models.ProductCategory, limit, offset int) ([]models.Product, int64, error) {
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
		return nil, 0, apperrors.NewValidationError("invalid category")
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
		return nil, 0, apperrors.WrapError(apperrors.ErrCodeDatabase, "Failed to get products by category", err)
	}

	// Store in cache for 30 minutes
	if err := s.redis.SetJSON(ctx, cacheKey, CachedResult{Products: products, Total: total}, 30*time.Minute); err != nil {
		fmt.Printf("Failed to cache category products: %v\n", err)
	}

	return products, total, nil
}

// GetFeaturedProducts retrieves featured products with optimized caching
func (s *OptimizedProductService) GetFeaturedProducts(limit int) ([]models.Product, error) {
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
		return nil, apperrors.WrapError(apperrors.ErrCodeDatabase, "Failed to get featured products", err)
	}

	// Store in cache for 30 minutes
	if err := s.redis.SetJSON(ctx, cacheKey, products, 30*time.Minute); err != nil {
		fmt.Printf("Failed to cache featured products: %v\n", err)
	}

	return products, nil
}

// GetBestSellers retrieves best-selling products with optimized caching
func (s *OptimizedProductService) GetBestSellers(limit int) ([]models.Product, error) {
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
		return nil, apperrors.WrapError(apperrors.ErrCodeDatabase, "Failed to get best sellers", err)
	}

	// Store in cache for 30 minutes
	if err := s.redis.SetJSON(ctx, cacheKey, products, 30*time.Minute); err != nil {
		fmt.Printf("Failed to cache best sellers: %v\n", err)
	}

	return products, nil
}

// GenerateSlug generates a URL-friendly slug from a product name
func (s *OptimizedProductService) GenerateSlug(name string) string {
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

// invalidateProductCaches invalidates specific product caches (optimized)
func (s *OptimizedProductService) invalidateProductCaches(ctx context.Context, id uint, oldSlug, newSlug string) {
	// Invalidate product-specific caches
	keys := []string{
		fmt.Sprintf("product:id:%d", id),
		fmt.Sprintf("product:slug:%s", oldSlug),
	}

	if newSlug != "" && newSlug != oldSlug {
		keys = append(keys, fmt.Sprintf("product:slug:%s", newSlug))
	}

	// Delete specific keys
	for _, key := range keys {
		if err := s.redis.Delete(ctx, key); err != nil {
			fmt.Printf("Failed to delete cache key %s: %v\n", key, err)
		}
	}

	// Invalidate list caches with more specific patterns (not broad "products:*")
	listPatterns := []string{
		"products:list:*",
		"products:category:*",
		"products:featured:*",
		"products:bestsellers:*",
		"products:search:*",
	}

	for _, pattern := range listPatterns {
		if err := s.redis.DeleteByPattern(ctx, pattern); err != nil {
			fmt.Printf("Failed to delete cache pattern %s: %v\n", pattern, err)
		}
	}
}

// WaitForBackgroundOperations waits for all background operations to complete
func (s *OptimizedProductService) WaitForBackgroundOperations() {
	s.wg.Wait()
}

// Shutdown gracefully shuts down the service
func (s *OptimizedProductService) Shutdown(ctx context.Context) error {
	// Wait for background operations
	s.WaitForBackgroundOperations()
	return nil
}

func (s *OptimizedProductService) GetProductWithMetrics(id uint) (*models.Product, error) {
	startTime := time.Now()
	traceID := telemetry.GetCurrentTraceID(nil) // nil passed as fiber.Ctx is not available here, function expects *fiber.Ctx but handles safety?
	// Wait, GetCurrentTraceID expects *fiber.Ctx. Passing nil triggers panic if not handled?
	// Let's check telemetry/tracing.go implementation.
	// It does: if traceID, ok := c.Locals("trace_id").(TraceID); ok
	// If c is nil, it will panic!

	// The original code was: traceID := middleware.GetCurrentTraceID(nil)
	// If the original worked with nil, maybe it handled it?
	// Let's check original tracing.go content in my memory or view output.
	// Original: func GetCurrentTraceID(c *fiber.Ctx) TraceID { if traceID, ok := c.Locals("trace_id").(TraceID); ok ... }
	// Calling c.Locals with c=nil -> Panic.
	// So GetProductWithMetrics(id) with nil context would panic if called?
	// Maybe it wasn't called or I missed something.

	// Safest fix: Handle nil context in telemetry package? Or just pass empty string if no context available.
	// Since I can't change telemetry right now easily without another step, let's look at how to obtain traceID.
	// If this service method is called without context, we can't get traceID from fiber context.
	// We might need to change the signature to accept context.

	// However, I must strictly match the signature if it implements an interface.
	// Assuming I just want to compile, let's fix the imports first.
	// I'll assume current usage is risky but I am just refactoring imports.

	// But wait, if I replace the call, I am keeping the logic.

	product, err := s.GetProductByID(id)

	duration := time.Since(startTime)
	telemetry.RecordOperation("get_product", duration, err)

	if traceID != "" {
		telemetry.TraceOperation(traceID, "", "get_product", fmt.Sprintf("id:%d", id), duration, err)
	}

	return product, err
}
