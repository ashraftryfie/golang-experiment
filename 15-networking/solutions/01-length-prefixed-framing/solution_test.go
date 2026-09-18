package framing

import (
	"bytes"
	"errors"
	"io"
	"testing"
)

func TestFraming_RoundTrip(t *testing.T) {
	buf := new(bytes.Buffer)
	payload := []byte("high-performance distributed protocol payload")

	if err := WriteFrame(buf, payload); err != nil {
		t.Fatalf("WriteFrame failed: %v", err)
	}

	readPayload, err := ReadFrame(buf, 1024)
	if err != nil {
		t.Fatalf("ReadFrame failed: %v", err)
	}

	if !bytes.Equal(payload, readPayload) {
		t.Fatalf("payload mismatch: got %s, want %s", readPayload, payload)
	}
}

func TestFraming_EmptyPayload(t *testing.T) {
	buf := new(bytes.Buffer)
	if err := WriteFrame(buf, []byte{}); err != nil {
		t.Fatalf("WriteFrame empty failed: %v", err)
	}

	readPayload, err := ReadFrame(buf, 1024)
	if err != nil {
		t.Fatalf("ReadFrame empty failed: %v", err)
	}

	if len(readPayload) != 0 {
		t.Fatalf("expected empty slice, got %d bytes", len(readPayload))
	}
}

func TestFraming_MultipleFramesStream(t *testing.T) {
	buf := new(bytes.Buffer)
	messages := [][]byte{
		[]byte("frame-1"),
		[]byte("frame-2-longer-text"),
		[]byte("frame-3"),
	}

	for _, msg := range messages {
		if err := WriteFrame(buf, msg); err != nil {
			t.Fatalf("failed writing %s: %v", msg, err)
		}
	}

	for i, expected := range messages {
		actual, err := ReadFrame(buf, 1024)
		if err != nil {
			t.Fatalf("failed reading frame %d: %v", i, err)
		}
		if !bytes.Equal(expected, actual) {
			t.Fatalf("frame %d mismatch: got %s, want %s", i, actual, expected)
		}
	}
}

func TestFraming_ExceedsMaxFrameSize(t *testing.T) {
	buf := new(bytes.Buffer)
	payload := bytes.Repeat([]byte("A"), 500)

	if err := WriteFrame(buf, payload); err != nil {
		t.Fatalf("WriteFrame failed: %v", err)
	}

	// Max permitted frame size is 200, so reading 500 bytes should trigger ErrFrameTooLarge
	_, err := ReadFrame(buf, 200)
	if !errors.Is(err, ErrFrameTooLarge) {
		t.Fatalf("expected ErrFrameTooLarge, got %v", err)
	}
}

func TestFraming_UnexpectedEOF(t *testing.T) {
	// Truncated header (only 2 bytes instead of 4)
	buf := bytes.NewReader([]byte{0x00, 0x01})
	_, err := ReadFrame(buf, 1024)
	if !errors.Is(err, io.ErrUnexpectedEOF) {
		t.Fatalf("expected io.ErrUnexpectedEOF, got %v", err)
	}
}
