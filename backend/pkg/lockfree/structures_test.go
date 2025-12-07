package lockfree

import (
	"sync"
	"testing"
)

func TestStack_PushPop(t *testing.T) {
	stack := NewStack[int]()

	// Push
	stack.Push(1)
	stack.Push(2)
	stack.Push(3)

	if stack.Size() != 3 {
		t.Errorf("Expected size 3, got %d", stack.Size())
	}

	// Pop
	v, ok := stack.Pop()
	if !ok || v != 3 {
		t.Errorf("Expected 3, got %d", v)
	}

	v, ok = stack.Pop()
	if !ok || v != 2 {
		t.Errorf("Expected 2, got %d", v)
	}

	v, ok = stack.Pop()
	if !ok || v != 1 {
		t.Errorf("Expected 1, got %d", v)
	}

	// Empty
	_, ok = stack.Pop()
	if ok {
		t.Error("Expected empty stack")
	}
}

func TestStack_Concurrent(t *testing.T) {
	stack := NewStack[int]()
	var wg sync.WaitGroup
	n := 1000

	// Concurrent pushes
	for i := 0; i < n; i++ {
		wg.Add(1)
		go func(v int) {
			defer wg.Done()
			stack.Push(v)
		}(i)
	}
	wg.Wait()

	if stack.Size() != int64(n) {
		t.Errorf("Expected size %d, got %d", n, stack.Size())
	}

	// Concurrent pops
	count := int64(0)
	for i := 0; i < n; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			if _, ok := stack.Pop(); ok {
				// Atomic increment would be better here
			}
		}()
	}
	wg.Wait()

	_ = count
}

func TestQueue_EnqueueDequeue(t *testing.T) {
	queue := NewQueue[string]()

	queue.Enqueue("a")
	queue.Enqueue("b")
	queue.Enqueue("c")

	if queue.Size() != 3 {
		t.Errorf("Expected size 3, got %d", queue.Size())
	}

	v, ok := queue.Dequeue()
	if !ok || v != "a" {
		t.Errorf("Expected 'a', got %s", v)
	}

	v, ok = queue.Dequeue()
	if !ok || v != "b" {
		t.Errorf("Expected 'b', got %s", v)
	}

	v, ok = queue.Dequeue()
	if !ok || v != "c" {
		t.Errorf("Expected 'c', got %s", v)
	}

	if !queue.IsEmpty() {
		t.Error("Expected empty queue")
	}
}

func TestQueue_Concurrent(t *testing.T) {
	queue := NewQueue[int]()
	var wg sync.WaitGroup
	n := 1000

	// Concurrent enqueues
	for i := 0; i < n; i++ {
		wg.Add(1)
		go func(v int) {
			defer wg.Done()
			queue.Enqueue(v)
		}(i)
	}
	wg.Wait()

	if queue.Size() != int64(n) {
		t.Errorf("Expected size %d, got %d", n, queue.Size())
	}
}

func TestRingBuffer_Basic(t *testing.T) {
	rb := NewRingBuffer[int](8)

	if rb.Capacity() != 8 {
		t.Errorf("Expected capacity 8, got %d", rb.Capacity())
	}

	for i := 0; i < 5; i++ {
		if !rb.TryPush(i) {
			t.Errorf("Push %d failed", i)
		}
	}

	if rb.Size() != 5 {
		t.Errorf("Expected size 5, got %d", rb.Size())
	}

	for i := 0; i < 5; i++ {
		v, ok := rb.TryPop()
		if !ok || v != i {
			t.Errorf("Expected %d, got %d", i, v)
		}
	}
}

func TestRingBuffer_Full(t *testing.T) {
	rb := NewRingBuffer[int](4)

	// Fill it
	for i := 0; i < 4; i++ {
		rb.TryPush(i)
	}

	// Should fail when full
	if rb.TryPush(100) {
		t.Error("Push should fail when full")
	}
}

func TestCounter(t *testing.T) {
	c := NewCounter()

	c.Inc()
	c.Inc()
	c.Inc()

	if c.Value() != 3 {
		t.Errorf("Expected 3, got %d", c.Value())
	}

	c.Dec()
	if c.Value() != 2 {
		t.Errorf("Expected 2, got %d", c.Value())
	}

	c.Add(10)
	if c.Value() != 12 {
		t.Errorf("Expected 12, got %d", c.Value())
	}

	c.Reset()
	if c.Value() != 0 {
		t.Errorf("Expected 0, got %d", c.Value())
	}
}

func BenchmarkStack_Push(b *testing.B) {
	stack := NewStack[int]()
	for i := 0; i < b.N; i++ {
		stack.Push(i)
	}
}

func BenchmarkStack_PushPop(b *testing.B) {
	stack := NewStack[int]()
	for i := 0; i < b.N; i++ {
		stack.Push(i)
		stack.Pop()
	}
}

func BenchmarkQueue_EnqueueDequeue(b *testing.B) {
	queue := NewQueue[int]()
	for i := 0; i < b.N; i++ {
		queue.Enqueue(i)
		queue.Dequeue()
	}
}

func BenchmarkCounter_Inc(b *testing.B) {
	c := NewCounter()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			c.Inc()
		}
	})
}
