package geometry

import (
	"errors"
	"math"
)

var (
	ErrInvalidDimensions = errors.New("dimensions must be positive")
	ErrEmptyShapes       = errors.New("shapes slice cannot be empty")
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
	if w <= 0 || h <= 0 {
		return nil, ErrInvalidDimensions
	}
	return &Rectangle{Width: w, Height: h}, nil
}

func (r Rectangle) Area() float64 {
	return r.Width * r.Height
}

func (r Rectangle) Perimeter() float64 {
	return 2 * (r.Width + r.Height)
}

type Circle struct {
	Radius float64
}

func NewCircle(radius float64) (*Circle, error) {
	if radius <= 0 {
		return nil, ErrInvalidDimensions
	}
	return &Circle{Radius: radius}, nil
}

func (c Circle) Area() float64 {
	return math.Pi * c.Radius * c.Radius
}

func (c Circle) Perimeter() float64 {
	return 2 * math.Pi * c.Radius
}

func TotalArea(shapes []Shape) float64 {
	var total float64
	for _, s := range shapes {
		total += s.Area()
	}
	return total
}

func LargestPerimeter(shapes []Shape) (Shape, error) {
	if len(shapes) == 0 {
		return nil, ErrEmptyShapes
	}

	largest := shapes[0]
	maxP := largest.Perimeter()

	for _, s := range shapes[1:] {
		p := s.Perimeter()
		if p > maxP {
			maxP = p
			largest = s
		}
	}
	return largest, nil
}
