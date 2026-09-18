package validator

import (
	"errors"
	"strings"
)

var (
	ErrEmptyID   = errors.New("id cannot be empty")
	ErrEmptyName = errors.New("name cannot be empty")
	ErrPriceZero = errors.New("price must be greater than zero")
)

// ValidateProduct checks product field constraints.
func ValidateProduct(id, name string, price float64) error {
	if strings.TrimSpace(id) == "" {
		return ErrEmptyID
	}
	if strings.TrimSpace(name) == "" {
		return ErrEmptyName
	}
	if price <= 0 {
		return ErrPriceZero
	}
	return nil
}
