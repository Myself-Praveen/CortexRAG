package workerpool

import (
	"sync"
)

// Task represents a function to be executed by a worker.
type Task func()

// Pool represents a bounded goroutine worker pool.
type Pool struct {
	tasks chan Task
	wg    sync.WaitGroup
}

// New creates a new worker pool with the specified number of workers.
func New(numWorkers int) *Pool {
	if numWorkers <= 0 {
		numWorkers = 1
	}
	
	p := &Pool{
		tasks: make(chan Task),
	}

	for i := 0; i < numWorkers; i++ {
		p.wg.Add(1)
		go p.worker()
	}

	return p
}

// worker listens for tasks on the tasks channel and executes them.
func (p *Pool) worker() {
	defer p.wg.Done()
	for task := range p.tasks {
		task()
	}
}

// Submit adds a task to the pool. Blocks if the tasks channel is unbuffered and full (if we added buffer later).
// Currently, it blocks until a worker is available.
func (p *Pool) Submit(t Task) {
	p.tasks <- t
}

// Shutdown closes the tasks channel and waits for all workers to finish.
func (p *Pool) Shutdown() {
	close(p.tasks)
	p.wg.Wait()
}
