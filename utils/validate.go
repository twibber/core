package utils

import (
	"errors"
	"reflect"
	"strings"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

// ParseAndValidate parses the request body into the provided struct and validates it.
// Returns a detailed error if validation fails.
func ParseAndValidate(c *fiber.Ctx, body any) error {
	// Parse the body into the provided struct.
	if err := c.BodyParser(body); err != nil {
		return err
	}

	// Start a new validator instance.
	v := validator.New()

	// Validate the struct.
	if err := v.Struct(body); err != nil {
		// Check if the error is a validation error.
		var validationErrors validator.ValidationErrors
		if errors.As(err, &validationErrors) {
			// Loop through the validation errors and create a detailed error response.
			details := &ErrorDetails{Fields: []ErrorField{}}
			for _, ve := range validationErrors {
				// Get the JSON tag name for the field and add it to the details.
				fieldName := getJSONTagName(body, ve.StructField())
				details.Fields = append(details.Fields, ErrorField{
					Name:   fieldName,
					Errors: []string{ve.Error()},
				})
			}

			// Return a new error with the details.
			return NewError(fiber.StatusBadRequest, "Validation failed for one or more fields.", details)
		}

		// Return the original error if it's not a validation error.
		return err
	}

	// No error, the request body is valid, and can continue.
	return nil
}

// getJSONTagName extracts the JSON tag name from the struct field. If not present, returns the field name.
func getJSONTagName(obj any, fieldName string) string {
	v := reflect.ValueOf(obj)
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}
	field, found := v.Type().FieldByName(fieldName)
	if !found {
		return fieldName
	}
	tag := field.Tag.Get("json")
	if tag == "" || tag == "-" {
		return fieldName
	}
	return strings.Split(tag, ",")[0]
}
