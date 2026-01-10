package handlers

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http/httptest"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/karima-store/internal/config"
	"github.com/karima-store/internal/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockAuthService
type MockAuthService struct {
	mock.Mock
}

func (m *MockAuthService) Register(email, password, fullName string) (*models.User, string, error) {
	args := m.Called(email, password, fullName)
	if args.Get(0) == nil {
		return nil, "", args.Error(2)
	}
	return args.Get(0).(*models.User), args.String(1), args.Error(2)
}

func (m *MockAuthService) Login(email, password string) (*models.User, string, error) {
	args := m.Called(email, password)
	if args.Get(0) == nil {
		return nil, "", args.Error(2)
	}
	return args.Get(0).(*models.User), args.String(1), args.Error(2)
}

func (m *MockAuthService) FindOrCreateByOAuth(provider, email, providerID, fullName, avatar string) (*models.User, string, error) {
	args := m.Called(provider, email, providerID, fullName, avatar)
	if args.Get(0) == nil {
		return nil, "", args.Error(2)
	}
	return args.Get(0).(*models.User), args.String(1), args.Error(2)
}

func (m *MockAuthService) GetUserByID(id uint) (*models.User, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.User), args.Error(1)
}

func TestAuthHandler_Register(t *testing.T) {
	mockService := new(MockAuthService)
	cfg := &config.Config{AppEnv: "test"}
	handler := NewAuthHandler(mockService, cfg)

	app := fiber.New()
	app.Post("/auth/register", handler.Register)

	t.Run("Success", func(t *testing.T) {
		reqBody := map[string]string{
			"email":     "test@example.com",
			"password":  "password123",
			"full_name": "Test User",
		}
		body, _ := json.Marshal(reqBody)
		req := httptest.NewRequest("POST", "/auth/register", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")

		mockUser := &models.User{Email: "test@example.com", FullName: "Test User"}
		mockService.On("Register", "test@example.com", "password123", "Test User").Return(mockUser, "test-token", nil)

		resp, err := app.Test(req)

		assert.NoError(t, err)
		assert.Equal(t, fiber.StatusOK, resp.StatusCode)

		var result map[string]interface{}
		json.NewDecoder(resp.Body).Decode(&result)
		assert.Equal(t, "test-token", result["token"])
		mockService.AssertExpectations(t)
	})

	t.Run("InvalidPayload", func(t *testing.T) {
		req := httptest.NewRequest("POST", "/auth/register", bytes.NewReader([]byte("invalid json")))
		req.Header.Set("Content-Type", "application/json")

		resp, err := app.Test(req)
		assert.NoError(t, err)
		assert.Equal(t, fiber.StatusBadRequest, resp.StatusCode)
	})
}

func TestAuthHandler_Login(t *testing.T) {
	mockService := new(MockAuthService)
	cfg := &config.Config{AppEnv: "test"}
	handler := NewAuthHandler(mockService, cfg)

	app := fiber.New()
	app.Post("/auth/login", handler.Login)

	t.Run("Success", func(t *testing.T) {
		reqBody := map[string]string{
			"email":    "test@example.com",
			"password": "password123",
		}
		body, _ := json.Marshal(reqBody)
		req := httptest.NewRequest("POST", "/auth/login", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")

		mockUser := &models.User{Email: "test@example.com"}
		mockService.On("Login", "test@example.com", "password123").Return(mockUser, "test-token", nil)

		resp, err := app.Test(req)

		assert.NoError(t, err)
		assert.Equal(t, fiber.StatusOK, resp.StatusCode)
		mockService.AssertExpectations(t)
	})

	t.Run("InvalidCredentials", func(t *testing.T) {
		reqBody := map[string]string{
			"email":    "test@example.com",
			"password": "wrongpassword",
		}
		body, _ := json.Marshal(reqBody)
		req := httptest.NewRequest("POST", "/auth/login", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")

		mockService.On("Login", "test@example.com", "wrongpassword").Return(nil, "", errors.New("invalid credentials"))

		resp, err := app.Test(req)

		assert.NoError(t, err)
		assert.Equal(t, fiber.StatusUnauthorized, resp.StatusCode)
	})
}

func TestAuthHandler_Me(t *testing.T) {
	mockService := new(MockAuthService)
	cfg := &config.Config{AppEnv: "test"}
	handler := NewAuthHandler(mockService, cfg)

	app := fiber.New()

	// Mock middleware setting user_id
	app.Use(func(c *fiber.Ctx) error {
		c.Locals("user_id", float64(1)) // JWT often parses numbers as float64
		return c.Next()
	})

	app.Get("/auth/me", handler.Me)

	t.Run("Success", func(t *testing.T) {
		mockUser := &models.User{ID: 1, Email: "test@example.com"}
		mockService.On("GetUserByID", uint(1)).Return(mockUser, nil)

		req := httptest.NewRequest("GET", "/auth/me", nil)
		resp, err := app.Test(req)

		assert.NoError(t, err)
		assert.Equal(t, fiber.StatusOK, resp.StatusCode)
		mockService.AssertExpectations(t)
	})
}
