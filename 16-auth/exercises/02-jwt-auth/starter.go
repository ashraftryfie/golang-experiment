package jwtauth

import (
	"errors"
)

var (
	ErrNotImplemented   = errors.New("exercise not implemented yet")
	ErrSecretTooShort   = errors.New("secret key must be at least 16 bytes")
	ErrEmptySubject     = errors.New("subject claim cannot be empty")
	ErrInvalidToken     = errors.New("invalid token format")
	ErrInvalidSignature = errors.New("invalid token signature")
	ErrTokenExpired     = errors.New("token has expired")
)

type Claims struct {
	Subject string `json:"sub"`
	Role    string `json:"role"`
	Expires int64  `json:"exp"`
}

// CreateToken creates and signs a JWT using HMAC-SHA256.
func CreateToken(claims Claims, secret []byte) (string, error) {
	return "", ErrNotImplemented
}

// VerifyToken verifies the signature and expiration of a JWT token, returning the parsed Claims.
func VerifyToken(tokenString string, secret []byte) (*Claims, error) {
	return nil, ErrNotImplemented
}
