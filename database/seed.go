package database

import (
	"BE/internal/domain"
	"log"
	"time"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

func Seed() {
	seedAdmin()
	seedWhatsappConfig()
	seedSiteSettings()
}

func seedAdmin() {
	var count int64
	DB.Model(&domain.User{}).Count(&count)
	if count > 0 {
		return
	}

	hash, _ := bcrypt.GenerateFromPassword([]byte("admin123"), bcrypt.DefaultCost)
	admin := domain.User{
		ID:           uuid.New(),
		Name:         "Admin MCreative",
		Email:        "admin@mcreative.id",
		PasswordHash: string(hash),
		Role:         "admin",
	}
	DB.Create(&admin)
	log.Println("Admin seeded: admin@mcreative.id / admin123")
}

func seedWhatsappConfig() {
	var count int64
	DB.Model(&domain.WhatsappConfig{}).Count(&count)
	if count > 0 {
		return
	}

	DB.Create(&domain.WhatsappConfig{
		ID:            uuid.New(),
		PhoneNumber:   "628xxxxxxxxx",
		OrderTemplate: "Halo, saya ingin memesan:\n\n*{product_name}*\nHarga: {price}\n\nMohon info lebih lanjut. Terima kasih!",
		HireTemplate:  "Halo! Saya ingin menggunakan jasa MCreative. 🪵\n\n*Nama:* {name}\n*Jenis Proyek:* {project_type}\n*Budget:* {budget}\n*Deskripsi:*\n{description}",
	})
	log.Println("WhatsApp config seeded")
}

func seedSiteSettings() {
	var count int64
	DB.Model(&domain.SiteSetting{}).Count(&count)
	if count > 0 {
		return
	}

	settings := []domain.SiteSetting{
		{ID: uuid.New(), Key: "site_name", Value: "MCreative", UpdatedAt: time.Now()},
		{ID: uuid.New(), Key: "site_tagline", Value: "CNC Machine Builder & Creator", UpdatedAt: time.Now()},
		{ID: uuid.New(), Key: "site_description", Value: "Kami membangun mesin CNC dan menciptakan karya dengan mesin tersebut.", UpdatedAt: time.Now()},
		{ID: uuid.New(), Key: "owner_name", Value: "Admin MCreative", UpdatedAt: time.Now()},
		{ID: uuid.New(), Key: "owner_email", Value: "admin@mcreative.id", UpdatedAt: time.Now()},
		{ID: uuid.New(), Key: "logo_url", Value: "", UpdatedAt: time.Now()},
	}
	DB.Create(&settings)
	log.Println("Site settings seeded")
}