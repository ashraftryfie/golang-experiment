package ringbuffer

import (
	"errors"
)

var ErrInvalidCapacity = errors.New("capacity must be greater than zero")

type RingBuffer struct {
	data     []int
	capacity int
	head     int
	tail     int
	size     int
}

// NewRingBuffer allocates a ring buffer of fixed capacity.
func NewRingBuffer(capacity int) (*RingBuffer, error) {
	if capacity <= 0 {
		return nil, ErrInvalidCapacity
	}
	// TODO: Initialize RingBuffer struct with pre-allocated slice
	return &RingBuffer{
		data:     make([]int, capacity),
		capacity: capacity,
	}, nil
}

func (r *RingBuffer) Push(val int) bool {
	// TODO: If full, return false
	// TODO: Write val at tail, advance tail circularly, increment size
	return false
}

func (r *RingBuffer) Pop() (int, bool) {
	// TODO: If empty, return 0, false
	// TODO: Read val at head, advance head circularly, decrement size
	return 0, false
}

func (r *RingBuffer) Peek() (int, bool) {
	// TODO: If empty, return 0, false
	// TODO: Return value at head without advancing
	return 0, false
}

func (r *RingBuffer) Size() int {
	return r.size
}

func (r *RingBuffer) Capacity() int {
	return r.capacity
}

func (r *RingBuffer) IsFull() bool {
	return r.size == r.capacity
}

func (r *RingBuffer) IsEmpty() bool {
	return r.size == 0
}
