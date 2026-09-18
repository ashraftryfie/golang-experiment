package main

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

type Claims struct {
	Subject string `json:"sub"`
	Role    string `json:"role"`
	Expires int64  `json:"exp"`
}

func createToken(claims Claims, secret []byte) (string, error) {
	headerJSON := `{"alg":"HS256","typ":"JWT"}`
	headerEnc := base64.RawURLEncoding.EncodeToString([]byte(headerJSON))

	payloadBytes, err := json.Marshal(claims)
	if err != nil {
		return "", err
	}
	payloadEnc := base64.RawURLEncoding.EncodeToString(payloadBytes)

	signedData := headerEnc + "." + payloadEnc
	h := hmac.New(sha256.New, secret)
	h.Write([]byte(signedData))
	sigEnc := base64.RawURLEncoding.EncodeToString(h.Sum(nil))

	return signedData + "." + sigEnc, nil
}

func verifyToken(tokenStr string, secret []byte) (*Claims, error) {
	parts := strings.Split(tokenStr, ".")
	if len(parts) != 3 {
		return nil, errors.New("invalid token format")
	}

	signedData := parts[0] + "." + parts[1]
	expectedSig, err := base64.RawURLEncoding.DecodeString(parts[2])
	if err != nil {
		return nil, errors.New("invalid signature encoding")
	}

	h := hmac.New(sha256.New, secret)
	h.Write([]byte(signedData))
	calculatedSig := h.Sum(nil)

	if subtle.ConstantTimeCompare(expectedSig, calculatedSig) != 1 {
		return nil, errors.New("invalid signature")
	}

	payloadBytes, err := base64.RawURLEncoding.DecodeString(parts[1])
	if err != nil {
		return nil, errors.New("invalid payload encoding")
	}

	var claims Claims
	if err := json.Unmarshal(payloadBytes, &claims); err != nil {
		return nil, errors.New("invalid claims json")
	}

	if time.Now().Unix() > claims.Expires {
		return nil, errors.New("token expired")
	}

	return &claims, nil
}

func main() {
	secret := []byte("production-super-secret-key-32bytes")
	claims := Claims{
		Subject: "user_42",
		Role:    "admin",
		Expires: time.Now().Add(1 * time.Hour).Unix(),
	}

	token, err := createToken(claims, secret)
	if err != nil {
		panic(err)
	}
	fmt.Printf("[JWT] Generated Token: %s\n", token)

	verified, err := verifyToken(token, secret)
	if err != nil {
		panic(err)
	}
	fmt.Printf("[JWT] Verified User: %s, Role: %s, Expires: %d\n", verified.Subject, verified.Role, verified.Expires)
}
