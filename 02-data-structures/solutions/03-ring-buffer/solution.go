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
	return &RingBuffer{
		data:     make([]int, capacity),
		capacity: capacity,
		head:     0,
		tail:     0,
		size:     0,
	}, nil
}

func (r *RingBuffer) Push(val int) bool {
	if r.IsFull() {
		return false
	}
	r.data[r.tail] = val
	r.tail = (r.tail + 1) % r.capacity
	r.size++
	return true
}

func (r *RingBuffer) Pop() (int, bool) {
	if r.IsEmpty() {
		return 0, false
	}
	val := r.data[r.head]
	r.head = (r.head + 1) % r.capacity
	r.size--
	return val, true
}

func (r *RingBuffer) Peek() (int, bool) {
	if r.IsEmpty() {
		return 0, false
	}
	return r.data[r.head], true
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
