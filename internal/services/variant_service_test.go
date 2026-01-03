package services

import (
	"testing"

	"github.com/karima-store/internal/models"
	"github.com/karima-store/internal/repository"
	"github.com/karima-store/internal/test_setup"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"
)

func setupVariantServiceTest(t *testing.T) (*gorm.DB, *models.Product, VariantService, func()) {
	db, cleanup := test_setup.SetupTestDB(t)

	// Clean up existing data
	db.Exec("DELETE FROM product_variants")
	db.Exec("DELETE FROM products")

	// Create a test product
	product := &models.Product{
		Name:     "Test Product",
		Price:    100.00,
		Category: models.CategoryTops,
		Stock:    100,
		Status:   models.StatusAvailable,
		Slug:     "test-product",
		SKU:      "test-product-sku",
		Weight:   0.5,
	}
	err := db.Create(product).Error
	require.NoError(t, err)

	variantRepo := repository.NewVariantRepository(db)
	productRepo := repository.NewProductRepository(db)
	service := NewVariantService(variantRepo, productRepo)

	return db, product, service, cleanup
}

func createTestVariant(productID uint, name, size, color string, price float64, stock int) *models.ProductVariant {
	return &models.ProductVariant{
		ProductID: productID,
		Name:      name,
		Size:      size,
		Color:     color,
		Price:     price,
		Stock:     stock,
	}
}

// ============================================
// SKU GENERATION TESTS (Color/Size Combinations)
// ============================================

func TestVariantService_GenerateSKU_BasicCombinations(t *testing.T) {
	_, _, service, cleanup := setupVariantServiceTest(t)
	defer cleanup()

	testCases := []struct {
		name          string
		productName   string
		size          string
		color         string
		expectedSKU   string
		description   string
	}{
		{
			name:        "Small Red",
			productName: "TShirt",
			size:        "S",
			color:       "Red",
			expectedSKU: "TSH-S-RED",
			description: "Basic small red variant",
		},
		{
			name:        "Medium Blue",
			productName: "TShirt",
			size:        "M",
			color:       "Blue",
			expectedSKU: "TSH-M-BLU",
			description: "Basic medium blue variant",
		},
		{
			name:        "Large Green",
			productName: "TShirt",
			size:        "L",
			color:       "Green",
			expectedSKU: "TSH-L-GRE",
			description: "Basic large green variant",
		},
		{
			name:        "Extra Large Black",
			productName: "TShirt",
			size:        "XL",
			color:       "Black",
			expectedSKU: "TSH-XL-BLA",
			description: "Extra large black variant",
		},
		{
			name:        "Double Extra Large White",
			productName: "TShirt",
			size:        "XXL",
			color:       "White",
			expectedSKU: "TSH-XX-WHI",
			description: "Double extra large white variant",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			sku := service.GenerateSKU(tc.productName, tc.size, tc.color)
			assert.Equal(t, tc.expectedSKU, sku, tc.description)
		})
	}
}

func TestVariantService_GenerateSKU_ColorVariations(t *testing.T) {
	_, _, service, cleanup := setupVariantServiceTest(t)
	defer cleanup()

	testCases := []struct {
		name        string
		productName string
		size        string
		color       string
		expectedSKU string
	}{
		{
			name:        "Red Color",
			productName: "Dress",
			size:        "M",
			color:       "Red",
			expectedSKU: "DRE-M-RED",
		},
		{
			name:        "Blue Color",
			productName: "Dress",
			size:        "M",
			color:       "Blue",
			expectedSKU: "DRE-M-BLU",
		},
		{
			name:        "Green Color",
			productName: "Dress",
			size:        "M",
			color:       "Green",
			expectedSKU: "DRE-M-GRE",
		},
		{
			name:        "Yellow Color",
			productName: "Dress",
			size:        "M",
			color:       "Yellow",
			expectedSKU: "DRE-M-YEL",
		},
		{
			name:        "Purple Color",
			productName: "Dress",
			size:        "M",
			color:       "Purple",
			expectedSKU: "DRE-M-PUR",
		},
		{
			name:        "Orange Color",
			productName: "Dress",
			size:        "M",
			color:       "Orange",
			expectedSKU: "DRE-M-ORA",
		},
		{
			name:        "Pink Color",
			productName: "Dress",
			size:        "M",
			color:       "Pink",
			expectedSKU: "DRE-M-PIN",
		},
		{
			name:        "Brown Color",
			productName: "Dress",
			size:        "M",
			color:       "Brown",
			expectedSKU: "DRE-M-BRO",
		},
		{
			name:        "Gray Color",
			productName: "Dress",
			size:        "M",
			color:       "Gray",
			expectedSKU: "DRE-M-GRA",
		},
		{
			name:        "Black Color",
			productName: "Dress",
			size:        "M",
			color:       "Black",
			expectedSKU: "DRE-M-BLA",
		},
		{
			name:        "White Color",
			productName: "Dress",
			size:        "M",
			color:       "White",
			expectedSKU: "DRE-M-WHI",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			sku := service.GenerateSKU(tc.productName, tc.size, tc.color)
			assert.Equal(t, tc.expectedSKU, sku)
		})
	}
}

