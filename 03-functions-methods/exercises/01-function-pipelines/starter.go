package pipelines

// Map applies fn to each element of nums.
func Map(nums []int, fn func(int) int) []int {
	// TODO: Allocate slice and apply fn
	return nil
}

// Filter returns elements satisfying predicate.
func Filter(nums []int, predicate func(int) bool) []int {
	// TODO: Filter elements into result slice
	return nil
}

// Reduce combines all elements starting with initial.
func Reduce(nums []int, initial int, accumulator func(acc, current int) int) int {
	// TODO: Accumulate and return final value
	return initial
}

// Compose returns f(g(x)).
func Compose(f, g func(int) int) func(int) int {
	// TODO: Return composed closure
	return nil
}
