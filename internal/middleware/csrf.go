package middleware

import (
	"crypto/rand"
	"crypto/subtle"
	"encoding/base64"
	"fmt"
	"karima_store/internal/errors"
	"strings"
	"sync"
	"time"

	"github.com/gofiber/fiber/v2"
)

// CSRFConfig holds configuration for CSRF protection
type CSRFConfig struct {
	TokenLength      int           // Length of CSRF token in bytes
	TokenExpiration  time.Duration // Token expiration time
	TokenHeader      string        // Header name for CSRF token
	TokenFormField   string        // Form field name for CSRF token
	TokenContextKey  string        // Context key for storing token
	CookieName       string        // Cookie name for CSRF token
	CookieSecure     bool          // Whether cookie should be secure (HTTPS only)
	CookieHTTPOnly   bool          // Whether cookie should be HTTP only
	CookieSameSite   string        // SameSite attribute for cookie
	ExcludedPaths    []string      // Paths to exclude from CSRF protection
	TrustedOrigins   []string      // Trusted origins for CORS
}

// DefaultCSRFConfig returns default CSRF configuration
func DefaultCSRFConfig() CSRFConfig {
	return CSRFConfig{
		TokenLength:     32,
		TokenExpiration: 24 * time.Hour,
		TokenHeader:     "X-CSRF-Token",
		TokenFormField:  "csrf_token",
		TokenContextKey: "csrf_token",
		CookieName:      "csrf_token",
		CookieSecure:    true,
		CookieHTTPOnly:  false,
		CookieSameSite:  "Strict",
		ExcludedPaths: []string{
			"/api/health",
			"/api/metrics",
			"/api/swagger",
			"/api/docs",
		},
		TrustedOrigins: []string{},
	}
}

// CSRFManager manages CSRF tokens
type CSRFManager struct {
	config      CSRFConfig
	tokens      map[string]tokenInfo
	mu          sync.RWMutex
	cleanupTick *time.Ticker
	done        chan struct{}
}

type tokenInfo struct {
	value     string
	expiresAt time.Time
	createdAt time.Time
}

// NewCSRFManager creates a new CSRF manager
func NewCSRFManager(config CSRFConfig) *CSRFManager {
	if config.TokenLength == 0 {
		config = DefaultCSRFConfig()
	}

	manager := &CSRFManager{
		config: config,
		tokens: make(map[string]tokenInfo),
		done:   make(chan struct{}),
	}

	// Start cleanup goroutine
	manager.startCleanup()

	return manager
}

// startCleanup starts the cleanup goroutine to remove expired tokens
func (m *CSRFManager) startCleanup() {
	m.cleanupTick = time.NewTicker(5 * time.Minute)
	go func() {
		for {
			select {
			case <-m.cleanupTick.C:
				m.cleanupExpiredTokens()
			case <-m.done:
				return
			}
		}
	}()
}

// cleanupExpiredTokens removes expired tokens from memory
func (m *CSRFManager) cleanupExpiredTokens() {
	m.mu.Lock()
	defer m.mu.Unlock()

	now := time.Now()
	for sessionID, token := range m.tokens {
		if now.After(token.expiresAt) {
			delete(m.tokens, sessionID)
		}
	}
}

// Stop stops the cleanup goroutine
func (m *CSRFManager) Stop() {
	close(m.done)
	if m.cleanupTick != nil {
		m.cleanupTick.Stop()
	}
}

// GenerateToken generates a new CSRF token for a session
func (m *CSRFManager) GenerateToken(sessionID string) (string, error) {
	// Generate random bytes
	bytes := make([]byte, m.config.TokenLength)
	if _, err := rand.Read(bytes); err != nil {
		return "", fmt.Errorf("failed to generate CSRF token: %w", err)
	}

	// Encode to base64
	token := base64.URLEncoding.EncodeToString(bytes)

	// Store token
	m.mu.Lock()
	m.tokens[sessionID] = tokenInfo{
		value:     token,
		expiresAt: time.Now().Add(m.config.TokenExpiration),
		createdAt: time.Now(),
	}
	m.mu.Unlock()

	return token, nil
}

// ValidateToken validates a CSRF token for a session
func (m *CSRFManager) ValidateToken(sessionID, token string) bool {
	m.mu.RLock()
	defer m.mu.RUnlock()

	storedToken, exists := m.tokens[sessionID]
	if !exists {
		return false
	}

	// Check if token is expired
	if time.Now().After(storedToken.expiresAt) {
		return false
	}

	// Use constant-time comparison to prevent timing attacks
	return subtle.ConstantTimeCompare([]byte(storedToken.value), []byte(token)) == 1
}

// RotateToken rotates an existing CSRF token for a session
func (m *CSRFManager) RotateToken(sessionID string) (string, error) {
	// Generate new token
	newToken, err := m.GenerateToken(sessionID)
	if err != nil {
		return "", err
	}

	return newToken, nil
}

// RevokeToken revokes a CSRF token for a session
func (m *CSRFManager) RevokeToken(sessionID string) {
	m.mu.Lock()
	defer m.mu.Unlock()
	delete(m.tokens, sessionID)
}

// GetTokenInfo returns information about a token
func (m *CSRFManager) GetTokenInfo(sessionID string) *tokenInfo {
	m.mu.RLock()
	defer m.mu.RUnlock()

	if token, exists := m.tokens[sessionID]; exists {
		return &token
	}
	return nil
}

