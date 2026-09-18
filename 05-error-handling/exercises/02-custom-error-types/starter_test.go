package validation

import (
	"errors"
	"testing"
)

func TestValidationErrors(t *testing.T) {
	err := ValidateRegistration("ab", "user@domain.com", 25)
	if errors.Is(err, ErrNotImplemented) {
		t.Skip("skipping: ValidateRegistration is not yet implemented (implement in starter.go)")
	}

	var valErr *ValidationError
	if !errors.As(err, &valErr) {
		t.Fatalf("expected error to be *ValidationError, got %v", err)
	}
	if valErr.Field != "username" {
		t.Errorf("expected field 'username', got %s", valErr.Field)
	}

	// Email error
	err = ValidateRegistration("validuser", "invalidemail", 20)
	if !errors.As(err, &valErr) || valErr.Field != "email" {
		t.Errorf("expected email validation error")
	}

	// Age error
	err = ValidateRegistration("validuser", "user@domain.com", 16)
	if !errors.As(err, &valErr) || valErr.Field != "age" {
		t.Errorf("expected age validation error")
	}

	// Valid registration
	err = ValidateRegistration("validuser", "user@domain.com", 20)
	if err != nil {
		t.Errorf("expected nil error for valid input, got %v", err)
	}
}
