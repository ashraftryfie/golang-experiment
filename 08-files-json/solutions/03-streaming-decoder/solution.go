package streamfilter_solution

import (
	"encoding/json"
	"fmt"
	"io"
)

// Transaction represents a monetary movement.
type Transaction struct {
	ID     string  `json:"id"`
	Amount float64 `json:"amount"`
	Sender string  `json:"sender"`
}

// FilterHighValue reads a JSON array from src, extracts transactions with Amount >= minAmount,
// and streams them into dst formatted as a valid JSON array without loading the whole list into RAM.
func FilterHighValue(dst io.Writer, src io.Reader, minAmount float64) (int, error) {
	dec := json.NewDecoder(src)

	// Read opening bracket '['
	t, err := dec.Token()
	if err != nil {
		return 0, fmt.Errorf("read opening token error: %w", err)
	}
	delim, ok := t.(json.Delim)
	if !ok || delim != '[' {
		return 0, fmt.Errorf("expected '[', got %v", t)
	}

	// Write opening '[' to dst
	if _, err := dst.Write([]byte("[")); err != nil {
		return 0, err
	}

	count := 0
	firstWritten := false

	for dec.More() {
		var tx Transaction
		if err := dec.Decode(&tx); err != nil {
			return count, fmt.Errorf("decode transaction error: %w", err)
		}

		if tx.Amount >= minAmount {
			if firstWritten {
				if _, err := dst.Write([]byte(",")); err != nil {
					return count, err
				}
			}
			encoded, err := json.Marshal(tx)
			if err != nil {
				return count, fmt.Errorf("marshal item error: %w", err)
			}
			if _, err := dst.Write(encoded); err != nil {
				return count, err
			}
			firstWritten = true
			count++
		}
	}

	// Consume closing bracket ']'
	_, err = dec.Token()
	if err != nil && err != io.EOF {
		return count, fmt.Errorf("read closing token error: %w", err)
	}

	// Write closing ']' to dst
	if _, err := dst.Write([]byte("]")); err != nil {
		return count, err
	}

	return count, nil
}
