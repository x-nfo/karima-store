package models

import (
	"time"

	"gorm.io/gorm"
)

type MediaType string

const (
	MediaTypeImage MediaType = "image"
	MediaTypeVideo MediaType = "video"
)

type MediaStatus string

const (
	MediaStatusActive   MediaStatus = "active"
	MediaStatusInactive MediaStatus = "inactive"
)

type Media struct {
	ID        uint           `json:"id" gorm:"primaryKey"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`

	// Basic Information
	Type        MediaType   `json:"type" gorm:"not null;size:20"`
	URL         string     `json:"url" gorm:"not null;size:500"`
	AltText     string     `json:"alt_text" gorm:"size:255"`
	Status      MediaStatus `json:"status" gorm:"not null;default:'active'"`
	Position    int        `json:"position" gorm:"default:0"`
	IsPrimary   bool       `json:"is_primary" gorm:"default:false"`

	// File Information
	FileName    string  `json:"file_name" gorm:"size:255"`
	FileSize    int64   `json:"file_size" gorm:"default:0"` // in bytes
	ContentType string  `json:"content_type" gorm:"size:100"`
	Width       int     `json:"width" gorm:"default:0"`
	Height      int     `json:"height" gorm:"default:0"`

	// Storage Information
	StorageProvider string `json:"storage_provider" gorm:"size:50"` // "local", "s3", "r2"
	StoragePath    string `json:"storage_path" gorm:"size:500"`

	// Relations
	ProductID uint `json:"product_id" gorm:"not null;index"`
	Product   *Product `json:"product,omitempty" gorm:"foreignKey:ProductID"`
}

func (Media) TableName() string {
	return "media"
}
