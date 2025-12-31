package rajaongkir

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/karima-store/internal/models"
)

const (
	DefaultTimeout = 30 * time.Second
)

// Client represents the RajaOngkir API client
type Client struct {
	apiKey     string
	baseURL    string
	httpClient *http.Client
}

// NewClient creates a new RajaOngkir API client
func NewClient(apiKey, baseURL string) *Client {
	return &Client{
		apiKey:  apiKey,
		baseURL: baseURL,
		httpClient: &http.Client{
			Timeout: DefaultTimeout,
		},
	}
}

// makeRequest makes an HTTP request to the RajaOngkir API
func (c *Client) makeRequest(method, endpoint string, body interface{}) ([]byte, error) {
	var reqBody io.Reader
	if body != nil {
		jsonBody, err := json.Marshal(body)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal request body: %w", err)
		}
		reqBody = bytes.NewBuffer(jsonBody)
	}

	url := c.baseURL + endpoint
	req, err := http.NewRequest(method, url, reqBody)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	// Set headers (Komerce uses x-api-key header)
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
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API request failed with status %d: %s", resp.StatusCode, string(respBody))
	}

	return respBody, nil
}

// GetProvinces retrieves all provinces from RajaOngkir
func (c *Client) GetProvinces() ([]models.RajaOngkirProvince, error) {
	respBody, err := c.makeRequest("GET", "/province", nil)
	if err != nil {
		return nil, err
	}

	var response models.RajaOngkirProvincesResponse
	if err := json.Unmarshal(respBody, &response); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	if response.RajaOngkir.Status.Code != 200 {
		return nil, fmt.Errorf("API returned error: %s", response.RajaOngkir.Status.Description)
	}

	return response.RajaOngkir.Results, nil
}

// GetProvinceByID retrieves a specific province by ID
func (c *Client) GetProvinceByID(provinceID string) (*models.RajaOngkirProvince, error) {
	respBody, err := c.makeRequest("GET", fmt.Sprintf("/province?id=%s", provinceID), nil)
	if err != nil {
		return nil, err
	}

	var response models.RajaOngkirProvincesResponse
	if err := json.Unmarshal(respBody, &response); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	if response.RajaOngkir.Status.Code != 200 {
		return nil, fmt.Errorf("API returned error: %s", response.RajaOngkir.Status.Description)
	}

	if len(response.RajaOngkir.Results) == 0 {
		return nil, fmt.Errorf("province not found")
	}

	return &response.RajaOngkir.Results[0], nil
}

// GetCities retrieves all cities from RajaOngkir
func (c *Client) GetCities() ([]models.RajaOngkirCity, error) {
	respBody, err := c.makeRequest("GET", "/city", nil)
	if err != nil {
		return nil, err
	}

	var response models.RajaOngkirCitiesResponse
	if err := json.Unmarshal(respBody, &response); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	if response.RajaOngkir.Status.Code != 200 {
		return nil, fmt.Errorf("API returned error: %s", response.RajaOngkir.Status.Description)
	}

	return response.RajaOngkir.Results, nil
}

// GetCitiesByProvince retrieves cities for a specific province
func (c *Client) GetCitiesByProvince(provinceID string) ([]models.RajaOngkirCity, error) {
	respBody, err := c.makeRequest("GET", fmt.Sprintf("/city?province=%s", provinceID), nil)
	if err != nil {
		return nil, err
	}

	var response models.RajaOngkirCitiesResponse
	if err := json.Unmarshal(respBody, &response); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	if response.RajaOngkir.Status.Code != 200 {
		return nil, fmt.Errorf("API returned error: %s", response.RajaOngkir.Status.Description)
	}

	return response.RajaOngkir.Results, nil
}

// GetCityByID retrieves a specific city by ID
func (c *Client) GetCityByID(cityID string) (*models.RajaOngkirCity, error) {
	respBody, err := c.makeRequest("GET", fmt.Sprintf("/city?id=%s", cityID), nil)
	if err != nil {
		return nil, err
	}

	var response models.RajaOngkirCitiesResponse
	if err := json.Unmarshal(respBody, &response); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	if response.RajaOngkir.Status.Code != 200 {
		return nil, fmt.Errorf("API returned error: %s", response.RajaOngkir.Status.Description)
	}

	if len(response.RajaOngkir.Results) == 0 {
		return nil, fmt.Errorf("city not found")
	}

	return &response.RajaOngkir.Results[0], nil
}

