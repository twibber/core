package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/twibber/core/app/handlers/posts"
	"github.com/twibber/core/app/middleware"
)

func PostRoutes(api fiber.Router) {
	api.Get("/", posts.ListPosts)                          // Get all posts
	api.Post("/", middleware.Auth(true), posts.CreatePost) // Require authentication and a verified account to create a post

	post := api.Group("/:post")
	{
		post.Get("/", posts.GetPost)                              // Get a single post by its ID
		post.Delete("/", middleware.Auth(true), posts.DeletePost) // Require authentication and a verified account to delete a post

		replies := post.Group("/replies")
		{
			replies.Get("/", posts.ListPostReplies)                     // List all replies to a post
			replies.Post("/", middleware.Auth(true), posts.CreateReply) // Require authentication and a verified account to create a reply
		}

		likes := post.Group("/likes")
		{
			likes.Get("/", middleware.Auth(true), posts.ListPostLikes) // List all likes on a post
			likes.Post("/", middleware.Auth(true), posts.LikePost)     // Like a post
			likes.Delete("/", middleware.Auth(true), posts.UnlikePost) // Unlike a post
		}
	}
}
