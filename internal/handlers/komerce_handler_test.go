package handlers

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/karima-store/internal/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockKomerceService is a mock implementation of KomerceService
type MockKomerceService struct {
	mock.Mock
}

func (m *MockKomerceService) SearchDestination(keyword string) ([]models.KomerceDestination, error) {
	args := m.Called(keyword)
	return args.Get(0).([]models.KomerceDestination), args.Error(1)
}

func (m *MockKomerceService) CalculateShippingCost(shipperDestID, receiverDestID string, weight float64, itemValue int, cod string) (*models.KomerceCalculateResponse, error) {
	args := m.Called(shipperDestID, receiverDestID, weight, itemValue, cod)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.KomerceCalculateResponse), args.Error(1)
}

func (m *MockKomerceService) CreateOrder(order models.KomerceCreateOrderRequest) (*models.KomerceCreateOrderResponse, error) {
	args := m.Called(order)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.KomerceCreateOrderResponse), args.Error(1)
}

func (m *MockKomerceService) GetOrderDetail(orderNo string) (*models.KomerceOrderDetailResponse, error) {
	args := m.Called(orderNo)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.KomerceOrderDetailResponse), args.Error(1)
}

func (m *MockKomerceService) CancelOrder(orderNo string) error {
	args := m.Called(orderNo)
	return args.Error(0)
}

func (m *MockKomerceService) RequestPickup(vehicle, time, date string, orders []string) error {
	args := m.Called(vehicle, time, date, orders)
	return args.Error(0)
}

func (m *MockKomerceService) PrintLabel(orderNo, page string) (string, error) {
	args := m.Called(orderNo, page)
	return args.String(0), args.Error(1)
}

func (m *MockKomerceService) TrackOrder(shipping, airwayBill string) (*models.KomerceTrackingResponse, error) {
	args := m.Called(shipping, airwayBill)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.KomerceTrackingResponse), args.Error(1)
}

func TestKomerceHandler_CalculateShippingCost(t *testing.T) {
	tests := []struct {
		name           string
		requestBody    interface{}
		setupMock      func(*MockKomerceService)
		expectedStatus int
		expectedError  bool
	}{
		{
			name: "Success",
			requestBody: CalculateShippingCostRequest{
				ShipperDestinationID:  "123",
				ReceiverDestinationID: "456",
				Weight:                1.0,
				ItemValue:             10000,
				COD:                   "no",
			},
			setupMock: func(m *MockKomerceService) {
				m.On("CalculateShippingCost", "123", "456", 1.0, 10000, "no").Return(&models.KomerceCalculateResponse{
					Meta: models.KomerceMeta{Status: "success"},
					Data: models.KomerceCalculateData{
						CalculateReguler: []models.KomerceShippingOption{
							{ServiceName: "REG", ShippingCost: 10000},
						},
					},
				}, nil)
			},
			expectedStatus: http.StatusOK,
			expectedError:  false,
		},
		{
			name:           "Invalid Request Body",
			requestBody:    "invalid-json",
			setupMock:      func(m *MockKomerceService) {},
			expectedStatus: http.StatusBadRequest, // Handler returns 400 for bad parsing
			expectedError:  true,
		},
		{
			name: "Service Error",
			requestBody: CalculateShippingCostRequest{
				ShipperDestinationID:  "123",
				ReceiverDestinationID: "456",
				Weight:                1.0,
				ItemValue:             10000,
				COD:                   "no",
			},
			setupMock: func(m *MockKomerceService) {
				m.On("CalculateShippingCost", "123", "456", 1.0, 10000, "no").Return((*models.KomerceCalculateResponse)(nil), errors.New("service error"))
			},
			expectedStatus: http.StatusInternalServerError,
			expectedError:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockService := new(MockKomerceService)
			tt.setupMock(mockService)

			handler := NewKomerceHandler(mockService)
			app := fiber.New()
			app.Post("/shipping/calculate", handler.CalculateShippingCost)

			var body []byte
			if s, ok := tt.requestBody.(string); ok && s == "invalid-json" {
				body = []byte(tt.requestBody.(string))
			} else {
				body, _ = json.Marshal(tt.requestBody)
			}

			req := httptest.NewRequest("POST", "/shipping/calculate", bytes.NewReader(body))
			req.Header.Set("Content-Type", "application/json")
			resp, _ := app.Test(req)

			assert.Equal(t, tt.expectedStatus, resp.StatusCode)
			mockService.AssertExpectations(t)
		})
	}
}
