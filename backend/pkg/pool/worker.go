package pool

import (
	"context"
	"sync"
	"sync/atomic"
	"time"
)

// Task represents a unit of work
type Task func(ctx context.Context) error

// WorkerPool manages a pool of worker goroutines
type WorkerPool struct {
	workers     int
	taskQueue   chan Task
	results     chan error
	wg          sync.WaitGroup
	ctx         context.Context
	cancel      context.CancelFunc
	running     int32
	processed   int64
	errors      int64
}

// WorkerPoolConfig configures the worker pool
type WorkerPoolConfig struct {
	Workers      int
	QueueSize    int
	GracefulStop time.Duration
}

// DefaultWorkerConfig returns sensible defaults
func DefaultWorkerConfig() WorkerPoolConfig {
	return WorkerPoolConfig{
		Workers:      10,
		QueueSize:    100,
		GracefulStop: 30 * time.Second,
	}
}

// NewWorkerPool creates a new worker pool
func NewWorkerPool(cfg WorkerPoolConfig) *WorkerPool {
	ctx, cancel := context.WithCancel(context.Background())

	pool := &WorkerPool{
		workers:   cfg.Workers,
		taskQueue: make(chan Task, cfg.QueueSize),
		results:   make(chan error, cfg.QueueSize),
		ctx:       ctx,
		cancel:    cancel,
	}

	// Start workers
	for i := 0; i < cfg.Workers; i++ {
		pool.wg.Add(1)
		go pool.worker(i)
	}

	return pool
}

// worker processes tasks from the queue
func (p *WorkerPool) worker(id int) {
	defer p.wg.Done()

	for {
		select {
		case <-p.ctx.Done():
			return
		case task, ok := <-p.taskQueue:
			if !ok {
				return
			}

			atomic.AddInt32(&p.running, 1)
			err := task(p.ctx)
			atomic.AddInt32(&p.running, -1)
			atomic.AddInt64(&p.processed, 1)

			if err != nil {
				atomic.AddInt64(&p.errors, 1)
				select {
				case p.results <- err:
				default:
					// Results channel full, discard
				}
			}
		}
	}
}

// Submit adds a task to the queue
func (p *WorkerPool) Submit(task Task) error {
	select {
	case p.taskQueue <- task:
		return nil
	case <-p.ctx.Done():
		return context.Canceled
	}
}

// SubmitWait adds a task and waits for completion
func (p *WorkerPool) SubmitWait(ctx context.Context, task Task) error {
	done := make(chan error, 1)

	wrappedTask := func(c context.Context) error {
		err := task(c)
		done <- err
		return err
	}

	if err := p.Submit(wrappedTask); err != nil {
		return err
	}

	select {
	case err := <-done:
		return err
	case <-ctx.Done():
		return ctx.Err()
	}
}

// Stats returns pool statistics
func (p *WorkerPool) Stats() WorkerPoolStats {
	return WorkerPoolStats{
		Workers:      p.workers,
		Running:      int(atomic.LoadInt32(&p.running)),
		Processed:    atomic.LoadInt64(&p.processed),
		Errors:       atomic.LoadInt64(&p.errors),
		QueueLength:  len(p.taskQueue),
		QueueCapacity: cap(p.taskQueue),
	}
}

// WorkerPoolStats holds pool metrics
type WorkerPoolStats struct {
	Workers       int
	Running       int
	Processed     int64
	Errors        int64
	QueueLength   int
	QueueCapacity int
}

// Resize dynamically adjusts the number of workers
func (p *WorkerPool) Resize(newWorkers int) {
	if newWorkers <= 0 {
		return
	}

	diff := newWorkers - p.workers
	if diff == 0 {
		return
	}

	if diff > 0 {
		// Add workers
		for i := 0; i < diff; i++ {
			p.wg.Add(1)
			go p.worker(p.workers + i)
		}
	}
	// Note: Reducing workers requires more complex logic (graceful shutdown)

	p.workers = newWorkers
}

// Shutdown gracefully stops the worker pool
func (p *WorkerPool) Shutdown(timeout time.Duration) error {
	// Stop accepting new tasks
	close(p.taskQueue)

	// Wait for workers to finish with timeout
	done := make(chan struct{})
	go func() {
		p.wg.Wait()
		close(done)
	}()

	select {
	case <-done:
		return nil
	case <-time.After(timeout):
		p.cancel()
		return context.DeadlineExceeded
	}
}
