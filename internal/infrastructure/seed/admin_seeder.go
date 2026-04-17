package seed

import (
	"log"

	"github.com/savanyv/cartify/internal/model"
	"github.com/savanyv/cartify/internal/utils/helpers"
	"gorm.io/gorm"
)

func SeedAdmin(db *gorm.DB, bcryptService helpers.BcryptService) {
	var count int64
	db.Model(&model.User{}).Where("role = ?", model.RoleAdmin).Count(&count)

	if count > 0 {
		log.Println("⚠️ Admin already exists, skipping seed")
		return
	}

	hashedPassword, err := bcryptService.HashPassword("superadmin123!")
	if err != nil {
		log.Printf("Failed to hash password: %v", err)
		return
	}

	admin := model.User{
		Name: "Super Admin",
		Username: "superadmin",
		Email: "superadmin@cartify.com",
		Password: hashedPassword,
		Role: model.RoleAdmin,
	}

	if err := db.Create(&admin).Error; err != nil {
		log.Printf("Failed to create admin: %v", err)
		return
	}

	log.Println("✅ Admin created: superadmin@cartify.com / superadmin123!")
}