func TestVariantService_GenerateSKU_SizeVariations(t *testing.T) {
	_, _, service, cleanup := setupVariantServiceTest(t)
	defer cleanup()

	testCases := []struct {
		name        string
		productName string
		size        string
		color       string
		expectedSKU string
	}{
		{
			name:        "Extra Small",
			productName: "Jeans",
			size:        "XS",
			color:       "Blue",
			expectedSKU: "JEA-XS-BLU",
		},
		{
			name:        "Small",
			productName: "Jeans",
			size:        "S",
			color:       "Blue",
			expectedSKU: "JEA-S-BLU",
		},
		{
			name:        "Medium",
			productName: "Jeans",
			size:        "M",
			color:       "Blue",
			expectedSKU: "JEA-M-BLU",
		},
		{
			name:        "Large",
			productName: "Jeans",
			size:        "L",
			color:       "Blue",
			expectedSKU: "JEA-L-BLU",
		},
		{
			name:        "Extra Large",
			productName: "Jeans",
			size:        "XL",
			color:       "Blue",
			expectedSKU: "JEA-XL-BLU",
		},
		{
			name:        "Double Extra Large",
			productName: "Jeans",
			size:        "XXL",
			color:       "Blue",
			expectedSKU: "JEA-XX-BLU",
		},
		{
			name:        "Triple Extra Large",
			productName: "Jeans",
			size:        "XXXL",
			color:       "Blue",
			expectedSKU: "JEA-XX-BLU",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			sku := service.GenerateSKU(tc.productName, tc.size, tc.color)
			assert.Equal(t, tc.expectedSKU, sku)
		})
	}
}

func TestVariantService_GenerateSKU_SpecialCharacters(t *testing.T) {
	_, _, service, cleanup := setupVariantServiceTest(t)
	defer cleanup()

	testCases := []struct {
		name        string
		productName string
		size        string
		color       string
		expectedSKU string
		description string
	}{
		{
			name:        "Color with space",
			productName: "TShirt",
			size:        "M",
			color:       "Dark Blue",
			expectedSKU: "TSH-M-DAR",
			description: "Should handle spaces in color",
		},
		{
			name:        "Color with hyphen",
			productName: "TShirt",
			size:        "M",
			color:       "Dark-Blue",
			expectedSKU: "TSH-M-DAR",
			description: "Should handle hyphens in color",
		},
		{
			name:        "Color with special chars",
			productName: "TShirt",
			size:        "M",
			color:       "Navy@Blue!",
			expectedSKU: "TSH-M-NAV",
			description: "Should remove special characters",
		},
		{
			name:        "Size with special chars",
			productName: "TShirt",
			size:        "M+",
			color:       "Red",
			expectedSKU: "TSH-M-RED",
			description: "Should handle special characters in size",
		},
		{
			name:        "Product name with spaces",
			productName: "Long Sleeve Shirt",
			size:        "M",
			color:       "Red",
			expectedSKU: "LON-M-RED",
			description: "Should handle spaces in product name",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			sku := service.GenerateSKU(tc.productName, tc.size, tc.color)
			assert.Equal(t, tc.expectedSKU, sku, tc.description)
		})
	}
}

