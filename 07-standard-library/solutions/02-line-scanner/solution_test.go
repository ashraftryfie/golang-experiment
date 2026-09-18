package logscan

import (
	"strings"
	"testing"
)

func TestSolutionFilterLogs(t *testing.T) {
	data := "line 1 [WARN]\nline 2 [INFO]\nline 3 [WARN]"
	lines, err := FilterLogs(strings.NewReader(data), "[WARN]")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(lines) != 2 {
		t.Errorf("expected 2 matching lines, got %d", len(lines))
	}
}
