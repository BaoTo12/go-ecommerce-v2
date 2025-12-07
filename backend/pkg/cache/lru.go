package cache

import (
	"container/list"
	"sync"
	"time"
)

// Cache entry
type entry struct {
	key       string
	value     interface{}
	expiresAt time.Time
}

// LRUCache is a thread-safe LRU cache with TTL
type LRUCache struct {
	capacity int
	ttl      time.Duration
	items    map[string]*list.Element
	order    *list.List
	mu       sync.RWMutex
}

// NewLRUCache creates a new LRU cache
func NewLRUCache(capacity int, ttl time.Duration) *LRUCache {
	cache := &LRUCache{
		capacity: capacity,
		ttl:      ttl,
		items:    make(map[string]*list.Element),
		order:    list.New(),
	}
	// Start cleanup goroutine
	go cache.cleanup()
	return cache
}

// Get retrieves a value from the cache
func (c *LRUCache) Get(key string) (interface{}, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	if elem, ok := c.items[key]; ok {
		e := elem.Value.(*entry)
		if time.Now().Before(e.expiresAt) {
			c.order.MoveToFront(elem)
			return e.value, true
		}
	}
	return nil, false
}

// Set adds a value to the cache
func (c *LRUCache) Set(key string, value interface{}) {
	c.mu.Lock()
	defer c.mu.Unlock()

	// Update existing entry
	if elem, ok := c.items[key]; ok {
		c.order.MoveToFront(elem)
		e := elem.Value.(*entry)
		e.value = value
		e.expiresAt = time.Now().Add(c.ttl)
		return
	}

	// Evict if at capacity
	if c.order.Len() >= c.capacity {
		oldest := c.order.Back()
		if oldest != nil {
			c.order.Remove(oldest)
			delete(c.items, oldest.Value.(*entry).key)
		}
	}

	// Add new entry
	e := &entry{
		key:       key,
		value:     value,
		expiresAt: time.Now().Add(c.ttl),
	}
	elem := c.order.PushFront(e)
	c.items[key] = elem
}

// Delete removes a value from the cache
func (c *LRUCache) Delete(key string) {
	c.mu.Lock()
	defer c.mu.Unlock()

	if elem, ok := c.items[key]; ok {
		c.order.Remove(elem)
		delete(c.items, key)
	}
}

// Clear removes all entries from the cache
func (c *LRUCache) Clear() {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.items = make(map[string]*list.Element)
	c.order.Init()
}

// Size returns the number of items in the cache
func (c *LRUCache) Size() int {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.order.Len()
}

// cleanup removes expired entries periodically
func (c *LRUCache) cleanup() {
	ticker := time.NewTicker(1 * time.Minute)
	defer ticker.Stop()

	for range ticker.C {
		c.mu.Lock()
		now := time.Now()
		for key, elem := range c.items {
			e := elem.Value.(*entry)
			if now.After(e.expiresAt) {
				c.order.Remove(elem)
				delete(c.items, key)
			}
		}
		c.mu.Unlock()
	}
}

// CacheKey generates a cache key from parts
func CacheKey(parts ...string) string {
	key := ""
	for i, part := range parts {
		if i > 0 {
			key += ":"
		}
		key += part
	}
	return key
}
