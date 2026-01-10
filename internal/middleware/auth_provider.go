package middleware

import "github.com/gofiber/fiber/v2"

// AuthProvider interface for authentication middleware
// This provides a unified interface for authentication using Go Fiber
type AuthProvider interface {
	// ValidateToken validates the session token and sets user context
	ValidateToken() fiber.Handler

	// RequireRole checks if the authenticated user has one of the required roles
	RequireRole(roles ...string) fiber.Handler

	// RequireAdmin is a convenience method that checks if the user is an admin
	RequireAdmin() fiber.Handler

	// OptionalAuth validates session if present, but doesn't require it
	OptionalAuth() fiber.Handler
}
