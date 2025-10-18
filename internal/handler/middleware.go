package handler

import (
	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/logger"
	"github.com/gofiber/fiber/v3/middleware/recover"
)

// SetupMiddleware configures all application middleware
func SetupMiddleware(app *fiber.App) {
	// Recovery middleware to catch panics
	app.Use(recover.New(recover.Config{
		EnableStackTrace: true,
	}))

	// Logger middleware for HTTP requests
	app.Use(logger.New(logger.Config{
		Format: "[${time}] ${ip} ${status} - ${latency} ${method} ${path} ${error}\n",
	}))
}
