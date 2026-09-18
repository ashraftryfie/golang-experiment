# Exercise 01: Length-Prefixed Binary Framing

## Objective
Implement binary length-prefixed framing for TCP streams. Because TCP provides a stream of bytes without message delimiters, applications must frame messages so receivers can reconstruct individual messages reliably.

## Requirements
1. **Header Format**: A 4-byte BigEndian `uint32` encoding the payload size in bytes.
2. **`WriteFrame(w io.Writer, payload []byte) error`**:
   - Encodes length header and writes header + payload.
   - If payload length exceeds `math.MaxUint32`, return `ErrFrameTooLarge`.
3. **`ReadFrame(r io.Reader, maxFrameSize uint32) ([]byte, error)`**:
   - Uses `io.ReadFull` to read the 4-byte header.
   - Decodes payload size. If length > `maxFrameSize`, return `ErrFrameTooLarge` to prevent memory exhaustion attacks.
   - Uses `io.ReadFull` to read exactly the payload bytes.
   - Handles EOF properly.

## Starter & Tests
- Starter: `starter.go`
- Run starter tests: `go test -v ./15-networking/exercises/01-length-prefixed-framing`
