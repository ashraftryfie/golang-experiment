package pipelines

import (
	"reflect"
	"testing"
)

func TestSolutionPipelines(t *testing.T) {
	nums := []int{1, 2, 3, 4, 5}
	squaredEvens := Map(
		Filter(nums, func(n int) bool { return n%2 == 0 }),
		func(n int) int { return n * n },
	)

	want := []int{4, 16}
	if !reflect.DeepEqual(squaredEvens, want) {
		t.Errorf("got %v, want %v", squaredEvens, want)
	}

	sum := Reduce(squaredEvens, 0, func(a, b int) int { return a + b })
	if sum != 20 {
		t.Errorf("Reduce sum = %d, want 20", sum)
	}
}
