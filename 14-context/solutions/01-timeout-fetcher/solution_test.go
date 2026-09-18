package timeoutfetcher

import (
	"context"
	"errors"
	"fmt"
	"testing"
	"time"
)

func TestFetchURLs_Success(t *testing.T) {
	urls := []string{"https://api.example.com/1", "https://api.example.com/2", "https://api.example.com/3"}

	mockFetcher := func(ctx context.Context, url string) (string, error) {
		select {
		case <-ctx.Done():
			return "", ctx.Err()
		case <-time.After(10 * time.Millisecond):
			return "content for " + url, nil
		}
	}

	results, err := FetchURLs(context.Background(), urls, mockFetcher, 200*time.Millisecond)
	if err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}

	if len(results) != len(urls) {
		t.Fatalf("expected %d results, got %d", len(urls), len(results))
	}

	for i, r := range results {
		if r.URL != urls[i] {
			t.Errorf("expected URL %s, got %s", urls[i], r.URL)
		}
		expectedBody := fmt.Sprintf("content for %s", urls[i])
		if r.Body != expectedBody {
			t.Errorf("expected body %s, got %s", expectedBody, r.Body)
		}
		if r.Err != nil {
			t.Errorf("expected nil err, got %v", r.Err)
		}
	}
}

func TestFetchURLs_Timeout(t *testing.T) {
	urls := []string{"https://fast.example.com", "https://slow.example.com"}

	mockFetcher := func(ctx context.Context, url string) (string, error) {
		if url == "https://fast.example.com" {
			return "fast", nil
		}
		select {
		case <-ctx.Done():
			return "", ctx.Err()
		case <-time.After(200 * time.Millisecond):
			return "slow", nil
		}
	}

	results, err := FetchURLs(context.Background(), urls, mockFetcher, 50*time.Millisecond)
	if !errors.Is(err, context.DeadlineExceeded) {
		t.Fatalf("expected context.DeadlineExceeded, got %v", err)
	}

	if len(results) != 2 {
		t.Fatalf("expected 2 results, got %d", len(results))
	}

	if results[0].Body != "fast" || results[0].Err != nil {
		t.Errorf("fast request should have succeeded: %+v", results[0])
	}
	if !errors.Is(results[1].Err, context.DeadlineExceeded) {
		t.Errorf("slow request should have failed with DeadlineExceeded: %+v", results[1])
	}
}
