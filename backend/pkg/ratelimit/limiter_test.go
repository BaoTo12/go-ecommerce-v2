package ratelimit_test

import (
	"net/http"
	"net/http/httptest"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/titan-commerce/backend/pkg/ratelimit"
)

// ====================
// PHASE 1: UNIT TESTS
// ====================

func TestTokenBucket_Allow(t *testing.T) {
	tb := ratelimit.NewTokenBucket(10, 1) // 10 tokens, 1/sec refill

	// Should allow first 10 requests
	for i := 0; i < 10; i++ {
		if !tb.Allow() {
			t.Errorf("Expected request %d to be allowed", i)
		}
	}

	// 11th request should be denied
	if tb.Allow() {
		t.Error("Expected 11th request to be denied")
	}
}

func TestTokenBucket_Refill(t *testing.T) {
	tb := ratelimit.NewTokenBucket(5, 10) // 5 tokens, 10/sec refill

	// Drain all tokens
	for i := 0; i < 5; i++ {
		tb.Allow()
	}

	// Should be denied
	if tb.Allow() {
		t.Error("Expected to be denied after draining")
	}

	// Wait for refill
	time.Sleep(200 * time.Millisecond) // Should refill ~2 tokens

	// Should have tokens now
	if !tb.Allow() {
		t.Error("Expected to be allowed after refill")
	}
}

func TestTokenBucket_AllowN(t *testing.T) {
	tb := ratelimit.NewTokenBucket(10, 1)

	// Request 5 tokens
	if !tb.AllowN(5) {
		t.Error("Expected 5 tokens to be available")
	}

	// Request 6 more (should fail, only 5 left)
	if tb.AllowN(6) {
		t.Error("Expected to fail - not enough tokens")
	}

	// Request 5 more (should succeed)
	if !tb.AllowN(5) {
		t.Error("Expected 5 tokens to be available")
	}
}

func TestSlidingWindow_Allow(t *testing.T) {
	sw := ratelimit.NewSlidingWindow(5, 1*time.Second)

	// Allow first 5
	for i := 0; i < 5; i++ {
		if !sw.Allow() {
			t.Errorf("Expected request %d to be allowed", i)
		}
	}

	// 6th should be denied
	if sw.Allow() {
		t.Error("Expected 6th request to be denied")
	}

	// Wait for window to slide
	time.Sleep(1100 * time.Millisecond)

	// Should be allowed again
	if !sw.Allow() {
		t.Error("Expected to be allowed after window slides")
	}
}

func TestLeakyBucket_Allow(t *testing.T) {
	lb := ratelimit.NewLeakyBucket(5, 10) // 5 capacity, 10/sec leak

	// Fill bucket
	for i := 0; i < 5; i++ {
		if !lb.Allow() {
			t.Errorf("Expected request %d to be allowed", i)
		}
	}

	// Bucket full
	if lb.Allow() {
		t.Error("Expected to be denied - bucket full")
	}

	// Wait for leak
	time.Sleep(200 * time.Millisecond) // Should leak ~2

	if !lb.Allow() {
		t.Error("Expected to be allowed after leak")
	}
}

func TestPerClientLimiter(t *testing.T) {
	pcl := ratelimit.NewPerClientLimiter(5, 1)

	// Client A uses 5 tokens
	for i := 0; i < 5; i++ {
		if !pcl.Allow("clientA") {
			t.Error("Expected clientA to be allowed")
		}
	}

	// Client A should be denied
	if pcl.Allow("clientA") {
		t.Error("Expected clientA to be denied")
	}

	// Client B should still have tokens
	if !pcl.Allow("clientB") {
		t.Error("Expected clientB to be allowed")
	}
}

// ====================
// PHASE 2: MOCK + INTEGRATION TESTS
// ====================

