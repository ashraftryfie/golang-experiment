package geometry

import (
	"testing"
)

func TestSolutionGeometry(t *testing.T) {
	r, _ := NewRectangle(10, 5)
	c, _ := NewCircle(7)

	shapes := []Shape{r, c}
	total := TotalArea(shapes)
	if total <= 0 {
		t.Fatalf("total area should be positive")
	}

	largest, err := LargestPerimeter(shapes)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Rect P = 2*(10+5) = 30. Circle P = 2*pi*7 = ~43.98
	if _, isCircle := largest.(*Circle); !isCircle {
		t.Errorf("expected circle to have largest perimeter")
	}
}