func TestVariantService_GenerateSKU_EdgeCases(t *testing.T) {
	_, _, service, cleanup := setupVariantServiceTest(t)
	defer cleanup()

	testCases := []struct {
		name        string
		productName string
		size        string
		color       string
		expectedSKU string
		description string
	}{
		{
			name:        "Short product name",
			productName: "AB",
			size:        "M",
			color:       "Red",
			expectedSKU: "AB-M-RED",
			description: "Should handle short product names",
		},
		{
			name:        "Single character size",
			productName: "TShirt",
			size:        "L",
			color:       "Red",
			expectedSKU: "TSH-L-RED",
			description: "Should handle single character size",
		},
		{
			name:        "Short color",
			productName: "TShirt",
			size:        "M",
			color:       "R",
			expectedSKU: "TSH-M-R",
			description: "Should handle short color",
		},
		{
			name:        "Empty size",
			productName: "TShirt",
			size:        "",
			color:       "Red",
			expectedSKU: "TSH-RED",
			description: "Should handle empty size (hyphens are cleaned)",
		},
		{
			name:        "Empty color",
			productName: "TShirt",
			size:        "M",
			color:       "",
			expectedSKU: "TSH-M",
			description: "Should handle empty color (trailing hyphen is removed)",
		},
		{
			name:        "Numeric values",
			productName: "TShirt",
			size:        "36",
			color:       "1",
			expectedSKU: "TSH-36-1",
			description: "Should handle numeric values",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			sku := service.GenerateSKU(tc.productName, tc.size, tc.color)
			assert.Equal(t, tc.expectedSKU, sku, tc.description)
		})
	}
}

func TestVariantService_GenerateSKU_CaseInsensitivity(t *testing.T) {
	_, _, service, cleanup := setupVariantServiceTest(t)
	defer cleanup()

	testCases := []struct {
		name        string
		productName string
		size        string
		color       string
		expectedSKU string
	}{
		{
			name:        "Lowercase product name",
			productName: "tshirt",
			size:        "M",
			color:       "Red",
			expectedSKU: "TSH-M-RED",
		},
		{
			name:        "Uppercase product name",
			productName: "TSHIRT",
			size:        "M",
			color:       "Red",
			expectedSKU: "TSH-M-RED",
		},
		{
			name:        "Mixed case product name",
			productName: "TShIrT",
			size:        "M",
			color:       "Red",
			expectedSKU: "TSH-M-RED",
		},
		{
			name:        "Lowercase size",
			productName: "TShirt",
			size:        "m",
			color:       "Red",
			expectedSKU: "TSH-M-RED",
		},
		{
			name:        "Lowercase color",
			productName: "TShirt",
			size:        "M",
			color:       "red",
			expectedSKU: "TSH-M-RED",
		},
		{
			name:        "Mixed case color",
			productName: "TShirt",
			size:        "M",
			color:       "ReD",
			expectedSKU: "TSH-M-RED",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			sku := service.GenerateSKU(tc.productName, tc.size, tc.color)
			assert.Equal(t, tc.expectedSKU, sku)
		})
	}
}

// ============================================
// CREATE VARIANT TESTS
// ============================================

func TestVariantService_CreateVariant_Success(t *testing.T) {
	db, product, service, cleanup := setupVariantServiceTest(t)
	defer cleanup()

	variant := createTestVariant(product.ID, "Medium Red", "M", "Red", 110.00, 10)
	err := service.CreateVariant(variant)

	require.NoError(t, err)
	assert.NotZero(t, variant.ID)
	assert.NotEmpty(t, variant.SKU)

	// Verify in database
	var dbVariant models.ProductVariant
	err = db.First(&dbVariant, variant.ID).Error
	require.NoError(t, err)
	assert.Equal(t, "Medium Red", dbVariant.Name)
	assert.Equal(t, "M", dbVariant.Size)
	assert.Equal(t, "Red", dbVariant.Color)
}

func TestVariantService_CreateVariant_WithAutoGeneratedSKU(t *testing.T) {
	_, product, service, cleanup := setupVariantServiceTest(t)
	defer cleanup()

	variant := createTestVariant(product.ID, "Medium Blue", "M", "Blue", 110.00, 10)
	variant.SKU = "" // Empty SKU to trigger auto-generation

	err := service.CreateVariant(variant)

	require.NoError(t, err)
	assert.NotEmpty(t, variant.SKU)
	assert.Equal(t, "TES-M-BLU", variant.SKU) // TES from "Test Product"
}

