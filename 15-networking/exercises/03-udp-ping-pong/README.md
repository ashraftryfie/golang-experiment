# Exercise 03: Connectionless UDP Ping-Pong Service

## Objective
Implement a connectionless UDP server and client in Go that communicates via datagrams using `net.PacketConn`.

## Requirements
1. **UDP Server**:
   - `StartServer(addr string) (*PingPongServer, error)`
   - Listens on UDP using `net.ListenPacket("udp", addr)`.
   - In a background goroutine, reads packets up to 512 bytes.
   - For every packet starting with `"PING:<id>"`, replies to `peerAddr` with `"PONG:<id>"`.
   - `Close() error` closes the listener and halts the background loop cleanly.
2. **UDP Client Function**:
   - `SendPing(serverAddr string, id string, timeout time.Duration) (string, error)`
   - Dials or communicates with UDP `serverAddr`.
   - Sends `"PING:" + id`.
   - Enforces a read deadline of `time.Now().Add(timeout)`.
   - Returns the received ID if `"PONG:<id>"` matches, or an error if timeout or malformed.

## Starter & Tests
- Starter: `starter.go`
- Run starter tests: `go test -v ./15-networking/exercises/03-udp-ping-pong`
