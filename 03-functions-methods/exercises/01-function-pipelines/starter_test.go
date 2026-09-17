package pipelines

import (
	"reflect"
	"testing"
)

func TestMap(t *testing.T) {
	nums := []int{1, 2, 3, 4}
	double := func(n int) int { return n * 2 }
	got := Map(nums, double)
	if got == nil {
		t.Skip("skipping: Map is not yet implemented (implement in starter.go)")
	}

	want := []int{2, 4, 6, 8}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("Map() = %v, want %v", got, want)
	}
}

func TestFilter(t *testing.T) {
	nums := []int{1, 2, 3, 4, 5, 6}
	isEven := func(n int) bool { return n%2 == 0 }
	got := Filter(nums, isEven)
	if got == nil {
		t.Skip("skipping: Filter is not yet implemented (implement in starter.go)")
	}

	want := []int{2, 4, 6}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("Filter() = %v, want %v", got, want)
	}
}

func TestReduce(t *testing.T) {
	nums := []int{1, 2, 3, 4}
	sum := func(acc, cur int) int { return acc + cur }
	got := Reduce(nums, 0, sum)
	// If uncompleted returns 0 (which happens to be sum if empty, but for 1..4 sum is 10)
	if got == 0 {
		t.Skip("skipping: Reduce is not yet implemented (implement in starter.go)")
	}

	if got != 10 {
		t.Errorf("Reduce() = %v, want 10", got)
	}
}

func TestCompose(t *testing.T) {
	add1 := func(n int) int { return n + 1 }
	mult2 := func(n int) int { return n * 2 }
	comp := Compose(mult2, add1) // mult2(add1(x)) = (x + 1) * 2
	if comp == nil {
		t.Skip("skipping: Compose is not yet implemented (implement in starter.go)")
	}

	got := comp(3) // (3 + 1) * 2 = 8
	if got != 8 {
		t.Errorf("comp(3) = %v, want 8", got)
	}
}
