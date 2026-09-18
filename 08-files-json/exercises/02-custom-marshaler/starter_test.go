package epochtime

import (
	"encoding/json"
	"errors"
	"testing"
	"time"
)

type EventPayload struct {
	Name      string    `json:"name"`
	Timestamp EpochTime `json:"timestamp"`
}

func TestEpochTime(t *testing.T) {
	fixedTime := time.Date(2024, 3, 15, 12, 0, 0, 0, time.UTC)
	expectedUnix := fixedTime.Unix()

	event := EventPayload{
		Name:      "deploy",
		Timestamp: EpochTime{Time: fixedTime},
	}

	data, err := json.Marshal(event)
	if errors.Is(err, ErrNotImplemented) {
		t.Skip("skipping: MarshalJSON is not yet implemented (implement in starter.go)")
	}
	if err != nil {
		t.Fatalf("marshal error: %v", err)
	}

	// Verify unmarshaling back
	var decoded EventPayload
	if err := json.Unmarshal(data, &decoded); err != nil {
		t.Fatalf("unmarshal error: %v", err)
	}

	if decoded.Timestamp.Unix() != expectedUnix {
		t.Errorf("timestamp mismatch: got %d, want %d", decoded.Timestamp.Unix(), expectedUnix)
	}

	if decoded.Name != "deploy" {
		t.Errorf("name mismatch: got %s, want deploy", decoded.Name)
	}
}
