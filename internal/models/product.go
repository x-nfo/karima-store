package models

import (
	"time"

	"gorm.io/gorm"
)

type ProductCategory string

const (
	CategoryTops        ProductCategory = "tops"
	CategoryBottoms     ProductCategory = "bottoms"
	CategoryDresses     ProductCategory = "dresses"
	CategoryOuterwear   ProductCategory = "outerwear"
	CategoryFootwear    ProductCategory = "footwear"
	CategoryAccessories ProductCategory = "accessories"
)

type ProductStatus string

const (
	StatusAvailable    ProductStatus = "available"
	StatusOutOfStock   ProductStatus = "out_of_stock"
	StatusDiscontinued ProductStatus = "discontinued"
)

type Product struct {
	ID        uint           `json:"id" gorm:"primaryKey"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`

	// Basic Information
	Name        string          `json:"name" gorm:"not null;size:200"`
	Slug        string          `json:"slug" gorm:"uniqueIndex;not null;size:200"`
	Description string          `json:"description" gorm:"type:text"`
	Category    ProductCategory `json:"category" gorm:"not null;size:50"`

	// Pricing
	Price        float64 `json:"price" gorm:"not null"`
	ComparePrice float64 `json:"compare_price"` // Original price (for discount display)
	Discount     float64 `json:"discount"`      // Discount percentage

	// Inventory
	Stock  int           `json:"stock" gorm:"default:0"`
	Status ProductStatus `json:"status" gorm:"not null;default:'available'"`
	SKU    string        `json:"sku" gorm:"uniqueIndex;size:100"`

	// Media
	Media     []Media `json:"media,omitempty" gorm:"foreignKey:ProductID"`
	Thumbnail string  `json:"thumbnail" gorm:"size:255"`

	// Attributes
	Brand      string  `json:"brand" gorm:"size:100"`
	Color      string  `json:"color" gorm:"size:50"`
	Size       string  `json:"size" gorm:"size:50"`
	Material   string  `json:"material" gorm:"size:100"`
	Weight     float64 `json:"weight"`                     // in kg
	Dimensions string  `json:"dimensions" gorm:"size:100"` // LxWxH format

	// SEO
	MetaTitle       string `json:"meta_title" gorm:"size:200"`
	MetaDescription string `json:"meta_description" gorm:"size:500"`

	// Statistics
	ViewCount   int     `json:"view_count" gorm:"default:0"`
	SoldCount   int     `json:"sold_count" gorm:"default:0"`
	Rating      float64 `json:"rating" gorm:"default:0"`
	ReviewCount int     `json:"review_count" gorm:"default:0"`

	// Timestamps
	PublishedAt *time.Time `json:"published_at"`

	// Relations
	Variants   []ProductVariant `json:"variants,omitempty" gorm:"foreignKey:ProductID"`
	OrderItems []OrderItem      `json:"order_items,omitempty" gorm:"foreignKey:ProductID"`
	Reviews    []Review         `json:"reviews,omitempty" gorm:"foreignKey:ProductID"`
	Wishlist   []Wishlist       `json:"wishlist,omitempty" gorm:"foreignKey:ProductID"`
	FlashSales []FlashSale      `json:"flash_sales,omitempty" gorm:"many2many:flash_sale_products;"`
}

func (Product) TableName() string {
	return "products"
}

type ProductVariant struct {
	ID        uint      `json:"id" gorm:"primaryKey"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`

	ProductID uint    `json:"product_id" gorm:"not null;index"`
	Name      string  `json:"name" gorm:"not null;size:100"` // e.g., "Small - Red"
	Size      string  `json:"size" gorm:"size:50"`
	Color     string  `json:"color" gorm:"size:50"`
	Price     float64 `json:"price" gorm:"not null"`
	Stock     int     `json:"stock" gorm:"default:0"`
	SKU       string  `json:"sku" gorm:"uniqueIndex;size:100"`
}

func (ProductVariant) TableName() string {
	return "product_variants"
}
