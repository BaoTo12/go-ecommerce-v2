package bloom_test

import (
	"testing"

	"github.com/titan-commerce/backend/pkg/bloom"
)

// ====================
// PHASE 1: UNIT TESTS
// ====================

func TestBloomFilter_AddAndContains(t *testing.T) {
	bf := bloom.NewFilter(1000, 0.01)

	// Add elements
	bf.AddString("hello")
	bf.AddString("world")
	bf.AddString("test")

	// Test positive cases (should always return true)
	if !bf.ContainsString("hello") {
		t.Error("Expected 'hello' to be in filter")
	}
	if !bf.ContainsString("world") {
		t.Error("Expected 'world' to be in filter")
	}
	if !bf.ContainsString("test") {
		t.Error("Expected 'test' to be in filter")
	}

	// Test negative cases (should usually return false)
	// Note: May have false positives, but with 1% FP rate they're rare
	notInFilter := []string{"foo", "bar", "baz", "qux", "nope"}
	falsePositives := 0
	for _, s := range notInFilter {
		if bf.ContainsString(s) {
			falsePositives++
		}
	}
	// Allow up to 1 false positive out of 5 (20%, very generous)
	if falsePositives > 1 {
		t.Errorf("Too many false positives: %d", falsePositives)
	}
}

func TestBloomFilter_Count(t *testing.T) {
	bf := bloom.NewFilter(100, 0.01)

	if bf.Count() != 0 {
		t.Errorf("Expected count 0, got %d", bf.Count())
	}

	bf.AddString("item1")
	bf.AddString("item2")
	bf.AddString("item3")

	if bf.Count() != 3 {
		t.Errorf("Expected count 3, got %d", bf.Count())
	}
}

func TestBloomFilter_Reset(t *testing.T) {
	bf := bloom.NewFilter(100, 0.01)

	bf.AddString("test")
	if !bf.ContainsString("test") {
		t.Error("Expected 'test' to be in filter before reset")
	}

	bf.Reset()

	if bf.ContainsString("test") {
		t.Error("Expected 'test' to NOT be in filter after reset")
	}
	if bf.Count() != 0 {
		t.Error("Expected count 0 after reset")
	}
}

func TestCountingFilter_AddAndRemove(t *testing.T) {
	cf := bloom.NewCountingFilter(1000, 0.01)

	cf.Add([]byte("test"))
	if !cf.Contains([]byte("test")) {
		t.Error("Expected 'test' to be in filter")
	}

	cf.Remove([]byte("test"))
	if cf.Contains([]byte("test")) {
		t.Error("Expected 'test' to NOT be in filter after remove")
	}
}

func TestScalableFilter_AutoScale(t *testing.T) {
	sf := bloom.NewScalableFilter(10, 0.01)

	// Add more elements than initial capacity
	for i := 0; i < 100; i++ {
		sf.Add([]byte{byte(i)})
	}

	// Verify all elements are found
	for i := 0; i < 100; i++ {
		if !sf.Contains([]byte{byte(i)}) {
			t.Errorf("Expected element %d to be in filter", i)
		}
	}
}

// ====================
// PHASE 2: BENCHMARK TESTS
// ====================

func BenchmarkBloomFilter_Add(b *testing.B) {
	bf := bloom.NewFilter(1000000, 0.01)
	data := []byte("benchmark-test-data")

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		bf.Add(data)
	}
}

func BenchmarkBloomFilter_Contains(b *testing.B) {
	bf := bloom.NewFilter(1000000, 0.01)
	data := []byte("benchmark-test-data")
	bf.Add(data)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		bf.Contains(data)
	}
}

func BenchmarkBloomFilter_Parallel(b *testing.B) {
	bf := bloom.NewFilter(1000000, 0.01)

	b.RunParallel(func(pb *testing.PB) {
		i := 0
		for pb.Next() {
			bf.Add([]byte{byte(i % 256)})
			bf.Contains([]byte{byte((i + 1) % 256)})
			i++
		}
	})
}

// ====================
// PHASE 3: FUZZ TESTS
// ====================

func FuzzBloomFilter_AddContains(f *testing.F) {
	// Seed corpus
	f.Add([]byte("hello"))
	f.Add([]byte("world"))
	f.Add([]byte(""))
	f.Add([]byte{0, 1, 2, 3, 4, 5})

	bf := bloom.NewFilter(10000, 0.01)

	f.Fuzz(func(t *testing.T, data []byte) {
		// Add element
		bf.Add(data)

		// Must be found (no false negatives)
		if !bf.Contains(data) {
			t.Errorf("Added element not found: %v", data)
		}
	})
}

func FuzzCountingFilter(f *testing.F) {
	f.Add([]byte("test"))
	f.Add([]byte{255, 255, 255})

	cf := bloom.NewCountingFilter(1000, 0.01)

	f.Fuzz(func(t *testing.T, data []byte) {
		cf.Add(data)
		if !cf.Contains(data) {
			t.Error("Element not found after add")
		}
		cf.Remove(data)
	})
}
