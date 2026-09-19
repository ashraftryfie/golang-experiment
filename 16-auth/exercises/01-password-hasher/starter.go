package passwordhasher

import (
	"errors"
)

var (
	ErrNotImplemented = errors.New("exercise not implemented yet")
	ErrEmptyPassword  = errors.New("password cannot be empty")
)

// HashPassword hashes a raw password using bcrypt with the specified cost factor.
func HashPassword(password string, cost int) (string, error) {
	return "", ErrNotImplemented
}

// VerifyPassword verifies whether a raw password matches the bcrypt hash.
func VerifyPassword(hashedPassword, password string) bool {
	return false
}

// NeedsRehash checks if the hashedPassword was generated with a cost strictly less than targetCost.
func NeedsRehash(hashedPassword string, targetCost int) bool {
	return false
}
