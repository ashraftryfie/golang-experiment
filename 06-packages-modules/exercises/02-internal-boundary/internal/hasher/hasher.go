package hasher

import (
	"crypto/sha256"
	"encoding/hex"
)

// Hash generates a deterministic SHA-256 HMAC-like hash.
func Hash(secret, data string) string {
	h := sha256.New()
	h.Write([]byte(secret + ":" + data))
	return hex.EncodeToString(h.Sum(nil))
}
