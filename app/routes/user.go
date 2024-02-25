package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/twibber/core/app/handlers/users"
)

func UserRoutes(api fiber.Router) {
	api.Get("/", users.ListUsers) // Get all users
}
