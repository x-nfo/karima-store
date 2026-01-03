package services

import (
	"errors"
	"fmt"
	"time"

	"github.com/karima-store/internal/models"
	"github.com/karima-store/internal/repository"
)

type PricingService interface {
	CalculatePrice(req PriceCalculationRequest) (*PriceCalculationResponse, error)
	CalculateShippingCost(req ShippingCalculationRequest) (*ShippingCalculationResponse, error)
	CheckFreeShipping(orderAmount float64, regionCode string) (bool, error)
	CalculateOrderSummary(items []PriceCalculationRequest, shippingReq ShippingCalculationRequest, customerType CustomerType) (*OrderSummary, error)
	CalculateCouponDiscount(req CouponCalculationRequest) (float64, string, error)
	ApplyCouponToPriceCalculation(resp *PriceCalculationResponse, couponReq CouponCalculationRequest) error
}

type pricingService struct {
	productRepo      repository.ProductRepository
	variantRepo      repository.VariantRepository
	flashSaleRepo    repository.FlashSaleRepository
	couponRepo       repository.CouponRepository
	shippingZoneRepo repository.ShippingZoneRepository
	taxRate          float64 // Default tax rate (e.g., 0.11 for 11% VAT)
}

type CustomerType string

const (
	CustomerRetail   CustomerType = "retail"
	CustomerReseller CustomerType = "reseller"
)

type PriceCalculationRequest struct {
	ProductID    uint         `json:"product_id"`
	VariantID    *uint        `json:"variant_id,omitempty"`
	Quantity     int          `json:"quantity"`
	CustomerType CustomerType `json:"customer_type"`
}

type PriceCalculationResponse struct {
	BasePrice       float64 `json:"base_price"`
	FinalPrice      float64 `json:"final_price"`
	Discount        float64 `json:"discount"`
	DiscountType    string  `json:"discount_type"` // "none", "flash_sale", "reseller", "bulk", "coupon"
	OriginalPrice   float64 `json:"original_price"`
	Savings         float64 `json:"savings"`
	FlashSaleActive bool    `json:"flash_sale_active"`
	FlashSaleEnd    *string `json:"flash_sale_end,omitempty"`
	CouponApplied   bool    `json:"coupon_applied"`
	CouponCode      *string `json:"coupon_code,omitempty"`
	CouponDiscount  float64 `json:"coupon_discount"`
}

type CouponCalculationRequest struct {
	Code           string
	UserID         uint
	PurchaseAmount float64
	CustomerType   CustomerType
}

type ShippingCalculationRequest struct {
	Items        []ShippingItem `json:"items"`
	Destination  string         `json:"destination"`   // subdistrict_id for RajaOngkir
	ShippingType string         `json:"shipping_type"` // "jne", "tiki", "pos", etc.
}

type ShippingItem struct {
	ProductID  uint
	VariantID  *uint
	Quantity   int
	Weight     float64 // in kg
	Dimensions string  // LxWxH format
}

type ShippingCalculationResponse struct {
	TotalWeight   float64 `json:"total_weight"` // in kg
	ShippingCost  float64 `json:"shipping_cost"`
	ShippingType  string  `json:"shipping_type"`
	EstimatedDays int     `json:"estimated_days"`
}

type OrderSummary struct {
	Subtotal       float64 `json:"subtotal"`
	ShippingCost   float64 `json:"shipping_cost"`
	TotalWeight    float64 `json:"total_weight"`
	Total          float64 `json:"total"`
	ItemCount      int     `json:"item_count"`
	TotalDiscount  float64 `json:"total_discount"`
	CouponDiscount float64 `json:"coupon_discount"`
	TaxAmount      float64 `json:"tax_amount"`
	CouponApplied  bool    `json:"coupon_applied"`
	CouponCode     *string `json:"coupon_code,omitempty"`
}

func NewPricingService(
	productRepo repository.ProductRepository,
	variantRepo repository.VariantRepository,
	flashSaleRepo repository.FlashSaleRepository,
	couponRepo repository.CouponRepository,
	shippingZoneRepo repository.ShippingZoneRepository,
) PricingService {
	return &pricingService{
		productRepo:      productRepo,
		variantRepo:      variantRepo,
		flashSaleRepo:    flashSaleRepo,
		couponRepo:       couponRepo,
		shippingZoneRepo: shippingZoneRepo,
		taxRate:          0.11, // Default 11% VAT for Indonesia
	}
}

