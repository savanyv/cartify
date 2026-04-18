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

func productRegisterRoute(app fiber.Router, db *gorm.DB, jwtService helpers.JWTService) {
	productRepo := repository.NewProductRepository(db)
	productVariantRepo := repository.NewProductVariantRepository(db)
	productUsecase := usecase.NewProductUsecase(productRepo, productVariantRepo)
	productHandler := handlers.NewProductHandler(productUsecase)

	// ==================== PUBLIC ROUTES ====================
	public := app.Group("/")
	public.Get("/products", productHandler.GetAllProducts)
	public.Get("/products/:id", productHandler.GetProductByID)

	// ==================== ADMIN ROUTES ====================
	admin := app.Group("/admin", middlewares.JWTMiddleware(jwtService), middlewares.RoleMiddleware(model.RoleAdmin))
	admin.Post("/products", productHandler.CreateProduct)
	admin.Put("/products/:id", productHandler.UpdateProduct)
	admin.Delete("/products/:id", productHandler.DeleteProduct)
	admin.Post("/products/:product_id/variants", productHandler.CreateVariant)
	admin.Put("/products/variants/:id", productHandler.UpdateVariant)
}
