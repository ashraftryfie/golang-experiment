package idletimeout

import (
	"errors"
	"io"
	"net"
	"testing"
	"time"
)

func TestReadLinesWithIdleTimeout_Starter(t *testing.T) {
	serverConn, clientConn := net.Pipe()
	defer serverConn.Close()

	go func() {
		defer clientConn.Close()
		_, _ = io.WriteString(clientConn, "hello\n")
	}()

	var received []string
	lines, timedOut, err := ReadLinesWithIdleTimeout(serverConn, 50*time.Millisecond, func(line string) error {
		received = append(received, line)
		return nil
	})

	if errors.Is(err, ErrNotImplemented) {
		t.Skip("Skipping unimplemented exercise: ReadLinesWithIdleTimeout")
	}

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if lines != 1 || timedOut {
		t.Errorf("expected 1 line read and timedOut=false, got lines=%d, timedOut=%v", lines, timedOut)
	}
}
