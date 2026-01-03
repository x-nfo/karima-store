package repository

import (
	"fmt"
	"testing"

	"github.com/karima-store/internal/models"
	"github.com/karima-store/internal/test_setup"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"
)

func setupProductTest(t *testing.T) (*gorm.DB, func()) {
	db, cleanup := test_setup.SetupTestDB(t)

	// Clean up any existing data
	db.Exec("DELETE FROM products")

	return db, cleanup
}

func TestProductRepository_Create(t *testing.T) {
	db, cleanup := setupProductTest(t)
	defer cleanup()

	repo := NewProductRepository(db)

	product := &models.Product{
		Name:        "Test Product",
		Description: "Test Description",
		Price:       100.00,
		Category:    models.CategoryTops,
		Stock:       10,
		Status:      models.StatusAvailable,
		Slug:        "test-product", SKU: "test-product-sku",
	}

	err := repo.Create(product)
	require.NoError(t, err)
	assert.NotZero(t, product.ID)
	assert.Equal(t, "Test Product", product.Name)
}

func TestProductRepository_GetByID(t *testing.T) {
	db, cleanup := setupProductTest(t)
	defer cleanup()

	repo := NewProductRepository(db)

	// Create a product
	product := &models.Product{
		Name:        "Test Product",
		Description: "Test Description",
		Price:       100.00,
		Category:    models.CategoryTops,
		Stock:       10,
		Status:      models.StatusAvailable,
		Slug:        "test-product", SKU: "test-product-sku",
	}
	err := repo.Create(product)
	require.NoError(t, err)

	// Get the product
	fetchedProduct, err := repo.GetByID(product.ID)
	require.NoError(t, err)
	assert.Equal(t, product.ID, fetchedProduct.ID)
	assert.Equal(t, product.Name, fetchedProduct.Name)
}

func TestProductRepository_GetByID_NotFound(t *testing.T) {
	db, cleanup := setupProductTest(t)
	defer cleanup()

	repo := NewProductRepository(db)

	// Try to get a non-existent product
	_, err := repo.GetByID(99999)
	assert.Error(t, err)
	assert.Equal(t, gorm.ErrRecordNotFound, err)
}

func TestProductRepository_GetBySlug(t *testing.T) {
	db, cleanup := setupProductTest(t)
	defer cleanup()

	repo := NewProductRepository(db)

	// Create a product
	product := &models.Product{
		Name:        "Test Product",
		Description: "Test Description",
		Price:       100.00,
		Category:    models.CategoryTops,
		Stock:       10,
		Status:      models.StatusAvailable,
		Slug:        "test-product", SKU: "test-product-sku",
	}
	err := repo.Create(product)
	require.NoError(t, err)

	// Get the product by slug
	fetchedProduct, err := repo.GetBySlug("test-product")
	require.NoError(t, err)
	assert.Equal(t, product.ID, fetchedProduct.ID)
	assert.Equal(t, product.Slug, fetchedProduct.Slug)
}

func TestProductRepository_GetBySlug_NotFound(t *testing.T) {
	db, cleanup := setupProductTest(t)
	defer cleanup()

	repo := NewProductRepository(db)

	// Try to get a non-existent product
	_, err := repo.GetBySlug("non-existent-slug")
	assert.Error(t, err)
	assert.Equal(t, gorm.ErrRecordNotFound, err)
}

func TestProductRepository_GetAll(t *testing.T) {
	db, cleanup := setupProductTest(t)
	defer cleanup()

	repo := NewProductRepository(db)

	// Create multiple products
	for i := 1; i <= 5; i++ {
		product := &models.Product{
			Name:        fmt.Sprintf("Test Product %d", i),
			Description: "Test Description",
			Price:       float64(i * 100),
			Category:    models.CategoryTops,
			Stock:       10,
			Status:      models.StatusAvailable,
			Slug:        fmt.Sprintf("test-product-%d", i),
			SKU:         fmt.Sprintf("sku-%d", i),
		}
		err := repo.Create(product)
		require.NoError(t, err)
	}

	// Get all products with pagination
	products, total, err := repo.GetAll(10, 0, nil)
	require.NoError(t, err)
	assert.Len(t, products, 5)
	assert.Equal(t, int64(5), total)
}

