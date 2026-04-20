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

func cartRegisterRoute(app fiber.Router, db *gorm.DB, jwtService helpers.JWTService) {
	cartRepo := repository.NewCartRepository(db)
	productVariantRepo := repository.NewProductVariantRepository(db)

	cartUsecase := usecase.NewCartUsecase(cartRepo, productVariantRepo)
	cartHandler := handlers.NewCartHandler(cartUsecase)

	cart := app.Group("/cart", middlewares.JWTMiddleware(jwtService))

	cart.Get("/", cartHandler.GetCart)
	cart.Post("/", cartHandler.AddToCart)
	cart.Put("/items/:item_id", cartHandler.UpdateCartItem)
	cart.Delete("/items/:item_id", cartHandler.RemoveFromCart)
	cart.Delete("/clear", cartHandler.ClearCart)
}