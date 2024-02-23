package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/twibber/core/app/handlers"
	"github.com/twibber/core/app/middleware"
)

func PostRoutes(api fiber.Router) {
	api.Get("/", handlers.ListPosts)                          // Get all posts
	api.Post("/", middleware.Auth(true), handlers.CreatePost) // Require authentication and a verified account to create a post

	post := api.Group("/:post")
	{
		post.Get("/", handlers.GetPost)                              // Get a single post by its ID
		post.Delete("/", middleware.Auth(true), handlers.DeletePost) // Require authentication and a verified account to delete a post
	}
}
