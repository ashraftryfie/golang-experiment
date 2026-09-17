package geometry

import (
	"errors"
)

var (
	ErrInvalidDimensions = errors.New("dimensions must be positive")
	ErrEmptyShapes       = errors.New("shapes slice cannot be empty")
	ErrNotImplemented    = errors.New("TODO: implement")
)

type Shape interface {
	Area() float64
	Perimeter() float64
}

type Rectangle struct {
	Width  float64
	Height float64
}

func NewRectangle(w, h float64) (*Rectangle, error) {
	// TODO: Validate w > 0 && h > 0
	return nil, ErrNotImplemented
}

func (r Rectangle) Area() float64 {
	return 0
}

func (r Rectangle) Perimeter() float64 {
	return 0
}

type Circle struct {
	Radius float64
}

func NewCircle(radius float64) (*Circle, error) {
	// TODO: Validate radius > 0
	return nil, ErrNotImplemented
}

func (c Circle) Area() float64 {
	return 0
}

func (c Circle) Perimeter() float64 {
	return 0
}

func TotalArea(shapes []Shape) float64 {
	// TODO: Sum areas of all shapes
	return 0
}

func LargestPerimeter(shapes []Shape) (Shape, error) {
	// TODO: Find shape with largest perimeter
	return nil, ErrNotImplemented
}
