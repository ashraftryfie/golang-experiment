package fuzzer

import (
	"errors"
	"testing"
)

func TestParseHeader(t *testing.T) {
	valid := []byte{1, 5, 0, 100}
	hdr, err := ParseHeader(valid)
	if errors.Is(err, ErrNotImplemented) {
		t.Skip("skipping: ParseHeader is not yet implemented (implement in starter.go)")
	}
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if hdr.Version != 1 || hdr.Type != 5 || hdr.Length != 100 {
		t.Errorf("header fields mismatch: %+v", hdr)
	}
}

// FuzzParseHeader performs native fuzz testing ensuring ParseHeader never panics.
func FuzzParseHeader(f *testing.F) {
	// Seed corpus
	f.Add([]byte{1, 2, 0, 10})
	f.Add([]byte{})
	f.Add([]byte{255, 255})
	f.Add([]byte{0, 1, 0, 0})

	f.Fuzz(func(t *testing.T, data []byte) {
		hdr, err := ParseHeader(data)
		if errors.Is(err, ErrNotImplemented) {
			t.Skip("skipping: ParseHeader is not yet implemented")
		}
		// Invariant: If err is nil, version must be >= 1 and len(data) >= 4
		if err == nil {
			if len(data) < 4 {
				t.Errorf("accepted packet with len < 4: %d", len(data))
			}
			if hdr.Version == 0 {
				t.Errorf("accepted packet with version 0")
			}
		}
	})
}
