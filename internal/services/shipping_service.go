package services

import (
	"fmt"

	"github.com/karima-store/internal/models"
	"github.com/karima-store/internal/rajaongkir"
)

// ShippingOption represents a simplified shipping option for API responses
type ShippingOption struct {
	Courier     string  `json:"courier"`
	Service     string  `json:"service"`
	Description string  `json:"description"`
	Cost        float64 `json:"cost"`
	ETD         string  `json:"etd"`
}

// ShippingService handles shipping-related operations
type ShippingService interface {
	GetAllProvinces() ([]models.RajaOngkirProvince, error)
	GetProvinceByID(provinceID string) (*models.RajaOngkirProvince, error)
	GetAllCities() ([]models.RajaOngkirCity, error)
	GetCitiesByProvince(provinceID string) ([]models.RajaOngkirCity, error)
	GetCityByID(cityID string) (*models.RajaOngkirCity, error)
	GetSubdistricts(cityID string) ([]models.RajaOngkirSubdistrict, error)
	GetSubdistrictByID(subdistrictID string) (*models.RajaOngkirSubdistrict, error)
	CalculateShippingCost(req models.RajaOngkirCostRequest) (*models.RajaOngkirResponse, error)
	CalculateShippingCostsForMultipleCouriers(
		origin string,
		originType string,
		destination string,
		destinationType string,
		weight int,
		couriers []string,
	) ([]models.RajaOngkirCourier, error)
	GetShippingOptions(origin, destination string, weight int) ([]models.RajaOngkirCourier, error)
	GetCheapestShippingOption(origin, destination string, weight int) (*ShippingOption, error)
	GetFastestShippingOption(origin, destination string, weight int) (*ShippingOption, error)
}

type shippingService struct {
	rajaongkirClient *rajaongkir.Client
}

// NewShippingService creates a new shipping service
func NewShippingService(rajaongkirClient *rajaongkir.Client) ShippingService {
	return &shippingService{
		rajaongkirClient: rajaongkirClient,
	}
}

// GetAllProvinces retrieves all provinces from RajaOngkir
func (s *shippingService) GetAllProvinces() ([]models.RajaOngkirProvince, error) {
	return s.rajaongkirClient.GetProvinces()
}

// GetProvinceByID retrieves a specific province by ID
func (s *shippingService) GetProvinceByID(provinceID string) (*models.RajaOngkirProvince, error) {
	return s.rajaongkirClient.GetProvinceByID(provinceID)
}

// GetAllCities retrieves all cities from RajaOngkir
func (s *shippingService) GetAllCities() ([]models.RajaOngkirCity, error) {
	return s.rajaongkirClient.GetCities()
}

// GetCitiesByProvince retrieves cities for a specific province
func (s *shippingService) GetCitiesByProvince(provinceID string) ([]models.RajaOngkirCity, error) {
	return s.rajaongkirClient.GetCitiesByProvince(provinceID)
}

// GetCityByID retrieves a specific city by ID
func (s *shippingService) GetCityByID(cityID string) (*models.RajaOngkirCity, error) {
	return s.rajaongkirClient.GetCityByID(cityID)
}

// GetSubdistricts retrieves all subdistricts for a specific city
func (s *shippingService) GetSubdistricts(cityID string) ([]models.RajaOngkirSubdistrict, error) {
	return s.rajaongkirClient.GetSubdistricts(cityID)
}

// GetSubdistrictByID retrieves a specific subdistrict by ID
func (s *shippingService) GetSubdistrictByID(subdistrictID string) (*models.RajaOngkirSubdistrict, error) {
	return s.rajaongkirClient.GetSubdistrictByID(subdistrictID)
}

// CalculateShippingCost calculates shipping cost between origin and destination
func (s *shippingService) CalculateShippingCost(req models.RajaOngkirCostRequest) (*models.RajaOngkirResponse, error) {
	// Validate request
	if req.Origin == "" {
		return nil, fmt.Errorf("origin is required")
	}
	if req.Destination == "" {
		return nil, fmt.Errorf("destination is required")
	}
	if req.Weight <= 0 {
		return nil, fmt.Errorf("weight must be greater than 0")
	}
	if req.Courier == "" {
		return nil, fmt.Errorf("courier is required")
	}

	// Set default origin type if not provided
	if req.OriginType == "" {
		req.OriginType = "city"
	}

	// Set default destination type if not provided
	if req.DestinationType == "" {
		req.DestinationType = "city"
	}

	return s.rajaongkirClient.GetShippingCost(req)
}

