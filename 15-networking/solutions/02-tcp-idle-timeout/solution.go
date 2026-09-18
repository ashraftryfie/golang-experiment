package idletimeout

import (
	"bufio"
	"errors"
	"io"
	"net"
	"os"
	"strings"
	"time"
)

// ReadLinesWithIdleTimeout reads lines from conn until EOF or until the idleTimeout is exceeded between reads.
// It returns the number of lines read, a boolean indicating if it timed out, and any non-timeout error.
func ReadLinesWithIdleTimeout(conn net.Conn, idleTimeout time.Duration, onLine func(line string) error) (int, bool, error) {
	reader := bufio.NewReader(conn)
	linesRead := 0

	for {
		// Set sliding deadline before each line read
		if err := conn.SetReadDeadline(time.Now().Add(idleTimeout)); err != nil {
			return linesRead, false, err
		}

		line, err := reader.ReadString('\n')
		if err != nil {
			// Check if this was caused by a timeout
			var netErr net.Error
			if (errors.As(err, &netErr) && netErr.Timeout()) || errors.Is(err, os.ErrDeadlineExceeded) {
				return linesRead, true, nil
			}

			// If EOF occurred with some trailing content
			if errors.Is(err, io.EOF) {
				cleanLine := strings.TrimRight(line, "\r\n")
				if len(cleanLine) > 0 {
					if cbErr := onLine(cleanLine); cbErr != nil {
						return linesRead, false, cbErr
					}
					linesRead++
				}
				return linesRead, false, nil
			}

			return linesRead, false, err
		}

		cleanLine := strings.TrimRight(line, "\r\n")
		if err := onLine(cleanLine); err != nil {
			return linesRead, false, err
		}
		linesRead++
	}
}
