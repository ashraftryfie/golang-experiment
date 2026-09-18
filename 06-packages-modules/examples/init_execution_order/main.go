package main

import "fmt"

var initOrder = []string{}

func track(step string) int {
	initOrder = append(initOrder, step)
	return len(initOrder)
}

// 1. Package-level variables initialize first
var firstVar = track("1. Package variable firstVar evaluated")
var secondVar = track("2. Package variable secondVar evaluated")

// 2. init() functions execute second (in source file order)
func init() {
	track("3. init() block A executed")
}

func init() {
	track("4. init() block B executed")
}

// 3. main() executes third
func main() {
	track("5. main() function entered")

	fmt.Println("=== Execution Sequence Verified ===")
	for _, step := range initOrder {
		fmt.Println("  ", step)
	}
}
