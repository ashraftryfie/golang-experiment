package idletimeout

import (
	"errors"
	"io"
	"net"
	"testing"
	"time"
)

func TestReadLinesWithIdleTimeout_Success(t *testing.T) {
	serverConn, clientConn := net.Pipe()
	defer serverConn.Close()

	go func() {
		defer clientConn.Close()
		_, _ = io.WriteString(clientConn, "line1\nline2\r\nline3\n")
	}()

	var lines []string
	count, timedOut, err := ReadLinesWithIdleTimeout(serverConn, 200*time.Millisecond, func(line string) error {
		lines = append(lines, line)
		return nil
	})

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if timedOut {
		t.Fatalf("expected timedOut=false")
	}
	if count != 3 {
		t.Fatalf("expected 3 lines, got %d", count)
	}
	if lines[0] != "line1" || lines[1] != "line2" || lines[2] != "line3" {
		t.Fatalf("unexpected line contents: %v", lines)
	}
}

func TestReadLinesWithIdleTimeout_IdleTimeoutTriggered(t *testing.T) {
	serverConn, clientConn := net.Pipe()
	defer serverConn.Close()
	defer clientConn.Close()

	go func() {
		// Send one line and then sleep longer than the idle timeout
		_, _ = io.WriteString(clientConn, "first line\n")
		time.Sleep(100 * time.Millisecond)
	}()

	var lines []string
	count, timedOut, err := ReadLinesWithIdleTimeout(serverConn, 30*time.Millisecond, func(line string) error {
		lines = append(lines, line)
		return nil
	})

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !timedOut {
		t.Fatalf("expected timedOut=true")
	}
	if count != 1 {
		t.Fatalf("expected 1 line read before timeout, got %d", count)
	}
}

func TestReadLinesWithIdleTimeout_CallbackError(t *testing.T) {
	serverConn, clientConn := net.Pipe()
	defer serverConn.Close()

	go func() {
		defer clientConn.Close()
		_, _ = io.WriteString(clientConn, "ok\nbad\nignored\n")
	}()

	errRejected := errors.New("rejected")
	count, timedOut, err := ReadLinesWithIdleTimeout(serverConn, 100*time.Millisecond, func(line string) error {
		if line == "bad" {
			return errRejected
		}
		return nil
	})

	if !errors.Is(err, errRejected) {
		t.Fatalf("expected errRejected, got %v", err)
	}
	if timedOut {
		t.Fatalf("expected timedOut=false")
	}
	if count != 1 {
		t.Fatalf("expected count 1, got %d", count)
	}
}
