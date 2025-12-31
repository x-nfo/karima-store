package models

// RajaOngkir Province Model
type RajaOngkirProvince struct {
	ProvinceID string `json:"province_id"`
	Province   string `json:"province"`
}

// RajaOngkir City Model
type RajaOngkirCity struct {
	CityID     string `json:"city_id"`
	ProvinceID string `json:"province_id"`
	Province   string `json:"province"`
	Type       string `json:"type"`
	CityName   string `json:"city_name"`
	PostalCode string `json:"postal_code"`
}

// RajaOngkir Subdistrict Model
type RajaOngkirSubdistrict struct {
	SubdistrictID string `json:"subdistrict_id"`
	ProvinceID    string `json:"province_id"`
	Province      string `json:"province"`
	CityID        string `json:"city_id"`
	City          string `json:"city"`
	Type          string `json:"type"`
	SubdistrictName string `json:"subdistrict_name"`
}

// RajaOngkir Cost Request
type RajaOngkirCostRequest struct {
	Origin      string `json:"origin"`      // City ID
	OriginType  string `json:"originType"`  // "city" or "subdistrict"
	Destination string `json:"destination"` // City ID
	DestinationType string `json:"destinationType"` // "city" or "subdistrict"
	Weight      int    `json:"weight"`      // in grams
	Courier     string `json:"courier"`     // "jne", "tiki", "pos", "sicepat", etc.
}

// RajaOngkir Cost Detail
type RajaOngkirCostDetail struct {
	Description string  `json:"description"`
	Value       float64 `json:"value"` // in IDR
	ETD         string  `json:"etd"`   // Estimated time of delivery
	Note        string  `json:"note"`
}

// RajaOngkir Service
type RajaOngkirService struct {
	Service     string               `json:"service"`
	Description string               `json:"description"`
	Costs       []RajaOngkirCostDetail `json:"costs"`
}

// RajaOngkir Courier
type RajaOngkirCourier struct {
	Code        string               `json:"code"` // "jne", "tiki", "pos", etc.
	Name        string               `json:"name"`
	Costs       []RajaOngkirService  `json:"costs"`
}

// RajaOngkir API Response Wrapper
type RajaOngkirResponse struct {
	RajaOngkir struct {
		Query struct {
			Origin string `json:"origin"`
			Destination string `json:"destination"`
			Weight int `json:"weight"`
			Courier string `json:"courier"`
		} `json:"query"`
		Status struct {
			Code        int    `json:"code"`
			Description string `json:"description"`
		} `json:"status"`
		OriginDetails      interface{} `json:"originDetails"`
		DestinationDetails interface{} `json:"destinationDetails"`
		Results            []RajaOngkirCourier `json:"results"`
	} `json:"rajaongkir"`
}

// RajaOngkir Provinces Response
type RajaOngkirProvincesResponse struct {
	RajaOngkir struct {
		Query struct {
		} `json:"query"`
		Status struct {
			Code        int    `json:"code"`
			Description string `json:"description"`
		} `json:"status"`
		Results []RajaOngkirProvince `json:"results"`
	} `json:"rajaongkir"`
}

// RajaOngkir Cities Response
type RajaOngkirCitiesResponse struct {
	RajaOngkir struct {
		Query struct {
			Province string `json:"province"`
		} `json:"query"`
		Status struct {
			Code        int    `json:"code"`
			Description string `json:"description"`
		} `json:"status"`
		Results []RajaOngkirCity `json:"results"`
	} `json:"rajaongkir"`
}

// RajaOngkir Subdistricts Response
type RajaOngkirSubdistrictsResponse struct {
	RajaOngkir struct {
		Query struct {
			City string `json:"city"`
		} `json:"query"`
		Status struct {
			Code        int    `json:"code"`
			Description string `json:"description"`
		} `json:"status"`
		Results []RajaOngkirSubdistrict `json:"results"`
	} `json:"rajaongkir"`
}

// Shipping Cost Response (Internal use)
type ShippingCostResponse struct {
	Origin      string `json:"origin"`
	Destination string `json:"destination"`
	Weight      int    `json:"weight"`
	Couriers    []RajaOngkirCourier `json:"couriers"`
}