// CalculatePrice calculates the final price based on customer type, quantity, and active flash sales
func (s *pricingService) CalculatePrice(req PriceCalculationRequest) (*PriceCalculationResponse, error) {
	// Validate input
	if req.Quantity <= 0 {
		return nil, errors.New("quantity must be greater than 0")
	}

	// Get product
	product, err := s.productRepo.GetByID(req.ProductID)
	if err != nil {
		return nil, fmt.Errorf("product not found: %w", err)
	}

	var basePrice float64
	var variant *models.ProductVariant

	// Use variant price if variant ID is provided
	if req.VariantID != nil && *req.VariantID > 0 {
		variant, err = s.variantRepo.GetByID(*req.VariantID)
		if err != nil {
			return nil, fmt.Errorf("variant not found: %w", err)
		}
		if variant.ProductID != req.ProductID {
			return nil, errors.New("variant does not belong to the specified product")
		}
		basePrice = variant.Price
	} else {
		basePrice = product.Price
	}

	// Check for flash sale
	flashSalePrice, flashSaleEnd, isFlashSaleActive := s.checkFlashSale(product.ID, req.VariantID)

	// Calculate final price
	response := &PriceCalculationResponse{
		BasePrice:     basePrice,
		OriginalPrice: basePrice,
	}

	// Priority: Flash Sale > Reseller Tiering > Bulk Discount > Retail
	if isFlashSaleActive && flashSalePrice > 0 {
		// Flash sale price takes precedence
		response.FinalPrice = flashSalePrice
		response.Discount = basePrice - flashSalePrice
		response.DiscountType = "flash_sale"
		response.FlashSaleActive = true
		response.FlashSaleEnd = flashSaleEnd
	} else if req.CustomerType == CustomerReseller {
		// Apply reseller tiering
		resellerPrice, discount := s.calculateResellerPrice(basePrice, req.Quantity)
		response.FinalPrice = resellerPrice
		response.Discount = discount
		response.DiscountType = "reseller"
	} else {
		// Apply bulk discount for retail customers
		bulkPrice, discount := s.calculateBulkDiscount(basePrice, req.Quantity)
		response.FinalPrice = bulkPrice
		response.Discount = discount
		if discount > 0 {
			response.DiscountType = "bulk"
		} else {
			response.DiscountType = "none"
		}
	}

	// Calculate savings
	response.Savings = response.OriginalPrice - response.FinalPrice

	// Apply quantity multiplier
	response.FinalPrice = response.FinalPrice * float64(req.Quantity)
	response.BasePrice = response.BasePrice * float64(req.Quantity)
	response.Savings = response.Savings * float64(req.Quantity)

	return response, nil
}

// checkFlashSale checks if there's an active flash sale for the product/variant
func (s *pricingService) checkFlashSale(productID uint, variantID *uint) (float64, *string, bool) {
	flashSales, err := s.flashSaleRepo.GetActiveFlashSales()
	if err != nil || len(flashSales) == 0 {
		return 0, nil, false
	}

	now := time.Now()

	for _, fs := range flashSales {
		// Check if flash sale is currently active
		if fs.Status != models.FlashSaleActive {
			continue
		}
		if now.Before(fs.StartTime) || now.After(fs.EndTime) {
			continue
		}

		// Check if product is in flash sale
		for _, fsp := range fs.Products {
			if fsp.ID == productID {
				// Get flash sale product details
				flashSaleProducts, err := s.flashSaleRepo.GetFlashSaleProducts(fs.ID)
				if err != nil {
					continue
				}

				for _, fspDetail := range flashSaleProducts {
					if fspDetail.ProductID == productID {
						return fspDetail.FlashSalePrice, formatTime(fs.EndTime), true
					}
				}
			}
		}
	}

	return 0, nil, false
}

// calculateResellerPrice applies reseller tiering based on quantity
func (s *pricingService) calculateResellerPrice(basePrice float64, quantity int) (float64, float64) {
	var discountPercent float64

	// Reseller tiering based on quantity
	switch {
	case quantity >= 100:
		discountPercent = 0.30 // 30% discount for 100+ items
	case quantity >= 50:
		discountPercent = 0.25 // 25% discount for 50-99 items
	case quantity >= 20:
		discountPercent = 0.20 // 20% discount for 20-49 items
	case quantity >= 10:
		discountPercent = 0.15 // 15% discount for 10-19 items
	case quantity >= 5:
		discountPercent = 0.10 // 10% discount for 5-9 items
	default:
		discountPercent = 0.05 // 5% discount for 1-4 items
	}

	discountAmount := basePrice * discountPercent
	finalPrice := basePrice - discountAmount

	return finalPrice, discountAmount
}

