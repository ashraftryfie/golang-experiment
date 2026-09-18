package validation

import (
	"fmt"
	"strings"
)

type ValidationError struct {
	Field   string
	Message string
	Value   any
}

func (e *ValidationError) Error() string {
	return fmt.Sprintf("validation failed on field '%s' (value: %v): %s", e.Field, e.Value, e.Message)
}

func ValidateRegistration(username, email string, age int) error {
	if len(username) < 3 {
		return &ValidationError{
			Field:   "username",
			Message: "must be at least 3 characters",
			Value:   username,
		}
	}
	if !strings.Contains(email, "@") {
		return &ValidationError{
			Field:   "email",
			Message: "must contain '@'",
			Value:   email,
		}
	}
	if age < 18 {
		return &ValidationError{
			Field:   "age",
			Message: "must be at least 18 years old",
			Value:   age,
		}
	}
	return nil
}
