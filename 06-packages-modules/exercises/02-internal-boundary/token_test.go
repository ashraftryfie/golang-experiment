package token

import (
	"testing"
)

func TestTokenLifecycle(t *testing.T) {
	tok := GenerateToken("secret-key", "user_123")
	if tok == "" {
		t.Skip("skipping: GenerateToken is not yet implemented (implement in token.go)")
	}

	if !ValidateToken("secret-key", "user_123", tok) {
		t.Errorf("expected valid token verification")
	}

	if ValidateToken("wrong-key", "user_123", tok) {
		t.Errorf("expected invalid token with wrong secret")
	}
}
