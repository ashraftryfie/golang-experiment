# Stage 13: Concurrency & CSP Channels

Welcome to **Stage 13** of your Go engineering journey. In this stage, you master Go's signature concurrency model: Communicating Sequential Processes (CSP), goroutines, unbuffered vs buffered channels, channel ownership rules, `select` multiplexing, and bounded worker pools.

---

## 🧭 Mental Model

> *"Do not communicate by sharing memory; instead, share memory by communicating."* — Effective Go

### 1. Goroutines & The M:N Scheduler
- A goroutine is a lightweight user-space thread managed by the Go runtime, starting with only ~2KB of stack that grows and shrinks dynamically.
- The Go scheduler maps $M$ goroutines onto $N$ OS threads across $P$ logical processor cores using a cooperative work-stealing algorithm.

### 2. Channel Types & Rendezvous Semantics
- **Unbuffered Channel (`make(chan T)`)**:
  - Synchronous handoff. The sender blocks until a receiver is ready to receive, and vice-versa. Acts as a point of synchronization ("rendezvous").
- **Buffered Channel (`make(chan T, capacity)`)**:
  - Asynchronous queue. Sender blocks only when the buffer is full; receiver blocks only when the buffer is empty.

### 3. Channel States Matrix
Crucial for every Go engineer to memorize:

| Operation | `nil` Channel | Open Channel | Closed Channel |
| :--- | :--- | :--- | :--- |
| **Read (`<-ch`)** | Blocks forever | Returns value (or blocks if empty) | Returns zero value + `ok=false` immediately |
| **Write (`ch <- v`)** | Blocks forever | Sends value (or blocks if full) | **Panics!** (`send on closed channel`) |
| **Close (`close(ch)`)** | **Panics!** | Closes channel | **Panics!** (`close of closed channel`) |

### 4. Golden Rule of Channel Ownership
- The goroutine that **creates and writes** to a channel is its **owner**.
- Only the owner should close the channel.
- Receivers should NEVER close channels.
- If there are multiple concurrent writers, coordinate closing using a `sync.WaitGroup` or an explicit stop channel.

### 5. Multiplexing with `select`
```go
select {
case msg := <-ch1:
    fmt.Println("Received from ch1:", msg)
case ch2 <- val:
    fmt.Println("Sent to ch2")
case <-time.After(1 * time.Second):
    fmt.Println("Timed out")
default:
    fmt.Println("No channel was ready; non-blocking fallthrough")
}
```

---

## 🎯 Learning Objectives

By the end of this stage, you will:
- [x] Coordinate concurrent operations using `sync.WaitGroup`, `sync.Mutex`, and `sync.Once`.
- [x] Implement fan-out (parallel worker dispatch) and fan-in (channel multiplexing).
- [x] Build bounded worker pools to prevent unbounded resource exhaustion.
- [x] Eliminate goroutine leaks using cancellation signals and timeouts.
- [x] Verify thread safety using the Go race detector (`go test -race`).

---

## 🗂️ Stage Structure

```
13-concurrency/
├── README.md                           # Stage curriculum & CSP mental model
├── examples/
│   ├── worker_pool/main.go             # Canonical bounded worker pool
│   └── select_multiplex/main.go        # Multiplexing streams with timeouts
├── exercises/
│   ├── 01-fan-out-fan-in/              # Parallel task dispatch & result consolidation
│   ├── 02-bounded-worker-pool/         # Fixed-size concurrency worker queue
│   └── 03-rate-limited-batcher/        # Time- and size-windowed batch processor
├── solutions/
│   ├── 01-fan-out-fan-in/              # Reference implementation
│   ├── 02-bounded-worker-pool/         # Reference implementation
│   └── 03-rate-limited-batcher/        # Reference implementation
└── challenges/
    └── concurrent-crawler/             # Depth-limited concurrent crawler with visited set
```

---

## ⚠️ Common Pitfalls & Gotchas

1. **Goroutine Leak**: Spawning a goroutine that sends to an unbuffered channel that has no receiver. The goroutine stays alive forever, leaking memory and resources.
2. **`time.After` in a Loop**: Calling `time.After(duration)` inside a tight `select` loop allocates a new timer on every iteration that cannot be garbage collected until the duration expires. Use `time.NewTimer` or `time.NewTicker` and call `Reset()` or `Stop()`.
3. **Closing Closed Channels**: Attempting to close a channel twice causes an immediate panic.

---

## 🧪 Verification Commands

```powershell
# Run all tests in this stage
go test -v ./13-concurrency/...

# Check formatting and vet
go vet ./13-concurrency/...
go fmt ./13-concurrency/...
```

---

## 🔗 Connections
- **Prerequisites**: [Stage 12: Testing & Benchmarking](../12-testing/README.md).
- **Next Stage**: [Stage 14: Context & Cancellation](../14-context/README.md) (`context.Context`, timeouts, deadlines, request-scoped cancellation).
