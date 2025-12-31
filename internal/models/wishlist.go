package models

import (
	"time"

	"gorm.io/gorm"
)

type Wishlist struct {
	ID        uint           `json:"id" gorm:"primaryKey"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`

	UserID    uint    `json:"user_id" gorm:"not null;index"`
	User      User    `json:"user,omitempty" gorm:"foreignKey:UserID"`
	ProductID uint    `json:"product_id" gorm:"not null;index"`
	Product   Product `json:"product,omitempty" gorm:"foreignKey:ProductID"`

	// Unique constraint to prevent duplicate wishlist items
	// This is handled by uniqueIndex on (user_id, product_id)
}

func (Wishlist) TableName() string {
	return "wishlists"
}
