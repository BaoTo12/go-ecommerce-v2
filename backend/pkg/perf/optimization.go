package perf

import (
	"runtime"
	"sync"
	"sync/atomic"
	"time"
	"unsafe"
)

// Arena is a memory arena for reducing GC pressure
// Uses bump allocation for fast allocations
type Arena struct {
	blocks    [][]byte
	current   []byte
	offset    int
	blockSize int
	mu        sync.Mutex
}

// NewArena creates a new memory arena
func NewArena(blockSize int) *Arena {
	if blockSize <= 0 {
		blockSize = 64 * 1024 // 64KB default
	}
	a := &Arena{
		blockSize: blockSize,
		blocks:    make([][]byte, 0, 4),
	}
	a.grow()
	return a
}

func (a *Arena) grow() {
	block := make([]byte, a.blockSize)
	a.blocks = append(a.blocks, block)
	a.current = block
	a.offset = 0
}

// Alloc allocates n bytes from the arena
func (a *Arena) Alloc(n int) []byte {
	a.mu.Lock()
	defer a.mu.Unlock()

	// Align to 8 bytes
	n = (n + 7) &^ 7

	if a.offset+n > len(a.current) {
		if n > a.blockSize {
			// Large allocation: create dedicated block
			block := make([]byte, n)
			a.blocks = append(a.blocks, block)
			return block
		}
		a.grow()
	}

	ptr := a.current[a.offset : a.offset+n]
	a.offset += n
	return ptr
}

// AllocString allocates a string in the arena
func (a *Arena) AllocString(s string) string {
	b := a.Alloc(len(s))
	copy(b, s)
	return *(*string)(unsafe.Pointer(&b))
}

// Reset resets the arena for reuse
func (a *Arena) Reset() {
	a.mu.Lock()
	defer a.mu.Unlock()
	if len(a.blocks) > 0 {
		a.current = a.blocks[0]
		a.offset = 0
		// Keep first block, release others
		a.blocks = a.blocks[:1]
	}
}

// Stats returns arena statistics
func (a *Arena) Stats() ArenaStats {
	a.mu.Lock()
	defer a.mu.Unlock()
	totalSize := 0
	for _, b := range a.blocks {
		totalSize += len(b)
	}
	return ArenaStats{
		Blocks:    len(a.blocks),
		TotalSize: totalSize,
		Used:      a.offset + (len(a.blocks)-1)*a.blockSize,
	}
}

// ArenaStats holds arena statistics
type ArenaStats struct {
	Blocks    int
	TotalSize int
	Used      int
}

// StringInterner interns strings to reduce memory
type StringInterner struct {
	strings sync.Map
	arena   *Arena
}

// NewStringInterner creates a string interner
func NewStringInterner() *StringInterner {
	return &StringInterner{
		arena: NewArena(16 * 1024),
	}
}

// Intern returns an interned copy of the string
func (si *StringInterner) Intern(s string) string {
	if existing, ok := si.strings.Load(s); ok {
		return existing.(string)
	}
	interned := si.arena.AllocString(s)
	actual, _ := si.strings.LoadOrStore(s, interned)
	return actual.(string)
}

// CPUOptimizer provides CPU-bound optimizations
type CPUOptimizer struct {
	numCPU      int
	parallelism int
}

// NewCPUOptimizer creates a CPU optimizer
func NewCPUOptimizer() *CPUOptimizer {
	return &CPUOptimizer{
		numCPU:      runtime.NumCPU(),
		parallelism: runtime.GOMAXPROCS(0),
	}
}

// ParallelFor executes fn for each item in parallel
func (c *CPUOptimizer) ParallelFor(n int, fn func(start, end int)) {
	if n <= 0 {
		return
	}

	workers := c.parallelism
	if n < workers {
		workers = n
	}

	chunkSize := (n + workers - 1) / workers
	var wg sync.WaitGroup

	for i := 0; i < workers; i++ {
		start := i * chunkSize
		end := start + chunkSize
		if end > n {
			end = n
		}
		if start >= end {
			continue
		}

		wg.Add(1)
		go func(s, e int) {
			defer wg.Done()
			fn(s, e)
		}(start, end)
	}

	wg.Wait()
}

// BatchProcessor processes items in batches
type BatchProcessor[T any] struct {
	batchSize int
	processor func([]T) error
	buffer    []T
	mu        sync.Mutex
	flushCh   chan struct{}
	done      chan struct{}
}

// NewBatchProcessor creates a batch processor
func NewBatchProcessor[T any](batchSize int, flushInterval time.Duration, processor func([]T) error) *BatchProcessor[T] {
	bp := &BatchProcessor[T]{
		batchSize: batchSize,
		processor: processor,
		buffer:    make([]T, 0, batchSize),
		flushCh:   make(chan struct{}, 1),
		done:      make(chan struct{}),
	}

	go bp.flushLoop(flushInterval)
	return bp
}

// Add adds an item to the batch
func (bp *BatchProcessor[T]) Add(item T) error {
	bp.mu.Lock()
	bp.buffer = append(bp.buffer, item)
	shouldFlush := len(bp.buffer) >= bp.batchSize
	bp.mu.Unlock()

	if shouldFlush {
		return bp.Flush()
	}
	return nil
}

// Flush flushes the current batch
func (bp *BatchProcessor[T]) Flush() error {
	bp.mu.Lock()
	if len(bp.buffer) == 0 {
		bp.mu.Unlock()
		return nil
	}
	batch := bp.buffer
	bp.buffer = make([]T, 0, bp.batchSize)
	bp.mu.Unlock()

	return bp.processor(batch)
}

func (bp *BatchProcessor[T]) flushLoop(interval time.Duration) {
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			bp.Flush()
		case <-bp.flushCh:
			bp.Flush()
		case <-bp.done:
			return
		}
	}
}

// Close closes the batch processor
func (bp *BatchProcessor[T]) Close() error {
	close(bp.done)
	return bp.Flush()
}

// AtomicCounter is a lock-free counter
type AtomicCounter struct {
	value int64
	_     [56]byte // padding to avoid false sharing
}

// NewAtomicCounter creates an atomic counter
func NewAtomicCounter() *AtomicCounter {
	return &AtomicCounter{}
}

// Inc increments the counter
func (c *AtomicCounter) Inc() int64 {
	return atomic.AddInt64(&c.value, 1)
}

// Dec decrements the counter
func (c *AtomicCounter) Dec() int64 {
	return atomic.AddInt64(&c.value, -1)
}

// Add adds n to the counter
func (c *AtomicCounter) Add(n int64) int64 {
	return atomic.AddInt64(&c.value, n)
}

// Value returns the current value
func (c *AtomicCounter) Value() int64 {
	return atomic.LoadInt64(&c.value)
}

// Reset resets the counter
func (c *AtomicCounter) Reset() {
	atomic.StoreInt64(&c.value, 0)
}
