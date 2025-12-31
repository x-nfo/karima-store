package models

import (
	"time"

	"gorm.io/gorm"
)

type Cart struct {
	ID        uint           `json:"id" gorm:"primaryKey"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`

	UserID uint `json:"user_id" gorm:"uniqueIndex;not null;index"`
	User   User `json:"user,omitempty" gorm:"foreignKey:UserID"`

	// Relations
	Items []CartItem `json:"items,omitempty" gorm:"foreignKey:CartID"`
}

func (Cart) TableName() string {
	return "carts"
}

type CartItem struct {
	ID        uint      `json:"id" gorm:"primaryKey"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`

	CartID    uint `json:"cart_id" gorm:"not null;index"`
	Cart      Cart `json:"cart,omitempty" gorm:"foreignKey:CartID"`

	ProductID       uint   `json:"product_id" gorm:"not null;index"`
	Product         Product `json:"product,omitempty" gorm:"foreignKey:ProductID"`
	ProductVariantID *uint  `json:"product_variant_id" gorm:"index"`
	ProductVariant  *ProductVariant `json:"product_variant,omitempty" gorm:"foreignKey:ProductVariantID"`

	// Product snapshot
	ProductName  string  `json:"product_name" gorm:"not null;size:200"`
	ProductSKU   string  `json:"product_sku" gorm:"size:100"`
	ProductImage string  `json:"product_image" gorm:"size:500"`
	UnitPrice    float64 `json:"unit_price" gorm:"not null"`

	// Cart item details
	Quantity    int     `json:"quantity" gorm:"not null;default:1"`
	TotalPrice  float64 `json:"total_price" gorm:"not null"`

	// Variant info (if applicable)
	VariantName string `json:"variant_name" gorm:"size:100"`
	VariantSize string `json:"variant_size" gorm:"size:50"`
	VariantColor string `json:"variant_color" gorm:"size:50"`
}

func (CartItem) TableName() string {
	return "cart_items"
}
