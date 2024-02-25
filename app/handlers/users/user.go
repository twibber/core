package users

import (
	"github.com/gofiber/fiber/v2"
	"github.com/twibber/core/app/models"
	"github.com/twibber/core/db"
)

func ListUsers(c *fiber.Ctx) error {
	var users []models.User
	if err := db.DB.
		Omit("Email"). // Omit the email field for security and privacy reasons
		Find(&users).Error; err != nil {
		return err
	}

	return c.JSON(users)
}

type ExtendedUser struct {
	models.User
}

type UserCounts struct {
	Followers int // Number of followers
	Following int // Number of users followed
	Posts     int // Number of posts made
	Likes     int // Received likes
}

// GetUser returns the specified user's profile
func GetUser(c *fiber.Ctx) error {
	// Get user using the ID provided in the request
	var user models.User
	if err := db.DB.
		Where(&models.User{BaseModel: models.BaseModel{ID: c.Params("user")}}).
		First(&user).Error; err != nil {
		return err
	}

	return c.JSON(user)
}

func FollowUser(c *fiber.Ctx) error {
	var user models.User
	if err := db.DB.
		Where(&models.User{BaseModel: models.BaseModel{ID: c.Params("user")}}).
		First(&user).Error; err != nil {
		return err
	}

	// Check if the user is already being followed
	var follow models.Follow
	if err := db.DB.
		Where(&models.Follow{FollowerID: c.Locals("user").(string), FollowingID: user.ID}).
		First(&follow).Error; err != nil {
		return err
	}

	// If the user is already being followed, return an error
	if follow.ID != "" {
		return fiber.NewError(fiber.StatusConflict, "You are already following this user")
	}

	// Create the follow
	if err := db.DB.Create(&models.Follow{
		FollowerID:  c.Locals("user").(string),
		FollowingID: user.ID,
	}).Error; err != nil {
		return err
	}

	return c.SendStatus(fiber.StatusOK)
}

func UnfollowUser(c *fiber.Ctx) error {
	return nil
}

func GetFollowers(c *fiber.Ctx) error {
	return nil
}

func GetFollowing(c *fiber.Ctx) error {
	return nil
}

func GetPosts(c *fiber.Ctx) error {
	return nil
}

func GetLikes(c *fiber.Ctx) error {
	return nil
}
