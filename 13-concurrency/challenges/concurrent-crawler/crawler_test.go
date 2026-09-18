package crawler

import (
	"fmt"
	"sort"
	"sync/atomic"
	"testing"
	"time"
)

type mockFetcher struct {
	graph       map[string][]string
	activeCount int64
	maxActive   int64
}

func (m *mockFetcher) Fetch(url string) (string, []string, error) {
	current := atomic.AddInt64(&m.activeCount, 1)

	// Track peak concurrency
	for {
		max := atomic.LoadInt64(&m.maxActive)
		if current <= max || atomic.CompareAndSwapInt64(&m.maxActive, max, current) {
			break
		}
	}

	time.Sleep(10 * time.Millisecond) // Simulate network I/O
	atomic.AddInt64(&m.activeCount, -1)

	links, ok := m.graph[url]
	if !ok {
		return "", nil, fmt.Errorf("404: %s", url)
	}
	return "mock body", links, nil
}

func TestCrawlBoundedConcurrency(t *testing.T) {
	fetcher := &mockFetcher{
		graph: map[string][]string{
			"https://golang.org": {
				"https://golang.org/pkg",
				"https://golang.org/cmd",
			},
			"https://golang.org/pkg": {
				"https://golang.org/pkg/fmt",
				"https://golang.org/pkg/sync",
			},
			"https://golang.org/cmd": {
				"https://golang.org", // cycle test
			},
			"https://golang.org/pkg/fmt":  {},
			"https://golang.org/pkg/sync": {},
		},
	}

	maxConcurrency := 2
	crawled := Crawl("https://golang.org", 2, maxConcurrency, fetcher)

	if len(crawled) == 0 {
		t.Fatalf("expected non-empty crawled URLs")
	}

	// Verify peak concurrency never exceeded the limit
	if peak := atomic.LoadInt64(&fetcher.maxActive); peak > int64(maxConcurrency) {
		t.Errorf("peak concurrency %d exceeded max allowed %d", peak, maxConcurrency)
	}

	// Verify cycle was visited at most once
	counts := make(map[string]int)
	for _, u := range crawled {
		counts[u]++
		if counts[u] > 1 {
			t.Errorf("URL %s was crawled %d times (cycle bug)", u, counts[u])
		}
	}

	sort.Strings(crawled)
	expected := []string{
		"https://golang.org",
		"https://golang.org/cmd",
		"https://golang.org/pkg",
		"https://golang.org/pkg/fmt",
		"https://golang.org/pkg/sync",
	}

	if len(crawled) != len(expected) {
		t.Errorf("got %d crawled URLs, want %d: %+v", len(crawled), len(expected), crawled)
	}
}
