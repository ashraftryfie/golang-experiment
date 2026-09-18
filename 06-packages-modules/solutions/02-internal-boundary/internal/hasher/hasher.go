package hasher

import (
	"crypto/sha256"
	"encoding/hex"
)

func Hash(secret, data string) string {
	h := sha256.New()
	h.Write([]byte(secret + ":" + data))
	return hex.EncodeToString(h.Sum(nil))
}
