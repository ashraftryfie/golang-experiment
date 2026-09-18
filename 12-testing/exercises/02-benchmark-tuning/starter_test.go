package dedup

import (
	"reflect"
	"testing"
)

func TestDeduplicate(t *testing.T) {
	input := []int{4, 2, 4, 3, 2, 1, 5, 1}
	expected := []int{4, 2, 3, 1, 5}

	got := Deduplicate(input)
	if got == nil {
		t.Skip("skipping: Deduplicate is not yet implemented (implement in starter.go)")
	}

	if !reflect.DeepEqual(got, expected) {
		t.Errorf("got %v, want %v", got, expected)
	}
}

func BenchmarkDeduplicate(b *testing.B) {
	data := make([]int, 1000)
	for i := 0; i < len(data); i++ {
		data[i] = i % 50 // frequent duplicates
	}

	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_ = Deduplicate(data)
	}
}
