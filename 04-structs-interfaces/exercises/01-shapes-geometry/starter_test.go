package geometry

import (
	"errors"
	"math"
	"testing"
)

func almostEqual(a, b float64) bool {
	return math.Abs(a-b) < 1e-4
}

func TestShapes(t *testing.T) {
	rect, err := NewRectangle(4, 5)
	if errors.Is(err, ErrNotImplemented) {
		t.Skip("skipping: NewRectangle is not yet implemented (implement in starter.go)")
	}
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if rect.Area() != 20.0 {
		t.Errorf("rect area = %v, want 20.0", rect.Area())
	}
	if rect.Perimeter() != 18.0 {
		t.Errorf("rect perimeter = %v, want 18.0", rect.Perimeter())
	}

	circ, err := NewCircle(3)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	shapes := []Shape{rect, circ}
	total := TotalArea(shapes)
	expectedTotal := 20.0 + (math.Pi * 9.0)
	if !almostEqual(total, expectedTotal) {
		t.Errorf("TotalArea() = %v, want %v", total, expectedTotal)
	}
}

func TestValidationErrors(t *testing.T) {
	if _, err := NewRectangle(-1, 5); !errors.Is(err, ErrInvalidDimensions) {
		if errors.Is(err, ErrNotImplemented) {
			t.Skip("skipping: NewRectangle is not yet implemented (implement in starter.go)")
		}
		t.Errorf("expected ErrInvalidDimensions, got %v", err)
	}
	if _, err := LargestPerimeter([]Shape{}); !errors.Is(err, ErrEmptyShapes) {
		if errors.Is(err, ErrNotImplemented) {
			t.Skip("skipping: LargestPerimeter is not yet implemented (implement in starter.go)")
		}
		t.Errorf("expected ErrEmptyShapes, got %v", err)
	}
}
