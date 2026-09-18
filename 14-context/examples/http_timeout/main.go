package main

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"time"
)

func main() {
	// Create a test HTTP server that simulates a slow endpoint
	slowServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		select {
		case <-r.Context().Done():
			// Client disconnected or timed out
			fmt.Println("[Server] Client canceled or timed out, aborting handler")
			return
		case <-time.After(200 * time.Millisecond):
			// Simulate slow processing
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write([]byte("response from server"))
		}
	}))
	defer slowServer.Close()

	// 1. Client request with short timeout (should fail)
	ctxShort, cancelShort := context.WithTimeout(context.Background(), 50*time.Millisecond)
	defer cancelShort()

	reqShort, err := http.NewRequestWithContext(ctxShort, http.MethodGet, slowServer.URL, nil)
	if err != nil {
		panic(err)
	}

	client := &http.Client{}
	_, err = client.Do(reqShort)
	if err != nil {
		fmt.Printf("[Client] Request 1 failed as expected: %v\n", err)
	}

	// 2. Client request with generous timeout (should succeed)
	ctxLong, cancelLong := context.WithTimeout(context.Background(), 500*time.Millisecond)
	defer cancelLong()

	reqLong, err := http.NewRequestWithContext(ctxLong, http.MethodGet, slowServer.URL, nil)
	if err != nil {
		panic(err)
	}

	resp, err := client.Do(reqLong)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	fmt.Printf("[Client] Request 2 succeeded with status: %d\n", resp.StatusCode)
}
