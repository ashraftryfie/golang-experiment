package passwordhasher

import (
	"errors"

	"golang.org/x/crypto/bcrypt"
)

var (
	ErrEmptyPassword = errors.New("password cannot be empty")
)

// HashPassword hashes a raw password using bcrypt with the specified cost factor.
func HashPassword(password string, cost int) (string, error) {
	if password == "" {
		return "", ErrEmptyPassword
	}

	if cost < bcrypt.MinCost || cost > bcrypt.MaxCost {
		cost = bcrypt.DefaultCost
	}

	hashBytes, err := bcrypt.GenerateFromPassword([]byte(password), cost)
	if err != nil {
		return "", err
	}

	return string(hashBytes), nil
}

// VerifyPassword verifies whether a raw password matches the bcrypt hash.
func VerifyPassword(hashedPassword, password string) bool {
	if hashedPassword == "" || password == "" {
		return false
	}
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	return err == nil
}

// NeedsRehash checks if the hashedPassword was generated with a cost strictly less than targetCost.
func NeedsRehash(hashedPassword string, targetCost int) bool {
	cost, err := bcrypt.Cost([]byte(hashedPassword))
	if err != nil {
		return true
	}
	return cost < targetCost
}
