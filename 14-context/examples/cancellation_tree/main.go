package main

import (
	"context"
	"fmt"
	"sync"
	"time"
)

// worker simulates an asynchronous background worker that respects context cancellation.
func worker(ctx context.Context, id int, wg *sync.WaitGroup) {
	defer wg.Done()
	fmt.Printf("[Worker %d] Started\n", id)

	for {
		select {
		case <-ctx.Done():
			fmt.Printf("[Worker %d] Stopped due to: %v\n", id, ctx.Err())
			return
		case <-time.After(100 * time.Millisecond):
			fmt.Printf("[Worker %d] Completed task slice\n", id)
		}
	}
}

func main() {
	// Root context
	rootCtx := context.Background()

	// Parent context with manual cancellation
	parentCtx, parentCancel := context.WithCancel(rootCtx)

	// Child context derived from parent with 350ms timeout
	childTimeoutCtx, childCancel := context.WithTimeout(parentCtx, 350*time.Millisecond)
	defer childCancel()

	var wg sync.WaitGroup

	// Launch worker attached to child timeout context
	wg.Add(1)
	go worker(childTimeoutCtx, 1, &wg)

	// Launch worker attached directly to parent context
	wg.Add(1)
	go worker(parentCtx, 2, &wg)

	// Wait 250ms and cancel parent manually before child timeout hits
	time.Sleep(250 * time.Millisecond)
	fmt.Println("--> Canceling parent context explicitly...")
	parentCancel()

	wg.Wait()
	fmt.Println("All workers exited cleanly. No goroutines leaked.")
}
