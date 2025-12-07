package metrics

import (
	"net/http"
	"runtime"
	"sync"
	"sync/atomic"
	"time"
)

// Counter is a monotonically increasing counter
type Counter struct {
	value uint64
}

// NewCounter creates a new counter
func NewCounter() *Counter {
	return &Counter{}
}

// Inc increments the counter by 1
func (c *Counter) Inc() {
	atomic.AddUint64(&c.value, 1)
}

// Add increments the counter by n
func (c *Counter) Add(n uint64) {
	atomic.AddUint64(&c.value, n)
}

// Value returns the current count
func (c *Counter) Value() uint64 {
	return atomic.LoadUint64(&c.value)
}

// Gauge represents a value that can go up and down
type Gauge struct {
	value int64
}

// NewGauge creates a new gauge
func NewGauge() *Gauge {
	return &Gauge{}
}

// Set sets the gauge value
func (g *Gauge) Set(n int64) {
	atomic.StoreInt64(&g.value, n)
}

// Inc increments the gauge by 1
func (g *Gauge) Inc() {
	atomic.AddInt64(&g.value, 1)
}

// Dec decrements the gauge by 1
func (g *Gauge) Dec() {
	atomic.AddInt64(&g.value, -1)
}

// Value returns the current value
func (g *Gauge) Value() int64 {
	return atomic.LoadInt64(&g.value)
}

// Histogram tracks distribution of values
type Histogram struct {
	buckets     []uint64
	boundaries  []float64
	sum         float64
	count       uint64
	mu          sync.Mutex
}

// NewHistogram creates a histogram with specified bucket boundaries
func NewHistogram(boundaries []float64) *Histogram {
	return &Histogram{
		buckets:    make([]uint64, len(boundaries)+1),
		boundaries: boundaries,
	}
}

// DefaultLatencyBuckets returns common latency buckets in seconds
func DefaultLatencyBuckets() []float64 {
	return []float64{0.001, 0.005, 0.01, 0.025, 0.05, 0.1, 0.25, 0.5, 1, 2.5, 5, 10}
}

// Observe records a value
func (h *Histogram) Observe(value float64) {
	h.mu.Lock()
	defer h.mu.Unlock()

	h.count++
	h.sum += value

	// Find bucket
	for i, boundary := range h.boundaries {
		if value <= boundary {
			h.buckets[i]++
			return
		}
	}
	h.buckets[len(h.buckets)-1]++
}

// Percentile returns the approximate percentile value
func (h *Histogram) Percentile(p float64) float64 {
	h.mu.Lock()
	defer h.mu.Unlock()

	if h.count == 0 {
		return 0
	}

	target := uint64(float64(h.count) * p)
	var cumulative uint64
	for i, count := range h.buckets {
		cumulative += count
		if cumulative >= target {
			if i == 0 {
				return h.boundaries[0]
			}
			if i >= len(h.boundaries) {
				return h.boundaries[len(h.boundaries)-1]
			}
			return h.boundaries[i]
		}
	}
	return h.boundaries[len(h.boundaries)-1]
}

// Summary returns histogram statistics
func (h *Histogram) Summary() HistogramSummary {
	h.mu.Lock()
	defer h.mu.Unlock()

	avg := float64(0)
	if h.count > 0 {
		avg = h.sum / float64(h.count)
	}
	return HistogramSummary{
		Count: h.count,
		Sum:   h.sum,
		Avg:   avg,
	}
}

// HistogramSummary contains histogram statistics
type HistogramSummary struct {
	Count uint64
	Sum   float64
	Avg   float64
}

// Timer measures elapsed time
type Timer struct {
	histogram *Histogram
	start     time.Time
}

// NewTimer starts a new timer
func NewTimer(h *Histogram) *Timer {
	return &Timer{
		histogram: h,
		start:     time.Now(),
	}
}

// ObserveDuration records the elapsed time
func (t *Timer) ObserveDuration() {
	t.histogram.Observe(time.Since(t.start).Seconds())
}

// Registry holds all metrics
type Registry struct {
	counters   map[string]*Counter
	gauges     map[string]*Gauge
	histograms map[string]*Histogram
	mu         sync.RWMutex
}

// NewRegistry creates a new metric registry
func NewRegistry() *Registry {
	return &Registry{
		counters:   make(map[string]*Counter),
		gauges:     make(map[string]*Gauge),
		histograms: make(map[string]*Histogram),
	}
}

// Counter gets or creates a counter
func (r *Registry) Counter(name string) *Counter {
	r.mu.Lock()
	defer r.mu.Unlock()

	if c, ok := r.counters[name]; ok {
		return c
	}
	c := NewCounter()
	r.counters[name] = c
	return c
}

// Gauge gets or creates a gauge
func (r *Registry) Gauge(name string) *Gauge {
	r.mu.Lock()
	defer r.mu.Unlock()

	if g, ok := r.gauges[name]; ok {
		return g
	}
	g := NewGauge()
	r.gauges[name] = g
	return g
}

// Histogram gets or creates a histogram
func (r *Registry) Histogram(name string, buckets []float64) *Histogram {
	r.mu.Lock()
	defer r.mu.Unlock()

	if h, ok := r.histograms[name]; ok {
		return h
	}
	h := NewHistogram(buckets)
	r.histograms[name] = h
	return h
}

// Snapshot returns all current metric values
func (r *Registry) Snapshot() map[string]interface{} {
	r.mu.RLock()
	defer r.mu.RUnlock()

	snapshot := make(map[string]interface{})

	for name, c := range r.counters {
		snapshot[name] = c.Value()
	}
	for name, g := range r.gauges {
		snapshot[name] = g.Value()
	}
	for name, h := range r.histograms {
		snapshot[name+"_summary"] = h.Summary()
		snapshot[name+"_p50"] = h.Percentile(0.5)
		snapshot[name+"_p95"] = h.Percentile(0.95)
		snapshot[name+"_p99"] = h.Percentile(0.99)
	}

	return snapshot
}

// RuntimeMetrics collects Go runtime metrics
type RuntimeMetrics struct {
	registry *Registry
}

// NewRuntimeMetrics creates runtime metric collectors
func NewRuntimeMetrics(r *Registry) *RuntimeMetrics {
	rm := &RuntimeMetrics{registry: r}
	go rm.collect()
	return rm
}

func (rm *RuntimeMetrics) collect() {
	ticker := time.NewTicker(10 * time.Second)
	defer ticker.Stop()

	for range ticker.C {
		var m runtime.MemStats
		runtime.ReadMemStats(&m)

		rm.registry.Gauge("go_goroutines").Set(int64(runtime.NumGoroutine()))
		rm.registry.Gauge("go_heap_alloc_bytes").Set(int64(m.HeapAlloc))
		rm.registry.Gauge("go_heap_objects").Set(int64(m.HeapObjects))
		rm.registry.Gauge("go_gc_runs_total").Set(int64(m.NumGC))
	}
}

// HTTPMiddleware creates metrics middleware
func HTTPMiddleware(r *Registry) func(http.Handler) http.Handler {
	requestCounter := r.Counter("http_requests_total")
	requestDuration := r.Histogram("http_request_duration_seconds", DefaultLatencyBuckets())
	activeRequests := r.Gauge("http_requests_active")

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
			requestCounter.Inc()
			activeRequests.Inc()
			defer activeRequests.Dec()

			timer := NewTimer(requestDuration)
			defer timer.ObserveDuration()

			next.ServeHTTP(w, req)
		})
	}
}

// Global registry
var DefaultRegistry = NewRegistry()
