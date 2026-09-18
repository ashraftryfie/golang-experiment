package framing

import (
	"bytes"
	"errors"
	"testing"
)

func TestFraming_Starter(t *testing.T) {
	buf := new(bytes.Buffer)
	payload := []byte("hello network framing")

	err := WriteFrame(buf, payload)
	if errors.Is(err, ErrNotImplemented) {
		t.Skip("Skipping unimplemented exercise: WriteFrame")
	}
	if err != nil {
		t.Fatalf("unexpected write error: %v", err)
	}

	readPayload, err := ReadFrame(buf, 1024)
	if errors.Is(err, ErrNotImplemented) {
		t.Skip("Skipping unimplemented exercise: ReadFrame")
	}
	if err != nil {
		t.Fatalf("unexpected read error: %v", err)
	}

	if !bytes.Equal(payload, readPayload) {
		t.Errorf("expected %s, got %s", payload, readPayload)
	}
}
