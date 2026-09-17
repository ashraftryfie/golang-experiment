# Exercise 01: Geometric Shapes & Polymorphism (Tier 1 - Easy)

## 🎯 Problem Statement
Implement polymorphic geometric shapes using the Go `Shape` interface:

```go
type Shape interface {
    Area() float64
    Perimeter() float64
}
```

### Requirements
1. Implement `Rectangle`:
   - Fields: `Width float64`, `Height float64`.
   - Returns error if either dimension is `<= 0`.
2. Implement `Circle`:
   - Field: `Radius float64`.
   - Returns error if radius is `<= 0`.
3. Implement `TotalArea(shapes []Shape) float64`:
   - Calculates the sum of all shapes' areas.
4. Implement `LargestPerimeter(shapes []Shape) (Shape, error)`:
   - Returns the shape with the maximum perimeter.
   - Returns error if `shapes` slice is empty.

---

## 🛠️ Instructions
1. Open [`starter.go`](./starter.go).
2. Implement shapes and helper functions.
3. Run tests:
   ```powershell
   go test -v ./04-structs-interfaces/exercises/01-shapes-geometry
   ```
