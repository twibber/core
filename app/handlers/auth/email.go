package auth

import (
	"github.com/gofiber/fiber/v2"
	"github.com/twibber/core/app/models"
	"github.com/twibber/core/db"
	"github.com/twibber/core/mail"
	"github.com/twibber/core/utils"
	"log/slog"
	"net/http"
	"time"
)

// LoginForm is used to parse the request body for login attempts.
type LoginForm struct {
	Email    string `json:"email" validate:"required,email,max=255"`
	Password string `json:"password" validate:"required,min=8"`
}

// RegisterForm is used to parse the request body for registration attempts.
type RegisterForm struct {
	DisplayName string `json:"display_name" validate:"required,max=512"`
	Username    string `json:"username" validate:"required,alpha,min=3,max=64,lowercase"`

	// Reuse the LoginForm struct to reuse the validation rules as well as the fields.
	LoginForm
}

// Register handles the registration of new users.
func Register(c *fiber.Ctx) error {
	// Get the request body and validate it.
	var body RegisterForm
	if err := utils.ParseAndValidate(c, &body); err != nil {
		return err
	}

	// Count all users with the same email.
	var emailCount int64
	if err := db.DB.Model(models.User{}).Where(models.User{
		Email: body.Email,
	}).Count(&emailCount).Error; err != nil {
		return err
	}

	// If the email already exists in the database, return a conflict error.
	if emailCount > 0 {
		return utils.NewError(http.StatusConflict, "The email address provided has already been registered.", &utils.ErrorDetails{
			Fields: []utils.ErrorField{
				{
					Name:   "email",
					Errors: []string{"The email address provided has already been registered."},
				},
			},
		})
	}

	// Count all users with the same email.
	var usernameCount int64
	if err := db.DB.Model(models.User{}).Where(models.User{
		Username: body.Username,
	}).Count(&usernameCount).Error; err != nil {
		return err
	}

	// If the username already exists in the database, return a conflict error.
	if usernameCount > 0 {
		return utils.NewError(http.StatusConflict, "The username provided has already been registered.", &utils.ErrorDetails{
			Fields: []utils.ErrorField{
				{
					Name:   "username",
					Errors: []string{"The username provided has already been registered."},
				},
			},
		})
	}

	// Create a password hash
	hashedPassword, err := utils.CreateHash(body.Password)
	if err != nil {
		return err
	}

	// Generate a session token
	token := utils.GenerateString(64)

	// Generate a TOTP secret
	totpSecret, err := utils.GenerateSecureRandomBase32(32)
	if err != nil {
		return err
	}

	// Generate a verification code
	code, err := utils.GenerateTOTP(totpSecret, utils.EmailVerification)
	if err != nil {
		return err
	}

	// Define the expiration time
	exp := time.Now().Add(utils.AuthDuration)

	// Create the user and the connection
	user := models.User{
		DisplayName: body.DisplayName,
		Username:    body.Username,
		Email:       body.Email,
		Connections: []models.Connection{
			{
				BaseModel: models.BaseModel{
					ID: models.ProviderEmailType.WithID(body.Email), // use the email as the connection ID with the type of connection
				},
				Password:   hashedPassword,
				Verified:   false,
				TOTPVerify: totpSecret,
				Sessions: []models.Session{
					{
						BaseModel: models.BaseModel{
							ID: token, // use the token as the session ID
						},
						ExpiresAt: exp, // use the cookie duration
					},
				},
			},
		},
	}

	// Create the user and the connection
	if err := db.DB.Create(&user).Error; err != nil {
		return err
	}

	// concurrently send verification email to user, this will not block the response to improve load times
	go func() {
		err := mail.VerifyDTO{
			Defaults: mail.Defaults{
				Email: user.Email,
				Name:  user.Username,
			},
			Code: code,
		}.Send()
		if err != nil {
			slog.With("email", user.Email).Error("failed to send verification email")
		}
	}()

	// Set the Authorization cookie
	utils.SetAuthCookie(c, token, exp)

	return c.SendStatus(http.StatusCreated)
}

func Login(c *fiber.Ctx) error {
	// Get the request body and validate it.
	var body LoginForm
	if err := utils.ParseAndValidate(c, &body); err != nil {
		return err
	}

	// Attempt to find the connection by email.
	var connection models.Connection
	if err := db.DB.Where(models.Connection{
		BaseModel: models.BaseModel{ID: models.ProviderEmailType.WithID(body.Email)},
	}).First(&connection).Error; err != nil {
		return err
	}

	// Compare the password hash with the provided password.
	match, err := utils.CompareHash(body.Password, connection.Password)
	if err != nil {
		return err
	}

	// If the password does not match, return the pre-defined error.
	if !match {
		return utils.ErrInvalidCredentials
	}

	// Generate a new session token
	token := utils.GenerateString(64)

	// Define the expiration time
	exp := time.Now().Add(utils.AuthDuration)

	// Create a new session
	if err := db.DB.Create(&models.Session{
		BaseModel: models.BaseModel{
			ID: token,
		},
		ConnectionID: connection.ID,
		ExpiresAt:    exp,
	}).Error; err != nil {
		return err
	}

	// Set the Authorization cookie
	utils.SetAuthCookie(c, token, exp)

	return c.SendStatus(http.StatusCreated)
}
