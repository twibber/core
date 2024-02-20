package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/twibber/core/app/handlers/auth"
)

func AuthRoutes(api fiber.Router) {
	// Authentication Flow
	api.Post("/login", auth.Login)
	api.Post("/register", auth.Register)
}
