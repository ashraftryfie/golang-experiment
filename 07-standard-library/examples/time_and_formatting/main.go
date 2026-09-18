package main

import (
	"fmt"
	"time"
)

func main() {
	now := time.Now()

	fmt.Println("=== 1. Reference Date Mnemonic ===")
	// Remember: 1 2 3 4 5 6 (Mon Jan 2 15:04:05 2006)
	fmt.Println("ISO Date:      ", now.Format("2006-01-02"))
	fmt.Println("Readable Date: ", now.Format("Monday, 02 Jan 2006"))
	fmt.Println("RFC3339:       ", now.Format(time.RFC3339))
	fmt.Println("Timestamp Log: ", now.Format("2006/01/02 15:04:05.000"))

	fmt.Println("\n=== 2. Duration Arithmetic ===")
	future := now.Add(2*time.Hour + 30*time.Minute)
	diff := future.Sub(now)
	fmt.Printf("Difference: %v (minutes: %.0f)\n", diff, diff.Minutes())

	fmt.Println("\n=== 3. Parsing Time ===")
	parsed, err := time.Parse("2006-01-02", "2026-09-18")
	if err != nil {
		panic(err)
	}
	fmt.Printf("Parsed successfully: %s (weekday: %s)\n", parsed.Format("2006-01-02"), parsed.Weekday())
}
