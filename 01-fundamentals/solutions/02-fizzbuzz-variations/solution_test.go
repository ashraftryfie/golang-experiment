package fizzbuzz

import (
	"reflect"
	"testing"
)

func TestSolutionClassicFizzBuzz(t *testing.T) {
	got, err := ClassicFizzBuzz(5)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	want := []string{"1", "2", "Fizz", "4", "Buzz"}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("got %v, want %v", got, want)
	}
}

func TestSolutionCustomFizzBuzz(t *testing.T) {
	rules := []Rule{
		{Divisor: 2, Word: "Tick"},
		{Divisor: 4, Word: "Tock"},
	}
	got, err := CustomFizzBuzz(4, rules)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	want := []string{"1", "Tick", "3", "TickTock"}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("got %v, want %v", got, want)
	}
}
