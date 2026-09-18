package workerpool

import (
	"sync/atomic"
	"testing"
	"time"
)

func TestWorkerPool(t *testing.T) {
	pool := NewWorkerPool(4, 20)
	if pool == nil {
		t.Skip("skipping: NewWorkerPool is not yet implemented (implement in starter.go)")
	}

	var counter int64
	totalTasks := 50

	for i := 0; i < totalTasks; i++ {
		ok := pool.Submit(func() {
			time.Sleep(2 * time.Millisecond)
			atomic.AddInt64(&counter, 1)
		})
		if !ok {
			t.Fatalf("failed to submit task %d", i)
		}
	}

	pool.Shutdown()

	if atomic.LoadInt64(&counter) != int64(totalTasks) {
		t.Errorf("expected %d completed tasks, got %d", totalTasks, counter)
	}

	// Submitting after shutdown must return false
	if ok := pool.Submit(func() {}); ok {
		t.Errorf("expected Submit() to return false after Shutdown()")
	}
}
