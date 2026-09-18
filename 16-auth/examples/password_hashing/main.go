package main

import (
	"fmt"
	"time"

	"golang.org/x/crypto/bcrypt"
)

func main() {
	rawPassword := "SuperSecretPassword123!"

	// Generate bcrypt hash with cost 10
	start := time.Now()
	hash, err := bcrypt.GenerateFromPassword([]byte(rawPassword), 10)
	if err != nil {
		panic(err)
	}
	elapsed := time.Since(start)

	fmt.Printf("[Bcrypt] Hashed password in %v\n", elapsed)
	fmt.Printf("[Bcrypt] Hash: %s\n", string(hash))

	// Verify correct password
	err = bcrypt.CompareHashAndPassword(hash, []byte(rawPassword))
	if err == nil {
		fmt.Println("[Bcrypt] Password verification: SUCCESS")
	} else {
		fmt.Printf("[Bcrypt] Password verification: FAILED (%v)\n", err)
	}

	// Verify wrong password
	err = bcrypt.CompareHashAndPassword(hash, []byte("WrongPassword"))
	if err != nil {
		fmt.Println("[Bcrypt] Rejection of invalid password: SUCCESS")
	}
}
