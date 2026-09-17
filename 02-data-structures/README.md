# Stage 02: Data Structures (Arrays, Slices, Maps)

## 🎯 Goal
Deeply understand Go's memory layout for collections: why arrays are fixed-size values, how slice headers peer into backing arrays, how `append()` manages capacity and re-allocation, how to prevent slice memory leaks, and how hash maps work under the hood.

---

## 🧠 Core Concepts

### 1. Arrays vs. Slices
- **Array (`[N]T`)**: Fixed length, size is part of the type (`[3]int` != `[4]int`). Arrays are values in Go: passing an array to a function creates a complete copy of all its elements.
- **Slice (`[]T`)**: Dynamically sized window into an underlying backing array. Slices are cheap to pass by value because a slice is just a 24-byte header.

### 2. The Slice Header Internals
A slice in the Go runtime (`reflect.SliceHeader`) is a struct with exactly 3 machine words (24 bytes on 64-bit systems):
```go
type SliceHeader struct {
    Data uintptr // pointer to first element in backing array
    Len  int     // number of elements currently present
    Cap  int     // total capacity before re-allocation is required
}
```

```text
Slice Header (24 bytes)
┌──────────────┬──────────────┬──────────────┐
│  Data (ptr)  │   Len = 3    │   Cap = 5    │
└──────┬───────┴──────────────┴──────────────┘
       │
       ▼
Backing Array in Memory
┌──────┬──────┬──────┬──────┬──────┐
│  10  │  20  │  30  │  0   │  0   │
└──────┴──────┴──────┴──────┴──────┘
[─────── Len ────────]
[───────────── Cap ────────────────]
```

### 3. Append & Backing Array Re-allocation
- When `len == cap`, calling `append()` forces Go's runtime to allocate a brand-new, larger backing array (typically ~2x growth for small slices, ~1.25x for larger slices) and copy existing items.
- If multiple slices share the same backing array and `len < cap`, mutating an element via one slice mutates the backing array for all!

### 4. Maps
- Hash table mapping `comparable` keys to values.
- Reading from a `nil` map safely yields the value type's zero value.
- Writing to a `nil` map causes an immediate **panic**. Always initialize maps with `make(map[K]V)` or a map literal `{}`.
- Iteration order over maps is randomized by design in Go to prevent developers from relying on hash traversal order.

---

## 🔍 Mental Model

Think of a slice as a magnifying glass looking at a physical conveyor belt (the backing array):
- The magnifying glass knows where it starts (`Data`), how many items it shows right now (`Len`), and how long the conveyor belt extends before falling off the edge (`Cap`).
- If you copy the magnifying glass, both people look at the same conveyor belt! Changing an item changes what both people see.
- If you add an item that doesn't fit on the conveyor belt, Go builds an entirely new conveyor belt in a new warehouse, copies everything over, and points your magnifying glass there. The other person is still looking at the old conveyor belt!

---

## 💻 Examples

See executable code in:
- [`examples/slices_internals/main.go`](./examples/slices_internals/main.go)
- [`examples/maps_usage/main.go`](./examples/maps_usage/main.go)

---

## 🧪 Exercises

### 1. Slice Manipulator & Deduplication (Tier 1 - Easy)
- **Path**: [`exercises/01-slice-ops/`](./exercises/01-slice-ops/)
- **Concepts**: In-place filtering without extra allocations, deduplication, element deletion preserving order.
- **Run tests**:
  ```powershell
  go test -v ./02-data-structures/exercises/01-slice-ops
  ```

### 2. Frequency Counter & Top-K (Tier 2 - Medium)
- **Path**: [`exercises/02-word-frequency/`](./exercises/02-word-frequency/)
- **Concepts**: Maps for counting, comma-ok idiom, sorting keys with `sort.Slice`, deterministic output.
- **Run tests**:
  ```powershell
  go test -v ./02-data-structures/exercises/02-word-frequency
  ```

### 3. Circular Ring Buffer (Tier 3 - Hard)
- **Path**: [`exercises/03-ring-buffer/`](./exercises/03-ring-buffer/)
- **Concepts**: Bounded circular buffer over a fixed-capacity slice, wrapping head/tail pointers, FIFO eviction.
- **Run tests**:
  ```powershell
  go test -v ./02-data-structures/exercises/03-ring-buffer
  ```

---

## 🧩 Challenge

Build an in-memory **Log Metrics Aggregator** in [`challenges/metrics-aggregator/`](./challenges/metrics-aggregator/):
- Processes log event streams (`timestamp`, `service`, `status_code`, `latency_ms`).
- Aggregates request counts, error rates (`5xx`), and average latencies per service using maps and slices.
- Emits formatted statistical summaries.

---

## ⚠️ Common Mistakes

1. **Writing to a `nil` Map**:
   ```go
   var m map[string]int
   m["key"] = 1 // PANIC: assignment to entry in nil map
   ```
2. **Sub-slice Memory Leak**:
   ```go
   // Slicing 5 bytes from a 100 MB file keeps the entire 100 MB array in memory!
   func getPrefix() []byte {
       large := read100MBFile()
       return large[:5] // Backing array cannot be garbage collected!
   }
   // FIX: copy to a fresh slice
   func getPrefix() []byte {
       large := read100MBFile()
       prefix := make([]byte, 5)
       copy(prefix, large[:5])
       return prefix
   }
   ```
3. **Assuming Map Iteration is Ordered**: Map keys in Go are randomized on each `for range`. If you need ordered iteration, collect the keys into a slice, sort the slice, and iterate the sorted keys.

---

## 🔄 Review

1. What are the three fields inside a 64-bit Go slice header, and what is its total memory size?
2. Why does `[5]int` have a different type than `[6]int` in Go?
3. What happens under the hood when `append()` exceeds the slice's capacity?
4. How do you distinguish between a missing map key vs. a key whose value is the zero value?
5. How does the 3-index slice expression `s[low:high:max]` protect against unintended mutations to the original slice's capacity?

---

## ✅ Completion Criteria

- [ ] All 3 exercises in `exercises/` pass unit tests (`go test -v ./...`).
- [ ] Challenge in `challenges/metrics-aggregator/` passes tests.
- [ ] Understand slice headers, capacity growth, and nil map behavior.
- [ ] Progress updated in `learning/progress.yaml`.

---

## 🔗 Connections
- **Prerequisites**: [Stage 01: Go Fundamentals](../01-fundamentals/README.md).
- **Next Stage**: [Stage 03: Functions & Methods](../03-functions-methods/README.md) (Closures, defer, receiver methods).
