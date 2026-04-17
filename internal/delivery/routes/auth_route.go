package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/savanyv/cartify/internal/delivery/handlers"
	"github.com/savanyv/cartify/internal/middlewares"
	"github.com/savanyv/cartify/internal/repository"
	"github.com/savanyv/cartify/internal/usecase"
	"github.com/savanyv/cartify/internal/utils/helpers"
	"gorm.io/gorm"
)

func authRegisterRoute(app fiber.Router, db *gorm.DB, jwtService helpers.JWTService, bcryptService helpers.BcryptService) {
	userRepo := repository.NewUserRepository(db)

	authUsecase := usecase.NewAuthUsecase(userRepo, jwtService, bcryptService)
	authHandler := handlers.NewAuthHandler(authUsecase)

	auth := app.Group("/auth")
	auth.Post("/register", authHandler.Register)
	auth.Post("/login", authHandler.Login)
	auth.Post("/refresh", authHandler.RefreshToken)

	protected := app.Group("/", middlewares.JWTMiddleware(jwtService))

	user := protected.Group("/user")
	user.Get("/profile", authHandler.GetProfile)
	user.Post("/change-password", authHandler.ChangePassword)
	user.Post("/logout", authHandler.Logout)
}