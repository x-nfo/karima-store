package services

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/karima-store/internal/komerce"
	"github.com/karima-store/internal/models"
	"github.com/stretchr/testify/assert"
)

func TestKomerceService_SearchDestination(t *testing.T) {
	// Mock Komerce API
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/tariff/api/v1/destination/search", r.URL.Path)
		assert.Equal(t, "jakarta", r.URL.Query().Get("keyword"))

		response := map[string]interface{}{
			"meta": map[string]interface{}{
				"status":  "success",
				"message": "Data found",
			},
			"data": []interface{}{
				map[string]interface{}{
					"id":               "1",
					"label":            "Jakarta Selatan",
					"subdistrict_name": "Kebayoran Baru",
				},
			},
		}
		json.NewEncoder(w).Encode(response)
	}))
	defer server.Close()

	client := komerce.NewClient("test-key", server.URL)
	service := NewKomerceService(client)

	destinations, err := service.SearchDestination("jakarta")
	assert.NoError(t, err)
	assert.Len(t, destinations, 1)
	assert.Equal(t, "Jakarta Selatan", destinations[0].Label)
}

func TestKomerceService_SearchDestination_Error(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer server.Close()

	client := komerce.NewClient("test-key", server.URL)
	service := NewKomerceService(client)

	_, err := service.SearchDestination("jakarta")
	assert.Error(t, err)
}

func TestKomerceService_CalculateShippingCost(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/tariff/api/v1/calculate", r.URL.Path)

		response := map[string]interface{}{
			"meta": map[string]interface{}{
				"status": "success",
			},
			"data": map[string]interface{}{
				"calculate_reguler": []interface{}{
					map[string]interface{}{
						"shipping_name": "JNE",
						"service_name":  "REG",
						"shipping_cost": 10000,
					},
				},
				"calculate_cargo": []interface{}{},
			},
		}
		json.NewEncoder(w).Encode(response)
	}))
	defer server.Close()

	client := komerce.NewClient("test-key", server.URL)
	service := NewKomerceService(client)

	resp, err := service.CalculateShippingCost("1", "2", 1.0, 50000, "false")
	assert.NoError(t, err)
	assert.NotNil(t, resp)
}

func TestKomerceService_CreateOrder(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/order/api/v1/orders/store", r.URL.Path)
		assert.Equal(t, "POST", r.Method)

		response := map[string]interface{}{
			"meta": map[string]interface{}{
				"status": "success",
			},
			"data": map[string]interface{}{
				"order_no": "ORD-123",
			},
		}
		json.NewEncoder(w).Encode(response)
	}))
	defer server.Close()

	client := komerce.NewClient("test-key", server.URL)
	service := NewKomerceService(client)

	order := models.KomerceCreateOrderRequest{
		OrderDate:             "2023-01-01",
		BrandName:             "Karima",
		ShipperName:           "Sender",
		ShipperPhone:          "081",
		ShipperDestinationID:  "1",
		ShipperAddress:        "Addr",
		ShipperEmail:          "sender@test.com",
		ReceiverName:          "Receiver",
		ReceiverPhone:         "082",
		ReceiverDestinationID: "2",
		ReceiverAddress:       "Addr2",
		Shipping:              "JNE",
		ShippingCost:          10000,
		GrandTotal:            100000,
	}

	resp, err := service.CreateOrder(order)
	assert.NoError(t, err)
	assert.NotNil(t, resp)
}

func TestKomerceService_CreateOrder_ValidationRequest(t *testing.T) {
	service := NewKomerceService(nil)

	// Missing required fields
	order := models.KomerceCreateOrderRequest{}
	_, err := service.CreateOrder(order)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "required")
}
