package routes

import (
	"context"
	"time"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

type HealthChecker struct {
	db *gorm.DB
	startTime time.Time
}

func NewHealthChecker(db *gorm.DB) *HealthChecker {
	return &HealthChecker{
		db: db,
		startTime: time.Now(),
	}
}

func (h *HealthChecker) HealthCheck(c *fiber.Ctx) error {
	dbStatus := "up"
	ctx, cancel := context.WithTimeout(context.Background(), 2 * time.Second)
	defer cancel()

	var result int
	if err := h.db.WithContext(ctx).Raw("SELECT 1").Scan(&result).Error; err != nil {
		dbStatus = "down"
	}

	status := "ok"
	statusCode := fiber.StatusOK
	if dbStatus == "down" {
		status = "degraded"
		statusCode = fiber.StatusServiceUnavailable
	}

	return c.Status(statusCode).JSON(fiber.Map{
		"status": status,
		"service": "cartify-api",
		"timestamp": time.Now().Unix(),
		"uptime": time.Since(h.startTime).String(),
		"database": dbStatus,
	})
}