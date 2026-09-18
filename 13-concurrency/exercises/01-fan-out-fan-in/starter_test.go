package fanoutfanin

import (
	"errors"
	"sort"
	"testing"
)

func TestFanOutFanIn(t *testing.T) {
	inputs := []int{1, 2, 3, 4, 5, 6, 7, 8}
	square := func(x int) int { return x * x }

	results, err := FanOutFanIn(inputs, 3, square)
	if errors.Is(err, ErrNotImplemented) {
		t.Skip("skipping: FanOutFanIn is not yet implemented (implement in starter.go)")
	}
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(results) != len(inputs) {
		t.Fatalf("expected %d results, got %d", len(inputs), len(results))
	}

	// Because concurrent ordering can vary, sort before comparing
	sort.Ints(results)
	expected := []int{1, 4, 9, 16, 25, 36, 49, 64}

	for i := range expected {
		if results[i] != expected[i] {
			t.Errorf("at index %d: got %d, want %d", i, results[i], expected[i])
		}
	}
}

func TestFanOutFanInEmpty(t *testing.T) {
	results, err := FanOutFanIn([]int{}, 4, func(x int) int { return x })
	if errors.Is(err, ErrNotImplemented) {
		t.Skip("skipping: FanOutFanIn is not yet implemented")
	}
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if results == nil || len(results) != 0 {
		t.Errorf("expected empty slice, got %v", results)
	}
}
