package selfhealth

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestSelfHealth_Starter(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
	defer ts.Close()

	err := RunHealthProbe(ts.URL, 100*time.Millisecond)
	if errors.Is(err, ErrNotImplemented) {
		t.Skip("Skipping unimplemented exercise: RunHealthProbe")
	}
	if err != nil {
		t.Fatalf("unexpected probe error: %v", err)
	}
}
