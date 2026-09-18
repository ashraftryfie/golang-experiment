package streamfilter

import (
	"errors"
	"io"
)

var ErrNotImplemented = errors.New("TODO: implement FilterHighValue")

// Transaction represents a monetary movement.
type Transaction struct {
	ID     string  `json:"id"`
	Amount float64 `json:"amount"`
	Sender string  `json:"sender"`
}

// FilterHighValue reads a JSON array from src, extracts transactions with Amount >= minAmount,
// and streams them into dst formatted as a valid JSON array.
func FilterHighValue(dst io.Writer, src io.Reader, minAmount float64) (int, error) {
	// TODO: Use json.NewDecoder(src), token iteration, and write matching items to dst
	return 0, ErrNotImplemented
}
