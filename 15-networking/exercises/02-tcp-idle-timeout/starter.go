package idletimeout

import (
	"errors"
	"net"
	"time"
)

var ErrNotImplemented = errors.New("exercise not implemented yet")

// ReadLinesWithIdleTimeout reads lines from conn until EOF or until the idleTimeout is exceeded between reads.
// It returns the number of lines read, a boolean indicating if it timed out, and any non-timeout error.
func ReadLinesWithIdleTimeout(conn net.Conn, idleTimeout time.Duration, onLine func(line string) error) (int, bool, error) {
	return 0, false, ErrNotImplemented
}
