package middleware

import "github.com/gofiber/fiber/v2"

// KratosMiddleware interface for authentication middleware
// This allows us to use Ory Kratos for authentication
type KratosMiddleware interface {
	// ValidateToken validates the session token and sets user context
	ValidateToken() fiber.Handler

	// RequireRole checks if the authenticated user has one of the required roles
	RequireRole(roles ...string) fiber.Handler

	// RequireAdmin is a convenience method that checks if the user is an admin
	RequireAdmin() fiber.Handler

	// OptionalAuth validates session if present, but doesn't require it
	OptionalAuth() fiber.Handler
}
