package models

// Komerce Destination Search Response
type KomerceDestination struct {
	ID              string `json:"id"`
	Label           string `json:"label"`
	SubdistrictName string `json:"subdistrict_name"`
	DistrictName    string `json:"district_name"`
	CityName        string `json:"city_name"`
	ZipCode         string `json:"zip_code"`
}

// Komerce Calculate Shipping Request
type KomerceCalculateRequest struct {
	ShipperDestinationID  string  `json:"shipper_destination_id"`
	ReceiverDestinationID string  `json:"receiver_destination_id"`
	Weight                float64 `json:"weight"`
	ItemValue             int     `json:"item_value"`
	COD                   string  `json:"cod"`
}

// Komerce Shipping Courier
type KomerceCourier struct {
	Name    string                           `json:"name"`
	Regular map[string]KomerceShippingOption `json:"regular"`
	Cargo   map[string]KomerceShippingOption `json:"cargo"`
}

// Komerce Shipping Option from Calculate API
type KomerceShippingOption struct {
	ShippingName     string `json:"shipping_name"`
	ServiceName      string `json:"service_name"`
	Weight           int    `json:"weight"`
	IsCOD            bool   `json:"is_cod"`
	ShippingCost     int    `json:"shipping_cost"`
	ShippingCashback int    `json:"shipping_cashback"`
	ShippingCostNet  int    `json:"shipping_cost_net"`
	GrandTotal       int    `json:"grandtotal"`
	ServiceFee       int    `json:"service_fee"`
	NetIncome        int    `json:"net_income"`
	ETD              string `json:"etd"`
}

// Komerce Calculate Response Data
type KomerceCalculateData struct {
	CalculateReguler []KomerceShippingOption `json:"calculate_reguler"`
	CalculateCargo   []KomerceShippingOption `json:"calculate_cargo"`
}

// Komerce Calculate Response
type KomerceCalculateResponse struct {
	Meta KomerceMeta          `json:"meta"`
	Data KomerceCalculateData `json:"data"`
}

// Komerce Order Detail
type KomerceOrderDetail struct {
	ProductName        string `json:"product_name"`
	ProductVariantName string `json:"product_variant_name"`
	ProductPrice       int    `json:"product_price"`
	ProductWidth       int    `json:"product_width"`
	ProductHeight      int    `json:"product_height"`
	ProductLength      int    `json:"product_length"`
	ProductWeight      int    `json:"product_weight"`
	Qty                int    `json:"qty"`
	Subtotal           int    `json:"subtotal"`
}

// Komerce Create Order Request
type KomerceCreateOrderRequest struct {
	OrderDate             string               `json:"order_date"`
	BrandName             string               `json:"brand_name"`
	ShipperName           string               `json:"shipper_name"`
	ShipperPhone          string               `json:"shipper_phone"`
	ShipperDestinationID  string               `json:"shipper_destination_id"`
	ShipperAddress        string               `json:"shipper_address"`
	OriginPinPoint        string               `json:"origin_pin_point"`
	ShipperEmail          string               `json:"shipper_email"`
	ReceiverName          string               `json:"receiver_name"`
	ReceiverPhone         string               `json:"receiver_phone"`
	ReceiverDestinationID string               `json:"receiver_destination_id"`
	ReceiverAddress       string               `json:"receiver_address"`
	DestinationPinPoint   string               `json:"destination_pin_point"`
	Shipping              string               `json:"shipping"`
	ShippingType          string               `json:"shipping_type"`
	ShippingCost          int                  `json:"shipping_cost"`
	ShippingCashback      int                  `json:"shipping_cashback"`
	PaymentMethod         string               `json:"payment_method"`
	ServiceFee            int                  `json:"service_fee"`
	AdditionalCost        int                  `json:"additional_cost"`
	GrandTotal            int                  `json:"grand_total"`
	CODValue              int                  `json:"cod_value"`
	InsuranceValue        int                  `json:"insurance_value"`
	OrderDetails          []KomerceOrderDetail `json:"order_details"`
}

