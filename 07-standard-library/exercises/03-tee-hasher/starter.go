package hasher

import (
	"errors"
	"io"
)

var ErrNotImplemented = errors.New("TODO: implement CopyAndHash")

// StreamResult captures the outcome of a streaming copy and hash operation.
type StreamResult struct {
	BytesWritten int64
	SHA256Hex    string
}

// CopyAndHash streams data from src to dst while computing the SHA-256 digest in a single pass.
// Requirements:
// 1. Create a sha256.New() hasher.
// 2. Use io.TeeReader(src, hasher) so reading from the tee reader automatically writes to the hasher.
// 3. Use io.Copy(dst, teeReader) to stream data without buffering the entire payload into memory.
// 4. Return total BytesWritten and formatted hex hash (hex.EncodeToString(hasher.Sum(nil))).
func CopyAndHash(dst io.Writer, src io.Reader) (StreamResult, error) {
	// TODO: Implement using io.TeeReader, sha256.New(), and io.Copy
	return StreamResult{}, ErrNotImplemented
}
