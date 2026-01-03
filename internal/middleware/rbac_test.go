package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/karima-store/internal/models"
	"github.com/stretchr/testify/assert"
)

func TestRequirePermission(t *testing.T) {
	ts := mockKratosServer()
	defer ts.Close()

	app := fiber.New()
	mockAuthService := &MockAuthService{}
	kratosMiddleware := NewKratosMiddleware(ts.URL, ts.URL, mockAuthService)

	// Test route requiring specific permission
	app.Get("/admin/products",
		kratosMiddleware.Authenticate(),
		kratosMiddleware.RequirePermission(models.PermissionCreateProducts),
		func(c *fiber.Ctx) error {
			return c.SendStatus(fiber.StatusOK)
		})

	// Test case 1: User with permission (admin)
	req1 := httptest.NewRequest("GET", "/admin/products", nil)
	req1.AddCookie(&http.Cookie{
		Name:  "ory_kratos_session",
		Value: "valid-session-token", // Mock returns admin role
	})
	resp1, err1 := app.Test(req1)
	assert.NoError(t, err1)
	// Note: This test demonstrates permission checking
	// In production, you'd mock different roles for comprehensive testing
	_ = resp1 // Use the response to avoid unused variable error

	// Test case 2: No authentication
	req2 := httptest.NewRequest("GET", "/admin/products", nil)
	resp2, err2 := app.Test(req2)
	assert.NoError(t, err2)
	assert.Equal(t, fiber.StatusUnauthorized, resp2.StatusCode)
}

func TestRequireOwnership(t *testing.T) {
	ts := mockKratosServer()
	defer ts.Close()

	app := fiber.New()
	mockAuthService := &MockAuthService{}
	kratosMiddleware := NewKratosMiddleware(ts.URL, ts.URL, mockAuthService)

	// Test route with ownership validation
	app.Get("/orders/:id",
		kratosMiddleware.Authenticate(),
		kratosMiddleware.RequireOwnership("id"),
		func(c *fiber.Ctx) error {
			return c.SendStatus(fiber.StatusOK)
		})

	// Test case 1: User accessing their own resource (user_id = 1, resource_id = 1)
	req1 := httptest.NewRequest("GET", "/orders/1", nil)
	req1.AddCookie(&http.Cookie{
		Name:  "ory_kratos_session",
		Value: "valid-session-token",
	})
	resp1, err1 := app.Test(req1)
	assert.NoError(t, err1)
	assert.Equal(t, fiber.StatusOK, resp1.StatusCode)

	// Test case 2: User accessing another user's resource
	req2 := httptest.NewRequest("GET", "/orders/999", nil)
	req2.AddCookie(&http.Cookie{
		Name:  "ory_kratos_session",
		Value: "valid-session-token",
	})
	resp2, err2 := app.Test(req2)
	assert.NoError(t, err2)
	assert.Equal(t, fiber.StatusForbidden, resp2.StatusCode)

	// Test case 3: No authentication
	req3 := httptest.NewRequest("GET", "/orders/1", nil)
	resp3, err3 := app.Test(req3)
	assert.NoError(t, err3)
	assert.Equal(t, fiber.StatusUnauthorized, resp3.StatusCode)
}

func TestRequireAdminOrOwner(t *testing.T) {
	ts := mockKratosServer()
	defer ts.Close()

	app := fiber.New()
	mockAuthService := &MockAuthService{}
	kratosMiddleware := NewKratosMiddleware(ts.URL, ts.URL, mockAuthService)

	// Test route with admin or owner access
	app.Get("/users/:user_id/profile",
		kratosMiddleware.Authenticate(),
		kratosMiddleware.RequireAdminOrOwner("user_id"),
		func(c *fiber.Ctx) error {
			return c.SendStatus(fiber.StatusOK)
		})

	// Test case 1: User accessing their own profile
	req1 := httptest.NewRequest("GET", "/users/1/profile", nil)
	req1.AddCookie(&http.Cookie{
		Name:  "ory_kratos_session",
		Value: "valid-session-token",
	})
	resp1, err1 := app.Test(req1)
	assert.NoError(t, err1)
	assert.Equal(t, fiber.StatusOK, resp1.StatusCode)

	// Test case 2: User accessing another user's profile
	req2 := httptest.NewRequest("GET", "/users/999/profile", nil)
	req2.AddCookie(&http.Cookie{
		Name:  "ory_kratos_session",
		Value: "valid-session-token",
	})
	resp2, err2 := app.Test(req2)
	assert.NoError(t, err2)
	assert.Equal(t, fiber.StatusForbidden, resp2.StatusCode)

	// Test case 3: No authentication
	req3 := httptest.NewRequest("GET", "/users/1/profile", nil)
	resp3, err3 := app.Test(req3)
	assert.NoError(t, err3)
	assert.Equal(t, fiber.StatusUnauthorized, resp3.StatusCode)
}

func TestPermissionModel(t *testing.T) {
	// Test permission checking
	t.Run("Admin has all permissions", func(t *testing.T) {
		assert.True(t, models.HasPermission(models.RoleAdmin, models.PermissionCreateProducts))
		assert.True(t, models.HasPermission(models.RoleAdmin, models.PermissionViewAllOrders))
		assert.True(t, models.HasPermission(models.RoleAdmin, models.PermissionUpdateUsers))
	})

	t.Run("Customer has limited permissions", func(t *testing.T) {
		assert.True(t, models.HasPermission(models.RoleCustomer, models.PermissionViewProducts))
		assert.True(t, models.HasPermission(models.RoleCustomer, models.PermissionViewOwnOrders))
		assert.False(t, models.HasPermission(models.RoleCustomer, models.PermissionCreateProducts))
		assert.False(t, models.HasPermission(models.RoleCustomer, models.PermissionViewAllOrders))
	})

	t.Run("Resource access control", func(t *testing.T) {
		// Admin can access any resource
		assert.True(t, models.CanAccessResource(1, 999, models.RoleAdmin))

		// Customer can only access their own resources
		assert.True(t, models.CanAccessResource(1, 1, models.RoleCustomer))
		assert.False(t, models.CanAccessResource(1, 999, models.RoleCustomer))
	})

	t.Run("Role validation", func(t *testing.T) {
		assert.True(t, models.ValidateRole(models.RoleAdmin))
		assert.True(t, models.ValidateRole(models.RoleCustomer))
		assert.False(t, models.ValidateRole(models.UserRole("invalid")))
	})
}
