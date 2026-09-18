package batcher

import (
	"time"
)

type Batcher struct {
	// TODO: Add fields for buffering, ticker, channels, and synchronization
}

// NewBatcher starts a background flusher that invokes flushFn on size or timeout.
func NewBatcher(maxSize int, maxWait time.Duration, flushFn func([]string)) *Batcher {
	// TODO: Initialize struct and background select loop
	return nil
}

// Add appends an item to the current batch.
func (b *Batcher) Add(item string) {
	// TODO: Send item or append under lock, flush if threshold reached
}

// FlushAndClose flushes any pending items and cleanly terminates the batcher.
func (b *Batcher) FlushAndClose() {
	// TODO: Flush remainder, stop ticker, wait for background goroutine
}
