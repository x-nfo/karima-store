package models

import (
	"time"

	"gorm.io/gorm"
)

type OrderStatus string

const (
	StatusPending    OrderStatus = "pending"
	StatusConfirmed  OrderStatus = "confirmed"
	StatusProcessing OrderStatus = "processing"
	StatusShipped    OrderStatus = "shipped"
	StatusDelivered  OrderStatus = "delivered"
	StatusCancelled  OrderStatus = "cancelled"
	StatusRefunded   OrderStatus = "refunded"
)

type PaymentStatus string

const (
	PaymentPending   PaymentStatus = "pending"
	PaymentPaid      PaymentStatus = "paid"
	PaymentFailed    PaymentStatus = "failed"
	PaymentRefunded  PaymentStatus = "refunded"
)

type PaymentMethod string

const (
	PaymentBankTransfer PaymentMethod = "bank_transfer"
	PaymentCreditCard   PaymentMethod = "credit_card"
	PaymentEWallet      PaymentMethod = "e_wallet"
	PaymentCOD          PaymentMethod = "cod"
)

type Order struct {
	ID        uint           `json:"id" gorm:"primaryKey"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`

	// Order Information
	OrderNumber string       `json:"order_number" gorm:"uniqueIndex;not null;size:50"`
	UserID      uint         `json:"user_id" gorm:"not null;index"`
	User        User         `json:"user,omitempty" gorm:"foreignKey:UserID"`

	// Status
	Status       OrderStatus   `json:"status" gorm:"not null;default:'pending'"`
	PaymentStatus PaymentStatus `json:"payment_status" gorm:"not null;default:'pending'"`
	PaymentMethod PaymentMethod `json:"payment_method" gorm:"not null;size:50"`

	// Pricing
	Subtotal      float64 `json:"subtotal" gorm:"not null"`
	Discount      float64 `json:"discount" gorm:"default:0"`
	ShippingCost  float64 `json:"shipping_cost" gorm:"default:0"`
	Tax           float64 `json:"tax" gorm:"default:0"`
	TotalAmount   float64 `json:"total_amount" gorm:"not null"`

	// Shipping Information
	ShippingName    string `json:"shipping_name" gorm:"not null;size:100"`
	ShippingPhone   string `json:"shipping_phone" gorm:"not null;size:20"`
	ShippingAddress string `json:"shipping_address" gorm:"not null;size:255"`
	ShippingCity    string `json:"shipping_city" gorm:"not null;size:100"`
	ShippingProvince string `json:"shipping_province" gorm:"not null;size:100"`
	ShippingPostalCode string `json:"shipping_postal_code" gorm:"not null;size:10"`

	// Tracking
	TrackingNumber string `json:"tracking_number" gorm:"size:100"`
	ShippingProvider string `json:"shipping_provider" gorm:"size:100"`

	// Timestamps
	ConfirmedAt   *time.Time `json:"confirmed_at"`
	ShippedAt     *time.Time `json:"shipped_at"`
	DeliveredAt   *time.Time `json:"delivered_at"`
	CancelledAt   *time.Time `json:"cancelled_at"`
	CancelReason  string     `json:"cancel_reason" gorm:"size:500"`

	// Notes
	CustomerNotes string `json:"customer_notes" gorm:"type:text"`
	AdminNotes     string `json:"admin_notes" gorm:"type:text"`

	// Relations
	Items []OrderItem `json:"items,omitempty" gorm:"foreignKey:OrderID"`
}

func (Order) TableName() string {
	return "orders"
}

type OrderItem struct {
	ID        uint      `json:"id" gorm:"primaryKey"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`

	OrderID   uint    `json:"order_id" gorm:"not null;index"`
	ProductID uint    `json:"product_id" gorm:"not null;index"`
	Product   Product `json:"product,omitempty" gorm:"foreignKey:ProductID"`

	// Product snapshot (to preserve price/details at time of order)
	ProductName  string  `json:"product_name" gorm:"not null;size:200"`
	ProductSKU   string  `json:"product_sku" gorm:"size:100"`
	ProductImage string  `json:"product_image" gorm:"size:500"`

	// Order details
	Quantity    int     `json:"quantity" gorm:"not null"`
	UnitPrice    float64 `json:"unit_price" gorm:"not null"`
	TotalPrice   float64 `json:"total_price" gorm:"not null"`

	// Variant info (if applicable)
	VariantName string `json:"variant_name" gorm:"size:100"`
	VariantSize string `json:"variant_size" gorm:"size:50"`
	VariantColor string `json:"variant_color" gorm:"size:50"`
}

func (OrderItem) TableName() string {
	return "order_items"
}
