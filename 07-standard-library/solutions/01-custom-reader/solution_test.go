package rot13

import (
	"io"
	"strings"
	"testing"
)

func TestSolutionRot13(t *testing.T) {
	input := "The Quick Brown Fox Jumps Over 13 Lazy Dogs."
	r := NewRot13Reader(strings.NewReader(input))
	encoded, err := io.ReadAll(r)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	r2 := NewRot13Reader(strings.NewReader(string(encoded)))
	decoded, err := io.ReadAll(r2)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if string(decoded) != input {
		t.Errorf("got %q, want %q", string(decoded), input)
	}
}
