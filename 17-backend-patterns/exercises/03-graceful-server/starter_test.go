package gracefulserver

import (
	"errors"
	"net/http"
	"testing"
	"time"
)

func TestGracefulServer_Starter(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/ping", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	srv := NewGracefulServer("127.0.0.1:0", mux)
	if srv == nil {
		t.Skip("Skipping unimplemented exercise: NewGracefulServer")
	}

	err := srv.Start()
	if errors.Is(err, ErrNotImplemented) {
		t.Skip("Skipping unimplemented exercise: Start")
	}
	if err != nil {
		t.Fatalf("unexpected error starting server: %v", err)
	}

	_ = srv.Shutdown(100 * time.Millisecond)
}
