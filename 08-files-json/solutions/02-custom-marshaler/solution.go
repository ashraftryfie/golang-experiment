package epochtime_solution

import (
	"fmt"
	"strconv"
	"strings"
	"time"
)

// EpochTime wraps time.Time to serialize as an integer Unix timestamp (seconds).
type EpochTime struct {
	time.Time
}

// MarshalJSON converts the embedded time into a JSON integer byte slice.
func (t EpochTime) MarshalJSON() ([]byte, error) {
	return []byte(fmt.Sprintf("%d", t.Unix())), nil
}

// UnmarshalJSON parses a JSON integer byte slice into time.Time.
func (t *EpochTime) UnmarshalJSON(data []byte) error {
	str := strings.TrimSpace(string(data))
	if str == "null" || str == `""` || str == "" {
		t.Time = time.Time{}
		return nil
	}

	sec, err := strconv.ParseInt(str, 10, 64)
	if err != nil {
		return fmt.Errorf("invalid epoch timestamp %q: %w", str, err)
	}

	t.Time = time.Unix(sec, 0).UTC()
	return nil
}
