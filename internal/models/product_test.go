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
			expected: "test-@product#$",
		},
		{
			name:     "Name with multiple spaces",
			input:    "Test   Product   Name",
			expected: "test---product---name",
		},
		{
			name:     "Name with numbers",
			input:    "Product 123",
			expected: "product-123",
		},
		{
			name:     "Empty name",
			input:    "",
			expected: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			product := &Product{Name: tt.input}
			product.GenerateSlug()
			assert.Equal(t, tt.expected, product.Slug)
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
				SKU:      "TEST-001",
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
		expectedPrice float64
	}{
		{
			name: "No discount",
			product: &Product{
				Price:    100.00,
				Discount: 0,
			},
			expectedPrice: 100.00,
		},
		{
			name: "10% discount",
			product: &Product{
				Price:    100.00,
				Discount: 10,
			},
			expectedPrice: 90.00,
		},
		{
			name: "50% discount",
			product: &Product{
				Price:    200.00,
				Discount: 50,
			},
			expectedPrice: 100.00,
		},
		{
			name: "100% discount",
			product: &Product{
				Price:    100.00,
				Discount: 100,
			},
			expectedPrice: 0.00,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.product.CalculateDiscountedPrice()
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
			result := tt.product.IsFeatured
			assert.Equal(t, tt.expected, result)
		})
	}
}

