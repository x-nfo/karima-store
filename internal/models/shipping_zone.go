package models

import (
	"time"

	"gorm.io/gorm"
)

type ShippingZoneStatus string

const (
	ShippingZoneActive   ShippingZoneStatus = "active"
	ShippingZoneInactive ShippingZoneStatus = "inactive"
)

type ShippingZone struct {
	ID        uint             `json:"id" gorm:"primaryKey"`
	CreatedAt time.Time        `json:"created_at"`
	UpdatedAt time.Time        `json:"updated_at"`
	DeletedAt gorm.DeletedAt   `json:"-" gorm:"index"`

	// Basic Information
	Name        string             `json:"name" gorm:"not null;size:200"`
	Description string             `json:"description" gorm:"type:text"`
	Status      ShippingZoneStatus `json:"status" gorm:"not null;default:'active'"`

	// Zone Definition
	Regions     []string `json:"regions" gorm:"-"` // List of region codes (e.g., ["ID-JK", "ID-JB"])
	ExcludeRegions []string `json:"exclude_regions" gorm:"-"`    // Regions to exclude from this zone

	// Free Shipping
	FreeShippingEnabled bool    `json:"free_shipping_enabled" gorm:"default:false"`
	FreeShippingThreshold float64 `json:"free_shipping_threshold"` // Minimum order amount for free shipping

	// Shipping Rates (base rates per provider)
	JNEBaseRate    float64 `json:"jne_base_rate" gorm:"default:15000"`    // IDR per kg
	TIKIBaseRate   float64 `json:"tiki_base_rate" gorm:"default:16000"`   // IDR per kg
	POSBaseRate    float64 `json:"pos_base_rate" gorm:"default:14000"`    // IDR per kg
	SiCepatBaseRate float64 `json:"sicepat_base_rate" gorm:"default:13000"` // IDR per kg

	// Additional Costs
	HandlingFee    float64 `json:"handling_fee" gorm:"default:0"`
	MinimumCost    float64 `json:"minimum_cost" gorm:"default:9000"`

	// Validity
	ValidFrom   *time.Time `json:"valid_from"`
	ValidUntil  *time.Time `json:"valid_until"`
}

func (ShippingZone) TableName() string {
	return "shipping_zones"
}
