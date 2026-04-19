package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/savanyv/cartify/internal/utils/helpers"
	"gorm.io/gorm"
)

func RegisterRoutes(app fiber.Router, db *gorm.DB) {
	api := app.Group("/api/v1")

	healthChecker := NewHealthChecker(db)
	app.Get("/health", healthChecker.HealthCheck)

	jwtService := helpers.NewJWTService()
	bcryptService := helpers.NewBcryptService()


	productRegisterRoute(api, db, jwtService)
	cartRegisterRoute(api, db, jwtService)
	authRegisterRoute(api, db, jwtService, bcryptService)
	orderRegisterRoute(api, db, jwtService)
}
