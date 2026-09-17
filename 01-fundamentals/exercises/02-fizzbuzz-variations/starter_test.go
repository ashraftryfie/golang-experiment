package fizzbuzz

import (
	"errors"
	"reflect"
	"testing"
)

func TestClassicFizzBuzz(t *testing.T) {
	tests := []struct {
		name      string
		n         int
		want      []string
		wantError error
	}{
		{
			name: "up to 5",
			n:    5,
			want: []string{"1", "2", "Fizz", "4", "Buzz"},
		},
		{
			name: "up to 15 contains FizzBuzz",
			n:    15,
			want: []string{
				"1", "2", "Fizz", "4", "Buzz",
				"Fizz", "7", "8", "Fizz", "Buzz",
				"11", "Fizz", "13", "14", "FizzBuzz",
			},
		},
		{
			name:      "invalid n zero",
			n:         0,
			wantError: ErrInvalidRange,
		},
		{
			name:      "invalid n negative",
			n:         -5,
			wantError: ErrInvalidRange,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got, err := ClassicFizzBuzz(tc.n)
			if errors.Is(err, ErrNotImplemented) {
				t.Skip("skipping: ClassicFizzBuzz is not yet implemented")
			}
			if tc.wantError != nil {
				if !errors.Is(err, tc.wantError) {
					t.Fatalf("expected error %v, got %v", tc.wantError, err)
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if !reflect.DeepEqual(got, tc.want) {
				t.Errorf("got %v, want %v", got, tc.want)
			}
		})
	}
}

func TestCustomFizzBuzz(t *testing.T) {
	rules := []Rule{
		{Divisor: 2, Word: "Foo"},
		{Divisor: 3, Word: "Bar"},
	}

	got, err := CustomFizzBuzz(6, rules)
	if errors.Is(err, ErrNotImplemented) {
		t.Skip("skipping: CustomFizzBuzz is not yet implemented")
	}
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	want := []string{"1", "Foo", "Bar", "Foo", "5", "FooBar"}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("CustomFizzBuzz() = %v, want %v", got, want)
	}

	// Test invalid rule divisor
	badRules := []Rule{{Divisor: 0, Word: "Crash"}}
	if _, err := CustomFizzBuzz(5, badRules); !errors.Is(err, ErrInvalidRule) {
		t.Errorf("expected ErrInvalidRule, got %v", err)
	}
}
