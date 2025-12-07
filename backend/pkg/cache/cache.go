package cache

import (
	"context"
	"encoding/json"
	"errors"
	"sync"
	"time"
)

/*
DISTRIBUTED CACHE WITH CACHE-ASIDE PATTERN

Implements:
- Multi-level caching (L1 local + L2 distributed)
- Cache-aside pattern
- Cache stampede protection
- TTL and sliding expiration
- Cache tags for invalidation
- Warm-up and pre-loading
*/

var (
	ErrCacheMiss = errors.New("cache miss")
	ErrCacheExpired = errors.New("cache expired")
)

// Entry represents a cache entry
type Entry struct {
	Key       string
	Value     []byte
	Tags      []string
	CreatedAt time.Time
	ExpiresAt time.Time
	Sliding   time.Duration // Sliding expiration
	Hits      int64
}

// CacheStore interface for cache backends
type CacheStore interface {
	Get(ctx context.Context, key string) (*Entry, error)
	Set(ctx context.Context, entry *Entry) error
	Delete(ctx context.Context, key string) error
	DeleteByTag(ctx context.Context, tag string) error
	Clear(ctx context.Context) error
}

// MultiLevelCache implements multi-level caching
type MultiLevelCache struct {
	l1 CacheStore // Local cache
	l2 CacheStore // Distributed cache
}

// NewMultiLevelCache creates a multi-level cache
func NewMultiLevelCache(l1, l2 CacheStore) *MultiLevelCache {
	return &MultiLevelCache{l1: l1, l2: l2}
}

func (c *MultiLevelCache) Get(ctx context.Context, key string) (*Entry, error) {
	// Try L1 first
	if c.l1 != nil {
		entry, err := c.l1.Get(ctx, key)
		if err == nil && entry != nil {
			return entry, nil
		}
	}

	// Try L2
	if c.l2 != nil {
		entry, err := c.l2.Get(ctx, key)
		if err == nil && entry != nil {
			// Populate L1
			if c.l1 != nil {
				c.l1.Set(ctx, entry)
			}
			return entry, nil
		}
	}

	return nil, ErrCacheMiss
}

func (c *MultiLevelCache) Set(ctx context.Context, entry *Entry) error {
	// Set in both levels
	if c.l1 != nil {
		c.l1.Set(ctx, entry)
	}
	if c.l2 != nil {
		return c.l2.Set(ctx, entry)
	}
	return nil
}

func (c *MultiLevelCache) Delete(ctx context.Context, key string) error {
	if c.l1 != nil {
		c.l1.Delete(ctx, key)
	}
	if c.l2 != nil {
		return c.l2.Delete(ctx, key)
	}
	return nil
}

func (c *MultiLevelCache) DeleteByTag(ctx context.Context, tag string) error {
	if c.l1 != nil {
		c.l1.DeleteByTag(ctx, tag)
	}
	if c.l2 != nil {
		return c.l2.DeleteByTag(ctx, tag)
	}
	return nil
}

func (c *MultiLevelCache) Clear(ctx context.Context) error {
	if c.l1 != nil {
		c.l1.Clear(ctx)
	}
	if c.l2 != nil {
		return c.l2.Clear(ctx)
	}
	return nil
}

// CacheAside implements the cache-aside pattern
type CacheAside[T any] struct {
	cache     CacheStore
	loader    func(ctx context.Context, key string) (T, error)
	ttl       time.Duration
	singleflight map[string]*call[T]
	mu        sync.Mutex
}

type call[T any] struct {
	value T
	err   error
	done  chan struct{}
}

// NewCacheAside creates a cache-aside instance
func NewCacheAside[T any](cache CacheStore, loader func(ctx context.Context, key string) (T, error), ttl time.Duration) *CacheAside[T] {
	return &CacheAside[T]{
		cache:        cache,
		loader:       loader,
		ttl:          ttl,
		singleflight: make(map[string]*call[T]),
	}
}

// Get retrieves from cache or loads from source
func (c *CacheAside[T]) Get(ctx context.Context, key string) (T, error) {
	// Try cache first
	entry, err := c.cache.Get(ctx, key)
	if err == nil && entry != nil && time.Now().Before(entry.ExpiresAt) {
		var value T
		if err := json.Unmarshal(entry.Value, &value); err == nil {
			return value, nil
		}
	}

	// Cache miss - load from source with singleflight
	return c.loadWithSingleflight(ctx, key)
}

