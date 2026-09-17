package main

import (
	"fmt"
	"time"
)

// MakeCounter returns a closure encapsulating a private 'count' variable.
func MakeCounter(initial int) func() int {
	count := initial
	return func() int {
		count++
		return count
	}
}

// TrackExecution demonstrates using defer for timing.
func TrackExecution(task string) func() {
	start := time.Now()
	fmt.Printf("[%s] Starting...\n", task)
	return func() {
		fmt.Printf("[%s] Completed in %v\n", task, time.Since(start))
	}
}

func deferOrderExample() {
	fmt.Println("\n=== Defer LIFO Execution Order ===")
	defer fmt.Println("Defer 1 (registered first -> executes last)")
	defer fmt.Println("Defer 2 (registered second)")
	defer fmt.Println("Defer 3 (registered third -> executes first)")
	fmt.Println("Surrounding function body finishing...")
}

func main() {
	fmt.Println("=== 1. Closures ===")
	c1 := MakeCounter(0)
	c2 := MakeCounter(100)

	fmt.Printf("c1: %d, %d, %d\n", c1(), c1(), c1())
	fmt.Printf("c2: %d, %d\n", c2(), c2())

	fmt.Println("\n=== 2. Defer Timing ===")
	done := TrackExecution("Database Query")
	time.Sleep(10 * time.Millisecond)
	done()

	deferOrderExample()
}
