package database

import (
	"BE/config"
	"BE/internal/domain"
	"log"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DB *gorm.DB

func Connect() {
	db, err := gorm.Open(postgres.Open(config.App.DBDSN), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	DB = db
	log.Println("Database connected successfully")
}

func Migrate() {
	err := DB.AutoMigrate(
		&domain.User{},
		&domain.ProductCategory{},
		&domain.Product{},
		&domain.Portfolio{},
		&domain.BlogCategory{},
		&domain.Article{},
		&domain.Testimonial{},
		&domain.SiteSetting{},
		&domain.WhatsappConfig{},
		&domain.ChatSession{},
		&domain.ChatMessage{},
		&domain.OrderLog{},
	)
	if err != nil {
		log.Fatalf("Failed to migrate database: %v", err)
	}
	log.Println("Database migrated successfully")
}