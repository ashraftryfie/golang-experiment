package timeoutfetcher

import (
	"context"
	"errors"
	"fmt"
	"testing"
	"time"
)

func TestFetchURLs_Starter(t *testing.T) {
	urls := []string{"https://api.example.com/1", "https://api.example.com/2"}

	mockFetcher := func(ctx context.Context, url string) (string, error) {
		select {
		case <-ctx.Done():
			return "", ctx.Err()
		case <-time.After(20 * time.Millisecond):
			return "content for " + url, nil
		}
	}

	results, err := FetchURLs(context.Background(), urls, mockFetcher, 100*time.Millisecond)
	if errors.Is(err, ErrNotImplemented) {
		t.Skip("Skipping unimplemented exercise: FetchURLs")
	}

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(results) != len(urls) {
		t.Fatalf("expected %d results, got %d", len(urls), len(results))
	}
	for i, r := range results {
		if r.URL != urls[i] {
			t.Errorf("expected URL %s at index %d, got %s", urls[i], i, r.URL)
		}
		expectedBody := fmt.Sprintf("content for %s", urls[i])
		if r.Body != expectedBody {
			t.Errorf("expected body %s, got %s", expectedBody, r.Body)
		}
	}
}
