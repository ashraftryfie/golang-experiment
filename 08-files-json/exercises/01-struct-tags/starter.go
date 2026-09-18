package useraccount

import (
	"errors"
)

var ErrNotImplemented = errors.New("TODO: implement struct tags and SerializeUser")

// UserAccount models a registered system user.
// Add appropriate `json:"..."` tags to satisfy:
// - ID -> "id"
// - Username -> "username"
// - Email -> "email"
// - Role -> "role"
// - PhoneNumber -> "phone_number,omitempty"
// - IsActive -> "is_active"
// - PasswordHash -> ignore from JSON ("-")
// - APISecret -> ignore from JSON ("-")
type UserAccount struct {
	ID           int64
	Username     string
	Email        string
	Role         string
	PhoneNumber  string
	IsActive     bool
	PasswordHash string
	APISecret    string
}

// SerializeUser converts a UserAccount to formatted JSON bytes.
func SerializeUser(u UserAccount) ([]byte, error) {
	// TODO: Use json.Marshal
	return nil, ErrNotImplemented
}
