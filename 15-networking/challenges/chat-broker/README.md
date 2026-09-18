# Engineering Challenge: Concurrent TCP Chat Broker

## Challenge Overview
Real-time collaborative applications, pub-sub messaging fabrics, and gaming backends rely on low-latency, persistent TCP connections with multi-client multiplexing.

You must build a high-performance **Concurrent TCP Chat Server / Pub-Sub Broker** with:
1. **Connection Lifecycle Management**:
   - Accepts concurrent TCP connections.
   - Initial handshake: First line received must be the user's nickname.
   - Prevents duplicate nicknames; returns an error and closes connection if collision occurs.
2. **Real-time Broadcast**:
   - When a user joins: broadcasts `JOIN: <nick>\n` to all other clients.
   - When a user sends text: broadcasts `MSG: <nick>: <message>\n` to all other clients.
   - When a user leaves (via `/quit` or EOF): broadcasts `LEAVE: <nick>\n` to remaining clients.
3. **Thread-Safe Architecture**:
   - Safe concurrent registration, message broadcasting, and deregistration.
   - Outbound message queues preventing slow writers from blocking the entire broker.
4. **Graceful Shutdown**:
   - `Stop()` cleanly closes listener, disconnects all active clients, and waits for all goroutines to terminate without leaking sockets or goroutines.

---

## Protocol Specification

| Action | Sent by Client | Server Response to Client | Broadcast to Others |
| :--- | :--- | :--- | :--- |
| Handshake | `<nick>\n` | `OK\n` | `JOIN: <nick>\n` |
| Duplicate Handshake | `<nick>\n` | `ERR: nickname already taken\n` | *(none)* |
| Message | `<text>\n` | *(none)* | `MSG: <nick>: <text>\n` |
| Disconnect | `/quit\n` | `BYE\n` | `LEAVE: <nick>\n` |

---

## Running Tests
Run tests with race detection:
```bash
go test -v -race ./15-networking/challenges/chat-broker/...
```
