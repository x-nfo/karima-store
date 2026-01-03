package middleware

import (
	"fmt"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/karima-store/internal/models"
)

// RequirePermission checks if the authenticated user has a specific permission
func (m *KratosAuthProvider) RequirePermission(permission models.Permission) fiber.Handler {
	return func(c *fiber.Ctx) error {
		userRole := c.Locals("user_role")
		if userRole == nil {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "User not authenticated",
				"code":  "UNAUTHORIZED",
			})
		}

		// Convert userRole to models.UserRole
		role, ok := userRole.(models.UserRole)
		if !ok {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Invalid user role type",
				"code":  "INTERNAL_SERVER_ERROR",
			})
		}

		// Check if user has the required permission
		if !models.HasPermission(role, permission) {
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
				"error": fmt.Sprintf("Insufficient permissions. Required: %s", permission),
				"code":  "FORBIDDEN",
			})
		}

		return c.Next()
	}
}

// RequireOwnership validates that the authenticated user owns the resource
// The resource ID should be in the route params (e.g., /orders/:id)
// paramName is the name of the parameter containing the resource owner ID (default: "id")
func (m *KratosAuthProvider) RequireOwnership(paramName string) fiber.Handler {
	if paramName == "" {
		paramName = "id"
	}

	return func(c *fiber.Ctx) error {
		// Get authenticated user
		localUserID := c.Locals("local_user_id")
		if localUserID == nil {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "User not authenticated",
				"code":  "UNAUTHORIZED",
			})
		}

		userID, ok := localUserID.(uint)
		if !ok {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Invalid user ID type",
				"code":  "INTERNAL_SERVER_ERROR",
			})
		}

		// Get user role
		userRole := c.Locals("user_role")
		role, ok := userRole.(models.UserRole)
		if !ok {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Invalid user role type",
				"code":  "INTERNAL_SERVER_ERROR",
			})
		}

		// Admins can access all resources
		if models.IsAdmin(role) {
			return c.Next()
		}

		// Get resource owner ID from params
		resourceIDStr := c.Params(paramName)
		if resourceIDStr == "" {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Resource ID not provided",
				"code":  "BAD_REQUEST",
			})
		}

		resourceID, err := strconv.ParseUint(resourceIDStr, 10, 32)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Invalid resource ID format",
				"code":  "BAD_REQUEST",
			})
		}

		// For now, we assume the resource ID in params is the owner ID
		// In a real implementation, you'd query the database to get the actual owner
		// This is a simplified version - actual implementation should be in service layer
		resourceOwnerID := uint(resourceID)

		// Check if user owns the resource
		if !models.CanAccessResource(userID, resourceOwnerID, role) {
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
				"error": "You don't have permission to access this resource",
				"code":  "FORBIDDEN",
			})
		}

		return c.Next()
	}
}

// RequireAdminOrOwner allows access if user is admin OR owns the resource
// This combines admin check with ownership validation
func (m *KratosAuthProvider) RequireAdminOrOwner(ownerIDParam string) fiber.Handler {
	if ownerIDParam == "" {
		ownerIDParam = "user_id"
	}

	return func(c *fiber.Ctx) error {
		// Get authenticated user
		localUserID := c.Locals("local_user_id")
		if localUserID == nil {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "User not authenticated",
				"code":  "UNAUTHORIZED",
			})
		}

		userID, ok := localUserID.(uint)
		if !ok {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Invalid user ID type",
				"code":  "INTERNAL_SERVER_ERROR",
			})
		}

		// Get user role
		userRole := c.Locals("user_role")
		role, ok := userRole.(models.UserRole)
		if !ok {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Invalid user role type",
				"code":  "INTERNAL_SERVER_ERROR",
			})
		}

		// Admins can access all resources
		if models.IsAdmin(role) {
			return c.Next()
		}

		// Get resource owner ID from params
		ownerIDStr := c.Params(ownerIDParam)
		if ownerIDStr == "" {
			// If no owner ID in params, check query params
			ownerIDStr = c.Query(ownerIDParam)
		}

		if ownerIDStr == "" {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Owner ID not provided",
				"code":  "BAD_REQUEST",
			})
		}

		ownerID, err := strconv.ParseUint(ownerIDStr, 10, 32)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Invalid owner ID format",
				"code":  "BAD_REQUEST",
			})
		}

		// Check if user is the owner
		if userID != uint(ownerID) {
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
				"error": "You don't have permission to access this resource",
				"code":  "FORBIDDEN",
			})
		}

		return c.Next()
	}
}
