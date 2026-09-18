# Exercise 03: Go Native Fuzz Testing with `testing.F` (Tier 3 - Hard)

## 🎯 Problem Statement
Parsers that handle raw byte streams from the network often suffer from index-out-of-range panics when handling malformed or truncated inputs.

Implement `ParseHeader`:
```go
type PacketHeader struct {
    Version byte
    Type    byte
    Length  uint16
}

func ParseHeader(data []byte) (PacketHeader, error)
```
And write a Go 1.18+ native fuzz test in [`starter_test.go`](./starter_test.go) verifying that `ParseHeader` never crashes or panics on arbitrary byte input.

### Requirements
1. `data` must have at least 4 bytes:
   - Byte 0: Version (must be $\ge 1$)
   - Byte 1: Type
   - Bytes 2-3: Big-endian `uint16` Length (`binary.BigEndian.Uint16(data[2:4])`)
2. Return a descriptive error if `len(data) < 4` or `data[0] < 1`.
3. In `FuzzParseHeader`:
   - Seed corpus with valid and truncated samples using `f.Add(...)`.
   - In `f.Fuzz(...)`, call `ParseHeader(data)`.
   - Assert that no panic occurs.

## 🛠️ Instructions
1. Open [`starter.go`](./starter.go) and implement `ParseHeader`.
2. Open [`starter_test.go`](./starter_test.go) and implement `FuzzParseHeader`.
3. Run tests:
   ```powershell
   go test -v ./12-testing/exercises/03-native-fuzzer
   ```
