package rot13

import (
	"errors"
	"io"
	"strings"
	"testing"
)

func TestRot13Reader(t *testing.T) {
	input := "Hello, World! 123"
	expected := "Uryyb, Jbeyq! 123"

	r := NewRot13Reader(strings.NewReader(input))
	out, err := io.ReadAll(r)
	if errors.Is(err, ErrNotImplemented) {
		t.Skip("skipping: Read is not yet implemented (implement in starter.go)")
	}
	if err != nil {
		t.Fatalf("unexpected error reading: %v", err)
	}

	if string(out) != expected {
		t.Errorf("got %q, want %q", string(out), expected)
	}

	// Double ROT13 returns original
	r2 := NewRot13Reader(strings.NewReader(expected))
	original, err := io.ReadAll(r2)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if string(original) != input {
		t.Errorf("double rot13 got %q, want %q", string(original), input)
	}
}
