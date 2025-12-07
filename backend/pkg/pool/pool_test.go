package pool_test

import (
	"context"
	"errors"
	"sync/atomic"
	"testing"
	"time"

	"github.com/titan-commerce/backend/pkg/pool"
)

func TestWorkerPool_Submit(t *testing.T) {
	cfg := pool.WorkerPoolConfig{
		Workers:   4,
		QueueSize: 10,
	}
	wp := pool.NewWorkerPool(cfg)
	defer wp.Shutdown(5 * time.Second)

	var counter int32
	for i := 0; i < 10; i++ {
		err := wp.Submit(func(ctx context.Context) error {
			atomic.AddInt32(&counter, 1)
			return nil
		})
		if err != nil {
			t.Errorf("Submit failed: %v", err)
		}
	}

	// Wait for tasks to complete
	time.Sleep(100 * time.Millisecond)

	if atomic.LoadInt32(&counter) != 10 {
		t.Errorf("Expected 10 tasks processed, got %d", counter)
	}
}

func TestWorkerPool_SubmitWait(t *testing.T) {
	cfg := pool.DefaultWorkerConfig()
	wp := pool.NewWorkerPool(cfg)
	defer wp.Shutdown(5 * time.Second)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	result := 0
	err := wp.SubmitWait(ctx, func(ctx context.Context) error {
		result = 42
		return nil
	})

	if err != nil {
		t.Errorf("SubmitWait failed: %v", err)
	}
	if result != 42 {
		t.Errorf("Expected result 42, got %d", result)
	}
}

func TestWorkerPool_Stats(t *testing.T) {
	cfg := pool.WorkerPoolConfig{
		Workers:   5,
		QueueSize: 20,
	}
	wp := pool.NewWorkerPool(cfg)
	defer wp.Shutdown(5 * time.Second)

	// Submit some tasks
	for i := 0; i < 5; i++ {
		wp.Submit(func(ctx context.Context) error {
			time.Sleep(10 * time.Millisecond)
			return nil
		})
	}

	stats := wp.Stats()
	if stats.Workers != 5 {
		t.Errorf("Expected 5 workers, got %d", stats.Workers)
	}
	if stats.QueueCapacity != 20 {
		t.Errorf("Expected queue capacity 20, got %d", stats.QueueCapacity)
	}
}

func TestWorkerPool_ErrorHandling(t *testing.T) {
	cfg := pool.DefaultWorkerConfig()
	wp := pool.NewWorkerPool(cfg)
	defer wp.Shutdown(5 * time.Second)

	expectedErr := errors.New("test error")

	// Submit task that returns error
	wp.Submit(func(ctx context.Context) error {
		return expectedErr
	})

	time.Sleep(50 * time.Millisecond)

	stats := wp.Stats()
	if stats.Errors == 0 {
		t.Error("Expected errors to be tracked")
	}
}

func TestObjectPool_Generic(t *testing.T) {
	type TestObject struct {
		Value int
	}

	objPool := pool.NewObjectPool(func() *TestObject {
		return &TestObject{Value: 0}
	})

	// Get object
	obj := objPool.Get()
	if obj == nil {
		t.Fatal("Expected non-nil object")
	}

	obj.Value = 42
	objPool.Put(obj)

	// Object should be reused
	obj2 := objPool.Get()
	if obj2.Value != 42 {
		// Note: sync.Pool may not return the same object, this is expected
		// Just verify we get a valid object
		if obj2 == nil {
			t.Error("Expected non-nil object from pool")
		}
	}
}

func TestBufferPool(t *testing.T) {
	bufPool := pool.NewBufferPool(1024)

	buf := bufPool.Get()
	if buf == nil {
		t.Fatal("Expected non-nil buffer")
	}

	buf.WriteString("Hello, World!")
	if buf.String() != "Hello, World!" {
		t.Errorf("Expected 'Hello, World!', got %s", buf.String())
	}

	bufPool.Put(buf)

	// Get another buffer - should be reset
	buf2 := bufPool.Get()
	if buf2.Len() != 0 {
		t.Errorf("Expected empty buffer, got length %d", buf2.Len())
	}
}

func TestSlicePool(t *testing.T) {
	slicePool := pool.NewSlicePool[int](10)

	s := slicePool.Get()
	if s == nil {
		t.Fatal("Expected non-nil slice")
	}
	if len(s) != 0 {
		t.Errorf("Expected empty slice, got length %d", len(s))
	}
	if cap(s) < 10 {
		t.Errorf("Expected capacity >= 10, got %d", cap(s))
	}

	// Use the slice
	s = append(s, 1, 2, 3)
	slicePool.Put(s)
}

func BenchmarkWorkerPool_Submit(b *testing.B) {
	cfg := pool.WorkerPoolConfig{
		Workers:   10,
		QueueSize: 1000,
	}
	wp := pool.NewWorkerPool(cfg)
	defer wp.Shutdown(5 * time.Second)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		wp.Submit(func(ctx context.Context) error {
			return nil
		})
	}
}

func BenchmarkBufferPool(b *testing.B) {
	bufPool := pool.NewBufferPool(4096)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		buf := bufPool.Get()
		buf.WriteString("Hello, World!")
		bufPool.Put(buf)
	}
}