// GetSubdistricts retrieves all subdistricts for a specific city
func (c *Client) GetSubdistricts(cityID string) ([]models.RajaOngkirSubdistrict, error) {
	respBody, err := c.makeRequest("GET", fmt.Sprintf("/subdistrict?city=%s", cityID), nil)
	if err != nil {
		return nil, err
	}

	var response models.RajaOngkirSubdistrictsResponse
	if err := json.Unmarshal(respBody, &response); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	if response.RajaOngkir.Status.Code != 200 {
		return nil, fmt.Errorf("API returned error: %s", response.RajaOngkir.Status.Description)
	}

	return response.RajaOngkir.Results, nil
}

// GetSubdistrictByID retrieves a specific subdistrict by ID
func (c *Client) GetSubdistrictByID(subdistrictID string) (*models.RajaOngkirSubdistrict, error) {
	respBody, err := c.makeRequest("GET", fmt.Sprintf("/subdistrict?id=%s", subdistrictID), nil)
	if err != nil {
		return nil, err
	}

	var response models.RajaOngkirSubdistrictsResponse
	if err := json.Unmarshal(respBody, &response); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	if response.RajaOngkir.Status.Code != 200 {
		return nil, fmt.Errorf("API returned error: %s", response.RajaOngkir.Status.Description)
	}

	if len(response.RajaOngkir.Results) == 0 {
		return nil, fmt.Errorf("subdistrict not found")
	}

	return &response.RajaOngkir.Results[0], nil
}

// GetShippingCost calculates shipping cost between origin and destination
func (c *Client) GetShippingCost(req models.RajaOngkirCostRequest) (*models.RajaOngkirResponse, error) {
	respBody, err := c.makeRequest("POST", "/cost", req)
	if err != nil {
		return nil, err
	}

	var response models.RajaOngkirResponse
	if err := json.Unmarshal(respBody, &response); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	if response.RajaOngkir.Status.Code != 200 {
		return nil, fmt.Errorf("API returned error: %s", response.RajaOngkir.Status.Description)
	}

	return &response, nil
}

// GetShippingCostsForMultipleCouriers calculates shipping costs for multiple couriers
func (c *Client) GetShippingCostsForMultipleCouriers(
	origin string,
	originType string,
	destination string,
	destinationType string,
	weight int,
	couriers []string,
) ([]models.RajaOngkirCourier, error) {
	var allCouriers []models.RajaOngkirCourier

	for _, courier := range couriers {
		req := models.RajaOngkirCostRequest{
			Origin:          origin,
			OriginType:      originType,
			Destination:     destination,
			DestinationType: destinationType,
			Weight:          weight,
			Courier:         courier,
		}

		resp, err := c.GetShippingCost(req)
		if err != nil {
			// Log error but continue with other couriers
			continue
		}

		allCouriers = append(allCouriers, resp.RajaOngkir.Results...)
	}

	return allCouriers, nil
}

// GetInternationalShippingCost calculates international shipping cost (for PRO account)
func (c *Client) GetInternationalShippingCost(origin string, destination string, weight int, courier string) (*models.RajaOngkirResponse, error) {
	req := models.RajaOngkirCostRequest{
		Origin:      origin,
		Destination: destination,
		Weight:      weight,
		Courier:     courier,
	}

	respBody, err := c.makeRequest("POST", "/v2/internationalCost", req)
	if err != nil {
		return nil, err
	}

	var response models.RajaOngkirResponse
	if err := json.Unmarshal(respBody, &response); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	if response.RajaOngkir.Status.Code != 200 {
		return nil, fmt.Errorf("API returned error: %s", response.RajaOngkir.Status.Description)
	}

	return &response, nil
}

// GetWaybill tracks a shipment by waybill number (for PRO account)
func (c *Client) GetWaybill(waybillNumber, courier string) ([]byte, error) {
	reqBody := map[string]string{
		"waybill": waybillNumber,
		"courier": courier,
	}

	respBody, err := c.makeRequest("POST", "/waybill", reqBody)
	if err != nil {
		return nil, err
	}

	return respBody, nil
}