// CSRF creates CSRF protection middleware
func CSRF(config CSRFConfig) fiber.Handler {
	if config.TokenLength == 0 {
		config = DefaultCSRFConfig()
	}

	manager := NewCSRFManager(config)

	return func(c *fiber.Ctx) error {
		// Check if path is excluded
		if isExcludedPath(c.Path(), config.ExcludedPaths) {
			return c.Next()
		}

		// Skip CSRF for GET, HEAD, OPTIONS, TRACE methods
		if c.Method() == "GET" || c.Method() == "HEAD" || c.Method() == "OPTIONS" || c.Method() == "TRACE" {
			// Generate and set CSRF token for safe methods
			sessionID := getSessionID(c)
			if sessionID != "" {
				token, err := manager.GenerateToken(sessionID)
				if err != nil {
					return errors.NewInternalError("Failed to generate CSRF token")
				}

				// Set token in cookie
				c.Cookie(&fiber.Cookie{
					Name:     config.CookieName,
					Value:    token,
					Secure:   config.CookieSecure,
					HTTPOnly: config.CookieHTTPOnly,
					SameSite: getSameSite(config.CookieSameSite),
					MaxAge:   int(config.TokenExpiration.Seconds()),
					Path:     "/",
				})

				// Set token in context
				c.Locals(config.TokenContextKey, token)
			}
			return c.Next()
		}

		// Validate CSRF token for state-changing methods
		sessionID := getSessionID(c)
		if sessionID == "" {
			return errors.NewUnauthorizedError("No session found")
		}

		// Get token from header or form
		token := c.Get(config.TokenHeader)
		if token == "" {
			token = c.FormValue(config.TokenFormField)
		}

		if token == "" {
			return errors.NewValidationError("CSRF token is required")
		}

		// Validate token
		if !manager.ValidateToken(sessionID, token) {
			return errors.NewValidationError("Invalid CSRF token")
		}

		// Rotate token after successful validation (optional)
		// This provides additional security by preventing token reuse
		newToken, err := manager.RotateToken(sessionID)
		if err != nil {
			return errors.NewInternalError("Failed to rotate CSRF token")
		}

		// Update cookie with new token
		c.Cookie(&fiber.Cookie{
			Name:     config.CookieName,
			Value:    newToken,
			Secure:   config.CookieSecure,
			HTTPOnly: config.CookieHTTPOnly,
			SameSite: getSameSite(config.CookieSameSite),
			MaxAge:   int(config.TokenExpiration.Seconds()),
			Path:     "/",
		})

		// Set new token in context
		c.Locals(config.TokenContextKey, newToken)

		return c.Next()
	}
}

// CSRFTokenHandler returns a handler that generates and returns a CSRF token
func CSRFTokenHandler(manager *CSRFManager, config CSRFConfig) fiber.Handler {
	return func(c *fiber.Ctx) error {
		sessionID := getSessionID(c)
		if sessionID == "" {
			return errors.NewUnauthorizedError("No session found")
		}

		token, err := manager.GenerateToken(sessionID)
		if err != nil {
			return errors.NewInternalError("Failed to generate CSRF token")
		}

		// Set token in cookie
		c.Cookie(&fiber.Cookie{
			Name:     config.CookieName,
			Value:    token,
			Secure:   config.CookieSecure,
			HTTPOnly: config.CookieHTTPOnly,
			SameSite: getSameSite(config.CookieSameSite),
			MaxAge:   int(config.TokenExpiration.Seconds()),
			Path:     "/",
		})

		return c.JSON(fiber.Map{
			"csrf_token": token,
		})
	}
}

// getSessionID extracts session ID from context
func getSessionID(c *fiber.Ctx) string {
	// Try to get session ID from various sources
	// This can be customized based on your authentication system
	
	// Try from context (set by auth middleware)
	if sessionID := c.Locals("session_id"); sessionID != nil {
		if sid, ok := sessionID.(string); ok {
			return sid
		}
	}

	// Try from user ID
	if userID := c.Locals("user_id"); userID != nil {
		if uid, ok := userID.(string); ok {
			return uid
		}
	}

	// Try from header
	if sessionID := c.Get("X-Session-ID"); sessionID != "" {
		return sessionID
	}

	return ""
}

// isExcludedPath checks if a path is excluded from CSRF protection
func isExcludedPath(path string, excludedPaths []string) bool {
	for _, excluded := range excludedPaths {
		if strings.HasPrefix(path, excluded) {
			return true
		}
	}
	return false
}

// getSameSite converts string to fiber CookieSameSite
func getSameSite(sameSite string) fiber.CookieSameSite {
	switch strings.ToLower(sameSite) {
	case "strict":
		return fiber.CookieSameSiteStrictMode
	case "lax":
		return fiber.CookieSameSiteLaxMode
	case "none":
		return fiber.CookieSameSiteNoneMode
	default:
		return fiber.CookieSameSiteStrictMode
	}
}

// CSRFMiddleware creates a CSRF middleware with default configuration
func CSRFMiddleware() fiber.Handler {
	return CSRF(DefaultCSRFConfig())
}

// ValidateCSRFToken validates CSRF token from request
func ValidateCSRFToken(c *fiber.Ctx, config CSRFConfig, manager *CSRFManager) error {
	sessionID := getSessionID(c)
	if sessionID == "" {
		return errors.NewUnauthorizedError("No session found")
	}

	// Get token from header or form
	token := c.Get(config.TokenHeader)
	if token == "" {
		token = c.FormValue(config.TokenFormField)
	}

	if token == "" {
		return errors.NewValidationError("CSRF token is required")
	}

	// Validate token
	if !manager.ValidateToken(sessionID, token) {
		return errors.NewValidationError("Invalid CSRF token")
	}

	return nil
}
