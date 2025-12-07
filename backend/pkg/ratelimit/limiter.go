package ratelimit

import (
	"context"
	"net/http"
	"sync"
	"time"
)

// TokenBucket implements the token bucket algorithm
type TokenBucket struct {
	capacity     int64         // Maximum tokens
	tokens       float64       // Current tokens
	refillRate   float64       // Tokens per second
	lastRefill   time.Time
	mu           sync.Mutex
}

// NewTokenBucket creates a new token bucket
func NewTokenBucket(capacity int64, refillRate float64) *TokenBucket {
	return &TokenBucket{
		capacity:   capacity,
		tokens:     float64(capacity),
		refillRate: refillRate,
		lastRefill: time.Now(),
	}
}

// Allow checks if a request is allowed and consumes a token
func (tb *TokenBucket) Allow() bool {
	return tb.AllowN(1)
}

// AllowN checks if n tokens are available
func (tb *TokenBucket) AllowN(n int64) bool {
	tb.mu.Lock()
	defer tb.mu.Unlock()

	tb.refill()

	if tb.tokens >= float64(n) {
		tb.tokens -= float64(n)
		return true
	}
	return false
}

// refill adds tokens based on elapsed time
func (tb *TokenBucket) refill() {
	now := time.Now()
	elapsed := now.Sub(tb.lastRefill).Seconds()
	tb.tokens += elapsed * tb.refillRate

	if tb.tokens > float64(tb.capacity) {
		tb.tokens = float64(tb.capacity)
	}
	tb.lastRefill = now
}

// Tokens returns current token count
func (tb *TokenBucket) Tokens() float64 {
	tb.mu.Lock()
	defer tb.mu.Unlock()
	tb.refill()
	return tb.tokens
}

// SlidingWindow implements sliding window rate limiting
type SlidingWindow struct {
	windowSize time.Duration
	limit      int
	requests   []time.Time
	mu         sync.Mutex
}

// NewSlidingWindow creates a sliding window limiter
func NewSlidingWindow(limit int, windowSize time.Duration) *SlidingWindow {
	return &SlidingWindow{
		windowSize: windowSize,
		limit:      limit,
		requests:   make([]time.Time, 0),
	}
}

// Allow checks if request is within rate limit
func (sw *SlidingWindow) Allow() bool {
	sw.mu.Lock()
	defer sw.mu.Unlock()

	now := time.Now()
	cutoff := now.Add(-sw.windowSize)

	// Remove old requests
	valid := make([]time.Time, 0)
	for _, t := range sw.requests {
		if t.After(cutoff) {
			valid = append(valid, t)
		}
	}
	sw.requests = valid

	if len(sw.requests) >= sw.limit {
		return false
	}

	sw.requests = append(sw.requests, now)
	return true
}

// LeakyBucket implements the leaky bucket algorithm
type LeakyBucket struct {
	capacity   int64
	leakRate   float64  // Items leaked per second
	water      float64  // Current water level
	lastLeak   time.Time
	mu         sync.Mutex
}

// NewLeakyBucket creates a new leaky bucket
func NewLeakyBucket(capacity int64, leakRate float64) *LeakyBucket {
	return &LeakyBucket{
		capacity: capacity,
		leakRate: leakRate,
		water:    0,
		lastLeak: time.Now(),
	}
}

// Allow adds water to the bucket if there's room
func (lb *LeakyBucket) Allow() bool {
	lb.mu.Lock()
	defer lb.mu.Unlock()

	lb.leak()

	if lb.water < float64(lb.capacity) {
		lb.water++
		return true
	}
	return false
}

func (lb *LeakyBucket) leak() {
	now := time.Now()
	elapsed := now.Sub(lb.lastLeak).Seconds()
	leaked := elapsed * lb.leakRate
	lb.water -= leaked
	if lb.water < 0 {
		lb.water = 0
	}
	lb.lastLeak = now
}

// AdaptiveRateLimiter adjusts limits based on system load
type AdaptiveRateLimiter struct {
	base        *TokenBucket
	minRate     float64
	maxRate     float64
	currentLoad func() float64 // Returns load 0.0-1.0
	mu          sync.RWMutex
}

// NewAdaptiveRateLimiter creates an adaptive limiter
func NewAdaptiveRateLimiter(capacity int64, minRate, maxRate float64, loadFn func() float64) *AdaptiveRateLimiter {
	return &AdaptiveRateLimiter{
		base:        NewTokenBucket(capacity, maxRate),
		minRate:     minRate,
		maxRate:     maxRate,
		currentLoad: loadFn,
	}
}

// Allow adjusts rate based on load and checks permission
func (arl *AdaptiveRateLimiter) Allow() bool {
	load := arl.currentLoad()
	
	// Adjust refill rate based on load
	arl.mu.Lock()
	arl.base.refillRate = arl.maxRate - (arl.maxRate-arl.minRate)*load
	arl.mu.Unlock()

	return arl.base.Allow()
}

// DistributedLimiter interface for distributed rate limiting
type DistributedLimiter interface {
	Allow(ctx context.Context, key string) (bool, error)
	AllowN(ctx context.Context, key string, n int64) (bool, error)
}

// PerClientLimiter tracks limits per client
type PerClientLimiter struct {
	limiters sync.Map
	capacity int64
	rate     float64
}

// NewPerClientLimiter creates a per-client rate limiter
func NewPerClientLimiter(capacity int64, rate float64) *PerClientLimiter {
	return &PerClientLimiter{
		capacity: capacity,
		rate:     rate,
	}
}

// Allow checks if client's request is allowed
func (pcl *PerClientLimiter) Allow(clientID string) bool {
	limiter, _ := pcl.limiters.LoadOrStore(clientID, NewTokenBucket(pcl.capacity, pcl.rate))
	return limiter.(*TokenBucket).Allow()
}

// HTTPMiddleware creates rate limiting middleware
func HTTPMiddleware(limiter *TokenBucket) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if !limiter.Allow() {
				w.Header().Set("Retry-After", "1")
				http.Error(w, "Rate limit exceeded", http.StatusTooManyRequests)
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}

// PerClientMiddleware creates per-client rate limiting
func PerClientMiddleware(limiter *PerClientLimiter, keyFn func(*http.Request) string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			clientID := keyFn(r)
			if !limiter.Allow(clientID) {
				w.Header().Set("Retry-After", "1")
				http.Error(w, "Rate limit exceeded", http.StatusTooManyRequests)
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}
