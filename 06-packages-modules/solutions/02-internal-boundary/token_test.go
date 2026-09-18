package token

import (
	"testing"
)

func TestSolutionToken(t *testing.T) {
	tok := GenerateToken("key123", "alice")
	if !ValidateToken("key123", "alice", tok) {
		t.Errorf("token validation failed")
	}
	if ValidateToken("key123", "bob", tok) {
		t.Errorf("token should not validate for wrong user")
	}
}
