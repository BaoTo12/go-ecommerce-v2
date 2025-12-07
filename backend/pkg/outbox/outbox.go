package outbox

import (
	"context"
	"encoding/json"
	"errors"
	"sync"
	"time"
)

/*
OUTBOX PATTERN - Reliable Messaging

Implements the Transactional Outbox pattern for reliable event publishing.
Events are stored in an outbox table within the same transaction as the
business data, then published asynchronously by a separate process.

This guarantees at-least-once delivery of events.
*/

// OutboxMessage represents a message in the outbox
type OutboxMessage struct {
	ID          string
	AggregateID string
	EventType   string
	Payload     []byte
	CreatedAt   time.Time
	ProcessedAt *time.Time
	Attempts    int
	LastError   string
	Status      MessageStatus
}

// MessageStatus represents the status of an outbox message
type MessageStatus string

const (
	StatusPending    MessageStatus = "PENDING"
	StatusProcessing MessageStatus = "PROCESSING"
	StatusCompleted  MessageStatus = "COMPLETED"
	StatusFailed     MessageStatus = "FAILED"
)

// OutboxStore persists outbox messages
type OutboxStore interface {
	// Save stores a new message (within the same transaction as business data)
	Save(ctx context.Context, msg *OutboxMessage) error
	
	// GetPending retrieves pending messages for processing
	GetPending(ctx context.Context, limit int) ([]*OutboxMessage, error)
	
	// MarkProcessing marks a message as being processed
	MarkProcessing(ctx context.Context, id string) error
	
	// MarkCompleted marks a message as successfully processed
	MarkCompleted(ctx context.Context, id string) error
	
	// MarkFailed marks a message as failed
	MarkFailed(ctx context.Context, id string, err string) error
}

// Publisher publishes messages to the message broker
type Publisher interface {
	Publish(ctx context.Context, topic string, key string, payload []byte) error
}

// OutboxProcessor processes outbox messages
type OutboxProcessor struct {
	store      OutboxStore
	publisher  Publisher
	batchSize  int
	pollInterval time.Duration
	maxRetries int
	
	running bool
	stopCh  chan struct{}
	wg      sync.WaitGroup
	mu      sync.Mutex
}

// ProcessorConfig configures the outbox processor
type ProcessorConfig struct {
	BatchSize    int
	PollInterval time.Duration
	MaxRetries   int
}

// DefaultProcessorConfig returns default configuration
func DefaultProcessorConfig() ProcessorConfig {
	return ProcessorConfig{
		BatchSize:    100,
		PollInterval: 1 * time.Second,
		MaxRetries:   5,
	}
}

// NewOutboxProcessor creates a new outbox processor
func NewOutboxProcessor(store OutboxStore, publisher Publisher, config ProcessorConfig) *OutboxProcessor {
	return &OutboxProcessor{
		store:        store,
		publisher:    publisher,
		batchSize:    config.BatchSize,
		pollInterval: config.PollInterval,
		maxRetries:   config.MaxRetries,
		stopCh:       make(chan struct{}),
	}
}

// Start starts the outbox processor
func (p *OutboxProcessor) Start(ctx context.Context) error {
	p.mu.Lock()
	if p.running {
		p.mu.Unlock()
		return errors.New("processor already running")
	}
	p.running = true
	p.mu.Unlock()

	p.wg.Add(1)
	go p.processLoop(ctx)

	return nil
}

// Stop stops the outbox processor
func (p *OutboxProcessor) Stop() {
	p.mu.Lock()
	if !p.running {
		p.mu.Unlock()
		return
	}
	p.running = false
	p.mu.Unlock()

	close(p.stopCh)
	p.wg.Wait()
}

func (p *OutboxProcessor) processLoop(ctx context.Context) {
	defer p.wg.Done()

	ticker := time.NewTicker(p.pollInterval)
	defer ticker.Stop()

	for {
		select {
		case <-p.stopCh:
			return
		case <-ctx.Done():
			return
		case <-ticker.C:
			p.processBatch(ctx)
		}
	}
}

