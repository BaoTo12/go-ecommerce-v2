package cdc

import (
	"context"
	"encoding/json"
	"sync"
	"time"
)

/*
CHANGE DATA CAPTURE (CDC) - Track Database Changes

Implements CDC pattern for capturing and streaming database changes.
Useful for:
- Event-driven architectures
- Real-time data synchronization
- Audit logging
- Cache invalidation
*/

// ChangeType represents the type of change
type ChangeType string

const (
	ChangeTypeInsert ChangeType = "INSERT"
	ChangeTypeUpdate ChangeType = "UPDATE"
	ChangeTypeDelete ChangeType = "DELETE"
)

// ChangeEvent represents a single change event
type ChangeEvent struct {
	ID          string                 `json:"id"`
	Type        ChangeType             `json:"type"`
	Table       string                 `json:"table"`
	Schema      string                 `json:"schema"`
	PrimaryKey  map[string]interface{} `json:"primary_key"`
	Before      map[string]interface{} `json:"before,omitempty"`
	After       map[string]interface{} `json:"after,omitempty"`
	Timestamp   time.Time              `json:"timestamp"`
	Transaction string                 `json:"transaction_id,omitempty"`
	Sequence    int64                  `json:"sequence"`
}

// ChangeHandler processes change events
type ChangeHandler func(event *ChangeEvent) error

// CDCSource captures changes from a data source
type CDCSource interface {
	Start(ctx context.Context) error
	Stop() error
	Subscribe(handler ChangeHandler) string
	Unsubscribe(id string)
}

// CDCSink receives change events
type CDCSink interface {
	Handle(ctx context.Context, event *ChangeEvent) error
}

// CDCPipeline connects sources to sinks
type CDCPipeline struct {
	source    CDCSource
	sinks     []CDCSink
	filters   []ChangeFilter
	transformers []ChangeTransformer
	
	running bool
	stopCh  chan struct{}
	wg      sync.WaitGroup
	mu      sync.Mutex
}

// ChangeFilter filters change events
type ChangeFilter func(event *ChangeEvent) bool

// ChangeTransformer transforms change events
type ChangeTransformer func(event *ChangeEvent) *ChangeEvent

// NewCDCPipeline creates a CDC pipeline
func NewCDCPipeline(source CDCSource) *CDCPipeline {
	return &CDCPipeline{
		source:       source,
		sinks:        make([]CDCSink, 0),
		filters:      make([]ChangeFilter, 0),
		transformers: make([]ChangeTransformer, 0),
		stopCh:       make(chan struct{}),
	}
}

// AddSink adds a sink to the pipeline
func (p *CDCPipeline) AddSink(sink CDCSink) *CDCPipeline {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.sinks = append(p.sinks, sink)
	return p
}

// AddFilter adds a filter to the pipeline
func (p *CDCPipeline) AddFilter(filter ChangeFilter) *CDCPipeline {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.filters = append(p.filters, filter)
	return p
}

// AddTransformer adds a transformer to the pipeline
func (p *CDCPipeline) AddTransformer(transformer ChangeTransformer) *CDCPipeline {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.transformers = append(p.transformers, transformer)
	return p
}

// Start starts the CDC pipeline
func (p *CDCPipeline) Start(ctx context.Context) error {
	p.mu.Lock()
	if p.running {
		p.mu.Unlock()
		return nil
	}
	p.running = true
	p.mu.Unlock()

	// Subscribe to source
	p.source.Subscribe(func(event *ChangeEvent) error {
		return p.processEvent(ctx, event)
	})

	return p.source.Start(ctx)
}

// Stop stops the CDC pipeline
func (p *CDCPipeline) Stop() error {
	p.mu.Lock()
	if !p.running {
		p.mu.Unlock()
		return nil
	}
	p.running = false
	p.mu.Unlock()

	close(p.stopCh)
	p.wg.Wait()
	return p.source.Stop()
}

