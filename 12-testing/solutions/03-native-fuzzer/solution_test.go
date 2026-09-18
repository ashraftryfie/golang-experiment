package fuzzer_solution

import (
	"testing"
)

func TestParseHeaderSolution(t *testing.T) {
	valid := []byte{1, 5, 0, 100}
	hdr, err := ParseHeader(valid)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if hdr.Version != 1 || hdr.Type != 5 || hdr.Length != 100 {
		t.Errorf("header fields mismatch: %+v", hdr)
	}

	// Too short
	_, err = ParseHeader([]byte{1, 2})
	if err != ErrTooShort {
		t.Errorf("expected ErrTooShort, got %v", err)
	}

	// Invalid version
	_, err = ParseHeader([]byte{0, 2, 0, 10})
	if err != ErrInvalidVersion {
		t.Errorf("expected ErrInvalidVersion, got %v", err)
	}
}

func FuzzParseHeaderSolution(f *testing.F) {
	f.Add([]byte{1, 2, 0, 10})
	f.Add([]byte{})
	f.Add([]byte{255, 255})
	f.Add([]byte{0, 1, 0, 0})

	f.Fuzz(func(t *testing.T, data []byte) {
		hdr, err := ParseHeader(data)
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
