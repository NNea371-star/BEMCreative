package config

import (
	"fmt"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	AppPort             string
	AppEnv              string
	DBHost              string
	DBPort              string
	DBUser              string
	DBPass              string
	DBName              string
	DBDSN               string
	JWTSecret           string
	JWTExpireHours      int
	StoragePath         string
	CloudinaryCloudName string
	CloudinaryAPIKey    string
	CloudinaryAPISecret string
}

var App *Config

func Load() {
	_ = godotenv.Load()

	expireHours, _ := strconv.Atoi(getEnv("JWT_EXPIRE_HOURS", "24"))

	App = &Config{
		AppPort:             getEnv("APP_PORT", "8080"),
		AppEnv:              getEnv("APP_ENV", "development"),
		DBHost:              getEnv("DB_HOST", "localhost"),
		DBPort:              getEnv("DB_PORT", "5432"),
		DBUser:              getEnv("DB_USER", "postgres"),
		DBPass:              getEnv("DB_PASS", "postgres"),
		DBName:              getEnv("DB_NAME", "cnc_mcreative"),
		JWTSecret:           getEnv("JWT_SECRET", "secret"),
		JWTExpireHours:      expireHours,
		StoragePath:         getEnv("STORAGE_PATH", "./storage"),
		CloudinaryCloudName: getEnv("CLOUDINARY_CLOUD_NAME", ""),
		CloudinaryAPIKey:    getEnv("CLOUDINARY_API_KEY", ""),
		CloudinaryAPISecret: getEnv("CLOUDINARY_API_SECRET", ""),
	}

	// ✅ Prioritaskan DATABASE_URL dari Railway jika tersedia
	if dbURL := os.Getenv("DATABASE_URL"); dbURL != "" {
		App.DBDSN = dbURL
	} else {
		// Fallback: bangun DSN dari komponen (untuk development lokal)
		sslMode := "disable"
		if App.AppEnv == "production" {
			sslMode = "require"
		}

		App.DBDSN = fmt.Sprintf(
			"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s TimeZone=Asia/Jakarta",
			App.DBHost, App.DBPort, App.DBUser, App.DBPass, App.DBName, sslMode,
		)
	}
}

func getEnv(key, fallback string) string {
	if val := os.Getenv(key); val != "" {
		return val
	}
	return fallback
}