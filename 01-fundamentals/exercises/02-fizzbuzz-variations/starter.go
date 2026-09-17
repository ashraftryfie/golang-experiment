package fizzbuzz

import (
	"errors"
)

// ErrNotImplemented is returned by starter templates before you implement them.
var ErrNotImplemented = errors.New("TODO: implement function")

// ErrInvalidRange is returned when n is less than 1.
var ErrInvalidRange = errors.New("n must be at least 1")

// ErrInvalidRule is returned when a rule contains a non-positive divisor.
var ErrInvalidRule = errors.New("rule divisor must be positive")

// Rule defines a custom divisor and its string replacement.
type Rule struct {
	Divisor int
	Word    string
}

// ClassicFizzBuzz generates classic FizzBuzz from 1 to n.
func ClassicFizzBuzz(n int) ([]string, error) {
	// TODO: Replace with your implementation
	return nil, ErrNotImplemented
}

// CustomFizzBuzz generates rule-based FizzBuzz with customizable divisors and words.
func CustomFizzBuzz(n int, rules []Rule) ([]string, error) {
	// TODO: Replace with your implementation
	return nil, ErrNotImplemented
}
