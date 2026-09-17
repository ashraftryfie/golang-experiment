package primes

import (
	"errors"
)

// ErrNotImplemented is returned by starter templates before you implement them.
var ErrNotImplemented = errors.New("TODO: implement function")

// ErrInvalidPrimeInput is returned when factorization input is less than 2.
var ErrInvalidPrimeInput = errors.New("input must be greater than or equal to 2")

// IsPrime returns true if n is a prime number.
func IsPrime(n int) bool {
	// TODO: Return false for n <= 1
	// TODO: Handle 2 and 3
	// TODO: Test factors up to sqrt(n)
	return false
}

// SieveOfEratosthenes returns all primes up to limit.
func SieveOfEratosthenes(limit int) []int {
	// TODO: Return empty slice if limit < 2
	// TODO: Implement boolean sieve
	return nil
}

// PrimeFactors returns the ascending prime factors of n.
// Returns ErrInvalidPrimeInput if n < 2.
func PrimeFactors(n int) ([]int, error) {
	// TODO: Validate n >= 2
	// TODO: Extract prime factors iteratively
	return nil, ErrNotImplemented
}
