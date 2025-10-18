package handler

import (
	"github.com/gofiber/fiber/v3"
)

// SetupRoutes configures all application routes
func SetupRoutes(app *fiber.App) {
	// Root routes
	app.Get("/", Welcome)

	// API v1 routes
	api := app.Group("/api/v1")

	// Health check
	api.Get("/health", HealthCheck)
}

// Welcome handler returns a welcome message
func Welcome(c fiber.Ctx) error {
	return c.JSON(fiber.Map{
		"message": "Welcome to User Management API",
		"version": "1.0.0",
	})
}

// HealthCheck handler returns the health status
func HealthCheck(c fiber.Ctx) error {
	return c.JSON(fiber.Map{
		"status": "healthy",
	})
}