func (c *CacheAside[T]) loadWithSingleflight(ctx context.Context, key string) (T, error) {
	c.mu.Lock()
	if existing, ok := c.singleflight[key]; ok {
		c.mu.Unlock()
		// Wait for existing call
		<-existing.done
		return existing.value, existing.err
	}

	// Start new call
	call := &call[T]{done: make(chan struct{})}
	c.singleflight[key] = call
	c.mu.Unlock()

	// Load from source
	call.value, call.err = c.loader(ctx, key)

	// Cache result if successful
	if call.err == nil {
		data, _ := json.Marshal(call.value)
		c.cache.Set(ctx, &Entry{
			Key:       key,
			Value:     data,
			CreatedAt: time.Now(),
			ExpiresAt: time.Now().Add(c.ttl),
		})
	}

	// Cleanup and notify waiters
	c.mu.Lock()
	delete(c.singleflight, key)
	c.mu.Unlock()
	close(call.done)

	return call.value, call.err
}

// Invalidate removes a key from cache
func (c *CacheAside[T]) Invalidate(ctx context.Context, key string) error {
	return c.cache.Delete(ctx, key)
}

// InMemoryCache is an in-memory cache implementation
type InMemoryCache struct {
	entries map[string]*Entry
	tags    map[string][]string // tag -> keys
	maxSize int
	mu      sync.RWMutex
}

// NewInMemoryCache creates an in-memory cache
func NewInMemoryCache(maxSize int) *InMemoryCache {
	return &InMemoryCache{
		entries: make(map[string]*Entry),
		tags:    make(map[string][]string),
		maxSize: maxSize,
	}
}

func (c *InMemoryCache) Get(ctx context.Context, key string) (*Entry, error) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	entry, ok := c.entries[key]
	if !ok {
		return nil, ErrCacheMiss
	}

	if time.Now().After(entry.ExpiresAt) {
		return nil, ErrCacheExpired
	}

	// Update sliding expiration
	if entry.Sliding > 0 {
		entry.ExpiresAt = time.Now().Add(entry.Sliding)
	}
	entry.Hits++

	return entry, nil
}

func (c *InMemoryCache) Set(ctx context.Context, entry *Entry) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	// Evict if needed
	if len(c.entries) >= c.maxSize {
		c.evictOne()
	}

	c.entries[entry.Key] = entry

	// Index tags
	for _, tag := range entry.Tags {
		c.tags[tag] = append(c.tags[tag], entry.Key)
	}

	return nil
}

func (c *InMemoryCache) Delete(ctx context.Context, key string) error {
	c.mu.Lock()
	defer c.mu.Unlock()
	delete(c.entries, key)
	return nil
}

func (c *InMemoryCache) DeleteByTag(ctx context.Context, tag string) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	keys := c.tags[tag]
	for _, key := range keys {
		delete(c.entries, key)
	}
	delete(c.tags, tag)

	return nil
}

func (c *InMemoryCache) Clear(ctx context.Context) error {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.entries = make(map[string]*Entry)
	c.tags = make(map[string][]string)
	return nil
}

func (c *InMemoryCache) evictOne() {
	// LRU eviction - remove oldest entry
	var oldest *Entry
	var oldestKey string

	for key, entry := range c.entries {
		if oldest == nil || entry.CreatedAt.Before(oldest.CreatedAt) {
			oldest = entry
			oldestKey = key
		}
	}

	if oldestKey != "" {
		delete(c.entries, oldestKey)
	}
}

// CacheWarmer pre-loads cache entries
type CacheWarmer struct {
	cache   CacheStore
	loaders map[string]func(ctx context.Context) ([]*Entry, error)
	mu      sync.RWMutex
}

// NewCacheWarmer creates a cache warmer
func NewCacheWarmer(cache CacheStore) *CacheWarmer {
	return &CacheWarmer{
		cache:   cache,
		loaders: make(map[string]func(ctx context.Context) ([]*Entry, error)),
	}
}

// Register registers a loader
func (w *CacheWarmer) Register(name string, loader func(ctx context.Context) ([]*Entry, error)) {
	w.mu.Lock()
	defer w.mu.Unlock()
	w.loaders[name] = loader
}

// WarmUp warms up the cache
func (w *CacheWarmer) WarmUp(ctx context.Context) error {
	w.mu.RLock()
	defer w.mu.RUnlock()

	for _, loader := range w.loaders {
		entries, err := loader(ctx)
		if err != nil {
			continue
		}
		for _, entry := range entries {
			w.cache.Set(ctx, entry)
		}
	}

	return nil
}

// Stats returns cache statistics
type Stats struct {
	Size      int
	Hits      int64
	Misses    int64
	Evictions int64
}

// GetStats returns stats for InMemoryCache
func (c *InMemoryCache) GetStats() Stats {
	c.mu.RLock()
	defer c.mu.RUnlock()

	var totalHits int64
	for _, entry := range c.entries {
		totalHits += entry.Hits
	}

	return Stats{
		Size: len(c.entries),
		Hits: totalHits,
	}
}
