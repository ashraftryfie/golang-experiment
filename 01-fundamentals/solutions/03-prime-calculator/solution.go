package primes

import (
	"errors"
	"math"
)

var ErrInvalidPrimeInput = errors.New("input must be greater than or equal to 2")

// IsPrime returns true if n is a prime number.
func IsPrime(n int) bool {
	if n <= 1 {
		return false
	}
	if n <= 3 {
		return true
	}
	if n%2 == 0 || n%3 == 0 {
		return false
	}
	limit := int(math.Sqrt(float64(n)))
	for i := 5; i <= limit; i += 6 {
		if n%i == 0 || n%(i+2) == 0 {
			return false
		}
	}
	return true
}

// SieveOfEratosthenes returns all primes up to limit.
func SieveOfEratosthenes(limit int) []int {
	if limit < 2 {
		return []int{}
	}

	isComposite := make([]bool, limit+1)
	sqrtLimit := int(math.Sqrt(float64(limit)))

	for p := 2; p <= sqrtLimit; p++ {
		if !isComposite[p] {
			for multiple := p * p; multiple <= limit; multiple += p {
				isComposite[multiple] = true
			}
		}
	}

	primes := make([]int, 0)
	for p := 2; p <= limit; p++ {
		if !isComposite[p] {
			primes = append(primes, p)
		}
	}
	return primes
}

// PrimeFactors returns the ascending prime factors of n.
func PrimeFactors(n int) ([]int, error) {
	if n < 2 {
		return nil, ErrInvalidPrimeInput
	}

	factors := make([]int, 0)

	// Pull out 2s
	for n%2 == 0 {
		factors = append(factors, 2)
		n /= 2
	}

	// Pull out odd factors
	for d := 3; d*d <= n; d += 2 {
		for n%d == 0 {
			factors = append(factors, d)
			n /= d
		}
	}

	// If n is still > 1, remaining n is prime
	if n > 1 {
		factors = append(factors, n)
	}

	return factors, nil
}