func (p *CDCPipeline) processEvent(ctx context.Context, event *ChangeEvent) error {
	// Apply filters
	for _, filter := range p.filters {
		if !filter(event) {
			return nil // Filtered out
		}
	}

	// Apply transformers
	for _, transformer := range p.transformers {
		event = transformer(event)
		if event == nil {
			return nil // Dropped by transformer
		}
	}

	// Send to all sinks
	for _, sink := range p.sinks {
		if err := sink.Handle(ctx, event); err != nil {
			// Log error but continue with other sinks
			continue
		}
	}

	return nil
}

// InMemoryCDCSource simulates CDC for testing
type InMemoryCDCSource struct {
	events    chan *ChangeEvent
	handlers  map[string]ChangeHandler
	sequence  int64
	running   bool
	stopCh    chan struct{}
	mu        sync.RWMutex
}

// NewInMemoryCDCSource creates an in-memory CDC source
func NewInMemoryCDCSource() *InMemoryCDCSource {
	return &InMemoryCDCSource{
		events:   make(chan *ChangeEvent, 1000),
		handlers: make(map[string]ChangeHandler),
		stopCh:   make(chan struct{}),
	}
}

func (s *InMemoryCDCSource) Start(ctx context.Context) error {
	s.mu.Lock()
	if s.running {
		s.mu.Unlock()
		return nil
	}
	s.running = true
	s.mu.Unlock()

	go s.processEvents(ctx)
	return nil
}

func (s *InMemoryCDCSource) Stop() error {
	s.mu.Lock()
	if !s.running {
		s.mu.Unlock()
		return nil
	}
	s.running = false
	s.mu.Unlock()
	close(s.stopCh)
	return nil
}

func (s *InMemoryCDCSource) processEvents(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			return
		case <-s.stopCh:
			return
		case event := <-s.events:
			s.mu.RLock()
			for _, handler := range s.handlers {
				handler(event)
			}
			s.mu.RUnlock()
		}
	}
}

func (s *InMemoryCDCSource) Subscribe(handler ChangeHandler) string {
	s.mu.Lock()
	defer s.mu.Unlock()
	id := generateCDCID()
	s.handlers[id] = handler
	return id
}

func (s *InMemoryCDCSource) Unsubscribe(id string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	delete(s.handlers, id)
}

// Emit emits a change event (for simulation/testing)
func (s *InMemoryCDCSource) Emit(event *ChangeEvent) {
	s.mu.Lock()
	s.sequence++
	event.Sequence = s.sequence
	event.Timestamp = time.Now()
	event.ID = generateCDCID()
	s.mu.Unlock()
	
	s.events <- event
}

// LoggingSink logs change events
type LoggingSink struct {
	prefix string
}

func NewLoggingSink(prefix string) *LoggingSink {
	return &LoggingSink{prefix: prefix}
}

func (s *LoggingSink) Handle(ctx context.Context, event *ChangeEvent) error {
	data, _ := json.Marshal(event)
	println(s.prefix, string(data))
	return nil
}

// KafkaSink sends events to Kafka (placeholder)
type KafkaSink struct {
	topic   string
	brokers []string
}

func NewKafkaSink(topic string, brokers []string) *KafkaSink {
	return &KafkaSink{
		topic:   topic,
		brokers: brokers,
	}
}

func (s *KafkaSink) Handle(ctx context.Context, event *ChangeEvent) error {
	// In real implementation, publish to Kafka
	// kafka.Publish(s.topic, event.PrimaryKey, event)
	return nil
}

// Common filters
func TableFilter(tables ...string) ChangeFilter {
	tableSet := make(map[string]bool)
	for _, t := range tables {
		tableSet[t] = true
	}
	return func(event *ChangeEvent) bool {
		return tableSet[event.Table]
	}
}

func TypeFilter(types ...ChangeType) ChangeFilter {
	typeSet := make(map[ChangeType]bool)
	for _, t := range types {
		typeSet[t] = true
	}
	return func(event *ChangeEvent) bool {
		return typeSet[event.Type]
	}
}

func generateCDCID() string {
	return "cdc_" + time.Now().Format("20060102150405.000000")
}
