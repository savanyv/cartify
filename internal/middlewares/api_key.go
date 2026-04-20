package middlewares

import (
	"github.com/gofiber/fiber/v2"
	"github.com/savanyv/cartify/internal/utils/helpers"
)

func APIKeyMiddleware(validAPIKey string, skipPath []string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		for _, path := range skipPath {
			if c.Path() == path {
				return c.Next()
			}
		}

		apiKey := c.Get("X-API-Key")
		if apiKey == "" {
			return helpers.Unauthorized(c, "Missing X-API-Key header")
		}

		if apiKey != validAPIKey {
			return helpers.Unauthorized(c, "Invalid API Key")
		}

		c.Locals("auth_type", "api_key")
		return c.Next()
	}
}
