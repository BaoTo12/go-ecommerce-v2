package lockfree

import (
	"sync/atomic"
	"unsafe"
)

/*
LOCK-FREE DATA STRUCTURES - High Performance Concurrency

Implements lock-free concurrent data structures using atomic operations.
These are significantly faster than mutex-based structures under high contention.

Includes:
- Lock-free Stack (Treiber Stack)
- Lock-free Queue (Michael-Scott Queue)
- Lock-free Ring Buffer
- Atomic Value with version (for ABA problem)
*/

// Stack is a lock-free stack (Treiber Stack algorithm)
type Stack[T any] struct {
	head unsafe.Pointer // *node[T]
	size int64
}

type node[T any] struct {
	value T
	next  unsafe.Pointer // *node[T]
}

// NewStack creates a new lock-free stack
func NewStack[T any]() *Stack[T] {
	return &Stack[T]{}
}

// Push adds an item to the stack
func (s *Stack[T]) Push(value T) {
	n := &node[T]{value: value}
	for {
		oldHead := atomic.LoadPointer(&s.head)
		n.next = oldHead
		if atomic.CompareAndSwapPointer(&s.head, oldHead, unsafe.Pointer(n)) {
			atomic.AddInt64(&s.size, 1)
			return
		}
	}
}

// Pop removes and returns the top item from the stack
func (s *Stack[T]) Pop() (T, bool) {
	for {
		oldHead := atomic.LoadPointer(&s.head)
		if oldHead == nil {
			var zero T
			return zero, false
		}
		n := (*node[T])(oldHead)
		if atomic.CompareAndSwapPointer(&s.head, oldHead, n.next) {
			atomic.AddInt64(&s.size, -1)
			return n.value, true
		}
	}
}

// Peek returns the top item without removing it
func (s *Stack[T]) Peek() (T, bool) {
	head := atomic.LoadPointer(&s.head)
	if head == nil {
		var zero T
		return zero, false
	}
	return (*node[T])(head).value, true
}

// Size returns the current size
func (s *Stack[T]) Size() int64 {
	return atomic.LoadInt64(&s.size)
}

// IsEmpty returns true if the stack is empty
func (s *Stack[T]) IsEmpty() bool {
	return atomic.LoadPointer(&s.head) == nil
}

// Queue is a lock-free queue (Michael-Scott algorithm)
type Queue[T any] struct {
	head unsafe.Pointer // *qnode[T]
	tail unsafe.Pointer // *qnode[T]
	size int64
}

type qnode[T any] struct {
	value T
	next  unsafe.Pointer // *qnode[T]
}

// NewQueue creates a new lock-free queue
func NewQueue[T any]() *Queue[T] {
	// Create dummy node
	dummy := &qnode[T]{}
	return &Queue[T]{
		head: unsafe.Pointer(dummy),
		tail: unsafe.Pointer(dummy),
	}
}

// Enqueue adds an item to the queue
func (q *Queue[T]) Enqueue(value T) {
	n := &qnode[T]{value: value}
	for {
		tail := atomic.LoadPointer(&q.tail)
		tailNode := (*qnode[T])(tail)
		next := atomic.LoadPointer(&tailNode.next)

		if tail == atomic.LoadPointer(&q.tail) {
			if next == nil {
				// Try to add new node
				if atomic.CompareAndSwapPointer(&tailNode.next, nil, unsafe.Pointer(n)) {
					// Try to move tail
					atomic.CompareAndSwapPointer(&q.tail, tail, unsafe.Pointer(n))
					atomic.AddInt64(&q.size, 1)
					return
				}
			} else {
				// Tail falling behind, try to advance it
				atomic.CompareAndSwapPointer(&q.tail, tail, next)
			}
		}
	}
}

// Dequeue removes and returns the front item
func (q *Queue[T]) Dequeue() (T, bool) {
	for {
		head := atomic.LoadPointer(&q.head)
		tail := atomic.LoadPointer(&q.tail)
		headNode := (*qnode[T])(head)
		next := atomic.LoadPointer(&headNode.next)

		if head == atomic.LoadPointer(&q.head) {
			if head == tail {
				if next == nil {
					var zero T
					return zero, false // Queue is empty
				}
				// Tail falling behind
				atomic.CompareAndSwapPointer(&q.tail, tail, next)
			} else {
				// Read value before CAS
				nextNode := (*qnode[T])(next)
				value := nextNode.value
				if atomic.CompareAndSwapPointer(&q.head, head, next) {
					atomic.AddInt64(&q.size, -1)
					return value, true
				}
			}
		}
	}
}

// Size returns the current size
func (q *Queue[T]) Size() int64 {
	return atomic.LoadInt64(&q.size)
}

