package logscan

import (
	"errors"
	"reflect"
	"strings"
	"testing"
)

func TestFilterLogs(t *testing.T) {
	logData := `2026-09-18 10:00:01 [INFO] Service started
2026-09-18 10:00:02 [ERROR] Database connection failed
2026-09-18 10:00:03 [WARN] Retrying connection
2026-09-18 10:00:04 [ERROR] Connection timed out
2026-09-18 10:00:05 [INFO] Worker pool idle`

	r := strings.NewReader(logData)
	got, err := FilterLogs(r, "[ERROR]")
	if errors.Is(err, ErrNotImplemented) {
		t.Skip("skipping: FilterLogs is not yet implemented (implement in starter.go)")
	}
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	want := []string{
		"2026-09-18 10:00:02 [ERROR] Database connection failed",
		"2026-09-18 10:00:04 [ERROR] Connection timed out",
	}

	if !reflect.DeepEqual(got, want) {
		t.Errorf("got %v, want %v", got, want)
	}

	// Empty target check
	if _, err := FilterLogs(strings.NewReader(""), ""); !errors.Is(err, ErrEmptyTargetLevel) {
		t.Errorf("expected ErrEmptyTargetLevel, got %v", err)
	}
}
