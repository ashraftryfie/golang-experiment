package eventdispatcher

import (
	"context"
	"errors"
	"sync"
	"sync/atomic"
)

var (
	ErrBufferFull = errors.New("event buffer is full")
	ErrClosed     = errors.New("dispatcher is closed")
)

type Event struct {
	ID      string
	Topic   string
	Payload any
}

type EventResult struct {
	Event     Event
	Processed bool
	Err       error
}

type EventHandler func(ctx context.Context, e Event) error

type Dispatcher struct {
	workers       int
	queue         chan Event
	results       chan EventResult
	handler       EventHandler
	mu            sync.Mutex
	closed        bool
	wg            sync.WaitGroup
	eventsIn      uint64
	eventsHandled uint64
}

func NewDispatcher(workers, bufferSize int, handler EventHandler) *Dispatcher {
	if workers <= 0 {
		workers = 4
	}
	if bufferSize <= 0 {
		bufferSize = 100
	}
	return &Dispatcher{
		workers: workers,
		queue:   make(chan Event, bufferSize),
		results: make(chan EventResult, bufferSize),
		handler: handler,
	}
}

func (d *Dispatcher) Start(ctx context.Context) {
	for i := 0; i < d.workers; i++ {
		d.wg.Add(1)
		go d.worker(ctx)
	}
}

func (d *Dispatcher) Dispatch(e Event) error {
	d.mu.Lock()
	if d.closed {
		d.mu.Unlock()
		return ErrClosed
	}
	d.mu.Unlock()

	select {
	case d.queue <- e:
		atomic.AddUint64(&d.eventsIn, 1)
		return nil
	default:
		return ErrBufferFull
	}
}

func (d *Dispatcher) Results() <-chan EventResult {
	return d.results
}

func (d *Dispatcher) worker(ctx context.Context) {
	defer d.wg.Done()

	for {
		select {
		case <-ctx.Done():
			return
		case event, ok := <-d.queue:
			if !ok {
				return
			}
			err := d.handler(ctx, event)
			atomic.AddUint64(&d.eventsHandled, 1)

			select {
			case <-ctx.Done():
				return
			case d.results <- EventResult{Event: event, Processed: err == nil, Err: err}:
			}
		}
	}
}

func (d *Dispatcher) Stop() {
	d.mu.Lock()
	if !d.closed {
		d.closed = true
		close(d.queue)
	}
	d.mu.Unlock()

	d.wg.Wait()
	close(d.results)
}

func (d *Dispatcher) Stats() (uint64, uint64) {
	return atomic.LoadUint64(&d.eventsIn), atomic.LoadUint64(&d.eventsHandled)
}
