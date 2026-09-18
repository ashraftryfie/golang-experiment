package composevalidator

import (
	"testing"
)

func TestComposeValidator_Valid(t *testing.T) {
	composeYML := `
services:
  app:
    image: app:v1
  postgres:
    image: postgres:16
    environment:
      POSTGRES_PASSWORD: ${DB_PASS}
    healthcheck:
      test: ["CMD", "pg_isready"]
`
	violations, err := ValidateCompose(composeYML, []string{"app", "postgres"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(violations) != 0 {
		t.Fatalf("expected 0 violations, got %d: %+v", len(violations), violations)
	}
}

func TestComposeValidator_Violations(t *testing.T) {
	composeYML := `
services:
  app:
    image: app:v1
  postgres_db:
    image: postgres:16
    environment:
      POSTGRES_PASSWORD: unencrypted_secret_123
`
	// Requires redis, which is missing
	violations, err := ValidateCompose(composeYML, []string{"app", "postgres_db", "redis"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	expectedRules := map[string]bool{
		RuleMissingService:      false,
		RuleHardcodedSecret:     false,
		RuleDatabaseHealthcheck: false,
	}

	for _, v := range violations {
		if _, ok := expectedRules[v.RuleID]; ok {
			expectedRules[v.RuleID] = true
		}
	}

	for rule, found := range expectedRules {
		if !found {
			t.Errorf("expected violation for rule %s, but none was reported", rule)
		}
	}
}
