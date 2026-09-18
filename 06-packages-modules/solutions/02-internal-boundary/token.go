package token

import (
	"github.com/ashraftryfie/golang-experiment/06-packages-modules/solutions/02-internal-boundary/internal/hasher"
)

func GenerateToken(secret, userID string) string {
	return hasher.Hash(secret, userID)
}

func ValidateToken(secret, userID, candidateToken string) bool {
	expected := hasher.Hash(secret, userID)
	return expected == candidateToken
}
