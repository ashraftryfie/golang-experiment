package main

import (
	"fmt"
	"math"
)

// Greeter interface
type Greeter interface {
	Greet() string
}

type Human struct {
	Name string
}

func (h Human) Greet() string {
	return "Hello, I am " + h.Name
}

type Robot struct {
	Model string
}

func (r Robot) Greet() string {
	return "Beep boop! I am model " + r.Model
}

// Broadcaster accepts any type satisfying Greeter
func Broadcast(g Greeter) {
	fmt.Println(g.Greet())
}

// TypeSwitchExample demonstrates runtime type inspection
func Inspect(val any) {
	switch v := val.(type) {
	case string:
		fmt.Printf("String of len %d: %q\n", len(v), v)
	case int:
		fmt.Printf("Integer: %d (squared: %d)\n", v, v*v)
	case float64:
		fmt.Printf("Float: %.2f (sqrt: %.2f)\n", v, math.Sqrt(v))
	case Greeter:
		fmt.Printf("Greeter interface: %s\n", v.Greet())
	default:
		fmt.Printf("Unknown type: %T\n", v)
	}
}

func main() {
	h := Human{Name: "Gopher"}
	r := Robot{Model: "RX-78"}

	Broadcast(h)
	Broadcast(r)

	fmt.Println("\n=== Type Switch Demonstrations ===")
	Inspect("hello world")
	Inspect(42)
	Inspect(25.0)
	Inspect(h)
}
