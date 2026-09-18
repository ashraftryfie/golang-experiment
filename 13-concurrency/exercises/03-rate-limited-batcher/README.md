# Exercise 03: Windowed & Size-Limited Concurrency Batcher (Tier 3 - Hard)

## 🎯 Problem Statement
In log forwarding, database inserts, and event pipelines, batching events reduces network roundtrips. A production batcher must flush under two conditions:
1. **Size trigger**: As soon as the buffer reaches `maxSize`.
2. **Time trigger**: If `maxWait` passes with items still sitting in the buffer.

Implement `Batcher`:
```go
type Batcher struct { ... }

func NewBatcher(maxSize int, maxWait time.Duration, flushFn func([]string)) *Batcher
func (b *Batcher) Add(item string)
func (b *Batcher) FlushAndClose()
```

### Requirements
1. `Add(item)` adds an item to the active batch.
2. If the batch reaches `maxSize`, call `flushFn(batch)` immediately and reset the buffer.
3. If `maxWait` elapses and the batch has at least one item, call `flushFn(batch)` and reset the buffer.
4. `FlushAndClose()` flushes any remaining items, terminates the background timer loop, and ensures no goroutines leak.
5. All operations must be thread-safe.

## 🛠️ Instructions
1. Open [`starter.go`](./starter.go).
2. Implement `Batcher`.
3. Run tests:
   ```powershell
   go test -v ./13-concurrency/exercises/03-rate-limited-batcher
   ```
