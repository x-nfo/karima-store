package models

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestProduct_GenerateSlug(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "Simple name",
			input:    "Test Product",
			expected: "test-product",
		},
		{
			name:     "Name with special characters",
			input:    "Test @Product#$",
			expected: "test-product",
		},
		{
			name:     "Name with multiple spaces",
			input:    "Test   Product   Name",
			expected: "test-product-name",
		},
		{
			name:     "Name with numbers",
			input:    "Product 123",
			expected: "product-123",
		},
		{
			name:     "Empty name",
			input:    "",
			expected: "product",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			product := &Product{Name: tt.input}
			result := product.GenerateSlug()
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestProduct_Validate(t *testing.T) {
	tests := []struct {
		name    string
		product *Product
		wantErr bool
	}{
		{
			name: "Valid product",
			product: &Product{
				Name:     "Test Product",
				Price:    100.00,
				Category: CategoryTops,
				Stock:    10,
				Status:   StatusAvailable,
			},
			wantErr: false,
		},
		{
			name: "Missing name",
			product: &Product{
				Price:    100.00,
				Category: CategoryTops,
				Stock:    10,
				Status:   StatusAvailable,
			},
			wantErr: true,
		},
		{
			name: "Invalid price",
			product: &Product{
				Name:     "Test Product",
				Price:    -10.00,
				Category: CategoryTops,
				Stock:    10,
				Status:   StatusAvailable,
			},
			wantErr: true,
		},
		{
			name: "Missing category",
			product: &Product{
				Name:   "Test Product",
				Price:  100.00,
				Stock:  10,
				Status: StatusAvailable,
			},
			wantErr: true,
		},
		{
			name: "Negative stock",
			product: &Product{
				Name:     "Test Product",
				Price:    100.00,
				Category: CategoryTops,
				Stock:    -5,
				Status:   StatusAvailable,
			},
			wantErr: true,
		},
		{
			name: "Invalid status",
			product: &Product{
				Name:     "Test Product",
				Price:    100.00,
				Category: CategoryTops,
				Stock:    10,
				Status:   "invalid",
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.product.Validate()
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestProduct_IsAvailable(t *testing.T) {
	tests := []struct {
		name     string
		product  *Product
		expected bool
	}{
		{
			name: "Available product",
			product: &Product{
				Status: StatusAvailable,
				Stock:  10,
			},
			expected: true,
		},
		{
			name: "Out of stock",
			product: &Product{
				Status: StatusAvailable,
				Stock:  0,
			},
			expected: false,
		},
		{
			name: "Unavailable status",
			product: &Product{
				Status: StatusUnavailable,
				Stock:  10,
			},
			expected: false,
		},
		{
			name: "Draft status",
			product: &Product{
				Status: StatusDraft,
				Stock:  10,
			},
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.product.IsAvailable()
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestProduct_HasStock(t *testing.T) {
	tests := []struct {
		name     string
		product  *Product
		quantity int
		expected bool
	}{
		{
			name: "Sufficient stock",
			product: &Product{
				Stock: 10,
			},
			quantity: 5,
			expected: true,
		},
		{
			name: "Insufficient stock",
			product: &Product{
				Stock: 5,
			},
			quantity: 10,
			expected: false,
		},
		{
			name: "Exact stock",
			product: &Product{
				Stock: 10,
			},
			quantity: 10,
			expected: true,
		},
		{
			name: "Zero stock",
			product: &Product{
				Stock: 0,
			},
			quantity: 1,
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.product.HasStock(tt.quantity)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestProduct_CalculateDiscountedPrice(t *testing.T) {
	tests := []struct {
		name          string
		product       *Product
		discount      float64
		expectedPrice float64
	}{
		{
			name: "No discount",
			product: &Product{
				Price: 100.00,
			},
			discount:      0,
			expectedPrice: 100.00,
		},
		{
			name: "10% discount",
			product: &Product{
				Price: 100.00,
			},
			discount:      10,
			expectedPrice: 90.00,
		},
		{
			name: "50% discount",
			product: &Product{
				Price: 200.00,
			},
			discount:      50,
			expectedPrice: 100.00,
		},
		{
			name: "100% discount",
			product: &Product{
				Price: 100.00,
			},
			discount:      100,
			expectedPrice: 0.00,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.product.CalculateDiscountedPrice(tt.discount)
			assert.Equal(t, tt.expectedPrice, result)
		})
	}
}

func TestProduct_IsFeatured(t *testing.T) {
	tests := []struct {
		name     string
		product  *Product
		expected bool
	}{
		{
			name: "Featured product",
			product: &Product{
				IsFeatured: true,
			},
			expected: true,
		},
		{
			name: "Not featured product",
			product: &Product{
				IsFeatured: false,
			},
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.product.IsFeatured()
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestProduct_GetStatus(t *testing.T) {
	tests := []struct {
		name     string
		product  *Product
		expected string
	}{
		{
			name: "Available status",
			product: &Product{
				Status: StatusAvailable,
			},
			expected: "Available",
		},
		{
			name: "Unavailable status",
			product: &Product{
				Status: StatusUnavailable,
			},
			expected: "Unavailable",
		},
		{
			name: "Draft status",
			product: &Product{
				Status: StatusDraft,
			},
			expected: "Draft",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.product.GetStatus()
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestProductCategory_String(t *testing.T) {
	tests := []struct {
		name     string
		category ProductCategory
		expected string
	}{
		{
			name:     "Tops category",
			category: CategoryTops,
			expected: "tops",
		},
		{
			name:     "Bottoms category",
			category: CategoryBottoms,
			expected: "bottoms",
		},
		{
			name:     "Dresses category",
			category: CategoryDresses,
			expected: "dresses",
		},
		{
			name:     "Outerwear category",
			category: CategoryOuterwear,
			expected: "outerwear",
		},
		{
			name:     "Footwear category",
			category: CategoryFootwear,
			expected: "footwear",
		},
		{
			name:     "Accessories category",
			category: CategoryAccessories,
			expected: "accessories",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.category.String()
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestProductCategory_Validate(t *testing.T) {
	tests := []struct {
		name     string
		category ProductCategory
		wantErr  bool
	}{
		{
			name:     "Valid category",
			category: CategoryTops,
			wantErr:  false,
		},
		{
			name:     "Invalid category",
			category: ProductCategory("invalid"),
			wantErr:  true,
		},
		{
			name:     "Empty category",
			category: ProductCategory(""),
			wantErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.category.Validate()
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
