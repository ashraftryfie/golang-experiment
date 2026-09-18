package gracefulserver

import (
	"context"
	"errors"
	"net/http"
	"sync/atomic"
	"testing"
	"time"
)

func TestGracefulServer_Lifecycle(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/work", func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(50 * time.Millisecond)
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("ok"))
	})

	srv := NewGracefulServer("127.0.0.1:0", mux)
	if err := srv.Start(); err != nil {
		t.Fatalf("Start failed: %v", err)
	}

	addr := srv.Addr()

	// Track background worker
	var taskExecuted int32
	var taskCanceled int32

	srv.TrackTask(func(ctx context.Context) {
		atomic.StoreInt32(&taskExecuted, 1)
		<-ctx.Done()
		atomic.StoreInt32(&taskCanceled, 1)
	})

	// Fire an HTTP request in goroutine
	reqDone := make(chan struct{})
	var statusCode int
	go func() {
		defer close(reqDone)
		resp, err := http.Get("http://" + addr + "/work")
		if err == nil {
			statusCode = resp.StatusCode
			resp.Body.Close()
		}
	}()

	// Allow request to land and background task to start
	time.Sleep(20 * time.Millisecond)

	// Initiate graceful shutdown with 200ms timeout
	if err := srv.Shutdown(200 * time.Millisecond); err != nil {
		t.Fatalf("Shutdown failed: %v", err)
	}

	<-reqDone
	if statusCode != http.StatusOK {
		t.Errorf("expected 200 OK for in-flight request, got %d", statusCode)
	}

	if atomic.LoadInt32(&taskExecuted) != 1 {
		t.Errorf("background task was not executed")
	}
	if atomic.LoadInt32(&taskCanceled) != 1 {
		t.Errorf("background task was not notified of cancellation")
	}

	// Repeated shutdown should return error
	err := srv.Shutdown(50 * time.Millisecond)
	if !errors.Is(err, ErrServerNotRunning) {
		t.Fatalf("expected ErrServerNotRunning, got %v", err)
	}
}
