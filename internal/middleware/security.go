package middleware

import (
	"github.com/gofiber/fiber/v2"
)

// SecurityHeaders adds security headers to all responses
func SecurityHeaders() fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Content Security Policy - Restrict resources the browser is allowed to load
		// Default to 'self' for production, can be relaxed for development
		c.Set("Content-Security-Policy",
			"default-src 'self'; "+
			"script-src 'self' 'unsafe-inline' 'unsafe-eval'; "+
			"style-src 'self' 'unsafe-inline'; "+
			"img-src 'self' data: https:; "+
			"font-src 'self' data:; "+
			"connect-src 'self'; "+
			"frame-ancestors 'none'; "+
			"base-uri 'self'; "+
			"form-action 'self'")

		// X-Content-Type-Options - Prevent MIME type sniffing
		c.Set("X-Content-Type-Options", "nosniff")

		// X-Frame-Options - Prevent clickjacking attacks
		c.Set("X-Frame-Options", "DENY")

		// X-XSS-Protection - Enable XSS filtering
		c.Set("X-XSS-Protection", "1; mode=block")

		// Strict-Transport-Security - Enforce HTTPS connections
		c.Set("Strict-Transport-Security", "max-age=31536000; includeSubDomains; preload")

		// Referrer-Policy - Control how much referrer information is sent
		c.Set("Referrer-Policy", "strict-origin-when-cross-origin")

		// Permissions-Policy - Control browser features and APIs
		c.Set("Permissions-Policy",
			"geolocation=(), "+
			"microphone=(), "+
			"camera=(), "+
			"payment=(), "+
			"usb=(), "+
			"magnetometer=(), "+
			"gyroscope=(), "+
			"accelerometer=()")

		// X-DNS-Prefetch-Control - Control DNS prefetching
		c.Set("X-DNS-Prefetch-Control", "off")

		// Cross-Origin-Embedder-Policy - Control cross-origin resource loading
		c.Set("Cross-Origin-Embedder-Policy", "require-corp")

		// Cross-Origin-Opener-Policy - Control cross-origin window access
		c.Set("Cross-Origin-Opener-Policy", "same-origin")

		// Cross-Origin-Resource-Policy - Control cross-origin resource sharing
		c.Set("Cross-Origin-Resource-Policy", "same-origin")

		return c.Next()
	}
}

// SecurityHeadersDevelopment adds relaxed security headers for development
func SecurityHeadersDevelopment() fiber.Handler {
	return func(c *fiber.Ctx) error {
		// More relaxed CSP for development
		c.Set("Content-Security-Policy",
			"default-src 'self'; "+
			"script-src 'self' 'unsafe-inline' 'unsafe-eval'; "+
			"style-src 'self' 'unsafe-inline'; "+
			"img-src 'self' data: https: http:; "+
			"font-src 'self' data:; "+
			"connect-src 'self' ws: wss:; "+
			"frame-ancestors 'self'")

		c.Set("X-Content-Type-Options", "nosniff")
		c.Set("X-Frame-Options", "SAMEORIGIN")
		c.Set("X-XSS-Protection", "1; mode=block")

		// Don't set HSTS in development
		c.Set("Referrer-Policy", "strict-origin-when-cross-origin")

		return c.Next()
	}
}
