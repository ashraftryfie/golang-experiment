package main

import "fmt"

func classifyNumber(n int) string {
	// If with initializer: 'remainder' is scoped exclusively to this if-else block
	if remainder := n % 2; remainder == 0 {
		return "even"
	} else {
		return "odd"
	}
}

func gradeScore(score int) string {
	// Expressionless switch: clean replacement for nested if-else ladders
	switch {
	case score >= 90:
		return "A"
	case score >= 80:
		return "B"
	case score >= 70:
		return "C"
	case score >= 60:
		return "D"
	default:
		return "F"
	}
}

func main() {
	fmt.Println("=== 1. If with Initializer ===")
	fmt.Printf("4 is %s\n", classifyNumber(4))
	fmt.Printf("7 is %s\n", classifyNumber(7))

	fmt.Println("\n=== 2. Switch Patterns ===")
	scores := []int{95, 82, 73, 58}
	for _, s := range scores {
		fmt.Printf("Score %d -> Grade %s\n", s, gradeScore(s))
	}

	fmt.Println("\n=== 3. For Loop Variants ===")
	// Variant A: Standard 3-component loop
	total := 0
	for i := 1; i <= 5; i++ {
		total += i
	}
	fmt.Println("Sum 1..5:", total)

	// Variant B: Condition-only loop (Go's while loop)
	countdown := 3
	for countdown > 0 {
		fmt.Printf("Countdown: %d\n", countdown)
		countdown--
	}
	fmt.Println("Blastoff!")
}
