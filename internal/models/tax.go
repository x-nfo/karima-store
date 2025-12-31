package models

import (
	"time"

	"gorm.io/gorm"
)

type TaxType string

const (
	TaxTypePercentage TaxType = "percentage"
	TaxTypeFixed      TaxType = "fixed"
)

type TaxStatus string

const (
	TaxStatusActive   TaxStatus = "active"
	TaxStatusInactive TaxStatus = "inactive"
)

type Tax struct {
	ID        uint           `json:"id" gorm:"primaryKey"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`

	// Basic Information
	Name        string    `json:"name" gorm:"not null;size:200"`
	Description string    `json:"description" gorm:"type:text"`
	Type        TaxType   `json:"type" gorm:"not null"`
	Status      TaxStatus `json:"status" gorm:"not null;default:'active'"`

	// Tax Rate
	Rate float64 `json:"rate" gorm:"not null"` // Percentage or fixed amount

	// Applicability
	IsDefault      bool     `json:"is_default" gorm:"default:false"`
	ApplyToShipping bool     `json:"apply_to_shipping" gorm:"default:false"` // Apply tax to shipping cost

	// Validity
	ValidFrom   *time.Time `json:"valid_from"`
	ValidUntil  *time.Time `json:"valid_until"`

	// Regional Restrictions
	ApplicableRegions []string `json:"applicable_regions,omitempty" gorm:"-"` // Region codes (e.g., "ID-JK", "ID-JB")
	ExcludeRegions    []string `json:"exclude_regions,omitempty" gorm:"-"`    // Regions to exclude
}

func (Tax) TableName() string {
	return "taxes"
}
