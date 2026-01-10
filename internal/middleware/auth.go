package middleware

import (
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"github.com/karima-store/internal/config"
)

type AuthMiddleware struct {
	cfg *config.Config
}

func NewAuthMiddleware(cfg *config.Config) *AuthMiddleware {
	return &AuthMiddleware{cfg: cfg}
}

// ValidateToken validates the JWT token and sets user context
func (m *AuthMiddleware) ValidateToken() fiber.Handler {
	return func(c *fiber.Ctx) error {
		// 1. Get token from Cookie or Header
		var tokenString string

		// Try Cookie
		tokenString = c.Cookies("jwt")

		// Try Header (Bearer )
		if tokenString == "" {
			authHeader := c.Get("Authorization")
			if strings.HasPrefix(authHeader, "Bearer ") {
				tokenString = strings.TrimPrefix(authHeader, "Bearer ")
			}
		}

		if tokenString == "" {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Unauthorized",
				"code":  "UNAUTHORIZED",
			})
		}

		// 2. Parse and Validate Token
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fiber.ErrUnauthorized
			}
			return []byte(m.cfg.JWTSecret), nil
		})

		if err != nil || !token.Valid {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Invalid or expired token",
				"code":  "UNAUTHORIZED",
			})
		}

		// 3. Set Claims to Locals
		if claims, ok := token.Claims.(jwt.MapClaims); ok {
			c.Locals("user_id", claims["user_id"])
			c.Locals("email", claims["email"])
			c.Locals("role", claims["role"])
		}

		return c.Next()
	}
}

// Protected middleware verifies the JWT token (alias for ValidateToken)
func (m *AuthMiddleware) Protected() fiber.Handler {
	return m.ValidateToken()
}

// RequireRole checks if the user has the required role (must be used after Protected)
func (m *AuthMiddleware) RequireRole(roles ...string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		userRole, ok := c.Locals("role").(string)
		if !ok {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Unauthorized",
			})
		}

		for _, role := range roles {
			if userRole == role {
				return c.Next()
			}
		}

		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"error": "Forbidden: Insufficient permissions",
		})
	}
}

// RequireAdmin is a convenience method that checks if the user is an admin
func (m *AuthMiddleware) RequireAdmin() fiber.Handler {
	return m.RequireRole("admin")
}

// OptionalAuth validates session if present, but doesn't require it
func (m *AuthMiddleware) OptionalAuth() fiber.Handler {
	return func(c *fiber.Ctx) error {
		// 1. Get token from Cookie or Header
		var tokenString string

		// Try Cookie
		tokenString = c.Cookies("jwt")

		// Try Header (Bearer )
		if tokenString == "" {
			authHeader := c.Get("Authorization")
			if strings.HasPrefix(authHeader, "Bearer ") {
				tokenString = strings.TrimPrefix(authHeader, "Bearer ")
			}
		}

		// If no token, continue without authentication
		if tokenString == "" {
			return c.Next()
		}

		// 2. Parse and Validate Token
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fiber.ErrUnauthorized
			}
			return []byte(m.cfg.JWTSecret), nil
		})

		if err != nil || !token.Valid {
			// Invalid token, but we allow continuation for optional auth
			return c.Next()
		}

		// 3. Set Claims to Locals
		if claims, ok := token.Claims.(jwt.MapClaims); ok {
			c.Locals("user_id", claims["user_id"])
			c.Locals("email", claims["email"])
			c.Locals("role", claims["role"])
		}

		return c.Next()
	}
}
