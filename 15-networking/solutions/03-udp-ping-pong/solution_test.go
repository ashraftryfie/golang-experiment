package udppingpong

import (
	"fmt"
	"sync"
	"testing"
	"time"
)

func TestPingPong_Success(t *testing.T) {
	server, err := StartServer("127.0.0.1:0")
	if err != nil {
		t.Fatalf("StartServer failed: %v", err)
	}
	defer server.Close()

	resp, err := SendPing(server.Addr().String(), "roundtrip-1", 500*time.Millisecond)
	if err != nil {
		t.Fatalf("SendPing failed: %v", err)
	}

	if resp != "roundtrip-1" {
		t.Fatalf("expected roundtrip-1, got %s", resp)
	}
}

func TestPingPong_ConcurrentPings(t *testing.T) {
	server, err := StartServer("127.0.0.1:0")
	if err != nil {
		t.Fatalf("StartServer failed: %v", err)
	}
	defer server.Close()

	var wg sync.WaitGroup
	serverAddr := server.Addr().String()

	for i := 0; i < 20; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			reqID := fmt.Sprintf("req-%d", id)
			resp, err := SendPing(serverAddr, reqID, 1*time.Second)
			if err != nil {
				t.Errorf("worker %d failed: %v", id, err)
				return
			}
			if resp != reqID {
				t.Errorf("worker %d expected %s, got %s", id, reqID, resp)
			}
		}(i)
	}

	wg.Wait()
}

func TestPingPong_Timeout(t *testing.T) {
	// Pick an unallocated address on 127.0.0.1 where nothing listens
	dummyAddr := "127.0.0.1:54321"

	_, err := SendPing(dummyAddr, "timeout-test", 50*time.Millisecond)
	if err == nil {
		t.Fatalf("expected timeout error, got nil")
	}
}