func (p *OutboxProcessor) processBatch(ctx context.Context) {
	messages, err := p.store.GetPending(ctx, p.batchSize)
	if err != nil {
		return
	}

	for _, msg := range messages {
		if err := p.processMessage(ctx, msg); err != nil {
			// Message will be retried on next poll
			continue
		}
	}
}

func (p *OutboxProcessor) processMessage(ctx context.Context, msg *OutboxMessage) error {
	// Mark as processing
	if err := p.store.MarkProcessing(ctx, msg.ID); err != nil {
		return err
	}

	// Check max retries
	if msg.Attempts >= p.maxRetries {
		p.store.MarkFailed(ctx, msg.ID, "max retries exceeded")
		return errors.New("max retries exceeded")
	}

	// Publish to message broker
	err := p.publisher.Publish(ctx, msg.EventType, msg.AggregateID, msg.Payload)
	if err != nil {
		p.store.MarkFailed(ctx, msg.ID, err.Error())
		return err
	}

	// Mark as completed
	return p.store.MarkCompleted(ctx, msg.ID)
}

// InMemoryOutboxStore is an in-memory implementation
type InMemoryOutboxStore struct {
	messages map[string]*OutboxMessage
	mu       sync.RWMutex
}

func NewInMemoryOutboxStore() *InMemoryOutboxStore {
	return &InMemoryOutboxStore{
		messages: make(map[string]*OutboxMessage),
	}
}

func (s *InMemoryOutboxStore) Save(ctx context.Context, msg *OutboxMessage) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	
	msg.ID = generateOutboxID()
	msg.CreatedAt = time.Now()
	msg.Status = StatusPending
	s.messages[msg.ID] = msg
	return nil
}

func (s *InMemoryOutboxStore) GetPending(ctx context.Context, limit int) ([]*OutboxMessage, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	
	pending := make([]*OutboxMessage, 0, limit)
	for _, msg := range s.messages {
		if msg.Status == StatusPending || msg.Status == StatusFailed {
			pending = append(pending, msg)
			if len(pending) >= limit {
				break
			}
		}
	}
	return pending, nil
}

func (s *InMemoryOutboxStore) MarkProcessing(ctx context.Context, id string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	
	msg, exists := s.messages[id]
	if !exists {
		return errors.New("message not found")
	}
	msg.Status = StatusProcessing
	msg.Attempts++
	return nil
}

func (s *InMemoryOutboxStore) MarkCompleted(ctx context.Context, id string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	
	msg, exists := s.messages[id]
	if !exists {
		return errors.New("message not found")
	}
	msg.Status = StatusCompleted
	now := time.Now()
	msg.ProcessedAt = &now
	return nil
}

func (s *InMemoryOutboxStore) MarkFailed(ctx context.Context, id string, errMsg string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	
	msg, exists := s.messages[id]
	if !exists {
		return errors.New("message not found")
	}
	msg.Status = StatusFailed
	msg.LastError = errMsg
	return nil
}

func generateOutboxID() string {
	return "outbox_" + time.Now().Format("20060102150405.000000")
}

// Helper to create outbox message from event
func NewOutboxMessage(aggregateID, eventType string, payload interface{}) (*OutboxMessage, error) {
	data, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}
	
	return &OutboxMessage{
		AggregateID: aggregateID,
		EventType:   eventType,
		Payload:     data,
	}, nil
}

// MockPublisher for testing
type MockPublisher struct {
	Published []struct {
		Topic   string
		Key     string
		Payload []byte
	}
	mu sync.Mutex
}

func (m *MockPublisher) Publish(ctx context.Context, topic, key string, payload []byte) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	
	m.Published = append(m.Published, struct {
		Topic   string
		Key     string
		Payload []byte
	}{topic, key, payload})
	return nil
}