func TestVariantService_CreateVariant_WithCustomSKU(t *testing.T) {
	_, product, service, cleanup := setupVariantServiceTest(t)
	defer cleanup()

	customSKU := "CUSTOM-SKU-123"
	variant := createTestVariant(product.ID, "Custom SKU", "M", "Red", 110.00, 10)
	variant.SKU = customSKU

	err := service.CreateVariant(variant)

	require.NoError(t, err)
	assert.Equal(t, customSKU, variant.SKU)
}

func TestVariantService_CreateVariant_DuplicateSKU(t *testing.T) {
	_, product, service, cleanup := setupVariantServiceTest(t)
	defer cleanup()

	// Create first variant
	variant1 := createTestVariant(product.ID, "First Variant", "M", "Red", 110.00, 10)
	variant1.SKU = "DUPLICATE-SKU"
	err := service.CreateVariant(variant1)
	require.NoError(t, err)

	// Try to create second variant with same SKU
	variant2 := createTestVariant(product.ID, "Second Variant", "L", "Blue", 120.00, 15)
	variant2.SKU = "DUPLICATE-SKU"
	err = service.CreateVariant(variant2)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "already exists")
}

func TestVariantService_CreateVariant_MissingProductID(t *testing.T) {
	_, _, service, cleanup := setupVariantServiceTest(t)
	defer cleanup()

	variant := createTestVariant(0, "No Product", "M", "Red", 110.00, 10)
	err := service.CreateVariant(variant)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "product ID is required")
}

func TestVariantService_CreateVariant_MissingName(t *testing.T) {
	_, product, service, cleanup := setupVariantServiceTest(t)
	defer cleanup()

	variant := createTestVariant(product.ID, "", "M", "Red", 110.00, 10)
	err := service.CreateVariant(variant)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "variant name is required")
}

func TestVariantService_CreateVariant_InvalidPrice(t *testing.T) {
	testCases := []struct {
		name        string
		price       float64
		description string
	}{
		{"Zero price", 0, "Price of zero should be invalid"},
		{"Negative price", -10.00, "Negative price should be invalid"},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			_, product, service, cleanup := setupVariantServiceTest(t)
			defer cleanup()

			variant := createTestVariant(product.ID, "Invalid Price", "M", "Red", tc.price, 10)
			err := service.CreateVariant(variant)

			assert.Error(t, err)
			assert.Contains(t, err.Error(), "price must be greater than 0", tc.description)
		})
	}
}

func TestVariantService_CreateVariant_NonExistentProduct(t *testing.T) {
	_, _, service, cleanup := setupVariantServiceTest(t)
	defer cleanup()

	variant := createTestVariant(99999, "Non-existent Product", "M", "Red", 110.00, 10)
	err := service.CreateVariant(variant)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "product not found")
}

// ============================================
// GET VARIANT TESTS
// ============================================

func TestVariantService_GetVariantByID_Success(t *testing.T) {
	_, product, service, cleanup := setupVariantServiceTest(t)
	defer cleanup()

	// Create variant
	variant := createTestVariant(product.ID, "Get By ID", "M", "Red", 110.00, 10)
	err := service.CreateVariant(variant)
	require.NoError(t, err)

	// Get by ID
	fetched, err := service.GetVariantByID(variant.ID)

	require.NoError(t, err)
	assert.Equal(t, variant.ID, fetched.ID)
	assert.Equal(t, "Get By ID", fetched.Name)
	assert.Equal(t, "M", fetched.Size)
	assert.Equal(t, "Red", fetched.Color)
}

func TestVariantService_GetVariantByID_NotFound(t *testing.T) {
	_, _, service, cleanup := setupVariantServiceTest(t)
	defer cleanup()

	_, err := service.GetVariantByID(99999)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "variant not found")
}

