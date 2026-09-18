package framing

import (
	"errors"
	"io"
)

var (
	ErrNotImplemented = errors.New("exercise not implemented yet")
	ErrFrameTooLarge   = errors.New("frame payload exceeds maximum permitted size")
)

// WriteFrame encodes payload length as a 4-byte BigEndian header and writes it followed by payload.
func WriteFrame(w io.Writer, payload []byte) error {
	return ErrNotImplemented
}

// ReadFrame reads a 4-byte BigEndian length header followed by the frame payload.
// It returns ErrFrameTooLarge if header length > maxFrameSize.
func ReadFrame(r io.Reader, maxFrameSize uint32) ([]byte, error) {
	return nil, ErrNotImplemented
}
