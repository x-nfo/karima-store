package models

// CheckoutRequest represents the checkout request
type CheckoutRequest struct {
	// Cart items
	Items []CheckoutItem `json:"items" validate:"required,min=1"`

	// Shipping information
	ShippingName    string `json:"shipping_name" validate:"required"`
	ShippingPhone   string `json:"shipping_phone" validate:"required"`
	ShippingAddress string `json:"shipping_address" validate:"required"`
	ShippingCity    string `json:"shipping_city" validate:"required"`
	ShippingProvince string `json:"shipping_province" validate:"required"`
	ShippingPostalCode string `json:"shipping_postal_code" validate:"required"`

	// Payment method
	PaymentMethod string `json:"payment_method" validate:"required,oneof=bank_transfer credit_card e_wallet cod"`

	// Customer info
	UserID uint `json:"user_id" validate:"required"`

	// Optional
	CustomerNotes string `json:"customer_notes"`
	CouponCode    string `json:"coupon_code"`
}

// CheckoutItem represents a cart item for checkout
type CheckoutItem struct {
	ProductID uint   `json:"product_id" validate:"required"`
	VariantID uint   `json:"variant_id"`
	Quantity  int    `json:"quantity" validate:"required,min=1"`
}

// CheckoutResponse represents the checkout response
type CheckoutResponse struct {
	OrderNumber    string  `json:"order_number"`
	OrderID        uint    `json:"order_id"`
	SnapToken      string  `json:"snap_token"`
	SnapURL        string  `json:"snap_url"`
	RedirectURL    string  `json:"redirect_url"`
	Amount         float64 `json:"amount"`
	ExpiryTime     string  `json:"expiry_time"`
}

// MidtransPaymentRequest represents the request to Midtrans Snap API
type MidtransPaymentRequest struct {
	TransactionDetails TransactionDetails `json:"transaction_details"`
	CustomerDetails   CustomerDetails   `json:"customer_details"`
	ItemDetails      []ItemDetail      `json:"item_details"`
	EnabledPayments  []string          `json:"enabled_payments"`
	CallbackURL      string           `json:"finish_redirect_url"`
	NotificationURL  string           `json:"notification_url"`
}

// TransactionDetails represents Midtrans transaction details
type TransactionDetails struct {
	OrderID     string  `json:"order_id"`
	GrossAmount  float64 `json:"gross_amount"`
}

// CustomerDetails represents Midtrans customer details
type CustomerDetails struct {
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Email     string `json:"email"`
	Phone     string `json:"phone"`
	BillingAddress  *Address `json:"billing_address"`
	ShippingAddress *Address `json:"shipping_address"`
}

// Address represents Midtrans address
type Address struct {
	FirstName    string `json:"first_name"`
	LastName     string `json:"last_name"`
	Address       string `json:"address"`
	City          string `json:"city"`
	PostalCode    string `json:"postal_code"`
	Phone         string `json:"phone"`
	CountryCode   string `json:"country_code"`
}

// ItemDetail represents Midtrans item details
type ItemDetail struct {
	ID       string  `json:"id"`
	Price    float64 `json:"price"`
	Quantity int     `json:"name"`
	Name     string  `json:"name"`
}

// MidtransPaymentNotification represents the webhook notification from Midtrans
type MidtransPaymentNotification struct {
	TransactionID       string  `json:"transaction_id"`
	StatusMessage       string  `json:"status_message"`
	PaymentType        string  `json:"payment_type"`
	OrderID            string  `json:"order_id"`
	GrossAmount        float64 `json:"gross_amount"`
	FraudStatus        string  `json:"fraud_status"`
	SignatureKey       string  `json:"signature_key"`
	StatusCode         string  `json:"status_code"`
	TransactionStatus  string  `json:"transaction_status"`
	TransactionTime    string  `json:"transaction_time"`
	ApprovalCode       string  `json:"approval_code"`
	PaymentAmounts     []PaymentAmount `json:"payment_amounts"`
	Bank               string  `json:"bank"`
	VANumbers          []VANumber `json:"va_numbers"`
}

// PaymentAmount represents payment amount breakdown
type PaymentAmount struct {
	PaymentType string  `json:"payment_type"`
	Amount      float64 `json:"amount"`
}

// VANumber represents virtual account number
type VANumber struct {
	Bank     string `json:"bank"`
	VANumber string `json:"va_number"`
}

// MidtransSnapResponse represents the response from Midtrans Snap API
type MidtransSnapResponse struct {
	Token       string `json:"token"`
	RedirectURL string `json:"redirect_url"`
}
