package main

import (
	"fmt"

	"github.com/ashraftryfie/golang-experiment/06-packages-modules/examples/visibility_and_internal/internal/secret"
)

func main() {
	// Calling exported function from internal package
	token := secret.EncryptToken("user-session-987")
	fmt.Println("Encrypted token via internal package:", token)

	// Note: secret.privateSalt is unexported and cannot be accessed here!
}
