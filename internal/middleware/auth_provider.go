package middleware

import "github.com/gofiber/fiber/v2"

// AuthProvider is an interface for authentication middleware
// This allows us to use different authentication providers (JWT, Kratos, etc.)
type AuthProvider interface {
	// Authenticate validates the user's credentials and sets user context
	Authenticate() fiber.Handler

	// RequireRole checks if the authenticated user has one of the required roles
	RequireRole(roles ...string) fiber.Handler

	// RequireAdmin is a convenience method that checks if the user is an admin
	RequireAdmin() fiber.Handler

	// OptionalAuth validates credentials if present, but doesn't require them
	OptionalAuth() fiber.Handler
}
