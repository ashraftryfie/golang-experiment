# Exercise 01: Fan-Out / Fan-In Parallel Pipeline (Tier 1 - Easy)

## 🎯 Problem Statement
When processing large slices of independent work items (e.g. image thumbnails, checksums), fan-out splits tasks across $N$ concurrent workers, while fan-in aggregates their results back into a single slice.

Implement `FanOutFanIn`:
```go
func FanOutFanIn(inputs []int, workerCount int, transform func(int) int) []int
```

### Requirements
1. If `len(inputs) == 0`, return an empty slice `[]int{}` immediately.
2. If `workerCount <= 0`, default `workerCount = 1`.
3. Create a `jobs` channel and distribute all inputs to workers.
4. Spawn `workerCount` goroutines that apply `transform(val)`.
5. Use `sync.WaitGroup` to wait for all workers to finish, then close the results channel.
6. Collect and return all results in a single slice.

## 🛠️ Instructions
1. Open [`starter.go`](./starter.go).
2. Implement `FanOutFanIn`.
3. Run tests:
   ```powershell
   go test -v ./13-concurrency/exercises/01-fan-out-fan-in
   ```
