package services

import (
	"encoding/json"
	"fmt"

	"github.com/karima-store/internal/komerce"
	"github.com/karima-store/internal/models"
)

// KomerceService handles Komerce API operations
type KomerceService interface {
	SearchDestination(keyword string) ([]models.KomerceDestination, error)
	CalculateShippingCost(shipperDestID, receiverDestID string, weight float64, itemValue int, cod string) (*models.KomerceCalculateResponse, error)
	CreateOrder(order models.KomerceCreateOrderRequest) (*models.KomerceCreateOrderResponse, error)
	GetOrderDetail(orderNo string) (*models.KomerceOrderDetailResponse, error)
	CancelOrder(orderNo string) error
	RequestPickup(vehicle, time, date string, orders []string) error
	PrintLabel(orderNo, page string) (string, error)
	TrackOrder(shipping, airwayBill string) (*models.KomerceTrackingResponse, error)
}

type komerceService struct {
	komerceClient *komerce.Client
}

// NewKomerceService creates a new Komerce service
func NewKomerceService(komerceClient *komerce.Client) KomerceService {
	return &komerceService{
		komerceClient: komerceClient,
	}
}

// SearchDestination searches for destination by keyword
func (s *komerceService) SearchDestination(keyword string) ([]models.KomerceDestination, error) {
	if keyword == "" {
		return nil, fmt.Errorf("keyword is required")
	}

	respBody, err := s.komerceClient.SearchDestination(keyword)
	if err != nil {
		return nil, err
	}

	var response models.KomerceResponse
	if err := json.Unmarshal(respBody, &response); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	if response.Meta.Status != "success" {
		return nil, fmt.Errorf("API returned error: %s", response.Meta.Message)
	}

	// Parse destinations from response
	destinations := []models.KomerceDestination{}
	if response.Data != nil {
		if data, ok := response.Data.([]interface{}); ok {
			for _, item := range data {
				if itemMap, ok := item.(map[string]interface{}); ok {
					dest := models.KomerceDestination{
						ID:              getString(itemMap, "id"),
						Label:           getString(itemMap, "label"),
						SubdistrictName: getString(itemMap, "subdistrict_name"),
						DistrictName:   getString(itemMap, "district_name"),
						CityName:       getString(itemMap, "city_name"),
						ZipCode:        getString(itemMap, "zip_code"),
					}
					destinations = append(destinations, dest)
				}
			}
		}
	}

	return destinations, nil
}

// CalculateShippingCost calculates shipping cost
func (s *komerceService) CalculateShippingCost(shipperDestID, receiverDestID string, weight float64, itemValue int, cod string) (*models.KomerceCalculateResponse, error) {
	// Validate required fields
	if shipperDestID == "" {
		return nil, fmt.Errorf("shipper_destination_id is required")
	}
	if receiverDestID == "" {
		return nil, fmt.Errorf("receiver_destination_id is required")
	}
	if weight <= 0 {
		return nil, fmt.Errorf("weight must be greater than 0")
	}
	if itemValue < 0 {
		return nil, fmt.Errorf("item_value must be greater than or equal to 0")
	}

	respBody, err := s.komerceClient.CalculateShippingCost(shipperDestID, receiverDestID, weight, itemValue, cod)
	if err != nil {
		return nil, err
	}

	var response models.KomerceCalculateResponse
	if err := json.Unmarshal(respBody, &response); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	if response.Meta.Status != "success" {
		return nil, fmt.Errorf("API returned error: %s", response.Meta.Message)
	}

	return &response, nil
}

// CreateOrder creates a new order
func (s *komerceService) CreateOrder(order models.KomerceCreateOrderRequest) (*models.KomerceCreateOrderResponse, error) {
	// Validate required fields
	if order.OrderDate == "" {
		return nil, fmt.Errorf("order_date is required")
	}
	if order.BrandName == "" {
		return nil, fmt.Errorf("brand_name is required")
	}
	if order.ShipperName == "" {
		return nil, fmt.Errorf("shipper_name is required")
	}
	if order.ShipperPhone == "" {
		return nil, fmt.Errorf("shipper_phone is required")
	}
	if order.ShipperDestinationID == "" {
		return nil, fmt.Errorf("shipper_destination_id is required")
	}
	if order.ShipperAddress == "" {
		return nil, fmt.Errorf("shipper_address is required")
	}
	if order.ShipperEmail == "" {
		return nil, fmt.Errorf("shipper_email is required")
	}
	if order.ReceiverName == "" {
		return nil, fmt.Errorf("receiver_name is required")
	}
	if order.ReceiverPhone == "" {
		return nil, fmt.Errorf("receiver_phone is required")
	}
	if order.ReceiverDestinationID == "" {
		return nil, fmt.Errorf("receiver_destination_id is required")
	}
	if order.ReceiverAddress == "" {
		return nil, fmt.Errorf("receiver_address is required")
	}
	if order.Shipping == "" {
		return nil, fmt.Errorf("shipping is required")
	}
	if order.ShippingCost < 0 {
		return nil, fmt.Errorf("shipping_cost must be greater than or equal to 0")
	}
	if order.GrandTotal < 0 {
		return nil, fmt.Errorf("grand_total must be greater than or equal to 0")
	}

	// Validate COD value equals grand total for COD payment
	if order.PaymentMethod == "COD" && order.CODValue != order.GrandTotal {
		return nil, fmt.Errorf("cod_value must be equal to grand_total for COD payment")
	}

	respBody, err := s.komerceClient.CreateOrder(order)
	if err != nil {
		return nil, err
	}

	var response models.KomerceCreateOrderResponse
	if err := json.Unmarshal(respBody, &response); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	if response.Meta.Status != "success" {
		return nil, fmt.Errorf("API returned error: %s", response.Meta.Message)
	}

	return &response, nil
}

