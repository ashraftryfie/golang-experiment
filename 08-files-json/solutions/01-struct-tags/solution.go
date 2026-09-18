package useraccount_solution

import (
	"encoding/json"
	"fmt"
)

// UserAccount models a registered system user with idiomatic JSON tags.
type UserAccount struct {
	ID           int64  `json:"id"`
	Username     string `json:"username"`
	Email        string `json:"email"`
	Role         string `json:"role"`
	PhoneNumber  string `json:"phone_number,omitempty"`
	IsActive     bool   `json:"is_active"`
	PasswordHash string `json:"-"`
	APISecret    string `json:"-"`
}

// SerializeUser converts a UserAccount to formatted JSON bytes.
func SerializeUser(u UserAccount) ([]byte, error) {
	data, err := json.Marshal(u)
	if err != nil {
		return nil, fmt.Errorf("serialize user failed: %w", err)
	}
	return data, nil
}
