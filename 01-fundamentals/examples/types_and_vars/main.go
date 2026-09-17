package main

import "fmt"

func main() {
	fmt.Println("=== 1. Zero Values ===")
	var defaultInt int
	var defaultFloat float64
	var defaultBool bool
	var defaultString string
	var defaultPointer *int

	fmt.Printf("int: %d\n", defaultInt)
	fmt.Printf("float64: %f\n", defaultFloat)
	fmt.Printf("bool: %t\n", defaultBool)
	fmt.Printf("string: %q\n", defaultString)
	fmt.Printf("pointer: %v\n", defaultPointer)

	fmt.Println("\n=== 2. Type Inference & Short Declaration ===")
	language := "Go"
	version := 1.22
	isAwesome := true
	fmt.Printf("Language: %s (type %T)\n", language, language)
	fmt.Printf("Version: %.2f (type %T)\n", version, version)
	fmt.Printf("Awesome: %t (type %T)\n", isAwesome, isAwesome)

	fmt.Println("\n=== 3. Explicit Type Conversions ===")
	var a int = 42
	var b float64 = 3.14
	// In Go, mixed-type arithmetic without explicit conversion is a compile error
	sum := float64(a) + b
	fmt.Printf("Sum of %d (converted to float) and %.2f = %.2f\n", a, b, sum)

	fmt.Println("\n=== 4. Untyped vs Typed Constants ===")
	const UntypedPi = 3.14159265358979323846 // untyped floating-point constant
	var floatVal float32 = UntypedPi         // adapts to float32
	var doubleVal float64 = UntypedPi        // adapts to float64
	fmt.Printf("float32: %f, float64: %f\n", floatVal, doubleVal)
}
