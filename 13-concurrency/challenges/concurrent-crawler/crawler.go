package crawler

import (
	"sync"
)

type Fetcher interface {
	Fetch(url string) (body string, urls []string, err error)
}

type crawlerState struct {
	mu      sync.Mutex
	visited map[string]bool
	crawled []string
	sem     chan struct{}
}

// Crawl recursively visits links up to maxDepth with a bounded concurrency pool.
func Crawl(startURL string, maxDepth, maxConcurrency int, fetcher Fetcher) []string {
	if maxConcurrency <= 0 {
		maxConcurrency = 1
	}

	state := &crawlerState{
		visited: make(map[string]bool),
		crawled: make([]string, 0),
		sem:     make(chan struct{}, maxConcurrency),
	}

	var wg sync.WaitGroup
	wg.Add(1)

	state.visited[startURL] = true
	go state.crawlURL(startURL, 0, maxDepth, fetcher, &wg)

	wg.Wait()
	return state.crawled
}

func (s *crawlerState) crawlURL(url string, depth, maxDepth int, fetcher Fetcher, wg *sync.WaitGroup) {
	defer wg.Done()

	// Acquire semaphore token
	s.sem <- struct{}{}
	_, childURLs, err := fetcher.Fetch(url)
	<-s.sem // Release token

	s.mu.Lock()
	s.crawled = append(s.crawled, url)
	s.mu.Unlock()

	if err != nil || depth >= maxDepth {
		return
	}

	for _, next := range childURLs {
		s.mu.Lock()
		if !s.visited[next] {
			s.visited[next] = true
			s.mu.Unlock()

			wg.Add(1)
			go s.crawlURL(next, depth+1, maxDepth, fetcher, wg)
		} else {
			s.mu.Unlock()
		}
	}
}
