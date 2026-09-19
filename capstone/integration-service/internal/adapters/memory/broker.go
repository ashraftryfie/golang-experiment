package memory

import (
	"context"
	"errors"
	"sync"

	"github.com/ashraftryfie/golang-experiment/capstone/integration-service/internal/domain"
)

var ErrQueueClosed = errors.New("queue broker is closed")

type ChannelJobQueue struct {
	ch     chan *domain.Job
	mu     sync.Mutex
	closed bool
}

func NewChannelJobQueue(bufferSize int) *ChannelJobQueue {
	if bufferSize <= 0 {
		bufferSize = 100
	}
	return &ChannelJobQueue{
		ch: make(chan *domain.Job, bufferSize),
	}
}

func (q *ChannelJobQueue) Enqueue(ctx context.Context, job *domain.Job) error {
	q.mu.Lock()
	if q.closed {
		q.mu.Unlock()
		return ErrQueueClosed
	}
	q.mu.Unlock()

	select {
	case <-ctx.Done():
		return ctx.Err()
	case q.ch <- cloneJob(job):
		return nil
	}
}

func (q *ChannelJobQueue) Dequeue(ctx context.Context) (*domain.Job, error) {
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	case job, ok := <-q.ch:
		if !ok {
			return nil, ErrQueueClosed
		}
		return job, nil
	}
}

func (q *ChannelJobQueue) Close() error {
	q.mu.Lock()
	defer q.mu.Unlock()

	if !q.closed {
		q.closed = true
		close(q.ch)
	}
	return nil
}
