package token

import (
	"github.com/ashraftryfie/golang-experiment/06-packages-modules/exercises/02-internal-boundary/internal/hasher"
)

// GenerateToken generates an auth token for user using internal hasher.
func GenerateToken(secret, userID string) string {
	// TODO: Return hashed token
	_ = hasher.Hash(secret, userID)
	return ""
}

// ValidateToken verifies whether candidate token matches computed hash.
func ValidateToken(secret, userID, candidateToken string) bool {
	// TODO: Compare candidateToken to generated hash
	return false
}
