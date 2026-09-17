package sliceops

import (
	"errors"
	"reflect"
	"testing"
)

func TestDeduplicate(t *testing.T) {
	input := []int{1, 2, 2, 3, 1, 4, 5, 4}
	got := Deduplicate(input)
	if got == nil {
		t.Skip("skipping: Deduplicate is not yet implemented (implement in starter.go)")
	}

	want := []int{1, 2, 3, 4, 5}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("Deduplicate() = %v, want %v", got, want)
	}
}

func TestFilterEven(t *testing.T) {
	input := []int{1, 2, 3, 4, 5, 6}
	got := FilterEven(input)
	if got == nil {
		t.Skip("skipping: FilterEven is not yet implemented (implement in starter.go)")
	}

	want := []int{2, 4, 6}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("FilterEven() = %v, want %v", got, want)
	}
}

func TestRemoveAtIndex(t *testing.T) {
	input := []int{10, 20, 30, 40}
	got, err := RemoveAtIndex(input, 2)
	if errors.Is(err, ErrNotImplemented) {
		t.Skip("skipping: RemoveAtIndex is not yet implemented (implement in starter.go)")
	}
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	want := []int{10, 20, 40}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("RemoveAtIndex() = %v, want %v", got, want)
	}

	// Test boundary error
	if _, err := RemoveAtIndex(input, 10); !errors.Is(err, ErrIndexOutOfRange) {
		t.Errorf("expected ErrIndexOutOfRange, got %v", err)
	}
}
