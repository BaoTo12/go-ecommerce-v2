package mq

import (
	"context"
	"encoding/json"
	"errors"
	"sync"
	"time"
)

/*
MESSAGE QUEUE ABSTRACTION

Provides a unified interface for message queues with:
- Multiple backend support (In-memory, Kafka, RabbitMQ, Redis)
- Dead letter queues
- Message acknowledgment
- Consumer groups
- Message retry with backoff
- Exactly-once semantics (idempotency)
*/

var (
	ErrQueueEmpty    = errors.New("queue empty")
	ErrQueueClosed   = errors.New("queue closed")
	ErrMessageNacked = errors.New("message nacked")
)

// Message represents a queue message
type Message struct {
	ID          string
	Topic       string
	Key         string
	Payload     []byte
	Headers     map[string]string
	Timestamp   time.Time
	Attempts    int
	MaxRetries  int
	RetryDelay  time.Duration
}

// Acknowledgment callback
type AckFunc func() error
type NackFunc func(requeue bool) error

// DeliveredMessage includes ack/nack callbacks
type DeliveredMessage struct {
	Message
	Ack  AckFunc
	Nack NackFunc
}

// Producer publishes messages
type Producer interface {
	Publish(ctx context.Context, msg *Message) error
	PublishBatch(ctx context.Context, msgs []*Message) error
	Close() error
}

// Consumer receives messages
type Consumer interface {
	Subscribe(topics []string) error
	Consume(ctx context.Context) (*DeliveredMessage, error)
	Close() error
}

// MessageHandler processes messages
type MessageHandler func(ctx context.Context, msg *Message) error

// MessageBroker is the unified interface
type MessageBroker interface {
	CreateProducer(config ProducerConfig) (Producer, error)
	CreateConsumer(config ConsumerConfig) (Consumer, error)
	CreateTopic(ctx context.Context, topic string, partitions int) error
	DeleteTopic(ctx context.Context, topic string) error
	Close() error
}

// ProducerConfig configures producers
type ProducerConfig struct {
	ClientID      string
	Acks          string // "none", "leader", "all"
	Idempotent    bool
	MaxRetries    int
	RetryInterval time.Duration
	BatchSize     int
	LingerMs      int
}

// ConsumerConfig configures consumers
type ConsumerConfig struct {
	GroupID           string
	Topics            []string
	AutoCommit        bool
	AutoCommitInterval time.Duration
	MaxPollRecords    int
	SessionTimeout    time.Duration
}

// InMemoryBroker is an in-memory implementation
type InMemoryBroker struct {
	topics    map[string]*topic
	consumers map[string][]*InMemoryConsumer
	mu        sync.RWMutex
	closed    bool
}

type topic struct {
	name       string
	messages   chan *Message
	dlq        chan *Message
	partitions int
}

// NewInMemoryBroker creates an in-memory broker
func NewInMemoryBroker() *InMemoryBroker {
	return &InMemoryBroker{
		topics:    make(map[string]*topic),
		consumers: make(map[string][]*InMemoryConsumer),
	}
}

func (b *InMemoryBroker) CreateTopic(ctx context.Context, name string, partitions int) error {
	b.mu.Lock()
	defer b.mu.Unlock()

	if _, exists := b.topics[name]; exists {
		return nil
	}

	b.topics[name] = &topic{
		name:       name,
		messages:   make(chan *Message, 10000),
		dlq:        make(chan *Message, 1000),
		partitions: partitions,
	}
	return nil
}

func (b *InMemoryBroker) DeleteTopic(ctx context.Context, name string) error {
	b.mu.Lock()
	defer b.mu.Unlock()
	delete(b.topics, name)
	return nil
}

func (b *InMemoryBroker) CreateProducer(config ProducerConfig) (Producer, error) {
	return &InMemoryProducer{broker: b, config: config}, nil
}

func (b *InMemoryBroker) CreateConsumer(config ConsumerConfig) (Consumer, error) {
	consumer := &InMemoryConsumer{
		broker: b,
		config: config,
		stopCh: make(chan struct{}),
	}

	b.mu.Lock()
	for _, topic := range config.Topics {
		b.consumers[topic] = append(b.consumers[topic], consumer)
	}
	b.mu.Unlock()

	return consumer, nil
}

func (b *InMemoryBroker) Close() error {
	b.mu.Lock()
	defer b.mu.Unlock()
	b.closed = true
	return nil
}

// InMemoryProducer produces to in-memory broker
type InMemoryProducer struct {
	broker *InMemoryBroker
	config ProducerConfig
}

func (p *InMemoryProducer) Publish(ctx context.Context, msg *Message) error {
	p.broker.mu.RLock()
	t, exists := p.broker.topics[msg.Topic]
	p.broker.mu.RUnlock()

	if !exists {
		p.broker.CreateTopic(ctx, msg.Topic, 1)
		p.broker.mu.RLock()
		t = p.broker.topics[msg.Topic]
		p.broker.mu.RUnlock()
	}

	if msg.ID == "" {
		msg.ID = generateMessageID()
	}
	msg.Timestamp = time.Now()

	select {
	case t.messages <- msg:
		return nil
	case <-ctx.Done():
		return ctx.Err()
	default:
		return errors.New("queue full")
	}
}