// Komerce Create Order Response
type KomerceCreateOrderResponse struct {
	Meta KomerceMeta      `json:"meta"`
	Data KomerceOrderData `json:"data"`
}

type KomerceOrderData struct {
	OrderID string `json:"order_id"`
	OrderNo string `json:"order_no"`
}

// Komerce Order Detail Response
type KomerceOrderDetailResponse struct {
	Meta KomerceMeta            `json:"meta"`
	Data KomerceFullOrderDetail `json:"data"`
}

type KomerceFullOrderDetail struct {
	OrderNo               string               `json:"order_no"`
	AWB                   string               `json:"awb"`
	OrderStatus           string               `json:"order_status"`
	OrderDate             string               `json:"order_date"`
	BrandName             string               `json:"brand_name"`
	ShipperName           string               `json:"shipper_name"`
	ShipperPhone          string               `json:"shipper_phone"`
	ShipperDestinationID  string               `json:"shipper_destination_id"`
	ShipperAddress        string               `json:"shipper_address"`
	ReceiverName          string               `json:"receiver_name"`
	ReceiverPhone         string               `json:"receiver_phone"`
	ReceiverDestinationID string               `json:"receiver_destination_id"`
	ReceiverAddress       string               `json:"receiver_address"`
	Shipping              string               `json:"shipping"`
	ShippingType          string               `json:"shipping_type"`
	PaymentMethod         string               `json:"payment_method"`
	ShippingCost          int                  `json:"shipping_cost"`
	ShippingCashback      int                  `json:"shipping_cashback"`
	ServiceFee            int                  `json:"service_fee"`
	AdditionalCost        int                  `json:"additional_cost"`
	GrandTotal            int                  `json:"grand_total"`
	CODValue              int                  `json:"cod_value"`
	InsuranceValue        int                  `json:"insurance_value"`
	Notes                 string               `json:"notes"`
	OriginPinPoint        string               `json:"origin_pin_point"`
	DestinationPinPoint   string               `json:"destination_pin_point"`
	BookingID             string               `json:"booking_id"`
	DriverName            string               `json:"driver_name"`
	DriverPhone           string               `json:"driver_phone"`
	CancelationReason     string               `json:"cancelation_reason"`
	LiveTrackingURL       string               `json:"live_tracking_url"`
	CommodityCode         string               `json:"commodity_code"`
	OrderDetails          []KomerceOrderDetail `json:"order_details"`
}

// Komerce Cancel Order Request
type KomerceCancelOrderRequest struct {
	OrderNo string `json:"order_no"`
}

// Komerce Pickup Request
type KomercePickupRequest struct {
	PickupVehicle string   `json:"pickup_vehicle"`
	PickupTime    string   `json:"pickup_time"`
	PickupDate    string   `json:"pickup_date"`
	Orders        []string `json:"orders"`
}

// Komerce Pickup Response
type KomercePickupResponse struct {
	Meta KomerceMeta         `json:"meta"`
	Data []KomercePickupData `json:"data"`
}

type KomercePickupData struct {
	Status  string `json:"status"`
	OrderNo string `json:"order_no"`
	AWB     string `json:"awb"`
}

// Komerce Tracking History
type KomerceTrackingHistory struct {
	Desc   string `json:"desc"`
	Date   string `json:"date"`
	Code   string `json:"code"`
	Status string `json:"status"`
}

// Komerce Tracking Response
type KomerceTrackingResponse struct {
	Meta KomerceMeta         `json:"meta"`
	Data KomerceTrackingData `json:"data"`
}

type KomerceTrackingData struct {
	AirwayBill string                   `json:"airway_bill"`
	LastStatus string                   `json:"last_status"`
	History    []KomerceTrackingHistory `json:"history"`
}

// Komerce Meta
type KomerceMeta struct {
	Message string `json:"message"`
	Code    int    `json:"code"`
	Status  string `json:"status"`
}

// Komerce API Response Wrapper
type KomerceResponse struct {
	Meta KomerceMeta `json:"meta"`
	Data interface{} `json:"data"`
}
