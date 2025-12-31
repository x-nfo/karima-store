package models

import (
	"time"

	"gorm.io/gorm"
)

type Review struct {
	ID        uint           `json:"id" gorm:"primaryKey"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`

	// Review Information
	ProductID uint    `json:"product_id" gorm:"not null;index"`
	Product   Product `json:"product,omitempty" gorm:"foreignKey:ProductID"`
	UserID    uint    `json:"user_id" gorm:"not null;index"`
	User      User    `json:"user,omitempty" gorm:"foreignKey:UserID"`

	// Rating & Content
	Rating    int    `json:"rating" gorm:"not null;min:1;max:5"`
	Title     string `json:"title" gorm:"size:200"`
	Comment   string `json:"comment" gorm:"type:text"`

	// Media
	Images    []ReviewImage `json:"images,omitempty" gorm:"foreignKey:ReviewID"`

	// Verification
	IsVerified bool `json:"is_verified" gorm:"default:false"` // Verified purchase

	// Moderation
	IsApproved bool   `json:"is_approved" gorm:"default:false"`
	AdminNotes string `json:"admin_notes" gorm:"size:500"`

	// Statistics
	HelpfulCount int `json:"helpful_count" gorm:"default:0"`
	NotHelpfulCount int `json:"not_helpful_count" gorm:"default:0"`
}

func (Review) TableName() string {
	return "reviews"
}

type ReviewImage struct {
	ID        uint      `json:"id" gorm:"primaryKey"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`

	ReviewID uint   `json:"review_id" gorm:"not null;index"`
	URL      string `json:"url" gorm:"not null;size:500"`
	AltText  string `json:"alt_text" gorm:"size:255"`
}

func (ReviewImage) TableName() string {
	return "review_images"
}
