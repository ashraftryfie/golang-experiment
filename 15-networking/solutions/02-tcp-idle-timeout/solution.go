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
		// Set sliding deadline before reading from network if buffer is empty
		if reader.Buffered() == 0 {
			if err := conn.SetReadDeadline(time.Now().Add(idleTimeout)); err != nil {
				if errors.Is(err, io.ErrClosedPipe) || errors.Is(err, io.EOF) {
					return linesRead, false, nil
				}
				return linesRead, false, err
			}
		}

		line, err := reader.ReadString('\n')
		if err != nil {
			var netErr net.Error
			if (errors.As(err, &netErr) && netErr.Timeout()) ||
				errors.Is(err, os.ErrDeadlineExceeded) {
				return linesRead, true, nil
			}

			// net.Pipe may return buffered data with io.ErrClosedPipe when
			// the peer closes immediately after writing.
			if errors.Is(err, io.EOF) || errors.Is(err, io.ErrClosedPipe) {
				if len(line) > 0 {
					cleanLine := strings.TrimRight(line, "\r\n")
					if cleanLine != "" {
						if cbErr := onLine(cleanLine); cbErr != nil {
							return linesRead, false, cbErr
						}
						linesRead++
					}
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
