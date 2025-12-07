package pool

import (
	"bytes"
	"sync"
)

// ObjectPool is a generic object pool using sync.Pool
type ObjectPool[T any] struct {
	pool *sync.Pool
	new  func() T
}

// NewObjectPool creates a new object pool
func NewObjectPool[T any](factory func() T) *ObjectPool[T] {
	return &ObjectPool[T]{
		pool: &sync.Pool{
			New: func() interface{} {
				return factory()
			},
		},
		new: factory,
	}
}

// Get retrieves an object from the pool
func (p *ObjectPool[T]) Get() T {
	return p.pool.Get().(T)
}

// Put returns an object to the pool
func (p *ObjectPool[T]) Put(obj T) {
	p.pool.Put(obj)
}

// BufferPool is a specialized pool for bytes.Buffer
type BufferPool struct {
	pool *sync.Pool
}

// NewBufferPool creates a new buffer pool
func NewBufferPool(initialSize int) *BufferPool {
	return &BufferPool{
		pool: &sync.Pool{
			New: func() interface{} {
				return bytes.NewBuffer(make([]byte, 0, initialSize))
			},
		},
	}
}

// Get retrieves a buffer from the pool
func (p *BufferPool) Get() *bytes.Buffer {
	buf := p.pool.Get().(*bytes.Buffer)
	buf.Reset()
	return buf
}

// Put returns a buffer to the pool
func (p *BufferPool) Put(buf *bytes.Buffer) {
	// Don't return buffers that have grown too large
	if buf.Cap() > 64*1024 { // 64KB limit
		return
	}
	p.pool.Put(buf)
}

// SlicePool pools slices of a specific type
type SlicePool[T any] struct {
	pool *sync.Pool
	cap  int
}

// NewSlicePool creates a pool for slices
func NewSlicePool[T any](capacity int) *SlicePool[T] {
	return &SlicePool[T]{
		pool: &sync.Pool{
			New: func() interface{} {
				s := make([]T, 0, capacity)
				return &s
			},
		},
		cap: capacity,
	}
}

// Get retrieves a slice from the pool
func (p *SlicePool[T]) Get() []T {
	s := p.pool.Get().(*[]T)
	return (*s)[:0]
}

// Put returns a slice to the pool
func (p *SlicePool[T]) Put(s []T) {
	// Only return slices that haven't grown too much
	if cap(s) <= p.cap*2 {
		sp := s[:0]
		p.pool.Put(&sp)
	}
}

// MapPool pools maps of specific types
type MapPool[K comparable, V any] struct {
	pool *sync.Pool
}

// NewMapPool creates a pool for maps
func NewMapPool[K comparable, V any]() *MapPool[K, V] {
	return &MapPool[K, V]{
		pool: &sync.Pool{
			New: func() interface{} {
				return make(map[K]V)
			},
		},
	}
}

// Get retrieves a map from the pool
func (p *MapPool[K, V]) Get() map[K]V {
	return p.pool.Get().(map[K]V)
}

// Put clears and returns a map to the pool
func (p *MapPool[K, V]) Put(m map[K]V) {
	// Clear the map
	for k := range m {
		delete(m, k)
	}
	p.pool.Put(m)
}

// Global pools for common types
var (
	// ByteBufferPool is a global pool for byte buffers
	ByteBufferPool = NewBufferPool(4096)

	// StringSlicePool is a global pool for string slices
	StringSlicePool = NewSlicePool[string](16)
)
