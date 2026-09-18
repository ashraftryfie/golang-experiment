# Stage 12: Testing, Fuzzing & Benchmarking

Welcome to **Stage 12** of your Go engineering journey. In this stage, you master the Go standard testing ecosystem: table-driven unit tests, test helpers (`t.Helper()`), parallel execution (`t.Parallel()`), micro-benchmarks (`testing.B`) with allocation profiling, and native fuzzing (`testing.F`).

---

## 🧭 Mental Model

Go eschews heavy assertion libraries in favor of explicit Go code using the standard library's `testing` package.

### 1. Canonical Table-Driven Test Pattern
```go
func TestCalculateDiscount(t *testing.T) {
    tests := []struct {
        name        string
        total       float64
        isVIP       bool
        wantPercent float64
        wantErr     bool
    }{
        {name: "regular below threshold", total: 50, isVIP: false, wantPercent: 0},
        {name: "vip flat discount", total: 50, isVIP: true, wantPercent: 0.10},
        {name: "negative total returns error", total: -10, wantErr: true},
    }

    for _, tc := range tests {
        t.Run(tc.name, func(t *testing.T) {
            got, err := CalculateDiscount(tc.total, tc.isVIP)
            if (err != nil) != tc.wantErr {
                t.Fatalf("CalculateDiscount() error = %v, wantErr %v", err, tc.wantErr)
            }
            if got != tc.wantPercent {
                t.Errorf("CalculateDiscount() = %v, want %v", got, tc.wantPercent)
            }
        })
    }
}
```

### 2. Test Helpers with `t.Helper()`
When writing helper functions that make assertions or set up fixtures, call `t.Helper()` first. This strips the helper stack frame so test failures report the line number of the actual test caller:
```go
func assertEqual[T comparable](t *testing.T, got, want T) {
    t.Helper() // Reports caller's line on failure!
    if got != want {
        t.Errorf("got %v, want %v", got, want)
    }
}
```

### 3. Benchmarking & Allocation Profiling (`testing.B`)
```go
func BenchmarkBuilder(b *testing.B) {
    b.ReportAllocs() // Tracks allocations per op
    b.ResetTimer()   // Excludes setup time

    for i := 0; i < b.N; i++ {
        var sb strings.Builder
        sb.WriteString("hello")
        _ = sb.String()
    }
}
```
Run with:
```powershell
go test -bench=. -benchmem ./...
```

### 4. Native Fuzzing (Go 1.18+)
Fuzzing automatically generates thousands of pseudo-random inputs to find panics, out-of-bounds index exceptions, or infinite loops:
```go
func FuzzParseJSON(f *testing.F) {
    f.Add([]byte(`{"valid": true}`)) // Seed corpus
    f.Fuzz(func(t *testing.T, data []byte) {
        var result map[string]any
        _ = json.Unmarshal(data, &result) // Must never panic!
    })
}
```
Run with:
```powershell
go test -fuzz=FuzzParseJSON -fuzztime=10s ./...
```

---

## 🎯 Learning Objectives

By the end of this stage, you will:
- [x] Write idiomatic table-driven unit tests with descriptive subtests (`t.Run`).
- [x] Implement test helpers with `t.Helper()` and test teardown using `t.Cleanup()`.
- [x] Profile code performance and allocations using `b.ReportAllocs()` and `-benchmem`.
- [x] Uncover edge-case panics and memory violations using Go native fuzz testing (`testing.F`).
- [x] Detect race conditions in concurrent code using `go test -race`.

---

## 🗂️ Stage Structure

```
12-testing/
├── README.md                           # Stage curriculum & testing theory
├── examples/
│   ├── table_tests/service_test.go     # Table-driven testing with t.Helper()
│   └── benchmarks/hasher_test.go       # Benchmarking strings.Builder vs concat
├── exercises/
│   ├── 01-table-suite/                 # Table-driven test suite for slugifier
│   ├── 02-benchmark-tuning/            # Deduplication benchmark & memory optimization
│   └── 03-native-fuzzer/               # Fuzz test discovering edge-case panics
├── solutions/
│   ├── 01-table-suite/                 # Reference implementation
│   ├── 02-benchmark-tuning/            # Reference implementation
│   └── 03-native-fuzzer/               # Reference implementation
└── challenges/
    └── resilient-parser/               # Production parser with table tests, bench & fuzzing
```

---

## 🧪 Verification Commands

```powershell
# Run all tests in this stage
go test -v ./12-testing/...

# Run benchmarks with memory profiling
go test -bench=. -benchmem ./12-testing/...

# Vet and check formatting
go vet ./12-testing/...
go fmt ./12-testing/...
```

---

## 🔗 Connections
- **Prerequisites**: [Stage 11: PostgreSQL & Storage](../11-postgresql/README.md).
- **Next Stage**: [Stage 13: Concurrency & Channels](../13-concurrency/README.md) (Goroutines, channels, worker pools, `sync.Mutex`).
