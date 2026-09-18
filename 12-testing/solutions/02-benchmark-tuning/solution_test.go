package dedup_solution

import (
	"reflect"
	"testing"
)

func TestDeduplicateSolution(t *testing.T) {
	input := []int{4, 2, 4, 3, 2, 1, 5, 1}
	expected := []int{4, 2, 3, 1, 5}

	got := Deduplicate(input)
	if !reflect.DeepEqual(got, expected) {
		t.Errorf("got %v, want %v", got, expected)
	}

	// Empty slice test
	empty := Deduplicate([]int{})
	if empty == nil || len(empty) != 0 {
		t.Errorf("expected empty non-nil slice, got %v", empty)
	}
}

func BenchmarkDeduplicateSolution(b *testing.B) {
	data := make([]int, 1000)
	for i := 0; i < len(data); i++ {
		data[i] = i % 50
	}

	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_ = Deduplicate(data)
	}
}
