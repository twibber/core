package middleware

import (
	"github.com/gofiber/fiber/v2"
	"github.com/twibber/core/app/models"
	"github.com/twibber/core/db"
	"github.com/twibber/core/utils"
	"log/slog"
	"time"
)

func Auth(verify bool) fiber.Handler {
	return func(c *fiber.Ctx) error {
		authCookie := c.Cookies(utils.AuthCookieName)

		// if no auth cookie is found, return unauthorised
		if authCookie == "" {
			slog.Debug("No auth cookie found")
			return utils.ErrUnauthorised
		}

		var session models.Session
		if err := db.DB.Where(models.Session{
			BaseModel: models.BaseModel{
				ID: authCookie,
			},
		}).
			Preload("Connection").
			Preload("Connection.User").
			First(&session).
			Error; err != nil {
			// even if it is another error, we will return unauthorised and clear the cookie
			utils.ClearAuth(c)
			return utils.ErrUnauthorised
		}

		// Check if token is expired
		if time.Now().After(session.ExpiresAt) {
			slog.With("session", session).Debug("Session expired")

			utils.ClearAuth(c)
			return utils.ErrUnauthorised
		}

		// check if user is verified if the action requires a verified user
		if !verify || session.Connection.Verified {
			// attach the session to the context
			c.Locals("session", session)

			// continue to the next handler
			return c.Next()
		} else {
			// return an error if the user is not verified and the action requires a verified user
			return utils.NewError(fiber.StatusForbidden,
				"You must verify your email address before performing this action.",
				nil,
				"UNVERIFIED",
			)
		}
	}
}