func TestProductRepository_GetAll_WithFilters(t *testing.T) {
	db, cleanup := setupProductTest(t)
	defer cleanup()

	repo := NewProductRepository(db)

	// Create products with different categories
	product1 := &models.Product{
		Name:     "Top Product",
		Price:    100.00,
		Category: models.CategoryTops,
		Stock:    10,
		Status:   models.StatusAvailable,
		Slug:     "top-product", SKU: "top-product-sku",
	}
	product2 := &models.Product{
		Name:     "Bottom Product",
		Price:    200.00,
		Category: models.CategoryBottoms,
		Stock:    10,
		Status:   models.StatusAvailable,
		Slug:     "bottom-product", SKU: "bottom-product-sku",
	}

	err := repo.Create(product1)
	require.NoError(t, err)
	err = repo.Create(product2)
	require.NoError(t, err)

	// Get products by category filter
	filters := map[string]interface{}{"category": models.CategoryTops}
	products, total, err := repo.GetAll(10, 0, filters)
	require.NoError(t, err)
	assert.Len(t, products, 1)
	assert.Equal(t, int64(1), total)
	assert.Equal(t, models.CategoryTops, products[0].Category)
}

func TestProductRepository_Update(t *testing.T) {
	db, cleanup := setupProductTest(t)
	defer cleanup()

	repo := NewProductRepository(db)

	// Create a product
	product := &models.Product{
		Name:        "Test Product",
		Description: "Test Description",
		Price:       100.00,
		Category:    models.CategoryTops,
		Stock:       10,
		Status:      models.StatusAvailable,
		Slug:        "test-product", SKU: "test-product-sku",
	}
	err := repo.Create(product)
	require.NoError(t, err)

	// Update the product
	product.Name = "Updated Product"
	product.Price = 150.00
	err = repo.Update(product)
	require.NoError(t, err)

	// Verify the update
	fetchedProduct, err := repo.GetByID(product.ID)
	require.NoError(t, err)
	assert.Equal(t, "Updated Product", fetchedProduct.Name)
	assert.Equal(t, 150.00, fetchedProduct.Price)
}

func TestProductRepository_Delete(t *testing.T) {
	db, cleanup := setupProductTest(t)
	defer cleanup()

	repo := NewProductRepository(db)

	// Create a product
	product := &models.Product{
		Name:        "Test Product",
		Description: "Test Description",
		Price:       100.00,
		Category:    models.CategoryTops,
		Stock:       10,
		Status:      models.StatusAvailable,
		Slug:        "test-product", SKU: "test-product-sku",
	}
	err := repo.Create(product)
	require.NoError(t, err)

	// Delete the product
	err = repo.Delete(product.ID)
	require.NoError(t, err)

	// Verify the deletion
	_, err = repo.GetByID(product.ID)
	assert.Error(t, err)
	assert.Equal(t, gorm.ErrRecordNotFound, err)
}

func TestProductRepository_UpdateStock(t *testing.T) {
	db, cleanup := setupProductTest(t)
	defer cleanup()

	repo := NewProductRepository(db)

	// Create a product
	product := &models.Product{
		Name:        "Test Product",
		Description: "Test Description",
		Price:       100.00,
		Category:    models.CategoryTops,
		Stock:       10,
		Status:      models.StatusAvailable,
		Slug:        "test-product", SKU: "test-product-sku",
	}
	err := repo.Create(product)
	require.NoError(t, err)

	// Update stock
	err = repo.UpdateStock(product.ID, -5)
	require.NoError(t, err)

	// Verify the update
	fetchedProduct, err := repo.GetByID(product.ID)
	require.NoError(t, err)
	assert.Equal(t, 5, fetchedProduct.Stock)
}

func TestProductRepository_UpdateStock_Insufficient(t *testing.T) {
	db, cleanup := setupProductTest(t)
	defer cleanup()

	repo := NewProductRepository(db)

	// Create a product with low stock
	product := &models.Product{
		Name:        "Test Product",
		Description: "Test Description",
		Price:       100.00,
		Category:    models.CategoryTops,
		Stock:       5,
		Status:      models.StatusAvailable,
		Slug:        "test-product", SKU: "test-product-sku",
	}
	err := repo.Create(product)
	require.NoError(t, err)

	// Try to update stock with insufficient quantity
	err = repo.UpdateStock(product.ID, -10)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "insufficient stock")
}

