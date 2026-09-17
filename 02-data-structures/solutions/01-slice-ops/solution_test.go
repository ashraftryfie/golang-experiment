package sliceops

import (
	"errors"
	"reflect"
	"testing"
)

func TestSolutionDeduplicate(t *testing.T) {
	input := []int{3, 1, 2, 3, 1, 4}
	got := Deduplicate(input)
	want := []int{3, 1, 2, 4}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("got %v, want %v", got, want)
	}
}

func TestSolutionFilterEven(t *testing.T) {
	input := []int{10, 11, 12, 13, 14}
	got := FilterEven(input)
	want := []int{10, 12, 14}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("got %v, want %v", got, want)
	}
}

func TestSolutionRemoveAtIndex(t *testing.T) {
	input := []int{1, 2, 3}
	got, err := RemoveAtIndex(input, 1)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	want := []int{1, 3}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("got %v, want %v", got, want)
	}

	_, err = RemoveAtIndex(input, -1)
	if !errors.Is(err, ErrIndexOutOfRange) {
		t.Errorf("expected ErrIndexOutOfRange, got %v", err)
	}
}
