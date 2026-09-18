# Exercise 02: Sliding TCP Idle Timeout

## Objective
Implement a line-based connection reader that maintains a sliding read deadline, protecting servers against slowloris attacks and orphaned TCP connections.

## Requirements
1. **Sliding Read Deadline**: Before each line read attempt, refresh the connection read deadline to `time.Now().Add(idleTimeout)`.
2. **Line Reading**: Read text line-by-line (delimited by `\n`). Trim trailing carriage returns and newlines before dispatching to `onLine`.
3. **Timeout Detection**:
   - If the read fails due to a network timeout (`netErr.Timeout() == true` or `os.ErrDeadlineExceeded`), terminate gracefully.
   - Return the count of complete lines read, `timedOut = true`, and `err = nil`.
4. **Clean Disconnect**:
   - If the client closes the connection (`io.EOF`), return lines read, `timedOut = false`, and `err = nil`.

## Starter & Tests
- Starter: `starter.go`
- Run starter tests: `go test -v ./15-networking/exercises/02-tcp-idle-timeout`
