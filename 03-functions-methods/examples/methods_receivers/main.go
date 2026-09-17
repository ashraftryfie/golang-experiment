package main

import "fmt"

type Point struct {
	X float64
	Y float64
}

// DistanceFromOrigin is a value receiver: operates on a read-only copy of Point.
func (p Point) DistanceSquared() float64 {
	return (p.X * p.X) + (p.Y * p.Y)
}

// Scale is a pointer receiver: mutates the actual caller struct in-place.
func (p *Point) Scale(factor float64) {
	p.X *= factor
	p.Y *= factor
}

func main() {
	p := Point{X: 3, Y: 4}
	fmt.Printf("Initial Point: (%v, %v)\n", p.X, p.Y)
	fmt.Printf("Distance Squared: %v\n", p.DistanceSquared())

	// Go automatically takes the address of p (&p) to call the pointer receiver
	p.Scale(2)
	fmt.Printf("After p.Scale(2): (%v, %v)\n", p.X, p.Y)
	fmt.Printf("New Distance Squared: %v\n", p.DistanceSquared())
}
