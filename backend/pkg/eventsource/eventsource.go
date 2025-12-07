package eventsource

import (
	"context"
	"encoding/json"
	"errors"
	"reflect"
	"sync"
	"time"
)

/*
EVENT SOURCING - Full Implementation

Implements complete event sourcing with:
- Event Store
- Aggregate Root base class
- Projections (read models)
- Snapshots for performance
- Event replay
*/

var (
	ErrAggregateNotFound = errors.New("aggregate not found")
	ErrConcurrencyConflict = errors.New("concurrency conflict")
)

// Event represents a domain event
type Event struct {
	ID            string                 `json:"id"`
	AggregateID   string                 `json:"aggregate_id"`
	AggregateType string                 `json:"aggregate_type"`
	Type          string                 `json:"type"`
	Version       int                    `json:"version"`
	Timestamp     time.Time              `json:"timestamp"`
	Data          map[string]interface{} `json:"data"`
	Metadata      map[string]string      `json:"metadata"`
}

// EventStore stores and retrieves events
type EventStore interface {
	// Append appends events to an aggregate stream
	Append(ctx context.Context, aggregateID string, expectedVersion int, events []Event) error
	
	// Load loads all events for an aggregate
	Load(ctx context.Context, aggregateID string) ([]Event, error)
	
	// LoadFrom loads events from a specific version
	LoadFrom(ctx context.Context, aggregateID string, fromVersion int) ([]Event, error)
	
	// LoadAll loads all events (for projections)
	LoadAll(ctx context.Context, fromSequence int64) ([]Event, error)
}

// Aggregate is the base for event-sourced aggregates
type Aggregate struct {
	ID      string
	Version int
	changes []Event
}

// AggregateRoot interface for aggregates
type AggregateRoot interface {
	GetID() string
	GetVersion() int
	GetChanges() []Event
	ClearChanges()
	Apply(event Event)
}

// GetID returns the aggregate ID
func (a *Aggregate) GetID() string { return a.ID }

// GetVersion returns the current version
func (a *Aggregate) GetVersion() int { return a.Version }

// GetChanges returns uncommitted changes
func (a *Aggregate) GetChanges() []Event { return a.changes }

// ClearChanges clears uncommitted changes
func (a *Aggregate) ClearChanges() { a.changes = nil }

// RaiseEvent raises a new event
func (a *Aggregate) RaiseEvent(eventType string, data map[string]interface{}) Event {
	a.Version++
	event := Event{
		ID:          generateEventID(),
		AggregateID: a.ID,
		Type:        eventType,
		Version:     a.Version,
		Timestamp:   time.Now(),
		Data:        data,
	}
	a.changes = append(a.changes, event)
	return event
}

// Repository handles aggregate persistence
type Repository[T AggregateRoot] struct {
	store   EventStore
	factory func() T
}

// NewRepository creates a repository
func NewRepository[T AggregateRoot](store EventStore, factory func() T) *Repository[T] {
	return &Repository[T]{
		store:   store,
		factory: factory,
	}
}

// Load loads an aggregate by ID
func (r *Repository[T]) Load(ctx context.Context, id string) (T, error) {
	events, err := r.store.Load(ctx, id)
	if err != nil {
		var zero T
		return zero, err
	}

	if len(events) == 0 {
		var zero T
		return zero, ErrAggregateNotFound
	}

	aggregate := r.factory()
	for _, event := range events {
		aggregate.Apply(event)
	}

	return aggregate, nil
}

// Save saves an aggregate
func (r *Repository[T]) Save(ctx context.Context, aggregate T) error {
	changes := aggregate.GetChanges()
	if len(changes) == 0 {
		return nil
	}

	expectedVersion := aggregate.GetVersion() - len(changes)
	err := r.store.Append(ctx, aggregate.GetID(), expectedVersion, changes)
	if err != nil {
		return err
	}

	aggregate.ClearChanges()
	return nil
}

// Projection builds read models from events
type Projection interface {
	Handle(ctx context.Context, event Event) error
	Reset(ctx context.Context) error
}

// ProjectionManager manages projections
type ProjectionManager struct {
	store       EventStore
	projections map[string]Projection
	positions   map[string]int64
	mu          sync.RWMutex
}

// NewProjectionManager creates a projection manager
func NewProjectionManager(store EventStore) *ProjectionManager {
	return &ProjectionManager{
		store:       store,
		projections: make(map[string]Projection),
		positions:   make(map[string]int64),
	}
}

// Register registers a projection
func (pm *ProjectionManager) Register(name string, projection Projection) {
	pm.mu.Lock()
	defer pm.mu.Unlock()
	pm.projections[name] = projection
	pm.positions[name] = 0
}

