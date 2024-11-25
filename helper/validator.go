package helper

import (
	"regexp"
	"strings"

	"github.com/go-playground/validator/v10"
)

var validate = validator.New()

type FieldError struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

func EmailOrPhoneValidator(input string) FieldError {
	if strings.Contains(input, "@") {
		// Validate as email
		if err := validate.Var(input, "required,email"); err != nil {
			return FieldError{
				Field:   "email",
				Message: "Invalid email address",
			}
		}
		return FieldError{}
	}

	// Regex for validating phone numbers with optional country code (+)
	phoneRegex := regexp.MustCompile(`^\+?[0-9]+$`)
	if phoneRegex.MatchString(input) {
		// Validate as a phone number
		if err := validate.Var(input, "required,numeric"); err != nil {
			return FieldError{
				Field:   "phone",
				Message: "Invalid phone number",
			}
		}
		return FieldError{}
	}

	// If input is neither email nor phone number
	return FieldError{
		Field:   "email_or_phone",
		Message: "Invalid email or phone number",
	}
}

// PasswordValidator validates the password based on common rules.
func PasswordValidator(password string) FieldError {
	// Define the validation rules
	err := validate.Var(password, "required,min=8,max=32")
	if err != nil {
		// Provide a detailed error message
		return FieldError{
			Field:   "password",
			Message: "Password must be at least 8 characters long",
		}
	}
	return FieldError{}
}
