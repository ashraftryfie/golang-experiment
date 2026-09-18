# Exercise 03: Simultaneous Stream Hashing with `io.TeeReader` (Tier 3 - Hard)

## 🎯 Problem Statement
Implement a streaming copier that copies data from an `io.Reader` to an `io.Writer` and simultaneously computes its SHA-256 checksum in a single pass without allocating a buffer for the whole stream:

```go
type StreamResult struct {
    BytesWritten int64
    SHA256Hex    string
}

func CopyAndHash(dst io.Writer, src io.Reader) (StreamResult, error)
```

### Requirements
1. Create a SHA-256 hasher using `sha256.New()`.
2. Use `io.TeeReader(src, hasher)` to compute the hash incrementally as bytes are read.
3. Stream from the TeeReader directly to `dst` using `io.Copy(dst, tee)`.
4. Return total `BytesWritten` and the formatted hex checksum.

---

## 🛠️ Instructions
1. Open [`starter.go`](./starter.go).
2. Implement `CopyAndHash`.
3. Run tests:
   ```powershell
   go test -v ./07-standard-library/exercises/03-tee-hasher
   ```
