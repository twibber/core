package utils

import (
	"github.com/gofiber/fiber/v2"
	"github.com/twibber/core/cfg"
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
