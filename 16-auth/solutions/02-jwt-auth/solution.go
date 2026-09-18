package jwtauth

import (
	"crypto/hmac"
	"crypto/sha256"
	"crypto/subtle"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"time"
)

var (
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
	if len(secret) < 16 {
		return "", ErrSecretTooShort
	}
	if strings.TrimSpace(claims.Subject) == "" {
		return "", ErrEmptySubject
	}

	headerJSON := `{"alg":"HS256","typ":"JWT"}`
	headerEnc := base64.RawURLEncoding.EncodeToString([]byte(headerJSON))

	payloadBytes, err := json.Marshal(claims)
	if err != nil {
		return "", fmt.Errorf("failed to marshal claims: %w", err)
	}
	payloadEnc := base64.RawURLEncoding.EncodeToString(payloadBytes)

	signedData := headerEnc + "." + payloadEnc
	h := hmac.New(sha256.New, secret)
	h.Write([]byte(signedData))
	sigEnc := base64.RawURLEncoding.EncodeToString(h.Sum(nil))

	return signedData + "." + sigEnc, nil
}

// VerifyToken verifies the signature and expiration of a JWT token, returning the parsed Claims.
func VerifyToken(tokenString string, secret []byte) (*Claims, error) {
	if len(secret) < 16 {
		return nil, ErrSecretTooShort
	}

	parts := strings.Split(tokenString, ".")
	if len(parts) != 3 {
		return nil, ErrInvalidToken
	}

	signedData := parts[0] + "." + parts[1]
	expectedSig, err := base64.RawURLEncoding.DecodeString(parts[2])
	if err != nil {
		return nil, ErrInvalidToken
	}

	h := hmac.New(sha256.New, secret)
	h.Write([]byte(signedData))
	calculatedSig := h.Sum(nil)

	if subtle.ConstantTimeCompare(expectedSig, calculatedSig) != 1 {
		return nil, ErrInvalidSignature
	}

	payloadBytes, err := base64.RawURLEncoding.DecodeString(parts[1])
	if err != nil {
		return nil, ErrInvalidToken
	}

	var claims Claims
	if err := json.Unmarshal(payloadBytes, &claims); err != nil {
		return nil, ErrInvalidToken
	}

	if claims.Expires > 0 && time.Now().Unix() > claims.Expires {
		return nil, ErrTokenExpired
	}

	return &claims, nil
}
