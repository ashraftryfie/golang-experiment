package main

import "fmt"

func main() {
	fmt.Println("=== 1. Map Initialization ===")
	// Initialized with make
	userAges := make(map[string]int)
	userAges["Alice"] = 30
	userAges["Bob"] = 25

	// Initialized with map literal
	ports := map[string]int{
		"http":  80,
		"https": 443,
		"ssh":   22,
	}
	fmt.Println("User ages:", userAges)
	fmt.Println("Ports:", ports)

	fmt.Println("\n=== 2. Comma-Ok Idiom ===")
	// Checking if key exists vs zero value
	age, exists := userAges["Charlie"]
	if !exists {
		fmt.Printf("'Charlie' is not in map (default zero-value is %d)\n", age)
	}

	userAges["Charlie"] = 0 // Explicitly set to 0
	age, exists = userAges["Charlie"]
	if exists {
		fmt.Printf("'Charlie' exists in map with value %d\n", age)
	}

	fmt.Println("\n=== 3. Deleting Keys ===")
	delete(ports, "ssh")
	fmt.Println("Ports after delete(ssh):", ports)
	// Deleting a non-existent key is a safe no-op
	delete(ports, "ftp")

	fmt.Println("\n=== 4. Nil Map Behavior ===")
	var nilMap map[string]int
	fmt.Println("Reading from nil map is safe. Value:", nilMap["any"]) // returns 0
	// nilMap["crash"] = 1 // UNCOMMENTING THIS WILL PANIC: assignment to entry in nil map!
}
