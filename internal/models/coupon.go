package models

import (
	"time"

	"gorm.io/gorm"
)

type CouponType string

const (
	CouponTypePercentage CouponType = "percentage"
	CouponTypeFixed      CouponType = "fixed"
)

type CouponStatus string

const (
	CouponStatusActive   CouponStatus = "active"
	CouponStatusInactive CouponStatus = "inactive"
	CouponStatusExpired  CouponStatus = "expired"
)

type Coupon struct {
	ID        uint           `json:"id" gorm:"primaryKey"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`

	// Basic Information
	Code        string       `json:"code" gorm:"uniqueIndex;not null;size:50"`
	Name        string       `json:"name" gorm:"not null;size:200"`
	Description string       `json:"description" gorm:"type:text"`
	Type        CouponType   `json:"type" gorm:"not null"`
	Status      CouponStatus `json:"status" gorm:"not null;default:'active'"`

	// Discount
	DiscountValue float64 `json:"discount_value" gorm:"not null"` // Percentage or fixed amount
	MaxDiscount   float64 `json:"max_discount"`                   // Maximum discount for percentage coupons

	// Usage Limits
	MinPurchaseAmount float64 `json:"min_purchase_amount"` // Minimum purchase amount to apply coupon
	MaxUsageCount     int     `json:"max_usage_count"`     // Maximum times coupon can be used
	UsageCount        int     `json:"usage_count" gorm:"default:0"`
	MaxUsagePerUser   int     `json:"max_usage_per_user"` // Maximum times a user can use this coupon

	// Validity
	ValidFrom  *time.Time `json:"valid_from"`
	ValidUntil *time.Time `json:"valid_until"`

	// Restrictions
	ApplicableProducts   []uint   `json:"applicable_products,omitempty" gorm:"-"`   // Product IDs this coupon applies to
	ExcludeProducts      []uint   `json:"exclude_products,omitempty" gorm:"-"`      // Product IDs this coupon excludes
	ApplicableCategories []string `json:"applicable_categories,omitempty" gorm:"-"` // Categories this coupon applies to

	// Customer Type Restrictions
	ForRetail   bool `json:"for_retail"`
	ForReseller bool `json:"for_reseller"`

	// Statistics
	TotalDiscountUsed float64 `json:"total_discount_used" gorm:"default:0"`
	OrderCount        int     `json:"order_count" gorm:"default:0"`
}

func (Coupon) TableName() string {
	return "coupons"
}

type CouponUsage struct {
	ID        uint      `json:"id" gorm:"primaryKey"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`

	CouponID uint   `json:"coupon_id" gorm:"not null;index"`
	Coupon   Coupon `json:"coupon,omitempty" gorm:"foreignKey:CouponID"`

	UserID uint `json:"user_id" gorm:"not null;index"`

	OrderID uint `json:"order_id" gorm:"not null;index"`

	DiscountAmount float64 `json:"discount_amount" gorm:"not null"`
}

func (CouponUsage) TableName() string {
	return "coupon_usages"
}
