package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/savanyv/cartify/internal/delivery/handlers"
	"github.com/savanyv/cartify/internal/middlewares"
	"github.com/savanyv/cartify/internal/model"
	"github.com/savanyv/cartify/internal/repository"
	"github.com/savanyv/cartify/internal/usecase"
	"github.com/savanyv/cartify/internal/utils/helpers"
	"gorm.io/gorm"
)

func orderRegisterRoute(app fiber.Router, db *gorm.DB, jwtService helpers.JWTService) {
	orderRepo := repository.NewOrderRepository(db)
	cartRepo := repository.NewCartRepository(db)
	productVariantRepo := repository.NewProductVariantRepository(db)

	orderUsecase := usecase.NewOrderUsecase(db, orderRepo, cartRepo, productVariantRepo)

	orderHandler := handlers.NewOrderHandler(orderUsecase)

	user := app.Group("/", middlewares.JWTMiddleware(jwtService))

	user.Post("/orders", orderHandler.CreateOrder)
	user.Get("/orders", orderHandler.GetUserOrders)
	user.Get("/orders/:id", orderHandler.GetOrderByID)

	admin := app.Group("/admin", middlewares.JWTMiddleware(jwtService), middlewares.RoleMiddleware(model.RoleAdmin))

	admin.Get("/orders", orderHandler.GetAllOrders)
	admin.Put("/orders/:id/status", orderHandler.UpdateOrderStatus)
}