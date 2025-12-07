package bloom

import (
	"hash"
	"hash/fnv"
	"math"
	"sync"
)

// Filter is a space-efficient probabilistic data structure
// for membership testing with false positive possibility but no false negatives
type Filter struct {
	bits     []bool
	size     uint
	hashFns  int
	mu       sync.RWMutex
	count    uint64
}

// NewFilter creates a new Bloom filter
// n = expected number of elements
// p = false positive probability (0.01 = 1%)
func NewFilter(n uint, p float64) *Filter {
	// Calculate optimal size: m = -n*ln(p)/(ln(2)^2)
	ln2Sq := math.Ln2 * math.Ln2
	m := uint(math.Ceil(-float64(n) * math.Log(p) / ln2Sq))

	// Calculate optimal hash functions: k = (m/n)*ln(2)
	k := int(math.Ceil(float64(m) / float64(n) * math.Ln2))

	return &Filter{
		bits:    make([]bool, m),
		size:    m,
		hashFns: k,
	}
}

// Add adds an element to the filter
func (f *Filter) Add(data []byte) {
	f.mu.Lock()
	defer f.mu.Unlock()

	for i := 0; i < f.hashFns; i++ {
		idx := f.hash(data, i) % f.size
		f.bits[idx] = true
	}
	f.count++
}

// AddString adds a string element
func (f *Filter) AddString(s string) {
	f.Add([]byte(s))
}

// Contains tests if an element might be in the set
// Returns true if element might exist (could be false positive)
// Returns false if element definitely doesn't exist
func (f *Filter) Contains(data []byte) bool {
	f.mu.RLock()
	defer f.mu.RUnlock()

	for i := 0; i < f.hashFns; i++ {
		idx := f.hash(data, i) % f.size
		if !f.bits[idx] {
			return false
		}
	}
	return true
}

// ContainsString tests if a string element might be in the set
func (f *Filter) ContainsString(s string) bool {
	return f.Contains([]byte(s))
}

// hash generates hash values using FNV with seed variation
func (f *Filter) hash(data []byte, seed int) uint {
	h := fnv.New64a()
	h.Write(data)
	h.Write([]byte{byte(seed)})
	return uint(h.Sum64())
}

// Count returns approximate number of elements added
func (f *Filter) Count() uint64 {
	f.mu.RLock()
	defer f.mu.RUnlock()
	return f.count
}

// EstimateFalsePositiveRate returns current false positive probability
func (f *Filter) EstimateFalsePositiveRate() float64 {
	f.mu.RLock()
	setBits := 0
	for _, bit := range f.bits {
		if bit {
			setBits++
		}
	}
	f.mu.RUnlock()

	if setBits == 0 {
		return 0
	}

	// p = (1 - e^(-kn/m))^k where n is count, m is size, k is hashFns
	fillRatio := float64(setBits) / float64(f.size)
	return math.Pow(fillRatio, float64(f.hashFns))
}

// Reset clears the filter
func (f *Filter) Reset() {
	f.mu.Lock()
	defer f.mu.Unlock()
	f.bits = make([]bool, f.size)
	f.count = 0
}

// CountingFilter allows deletions with counters instead of bits
type CountingFilter struct {
	counters []uint8
	size     uint
	hashFns  int
	mu       sync.RWMutex
}

// NewCountingFilter creates a counting bloom filter
func NewCountingFilter(n uint, p float64) *CountingFilter {
	ln2Sq := math.Ln2 * math.Ln2
	m := uint(math.Ceil(-float64(n) * math.Log(p) / ln2Sq))
	k := int(math.Ceil(float64(m) / float64(n) * math.Ln2))

	return &CountingFilter{
		counters: make([]uint8, m),
		size:     m,
		hashFns:  k,
	}
}

// Add increments counters for an element
func (f *CountingFilter) Add(data []byte) {
	f.mu.Lock()
	defer f.mu.Unlock()

	for i := 0; i < f.hashFns; i++ {
		idx := f.hash(data, i) % f.size
		if f.counters[idx] < 255 {
			f.counters[idx]++
		}
	}
}

// Remove decrements counters for an element
func (f *CountingFilter) Remove(data []byte) {
	f.mu.Lock()
	defer f.mu.Unlock()

	for i := 0; i < f.hashFns; i++ {
		idx := f.hash(data, i) % f.size
		if f.counters[idx] > 0 {
			f.counters[idx]--
		}
	}
}

// Contains tests membership
func (f *CountingFilter) Contains(data []byte) bool {
	f.mu.RLock()
	defer f.mu.RUnlock()

	for i := 0; i < f.hashFns; i++ {
		idx := f.hash(data, i) % f.size
		if f.counters[idx] == 0 {
			return false
		}
	}
	return true
}

func (f *CountingFilter) hash(data []byte, seed int) uint {
	h := fnv.New64a()
	h.Write(data)
	h.Write([]byte{byte(seed)})
	return uint(h.Sum64())
}

// ScalableFilter auto-scales by adding new filters
type ScalableFilter struct {
	filters  []*Filter
	n        uint
	p        float64
	growth   float64 // Growth rate for new filters
	mu       sync.RWMutex
}

// NewScalableFilter creates an auto-scaling bloom filter
func NewScalableFilter(initialN uint, p float64) *ScalableFilter {
	sf := &ScalableFilter{
		filters: make([]*Filter, 0),
		n:       initialN,
		p:       p,
		growth:  2.0, // Double size for each new filter
	}
	sf.addFilter()
	return sf
}

func (sf *ScalableFilter) addFilter() {
	n := sf.n
	if len(sf.filters) > 0 {
		n = uint(float64(n) * math.Pow(sf.growth, float64(len(sf.filters))))
	}
	// Reduce p for newer filters to maintain overall false positive rate
	p := sf.p * math.Pow(0.5, float64(len(sf.filters)+1))
	sf.filters = append(sf.filters, NewFilter(n, p))
}

// Add adds to the current filter, creating new one if needed
func (sf *ScalableFilter) Add(data []byte) {
	sf.mu.Lock()
	defer sf.mu.Unlock()

	currentFilter := sf.filters[len(sf.filters)-1]
	
	// Check if current filter is too full
	if currentFilter.EstimateFalsePositiveRate() > sf.p {
		sf.addFilter()
		currentFilter = sf.filters[len(sf.filters)-1]
	}

	currentFilter.Add(data)
}

// Contains checks all filters
func (sf *ScalableFilter) Contains(data []byte) bool {
	sf.mu.RLock()
	defer sf.mu.RUnlock()

	for _, f := range sf.filters {
		if f.Contains(data) {
			return true
		}
	}
	return false
}

// Hasher interface for custom hash functions
type Hasher interface {
	Hash(data []byte) uint64
}

// FNVHasher uses FNV-1a
type FNVHasher struct {
	h hash.Hash64
}

func NewFNVHasher() *FNVHasher {
	return &FNVHasher{h: fnv.New64a()}
}

func (f *FNVHasher) Hash(data []byte) uint64 {
	f.h.Reset()
	f.h.Write(data)
	return f.h.Sum64()
}
