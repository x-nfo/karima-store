package models

import (
	"time"

	"gorm.io/gorm"
)

type FlashSaleStatus string

const (
	FlashSaleUpcoming   FlashSaleStatus = "upcoming"
	FlashSaleActive     FlashSaleStatus = "active"
	FlashSaleEnded      FlashSaleStatus = "ended"
	FlashSaleCancelled  FlashSaleStatus = "cancelled"
)

type FlashSale struct {
	ID        uint           `json:"id" gorm:"primaryKey"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`

	// Basic Information
	Name        string          `json:"name" gorm:"not null;size:200"`
	Description string          `json:"description" gorm:"type:text"`
	Status      FlashSaleStatus `json:"status" gorm:"not null;default:'upcoming'"`

	// Timing
	StartTime   time.Time  `json:"start_time" gorm:"not null"`
	EndTime     time.Time  `json:"end_time" gorm:"not null"`

	// Discount
	DiscountPercentage float64 `json:"discount_percentage" gorm:"not null"` // e.g., 20 for 20% off

	// Limits
	MaxQuantityPerUser int `json:"max_quantity_per_user" gorm:"default:1"` // Max items per user
	TotalStockLimit    int `json:"total_stock_limit" gorm:"default:0"`      // 0 = unlimited

	// Statistics
	TotalSold   int     `json:"total_sold" gorm:"default:0"`
	TotalOrders int     `json:"total_orders" gorm:"default:0"`
	TotalRevenue float64 `json:"total_revenue" gorm:"default:0"`

	// Relations
	Products   []Product `json:"products,omitempty" gorm:"many2many:flash_sale_products;"`
}

func (FlashSale) TableName() string {
	return "flash_sales"
}

// FlashSaleProduct represents the many-to-many relationship between FlashSale and Product
type FlashSaleProduct struct {
	ID          uint      `json:"id" gorm:"primaryKey"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`

	FlashSaleID uint    `json:"flash_sale_id" gorm:"not null;index"`
	FlashSale   FlashSale `json:"flash_sale,omitempty" gorm:"foreignKey:FlashSaleID"`
	ProductID   uint    `json:"product_id" gorm:"not null;index"`
	Product     Product `json:"product,omitempty" gorm:"foreignKey:ProductID"`

	// Flash sale specific settings for this product
	FlashSalePrice float64 `json:"flash_sale_price" gorm:"not null"` // Special price during flash sale
	FlashSaleStock int     `json:"flash_sale_stock" gorm:"not null"` // Available stock for flash sale
	SoldCount     int     `json:"sold_count" gorm:"default:0"`      // How many sold
}

func (FlashSaleProduct) TableName() string {
	return "flash_sale_products"
}
