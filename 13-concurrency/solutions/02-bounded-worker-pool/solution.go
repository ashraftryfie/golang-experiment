package workerpool_solution

import (
	"sync"
)

type WorkerPool struct {
	tasks    chan func()
	wg       sync.WaitGroup
	mu       sync.RWMutex
	isClosed bool
	once     sync.Once
}

// NewWorkerPool starts workerCount workers listening on an internal buffered task channel.
func NewWorkerPool(workerCount, queueCapacity int) *WorkerPool {
	if workerCount <= 0 {
		workerCount = 1
	}
	if queueCapacity < 0 {
		queueCapacity = 0
	}

	p := &WorkerPool{
		tasks: make(chan func(), queueCapacity),
	}

	for i := 0; i < workerCount; i++ {
		p.wg.Add(1)
		go func() {
			defer p.wg.Done()
			for task := range p.tasks {
				task()
			}
		}()
	}

	return p
}

// Submit queues a task for execution, returning false if the pool is closed.
func (p *WorkerPool) Submit(task func()) bool {
	p.mu.RLock()
	defer p.mu.RUnlock()

	if p.isClosed {
		return false
	}

	p.tasks <- task
	return true
}

// Shutdown gracefully waits for all queued tasks to finish before returning.
func (p *WorkerPool) Shutdown() {
	p.once.Do(func() {
		p.mu.Lock()
		p.isClosed = true
		close(p.tasks)
		p.mu.Unlock()

		p.wg.Wait()
	})
}
