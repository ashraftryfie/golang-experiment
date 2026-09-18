package dedup_solution

// Deduplicate removes duplicate elements from nums, maintaining order of first appearance.
func Deduplicate(nums []int) []int {
	if len(nums) == 0 {
		return []int{}
	}

	seen := make(map[int]struct{}, len(nums))
	result := make([]int, 0, len(nums))

	for _, n := range nums {
		if _, ok := seen[n]; !ok {
			seen[n] = struct{}{}
			result = append(result, n)
		}
	}

	return result
}
