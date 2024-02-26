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

// GetUser returns the specified user's profile
func GetUser(c *fiber.Ctx) error {
	// Get user using the ID provided in the request
	var user models.User
	if err := db.DB.
		// Where(&models.User{BaseModel: models.BaseModel{ID: c.Params("user")}}).
		Where(models.User{Username: c.Params("user")}).
		First(&user).Error; err != nil {
		return err
	}

	return c.JSON(user)
}
