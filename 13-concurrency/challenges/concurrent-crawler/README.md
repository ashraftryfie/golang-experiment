# Challenge: Concurrent Bounded Web Crawler

## 🎯 Objective
Build a concurrent web crawler that explores links up to `maxDepth` while limiting concurrent network fetches to `maxConcurrency`, tracks visited URLs to prevent loops, and avoids data races.

## 📋 Requirements
1. **Fetcher Interface**:
   ```go
   type Fetcher interface {
       Fetch(url string) (body string, urls []string, err error)
   }
   ```
2. **Crawl Function**:
   ```go
   func Crawl(startURL string, maxDepth int, maxConcurrency int, fetcher Fetcher) []string
   ```
3. **Concurrency Constraints**:
   - Limit simultaneous `Fetch` calls to `maxConcurrency` using a buffered channel semaphore or bounded worker pool.
   - Use `sync.Mutex` or `sync.Map` to track visited URLs and prevent duplicate crawling.
   - Stop recursing once `depth >= maxDepth`.
   - Wait for all active crawler goroutines to complete using `sync.WaitGroup`.
   - Must pass `go test -race` with 0 data races.

## 🧪 Verification
```powershell
go test -v ./13-concurrency/challenges/concurrent-crawler/...
```
