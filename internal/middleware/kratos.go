package middleware

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/gofiber/fiber/v2"
)

// KratosMiddleware handles authentication via Ory Kratos
type KratosMiddleware struct {
	kratosPublicURL string
	kratosAdminURL  string
}

// NewKratosMiddleware creates a new Kratos middleware instance
func NewKratosMiddleware(publicURL, adminURL string) *KratosMiddleware {
	return &KratosMiddleware{
		kratosPublicURL: publicURL,
		kratosAdminURL:  adminURL,
	}
}

// KratosSession represents a Kratos session
type KratosSession struct {
	ID       string                 `json:"id"`
	Active   bool                   `json:"active"`
	Identity KratosIdentity         `json:"identity"`
	Traits   map[string]interface{} `json:"traits"`
}

// KratosIdentity represents a Kratos identity
type KratosIdentity struct {
	ID     string                 `json:"id"`
	Traits map[string]interface{} `json:"traits"`
}

// Authenticate validates Kratos session from cookie
func (m *KratosMiddleware) Authenticate() fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Get session cookie
		sessionCookie := c.Cookies("ory_kratos_session")
		if sessionCookie == "" {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "No session cookie found",
				"code":  "UNAUTHORIZED",
			})
		}

		// Validate session with Kratos
		session, err := m.validateSession(sessionCookie)
		if err != nil {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Invalid or expired session",
				"code":  "UNAUTHORIZED",
			})
		}

		if !session.Active {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Session is not active",
				"code":  "UNAUTHORIZED",
			})
		}

		// Extract user information from traits
		email, _ := session.Identity.Traits["email"].(string)
		role, _ := session.Identity.Traits["role"].(string)

		// Default role if not set
		if role == "" {
			role = "user"
		}

		// Set user information in context
		c.Locals("identity_id", session.Identity.ID)
		c.Locals("user_email", email)
		c.Locals("user_role", role)
		c.Locals("session", session)

		return c.Next()
	}
}

// RequireRole checks if user has required role
func (m *KratosMiddleware) RequireRole(roles ...string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		userRole := c.Locals("user_role")
		if userRole == nil {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "User not authenticated",
				"code":  "UNAUTHORIZED",
			})
		}

		// Check if user has any of the required roles
		hasRole := false
		for _, role := range roles {
			if userRole == role {
				hasRole = true
				break
			}
		}

		if !hasRole {
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
				"error": fmt.Sprintf("Insufficient permissions. Required roles: %v", roles),
				"code":  "FORBIDDEN",
			})
		}

		return c.Next()
	}
}

// RequireAdmin checks if user is admin
func (m *KratosMiddleware) RequireAdmin() fiber.Handler {
	return m.RequireRole("admin")
}

// OptionalAuth validates session if present, but doesn't require it
func (m *KratosMiddleware) OptionalAuth() fiber.Handler {
	return func(c *fiber.Ctx) error {
		sessionCookie := c.Cookies("ory_kratos_session")
		if sessionCookie == "" {
			return c.Next()
		}

		session, err := m.validateSession(sessionCookie)
		if err != nil || !session.Active {
			return c.Next()
		}

		// Extract user information
		email, _ := session.Identity.Traits["email"].(string)
		role, _ := session.Identity.Traits["role"].(string)

		if role == "" {
			role = "user"
		}

		// Set user information in context
		c.Locals("identity_id", session.Identity.ID)
		c.Locals("user_email", email)
		c.Locals("user_role", role)
		c.Locals("session", session)

		return c.Next()
	}
}

// validateSession calls Kratos whoami endpoint to validate session
func (m *KratosMiddleware) validateSession(sessionToken string) (*KratosSession, error) {
	req, err := http.NewRequest("GET", m.kratosPublicURL+"/sessions/whoami", nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	// Add session cookie to request
	req.AddCookie(&http.Cookie{
		Name:  "ory_kratos_session",
		Value: sessionToken,
	})

	// Alternative: Use X-Session-Token header (recommended for APIs)
	req.Header.Set("X-Session-Token", sessionToken)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to validate session: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("session validation failed with status %d: %s", resp.StatusCode, string(body))
	}

	var session KratosSession
	if err := json.NewDecoder(resp.Body).Decode(&session); err != nil {
		return nil, fmt.Errorf("failed to decode session response: %w", err)
	}

	return &session, nil
}

// GetIdentity retrieves full identity information from Kratos Admin API
func (m *KratosMiddleware) GetIdentity(identityID string) (*KratosIdentity, error) {
	url := fmt.Sprintf("%s/admin/identities/%s", m.kratosAdminURL, identityID)

	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to get identity: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to get identity, status: %d", resp.StatusCode)
	}

	var identity KratosIdentity
	if err := json.NewDecoder(resp.Body).Decode(&identity); err != nil {
		return nil, fmt.Errorf("failed to decode identity: %w", err)
	}

	return &identity, nil
}

// ValidateToken middleware for Bearer token authentication (for API clients)
// This is useful for mobile apps or external API consumers
func (m *KratosMiddleware) ValidateToken() fiber.Handler {
	return func(c *fiber.Ctx) error {
		authHeader := c.Get("Authorization")
		if authHeader == "" {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Authorization header is required",
				"code":  "UNAUTHORIZED",
			})
		}

		// Extract token from "Bearer <token>" format
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Invalid authorization header format. Use: Bearer <session_token>",
				"code":  "UNAUTHORIZED",
			})
		}

		sessionToken := parts[1]

		// Validate session with Kratos
		session, err := m.validateSession(sessionToken)
		if err != nil {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Invalid or expired session token",
				"code":  "UNAUTHORIZED",
			})
		}

		if !session.Active {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Session is not active",
				"code":  "UNAUTHORIZED",
			})
		}

		// Extract user information
		email, _ := session.Identity.Traits["email"].(string)
		role, _ := session.Identity.Traits["role"].(string)

		if role == "" {
			role = "user"
		}

		// Set user information in context
		c.Locals("identity_id", session.Identity.ID)
		c.Locals("user_email", email)
		c.Locals("user_role", role)
		c.Locals("session", session)

		return c.Next()
	}
}