func TestProductRepository_IncrementViewCount(t *testing.T) {
	db, cleanup := setupProductTest(t)
	defer cleanup()

	repo := NewProductRepository(db)

	// Create a product
	product := &models.Product{
		Name:        "Test Product",
		Description: "Test Description",
		Price:       100.00,
		Category:    models.CategoryTops,
		Stock:       10,
		Status:      models.StatusAvailable,
		Slug:        "test-product", SKU: "test-product-sku",
		ViewCount: 0,
	}
	err := repo.Create(product)
	require.NoError(t, err)

	// Increment view count
	err = repo.IncrementViewCount(product.ID)
	require.NoError(t, err)

	// Verify the increment
	fetchedProduct, err := repo.GetByID(product.ID)
	require.NoError(t, err)
	assert.Equal(t, 1, fetchedProduct.ViewCount)
}

func TestProductRepository_Search(t *testing.T) {
	db, cleanup := setupProductTest(t)
	defer cleanup()

	repo := NewProductRepository(db)

	// Create products
	product1 := &models.Product{
		Name:        "Red T-Shirt",
		Description: "A comfortable red t-shirt",
		Price:       100.00,
		Category:    models.CategoryTops,
		Stock:       10,
		Status:      models.StatusAvailable,
		Slug:        "red-t-shirt", SKU: "red-t-shirt-sku",
	}
	product2 := &models.Product{
		Name:        "Blue Jeans",
		Description: "Stylish blue jeans",
		Price:       200.00,
		Category:    models.CategoryBottoms,
		Stock:       10,
		Status:      models.StatusAvailable,
		Slug:        "blue-jeans", SKU: "blue-jeans-sku",
	}
	product3 := &models.Product{
		Name:        "Red Dress",
		Description: "Elegant red dress",
		Price:       300.00,
		Category:    models.CategoryDresses,
		Stock:       10,
		Status:      models.StatusAvailable,
		Slug:        "red-dress", SKU: "red-dress-sku",
	}

	err := repo.Create(product1)
	require.NoError(t, err)
	err = repo.Create(product2)
	require.NoError(t, err)
	err = repo.Create(product3)
	require.NoError(t, err)

	// Search for "red"
	products, total, err := repo.Search("red", 10, 0)
	require.NoError(t, err)
	assert.Len(t, products, 2)
	assert.Equal(t, int64(2), total)
}

func TestProductRepository_GetByCategory(t *testing.T) {
	db, cleanup := setupProductTest(t)
	defer cleanup()

	repo := NewProductRepository(db)

	// Create products in different categories
	product1 := &models.Product{
		Name:     "T-Shirt 1",
		Price:    100.00,
		Category: models.CategoryTops,
		Stock:    10,
		Status:   models.StatusAvailable,
		Slug:     "t-shirt-1", SKU: "t-shirt-1-sku",
	}
	product2 := &models.Product{
		Name:     "T-Shirt 2",
		Price:    150.00,
		Category: models.CategoryTops,
		Stock:    10,
		Status:   models.StatusAvailable,
		Slug:     "t-shirt-2", SKU: "t-shirt-2-sku",
	}
	product3 := &models.Product{
		Name:     "Jeans",
		Price:    200.00,
		Category: models.CategoryBottoms,
		Stock:    10,
		Status:   models.StatusAvailable,
		Slug:     "jeans", SKU: "jeans-sku",
	}

	err := repo.Create(product1)
	require.NoError(t, err)
	err = repo.Create(product2)
	require.NoError(t, err)
	err = repo.Create(product3)
	require.NoError(t, err)

	// Get products by category
	products, total, err := repo.GetByCategory(models.CategoryTops, 10, 0)
	require.NoError(t, err)
	assert.Len(t, products, 2)
	assert.Equal(t, int64(2), total)
}

