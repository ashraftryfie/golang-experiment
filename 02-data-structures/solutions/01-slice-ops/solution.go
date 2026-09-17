package sliceops

import (
	"errors"
)

var ErrIndexOutOfRange = errors.New("index out of range")

// Deduplicate returns a slice with duplicates removed in original order.
func Deduplicate(nums []int) []int {
	if len(nums) == 0 {
		return []int{}
	}

	seen := make(map[int]struct{}, len(nums))
	result := make([]int, 0, len(nums))

	for _, n := range nums {
		if _, exists := seen[n]; !exists {
			seen[n] = struct{}{}
			result = append(result, n)
		}
	}
	return result
}

// FilterEven filters nums in-place, returning only even numbers.
func FilterEven(nums []int) []int {
	out := nums[:0]
	for _, n := range nums {
		if n%2 == 0 {
			out = append(out, n)
		}
	}
	return out
}

// RemoveAtIndex removes element at index, preserving order.
func RemoveAtIndex(nums []int, index int) ([]int, error) {
	if index < 0 || index >= len(nums) {
		return nil, ErrIndexOutOfRange
	}
	return append(nums[:index], nums[index+1:]...), nil
}
