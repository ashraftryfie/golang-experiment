package hasher_solution

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
)

// StreamResult captures the outcome of a streaming copy and hash operation.
type StreamResult struct {
	BytesWritten int64
	SHA256Hex    string
}

// CopyAndHash streams data from src to dst while computing the SHA-256 digest in a single pass.
func CopyAndHash(dst io.Writer, src io.Reader) (StreamResult, error) {
	hasher := sha256.New()
	// io.TeeReader writes to hasher everything read from src
	tee := io.TeeReader(src, hasher)

	written, err := io.Copy(dst, tee)
	if err != nil {
		return StreamResult{}, fmt.Errorf("copy stream error: %w", err)
	}

	hashHex := hex.EncodeToString(hasher.Sum(nil))
	return StreamResult{
		BytesWritten: written,
		SHA256Hex:    hashHex,
	}, nil
}
