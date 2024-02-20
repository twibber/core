package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/twibber/core/app/middleware"
	"github.com/twibber/core/cfg"
	"github.com/twibber/core/utils"
	"log/slog"
	"strings"
)

// Configure sets up the Fiber application with various middleware and routes.
func Configure() *fiber.App {
	// Create a new fiber instance
	app := fiber.New(fiber.Config{
		DisableStartupMessage: true,
		ServerHeader:          cfg.Config.Name,
		// error handler
		ErrorHandler: utils.ErrorHandler,
	})

	// log a successful start
	app.Hooks().OnListen(func(data fiber.ListenData) error {
		slog.With(
			"port", data.Port,
			"host", data.Host,
		).Info("initiated http listener")
		return nil
	})

	// Apply the CORS middleware
	app.Use(cors.New(cors.Config{
		AllowOriginsFunc: func(origin string) bool {
			return strings.Contains(origin, cfg.Config.Domain)
		},
		AllowCredentials: true,
	}))

	// Define a simple route to test the server is running
	app.Get("/", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"message": "Hello, World!",
		})
	})

	// Initiate sub-routers
	AuthRoutes(app.Group("/auth"))
	AccountRoutes(app.Group("/account", middleware.Auth(false)))

	// Return the configured app for the webserver to start listening
	return app
}
