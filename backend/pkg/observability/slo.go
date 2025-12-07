package observability

import (
	"context"
	"math"
	"sync"
	"time"
)

/*
SLO/SLI MONITORING & ANOMALY DETECTION

Implements Service Level Objectives (SLO) and Service Level Indicators (SLI)
monitoring with automatic anomaly detection.

Features:
- Error budget tracking
- Latency percentiles
- Availability calculation
- Anomaly detection using z-score
- Alert generation
*/

// SLOConfig defines SLO configuration
type SLOConfig struct {
	Name          string
	Target        float64       // e.g., 0.999 for 99.9% availability
	Window        time.Duration // e.g., 30 days
	BurnRateAlert float64       // Alert when burn rate exceeds this
}

// SLI represents a Service Level Indicator
type SLI struct {
	Name        string
	TotalCount  int64
	GoodCount   int64
	BadCount    int64
	LastUpdated time.Time
}

// SLO represents a Service Level Objective
type SLO struct {
	config       SLOConfig
	sli          *SLI
	buckets      []timeBucket
	bucketSize   time.Duration
	currentBucket int
	mu           sync.RWMutex
}

type timeBucket struct {
	Total int64
	Good  int64
	Start time.Time
}

// NewSLO creates a new SLO
func NewSLO(config SLOConfig) *SLO {
	// Create buckets for rolling window
	bucketCount := int(config.Window / time.Hour)
	if bucketCount < 1 {
		bucketCount = 24 // Minimum 24 buckets
	}

	buckets := make([]timeBucket, bucketCount)
	now := time.Now()
	for i := range buckets {
		buckets[i].Start = now.Add(-time.Duration(bucketCount-i) * time.Hour)
	}

	return &SLO{
		config:     config,
		sli:        &SLI{Name: config.Name},
		buckets:    buckets,
		bucketSize: config.Window / time.Duration(bucketCount),
	}
}

// RecordSuccess records a successful event
func (s *SLO) RecordSuccess() {
	s.mu.Lock()
	defer s.mu.Unlock()
	
	s.rotateIfNeeded()
	s.sli.TotalCount++
	s.sli.GoodCount++
	s.buckets[s.currentBucket].Total++
	s.buckets[s.currentBucket].Good++
	s.sli.LastUpdated = time.Now()
}

// RecordFailure records a failed event
func (s *SLO) RecordFailure() {
	s.mu.Lock()
	defer s.mu.Unlock()
	
	s.rotateIfNeeded()
	s.sli.TotalCount++
	s.sli.BadCount++
	s.buckets[s.currentBucket].Total++
	s.sli.LastUpdated = time.Now()
}

func (s *SLO) rotateIfNeeded() {
	now := time.Now()
	currentStart := s.buckets[s.currentBucket].Start
	
	if now.Sub(currentStart) >= s.bucketSize {
		// Move to next bucket
		s.currentBucket = (s.currentBucket + 1) % len(s.buckets)
		s.buckets[s.currentBucket] = timeBucket{Start: now}
	}
}

// CurrentSLI returns the current SLI value
func (s *SLO) CurrentSLI() float64 {
	s.mu.RLock()
	defer s.mu.RUnlock()
	
	var total, good int64
	for _, b := range s.buckets {
		total += b.Total
		good += b.Good
	}
	
	if total == 0 {
		return 1.0 // No data = 100%
	}
	return float64(good) / float64(total)
}

// ErrorBudget returns the remaining error budget
func (s *SLO) ErrorBudget() float64 {
	current := s.CurrentSLI()
	allowed := 1 - s.config.Target // e.g., 0.001 for 99.9%
	consumed := 1 - current
	
	if allowed == 0 {
		return 0
	}
	return 1 - (consumed / allowed)
}

// BurnRate returns the current error budget burn rate
func (s *SLO) BurnRate() float64 {
	s.mu.RLock()
	defer s.mu.RUnlock()
	
	// Calculate burn rate over last hour
	now := time.Now()
	var total, bad int64
	
	for _, b := range s.buckets {
		if now.Sub(b.Start) <= time.Hour {
			total += b.Total
			bad += b.Total - b.Good
		}
	}
	
	if total == 0 {
		return 0
	}
	
	currentRate := float64(bad) / float64(total)
	allowedRate := 1 - s.config.Target
	
	if allowedRate == 0 {
		return 0
	}
	return currentRate / allowedRate
}

// IsBreaching returns true if SLO is being breached
func (s *SLO) IsBreaching() bool {
	return s.CurrentSLI() < s.config.Target
}

// ShouldAlert returns true if burn rate exceeds threshold
func (s *SLO) ShouldAlert() bool {
	return s.BurnRate() > s.config.BurnRateAlert
}

// AnomalyDetector detects anomalies using statistical methods
type AnomalyDetector struct {
	window     []float64
	windowSize int
	threshold  float64 // z-score threshold
	mu         sync.Mutex
}

// NewAnomalyDetector creates an anomaly detector
func NewAnomalyDetector(windowSize int, zScoreThreshold float64) *AnomalyDetector {
	return &AnomalyDetector{
		window:     make([]float64, 0, windowSize),
		windowSize: windowSize,
		threshold:  zScoreThreshold,
	}
}