// calculateBulkDiscount applies bulk discount for retail customers
func (s *pricingService) calculateBulkDiscount(basePrice float64, quantity int) (float64, float64) {
	var discountPercent float64

	// Bulk discount for retail customers (less aggressive than reseller)
	switch {
	case quantity >= 10:
		discountPercent = 0.10 // 10% discount for 10+ items
	case quantity >= 5:
		discountPercent = 0.05 // 5% discount for 5-9 items
	default:
		discountPercent = 0.00 // No discount for less than 5 items
	}

	discountAmount := basePrice * discountPercent
	finalPrice := basePrice - discountAmount

	return finalPrice, discountAmount
}

// CalculateShippingCost calculates shipping cost based on weight and destination
func (s *pricingService) CalculateShippingCost(req ShippingCalculationRequest) (*ShippingCalculationResponse, error) {
	if len(req.Items) == 0 {
		return nil, errors.New("no items provided")
	}

	// Calculate total weight
	totalWeight := s.calculateTotalWeight(req.Items)

	// Get shipping zone for destination (if region is provided)
	var shippingZone *models.ShippingZone
	if req.Destination != "" {
		zone, err := s.shippingZoneRepo.GetByRegion(req.Destination)
		if err == nil {
			shippingZone = zone
		}
	}

	// Calculate shipping cost using zone or default rates
	shippingCost := s.calculateShippingCostWithZone(totalWeight, req.ShippingType, shippingZone)

	// Estimate delivery days based on shipping type
	estimatedDays := s.estimateDeliveryDays(req.ShippingType)

	return &ShippingCalculationResponse{
		TotalWeight:   totalWeight,
		ShippingCost:  shippingCost,
		ShippingType:  req.ShippingType,
		EstimatedDays: estimatedDays,
	}, nil
}

// calculateTotalWeight calculates the total weight of all items
func (s *pricingService) calculateTotalWeight(items []ShippingItem) float64 {
	var totalWeight float64

	for _, item := range items {
		totalWeight += item.Weight * float64(item.Quantity)
	}

	return totalWeight
}

// calculateShippingCostWithZone calculates shipping cost using zone-specific rates
func (s *pricingService) calculateShippingCostWithZone(weight float64, shippingType string, zone *models.ShippingZone) float64 {
	var baseCostPerKg float64
	var minCost float64
	var handlingFee float64

	// Use zone rates if available, otherwise use default rates
	if zone != nil {
		switch shippingType {
		case "jne":
			baseCostPerKg = zone.JNEBaseRate
		case "tiki":
			baseCostPerKg = zone.TIKIBaseRate
		case "pos":
			baseCostPerKg = zone.POSBaseRate
		case "sicepat":
			baseCostPerKg = zone.SiCepatBaseRate
		default:
			baseCostPerKg = zone.JNEBaseRate // Default to JNE rate
		}
		minCost = zone.MinimumCost
		handlingFee = zone.HandlingFee
	} else {
		// Default rates
		switch shippingType {
		case "jne":
			baseCostPerKg = 15000 // IDR 15,000 per kg
		case "tiki":
			baseCostPerKg = 16000 // IDR 16,000 per kg
		case "pos":
			baseCostPerKg = 14000 // IDR 14,000 per kg
		case "sicepat":
			baseCostPerKg = 13000 // IDR 13,000 per kg
		default:
			baseCostPerKg = 15000 // Default: IDR 15,000 per kg
		}
		minCost = 9000.0 // IDR 9,000 minimum
		handlingFee = 0.0
	}

	// Calculate cost
	cost := weight * baseCostPerKg

	// Apply minimum cost
	if cost < minCost {
		cost = minCost
	}

	// Add handling fee
	cost += handlingFee

	return cost
}

