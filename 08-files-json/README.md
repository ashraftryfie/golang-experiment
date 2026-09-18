# Stage 08: Files, JSON & Serialization

Welcome to **Stage 08** of your Go engineering journey. In this stage, you transition from basic streaming abstractions to persistent file operations, robust JSON serialization, custom encoding interfaces, and safe filesystem operations.

---

## 🧭 Mental Model

Serialization and persistence in Go rest on two fundamental foundations:

1. **Reflection & Struct Tags (`reflect`, `encoding/json`)**:
   - Go uses compile-time struct tags parsed at runtime via reflection to direct marshaling:
     ```go
     type User struct {
         ID        int64     `json:"id"`
         Email     string    `json:"email,omitempty"`
         Password  string    `json:"-"` // Never serialize
         CreatedAt time.Time `json:"created_at"`
     }
     ```
   - Unexported fields (lowercase) are **ignored** by external packages like `encoding/json` regardless of struct tags.

2. **In-Memory (`json.Marshal`/`json.Unmarshal`) vs Streaming (`json.Encoder`/`json.Decoder`)**:
   - `json.Marshal(v)` copies the entire payload into a heap-allocated `[]byte`. Suitable for small in-memory payloads (< 100KB).
   - `json.NewDecoder(r).Decode(&v)` streams tokens on-the-fly directly from any `io.Reader` (e.g., HTTP request body or OS file) without buffering the full document into RAM. Crucial for production services handling large or unbounded input.

3. **Atomic File Writes & Permissions**:
   - Writing directly to a production file via `os.Create(dest)` risks leaving a corrupted, truncated file if the process crashes mid-write.
   - **Production Pattern**: Write to a temporary file on the same filesystem (`os.CreateTemp`), flush and sync (`f.Sync()`), close, and atomically rename (`os.Rename(tempName, dest)`).

4. **Custom Encoding Contracts**:
   - Implement `json.Marshaler` (`MarshalJSON() ([]byte, error)`) and `json.Unmarshaler` (`UnmarshalJSON([]byte) error`) to handle non-standard schemas, such as Unix timestamps or strongly-typed money/currency units.

---

## 🎯 Learning Objectives

By the end of this stage, you will:
- [x] Configure precise JSON serialization with tags: names, `omitempty`, and `-` omission.
- [x] Choose between `json.Unmarshal` and streaming `json.Decoder` based on payload size and source.
- [x] Implement `json.Marshaler` and `json.Unmarshaler` without triggering infinite recursion.
- [x] Perform atomic, crash-resilient file persistence using `os.CreateTemp` and `os.Rename`.
- [x] Handle file errors properly using `errors.Is(err, os.ErrNotExist)`.

---

## 🗂️ Stage Structure

```
08-files-json/
├── README.md                           # Stage curriculum & theory
├── examples/
│   ├── file_operations/main.go         # Safe atomic write, reading, permission flags
│   └── json_streaming/main.go          # Streaming json.Decoder & Token iterator
├── exercises/
│   ├── 01-struct-tags/                 # Struct tags, omit empty, sensitive fields
│   ├── 02-custom-marshaler/            # Custom Unix epoch timestamp marshaler
│   └── 03-streaming-decoder/           # Memory-bounded JSON array stream processor
├── solutions/
│   ├── 01-struct-tags/                 # Reference implementation
│   ├── 02-custom-marshaler/            # Reference implementation
│   └── 03-streaming-decoder/           # Reference implementation
└── challenges/
    └── config-migrator/                # Atomic config manager with schema migration
```

---

## ⚠️ Common Pitfalls & Gotchas

1. **Unexported Struct Fields**:
   ```go
   type Config struct {
       port int // ❌ Unexported: json package cannot read or write this field!
       Port int `json:"port"` // ✅ Exported: Accessible to serialization packages
   }
   ```

2. **Infinite Recursion in `MarshalJSON`**:
   ```go
   // ❌ Infinite recursive loop (crashes with stack overflow):
   func (u User) MarshalJSON() ([]byte, error) {
       return json.Marshal(u) // Calls u.MarshalJSON() infinitely!
   }

   // ✅ Idiomatic alias pattern:
   func (u User) MarshalJSON() ([]byte, error) {
       type Alias User // Drops methods, avoiding infinite recursion
       return json.Marshal(&struct{
           Alias
           ExtraField string `json:"extra"`
       }{
           Alias: Alias(u),
           ExtraField: "computed",
       })
   }
   ```

3. **Ignoring `f.Sync()` before Atomic Rename**:
   - `os.Rename` updates metadata, but data may still linger in OS page cache. Always call `f.Sync()` before `f.Close()` and `os.Rename`.

---

## 🧪 Verification Commands

```powershell
# Run all tests in this stage
go test -v ./08-files-json/...

# Lint and check style
go vet ./08-files-json/...
go fmt ./08-files-json/...
```

---

## 🔗 Connections
- **Prerequisites**: [Stage 07: Standard Library Deep Dive](../07-standard-library/README.md).
- **Next Stage**: [Stage 09: HTTP Fundamentals](../09-http/README.md) (`http.Handler`, `http.ResponseWriter`, JSON API endpoints).