// CalculateShippingCostsForMultipleCouriers calculates shipping costs for multiple couriers
func (s *shippingService) CalculateShippingCostsForMultipleCouriers(
	origin string,
	originType string,
	destination string,
	destinationType string,
	weight int,
	couriers []string,
) ([]models.RajaOngkirCourier, error) {
	// Validate parameters
	if origin == "" {
		return nil, fmt.Errorf("origin is required")
	}
	if destination == "" {
		return nil, fmt.Errorf("destination is required")
	}
	if weight <= 0 {
		return nil, fmt.Errorf("weight must be greater than 0")
	}
	if len(couriers) == 0 {
		return nil, fmt.Errorf("at least one courier must be specified")
	}

	// Set default origin type if not provided
	if originType == "" {
		originType = "city"
	}

	// Set default destination type if not provided
	if destinationType == "" {
		destinationType = "city"
	}

	return s.rajaongkirClient.GetShippingCostsForMultipleCouriers(
		origin,
		originType,
		destination,
		destinationType,
		weight,
		couriers,
	)
}

// GetShippingOptions retrieves all available shipping options for a given route
func (s *shippingService) GetShippingOptions(origin, destination string, weight int) ([]models.RajaOngkirCourier, error) {
	// Validate parameters
	if origin == "" {
		return nil, fmt.Errorf("origin is required")
	}
	if destination == "" {
		return nil, fmt.Errorf("destination is required")
	}
	if weight <= 0 {
		return nil, fmt.Errorf("weight must be greater than 0")
	}

	// List of supported couriers
	couriers := []string{"jne", "tiki", "pos", "sicepat"}

	return s.CalculateShippingCostsForMultipleCouriers(
		origin,
		"city",
		destination,
		"city",
		weight,
		couriers,
	)
}

// GetCheapestShippingOption returns the cheapest shipping option
func (s *shippingService) GetCheapestShippingOption(origin, destination string, weight int) (*ShippingOption, error) {
	couriers, err := s.GetShippingOptions(origin, destination, weight)
	if err != nil {
		return nil, err
	}

	var cheapest *ShippingOption
	minCost := float64(999999999)

	for _, courier := range couriers {
		for _, service := range courier.Costs {
			for _, cost := range service.Costs {
				if cost.Value < minCost {
					minCost = cost.Value
					cheapest = &ShippingOption{
						Courier:     courier.Name,
						Service:     service.Service,
						Description: service.Description,
						Cost:        cost.Value,
						ETD:         cost.ETD,
					}
				}
			}
		}
	}

	if cheapest == nil {
		return nil, fmt.Errorf("no shipping options available")
	}

	return cheapest, nil
}

// GetFastestShippingOption returns the fastest shipping option
func (s *shippingService) GetFastestShippingOption(origin, destination string, weight int) (*ShippingOption, error) {
	couriers, err := s.GetShippingOptions(origin, destination, weight)
	if err != nil {
		return nil, err
	}

	var fastest *ShippingOption
	minDays := 999

	for _, courier := range couriers {
		for _, service := range courier.Costs {
			for _, cost := range service.Costs {
				// Parse ETD to get estimated days (simplified)
				days := parseETD(cost.ETD)
				if days < minDays && days > 0 {
					minDays = days
					fastest = &ShippingOption{
						Courier:     courier.Name,
						Service:     service.Service,
						Description: service.Description,
						Cost:        cost.Value,
						ETD:         cost.ETD,
					}
				}
			}
		}
	}

	if fastest == nil {
		return nil, fmt.Errorf("no shipping options available")
	}

	return fastest, nil
}

// parseETD parses the ETD string to extract the number of days
// This is a simplified implementation
func parseETD(etd string) int {
	// ETD format examples: "1-2 hari", "2-3 hari", "3 hari", "1 hari"
	// This is a simplified parser - in production, use a more robust solution
	if etd == "" {
		return 999
	}

	// Simple extraction of first number
	var firstNum int
	fmt.Sscanf(etd, "%d", &firstNum)

	if firstNum == 0 {
		return 999
	}

	return firstNum
}
