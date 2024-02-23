package handlers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/twibber/core/app/models"
	"github.com/twibber/core/db"
	"github.com/twibber/core/utils"
	"gorm.io/gorm"
	"log/slog"
	"time"
)

type PostForm struct {
	Content string `json:"content" validate:"required,max=512"`
}

// CreatePost handles the creation of new posts.
func CreatePost(c *fiber.Ctx) error {
	// Get the request body and validate it.
	var body PostForm
	if err := utils.ParseAndValidate(c, &body); err != nil {
		return err
	}

	// Get the current session for our author.
	user := c.Locals("session").(models.Session)

	// Create the post.
	post := models.Post{
		AuthorID: user.Connection.UserID,
		Content:  body.Content,
	}
	if err := db.DB.Create(&post).Error; err != nil {
		return err
	}

	// Return the created post.
	return c.JSON(post)
}

// ListPosts handles the retrieval of all posts on the platform.
func ListPosts(c *fiber.Ctx) error {
	// Get all posts.
	var posts []models.Post
	if err := db.DB.
		Preload("Author", func(db *gorm.DB) *gorm.DB {
			return db.Omit("Email") // Omit the email of the author for privacy reasons.
		}).
		Order("created_at desc"). // Order the posts by their creation date.
		Find(&posts).Error; err != nil {
		return err
	}

	// Return the posts.
	return c.JSON(posts)
}

// GetPost handles the retrieval of a single post by its ID.
func GetPost(c *fiber.Ctx) error {
	// Get the post by its ID.
	var post models.Post
	if err := db.DB.
		Preload("Author", func(db *gorm.DB) *gorm.DB {
			return db.Omit("Email") // Omit the email of the author for privacy reasons.
		}).
		Where(models.Post{BaseModel: models.BaseModel{ID: c.Params("post")}}).
		First(&post).Error; err != nil {
		return err
	}

	// Return the post.
	return c.JSON(post)
}

// DeletePost handles the deletion of a single post by its ID as long as the author is the one making the request, and it was created within the last 5 minutes.
func DeletePost(c *fiber.Ctx) error {
	// Get the post by its ID.
	var post models.Post
	if err := db.DB.Where(models.Post{
		BaseModel: models.BaseModel{ID: c.Params("post")},
	}).First(&post).Error; err != nil {
		return err
	}

	// Get the current session for our author.
	user := c.Locals("session").(models.Session)

	// Check if the author of the post is the same as the author of the session.
	if post.AuthorID != user.Connection.UserID {
		return utils.NewError(fiber.StatusForbidden, "You are not the author of this post.", nil)
	}

	slog.With("post", post,
		"created_at", post.CreatedAt,
		"now", time.Now(),
		"since", time.Since(post.CreatedAt),
	).Debug("DeletePost")

	// Check if the post was created within the last 5 minutes.
	if time.Since(post.CreatedAt) > 5*time.Minute {
		return utils.NewError(fiber.StatusForbidden, "You can only delete posts created within the last 5 minutes.", nil)
	}

	// Delete the post.
	if err := db.DB.Delete(&post).Error; err != nil {
		return err
	}

	// Return the deleted post.
	return c.SendStatus(fiber.StatusOK)
}
