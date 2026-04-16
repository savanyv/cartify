package infrastructure

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/savanyv/cartify/config"
	"github.com/savanyv/cartify/internal/model"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func NewDB(cfg *config.Config) (*gorm.DB, error) {
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=UTC",
		cfg.DBHost, cfg.DBUser, cfg.DBPassword, cfg.DBName, cfg.DBPort,
	)

	gormLogger := setupGORMLogger(cfg)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: gormLogger,
	})
	if err != nil {
		log.Fatal("failed to connect to database: ", err)
	}

	log.Println("✅ Database connected")

	// Auto migrate
	if err := db.AutoMigrate(
		&model.User{},
		&model.Product{},
		&model.ProductVariant{},
		&model.Cart{},
		&model.CartItem{},
		&model.Order{},
		&model.OrderItem{},
	); err != nil {
		log.Fatal("Migration failed:", err)
	}

	log.Println("✅ Migration completed")

	return db, nil
}

func setupGORMLogger(cfg *config.Config) logger.Interface {
	logLevel := logger.Silent

	switch cfg.AppEnv {
	case "development":
		logLevel = logger.Info
	case "staging":
		logLevel = logger.Warn
	case "production":
		logLevel = logger.Error
	default:
		logLevel = logger.Info
	}

	return logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags),
		logger.Config{
			SlowThreshold:             200 * time.Millisecond,
			LogLevel:                  logLevel,
			IgnoreRecordNotFoundError: true,
			ParameterizedQueries:      cfg.IsProduction(),
			Colorful:                  cfg.IsDevelopment(),
		},
	)
}
