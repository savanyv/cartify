package app

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/savanyv/cartify/config"
	"github.com/savanyv/cartify/internal/delivery/routes"
	"github.com/savanyv/cartify/internal/infrastructure"
	"github.com/savanyv/cartify/internal/infrastructure/seed"
	"github.com/savanyv/cartify/internal/middlewares"
	"github.com/savanyv/cartify/internal/utils/helpers"
)

type Server struct {
	app    *fiber.App
	config *config.Config
}

func NewServer(cfg *config.Config) *Server {
	app := fiber.New(fiber.Config{
		AppName:      cfg.AppName,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  10 * time.Second,
		Prefork:      cfg.IsProduction(),
		ServerHeader: "Cartify",
	})

	return &Server{
		app:    app,
		config: cfg,
	}
}

func (s *Server) Start() error {
	db, err := infrastructure.NewDB(s.config)
	if err != nil {
		return fmt.Errorf("init database: %w", err)
	}

	s.setupMiddlewares()

	if s.config.IsDevelopment() {
		bcryptService := helpers.NewBcryptService()
		seed.SeedAdmin(db, bcryptService)
	}
	routes.RegisterRoutes(s.app, db)

	addr := fmt.Sprintf(":%s", s.config.AppPort)
	go func() {
		log.Printf("🚀 Server running on %s", addr)
		log.Printf("📝 Environment: %s", s.config.AppEnv)
		log.Printf("🔗 API Base URL: http://localhost%s", addr)
		if err := s.app.Listen(addr); err != nil {
			log.Printf("server error: %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
	<-quit

	log.Println("⏳ Shutting down server...")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := s.app.ShutdownWithContext(ctx); err != nil {
		return fmt.Errorf("shutdown server: %w", err)
	}

	log.Println("✅ Server shut down successfully")
	return nil
}

func (s *Server) setupMiddlewares() {
	skipPath := []string{"/health"}

	s.app.Use(middlewares.APIKeyMiddleware(s.config.APIKey, skipPath))
	s.app.Use(middlewares.SecurityHeadersMiddleware())
	s.app.Use(middlewares.RequestIDMiddleware())
	s.app.Use(middlewares.RecoveryMiddleware())
	s.app.Use(middlewares.LoggerMiddleware())
	s.app.Use(middlewares.CORSMiddleware())
	s.app.Use(middlewares.MethodValidationMiddleware())
	s.app.Use(middlewares.RateLimiter(100, 1*time.Minute))
}
