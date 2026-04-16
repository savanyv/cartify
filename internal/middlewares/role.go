package middlewares

import (
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/savanyv/cartify/internal/model"
)

func RoleMiddleware(allowedRoles ...model.Role) fiber.Handler {
	return func(c *fiber.Ctx) error {
		userRoleStr, ok := c.Locals("role").(string)
		if !ok {
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
				"error": "Access denied",
			})
		}
		userRoleStr = strings.ToLower(userRoleStr)

		for _, role := range allowedRoles {
			if userRoleStr == string(role) {
				return c.Next()
			}
		}

		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"error": "Forbidden",
			"message": "insufficient permissions",
		})
	}
}