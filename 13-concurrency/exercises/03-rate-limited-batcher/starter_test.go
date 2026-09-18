package batcher

import (
	"sync"
	"testing"
	"time"
)

func TestBatcherSizeFlush(t *testing.T) {
	var mu sync.Mutex
	var flushed [][]string

	flushFn := func(batch []string) {
		mu.Lock()
		defer mu.Unlock()
		copyBatch := make([]string, len(batch))
		copy(copyBatch, batch)
		flushed = append(flushed, copyBatch)
	}

	b := NewBatcher(3, 500*time.Millisecond, flushFn)
	if b == nil {
		t.Skip("skipping: NewBatcher is not yet implemented (implement in starter.go)")
	}

	b.Add("e1")
	b.Add("e2")
	b.Add("e3") // Should flush immediately on reaching size 3

	time.Sleep(50 * time.Millisecond)

	mu.Lock()
	if len(flushed) != 1 || len(flushed[0]) != 3 {
		mu.Unlock()
		t.Fatalf("expected 1 batch of size 3, got %+v", flushed)
	}
	mu.Unlock()

	b.Add("e4")
	b.FlushAndClose() // Should flush remainder "e4"

	mu.Lock()
	defer mu.Unlock()
	if len(flushed) != 2 || len(flushed[1]) != 1 {
		t.Fatalf("expected 2nd batch of size 1, got %+v", flushed)
	}
}

func TestBatcherTimeoutFlush(t *testing.T) {
	var mu sync.Mutex
	var flushed [][]string

	flushFn := func(batch []string) {
		mu.Lock()
		defer mu.Unlock()
		copyBatch := make([]string, len(batch))
		copy(copyBatch, batch)
		flushed = append(flushed, copyBatch)
	}

	b := NewBatcher(10, 50*time.Millisecond, flushFn)
	if b == nil {
		t.Skip("skipping: NewBatcher is not yet implemented")
	}

	b.Add("slow_item")
	time.Sleep(120 * time.Millisecond) // Wait for timer tick

	mu.Lock()
	if len(flushed) != 1 || len(flushed[0]) != 1 {
		mu.Unlock()
		t.Fatalf("expected timer to flush 1 item, got %+v", flushed)
	}
	mu.Unlock()

	b.FlushAndClose()
}
