package framing

import (
	"encoding/binary"
	"errors"
	"fmt"
	"io"
	"math"
)

var (
	ErrFrameTooLarge = errors.New("frame payload exceeds maximum permitted size")
)

// WriteFrame encodes payload length as a 4-byte BigEndian header and writes it followed by payload.
func WriteFrame(w io.Writer, payload []byte) error {
	if len(payload) > math.MaxUint32 {
		return ErrFrameTooLarge
	}

	header := make([]byte, 4)
	binary.BigEndian.PutUint32(header, uint32(len(payload)))

	// Write header
	if _, err := w.Write(header); err != nil {
		return fmt.Errorf("failed writing frame header: %w", err)
	}

	// Write payload
	if len(payload) > 0 {
		if _, err := w.Write(payload); err != nil {
			return fmt.Errorf("failed writing frame payload: %w", err)
		}
	}

	return nil
}

// ReadFrame reads a 4-byte BigEndian length header followed by the frame payload.
// It returns ErrFrameTooLarge if header length > maxFrameSize.
func ReadFrame(r io.Reader, maxFrameSize uint32) ([]byte, error) {
	header := make([]byte, 4)
	if _, err := io.ReadFull(r, header); err != nil {
		return nil, err
	}

	length := binary.BigEndian.Uint32(header)
	if length > maxFrameSize {
		return nil, ErrFrameTooLarge
	}

	if length == 0 {
		return []byte{}, nil
	}

	payload := make([]byte, length)
	if _, err := io.ReadFull(r, payload); err != nil {
		return nil, fmt.Errorf("failed reading frame payload: %w", err)
	}

	return payload, nil
}
