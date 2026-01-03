package middleware

import (
	"strings"

	"github.com/gofiber/fiber/v2"
)

// CORS middleware configuration
func CORS(allowedOrigins string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		origin := c.Get("Origin")
		allowOrigin := ""

		// Check if origin is allowed
		if allowedOrigins == "*" {
			allowOrigin = "*"
		} else {
			origins := strings.Split(allowedOrigins, ",")
			for _, o := range origins {
				if o == origin {
					allowOrigin = origin
					break
				}
			}
		}

		// Only set CORS headers if origin is allowed
		if allowOrigin != "" {
			c.Set("Access-Control-Allow-Origin", allowOrigin)
			c.Set("Access-Control-Allow-Methods", "GET,POST,PUT,DELETE,OPTIONS")
			c.Set("Access-Control-Allow-Headers", "Origin,Content-Type,Accept,Authorization,X-Requested-With")
			c.Set("Access-Control-Expose-Headers", "Content-Length,Content-Type")
			c.Set("Access-Control-Allow-Credentials", "true")
			c.Set("Access-Control-Max-Age", "86400")
		}

		// Handle preflight requests
		if c.Method() == "OPTIONS" {
			return c.SendStatus(fiber.StatusNoContent)
		}

		return c.Next()
	}
}
