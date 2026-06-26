package main

import (
	"BE/config"
	"BE/database"
	"BE/router"
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
)

func main() {
	// Load config
	config.Load()

	// Connect & migrate DB
	database.Connect()
	database.Migrate()

	// Seed data awal
	database.Seed()

	// Init Fiber
	app := fiber.New(fiber.Config{
		ErrorHandler: func(c *fiber.Ctx, err error) error {
			code := fiber.StatusInternalServerError
			if e, ok := err.(*fiber.Error); ok {
				code = e.Code
			}
			return c.Status(code).JSON(fiber.Map{"message": err.Error()})
		},
	})

	// Global middleware
	app.Use(logger.New())

	// Setup routes
	router.Setup(app)

	log.Printf("Server running on port %s", config.App.AppPort)
	log.Fatal(app.Listen(":" + config.App.AppPort))
}