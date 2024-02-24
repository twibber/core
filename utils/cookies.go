package utils

import (
	"github.com/gofiber/fiber/v2"
	"github.com/twibber/core/app/models"
	"github.com/twibber/core/cfg"
	"github.com/twibber/core/db"
	"time"
)

// AuthCookieName is the name of the cookie used to store the session token.
const AuthCookieName = "Authorization"

// AuthDuration is the duration of the Authorization cookie it's also used for the token expiration.
const AuthDuration = time.Hour * 24 * 7 // 1 week

// SetAuthCookie sets the Authorization cookie with the token and the duration.
func SetAuthCookie(c *fiber.Ctx, token string, expiration time.Time) {
	c.Cookie(&fiber.Cookie{
		Name:   AuthCookieName,
		Value:  token,
		Path:   "/",
		Domain: cfg.Config.Domain,
		// MaxAge:   AuthDuration,
		Expires:  expiration,
		HTTPOnly: true,
		SameSite: "lax",
	})
}

// ClearAuth clears the Authorization cookie by setting the MaxAge to 0 and replacing the value with an empty string.
func ClearAuth(c *fiber.Ctx) {
	c.Cookie(&fiber.Cookie{
		Name:     AuthCookieName,
		Value:    "",
		Path:     "/",
		Domain:   cfg.Config.Domain,
		MaxAge:   0,
		HTTPOnly: true,
		SameSite: "lax",
	})
}

// GetUserID returns the session from the Authorization cookie.
func GetUserID(c *fiber.Ctx) string {
	// Check if session is already attached to the context
	session, ok := c.Locals("session").(models.Session)

	// If session is already attached to the context, return the user id
	if ok {
		return session.Connection.UserID
	}

	// If session is not attached to the context, get it from the cookie
	authCookie := c.Cookies(AuthCookieName)
	if authCookie == "" {
		return ""
	}

	// Get the session from the database using the cookie with the connection preloaded to get the user id from
	if err := db.DB.Where(models.Session{
		BaseModel: models.BaseModel{
			ID: authCookie,
		},
	}).
		Preload("Connection").
		First(&session).
		Error; err != nil {
		return ""
	}

	// Check if token is expired
	if time.Now().After(session.ExpiresAt) {
		return ""
	}

	// return the user id
	return session.Connection.UserID
}
