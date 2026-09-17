package fizzbuzz

import (
	"errors"
	"strconv"
	"strings"
)

var ErrInvalidRange = errors.New("n must be at least 1")
var ErrInvalidRule = errors.New("rule divisor must be positive")

type Rule struct {
	Divisor int
	Word    string
}

// ClassicFizzBuzz generates classic FizzBuzz from 1 to n.
func ClassicFizzBuzz(n int) ([]string, error) {
	if n < 1 {
		return nil, ErrInvalidRange
	}

	result := make([]string, 0, n)
	for i := 1; i <= n; i++ {
		switch {
		case i%15 == 0:
			result = append(result, "FizzBuzz")
		case i%3 == 0:
			result = append(result, "Fizz")
		case i%5 == 0:
			result = append(result, "Buzz")
		default:
			result = append(result, strconv.Itoa(i))
		}
	}
	return result, nil
}

// CustomFizzBuzz generates rule-based FizzBuzz with customizable divisors and words.
func CustomFizzBuzz(n int, rules []Rule) ([]string, error) {
	if n < 1 {
		return nil, ErrInvalidRange
	}
	for _, r := range rules {
		if r.Divisor <= 0 {
			return nil, ErrInvalidRule
		}
	}

	result := make([]string, 0, n)
	var builder strings.Builder

	for i := 1; i <= n; i++ {
		builder.Reset()
		for _, r := range rules {
			if i%r.Divisor == 0 {
				builder.WriteString(r.Word)
			}
		}
		if builder.Len() > 0 {
			result = append(result, builder.String())
		} else {
			result = append(result, strconv.Itoa(i))
		}
	}
	return result, nil
}
