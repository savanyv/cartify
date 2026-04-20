package middlewares

import "github.com/gofiber/fiber/v2"

func SecurityHeadersMiddleware() fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Prevent MIME type sniffing
		c.Set("X-Content-Type-Options", "nosniff")

		// Prevent clickjacking
		c.Set("X-Frame-Options", "DENY")

		// Enable XSS Protection
		c.Set("X-XSS-Protection", "1; mode=block")

		// HTTP Strict Transport Security (HSTS)
		c.Set("Strict-Transport-Security", "max=age=31536000; includeSubDomains; preload")

		// Referrer Policy
		c.Set("Referrer-Policy", "strict-origin-when-cross-origin")

		// Content Security Policy (CSP) - basic
		c.Set("Content-Security-Policy", "default-src 'self'")

		// Premissions policy
		c.Set("Permissions-Policy", "geolocation=(), microphone=(), camera=()")

		return c.Next()
	}
}