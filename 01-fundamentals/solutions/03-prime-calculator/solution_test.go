package primes

import (
	"reflect"
	"testing"
)

func TestSolutionPrimes(t *testing.T) {
	if !IsPrime(29) {
		t.Errorf("expected 29 to be prime")
	}
	if IsPrime(30) {
		t.Errorf("expected 30 to not be prime")
	}

	primes := SieveOfEratosthenes(15)
	expected := []int{2, 3, 5, 7, 11, 13}
	if !reflect.DeepEqual(primes, expected) {
		t.Errorf("got %v, want %v", primes, expected)
	}

	factors, err := PrimeFactors(84)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	expectedFactors := []int{2, 2, 3, 7}
	if !reflect.DeepEqual(factors, expectedFactors) {
		t.Errorf("got %v, want %v", factors, expectedFactors)
	}
}
