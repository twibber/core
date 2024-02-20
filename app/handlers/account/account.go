package account

import (
	"github.com/gofiber/fiber/v2"
	"github.com/twibber/core/app/models"
	"github.com/twibber/core/db"
	"github.com/twibber/core/utils"
)

// GetSession returns the session of the currently authenticated user
func GetSession(c *fiber.Ctx) error {
	// Just extract the attached session from the context and return it
	return c.JSON(c.Locals("session").(models.Session))
}

// UpdateProfile Not implemented but we have the route defined
func UpdateProfile(c *fiber.Ctx) error {
	return utils.ErrNotImplemented
}

// Logout logs the user out by deleting the session from the database and clearing the auth cookie
func Logout(c *fiber.Ctx) error {
	// Get the session from the context
	session := c.Locals("session").(models.Session)

	// delete the session from the database
	if err := db.DB.Delete(&session).Error; err != nil {
		return err
	}

	// Clear the auth cookie
	utils.ClearAuth(c)

	// Return OK
	return c.SendStatus(fiber.StatusOK)
}
