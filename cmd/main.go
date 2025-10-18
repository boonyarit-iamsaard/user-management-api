package main

import (
	"log"

	"github.com/boonyarit-iamsaard/user-management-api/internal/handler"
	"github.com/gofiber/fiber/v3"
)

func main() {
	// Create a new Fiber instance
	app := fiber.New(fiber.Config{
		// Server configuration
	})

	// Setup middleware
	handler.SetupMiddleware(app)

	// Setup routes
	handler.SetupRoutes(app)

	// Start server
	log.Println("Starting server on :3000")
	if err := app.Listen(":3000"); err != nil {
		log.Fatal("Server failed to start:", err)
	}
}
