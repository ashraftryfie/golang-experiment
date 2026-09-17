package ringbuffer

import (
	"testing"
)

func TestSolutionRingBuffer(t *testing.T) {
	rb, err := NewRingBuffer(2)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !rb.Push(1) || !rb.Push(2) {
		t.Fatalf("expected pushes to succeed")
	}
	if rb.Push(3) {
		t.Fatalf("expected push to full buffer to fail")
	}

	val, ok := rb.Pop()
	if !ok || val != 1 {
		t.Fatalf("got (%d, %v), want (1, true)", val, ok)
	}

	if !rb.Push(3) {
		t.Fatalf("expected push after pop to succeed")
	}

	val, ok = rb.Pop()
	if !ok || val != 2 {
		t.Fatalf("got (%d, %v), want (2, true)", val, ok)
	}

	val, ok = rb.Pop()
	if !ok || val != 3 {
		t.Fatalf("got (%d, %v), want (3, true)", val, ok)
	}

	if !rb.IsEmpty() {
		t.Fatalf("expected buffer to be empty")
	}
}