// CheckFreeShipping checks if shipping is free based on order amount and zone
func (s *pricingService) CheckFreeShipping(orderAmount float64, regionCode string) (bool, error) {
	if regionCode == "" {
		return false, nil
	}

	zone, err := s.shippingZoneRepo.GetByRegion(regionCode)
	if err != nil {
		return false, nil
	}

	if zone.FreeShippingEnabled && orderAmount >= zone.FreeShippingThreshold {
		return true, nil
	}

	return false, nil
}

// estimateDeliveryDays estimates delivery days based on shipping type
func (s *pricingService) estimateDeliveryDays(shippingType string) int {
	switch shippingType {
	case "jne":
		return 2 // 2 days
	case "tiki":
		return 2 // 2 days
	case "pos":
		return 3 // 3 days
	case "sicepat":
		return 1 // 1 day
	default:
		return 2 // Default: 2 days
	}
}

// CalculateOrderSummary calculates the complete order summary including pricing and shipping
func (s *pricingService) CalculateOrderSummary(
	items []PriceCalculationRequest,
	shippingReq ShippingCalculationRequest,
	customerType CustomerType,
) (*OrderSummary, error) {
	if len(items) == 0 {
		return nil, errors.New("no items provided")
	}

	var subtotal float64
	var totalDiscount float64
	var itemCount int

	// Calculate price for each item
	for _, item := range items {
		itemCount += item.Quantity
		priceResp, err := s.CalculatePrice(item)
		if err != nil {
			return nil, fmt.Errorf("error calculating price for item: %w", err)
		}

		subtotal += priceResp.BasePrice
		totalDiscount += priceResp.Savings
	}

	// Calculate shipping cost
	shippingResp, err := s.CalculateShippingCost(shippingReq)
	if err != nil {
		return nil, fmt.Errorf("error calculating shipping cost: %w", err)
	}

	// Calculate tax amount (11% VAT by default)
	taxAmount := (subtotal - totalDiscount) * s.taxRate

	// Calculate total
	total := subtotal - totalDiscount + shippingResp.ShippingCost + taxAmount

	return &OrderSummary{
		Subtotal:      subtotal,
		ShippingCost:  shippingResp.ShippingCost,
		TotalWeight:   shippingResp.TotalWeight,
		Total:         total,
		ItemCount:     itemCount,
		TotalDiscount: totalDiscount,
		TaxAmount:     taxAmount,
	}, nil
}

// CalculateCouponDiscount calculates discount for a coupon code
func (s *pricingService) CalculateCouponDiscount(req CouponCalculationRequest) (float64, string, error) {
	if req.Code == "" {
		return 0, "", nil
	}

	coupon, err := s.couponRepo.ValidateCoupon(
		req.Code,
		req.UserID,
		req.PurchaseAmount,
		string(req.CustomerType),
	)
	if err != nil {
		return 0, "", errors.New("invalid or expired coupon code")
	}

	var discountAmount float64

	switch coupon.Type {
	case models.CouponTypePercentage:
		discountAmount = req.PurchaseAmount * (coupon.DiscountValue / 100)
		// Apply max discount limit if set
		if coupon.MaxDiscount > 0 && discountAmount > coupon.MaxDiscount {
			discountAmount = coupon.MaxDiscount
		}
	case models.CouponTypeFixed:
		discountAmount = coupon.DiscountValue
		// Ensure discount doesn't exceed purchase amount
		if discountAmount > req.PurchaseAmount {
			discountAmount = req.PurchaseAmount
		}
	}

	return discountAmount, coupon.Name, nil
}

// ApplyCouponToPriceCalculation applies coupon discount to price calculation
func (s *pricingService) ApplyCouponToPriceCalculation(
	resp *PriceCalculationResponse,
	couponReq CouponCalculationRequest,
) error {
	if couponReq.Code == "" {
		return nil
	}

	discountAmount, _, err := s.CalculateCouponDiscount(couponReq)
	if err != nil {
		return err
	}

	resp.FinalPrice = resp.FinalPrice - discountAmount
	resp.Discount = resp.Discount + discountAmount
	resp.Savings = resp.Savings + discountAmount
	resp.DiscountType = "coupon"
	resp.CouponApplied = true
	resp.CouponCode = &couponReq.Code
	resp.CouponDiscount = discountAmount

	return nil
}

// formatTime formats time to ISO 8601 string
func formatTime(t time.Time) *string {
	formatted := t.Format(time.RFC3339)
	return &formatted
}
