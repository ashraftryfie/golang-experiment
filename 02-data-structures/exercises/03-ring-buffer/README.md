# Exercise 03: Circular Ring Buffer (Tier 3 - Hard)

## 🎯 Problem Statement
Implement a bounded, fixed-capacity **Circular Ring Buffer** for integer data using a single contiguous slice.

A ring buffer is a high-performance FIFO data structure commonly used in network packet buffering, audio streaming, and lock-free queues. It uses a fixed memory allocation and wraps indices around using modulo arithmetic.

### Requirements
Implement the `RingBuffer` struct and its methods:
1. `NewRingBuffer(capacity int) (*RingBuffer, error)`:
   - Validates that `capacity > 0`. Returns `ErrInvalidCapacity` otherwise.
   - Pre-allocates an internal slice of fixed capacity.
2. `Push(val int) bool`:
   - Adds `val` to the tail of the buffer.
   - If the buffer is full, returns `false` (does not overwrite).
   - If successful, advances tail and returns `true`.
3. `Pop() (int, bool)`:
   - Removes and returns the value at the head of the buffer.
   - If the buffer is empty, returns `0, false`.
   - If successful, advances head and returns `val, true`.
4. `Peek() (int, bool)`:
   - Returns the value at the head without removing it.
   - If empty, returns `0, false`.
5. `Size() int`:
   - Returns the current number of items in the buffer.
6. `Capacity() int`:
   - Returns total buffer capacity.
7. `IsFull() bool` and `IsEmpty() bool`.

---

## 🛠️ Instructions
1. Open [`starter.go`](./starter.go).
2. Implement methods using modulo math: `(index + 1) % capacity`.
3. Run tests:
   ```powershell
   go test -v ./02-data-structures/exercises/03-ring-buffer
   ```
