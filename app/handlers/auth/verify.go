package auth

import (
	"github.com/gofiber/fiber/v2"
	"github.com/twibber/core/app/models"
	"github.com/twibber/core/db"
	"github.com/twibber/core/mail"
	"github.com/twibber/core/utils"
	"log/slog"
)

type VerifyForm struct {
	Code string `json:"code" validate:"required,len=6"`
}

func Verify(c *fiber.Ctx) error {
	// Get the user session from the context.
	var session = c.Locals("session").(models.Session)

	// Get the connection from the session. We will use this to update the verified status.
	connection := session.Connection

	// Parse the body and validate it.
	var body VerifyForm
	if err := utils.ParseAndValidate(c, &body); err != nil {
		return err
	}

	// Validate the code provided.
	if !utils.ValidateTOTP(connection.TOTPVerify, body.Code, utils.EmailVerification) {
		// Return an error if the code is invalid.
		return utils.NewError(fiber.StatusBadRequest, "Invalid code provided.", &utils.ErrorDetails{
			Fields: []utils.ErrorField{
				{
					Name:   "code",
					Errors: []string{"The code provided is invalid."},
				},
			},
		})
	} else {
		// Mark the connection as verified.
		connection.Verified = true
	}

	// Update the connection in the database.
	if err := db.DB.Updates(&connection).Error; err != nil {
		return err
	}

	// Return a successful response.
	return c.SendStatus(fiber.StatusOK)
}

// ResendCode is a data structure for resending verification emails,
// we don't require any input as the user should be logged in to do this.
func ResendCode(c *fiber.Ctx) error {
	// Get the user session from the context.
	var session = c.Locals("session").(models.Session)

	// Generate a new code.
	code, err := utils.GenerateTOTP(session.Connection.TOTPVerify, utils.EmailVerification)
	if err != nil {
		return err
	}

	// concurrently send verification email to user
	go func() {
		err := mail.VerifyDTO{
			Defaults: mail.Defaults{
				Email: session.Connection.User.Email,
				Name:  session.Connection.User.Username,
			},
			Code: code,
		}.Send()
		if err != nil {
			slog.With("email", session.Connection.User.Email).Error("failed to send verification email")
		}
	}()

	// Return a successful response.
	return c.SendStatus(fiber.StatusOK)
}
