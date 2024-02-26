package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/twibber/core/app/handlers/posts"
	"github.com/twibber/core/app/handlers/users"
)

func UserRoutes(api fiber.Router) {
	api.Get("/", users.ListUsers) // Get all users

	// User specific routes
	userRouter := api.Group("/:user")
	{
		userRouter.Get("/", users.GetUser)           // Get user profile
		userRouter.Get("/posts", posts.GetUserPosts) // Get user posts
	}
}
