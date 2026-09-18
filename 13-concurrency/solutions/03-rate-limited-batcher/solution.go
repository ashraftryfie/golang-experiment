package batcher_solution

import (
	"sync"
	"time"
)

type Batcher struct {
	maxSize  int
	flushFn  func([]string)
	incoming chan string
	flushNow chan struct{}
	done     chan struct{}
	wg       sync.WaitGroup
	once     sync.Once
}

// NewBatcher starts a background flusher that invokes flushFn on size or timeout.
func NewBatcher(maxSize int, maxWait time.Duration, flushFn func([]string)) *Batcher {
	b := &Batcher{
		maxSize:  maxSize,
		flushFn:  flushFn,
		incoming: make(chan string, maxSize*2),
		flushNow: make(chan struct{}, 1),
		done:     make(chan struct{}),
	}

	b.wg.Add(1)
	go b.run(maxWait)

	return b
}

func (b *Batcher) run(maxWait time.Duration) {
	defer b.wg.Done()

	ticker := time.NewTicker(maxWait)
	defer ticker.Stop()

	buffer := make([]string, 0, b.maxSize)

	flush := func() {
		if len(buffer) > 0 {
			b.flushFn(buffer)
			buffer = make([]string, 0, b.maxSize)
		}
	}

	for {
		select {
		case item, ok := <-b.incoming:
			if !ok {
				flush()
				return
			}
			buffer = append(buffer, item)
			if len(buffer) >= b.maxSize {
				flush()
			}

		case <-ticker.C:
			flush()

		case <-b.flushNow:
			flush()

		case <-b.done:
			// Drain incoming
			for item := range b.incoming {
				buffer = append(buffer, item)
			}
			flush()
			return
		}
	}
}

// Add appends an item to the current batch.
func (b *Batcher) Add(item string) {
	b.incoming <- item
}

// FlushAndClose flushes any pending items and cleanly terminates the batcher.
func (b *Batcher) FlushAndClose() {
	b.once.Do(func() {
		close(b.incoming)
		close(b.done)
		b.wg.Wait()
	})
}
