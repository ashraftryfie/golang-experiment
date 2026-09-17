package pipelines

// Map applies fn to each element of nums.
func Map(nums []int, fn func(int) int) []int {
	result := make([]int, len(nums))
	for i, n := range nums {
		result[i] = fn(n)
	}
	return result
}

// Filter returns elements satisfying predicate.
func Filter(nums []int, predicate func(int) bool) []int {
	result := make([]int, 0, len(nums))
	for _, n := range nums {
		if predicate(n) {
			result = append(result, n)
		}
	}
	return result
}

// Reduce combines all elements starting with initial.
func Reduce(nums []int, initial int, accumulator func(acc, current int) int) int {
	acc := initial
	for _, n := range nums {
		acc = accumulator(acc, n)
	}
	return acc
}

// Compose returns f(g(x)).
func Compose(f, g func(int) int) func(int) int {
	return func(x int) int {
		return f(g(x))
	}
}
