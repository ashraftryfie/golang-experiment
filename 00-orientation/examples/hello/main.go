package main

import "fmt"

// Greet returns a friendly welcome string.
func Greet(name string) string {
	if name == "" {
		name = "Gopher"
	}
	return fmt.Sprintf("Hello, %s! Welcome to Go.", name)
}

func main() {
	fmt.Println(Greet("Learner"))
}
