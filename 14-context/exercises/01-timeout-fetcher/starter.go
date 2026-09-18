package timeoutfetcher

import (
	"context"
	"errors"
	"time"
)

var ErrNotImplemented = errors.New("exercise not implemented yet")

type URLFetcher func(ctx context.Context, url string) (string, error)

type FetchResult struct {
	URL  string
	Body string
	Err  error
}

// FetchURLs concurrently queries all given URLs using fetcher within the specified timeout.
// It returns a slice of FetchResult matching the order of input urls.
func FetchURLs(ctx context.Context, urls []string, fetcher URLFetcher, timeout time.Duration) ([]FetchResult, error) {
	return nil, ErrNotImplemented
}
