package timeoutfetcher

import (
	"context"
	"sync"
	"time"
)

type URLFetcher func(ctx context.Context, url string) (string, error)

type FetchResult struct {
	URL  string
	Body string
	Err  error
}

// FetchURLs concurrently queries all given URLs using fetcher within the specified timeout.
// It returns a slice of FetchResult matching the order of input urls.
func FetchURLs(ctx context.Context, urls []string, fetcher URLFetcher, timeout time.Duration) ([]FetchResult, error) {
	ctxTimeout, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	results := make([]FetchResult, len(urls))
	var wg sync.WaitGroup

	for i, u := range urls {
		wg.Add(1)
		go func(idx int, targetURL string) {
			defer wg.Done()

			body, err := fetcher(ctxTimeout, targetURL)
			results[idx] = FetchResult{
				URL:  targetURL,
				Body: body,
				Err:  err,
			}
		}(i, u)
	}

	// Wait for all workers to finish.
	// Since workers respect ctxTimeout passed to fetcher, they are guaranteed to return.
	wg.Wait()

	return results, ctxTimeout.Err()
}
