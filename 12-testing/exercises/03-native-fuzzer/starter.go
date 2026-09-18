package fuzzer

import (
	"errors"
)

var (
	ErrNotImplemented = errors.New("TODO: implement ParseHeader")
	ErrTooShort       = errors.New("header must be at least 4 bytes")
	ErrInvalidVersion = errors.New("header version must be >= 1")
)

type PacketHeader struct {
	Version byte
	Type    byte
	Length  uint16
}

// ParseHeader extracts a 4-byte packet header safely.
func ParseHeader(data []byte) (PacketHeader, error) {
	// TODO: Safely check lengths, decode fields using binary.BigEndian, prevent panics
	return PacketHeader{}, ErrNotImplemented
}
