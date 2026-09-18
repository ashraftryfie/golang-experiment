package validation

import (
	"errors"
	"testing"
)

func TestSolutionValidation(t *testing.T) {
	err := ValidateRegistration("al", "valid@domain.com", 22)
	var valErr *ValidationError
	if !errors.As(err, &valErr) {
		t.Fatalf("expected *ValidationError")
	}
	if valErr.Field != "username" {
		t.Errorf("expected username error, got %s", valErr.Field)
	}

	// Valid
	if err := ValidateRegistration("alice", "alice@domain.com", 25); err != nil {
		t.Errorf("unexpected error: %v", err)
	}
}
