package models

import (
	"time"
)

type StockLog struct {
	ID            uint      `json:"id" gorm:"primaryKey"`
	ProductID     uint      `json:"product_id"`
	VariantID     *uint     `json:"variant_id"`
	ChangeAmount  int       `json:"change_amount"`
	PreviousStock int       `json:"previous_stock"`
	NewStock      int       `json:"new_stock"`
	Reason        string    `json:"reason"`
	ReferenceID   string    `json:"reference_id"` // E.g., Order Number
	CreatedAt     time.Time `json:"created_at"`
}
