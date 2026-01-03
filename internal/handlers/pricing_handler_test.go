package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/karima-store/internal/services"
	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockPricingService is a mock implementation of PricingService
type MockPricingService struct {
	mock.Mock
}

func (m *MockPricingService) CalculatePrice(req services.PriceCalculationRequest) (*services.PriceCalculationResponse, error) {
	args := m.Called(req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*services.PriceCalculationResponse), args.Error(1)
}

func (m *MockPricingService) CalculateShippingCost(req services.ShippingCalculationRequest) (*services.ShippingCalculationResponse, error) {
	args := m.Called(req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*services.ShippingCalculationResponse), args.Error(1)
}

func (m *MockPricingService) CalculateOrderSummary(items []services.PriceCalculationRequest, shipping services.ShippingCalculationRequest, customerType services.CustomerType) (*services.OrderSummary, error) {
	args := m.Called(items, shipping, customerType)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*services.OrderSummary), args.Error(1)
}

func (m *MockPricingService) CalculateCouponDiscount(req services.CouponCalculationRequest) (float64, string, error) {
	args := m.Called(req)
	return args.Get(0).(float64), args.String(1), args.Error(2)
}

func (m *MockPricingService) CheckFreeShipping(orderAmount float64, regionCode string) (bool, error) {
	args := m.Called(orderAmount, regionCode)
	return args.Bool(0), args.Error(1)
}

func (m *MockPricingService) ApplyCouponToPriceCalculation(resp *services.PriceCalculationResponse, couponReq services.CouponCalculationRequest) error {
	args := m.Called(resp, couponReq)
	return args.Error(0)
}

// MockRedisClient is a mock implementation of RedisClient
type MockRedisClient struct {
	mock.Mock
}

func (m *MockRedisClient) Get(ctx context.Context, key string) (string, error) {
	args := m.Called(ctx, key)
	return args.String(0), args.Error(1)
}

func (m *MockRedisClient) Set(ctx context.Context, key string, value interface{}, expiration time.Duration) error {
	args := m.Called(ctx, key, value, expiration)
	return args.Error(0)
}

func (m *MockRedisClient) GetJSON(ctx context.Context, key string, dest interface{}) error {
	return nil
}
func (m *MockRedisClient) SetJSON(ctx context.Context, key string, value interface{}, expiration time.Duration) error {
	return nil
}
func (m *MockRedisClient) Delete(ctx context.Context, keys ...string) error {
	return nil
}
func (m *MockRedisClient) Exists(ctx context.Context, keys ...string) (int64, error) {
	return 0, nil
}
func (m *MockRedisClient) FlushDB(ctx context.Context) error {
	return nil
}
func (m *MockRedisClient) DeleteByPattern(ctx context.Context, pattern string) error {
	return nil
}
func (m *MockRedisClient) HealthCheck(ctx context.Context) error {
	return nil
}
func (m *MockRedisClient) PoolStats() map[string]interface{} {
	return nil
}
func (m *MockRedisClient) Client() *redis.Client {
	return nil
}
func (m *MockRedisClient) Close() error {
	return nil
}

func TestPricingHandler_CalculatePrice(t *testing.T) {
	tests := []struct {
		name           string
		requestBody    interface{}
		setupMock      func(*MockPricingService)
		expectedStatus int
		expectedError  bool
	}{
		{
			name: "Success",
			requestBody: services.PriceCalculationRequest{
				ProductID:    1,
				Quantity:     2,
				CustomerType: services.CustomerRetail,
			},
			setupMock: func(m *MockPricingService) {
				m.On("CalculatePrice", mock.AnythingOfType("services.PriceCalculationRequest")).Return(&services.PriceCalculationResponse{
					FinalPrice: 20000,
				}, nil)
			},
			expectedStatus: http.StatusOK,
			expectedError:  false,
		},
		{
			name: "Service Error",
			requestBody: services.PriceCalculationRequest{
				ProductID:    1,
				Quantity:     2,
				CustomerType: services.CustomerRetail,
			},
			setupMock: func(m *MockPricingService) {
				m.On("CalculatePrice", mock.AnythingOfType("services.PriceCalculationRequest")).Return((*services.PriceCalculationResponse)(nil), errors.New("service error"))
			},
			expectedStatus: http.StatusNotFound, // Handler mapping for error
			expectedError:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockService := new(MockPricingService)
			mockRedis := new(MockRedisClient)
			tt.setupMock(mockService)

			handler := NewPricingHandler(mockService, mockRedis)
			app := fiber.New()
			app.Post("/calculate", handler.CalculatePrice)

			body, _ := json.Marshal(tt.requestBody)
			req := httptest.NewRequest("POST", "/calculate", bytes.NewReader(body))
			req.Header.Set("Content-Type", "application/json")
			resp, _ := app.Test(req)

			assert.Equal(t, tt.expectedStatus, resp.StatusCode)
			mockService.AssertExpectations(t)
		})
	}
}
