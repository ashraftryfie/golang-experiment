package sliceops

import (
	"errors"
)

var ErrNotImplemented = errors.New("TODO: implement function")
var ErrIndexOutOfRange = errors.New("index out of range")

// Deduplicate returns a slice with duplicates removed in original order.
func Deduplicate(nums []int) []int {
	// TODO: Implement using a map for O(N) lookup
	return nil
}

// FilterEven filters nums in-place, returning only even numbers.
func FilterEven(nums []int) []int {
	// TODO: Implement in-place using nums[:0]
	return nil
}

// RemoveAtIndex removes element at index, preserving order.
func RemoveAtIndex(nums []int, index int) ([]int, error) {
	// TODO: Validate index
	// TODO: Use sub-slicing to excise element at index
	return nil, ErrNotImplemented
}
