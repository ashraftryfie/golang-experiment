package fuzzer_solution

import (
	"encoding/binary"
	"errors"
)

var (
	ErrTooShort       = errors.New("header must be at least 4 bytes")
	ErrInvalidVersion = errors.New("header version must be >= 1")
)

type PacketHeader struct {
	Version byte
	Type    byte
	Length  uint16
}

// ParseHeader safely extracts a 4-byte packet header without panic risk.
func ParseHeader(data []byte) (PacketHeader, error) {
	if len(data) < 4 {
		return PacketHeader{}, ErrTooShort
	}

	if data[0] < 1 {
		return PacketHeader{}, ErrInvalidVersion
	}

	return PacketHeader{
		Version: data[0],
		Type:    data[1],
		Length:  binary.BigEndian.Uint16(data[2:4]),
	}, nil
}
