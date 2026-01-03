package middleware

import (
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/csrf"
	"github.com/karima-store/internal/config"
)

// CSRFConfig holds configuration for CSRF protection
type CSRFConfig struct {
	KeyLookup       string   // KeyLookup is a string in the form of "<source>:<key>" that is used to extract token from the request
	CookieName      string   // Name of the CSRF cookie
	CookieSecure    bool     // Whether cookie should be secure (HTTPS only)
	CookieHTTPOnly  bool     // Whether cookie should be HTTP only
	CookieSameSite  string   // SameSite attribute for cookie
	Expiration      int      // Token expiration time in seconds
	ContextKey      string   // Context key for storing token
	TrustedOrigins  []string // Trusted origins for CORS
	Next            func(c *fiber.Ctx) bool // Next defines a function to skip this middleware when returned true
}

// DefaultCSRFConfig returns default CSRF configuration
func DefaultCSRFConfig() CSRFConfig {
	return CSRFConfig{
		KeyLookup:      "header:X-CSRF-Token",
		CookieName:     "csrf_token",
		CookieSecure:   true,
		CookieHTTPOnly: false,
		CookieSameSite: "Strict",
		Expiration:     24 * 60 * 60, // 24 hours in seconds
		ContextKey:     "token",
		TrustedOrigins: []string{},
		Next: nil,
	}
}

// CSRF creates CSRF protection middleware using official Fiber CSRF middleware
func CSRF(cfg CSRFConfig) fiber.Handler {
	if cfg.KeyLookup == "" {
		cfg = DefaultCSRFConfig()
	}

	// Convert our config to Fiber's CSRF config
	fiberCSRFConfig := csrf.Config{
		KeyLookup:      cfg.KeyLookup,
		CookieName:     cfg.CookieName,
		CookieSecure:   cfg.CookieSecure,
		CookieHTTPOnly: cfg.CookieHTTPOnly,
		CookieSameSite: getSameSite(cfg.CookieSameSite),
		Expiration:     time.Duration(cfg.Expiration) * time.Second,
		ContextKey:     cfg.ContextKey,
		Next:           cfg.Next,
	}

	return csrf.New(fiberCSRFConfig)
}

// CSRFMiddleware creates a CSRF middleware with default configuration
func CSRFMiddleware() fiber.Handler {
	return CSRF(DefaultCSRFConfig())
}

// CSRFFromConfig creates CSRF middleware from application config
func CSRFFromConfig(appConfig *config.Config) fiber.Handler {
	cfg := DefaultCSRFConfig()

	// Override with app config if available
	if appConfig != nil && appConfig.AppEnv != "production" {
		cfg.CookieSecure = false
		cfg.CookieSameSite = "Lax"
	}

	return CSRF(cfg)
}

// getSameSite converts string to fiber CookieSameSite
func getSameSite(sameSite string) string {
	switch sameSite {
	case "Strict":
		return "Strict"
	case "Lax":
		return "Lax"
	case "None":
		return "None"
	default:
		return "Strict"
	}
}
