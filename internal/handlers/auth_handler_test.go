package handlers

import (
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

func (m *MockAuthService) SyncUser(kratosIdentity *models.KratosIdentity, email string) (*models.User, error) {
	args := m.Called(kratosIdentity, email)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.User), args.Error(1)
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
	cfg := &config.Config{
		KratosPublicURL: "http://kratos:4433",
	}
	handler := NewAuthHandler(mockService, cfg)

	app := fiber.New()
	app.Post("/auth/register", handler.Register)

	req := httptest.NewRequest("POST", "/auth/register", nil)
	resp, err := app.Test(req)

	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusSeeOther, resp.StatusCode) // 303 See Other
	assert.Equal(t, cfg.KratosPublicURL+"/self-service/registration/browser", resp.Header.Get("Location"))
}

func TestAuthHandler_Login(t *testing.T) {
	mockService := new(MockAuthService)
	cfg := &config.Config{
		KratosPublicURL: "http://kratos:4433",
	}
	handler := NewAuthHandler(mockService, cfg)

	app := fiber.New()
	app.Post("/auth/login", handler.Login)

	req := httptest.NewRequest("POST", "/auth/login", nil)
	resp, err := app.Test(req)

	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusSeeOther, resp.StatusCode)
	assert.Equal(t, cfg.KratosPublicURL+"/self-service/login/browser", resp.Header.Get("Location"))
}

func TestAuthHandler_Logout(t *testing.T) {
	mockService := new(MockAuthService)
	cfg := &config.Config{
		KratosPublicURL: "http://kratos:4433",
	}
	handler := NewAuthHandler(mockService, cfg)

	app := fiber.New()
	app.Post("/auth/logout", handler.Logout)

	req := httptest.NewRequest("POST", "/auth/logout", nil)
	resp, err := app.Test(req)

	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusSeeOther, resp.StatusCode)
	assert.Equal(t, cfg.KratosPublicURL+"/self-service/browser/flows/logout", resp.Header.Get("Location"))
}

func TestAuthHandler_Me_Success(t *testing.T) {
	mockService := new(MockAuthService)
	cfg := &config.Config{}
	handler := NewAuthHandler(mockService, cfg)

	app := fiber.New()

	// Middleware to inject session manually for testing (simulating KratosMiddleware)
	app.Use(func(c *fiber.Ctx) error {
		c.Locals("session", &models.KratosSession{
			Identity: models.KratosIdentity{
				ID: "test-kratos-id",
			},
		})
		c.Locals("user_email", "test@example.com")
		return c.Next()
	})

	app.Get("/auth/me", handler.Me)

	// Mock Service Response
	mockUser := &models.User{
		KratosID: "test-kratos-id",
		Email:    "test@example.com",
		FullName: "Test User",
	}
	mockService.On("SyncUser", mock.Anything, "test@example.com").Return(mockUser, nil)

	req := httptest.NewRequest("GET", "/auth/me", nil)
	resp, err := app.Test(req)

	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusOK, resp.StatusCode)

	var result map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&result)

	data := result["data"].(map[string]interface{})
	user := data["user"].(map[string]interface{})

	assert.Equal(t, "test-kratos-id", user["kratos_id"])
	assert.Equal(t, "test@example.com", user["email"])
	mockService.AssertExpectations(t)
}

func TestAuthHandler_Me_Unauthenticated(t *testing.T) {
	mockService := new(MockAuthService)
	cfg := &config.Config{}
	handler := NewAuthHandler(mockService, cfg)

	app := fiber.New()
	// No session injection middleware = unauthenticated

	app.Get("/auth/me", handler.Me)

	req := httptest.NewRequest("GET", "/auth/me", nil)
	resp, err := app.Test(req)

	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusUnauthorized, resp.StatusCode)
}

func TestAuthHandler_Me_SyncError(t *testing.T) {
	mockService := new(MockAuthService)
	cfg := &config.Config{}
	handler := NewAuthHandler(mockService, cfg)

	app := fiber.New()
	app.Use(func(c *fiber.Ctx) error {
		c.Locals("session", &models.KratosSession{
			Identity: models.KratosIdentity{ID: "test-kratos-id"},
		})
		c.Locals("user_email", "test@example.com")
		return c.Next()
	})

	app.Get("/auth/me", handler.Me)

	mockService.On("SyncUser", mock.Anything, "test@example.com").Return(nil, errors.New("sync failed"))

	req := httptest.NewRequest("GET", "/auth/me", nil)
	resp, err := app.Test(req)

	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusInternalServerError, resp.StatusCode)
	mockService.AssertExpectations(t)
}
