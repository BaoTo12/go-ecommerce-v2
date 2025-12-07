package coalesce

import (
	"context"
	"sync"
	"time"
)

// Request represents a coalesced request
type Request struct {
	Key      string
	Response chan Result
}

// Result holds the result of a coalesced operation
type Result struct {
	Value interface{}
	Error error
}

// Coalescer batches identical requests together
type Coalescer struct {
	inflight map[string]*call
	mu       sync.Mutex
}

type call struct {
	wg     sync.WaitGroup
	result Result
}

// New creates a new coalescer
func New() *Coalescer {
	return &Coalescer{
		inflight: make(map[string]*call),
	}
}

// Do executes fn only once for concurrent calls with the same key
func (c *Coalescer) Do(key string, fn func() (interface{}, error)) (interface{}, error) {
	c.mu.Lock()

	// Check if call is already in flight
	if call, ok := c.inflight[key]; ok {
		c.mu.Unlock()
		call.wg.Wait()
		return call.result.Value, call.result.Error
	}

	// Create new call
	call := &call{}
	call.wg.Add(1)
	c.inflight[key] = call
	c.mu.Unlock()

	// Execute function
	val, err := fn()
	call.result = Result{Value: val, Error: err}
	call.wg.Done()

	// Clean up after a delay to allow for late arrivals
	go func() {
		time.Sleep(10 * time.Millisecond)
		c.mu.Lock()
		delete(c.inflight, key)
		c.mu.Unlock()
	}()

	return val, err
}

// DoWithContext supports context cancellation
func (c *Coalescer) DoWithContext(ctx context.Context, key string, fn func() (interface{}, error)) (interface{}, error) {
	c.mu.Lock()

	if call, ok := c.inflight[key]; ok {
		c.mu.Unlock()

		// Wait with context
		done := make(chan struct{})
		go func() {
			call.wg.Wait()
			close(done)
		}()

		select {
		case <-done:
			return call.result.Value, call.result.Error
		case <-ctx.Done():
			return nil, ctx.Err()
		}
	}

	call := &call{}
	call.wg.Add(1)
	c.inflight[key] = call
	c.mu.Unlock()

	val, err := fn()
	call.result = Result{Value: val, Error: err}
	call.wg.Done()

	go func() {
		time.Sleep(10 * time.Millisecond)
		c.mu.Lock()
		delete(c.inflight, key)
		c.mu.Unlock()
	}()

	return val, err
}

// BatchCoalescer batches multiple requests together
type BatchCoalescer struct {
	batchSize   int
	batchDelay  time.Duration
	pending     map[string][]chan Result
	batchFn     func(keys []string) map[string]Result
	mu          sync.Mutex
	timer       *time.Timer
}

// NewBatchCoalescer creates a batch coalescer
func NewBatchCoalescer(batchSize int, batchDelay time.Duration, batchFn func(keys []string) map[string]Result) *BatchCoalescer {
	return &BatchCoalescer{
		batchSize:  batchSize,
		batchDelay: batchDelay,
		pending:    make(map[string][]chan Result),
		batchFn:    batchFn,
	}
}

// Get adds a key to the batch and returns a channel for the result
func (bc *BatchCoalescer) Get(key string) <-chan Result {
	bc.mu.Lock()
	defer bc.mu.Unlock()

	ch := make(chan Result, 1)
	bc.pending[key] = append(bc.pending[key], ch)

	// Check if batch is full
	totalWaiting := 0
	for _, chs := range bc.pending {
		totalWaiting += len(chs)
	}

	if totalWaiting >= bc.batchSize {
		bc.flush()
		return ch
	}

	// Start timer if not already running
	if bc.timer == nil && len(bc.pending) == 1 {
		bc.timer = time.AfterFunc(bc.batchDelay, func() {
			bc.mu.Lock()
			defer bc.mu.Unlock()
			bc.flush()
		})
	}

	return ch
}

func (bc *BatchCoalescer) flush() {
	if bc.timer != nil {
		bc.timer.Stop()
		bc.timer = nil
	}

	keys := make([]string, 0, len(bc.pending))
	for key := range bc.pending {
		keys = append(keys, key)
	}

	if len(keys) == 0 {
		return
	}

	pending := bc.pending
	bc.pending = make(map[string][]chan Result)

	// Execute batch function
	go func() {
		results := bc.batchFn(keys)
		for key, channels := range pending {
			result, ok := results[key]
			if !ok {
				result = Result{Error: ErrKeyNotFound}
			}
			for _, ch := range channels {
				ch <- result
				close(ch)
			}
		}
	}()
}

// Deduplicator prevents duplicate processing
type Deduplicator struct {
	seen map[string]time.Time
	ttl  time.Duration
	mu   sync.RWMutex
}

// NewDeduplicator creates a deduplicator
func NewDeduplicator(ttl time.Duration) *Deduplicator {
	d := &Deduplicator{
		seen: make(map[string]time.Time),
		ttl:  ttl,
	}
	go d.cleanup()
	return d
}

// IsDuplicate returns true if key was seen recently
func (d *Deduplicator) IsDuplicate(key string) bool {
	d.mu.Lock()
	defer d.mu.Unlock()

	if lastSeen, ok := d.seen[key]; ok {
		if time.Since(lastSeen) < d.ttl {
			return true
		}
	}

	d.seen[key] = time.Now()
	return false
}

func (d *Deduplicator) cleanup() {
	ticker := time.NewTicker(d.ttl)
	defer ticker.Stop()

	for range ticker.C {
		d.mu.Lock()
		cutoff := time.Now().Add(-d.ttl)
		for key, lastSeen := range d.seen {
			if lastSeen.Before(cutoff) {
				delete(d.seen, key)
			}
		}
		d.mu.Unlock()
	}
}

// Errors
var (
	ErrKeyNotFound = &CoalesceError{Message: "key not found in batch results"}
)

type CoalesceError struct {
	Message string
}

func (e *CoalesceError) Error() string {
	return e.Message
}