func TestVariantService_GetVariantBySKU_Success(t *testing.T) {
	_, product, service, cleanup := setupVariantServiceTest(t)
	defer cleanup()

	// Create variant with specific SKU
	variant := createTestVariant(product.ID, "Get By SKU", "M", "Red", 110.00, 10)
	variant.SKU = "FIND-BY-SKU"
	err := service.CreateVariant(variant)
	require.NoError(t, err)

	// Get by SKU
	fetched, err := service.GetVariantBySKU("FIND-BY-SKU")

	require.NoError(t, err)
	assert.Equal(t, variant.ID, fetched.ID)
	assert.Equal(t, "FIND-BY-SKU", fetched.SKU)
}

func TestVariantService_GetVariantBySKU_NotFound(t *testing.T) {
	_, _, service, cleanup := setupVariantServiceTest(t)
	defer cleanup()

	_, err := service.GetVariantBySKU("NON-EXISTENT-SKU")

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "variant not found")
}

func TestVariantService_GetVariantsByProductID_Success(t *testing.T) {
	_, product, service, cleanup := setupVariantServiceTest(t)
	defer cleanup()

	// Create multiple variants
	sizes := []string{"S", "M", "L"}
	colors := []string{"Red", "Blue"}

	for _, size := range sizes {
		for _, color := range colors {
			variant := createTestVariant(product.ID, size+"-"+color, size, color, 110.00, 10)
			err := service.CreateVariant(variant)
			require.NoError(t, err)
		}
	}

	// Get all variants for product
	variants, err := service.GetVariantsByProductID(product.ID)

	require.NoError(t, err)
	assert.Len(t, variants, 6) // 3 sizes x 2 colors

	// Verify all belong to the correct product
	for _, v := range variants {
		assert.Equal(t, product.ID, v.ProductID)
	}
}

func TestVariantService_GetVariantsByProductID_NonExistentProduct(t *testing.T) {
	_, _, service, cleanup := setupVariantServiceTest(t)
	defer cleanup()

	_, err := service.GetVariantsByProductID(99999)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "product not found")
}

// ============================================
// UPDATE VARIANT TESTS
// ============================================

func TestVariantService_UpdateVariant_Success(t *testing.T) {
	_, product, service, cleanup := setupVariantServiceTest(t)
	defer cleanup()

	// Create variant
	variant := createTestVariant(product.ID, "Original Name", "M", "Red", 110.00, 10)
	err := service.CreateVariant(variant)
	require.NoError(t, err)

	// Update variant
	updatedVariant := createTestVariant(product.ID, "Updated Name", "L", "Blue", 120.00, 15)
	err = service.UpdateVariant(variant.ID, updatedVariant)
	require.NoError(t, err)

	// Verify update
	fetched, err := service.GetVariantByID(variant.ID)
	require.NoError(t, err)
	assert.Equal(t, "Updated Name", fetched.Name)
	assert.Equal(t, "L", fetched.Size)
	assert.Equal(t, "Blue", fetched.Color)
	assert.Equal(t, 120.00, fetched.Price)
	assert.Equal(t, 15, fetched.Stock)
}

func TestVariantService_UpdateVariant_NotFound(t *testing.T) {
	_, product, service, cleanup := setupVariantServiceTest(t)
	defer cleanup()

	variant := createTestVariant(product.ID, "Update Test", "M", "Red", 110.00, 10)
	err := service.UpdateVariant(99999, variant)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "variant not found")
}

func TestVariantService_UpdateVariant_ChangeProduct(t *testing.T) {
	db, product, service, cleanup := setupVariantServiceTest(t)
	defer cleanup()

	// Create second product
	product2 := &models.Product{
		Name:     "Second Product",
		Price:    200.00,
		Category: models.CategoryBottoms,
		Stock:    50,
		Status:   models.StatusAvailable,
		Slug:     "second-product",
		SKU:      "second-product-sku",
		Weight:   0.5,
	}
	err := db.Create(product2).Error
	require.NoError(t, err)

	// Create variant for first product
	variant := createTestVariant(product.ID, "Change Product", "M", "Red", 110.00, 10)
	err = service.CreateVariant(variant)
	require.NoError(t, err)

	// Update to second product
	updatedVariant := createTestVariant(product2.ID, "Changed Product", "M", "Red", 110.00, 10)
	err = service.UpdateVariant(variant.ID, updatedVariant)
	require.NoError(t, err)

	// Verify
	fetched, err := service.GetVariantByID(variant.ID)
	require.NoError(t, err)
	assert.Equal(t, product2.ID, fetched.ProductID)
}