// Rebuild rebuilds all projections from scratch
func (pm *ProjectionManager) Rebuild(ctx context.Context) error {
	pm.mu.Lock()
	defer pm.mu.Unlock()

	// Reset all projections
	for _, proj := range pm.projections {
		if err := proj.Reset(ctx); err != nil {
			return err
		}
	}

	// Reset positions
	for name := range pm.positions {
		pm.positions[name] = 0
	}

	// Replay all events
	return pm.catchUp(ctx)
}

// CatchUp catches up projections with new events
func (pm *ProjectionManager) CatchUp(ctx context.Context) error {
	pm.mu.Lock()
	defer pm.mu.Unlock()
	return pm.catchUp(ctx)
}

func (pm *ProjectionManager) catchUp(ctx context.Context) error {
	// Find minimum position
	minPos := int64(0)
	for _, pos := range pm.positions {
		if minPos == 0 || pos < minPos {
			minPos = pos
		}
	}

	// Load events from minimum position
	events, err := pm.store.LoadAll(ctx, minPos)
	if err != nil {
		return err
	}

	// Apply events to projections
	for _, event := range events {
		for name, proj := range pm.projections {
			if err := proj.Handle(ctx, event); err != nil {
				return err
			}
			pm.positions[name]++
		}
	}

	return nil
}

// Snapshot stores aggregate state for faster loading
type Snapshot struct {
	AggregateID   string
	AggregateType string
	Version       int
	State         []byte
	Timestamp     time.Time
}

// SnapshotStore stores snapshots
type SnapshotStore interface {
	Save(ctx context.Context, snapshot Snapshot) error
	Load(ctx context.Context, aggregateID string) (*Snapshot, error)
}

// SnapshotRepository uses snapshots for faster loading
type SnapshotRepository[T AggregateRoot] struct {
	eventStore    EventStore
	snapshotStore SnapshotStore
	factory       func() T
	snapshotEvery int // Take snapshot every N events
}

// NewSnapshotRepository creates a snapshot repository
func NewSnapshotRepository[T AggregateRoot](
	eventStore EventStore,
	snapshotStore SnapshotStore,
	factory func() T,
	snapshotEvery int,
) *SnapshotRepository[T] {
	return &SnapshotRepository[T]{
		eventStore:    eventStore,
		snapshotStore: snapshotStore,
		factory:       factory,
		snapshotEvery: snapshotEvery,
	}
}

// Load loads aggregate using snapshot if available
func (r *SnapshotRepository[T]) Load(ctx context.Context, id string) (T, error) {
	aggregate := r.factory()

	// Try to load snapshot
	snapshot, err := r.snapshotStore.Load(ctx, id)
	if err == nil && snapshot != nil {
		// Restore from snapshot
		if err := json.Unmarshal(snapshot.State, aggregate); err == nil {
			// Load events after snapshot
			events, err := r.eventStore.LoadFrom(ctx, id, snapshot.Version+1)
			if err != nil {
				return aggregate, err
			}
			for _, event := range events {
				aggregate.Apply(event)
			}
			return aggregate, nil
		}
	}

	// No snapshot - load all events
	events, err := r.eventStore.Load(ctx, id)
	if err != nil {
		var zero T
		return zero, err
	}

	for _, event := range events {
		aggregate.Apply(event)
	}

	return aggregate, nil
}

// InMemoryEventStore is an in-memory implementation
type InMemoryEventStore struct {
	events   map[string][]Event
	allEvents []Event
	sequence int64
	mu       sync.RWMutex
}

func NewInMemoryEventStore() *InMemoryEventStore {
	return &InMemoryEventStore{
		events:    make(map[string][]Event),
		allEvents: make([]Event, 0),
	}
}

func (s *InMemoryEventStore) Append(ctx context.Context, aggregateID string, expectedVersion int, events []Event) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	existing := s.events[aggregateID]
	if len(existing) != expectedVersion {
		return ErrConcurrencyConflict
	}

	// Set aggregate type from first event
	aggregateType := ""
	if len(events) > 0 {
		aggregateType = reflect.TypeOf(events[0]).Name()
	}

	for i := range events {
		events[i].AggregateType = aggregateType
		s.sequence++
		s.allEvents = append(s.allEvents, events[i])
	}

	s.events[aggregateID] = append(existing, events...)
	return nil
}

func (s *InMemoryEventStore) Load(ctx context.Context, aggregateID string) ([]Event, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.events[aggregateID], nil
}

func (s *InMemoryEventStore) LoadFrom(ctx context.Context, aggregateID string, fromVersion int) ([]Event, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	events := s.events[aggregateID]
	if fromVersion >= len(events) {
		return nil, nil
	}
	return events[fromVersion:], nil
}

func (s *InMemoryEventStore) LoadAll(ctx context.Context, fromSequence int64) ([]Event, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	if fromSequence >= int64(len(s.allEvents)) {
		return nil, nil
	}
	return s.allEvents[fromSequence:], nil
}

func generateEventID() string {
	return "evt_" + time.Now().Format("20060102150405.000000")
}
