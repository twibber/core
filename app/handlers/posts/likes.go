package posts

import (
	"github.com/gofiber/fiber/v2"
	"github.com/twibber/core/app/models"
	"github.com/twibber/core/db"
	"github.com/twibber/core/utils"
	"gorm.io/gorm"
)

// LikePost handles the liking of a single post by its ID.
func LikePost(c *fiber.Ctx) error {
	// Get the post by its ID.
	var post models.Post
	if err := db.DB.
		Preload("Likes").
		Where(models.Post{
			BaseModel: models.BaseModel{ID: c.Params("post")},
		}).First(&post).Error; err != nil {
		return err
	}

	// Get the current session of the user that is liking the post.
	user := c.Locals("session").(models.Session)

	// Check if the user has already liked the post.
	for _, like := range post.Likes {
		if like.LikedByID == c.Locals("session").(models.Session).Connection.UserID {
			return utils.NewError(fiber.StatusConflict, "You have already liked this post.", nil)
		}
	}

	// Create the like.
	if err := db.DB.Create(&models.Like{
		LikedByID: user.Connection.UserID,
		PostID:    post.ID,
	}).Error; err != nil {
		return err
	}

	// Return the created like.
	return c.SendStatus(fiber.StatusOK)
}

// UnlikePost handles the unliking of a single post by its ID.
func UnlikePost(c *fiber.Ctx) error {
	// Get the post by its ID.
	var post models.Post
	if err := db.DB.Where(models.Post{
		BaseModel: models.BaseModel{ID: c.Params("post")},
	}).First(&post).Error; err != nil {
		return err
	}

	// Get the current session of the user that is unliking the post.
	user := c.Locals("session").(models.Session)

	// Delete the like.
	if err := db.DB.Where(models.Like{
		LikedByID: user.Connection.UserID,
		PostID:    post.ID,
	}).Delete(&models.Like{}).Error; err != nil {
		return err
	}

	// Return the deleted like.
	return c.SendStatus(fiber.StatusOK)
}

// ListPostLikes handles the retrieval of all likes on a single post by its ID.
func ListPostLikes(c *fiber.Ctx) error {
	// Get the post by its ID.
	var post models.Post
	if err := db.DB.Where(models.Post{
		BaseModel: models.BaseModel{ID: c.Params("post")},
	}).First(&post).Error; err != nil {
		return err
	}

	// Get all likes on the post.
	var likes []models.Like
	if err := db.DB.
		// Get the user that liked the post
		Preload("LikedBy", func(db *gorm.DB) *gorm.DB {
			return db.Omit("Email") // Omit the email of the author for privacy reasons.
		}).
		Where(models.Like{PostID: post.ID}).
		Find(&likes).Error; err != nil {
		return err
	}

	// Return the likes.
	return c.JSON(likes)
}
