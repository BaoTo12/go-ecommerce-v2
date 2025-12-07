package cqrs

import (
	"context"
	"sync"
	"time"

	"github.com/google/uuid"
)

// Event represents a domain event
type Event struct {
	ID            string                 `json:"id"`
	AggregateID   string                 `json:"aggregate_id"`
	AggregateType string                 `json:"aggregate_type"`
	EventType     string                 `json:"event_type"`
	Version       int                    `json:"version"`
	Payload       map[string]interface{} `json:"payload"`
	Metadata      map[string]string      `json:"metadata"`
	Timestamp     time.Time              `json:"timestamp"`
}

// NewEvent creates a new event
func NewEvent(aggregateID, aggregateType, eventType string, payload map[string]interface{}) *Event {
	return &Event{
		ID:            uuid.New().String(),
		AggregateID:   aggregateID,
		AggregateType: aggregateType,
		EventType:     eventType,
		Payload:       payload,
		Metadata:      make(map[string]string),
		Timestamp:     time.Now(),
	}
}

// EventStore persists events
type EventStore interface {
	Save(ctx context.Context, events ...*Event) error
	Load(ctx context.Context, aggregateID string) ([]*Event, error)
	LoadFrom(ctx context.Context, aggregateID string, version int) ([]*Event, error)
}

// EventHandler processes events
type EventHandler func(ctx context.Context, event *Event) error

// EventBus distributes events to handlers
type EventBus struct {
	handlers map[string][]EventHandler
	mu       sync.RWMutex
}

// NewEventBus creates a new event bus
func NewEventBus() *EventBus {
	return &EventBus{
		handlers: make(map[string][]EventHandler),
	}
}

// Subscribe registers a handler for an event type
func (b *EventBus) Subscribe(eventType string, handler EventHandler) {
	b.mu.Lock()
	defer b.mu.Unlock()
	b.handlers[eventType] = append(b.handlers[eventType], handler)
}

// Publish sends an event to all subscribed handlers
func (b *EventBus) Publish(ctx context.Context, event *Event) error {
	b.mu.RLock()
	handlers := b.handlers[event.EventType]
	b.mu.RUnlock()

	for _, handler := range handlers {
		if err := handler(ctx, event); err != nil {
			return err
		}
	}
	return nil
}

// PublishAsync sends events asynchronously
func (b *EventBus) PublishAsync(ctx context.Context, event *Event) {
	b.mu.RLock()
	handlers := b.handlers[event.EventType]
	b.mu.RUnlock()

	for _, handler := range handlers {
		go handler(ctx, event)
	}
}

// Command represents a command to execute
type Command interface {
	CommandName() string
	AggregateID() string
}

// CommandHandler processes commands
type CommandHandler interface {
	Handle(ctx context.Context, cmd Command) error
}

// CommandBus routes commands to handlers
type CommandBus struct {
	handlers map[string]CommandHandler
	mu       sync.RWMutex
}

// NewCommandBus creates a new command bus
func NewCommandBus() *CommandBus {
	return &CommandBus{
		handlers: make(map[string]CommandHandler),
	}
}

// Register adds a handler for a command type
func (b *CommandBus) Register(commandName string, handler CommandHandler) {
	b.mu.Lock()
	defer b.mu.Unlock()
	b.handlers[commandName] = handler
}

// Dispatch sends a command to its handler
func (b *CommandBus) Dispatch(ctx context.Context, cmd Command) error {
	b.mu.RLock()
	handler, ok := b.handlers[cmd.CommandName()]
	b.mu.RUnlock()

	if !ok {
		return ErrHandlerNotFound
	}
	return handler.Handle(ctx, cmd)
}

// Aggregate is the base for event-sourced aggregates
type Aggregate struct {
	ID      string
	Version int
	Changes []*Event
}

// Apply records an event
func (a *Aggregate) Apply(event *Event) {
	event.Version = a.Version + 1
	a.Version = event.Version
	a.Changes = append(a.Changes, event)
}

// ClearChanges removes uncommitted events
func (a *Aggregate) ClearChanges() {
	a.Changes = nil
}

// InMemoryEventStore is a simple in-memory event store
type InMemoryEventStore struct {
	events map[string][]*Event
	mu     sync.RWMutex
}

// NewInMemoryEventStore creates a new in-memory event store
func NewInMemoryEventStore() *InMemoryEventStore {
	return &InMemoryEventStore{
		events: make(map[string][]*Event),
	}
}

func (s *InMemoryEventStore) Save(ctx context.Context, events ...*Event) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	for _, e := range events {
		s.events[e.AggregateID] = append(s.events[e.AggregateID], e)
	}
	return nil
}

func (s *InMemoryEventStore) Load(ctx context.Context, aggregateID string) ([]*Event, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.events[aggregateID], nil
}

func (s *InMemoryEventStore) LoadFrom(ctx context.Context, aggregateID string, version int) ([]*Event, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	all := s.events[aggregateID]
	result := make([]*Event, 0)
	for _, e := range all {
		if e.Version > version {
			result = append(result, e)
		}
	}
	return result, nil
}

// Errors
var (
	ErrHandlerNotFound = &Error{Code: "HANDLER_NOT_FOUND", Message: "command handler not found"}
)

type Error struct {
	Code    string
	Message string
}

func (e *Error) Error() string {
	return e.Message
}
