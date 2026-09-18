# Exercise 02: Streaming Log Scanner (Tier 2 - Medium)

## 🎯 Problem Statement
Implement a memory-safe log scanner that processes an `io.Reader` line-by-line using `bufio.Scanner` to extract and collect lines matching a specific log level:

```go
type LogFilter struct {
    Level string // e.g. "[ERROR]", "[WARN]"
}

func FilterLogs(r io.Reader, targetLevel string) ([]string, error)
```

### Requirements
1. Use `bufio.NewScanner(r)` to read stream line-by-line without loading entire payload into memory.
2. If a line contains `targetLevel`, append it to the result slice.
3. Check `scanner.Err()` after the scanning loop finishes to catch any I/O errors.
4. If `targetLevel` is empty, return an error (`ErrEmptyTargetLevel`).

---

## 🛠️ Instructions
1. Open [`starter.go`](./starter.go).
2. Implement `FilterLogs`.
3. Run tests:
   ```powershell
   go test -v ./07-standard-library/exercises/02-line-scanner
   ```
