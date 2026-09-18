# Stage 07: Standard Library Deep Dive

## 🎯 Goal
Master the standard library's streaming primitives: `io.Reader` and `io.Writer`, buffered I/O with `bufio.Scanner` and `bufio.Writer`, composable stream pipelines (`io.TeeReader`, `io.MultiReader`, `io.LimitReader`), zero-allocation string formatting with `strings.Builder`, and the `time` package layout mechanics.

---

## 🧠 Core Concepts

### 1. The Heartbeat of Go: `io.Reader` and `io.Writer`
Everything in Go that streams data (files, TCP sockets, HTTP request/response bodies, in-memory buffers, gzip decoders) implements these two interfaces:
```go
type Reader interface {
    Read(p []byte) (n int, err error)
}

type Writer interface {
    Write(p []byte) (n int, err error)
}
```

#### Critical Rules of `io.Reader`:
1. `Read` fills the provided buffer slice `p` up to `len(p)` bytes.
2. It returns `n`, the number of bytes read ($0 \le n \le len(p)$).
3. **The `io.EOF` Rule**: An `io.Reader` may return a non-zero `n` along with `err == io.EOF` at the end of a stream! **Always process the `n` bytes read before checking for `err == io.EOF`**.

### 2. Composable Stream Combinators
- **`io.Copy(dst, src)`**: Continuously streams from `src` to `dst` in fixed 32 KB chunks with $O(1)$ memory usage.
- **`io.TeeReader(r, w)`**: Splits a stream like an audio tee: reading from the TeeReader reads from `r` and simultaneously writes a copy to `w`.
- **`io.MultiReader(r1, r2, ...)`**: Concatenates multiple readers sequentially into a single unified stream.
- **`io.LimitReader(r, n)`**: Caps reading at `n` bytes, protecting servers from malicious payload exhaustion.

### 3. Buffered I/O: `bufio`
Making a system call for every single byte read/written is extremely slow. `bufio` provides an in-memory buffer:
- `bufio.Scanner`: Line-by-line or token-by-token parsing.
- `bufio.Writer`: Batches writes in memory. **Must call `Flush()`** when finished to flush remaining buffered bytes to the underlying destination.

### 4. Time Formatting & The Reference Date
Go does not use cryptic formatting strings like `%Y-%m-%d`. Instead, Go uses a **mnemonic reference date**:
```text
Mon Jan 2 15:04:05 MST 2006 (MST is UTC-0700)
 0   1  2  3  4  5   6    7
```
- Month: `01`
- Day: `02`
- Hour: `15` (or `03` for 12-hour)
- Minute: `04`
- Second: `05`
- Year: `2006`

Example: `time.Now().Format("2006-01-02 15:04:05")`.

---

## 🔍 Mental Model

Think of `io.Reader` and `io.Writer` as standard **plumbing pipes**:

```text
┌─────────────────┐       ┌─────────────────┐       ┌─────────────────┐
│   os.File /     │       │   bufio.Writer  │       │  http.Response  │
│   HTTP Body     ├──────►│   (32 KB buffer)├──────►│   TCP Socket    │
│  (io.Reader)    │       │   (io.Writer)   │       │  (io.Writer)    │
└─────────────────┘       └─────────────────┘       └─────────────────┘
```
Because every component speaks the exact same byte interface, you can snap filters (compression, hashing, encryption) into the pipe without changing existing code.

---

## 💻 Examples
See executable code in:
- [`examples/streaming_io/main.go`](./examples/streaming_io/main.go)
- [`examples/time_and_formatting/main.go`](./examples/time_and_formatting/main.go)

---

## 🧪 Exercises

### 1. Custom ROT13 Cipher Stream (Tier 1 - Easy)
- **Path**: [`exercises/01-custom-reader/`](./exercises/01-custom-reader/)
- **Concepts**: Implementing the `io.Reader` interface, byte buffer transformation in-place, handling EOF.
- **Run tests**:
  ```powershell
  go test -v ./07-standard-library/exercises/01-custom-reader
  ```

### 2. High-Performance Log Line Scanner (Tier 2 - Medium)
- **Path**: [`exercises/02-line-scanner/`](./exercises/02-line-scanner/)
- **Concepts**: `bufio.Scanner`, custom split functions, memory-safe line reading without unbounded allocations.
- **Run tests**:
  ```powershell
  go test -v ./07-standard-library/exercises/02-line-scanner
  ```

### 3. Simultaneous Stream Hashing with TeeReader (Tier 3 - Hard)
- **Path**: [`exercises/03-tee-hasher/`](./exercises/03-tee-hasher/)
- **Concepts**: `io.TeeReader`, `crypto/sha256`, simultaneous payload copying and checksum verification in a single pass.
- **Run tests**:
  ```powershell
  go test -v ./07-standard-library/exercises/03-tee-hasher
  ```

---

## 🧩 Challenge

Build a **Streaming Gzip Archive Pipeline** in [`challenges/file-archive-streamer/`](./challenges/file-archive-streamer/):
- Compresses an incoming `io.Reader` stream with `compress/gzip` directly to an `io.Writer` destination.
- Measures total uncompressed bytes read, total compressed bytes written, and calculates the compression ratio.
- Never loads the whole payload into memory.

---

## ⚠️ Common Mistakes

1. **Ignoring bytes read when `err == io.EOF`**:
   ```go
   // ANTI-PATTERN: Misses the final chunk of data!
   n, err := r.Read(buf)
   if err != nil { return } // If err == io.EOF, n bytes were read and discarded!
   
   // IDIOMATIC:
   n, err := r.Read(buf)
   if n > 0 { process(buf[:n]) }
   if err != nil { /* check io.EOF */ }
   ```
2. **Forgetting to `Flush()` a `bufio.Writer`**: Data stays trapped in memory and is never written to disk or network!
3. **Using `+` for string concatenation in loops**: Creates huge allocation overheads. Use `strings.Builder` instead.

---

## 🔄 Review

1. What is the contract of `io.Reader.Read([]byte)` when reaching the end of input?
2. Why is `io.Copy(dst, src)` more memory-efficient than reading with `os.ReadFile` and writing with `os.WriteFile`?
3. How does Go's reference time `2006-01-02 15:04:05` make layout strings intuitive?
4. What happens if you do not call `Flush()` on a `bufio.Writer` before closing?

---

## ✅ Completion Criteria

- [ ] All 3 exercises pass unit tests (`go test -v ./...`).
- [ ] Challenge in `challenges/file-archive-streamer/` passes tests.
- [ ] Code passes `go vet` and `go fmt`.
- [ ] Progress updated in `learning/progress.yaml`.

---

## 🔗 Connections
- **Prerequisites**: [Stage 06: Packages & Modules](../06-packages-modules/README.md).
- **Next Stage**: [Stage 08: Files, JSON & Serialization](../08-files-json/README.md) (Streaming encoders vs Unmarshal, struct tags).
