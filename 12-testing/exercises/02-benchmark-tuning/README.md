# Exercise 02: Memory Allocation Profiling & Benchmark Tuning (Tier 2 - Medium)

## 🎯 Problem Statement
In high-throughput microservices, unconstrained memory allocations trigger frequent Garbage Collection (GC) pauses.

Implement `Deduplicate`:
```go
func Deduplicate(nums []int) []int
```
And benchmark it using `testing.B` with `b.ReportAllocs()`.

### Requirements
1. Return a slice containing only unique values in order of first appearance.
2. Return an empty slice `[]int{}` (not `nil`) if input is empty or nil.
3. In [`starter_test.go`](./starter_test.go), write `BenchmarkDeduplicate(b *testing.B)` using `b.ReportAllocs()` and `b.ResetTimer()`.

## 🛠️ Instructions
1. Open [`starter.go`](./starter.go) and implement `Deduplicate`.
2. Open [`starter_test.go`](./starter_test.go) and implement `BenchmarkDeduplicate`.
3. Run tests and benchmarks:
   ```powershell
   go test -v -bench=. -benchmem ./12-testing/exercises/02-benchmark-tuning
   ```
