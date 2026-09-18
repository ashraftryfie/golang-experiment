package hasher

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"strings"
	"testing"
)

func TestCopyAndHash(t *testing.T) {
	input := "Go standard library io.Reader & io.Writer streaming is zero-allocation friendly!"
	h := sha256.Sum256([]byte(input))
	expectedHash := hex.EncodeToString(h[:])

	src := strings.NewReader(input)
	var dst bytes.Buffer

	result, err := CopyAndHash(&dst, src)
	if errors.Is(err, ErrNotImplemented) {
		t.Skip("skipping: CopyAndHash is not yet implemented (implement in starter.go)")
	}
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if result.BytesWritten != int64(len(input)) {
		t.Errorf("got %d bytes written, want %d", result.BytesWritten, len(input))
	}

	if dst.String() != input {
		t.Errorf("destination content mismatch: got %q, want %q", dst.String(), input)
	}

	if result.SHA256Hex != expectedHash {
		t.Errorf("got hash %q, want %q", result.SHA256Hex, expectedHash)
	}
}
