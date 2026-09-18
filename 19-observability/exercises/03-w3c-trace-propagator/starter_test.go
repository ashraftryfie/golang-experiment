package w3ctrace

import (
	"errors"
	"net/http"
	"testing"
)

func TestTracePropagator_Starter(t *testing.T) {
	raw := "00-4bf92f3577b34da6a3ce929d0e0e4736-00f067aa0ba902b7-01"
	tc, err := ParseTraceparent(raw)
	if errors.Is(err, ErrNotImplemented) {
		t.Skip("Skipping unimplemented exercise: ParseTraceparent")
	}
	if err != nil {
		t.Fatalf("unexpected parse error: %v", err)
	}

	req, _ := http.NewRequest(http.MethodGet, "http://example.com", nil)
	InjectTraceparent(req, *tc)

	extracted, ok := ExtractTraceparent(req)
	if !ok || extracted.TraceID != "4bf92f3577b34da6a3ce929d0e0e4736" {
		t.Errorf("extracted trace mismatch: %+v", extracted)
	}
}
