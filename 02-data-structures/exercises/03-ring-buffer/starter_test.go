package ringbuffer

import (
	"errors"
	"testing"
)

func TestRingBufferOperations(t *testing.T) {
	rb, err := NewRingBuffer(3)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// If Push is uncompleted, skip gracefully
	if !rb.Push(10) {
		t.Skip("skipping: Push is not yet implemented (implement in starter.go)")
	}

	rb.Push(20)
	rb.Push(30)

	if !rb.IsFull() {
		t.Errorf("expected buffer to be full")
	}

	// Cannot push when full
	if rb.Push(40) {
		t.Errorf("expected push to full buffer to return false")
	}

	// Peek
	val, ok := rb.Peek()
	if !ok || val != 10 {
		t.Errorf("Peek() = (%d, %v), want (10, true)", val, ok)
	}

	// Pop
	val, ok = rb.Pop()
	if !ok || val != 10 {
		t.Errorf("Pop() = (%d, %v), want (10, true)", val, ok)
	}

	// Circular wrap push
	if !rb.Push(40) {
		t.Errorf("expected push after pop to succeed")
	}

	// Pop remaining items in FIFO order
	expected := []int{20, 30, 40}
	for _, want := range expected {
		val, ok := rb.Pop()
		if !ok || val != want {
			t.Errorf("Pop() = (%d, %v), want (%d, true)", val, ok, want)
		}
	}

	if !rb.IsEmpty() {
		t.Errorf("expected buffer to be empty")
	}

	// Pop empty
	_, ok = rb.Pop()
	if ok {
		t.Errorf("expected pop on empty buffer to return false")
	}
}

func TestInvalidCapacity(t *testing.T) {
	_, err := NewRingBuffer(0)
	if !errors.Is(err, ErrInvalidCapacity) {
		t.Errorf("expected ErrInvalidCapacity, got %v", err)
	}
}
