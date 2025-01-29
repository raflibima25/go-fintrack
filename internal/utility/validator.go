package utility

import (
	"fmt"
	"strings"

	"github.com/go-playground/validator/v10"
)

type ValidationError struct {
	Field string
	Tag   string
	Value string
}

func FormatValidationError(err validator.ValidationErrors) []ValidationError {
	var errors []ValidationError

	for _, e := range err {
		errors = append(errors, ValidationError{
			Field: strings.ToLower(e.Field()),
			Tag:   e.Tag(),
			Value: e.Param(),
		})
	}

	return errors
}

func GetReadableErrorMessage(err ValidationError) string {
	switch err.Field {
	case "emailorusername":
		if err.Tag == "required" {
			return "Email, username or password is required"
		}
	case "password":
		if err.Tag == "required" {
			return "Email, username or password is required"
		}
		if err.Tag == "min" {
			return "Password must be at least 8 characters"
		}
	}

	switch err.Tag {
	case "required":
		return fmt.Sprintf("%s is required", err.Field)
	case "email":
		return "Invalid email format"
	case "min":
		return fmt.Sprintf("%s must be at least %s characters", err.Field, err.Value)
	case "max":
		return fmt.Sprintf("%s cannot be longer than %s characters", err.Field, err.Value)
	case "alphanum":
		return fmt.Sprintf("%s can only contain letters and numbers", err.Field)
	case "eqfield":
		return "Password and confirm password must match"
	default:
		return fmt.Sprintf("Invalid value for %s", err.Field)
	}
}
