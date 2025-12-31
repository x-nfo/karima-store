package komerce

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"
)

const (
	DefaultTimeout = 30 * time.Second
)

// Client represents Komerce API client
type Client struct {
	apiKey     string
	baseURL    string
	httpClient *http.Client
}

// NewClient creates a new Komerce API client
func NewClient(apiKey, baseURL string) *Client {
	if baseURL == "" {
		baseURL = "https://api-sandbox.collaborator.komerce.id"
	}
	return &Client{
		apiKey:  apiKey,
		baseURL: baseURL,
		httpClient: &http.Client{
			Timeout: DefaultTimeout,
		},
	}
}

// makeRequest makes an HTTP request to Komerce API
func (c *Client) makeRequest(method, endpoint string, body interface{}, queryParams map[string]string) ([]byte, error) {
	var reqBody io.Reader
	if body != nil {
		jsonBody, err := json.Marshal(body)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal request body: %w", err)
		}
		reqBody = bytes.NewBuffer(jsonBody)
	}

	// Build URL with query parameters
	requestURL := c.baseURL + endpoint
	if len(queryParams) > 0 {
		values := url.Values{}
		for key, value := range queryParams {
			values.Add(key, value)
		}
		requestURL = requestURL + "?" + values.Encode()
	}

	req, err := http.NewRequest(method, requestURL, reqBody)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	// Set headers
	req.Header.Set("x-api-key", c.apiKey)
	req.Header.Set("Content-Type", "application/json")

	// Make request
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to make request: %w", err)
	}
	defer resp.Body.Close()

	// Read response
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	// Check for non-200 status codes
	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		return nil, fmt.Errorf("API request failed with status %d: %s", resp.StatusCode, string(respBody))
	}

	return respBody, nil
}

// SearchDestination searches for destination by postal code, village, sub-district, or district
func (c *Client) SearchDestination(keyword string) ([]byte, error) {
	queryParams := map[string]string{
		"keyword": keyword,
	}
	return c.makeRequest("GET", "/tariff/api/v1/destination/search", nil, queryParams)
}

// CalculateShippingCost calculates shipping cost
func (c *Client) CalculateShippingCost(shipperDestID, receiverDestID string, weight float64, itemValue int, cod string) ([]byte, error) {
	queryParams := map[string]string{
		"shipper_destination_id":  shipperDestID,
		"receiver_destination_id": receiverDestID,
		"weight":                   fmt.Sprintf("%.0f", weight),
		"item_value":               fmt.Sprintf("%d", itemValue),
		"cod":                      cod,
	}
	return c.makeRequest("GET", "/tariff/api/v1/calculate", nil, queryParams)
}

// CreateOrder creates a new order
func (c *Client) CreateOrder(order interface{}) ([]byte, error) {
	return c.makeRequest("POST", "/order/api/v1/orders/store", order, nil)
}

// GetOrderDetail gets order detail by order number
func (c *Client) GetOrderDetail(orderNo string) ([]byte, error) {
	queryParams := map[string]string{
		"order_no": orderNo,
	}
	return c.makeRequest("GET", "/order/api/v1/orders/detail", nil, queryParams)
}

// CancelOrder cancels an order
func (c *Client) CancelOrder(orderNo string) ([]byte, error) {
	reqBody := map[string]string{
		"order_no": orderNo,
	}
	return c.makeRequest("PUT", "/order/api/v1/orders/cancel", reqBody, nil)
}

// RequestPickup requests pickup for orders
func (c *Client) RequestPickup(vehicle, time, date string, orders []string) ([]byte, error) {
	reqBody := map[string]interface{}{
		"pickup_vehicle": vehicle,
		"pickup_time":    time,
		"pickup_date":    date,
		"orders":         orders,
	}
	return c.makeRequest("POST", "/order/api/v1/pickup/request", reqBody, nil)
}

// PrintLabel generates print label
func (c *Client) PrintLabel(orderNo, page string) ([]byte, error) {
	queryParams := map[string]string{
		"order_no": orderNo,
		"page":     page,
	}
	return c.makeRequest("POST", "/order/api/v1/orders/print-label", nil, queryParams)
}

// TrackOrder tracks order by shipping provider and airway bill
func (c *Client) TrackOrder(shipping, airwayBill string) ([]byte, error) {
	queryParams := map[string]string{
		"shipping":    shipping,
		"airway_bill": airwayBill,
	}
	return c.makeRequest("GET", "/order/api/v1/orders/history-airway-bill", nil, queryParams)
}
