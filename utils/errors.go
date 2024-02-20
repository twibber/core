package utils

import (
	"errors"
	"fmt"
	"gorm.io/gorm"
	"log/slog"
	"net/http"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/utils"
)

// Recurring Errors
var (
	ErrInternal           = NewError(http.StatusInternalServerError, "An internal server error occurred while attempting to process the request.", nil)
	ErrForbidden          = NewError(http.StatusForbidden, "You do not have permission to access the requested resource.", nil)
	ErrUnauthorised       = NewError(http.StatusUnauthorized, "You are not authorised to access this endpoint.", nil)
	ErrNotFound           = NewError(http.StatusNotFound, "The requested resource does not exist.", nil)
	ErrNotImplemented     = NewError(http.StatusNotImplemented, "A portion of this request has not been implemented.", nil)
	ErrInvalidCredentials = NewError(http.StatusUnauthorized, "Invalid credentials. Please try again.", &ErrorDetails{
		Fields: []ErrorField{
			{
				Name:   "email",
				Errors: []string{"Invalid credentials. Please try again."},
			},
			{
				Name:   "password",
				Errors: []string{"Invalid credentials. Please try again."},
			},
		},
	})
)

// Error is the structure for an error responses.
type Error struct {
	Status  int           `json:"-"`
	Code    string        `json:"code"`
	Message string        `json:"message"`
	Details *ErrorDetails `json:"details"`
}

// Error returns the data within the error for internal use.
func (e Error) Error() string {
	return fmt.Sprintf("Status: %d, Code: %s, Message: %s, Details: %v", e.Status, e.Code, e.Message, e.Details)
}

// ErrorDetails provides details about the error, such as fields but can be expanded to include more details.
type ErrorDetails struct {
	Fields []ErrorField `json:"fields"`
}

// ErrorField is a field that has an error, this is filled in by the validator.
type ErrorField struct {
	Name   string   `json:"name"`
	Errors []string `json:"errors"`
}

// NewError is used to create a new error with the given status, message, and details and optional code.
func NewError(status int, message string, details *ErrorDetails, code ...string) *Error {
	// If status is not an error status, change it to a 500 status.
	if status < 400 || status > 599 {
		status = fiber.StatusInternalServerError
	}

	// Convert the status code to a string if it's not provided.
	var statusCode string
	if len(code) > 0 {
		statusCode = code[0]
	} else {
		// Use fiber utils to get status message from status code
		statusCode = strings.ReplaceAll(strings.ToUpper(utils.StatusMessage(status)), " ", "_")
	}

	// Return the error with the given status, message, and details.
	return &Error{
		Status:  status,
		Code:    statusCode,
		Message: message,
		Details: details,
	}
}

func ErrorHandler(c *fiber.Ctx, err error) error {
	// Define the error as nil to start with.
	var e *Error = nil

	// Already handled errors, skip the rest of the function.
	if errors.As(err, &e) {
		return c.Status(e.Status).JSON(e)
	}

	// Handle GORM not found errors
	if errors.Is(err, gorm.ErrRecordNotFound) || errors.Is(err, gorm.ErrEmptySlice) {
		if strings.HasPrefix(c.Route().Path, "/auth") {
			// Custom error for the auth routes.
			e = ErrInvalidCredentials
		} else {
			// Default error for not found records.
			e = NewError(fiber.StatusNotFound, "The requested resource was not found.", nil)
		}
	}

	// Handle all fiber errors
	var fiberErr *fiber.Error
	if errors.As(err, &fiberErr) {
		e = NewError(fiberErr.Code, fiberErr.Message, nil)
	}

	// If no error struct was matched return the internal server error and log the error.
	if e == nil {
		slog.With("error", err).Error("unhandled error occurred")
		e = ErrInternal
	}

	// Finally, return the error.
	return c.Status(e.Status).JSON(e)
}