func TestVariantService_UpdateVariant_ChangeToNonExistentProduct(t *testing.T) {
	_, product, service, cleanup := setupVariantServiceTest(t)
	defer cleanup()

	// Create variant
	variant := createTestVariant(product.ID, "Change Product", "M", "Red", 110.00, 10)
	err := service.CreateVariant(variant)
	require.NoError(t, err)

	// Try to update to non-existent product
	updatedVariant := createTestVariant(99999, "Non-existent", "M", "Red", 110.00, 10)
	err = service.UpdateVariant(variant.ID, updatedVariant)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "product not found")
}

// ============================================
// DELETE VARIANT TESTS
// ============================================

func TestVariantService_DeleteVariant_Success(t *testing.T) {
	_, product, service, cleanup := setupVariantServiceTest(t)
	defer cleanup()

	// Create variant
	variant := createTestVariant(product.ID, "To Delete", "M", "Red", 110.00, 10)
	err := service.CreateVariant(variant)
	require.NoError(t, err)

	// Delete variant
	err = service.DeleteVariant(variant.ID)
	require.NoError(t, err)

	// Verify deletion
	_, err = service.GetVariantByID(variant.ID)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "variant not found")
}

func TestVariantService_DeleteVariant_NotFound(t *testing.T) {
	_, _, service, cleanup := setupVariantServiceTest(t)
	defer cleanup()

	err := service.DeleteVariant(99999)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "variant not found")
}

// ============================================
// UPDATE STOCK TESTS
// ============================================

func TestVariantService_UpdateVariantStock_Increase(t *testing.T) {
	_, product, service, cleanup := setupVariantServiceTest(t)
	defer cleanup()

	// Create variant with initial stock
	variant := createTestVariant(product.ID, "Stock Test", "M", "Red", 110.00, 50)
	err := service.CreateVariant(variant)
	require.NoError(t, err)

	// Increase stock
	err = service.UpdateVariantStock(variant.ID, 10)
	require.NoError(t, err)

	// Verify
	fetched, err := service.GetVariantByID(variant.ID)
	require.NoError(t, err)
	assert.Equal(t, 60, fetched.Stock)
}

func TestVariantService_UpdateVariantStock_Decrease(t *testing.T) {
	_, product, service, cleanup := setupVariantServiceTest(t)
	defer cleanup()

	// Create variant with initial stock
	variant := createTestVariant(product.ID, "Stock Test", "M", "Red", 110.00, 50)
	err := service.CreateVariant(variant)
	require.NoError(t, err)

	// Decrease stock
	err = service.UpdateVariantStock(variant.ID, -20)
	require.NoError(t, err)

	// Verify
	fetched, err := service.GetVariantByID(variant.ID)
	require.NoError(t, err)
	assert.Equal(t, 30, fetched.Stock)
}

func TestVariantService_UpdateVariantStock_InsufficientStock(t *testing.T) {
	_, product, service, cleanup := setupVariantServiceTest(t)
	defer cleanup()

	// Create variant with limited stock
	variant := createTestVariant(product.ID, "Low Stock", "M", "Red", 110.00, 10)
	err := service.CreateVariant(variant)
	require.NoError(t, err)

	// Try to decrease more than available
	err = service.UpdateVariantStock(variant.ID, -20)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "insufficient stock")

	// Verify stock unchanged
	fetched, err := service.GetVariantByID(variant.ID)
	require.NoError(t, err)
	assert.Equal(t, 10, fetched.Stock)
}

func TestVariantService_UpdateVariantStock_NotFound(t *testing.T) {
	_, _, service, cleanup := setupVariantServiceTest(t)
	defer cleanup()

	err := service.UpdateVariantStock(99999, 10)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "variant not found")
}

func TestVariantService_UpdateVariantStock_ZeroStock(t *testing.T) {
	_, product, service, cleanup := setupVariantServiceTest(t)
	defer cleanup()

	// Create variant with zero stock
	variant := createTestVariant(product.ID, "Zero Stock", "M", "Red", 110.00, 0)
	err := service.CreateVariant(variant)
	require.NoError(t, err)

	// Try to decrease from zero
	err = service.UpdateVariantStock(variant.ID, -10)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "insufficient stock")
}

