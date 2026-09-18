# Stage 15: Networking & Sockets

Welcome to **Stage 15** of your Go journey. While HTTP and gRPC are high-level abstractions, building real-time engines, custom binary protocols, database drivers, proxies, and game servers requires direct mastery of low-level networking via the standard library's `net` package.

---

## 1. Core Abstractions in `net`

The `net` package provides unified interfaces across OS networking primitives:

```go
type Listener interface {
    Accept() (Conn, error)
    Close() error
    Addr() Addr
}

type Conn interface {
    Read(b []byte) (n int, err error)
    Write(b []byte) (n int, err error)
    Close() error
    LocalAddr() Addr
    RemoteAddr() Addr
    SetDeadline(t time.Time) error
    SetReadDeadline(t time.Time) error
    SetWriteDeadline(t time.Time) error
}

type PacketConn interface {
    ReadFrom(p []byte) (n int, addr Addr, err error)
    WriteTo(p []byte, addr Addr) (n int, err error)
    Close() error
    LocalAddr() Addr
    SetDeadline(t time.Time) error
    SetReadDeadline(t time.Time) error
    SetWriteDeadline(t time.Time) error
}
```

---

## 2. Key Networking Mental Models

### A. TCP is a Byte Stream, NOT a Message Stream
TCP provides reliable, ordered transmission of bytes. It does **not** preserve write boundaries:
- If a client calls `Write([]byte("hello"))` and `Write([]byte("world"))`, the server might receive:
  - `"helloworld"` in a single `Read`.
  - `"he"`, `"llo"`, `"world"` across three `Read` calls.
- **Solution: Protocol Framing**:
  1. **Delimiter-based** (e.g. `\n` or `\r\n` read with `bufio.Scanner`).
  2. **Length-prefixed** (e.g. 4-byte uint32 BigEndian size followed by payload, read with `io.ReadFull`).

### B. Deadlines vs Timeouts
Go network connections do not take duration timeouts; they take absolute timestamps:
```go
// Set deadline 5 seconds in the future
conn.SetReadDeadline(time.Now().Add(5 * time.Second))

// To clear a deadline:
conn.SetReadDeadline(time.Time{})
```
If a deadline expires, the ongoing or subsequent `Read` or `Write` returns an error where `errors.Is(err, os.ErrDeadlineExceeded)` or `netErr.Timeout() == true`.

### C. Graceful Server Shutdown
Closing `listener.Close()` stops accepting new connections but leaves active connections open. A production server must:
1. Close the listener so no new sockets arrive.
2. Signal active connection goroutines (or close their sockets).
3. Wait for active handlers to drain via `sync.WaitGroup`.

---

## 3. Directory Layout

```
15-networking/
├── README.md
├── examples/
│   ├── tcp_echo/main.go           # Concurrent TCP server with deadlines
│   └── udp_heartbeat/main.go      # UDP packet sender and listener
├── exercises/
│   ├── 01-length-prefixed-framing # Binary framing encoder and decoder
│   ├── 02-tcp-idle-timeout        # Sliding idle read deadline handler
│   └── 03-udp-ping-pong           # Connectionless UDP client/server
├── solutions/
│   ├── 01-length-prefixed-framing
│   ├── 02-tcp-idle-timeout
│   └── 03-udp-ping-pong
└── challenges/
    └── chat-broker/               # Concurrent TCP broadcast chat server
```
