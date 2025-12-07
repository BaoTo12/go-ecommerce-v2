package cache

import (
	"context"
	"sync"
	"testing"
	"time"
)

func TestInMemoryCache_GetSet(t *testing.T) {
	cache := NewInMemoryCache(100)
	ctx := context.Background()

	entry := &Entry{
		Key:       "test-key",
		Value:     []byte("test-value"),
		ExpiresAt: time.Now().Add(time.Hour),
	}

	// Set
	err := cache.Set(ctx, entry)
	if err != nil {
		t.Fatalf("Set failed: %v", err)
	}

	// Get
	retrieved, err := cache.Get(ctx, "test-key")
	if err != nil {
		t.Fatalf("Get failed: %v", err)
	}

	if string(retrieved.Value) != "test-value" {
		t.Errorf("Expected 'test-value', got '%s'", string(retrieved.Value))
	}
}

func TestInMemoryCache_Expiration(t *testing.T) {
	cache := NewInMemoryCache(100)
	ctx := context.Background()

	entry := &Entry{
		Key:       "expires-soon",
		Value:     []byte("value"),
		ExpiresAt: time.Now().Add(50 * time.Millisecond),
	}

	cache.Set(ctx, entry)

	// Should exist
	_, err := cache.Get(ctx, "expires-soon")
	if err != nil {
		t.Error("Entry should exist before expiration")
	}

	// Wait for expiration
	time.Sleep(100 * time.Millisecond)

	// Should be expired
	_, err = cache.Get(ctx, "expires-soon")
	if err != ErrCacheExpired {
		t.Error("Entry should be expired")
	}
}

func TestInMemoryCache_Tags(t *testing.T) {
	cache := NewInMemoryCache(100)
	ctx := context.Background()

	// Add entries with tags
	cache.Set(ctx, &Entry{
		Key:       "product-1",
		Value:     []byte("product 1"),
		Tags:      []string{"products", "category-a"},
		ExpiresAt: time.Now().Add(time.Hour),
	})

	cache.Set(ctx, &Entry{
		Key:       "product-2",
		Value:     []byte("product 2"),
		Tags:      []string{"products", "category-b"},
		ExpiresAt: time.Now().Add(time.Hour),
	})

	cache.Set(ctx, &Entry{
		Key:       "user-1",
		Value:     []byte("user 1"),
		Tags:      []string{"users"},
		ExpiresAt: time.Now().Add(time.Hour),
	})

	// Delete by tag
	cache.DeleteByTag(ctx, "products")

	// Product entries should be gone
	_, err := cache.Get(ctx, "product-1")
	if err != ErrCacheMiss {
		t.Error("product-1 should be deleted")
	}

	// User entry should still exist
	_, err = cache.Get(ctx, "user-1")
	if err != nil {
		t.Error("user-1 should still exist")
	}
}

func TestInMemoryCache_Eviction(t *testing.T) {
	cache := NewInMemoryCache(3)
	ctx := context.Background()

	for i := 0; i < 5; i++ {
		cache.Set(ctx, &Entry{
			Key:       string(rune('a' + i)),
			Value:     []byte("value"),
			CreatedAt: time.Now(),
			ExpiresAt: time.Now().Add(time.Hour),
		})
	}

	stats := cache.GetStats()
	if stats.Size > 3 {
		t.Errorf("Cache size should not exceed 3, got %d", stats.Size)
	}
}

func TestCacheAside(t *testing.T) {
	cache := NewInMemoryCache(100)
	loadCount := 0

	loader := func(ctx context.Context, key string) (string, error) {
		loadCount++
		return "loaded-" + key, nil
	}

	cacheAside := NewCacheAside(cache, loader, 5*time.Minute)
	ctx := context.Background()

	// First get - should load from source
	value, err := cacheAside.Get(ctx, "test-key")
	if err != nil {
		t.Fatalf("Get failed: %v", err)
	}

	if value != "loaded-test-key" {
		t.Errorf("Expected 'loaded-test-key', got '%s'", value)
	}

	if loadCount != 1 {
		t.Error("Loader should be called once")
	}

	// Second get - should hit cache
	value, _ = cacheAside.Get(ctx, "test-key")
	if loadCount != 1 {
		t.Error("Loader should not be called again")
	}
}

func TestCacheAside_Singleflight(t *testing.T) {
	cache := NewInMemoryCache(100)
	loadCount := 0
	var mu sync.Mutex

	loader := func(ctx context.Context, key string) (string, error) {
		mu.Lock()
		loadCount++
		mu.Unlock()
		time.Sleep(50 * time.Millisecond) // Simulate slow load
		return "value", nil
	}

	cacheAside := NewCacheAside(cache, loader, 5*time.Minute)
	ctx := context.Background()

	// Multiple concurrent requests for same key
	var wg sync.WaitGroup
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			cacheAside.Get(ctx, "same-key")
		}()
	}
	wg.Wait()

	// Loader should only be called once (singleflight)
	if loadCount != 1 {
		t.Errorf("Loader should be called once due to singleflight, got %d", loadCount)
	}
}

func TestMultiLevelCache(t *testing.T) {
	l1 := NewInMemoryCache(10)
	l2 := NewInMemoryCache(100)
	mlc := NewMultiLevelCache(l1, l2)
	ctx := context.Background()

	// Set in L2 only
	l2.Set(ctx, &Entry{
		Key:       "test-key",
		Value:     []byte("L2 value"),
		ExpiresAt: time.Now().Add(time.Hour),
	})

	// Get - should fetch from L2 and populate L1
	entry, err := mlc.Get(ctx, "test-key")
	if err != nil {
		t.Fatalf("Get failed: %v", err)
	}

	if string(entry.Value) != "L2 value" {
		t.Errorf("Expected 'L2 value', got '%s'", string(entry.Value))
	}

	// L1 should now have it
	l1Entry, err := l1.Get(ctx, "test-key")
	if err != nil {
		t.Error("L1 should have the entry now")
	}
	if string(l1Entry.Value) != "L2 value" {
		t.Error("L1 should have correct value")
	}
}

func BenchmarkInMemoryCache_Get(b *testing.B) {
	cache := NewInMemoryCache(10000)
	ctx := context.Background()

	// Pre-populate
	for i := 0; i < 1000; i++ {
		cache.Set(ctx, &Entry{
			Key:       string(rune(i)),
			Value:     []byte("value"),
			ExpiresAt: time.Now().Add(time.Hour),
		})
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		cache.Get(ctx, string(rune(i%1000)))
	}
}

func BenchmarkInMemoryCache_Set(b *testing.B) {
	cache := NewInMemoryCache(10000)
	ctx := context.Background()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		cache.Set(ctx, &Entry{
			Key:       string(rune(i % 1000)),
			Value:     []byte("value"),
			ExpiresAt: time.Now().Add(time.Hour),
		})
	}
}