// ============================================
// COMPREHENSIVE VARIANT COMBINATION TESTS
// ============================================

func TestVariantService_ColorSizeMatrix(t *testing.T) {
	_, product, service, cleanup := setupVariantServiceTest(t)
	defer cleanup()

	// Test all combinations of sizes and colors
	sizes := []string{"XS", "S", "M", "L", "XL", "XXL"}
	colors := []string{"Red", "Blue", "Green", "Black", "White"}

	expectedCount := len(sizes) * len(colors)

	for _, size := range sizes {
		for _, color := range colors {
			variant := createTestVariant(product.ID, size+"-"+color, size, color, 110.00, 10)
			err := service.CreateVariant(variant)
			require.NoError(t, err, "Failed to create variant %s-%s", size, color)
		}
	}

	// Verify all variants created
	variants, err := service.GetVariantsByProductID(product.ID)
	require.NoError(t, err)
	assert.Len(t, variants, expectedCount)

	// Verify all combinations exist
	createdCombos := make(map[string]bool)
	for _, v := range variants {
		createdCombos[v.Size+"-"+v.Color] = true
	}

	for _, size := range sizes {
		for _, color := range colors {
			key := size + "-" + color
			assert.True(t, createdCombos[key], "Missing combination: %s", key)
		}
	}
}

func TestVariantService_MultipleProductsWithVariants(t *testing.T) {
	db, _, service, cleanup := setupVariantServiceTest(t)
	defer cleanup()

	// Create multiple products
	products := []*models.Product{
		{
			Name:     "T-Shirt",
			Price:    100.00,
			Category: models.CategoryTops,
			Stock:    100,
			Status:   models.StatusAvailable,
			Slug:     "t-shirt",
			SKU:      "t-shirt-sku",
			Weight:   0.5,
		},
		{
			Name:     "Jeans",
			Price:    150.00,
			Category: models.CategoryBottoms,
			Stock:    80,
			Status:   models.StatusAvailable,
			Slug:     "jeans",
			SKU:      "jeans-sku",
			Weight:   0.8,
		},
		{
			Name:     "Dress",
			Price:    200.00,
			Category: models.CategoryDresses,
			Stock:    60,
			Status:   models.StatusAvailable,
			Slug:     "dress",
			SKU:      "dress-sku",
			Weight:   0.6,
		},
	}

	for _, product := range products {
		err := db.Create(product).Error
		require.NoError(t, err)
	}

	// Create variants for each product
	for _, product := range products {
		sizes := []string{"S", "M", "L"}
		colors := []string{"Red", "Blue"}

		for _, size := range sizes {
			for _, color := range colors {
				variant := createTestVariant(product.ID, size+"-"+color, size, color, product.Price, 10)
				err := service.CreateVariant(variant)
				require.NoError(t, err, "Failed to create variant for product %s", product.Name)
			}
		}
	}

	// Verify each product has correct number of variants
	for _, product := range products {
		variants, err := service.GetVariantsByProductID(product.ID)
		require.NoError(t, err)
		assert.Len(t, variants, 6, "Product %s should have 6 variants", product.Name)

		// Verify all variants belong to correct product
		for _, v := range variants {
			assert.Equal(t, product.ID, v.ProductID)
		}
	}
}

// ============================================
// EDGE CASES AND ERROR SCENARIOS
// ============================================

func TestVariantService_ConcurrentVariantCreation(t *testing.T) {
	_, product, service, cleanup := setupVariantServiceTest(t)
	defer cleanup()

	// Try to create variants with same SKU simultaneously
	// This tests the duplicate SKU constraint
	variant1 := createTestVariant(product.ID, "Variant 1", "M", "Red", 110.00, 10)
	variant1.SKU = "CONCURRENT-SKU"

	variant2 := createTestVariant(product.ID, "Variant 2", "L", "Blue", 120.00, 15)
	variant2.SKU = "CONCURRENT-SKU"

	err1 := service.CreateVariant(variant1)
	require.NoError(t, err1)

	err2 := service.CreateVariant(variant2)
	assert.Error(t, err2)
	assert.Contains(t, err2.Error(), "already exists")
}

