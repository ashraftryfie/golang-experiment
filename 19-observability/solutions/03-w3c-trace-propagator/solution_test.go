package w3ctrace

import (
	"errors"
	"net/http"
	"testing"
)

func TestTraceparent_RoundTrip(t *testing.T) {
	valid := "00-4bf92f3577b34da6a3ce929d0e0e4736-00f067aa0ba902b7-01"
	tc, err := ParseTraceparent(valid)
	if err != nil {
		t.Fatalf("ParseTraceparent failed: %v", err)
	}

	if tc.String() != valid {
		t.Fatalf("String() mismatch: got %s, want %s", tc.String(), valid)
	}
}

func TestTraceparent_ValidationErrors(t *testing.T) {
	// Unsupported version
	_, err := ParseTraceparent("01-4bf92f3577b34da6a3ce929d0e0e4736-00f067aa0ba902b7-01")
	if !errors.Is(err, ErrUnsupportedVer) {
		t.Fatalf("expected ErrUnsupportedVer, got %v", err)
	}

	// All zeros trace ID
	_, err = ParseTraceparent("00-00000000000000000000000000000000-00f067aa0ba902b7-01")
	if !errors.Is(err, ErrInvalidTraceID) {
		t.Fatalf("expected ErrInvalidTraceID, got %v", err)
	}

	// All zeros parent ID
	_, err = ParseTraceparent("00-4bf92f3577b34da6a3ce929d0e0e4736-0000000000000000-01")
	if !errors.Is(err, ErrInvalidParentID) {
		t.Fatalf("expected ErrInvalidParentID, got %v", err)
	}

	// Malformed parts
	_, err = ParseTraceparent("00-short-01")
	if !errors.Is(err, ErrInvalidFormat) {
		t.Fatalf("expected ErrInvalidFormat, got %v", err)
	}
}

func TestTraceparent_HTTPPropagation(t *testing.T) {
	tc := GenerateNewTraceContext()
	req, _ := http.NewRequest(http.MethodGet, "http://api.internal/v1/resource", nil)

	InjectTraceparent(req, tc)

	extracted, ok := ExtractTraceparent(req)
	if !ok {
		t.Fatalf("expected extraction to succeed")
	}

	if extracted.TraceID != tc.TraceID || extracted.ParentID != tc.ParentID {
		t.Errorf("extracted mismatch: got %+v, want %+v", extracted, tc)
	}
}
