package routes

import (
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func RegisterRoutes(app fiber.Router, db *gorm.DB) {
	// jwtService := helpers.NewJWTService()

	healthChecker := NewHealthChecker(db)
	app.Get("/health", healthChecker.HealthCheck)
}