func TestVariantService_VariantWithZeroPrice(t *testing.T) {
	_, product, service, cleanup := setupVariantServiceTest(t)
	defer cleanup()

	variant := createTestVariant(product.ID, "Zero Price", "M", "Red", 0, 10)
	err := service.CreateVariant(variant)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "price must be greater than 0")
}

func TestVariantService_VariantWithNegativeStock(t *testing.T) {
	_, product, service, cleanup := setupVariantServiceTest(t)
	defer cleanup()

	variant := createTestVariant(product.ID, "Negative Stock", "M", "Red", 110.00, -10)
	err := service.CreateVariant(variant)

	// This might pass at creation but should fail when trying to use the variant
	// The service doesn't validate stock on creation, only on updates
	require.NoError(t, err)
}

func TestVariantService_UpdateVariantToZeroStock(t *testing.T) {
	_, product, service, cleanup := setupVariantServiceTest(t)
	defer cleanup()

	variant := createTestVariant(product.ID, "Zero Stock Update", "M", "Red", 110.00, 10)
	err := service.CreateVariant(variant)
	require.NoError(t, err)

	// Update to zero stock
	err = service.UpdateVariantStock(variant.ID, -10)
	require.NoError(t, err)

	// Verify
	fetched, err := service.GetVariantByID(variant.ID)
	require.NoError(t, err)
	assert.Equal(t, 0, fetched.Stock)
}

func TestVariantService_LongVariantName(t *testing.T) {
	_, product, service, cleanup := setupVariantServiceTest(t)
	defer cleanup()

	longName := "This is a very long variant name that exceeds normal length but should still work"
	variant := createTestVariant(product.ID, longName, "M", "Red", 110.00, 10)
	err := service.CreateVariant(variant)

	// This might fail due to database constraint
	// The model has size:100 constraint
	if err != nil {
		assert.Contains(t, err.Error(), "too long")
	} else {
		assert.NotZero(t, variant.ID)
	}
}

func TestVariantService_EmptySizeAndColor(t *testing.T) {
	_, product, service, cleanup := setupVariantServiceTest(t)
	defer cleanup()

	variant := createTestVariant(product.ID, "No Size or Color", "", "", 110.00, 10)
	err := service.CreateVariant(variant)

	// Should succeed as size and color are optional
	require.NoError(t, err)
	assert.NotZero(t, variant.ID)
}

func TestVariantService_VariantWithoutPriceOverride(t *testing.T) {
	_, product, service, cleanup := setupVariantServiceTest(t)
	defer cleanup()

	// Create variant with same price as product
	variant := createTestVariant(product.ID, "Same Price", "M", "Red", 100.00, 10)
	err := service.CreateVariant(variant)

	require.NoError(t, err)
	assert.Equal(t, 100.00, variant.Price)
}

func TestVariantService_VariantWithPriceOverride(t *testing.T) {
	_, product, service, cleanup := setupVariantServiceTest(t)
	defer cleanup()

	// Create variant with different price
	variant := createTestVariant(product.ID, "Price Override", "M", "Red", 120.00, 10)
	err := service.CreateVariant(variant)

	require.NoError(t, err)
	assert.Equal(t, 120.00, variant.Price)
	assert.NotEqual(t, product.Price, variant.Price)
}

// ============================================
// SERVICE INITIALIZATION TESTS
// ============================================

func TestNewVariantService(t *testing.T) {
	db, cleanup := test_setup.SetupTestDB(t)
	defer cleanup()

	variantRepo := repository.NewVariantRepository(db)
	productRepo := repository.NewProductRepository(db)

	service := NewVariantService(variantRepo, productRepo)

	assert.NotNil(t, service)
}

func TestVariantService_ImplementsInterface(t *testing.T) {
	db, cleanup := test_setup.SetupTestDB(t)
	defer cleanup()

	variantRepo := repository.NewVariantRepository(db)
	productRepo := repository.NewProductRepository(db)

	service := NewVariantService(variantRepo, productRepo)

	// Verify service implements the interface
	var _ VariantService = service
}