// GetOrderDetail gets order detail by order number
func (s *komerceService) GetOrderDetail(orderNo string) (*models.KomerceOrderDetailResponse, error) {
	if orderNo == "" {
		return nil, fmt.Errorf("order_no is required")
	}

	respBody, err := s.komerceClient.GetOrderDetail(orderNo)
	if err != nil {
		return nil, err
	}

	var response models.KomerceOrderDetailResponse
	if err := json.Unmarshal(respBody, &response); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	if response.Meta.Status != "success" {
		return nil, fmt.Errorf("API returned error: %s", response.Meta.Message)
	}

	return &response, nil
}

// CancelOrder cancels an order
func (s *komerceService) CancelOrder(orderNo string) error {
	if orderNo == "" {
		return fmt.Errorf("order_no is required")
	}

	respBody, err := s.komerceClient.CancelOrder(orderNo)
	if err != nil {
		return err
	}

	var response models.KomerceResponse
	if err := json.Unmarshal(respBody, &response); err != nil {
		return fmt.Errorf("failed to unmarshal response: %w", err)
	}

	if response.Meta.Status != "success" {
		return fmt.Errorf("API returned error: %s", response.Meta.Message)
	}

	return nil
}

// RequestPickup requests pickup for orders
func (s *komerceService) RequestPickup(vehicle, time, date string, orders []string) error {
	if vehicle == "" {
		return fmt.Errorf("pickup_vehicle is required")
	}
	if time == "" {
		return fmt.Errorf("pickup_time is required")
	}
	if date == "" {
		return fmt.Errorf("pickup_date is required")
	}
	if len(orders) == 0 {
		return fmt.Errorf("at least one order is required")
	}

	_, err := s.komerceClient.RequestPickup(vehicle, time, date, orders)
	if err != nil {
		return err
	}

	return nil
}

// PrintLabel generates print label
func (s *komerceService) PrintLabel(orderNo, page string) (string, error) {
	if orderNo == "" {
		return "", fmt.Errorf("order_no is required")
	}

	respBody, err := s.komerceClient.PrintLabel(orderNo, page)
	if err != nil {
		return "", err
	}

	var response models.KomerceResponse
	if err := json.Unmarshal(respBody, &response); err != nil {
		return "", fmt.Errorf("failed to unmarshal response: %w", err)
	}

	if response.Meta.Status != "success" {
		return "", fmt.Errorf("API returned error: %s", response.Meta.Message)
	}

	// Extract path from response
	if response.Data != nil {
		if dataMap, ok := response.Data.(map[string]interface{}); ok {
			if path, exists := dataMap["path"]; exists {
				if pathStr, ok := path.(string); ok {
					return pathStr, nil
				}
			}
		}
	}

	return "", fmt.Errorf("path not found in response")
}

// TrackOrder tracks order by shipping provider and airway bill
func (s *komerceService) TrackOrder(shipping, airwayBill string) (*models.KomerceTrackingResponse, error) {
	if shipping == "" {
		return nil, fmt.Errorf("shipping is required")
	}
	if airwayBill == "" {
		return nil, fmt.Errorf("airway_bill is required")
	}

	respBody, err := s.komerceClient.TrackOrder(shipping, airwayBill)
	if err != nil {
		return nil, err
	}

	var response models.KomerceTrackingResponse
	if err := json.Unmarshal(respBody, &response); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	if response.Meta.Status != "success" {
		return nil, fmt.Errorf("API returned error: %s", response.Meta.Message)
	}

	return &response, nil
}

// Helper function to safely get string from map
func getString(m map[string]interface{}, key string) string {
	if val, exists := m[key]; exists {
		if str, ok := val.(string); ok {
			return str
		}
	}
	return ""
}