func (p *InMemoryProducer) PublishBatch(ctx context.Context, msgs []*Message) error {
	for _, msg := range msgs {
		if err := p.Publish(ctx, msg); err != nil {
			return err
		}
	}
	return nil
}

func (p *InMemoryProducer) Close() error {
	return nil
}

// InMemoryConsumer consumes from in-memory broker
type InMemoryConsumer struct {
	broker *InMemoryBroker
	config ConsumerConfig
	stopCh chan struct{}
	closed bool
}

func (c *InMemoryConsumer) Subscribe(topics []string) error {
	c.config.Topics = topics
	return nil
}

func (c *InMemoryConsumer) Consume(ctx context.Context) (*DeliveredMessage, error) {
	if c.closed {
		return nil, ErrQueueClosed
	}

	for _, topicName := range c.config.Topics {
		c.broker.mu.RLock()
		t, exists := c.broker.topics[topicName]
		c.broker.mu.RUnlock()

		if !exists {
			continue
		}

		select {
		case msg := <-t.messages:
			delivered := &DeliveredMessage{
				Message: *msg,
				Ack: func() error {
					return nil
				},
				Nack: func(requeue bool) error {
					if requeue && msg.Attempts < msg.MaxRetries {
						msg.Attempts++
						t.messages <- msg
					} else {
						t.dlq <- msg
					}
					return nil
				},
			}
			return delivered, nil
		case <-ctx.Done():
			return nil, ctx.Err()
		case <-c.stopCh:
			return nil, ErrQueueClosed
		default:
			continue
		}
	}

	// No messages, wait a bit
	select {
	case <-time.After(100 * time.Millisecond):
		return nil, ErrQueueEmpty
	case <-ctx.Done():
		return nil, ctx.Err()
	}
}

func (c *InMemoryConsumer) Close() error {
	c.closed = true
	close(c.stopCh)
	return nil
}

// ConsumerWorker manages message consumption with handlers
type ConsumerWorker struct {
	consumer    Consumer
	handler     MessageHandler
	concurrency int
	wg          sync.WaitGroup
	stopCh      chan struct{}
}

// NewConsumerWorker creates a consumer worker
func NewConsumerWorker(consumer Consumer, handler MessageHandler, concurrency int) *ConsumerWorker {
	return &ConsumerWorker{
		consumer:    consumer,
		handler:     handler,
		concurrency: concurrency,
		stopCh:      make(chan struct{}),
	}
}

// Start starts the worker
func (w *ConsumerWorker) Start(ctx context.Context) {
	for i := 0; i < w.concurrency; i++ {
		w.wg.Add(1)
		go w.work(ctx)
	}
}

// Stop stops the worker
func (w *ConsumerWorker) Stop() {
	close(w.stopCh)
	w.wg.Wait()
	w.consumer.Close()
}

func (w *ConsumerWorker) work(ctx context.Context) {
	defer w.wg.Done()

	for {
		select {
		case <-ctx.Done():
			return
		case <-w.stopCh:
			return
		default:
			msg, err := w.consumer.Consume(ctx)
			if err != nil {
				if err == ErrQueueEmpty {
					continue
				}
				return
			}

			if err := w.handler(ctx, &msg.Message); err != nil {
				msg.Nack(true)
			} else {
				msg.Ack()
			}
		}
	}
}

// IdempotencyStore prevents duplicate processing
type IdempotencyStore interface {
	CheckAndSet(ctx context.Context, key string, ttl time.Duration) (bool, error)
}

// InMemoryIdempotencyStore is an in-memory implementation
type InMemoryIdempotencyStore struct {
	keys map[string]time.Time
	mu   sync.Mutex
}

func NewInMemoryIdempotencyStore() *InMemoryIdempotencyStore {
	return &InMemoryIdempotencyStore{
		keys: make(map[string]time.Time),
	}
}

func (s *InMemoryIdempotencyStore) CheckAndSet(ctx context.Context, key string, ttl time.Duration) (bool, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if expires, exists := s.keys[key]; exists {
		if time.Now().Before(expires) {
			return false, nil // Already processed
		}
	}

	s.keys[key] = time.Now().Add(ttl)
	return true, nil
}

// IdempotentHandler wraps a handler with idempotency
func IdempotentHandler(store IdempotencyStore, handler MessageHandler, ttl time.Duration) MessageHandler {
	return func(ctx context.Context, msg *Message) error {
		isNew, err := store.CheckAndSet(ctx, msg.ID, ttl)
		if err != nil {
			return err
		}
		if !isNew {
			return nil // Already processed
		}
		return handler(ctx, msg)
	}
}

// Helper to publish JSON messages
func PublishJSON(ctx context.Context, producer Producer, topic, key string, payload interface{}) error {
	data, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	return producer.Publish(ctx, &Message{
		Topic:   topic,
		Key:     key,
		Payload: data,
		Headers: map[string]string{
			"content-type": "application/json",
		},
	})
}

func generateMessageID() string {
	return "msg_" + time.Now().Format("20060102150405.000000")
}
