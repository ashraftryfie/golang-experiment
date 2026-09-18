package epochtime

import (
	"errors"
	"time"
)

var ErrNotImplemented = errors.New("TODO: implement EpochTime MarshalJSON / UnmarshalJSON")

// EpochTime wraps time.Time to serialize as an integer Unix timestamp (seconds).
type EpochTime struct {
	time.Time
}

// MarshalJSON converts the embedded time into a JSON integer byte slice.
func (t EpochTime) MarshalJSON() ([]byte, error) {
	// TODO: Return integer Unix seconds encoded in []byte
	return nil, ErrNotImplemented
}

// UnmarshalJSON parses a JSON integer byte slice into time.Time.
func (t *EpochTime) UnmarshalJSON(data []byte) error {
	// TODO: Parse integer from data and set t.Time = time.Unix(sec, 0)
	return ErrNotImplemented
}
