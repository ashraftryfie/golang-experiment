package composevalidator

import (
	"errors"
	"testing"
)

func TestComposeValidator_Starter(t *testing.T) {
	composeYML := `
services:
  app:
    image: myapp:latest
  postgres:
    image: postgres:16
    environment:
      POSTGRES_PASSWORD: ${POSTGRES_PASSWORD}
    healthcheck:
      test: ["CMD-SHELL", "pg_isready"]
`
	violations, err := ValidateCompose(composeYML, []string{"app", "postgres"})
	if errors.Is(err, ErrNotImplemented) {
		t.Skip("Skipping unimplemented exercise: ValidateCompose")
	}
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(violations) != 0 {
		t.Errorf("expected 0 violations, got %d", len(violations))
	}
}
