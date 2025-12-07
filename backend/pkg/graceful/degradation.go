package graceful

import (
	"context"
	"net/http"
	"sync"
	"sync/atomic"
	"time"
)

// Degradation levels
type Level int

const (
	LevelNormal Level = iota
	LevelDegraded
	LevelEmergency
	LevelMaintenance
)

func (l Level) String() string {
	switch l {
	case LevelNormal:
		return "normal"
	case LevelDegraded:
		return "degraded"
	case LevelEmergency:
		return "emergency"
	case LevelMaintenance:
		return "maintenance"
	default:
		return "unknown"
	}
}

// Degrader manages graceful degradation
type Degrader struct {
	level       int32
	features    map[string]Level // Feature -> minimum level required
	healthCheck func() bool
	mu          sync.RWMutex
}

// NewDegrader creates a new degrader
func NewDegrader() *Degrader {
	return &Degrader{
		level:    int32(LevelNormal),
		features: make(map[string]Level),
	}
}

// SetLevel sets the current degradation level
func (d *Degrader) SetLevel(level Level) {
	atomic.StoreInt32(&d.level, int32(level))
}

// Level returns the current degradation level
func (d *Degrader) Level() Level {
	return Level(atomic.LoadInt32(&d.level))
}

// RegisterFeature registers a feature with its minimum required level
func (d *Degrader) RegisterFeature(name string, minLevel Level) {
	d.mu.Lock()
	defer d.mu.Unlock()
	d.features[name] = minLevel
}

// IsFeatureEnabled checks if a feature is enabled at current level
func (d *Degrader) IsFeatureEnabled(name string) bool {
	d.mu.RLock()
	minLevel, exists := d.features[name]
	d.mu.RUnlock()

	if !exists {
		return true // Unknown features are enabled by default
	}

	return d.Level() <= minLevel
}

// AutoDegrade monitors health and adjusts level
func (d *Degrader) AutoDegrade(ctx context.Context, healthCheck func() bool, interval time.Duration) {
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	consecutiveFailures := 0

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			if healthCheck() {
				consecutiveFailures = 0
				if d.Level() != LevelNormal {
					d.SetLevel(LevelNormal)
				}
			} else {
				consecutiveFailures++
				switch {
				case consecutiveFailures >= 5:
					d.SetLevel(LevelEmergency)
				case consecutiveFailures >= 2:
					d.SetLevel(LevelDegraded)
				}
			}
		}
	}
}

// FeatureMiddleware disables features based on degradation level
func FeatureMiddleware(d *Degrader, feature string, fallback http.Handler) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if d.IsFeatureEnabled(feature) {
				next.ServeHTTP(w, r)
			} else if fallback != nil {
				fallback.ServeHTTP(w, r)
			} else {
				w.Header().Set("X-Feature-Disabled", feature)
				http.Error(w, "Service temporarily unavailable", http.StatusServiceUnavailable)
			}
		})
	}
}

// LoadShedder drops requests under heavy load
type LoadShedder struct {
	maxConcurrent int32
	current       int32
	waitQueue     int32
	maxQueue      int32
}

// NewLoadShedder creates a new load shedder
func NewLoadShedder(maxConcurrent, maxQueue int) *LoadShedder {
	return &LoadShedder{
		maxConcurrent: int32(maxConcurrent),
		maxQueue:      int32(maxQueue),
	}
}

// Acquire attempts to acquire a slot
func (ls *LoadShedder) Acquire() bool {
	current := atomic.LoadInt32(&ls.current)
	if current >= ls.maxConcurrent {
		// Try to enter wait queue
		waiting := atomic.LoadInt32(&ls.waitQueue)
		if waiting >= ls.maxQueue {
			return false // Shed load
		}
		atomic.AddInt32(&ls.waitQueue, 1)
		// In real impl, would wait here
		atomic.AddInt32(&ls.waitQueue, -1)
	}
	atomic.AddInt32(&ls.current, 1)
	return true
}

// Release releases a slot
func (ls *LoadShedder) Release() {
	atomic.AddInt32(&ls.current, -1)
}

// Stats returns load shedder statistics
func (ls *LoadShedder) Stats() LoadShedderStats {
	return LoadShedderStats{
		Current:       int(atomic.LoadInt32(&ls.current)),
		MaxConcurrent: int(ls.maxConcurrent),
		Waiting:       int(atomic.LoadInt32(&ls.waitQueue)),
		MaxQueue:      int(ls.maxQueue),
	}
}

// LoadShedderStats contains load shedder metrics
type LoadShedderStats struct {
	Current       int
	MaxConcurrent int
	Waiting       int
	MaxQueue      int
}

// LoadSheddingMiddleware drops requests under load
func LoadSheddingMiddleware(ls *LoadShedder) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if !ls.Acquire() {
				w.Header().Set("Retry-After", "5")
				http.Error(w, "Service overloaded", http.StatusServiceUnavailable)
				return
			}
			defer ls.Release()
			next.ServeHTTP(w, r)
		})
	}
}

// Fallback executes fallback on error
type Fallback struct{}

// Execute runs fn with fallback on error
func (f *Fallback) Execute(fn func() (interface{}, error), fallbackFn func() (interface{}, error)) (interface{}, error) {
	result, err := fn()
	if err != nil {
		return fallbackFn()
	}
	return result, nil
}

// ExecuteWithTimeout adds timeout
func (f *Fallback) ExecuteWithTimeout(ctx context.Context, timeout time.Duration, fn func() (interface{}, error), fallbackFn func() (interface{}, error)) (interface{}, error) {
	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	done := make(chan struct {
		result interface{}
		err    error
	}, 1)

	go func() {
		result, err := fn()
		done <- struct {
			result interface{}
			err    error
		}{result, err}
	}()

	select {
	case r := <-done:
		if r.err != nil {
			return fallbackFn()
		}
		return r.result, nil
	case <-ctx.Done():
		return fallbackFn()
	}
}
