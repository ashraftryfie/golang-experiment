package restmodels_solution

import (
	"strings"
)

type ValidationError struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

type CreateUserRequest struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Age      int    `json:"age"`
}

// Validate checks business constraints on CreateUserRequest and returns field errors.
func (r *CreateUserRequest) Validate() []ValidationError {
	var errs []ValidationError

	// Validate Username
	trimmedUser := strings.TrimSpace(r.Username)
	if len(trimmedUser) < 3 || len(trimmedUser) > 30 || strings.Contains(r.Username, " ") {
		errs = append(errs, ValidationError{
			Field:   "username",
			Message: "username must be 3-30 characters without whitespace",
		})
	}

	// Validate Email
	if !strings.Contains(r.Email, "@") || !strings.Contains(r.Email, ".") {
		errs = append(errs, ValidationError{
			Field:   "email",
			Message: "must be a valid email address",
		})
	}

	// Validate Age
	if r.Age < 18 || r.Age > 120 {
		errs = append(errs, ValidationError{
			Field:   "age",
			Message: "age must be between 18 and 120",
		})
	}

	return errs
}
