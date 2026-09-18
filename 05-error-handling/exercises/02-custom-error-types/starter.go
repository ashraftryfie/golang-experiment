package validation

import (
	"errors"
)

var ErrNotImplemented = errors.New("TODO: implement")

type ValidationError struct {
	Field   string
	Message string
	Value   any
}

func (e *ValidationError) Error() string {
	// TODO: Format descriptive error string
	return ""
}

func ValidateRegistration(username, email string, age int) error {
	// TODO: Validate fields and return appropriate *ValidationError
	return ErrNotImplemented
}