func TestProductRepository_GetFeatured(t *testing.T) {
	db, cleanup := setupProductTest(t)
	defer cleanup()

	repo := NewProductRepository(db)

	// Create featured and non-featured products
	product1 := &models.Product{
		Name:       "Featured Product 1",
		Price:      100.00,
		Category:   models.CategoryTops,
		Stock:      10,
		Status:     models.StatusAvailable,
		IsFeatured: true,
		Slug:       "featured-product-1", SKU: "featured-product-1-sku",
	}
	product2 := &models.Product{
		Name:       "Featured Product 2",
		Price:      150.00,
		Category:   models.CategoryTops,
		Stock:      10,
		Status:     models.StatusAvailable,
		IsFeatured: true,
		Slug:       "featured-product-2", SKU: "featured-product-2-sku",
	}
	product3 := &models.Product{
		Name:       "Regular Product",
		Price:      200.00,
		Category:   models.CategoryBottoms,
		Stock:      10,
		Status:     models.StatusAvailable,
		IsFeatured: false,
		Slug:       "regular-product", SKU: "regular-product-sku",
	}

	err := repo.Create(product1)
	require.NoError(t, err)
	err = repo.Create(product2)
	require.NoError(t, err)
	err = repo.Create(product3)
	require.NoError(t, err)

	// Get featured products
	products, err := repo.GetFeatured(10)
	require.NoError(t, err)
	assert.Len(t, products, 2)
}

func TestProductRepository_GetBestSellers(t *testing.T) {
	db, cleanup := setupProductTest(t)
	defer cleanup()

	repo := NewProductRepository(db)

	// Create products with different sales counts
	product1 := &models.Product{
		Name:      "Best Seller 1",
		Price:     100.00,
		Category:  models.CategoryTops,
		Stock:     10,
		Status:    models.StatusAvailable,
		SoldCount: 100,
		Slug:      "best-seller-1", SKU: "best-seller-1-sku",
	}
	product2 := &models.Product{
		Name:      "Best Seller 2",
		Price:     150.00,
		Category:  models.CategoryTops,
		Stock:     10,
		Status:    models.StatusAvailable,
		SoldCount: 50,
		Slug:      "best-seller-2", SKU: "best-seller-2-sku",
	}
	product3 := &models.Product{
		Name:      "Regular Product",
		Price:     200.00,
		Category:  models.CategoryBottoms,
		Stock:     10,
		Status:    models.StatusAvailable,
		SoldCount: 10,
		Slug:      "regular-product", SKU: "regular-product-sku",
	}

	err := repo.Create(product1)
	require.NoError(t, err)
	err = repo.Create(product2)
	require.NoError(t, err)
	err = repo.Create(product3)
	require.NoError(t, err)

	// Get best sellers
	products, err := repo.GetBestSellers(10)
	require.NoError(t, err)
	assert.Len(t, products, 3)
	assert.Equal(t, "Best Seller 1", products[0].Name)
}

// TODO: Implement GetAllWithPreload method in repository
// func TestProductRepository_GetAllWithPreload(t *testing.T) {
// 	db, cleanup := setupProductTest(t)
// 	defer cleanup()
//
// 	repo := NewProductRepository(db)
//
// 	// Create a product with variants
// 	product := &models.Product{
// 		Name:        "Test Product",
// 		Description: "Test Description",
// 		Price:       100.00,
// 		Category:    models.CategoryTops,
// 		Stock:       10,
// 		Status:      models.StatusAvailable,
// 		Slug: "test-product", SKU: "test-product-sku",
// 	}
// 	err := repo.Create(product)
// 	require.NoError(t, err)
//
// 	// Get products with preloaded variants
// 	products, total, err := repo.GetAllWithPreload(10, 0, nil)
// 	require.NoError(t, err)
// 	assert.Len(t, products, 1)
// 	assert.Equal(t, int64(1), total)
// }

// TODO: Implement GetBatchWithVariants method in repository
// func TestProductRepository_GetBatchWithVariants(t *testing.T) {
// 	db, cleanup := setupProductTest(t)
// 	defer cleanup()
//
// 	repo := NewProductRepository(db)
//
// 	// Create multiple products
// 	var productIDs []uint
// 	for i := 1; i <= 3; i++ {
// 		product := &models.Product{
// 			Name:        fmt.Sprintf("Test Product %d", i),
// 			Price:       float64(i * 100),
// 			Category:    models.CategoryTops,
// 			Stock:       10,
// 			Status:      models.StatusAvailable,
// 			Slug:        fmt.Sprintf("test-product-%d", i),
// 		}
// 		err := repo.Create(product)
// 		require.NoError(t, err)
// 		productIDs = append(productIDs, product.ID)
// 	}
//
// 	// Get batch of products
// 	products, err := repo.GetBatchWithVariants(productIDs)
// 	require.NoError(t, err)
// 	assert.Len(t, products, 3)
// }
