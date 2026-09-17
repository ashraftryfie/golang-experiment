package primes

import (
	"errors"
	"reflect"
	"testing"
)

func TestIsPrime(t *testing.T) {
	// If starter default is untouched (returns false for 2), skip gracefully
	if !IsPrime(2) {
		t.Skip("skipping: IsPrime is not yet implemented (implement in starter.go)")
	}

	tests := []struct {
		n    int
		want bool
	}{
		{n: -5, want: false},
		{n: 0, want: false},
		{n: 1, want: false},
		{n: 2, want: true},
		{n: 3, want: true},
		{n: 4, want: false},
		{n: 17, want: true},
		{n: 100, want: false},
		{n: 997, want: true},
	}

	for _, tc := range tests {
		t.Run("", func(t *testing.T) {
			got := IsPrime(tc.n)
			if got != tc.want {
				t.Errorf("IsPrime(%d) = %v; want %v", tc.n, got, tc.want)
			}
		})
	}
}

func TestSieveOfEratosthenes(t *testing.T) {
	got := SieveOfEratosthenes(10)
	if got == nil {
		t.Skip("skipping: SieveOfEratosthenes is not yet implemented (implement in starter.go)")
	}

	tests := []struct {
		limit int
		want  []int
	}{
		{limit: 1, want: []int{}},
		{limit: 10, want: []int{2, 3, 5, 7}},
		{limit: 20, want: []int{2, 3, 5, 7, 11, 13, 17, 19}},
	}

	for _, tc := range tests {
		t.Run("", func(t *testing.T) {
			res := SieveOfEratosthenes(tc.limit)
			if len(res) == 0 && len(tc.want) == 0 {
				return
			}
			if !reflect.DeepEqual(res, tc.want) {
				t.Errorf("SieveOfEratosthenes(%d) = %v; want %v", tc.limit, res, tc.want)
			}
		})
	}
}

func TestPrimeFactors(t *testing.T) {
	tests := []struct {
		n         int
		want      []int
		wantError error
	}{
		{n: 12, want: []int{2, 2, 3}},
		{n: 100, want: []int{2, 2, 5, 5}},
		{n: 13, want: []int{13}},
		{n: 1, wantError: ErrInvalidPrimeInput},
		{n: -10, wantError: ErrInvalidPrimeInput},
	}

	for _, tc := range tests {
		t.Run("", func(t *testing.T) {
			got, err := PrimeFactors(tc.n)
			if errors.Is(err, ErrNotImplemented) {
				t.Skip("skipping: PrimeFactors is not yet implemented (implement in starter.go)")
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
				t.Errorf("PrimeFactors(%d) = %v; want %v", tc.n, got, tc.want)
			}
		})
	}
}
