package account

import (
	"github.com/gofiber/fiber/v2"
	"github.com/twibber/core/app/models"
	"github.com/twibber/core/db"
	"github.com/twibber/core/utils"
)

type UpdatePasswordDTO struct {
	CurrentPassword string `json:"password" validate:"required"`
	Password        string `json:"new_password" validate:"required"`
}

func UpdatePassword(c *fiber.Ctx) error {
	session := c.Locals("session").(models.Session)
	connection := session.Connection

	var dto UpdatePasswordDTO
	if err := utils.ParseAndValidate(c, &dto); err != nil {
		return err
	}

	match, err := utils.CompareHash(dto.CurrentPassword, connection.Password)
	if err != nil {
		return err
	}

	if !match {
		return utils.NewError(fiber.StatusBadRequest, "The current password provided is incorrect.", &utils.ErrorDetails{
			Fields: []utils.ErrorField{
				{
					Name:   "password",
					Errors: []string{"The current password provided is incorrect."},
				},
			},
		})
	}

	hash, err := utils.CreateHash(dto.Password)
	if err != nil {
		return err
	}
	connection.Password = hash

	if err := db.DB.Save(&connection).Error; err != nil {
		return err
	}

	return c.SendStatus(fiber.StatusOK)
}