func TestHTTPMiddleware_Integration(t *testing.T) {
	limiter := ratelimit.NewTokenBucket(3, 0) // 3 tokens, no refill

	handler := ratelimit.HTTPMiddleware(limiter)(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))

	// First 3 requests should succeed
	for i := 0; i < 3; i++ {
		req := httptest.NewRequest("GET", "/", nil)
		rec := httptest.NewRecorder()
		handler.ServeHTTP(rec, req)

		if rec.Code != http.StatusOK {
			t.Errorf("Request %d: expected 200, got %d", i, rec.Code)
		}
	}

	// 4th request should be rate limited
	req := httptest.NewRequest("GET", "/", nil)
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusTooManyRequests {
		t.Errorf("Expected 429, got %d", rec.Code)
	}

	// Check Retry-After header
	if rec.Header().Get("Retry-After") == "" {
		t.Error("Expected Retry-After header")
	}
}

func TestPerClientMiddleware_Integration(t *testing.T) {
	limiter := ratelimit.NewPerClientLimiter(2, 0)

	keyFn := func(r *http.Request) string {
		return r.Header.Get("X-Client-ID")
	}

	handler := ratelimit.PerClientMiddleware(limiter, keyFn)(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))

	// Client A makes 2 requests
	for i := 0; i < 2; i++ {
		req := httptest.NewRequest("GET", "/", nil)
		req.Header.Set("X-Client-ID", "clientA")
		rec := httptest.NewRecorder()
		handler.ServeHTTP(rec, req)

		if rec.Code != http.StatusOK {
			t.Errorf("ClientA request %d: expected 200, got %d", i, rec.Code)
		}
	}

	// Client A's 3rd request should fail
	req := httptest.NewRequest("GET", "/", nil)
	req.Header.Set("X-Client-ID", "clientA")
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusTooManyRequests {
		t.Errorf("Expected 429 for clientA, got %d", rec.Code)
	}

	// Client B should still work
	req = httptest.NewRequest("GET", "/", nil)
	req.Header.Set("X-Client-ID", "clientB")
	rec = httptest.NewRecorder()
	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("Expected 200 for clientB, got %d", rec.Code)
	}
}

// ====================
// PHASE 2: BENCHMARK TESTS
// ====================

func BenchmarkTokenBucket_Allow(b *testing.B) {
	tb := ratelimit.NewTokenBucket(1000000, 1000000)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		tb.Allow()
	}
}

func BenchmarkTokenBucket_Parallel(b *testing.B) {
	tb := ratelimit.NewTokenBucket(1000000, 1000000)

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			tb.Allow()
		}
	})
}

func BenchmarkSlidingWindow_Allow(b *testing.B) {
	sw := ratelimit.NewSlidingWindow(1000000, time.Hour)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		sw.Allow()
	}
}

func BenchmarkPerClientLimiter(b *testing.B) {
	pcl := ratelimit.NewPerClientLimiter(1000000, 1000000)
	clients := []string{"client1", "client2", "client3", "client4", "client5"}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		pcl.Allow(clients[i%5])
	}
}

// ====================
// PHASE 3: FUZZ TESTS
// ====================

func FuzzTokenBucket_AllowN(f *testing.F) {
	f.Add(int64(1))
	f.Add(int64(10))
	f.Add(int64(100))
	f.Add(int64(0))

	tb := ratelimit.NewTokenBucket(1000, 10)

	f.Fuzz(func(t *testing.T, n int64) {
		if n < 0 {
			n = -n
		}
		// Should not panic
		tb.AllowN(n)
	})
}

// ====================
// PHASE 3: CONCURRENCY/STRESS TESTS
// ====================

func TestTokenBucket_ConcurrentAccess(t *testing.T) {
	tb := ratelimit.NewTokenBucket(1000, 100)

	var allowed int64
	var denied int64
	var wg sync.WaitGroup

	// 100 goroutines making 100 requests each
	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for j := 0; j < 100; j++ {
				if tb.Allow() {
					atomic.AddInt64(&allowed, 1)
				} else {
					atomic.AddInt64(&denied, 1)
				}
			}
		}()
	}

	wg.Wait()

	total := allowed + denied
	if total != 10000 {
		t.Errorf("Expected 10000 total requests, got %d", total)
	}

	// With 1000 tokens and 100/sec refill over ~1 second, should allow ~1100
	if allowed > 2000 {
		t.Errorf("Allowed too many: %d", allowed)
	}
}
