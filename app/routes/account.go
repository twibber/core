package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/twibber/core/app/handlers/account"
	"github.com/twibber/core/app/handlers/auth"
)

func AccountRoutes(api fiber.Router) {
	// Account
	api.Get("/", account.GetSession)
	api.Patch("/password", account.UpdatePassword)
	api.Patch("/", account.UpdateProfile)

	// Verification Flow
	api.Post("/verify", auth.Verify)
	api.Post("/resend", auth.ResendCode)

	// Logout
	api.Post("/logout", account.Logout)
}
