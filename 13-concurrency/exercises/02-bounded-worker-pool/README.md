# Exercise 02: Bounded Worker Pool with Graceful Shutdown (Tier 2 - Medium)

## 🎯 Problem Statement
Spawning an unbounded number of goroutines under heavy client traffic risks crashing the application through memory exhaustion and CPU thrashing.

Implement a bounded worker pool:
```go
type WorkerPool struct { ... }

func NewWorkerPool(workerCount int, queueCapacity int) *WorkerPool
func (p *WorkerPool) Submit(task func()) bool
func (p *WorkerPool) Shutdown()
```

### Requirements
1. `NewWorkerPool`: Spawns `workerCount` persistent worker goroutines listening on an internal task queue channel.
2. `Submit(task)`:
   - If the pool is running, queues `task` for execution and returns `true`.
   - If the pool is shutting down or already stopped, returns `false` without executing or queueing.
3. `Shutdown()`:
   - Closes the task queue so workers exit once the queue is drained.
   - Blocks until all queued tasks finish executing using `sync.WaitGroup`.
   - Safe to call multiple times (`sync.Once`).

## 🛠️ Instructions
1. Open [`starter.go`](./starter.go).
2. Implement `WorkerPool`.
3. Run tests:
   ```powershell
   go test -v ./13-concurrency/exercises/02-bounded-worker-pool
   ```