// AddDataPoint adds a data point and returns if it's anomalous
func (d *AnomalyDetector) AddDataPoint(value float64) (isAnomaly bool, zScore float64) {
	d.mu.Lock()
	defer d.mu.Unlock()

	// Need at least a few data points
	if len(d.window) < 3 {
		d.window = append(d.window, value)
		return false, 0
	}

	// Calculate z-score
	mean := d.mean()
	stdDev := d.stdDev(mean)
	
	if stdDev == 0 {
		d.addToWindow(value)
		return false, 0
	}

	zScore = (value - mean) / stdDev
	isAnomaly = math.Abs(zScore) > d.threshold

	d.addToWindow(value)
	return isAnomaly, zScore
}

func (d *AnomalyDetector) addToWindow(value float64) {
	if len(d.window) >= d.windowSize {
		d.window = d.window[1:]
	}
	d.window = append(d.window, value)
}

func (d *AnomalyDetector) mean() float64 {
	sum := 0.0
	for _, v := range d.window {
		sum += v
	}
	return sum / float64(len(d.window))
}

func (d *AnomalyDetector) stdDev(mean float64) float64 {
	sumSquares := 0.0
	for _, v := range d.window {
		diff := v - mean
		sumSquares += diff * diff
	}
	return math.Sqrt(sumSquares / float64(len(d.window)))
}

// LatencyTracker tracks latency percentiles
type LatencyTracker struct {
	values  []float64
	maxSize int
	mu      sync.Mutex
}

// NewLatencyTracker creates a latency tracker
func NewLatencyTracker(maxSize int) *LatencyTracker {
	return &LatencyTracker{
		values:  make([]float64, 0, maxSize),
		maxSize: maxSize,
	}
}

// Record records a latency value in milliseconds
func (t *LatencyTracker) Record(latencyMs float64) {
	t.mu.Lock()
	defer t.mu.Unlock()
	
	if len(t.values) >= t.maxSize {
		t.values = t.values[1:]
	}
	t.values = append(t.values, latencyMs)
}

// Percentile returns the Nth percentile
func (t *LatencyTracker) Percentile(p float64) float64 {
	t.mu.Lock()
	defer t.mu.Unlock()
	
	if len(t.values) == 0 {
		return 0
	}

	// Sort a copy
	sorted := make([]float64, len(t.values))
	copy(sorted, t.values)
	
	// Simple insertion sort
	for i := 1; i < len(sorted); i++ {
		for j := i; j > 0 && sorted[j] < sorted[j-1]; j-- {
			sorted[j], sorted[j-1] = sorted[j-1], sorted[j]
		}
	}

	idx := int(float64(len(sorted)-1) * p / 100)
	return sorted[idx]
}

// P50, P90, P99 convenience methods
func (t *LatencyTracker) P50() float64 { return t.Percentile(50) }
func (t *LatencyTracker) P90() float64 { return t.Percentile(90) }
func (t *LatencyTracker) P99() float64 { return t.Percentile(99) }

// HealthChecker performs health checks
type HealthChecker struct {
	checks   map[string]HealthCheck
	results  map[string]HealthResult
	interval time.Duration
	mu       sync.RWMutex
}

// HealthCheck is a function that performs a health check
type HealthCheck func(ctx context.Context) error

// HealthResult represents the result of a health check
type HealthResult struct {
	Name      string
	Healthy   bool
	Message   string
	CheckedAt time.Time
	Duration  time.Duration
}

// NewHealthChecker creates a health checker
func NewHealthChecker(interval time.Duration) *HealthChecker {
	return &HealthChecker{
		checks:   make(map[string]HealthCheck),
		results:  make(map[string]HealthResult),
		interval: interval,
	}
}

// Register registers a health check
func (h *HealthChecker) Register(name string, check HealthCheck) {
	h.mu.Lock()
	defer h.mu.Unlock()
	h.checks[name] = check
}

// Start starts periodic health checks
func (h *HealthChecker) Start(ctx context.Context) {
	go func() {
		ticker := time.NewTicker(h.interval)
		defer ticker.Stop()
		
		for {
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
				h.runChecks(ctx)
			}
		}
	}()
}

func (h *HealthChecker) runChecks(ctx context.Context) {
	h.mu.RLock()
	checks := make(map[string]HealthCheck)
	for k, v := range h.checks {
		checks[k] = v
	}
	h.mu.RUnlock()

	results := make(map[string]HealthResult)
	for name, check := range checks {
		start := time.Now()
		err := check(ctx)
		duration := time.Since(start)
		
		result := HealthResult{
			Name:      name,
			Healthy:   err == nil,
			CheckedAt: start,
			Duration:  duration,
		}
		if err != nil {
			result.Message = err.Error()
		}
		results[name] = result
	}

	h.mu.Lock()
	h.results = results
	h.mu.Unlock()
}

// IsHealthy returns true if all checks pass
func (h *HealthChecker) IsHealthy() bool {
	h.mu.RLock()
	defer h.mu.RUnlock()
	
	for _, r := range h.results {
		if !r.Healthy {
			return false
		}
	}
	return true
}

// GetResults returns all health check results
func (h *HealthChecker) GetResults() map[string]HealthResult {
	h.mu.RLock()
	defer h.mu.RUnlock()
	
	results := make(map[string]HealthResult)
	for k, v := range h.results {
		results[k] = v
	}
	return results
}