// IsEmpty returns true if the queue is empty
func (q *Queue[T]) IsEmpty() bool {
	head := atomic.LoadPointer(&q.head)
	headNode := (*qnode[T])(head)
	return atomic.LoadPointer(&headNode.next) == nil
}

// RingBuffer is a lock-free ring buffer
type RingBuffer[T any] struct {
	buffer  []T
	mask    int64
	head    int64
	tail    int64
}

// NewRingBuffer creates a ring buffer with given size (must be power of 2)
func NewRingBuffer[T any](size int) *RingBuffer[T] {
	// Round up to power of 2
	n := 1
	for n < size {
		n <<= 1
	}
	
	return &RingBuffer[T]{
		buffer: make([]T, n),
		mask:   int64(n - 1),
	}
}

// TryPush attempts to add an item (non-blocking)
func (r *RingBuffer[T]) TryPush(value T) bool {
	for {
		head := atomic.LoadInt64(&r.head)
		tail := atomic.LoadInt64(&r.tail)
		
		// Check if full
		if head-tail >= int64(len(r.buffer)) {
			return false
		}
		
		if atomic.CompareAndSwapInt64(&r.head, head, head+1) {
			r.buffer[head&r.mask] = value
			return true
		}
	}
}

// TryPop attempts to remove an item (non-blocking)
func (r *RingBuffer[T]) TryPop() (T, bool) {
	for {
		head := atomic.LoadInt64(&r.head)
		tail := atomic.LoadInt64(&r.tail)
		
		// Check if empty
		if tail >= head {
			var zero T
			return zero, false
		}
		
		value := r.buffer[tail&r.mask]
		if atomic.CompareAndSwapInt64(&r.tail, tail, tail+1) {
			return value, true
		}
	}
}

// Size returns current size
func (r *RingBuffer[T]) Size() int64 {
	return atomic.LoadInt64(&r.head) - atomic.LoadInt64(&r.tail)
}

// Capacity returns the buffer capacity
func (r *RingBuffer[T]) Capacity() int {
	return len(r.buffer)
}

// AtomicValue wraps a value with version to solve ABA problem
type AtomicValue[T any] struct {
	ptr unsafe.Pointer // *versionedValue[T]
}

type versionedValue[T any] struct {
	value   T
	version uint64
}

// NewAtomicValue creates a new atomic value
func NewAtomicValue[T any](initial T) *AtomicValue[T] {
	return &AtomicValue[T]{
		ptr: unsafe.Pointer(&versionedValue[T]{value: initial, version: 0}),
	}
}

// Load returns the current value
func (a *AtomicValue[T]) Load() T {
	return (*versionedValue[T])(atomic.LoadPointer(&a.ptr)).value
}

// Store sets a new value
func (a *AtomicValue[T]) Store(value T) {
	for {
		old := atomic.LoadPointer(&a.ptr)
		oldV := (*versionedValue[T])(old)
		newV := &versionedValue[T]{
			value:   value,
			version: oldV.version + 1,
		}
		if atomic.CompareAndSwapPointer(&a.ptr, old, unsafe.Pointer(newV)) {
			return
		}
	}
}

// CompareAndSwap atomically swaps if current == old
func (a *AtomicValue[T]) CompareAndSwap(old, new T) bool {
	for {
		current := atomic.LoadPointer(&a.ptr)
		currentV := (*versionedValue[T])(current)
		
		// Compare values (need comparable type)
		if any(currentV.value) != any(old) {
			return false
		}
		
		newV := &versionedValue[T]{
			value:   new,
			version: currentV.version + 1,
		}
		if atomic.CompareAndSwapPointer(&a.ptr, current, unsafe.Pointer(newV)) {
			return true
		}
	}
}

// Counter is a lock-free counter with padding to avoid false sharing
type Counter struct {
	value int64
	_     [56]byte // Padding to cache line size
}

// NewCounter creates a new counter
func NewCounter() *Counter {
	return &Counter{}
}

// Inc increments and returns new value
func (c *Counter) Inc() int64 {
	return atomic.AddInt64(&c.value, 1)
}

// Dec decrements and returns new value
func (c *Counter) Dec() int64 {
	return atomic.AddInt64(&c.value, -1)
}

// Add adds delta and returns new value
func (c *Counter) Add(delta int64) int64 {
	return atomic.AddInt64(&c.value, delta)
}

// Value returns current value
func (c *Counter) Value() int64 {
	return atomic.LoadInt64(&c.value)
}

// Reset resets to zero
func (c *Counter) Reset() {
	atomic.StoreInt64(&c.value, 0)
}
