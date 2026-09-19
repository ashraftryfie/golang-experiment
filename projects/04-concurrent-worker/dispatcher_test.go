package eventdispatcher

import (
	"context"
	"fmt"
	"sync"
	"testing"
	"time"
)

func TestDispatcher_FanOutFanIn(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	handler := func(_ context.Context, e Event) error {
		time.Sleep(1 * time.Millisecond)
		return nil
	}

	d := NewDispatcher(4, 50, handler)
	d.Start(ctx)

	// Collect results in background
	var processedCount int
	var countMu sync.Mutex
	done := make(chan struct{})

	go func() {
		for res := range d.Results() {
			if res.Processed {
				countMu.Lock()
				processedCount++
				countMu.Unlock()
			}
		}
		close(done)
	}()

	// Dispatch 20 events
	for i := 0; i < 20; i++ {
		err := d.Dispatch(Event{ID: fmt.Sprintf("evt-%d", i), Topic: "orders"})
		if err != nil {
			t.Fatalf("dispatch error: %v", err)
		}
	}

	d.Stop()
	<-done

	countMu.Lock()
	count := processedCount
	countMu.Unlock()

	if count != 20 {
		t.Fatalf("expected 20 processed events, got %d", count)
	}

	in, handled := d.Stats()
	if in != 20 || handled != 20 {
		t.Fatalf("expected in=20, handled=20, got in=%d, handled=%d", in, handled)
	}
}
