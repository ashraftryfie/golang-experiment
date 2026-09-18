# Exercise 01: Custom ROT13 Cipher Stream (Tier 1 - Easy)

## 🎯 Problem Statement
Implement a custom `io.Reader` wrapper that applies the classic ROT13 cipher to bytes read from an underlying stream:

```go
type Rot13Reader struct {
    r io.Reader
}

func NewRot13Reader(r io.Reader) *Rot13Reader
func (rot *Rot13Reader) Read(p []byte) (n int, err error)
```

### ROT13 Rules
- For lowercase letters `'a'` through `'z'`, shift by 13 positions circularly.
- For uppercase letters `'A'` through `'Z'`, shift by 13 positions circularly.
- Non-alphabetic characters (spaces, punctuation, digits) remain unchanged.
- Must correctly delegate reading to the underlying `io.Reader` and handle `io.EOF`.

---

## 🛠️ Instructions
1. Open [`starter.go`](./starter.go).
2. Implement `Read` by reading into `p` and transforming bytes in-place up to `n`.
3. Run tests:
   ```powershell
   go test -v ./07-standard-library/exercises/01-custom-reader
   ```
