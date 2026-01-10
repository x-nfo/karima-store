package middleware

import (
	"net/http/httptest"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/karima-store/internal/config"
	"github.com/karima-store/internal/models"
	"github.com/stretchr/testify/assert"
)

func TestRequirePermission(t *testing.T) {
	middleware := NewAuthMiddleware(&config.Config{})

	t.Run("Allowed_Admin", func(t *testing.T) {
		app := fiber.New()
		app.Use(func(c *fiber.Ctx) error {
			c.Locals("role", models.RoleAdmin)
			return c.Next()
		})
		app.Get("/test", middleware.RequirePermission(models.PermissionCreateProducts), func(c *fiber.Ctx) error {
			return c.SendStatus(fiber.StatusOK)
		})

		req := httptest.NewRequest("GET", "/test", nil)
		resp, err := app.Test(req)
		assert.NoError(t, err)
		assert.Equal(t, fiber.StatusOK, resp.StatusCode)
	})

	t.Run("Denied_Customer", func(t *testing.T) {
		app := fiber.New()
		app.Use(func(c *fiber.Ctx) error {
			c.Locals("role", models.RoleCustomer)
			return c.Next()
		})
		app.Get("/test", middleware.RequirePermission(models.PermissionCreateProducts), func(c *fiber.Ctx) error {
			return c.SendStatus(fiber.StatusOK)
		})

		req := httptest.NewRequest("GET", "/test", nil)
		resp, err := app.Test(req)
		assert.NoError(t, err)
		assert.Equal(t, fiber.StatusForbidden, resp.StatusCode)
	})

	t.Run("Unauthorized_NoRole", func(t *testing.T) {
		app := fiber.New()
		app.Get("/test", middleware.RequirePermission(models.PermissionCreateProducts), func(c *fiber.Ctx) error {
			return c.SendStatus(fiber.StatusOK)
		})

		req := httptest.NewRequest("GET", "/test", nil)
		resp, err := app.Test(req)
		assert.NoError(t, err)
		assert.Equal(t, fiber.StatusUnauthorized, resp.StatusCode)
	})
}

func TestRequireOwnership(t *testing.T) {
	middleware := NewAuthMiddleware(&config.Config{})

	t.Run("Allowed_Owner", func(t *testing.T) {
		app := fiber.New()
		app.Use(func(c *fiber.Ctx) error {
			c.Locals("user_id", float64(1))
			c.Locals("role", models.RoleCustomer)
			return c.Next()
		})
		app.Get("/orders/:id", middleware.RequireOwnership("id"), func(c *fiber.Ctx) error {
			return c.SendStatus(fiber.StatusOK)
		})

		req := httptest.NewRequest("GET", "/orders/1", nil)
		resp, err := app.Test(req)
		assert.NoError(t, err)
		assert.Equal(t, fiber.StatusOK, resp.StatusCode)
	})

	t.Run("Denied_NotOwner", func(t *testing.T) {
		app := fiber.New()
		app.Use(func(c *fiber.Ctx) error {
			c.Locals("user_id", float64(1)) // User 1
			c.Locals("role", models.RoleCustomer)
			return c.Next()
		})
		app.Get("/orders/:id", middleware.RequireOwnership("id"), func(c *fiber.Ctx) error {
			return c.SendStatus(fiber.StatusOK)
		})

		req := httptest.NewRequest("GET", "/orders/2", nil) // Accessing User 2 (implied by ID)
		// Note: Requirement is simplified here. In real app, service layer checks order ownership.
		// The middleware assumes param ID matches user ID logic for simplicity or basic checking.
		// Actually, standard RequireOwnership usually checks if resource.OwnerID == CurrentUser.ID.
		// But the simplified middleware in rbac.go assumes:
		// "For now, we assume the resource ID in params is the owner ID" (line 101 in rbac.go)
		// So accessing /orders/2 means "Owner is 2". User is 1. 1 != 2 -> Forbidden.

		resp, err := app.Test(req)
		assert.NoError(t, err)
		assert.Equal(t, fiber.StatusForbidden, resp.StatusCode)
	})

	t.Run("Allowed_Admin_Override", func(t *testing.T) {
		app := fiber.New()
		app.Use(func(c *fiber.Ctx) error {
			c.Locals("user_id", float64(1))
			c.Locals("role", models.RoleAdmin)
			return c.Next()
		})
		app.Get("/orders/:id", middleware.RequireOwnership("id"), func(c *fiber.Ctx) error {
			return c.SendStatus(fiber.StatusOK)
		})

		req := httptest.NewRequest("GET", "/orders/2", nil)
		resp, err := app.Test(req)
		assert.NoError(t, err)
		assert.Equal(t, fiber.StatusOK, resp.StatusCode)
	})
}

func TestRequireAdminOrOwner(t *testing.T) {
	middleware := NewAuthMiddleware(&config.Config{})

	t.Run("Allowed_Owner", func(t *testing.T) {
		app := fiber.New()
		app.Use(func(c *fiber.Ctx) error {
			c.Locals("user_id", float64(1))
			c.Locals("role", models.RoleCustomer)
			return c.Next()
		})
		app.Get("/users/:user_id/profile", middleware.RequireAdminOrOwner("user_id"), func(c *fiber.Ctx) error {
			return c.SendStatus(fiber.StatusOK)
		})

		req := httptest.NewRequest("GET", "/users/1/profile", nil)
		resp, err := app.Test(req)
		assert.NoError(t, err)
		assert.Equal(t, fiber.StatusOK, resp.StatusCode)
	})

	t.Run("Denied_NotOwner", func(t *testing.T) {
		app := fiber.New()
		app.Use(func(c *fiber.Ctx) error {
			c.Locals("user_id", float64(1))
			c.Locals("role", models.RoleCustomer)
			return c.Next()
		})
		app.Get("/users/:user_id/profile", middleware.RequireAdminOrOwner("user_id"), func(c *fiber.Ctx) error {
			return c.SendStatus(fiber.StatusOK)
		})

		req := httptest.NewRequest("GET", "/users/2/profile", nil)
		resp, err := app.Test(req)
		assert.NoError(t, err)
		assert.Equal(t, fiber.StatusForbidden, resp.StatusCode)
	})

	t.Run("Allowed_Admin", func(t *testing.T) {
		app := fiber.New()
		app.Use(func(c *fiber.Ctx) error {
			c.Locals("user_id", float64(1))
			c.Locals("role", models.RoleAdmin)
			return c.Next()
		})
		app.Get("/users/:user_id/profile", middleware.RequireAdminOrOwner("user_id"), func(c *fiber.Ctx) error {
			return c.SendStatus(fiber.StatusOK)
		})

		req := httptest.NewRequest("GET", "/users/2/profile", nil)
		resp, err := app.Test(req)
		assert.NoError(t, err)
		assert.Equal(t, fiber.StatusOK, resp.StatusCode)
	})
}

func TestPermissionModel(t *testing.T) {
	// Test permission checking logic (pure model test)
	t.Run("Admin has all permissions", func(t *testing.T) {
		assert.True(t, models.HasPermission(models.RoleAdmin, models.PermissionCreateProducts))
	})
	t.Run("Customer has limited permissions", func(t *testing.T) {
		assert.True(t, models.HasPermission(models.RoleCustomer, models.PermissionViewProducts))
		assert.False(t, models.HasPermission(models.RoleCustomer, models.PermissionCreateProducts))
	})
}
