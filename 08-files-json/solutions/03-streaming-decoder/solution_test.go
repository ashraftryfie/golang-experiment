package streamfilter_solution

import (
	"bytes"
	"encoding/json"
	"strings"
	"testing"
)

func TestFilterHighValueSolution(t *testing.T) {
	rawInput := `[
		{"id": "tx-1", "amount": 50.00, "sender": "alice"},
		{"id": "tx-2", "amount": 1000.00, "sender": "bob"},
		{"id": "tx-3", "amount": 150.50, "sender": "charlie"},
		{"id": "tx-4", "amount": 5000.00, "sender": "diana"}
	]`

	var output bytes.Buffer
	count, err := FilterHighValue(&output, strings.NewReader(rawInput), 200.00)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if count != 2 {
		t.Errorf("expected 2 high value transactions, got %d", count)
	}

	var filtered []Transaction
	if err := json.Unmarshal(output.Bytes(), &filtered); err != nil {
		t.Fatalf("output is not valid JSON: %v, body=%s", err, output.String())
	}

	if len(filtered) != 2 {
		t.Fatalf("unmarshaled len = %d, want 2", len(filtered))
	}

	if filtered[0].ID != "tx-2" || filtered[1].ID != "tx-4" {
		t.Errorf("filtered items mismatch: %+v", filtered)
	}
}
