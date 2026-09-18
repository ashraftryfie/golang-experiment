package udppingpong

import (
	"errors"
	"testing"
	"time"
)

func TestPingPong_Starter(t *testing.T) {
	server, err := StartServer("127.0.0.1:0")
	if errors.Is(err, ErrNotImplemented) {
		t.Skip("Skipping unimplemented exercise: StartServer")
	}
	if err != nil {
		t.Fatalf("unexpected error starting server: %v", err)
	}
	defer server.Close()

	resp, err := SendPing(server.Addr().String(), "test-99", 200*time.Millisecond)
	if errors.Is(err, ErrNotImplemented) {
		t.Skip("Skipping unimplemented exercise: SendPing")
	}
	if err != nil {
		t.Fatalf("unexpected ping error: %v", err)
	}
	if resp != "test-99" {
		t.Errorf("expected pong response test-99, got %s", resp)
	}
}
