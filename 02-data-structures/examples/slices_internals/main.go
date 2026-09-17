package main

import "fmt"

func inspectSlice(name string, s []int) {
	fmt.Printf("%-10s: len=%d, cap=%d, elements=%v\n", name, len(s), cap(s), s)
}

func main() {
	fmt.Println("=== 1. Slice Length vs. Capacity ===")
	// make([]int, len, cap)
	s := make([]int, 3, 5)
	s[0], s[1], s[2] = 10, 20, 30
	inspectSlice("initial s", s)

	fmt.Println("\n=== 2. Append within Existing Capacity ===")
	s = append(s, 40)
	inspectSlice("appended 40", s)

	fmt.Println("\n=== 3. Sub-slicing Shared Backing Array ===")
	sub := s[1:3] // Elements at index 1 and 2: [20, 30]
	inspectSlice("sub [1:3]", sub)

	// Mutating sub mutates s!
	sub[0] = 999
	fmt.Println("\nAfter mutating sub[0] = 999:")
	inspectSlice("sub", sub)
	inspectSlice("original s", s)

	fmt.Println("\n=== 4. Append Exceeding Capacity (Re-allocation) ===")
	s = append(s, 50) // len is now 5 (cap was 5)
	inspectSlice("s full cap", s)

	s = append(s, 60) // cap exceeded! Backing array re-allocated!
	inspectSlice("s after growth", s)

	// Now mutating s will NOT affect sub, because s points to a new backing array
	s[1] = 111
	fmt.Println("\nAfter mutating s[1] = 111 on new backing array:")
	inspectSlice("s", s)
	inspectSlice("sub (unchanged)", sub)

	fmt.Println("\n=== 5. Safe Copying ===")
	original := []int{1, 2, 3, 4, 5}
	cloned := make([]int, len(original))
	copy(cloned, original)
	cloned[0] = 42
	fmt.Printf("Original: %v\n", original)
	fmt.Printf("Cloned:   %v\n", cloned)
}
