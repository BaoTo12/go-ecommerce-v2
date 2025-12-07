package saga

import (
	"context"
	"encoding/json"
	"errors"
	"sync"
	"time"
)

/*
SAGA ORCHESTRATOR - Distributed Transaction Pattern

Implements the Saga pattern for managing distributed transactions
across multiple microservices with compensation (rollback) support.

Example: Order Creation Saga
1. Reserve inventory → Compensate: Release inventory
2. Process payment → Compensate: Refund payment
3. Create shipment → Compensate: Cancel shipment
4. Send notification
*/

var (
	ErrSagaFailed    = errors.New("saga execution failed")
	ErrStepFailed    = errors.New("step execution failed")
	ErrCompensation  = errors.New("compensation failed")
	ErrTimeout       = errors.New("saga timeout")
)

// SagaState represents the current state of a saga
type SagaState string

const (
	SagaStateNew        SagaState = "NEW"
	SagaStateRunning    SagaState = "RUNNING"
	SagaStateCompleted  SagaState = "COMPLETED"
	SagaStateFailed     SagaState = "FAILED"
	SagaStateCompensating SagaState = "COMPENSATING"
	SagaStateCompensated SagaState = "COMPENSATED"
)

// Step represents a single step in a saga
type Step struct {
	Name        string
	Execute     func(ctx context.Context, data map[string]interface{}) error
	Compensate  func(ctx context.Context, data map[string]interface{}) error
	Timeout     time.Duration
	RetryCount  int
	RetryDelay  time.Duration
}

// SagaDefinition defines a saga with its steps
type SagaDefinition struct {
	Name    string
	Steps   []Step
	Timeout time.Duration
}

// SagaExecution represents a running saga instance
type SagaExecution struct {
	ID            string
	DefinitionID  string
	State         SagaState
	CurrentStep   int
	Data          map[string]interface{}
	CompletedSteps []string
	Errors        []SagaError
	StartedAt     time.Time
	CompletedAt   *time.Time
	mu            sync.Mutex
}

// SagaError represents an error that occurred during saga execution
type SagaError struct {
	Step      string
	Error     string
	Timestamp time.Time
	IsCompensation bool
}

// SagaStore persists saga state
type SagaStore interface {
	Save(ctx context.Context, saga *SagaExecution) error
	Load(ctx context.Context, id string) (*SagaExecution, error)
	ListPending(ctx context.Context) ([]*SagaExecution, error)
}

// Orchestrator manages saga executions
type Orchestrator struct {
	definitions map[string]*SagaDefinition
	store       SagaStore
	mu          sync.RWMutex
}

// NewOrchestrator creates a saga orchestrator
func NewOrchestrator(store SagaStore) *Orchestrator {
	return &Orchestrator{
		definitions: make(map[string]*SagaDefinition),
		store:       store,
	}
}

// RegisterSaga registers a saga definition
func (o *Orchestrator) RegisterSaga(def *SagaDefinition) {
	o.mu.Lock()
	defer o.mu.Unlock()
	o.definitions[def.Name] = def
}

// Execute starts a new saga execution
func (o *Orchestrator) Execute(ctx context.Context, sagaName string, data map[string]interface{}) (*SagaExecution, error) {
	o.mu.RLock()
	def, exists := o.definitions[sagaName]
	o.mu.RUnlock()

	if !exists {
		return nil, errors.New("saga not found: " + sagaName)
	}

	// Create execution
	exec := &SagaExecution{
		ID:           generateSagaID(),
		DefinitionID: sagaName,
		State:        SagaStateNew,
		CurrentStep:  0,
		Data:         data,
		CompletedSteps: make([]string, 0),
		Errors:       make([]SagaError, 0),
		StartedAt:    time.Now(),
	}

	// Save initial state
	if err := o.store.Save(ctx, exec); err != nil {
		return nil, err
	}

	// Execute saga
	go o.executeSaga(ctx, def, exec)

	return exec, nil
}

func (o *Orchestrator) executeSaga(ctx context.Context, def *SagaDefinition, exec *SagaExecution) {
	exec.mu.Lock()
	exec.State = SagaStateRunning
	exec.mu.Unlock()
	o.store.Save(ctx, exec)

	// Create saga-level timeout
	sagaCtx := ctx
	if def.Timeout > 0 {
		var cancel context.CancelFunc
		sagaCtx, cancel = context.WithTimeout(ctx, def.Timeout)
		defer cancel()
	}

	// Execute each step
	for i, step := range def.Steps {
		exec.mu.Lock()
		exec.CurrentStep = i
		exec.mu.Unlock()

		err := o.executeStep(sagaCtx, step, exec)
		if err != nil {
			exec.mu.Lock()
			exec.Errors = append(exec.Errors, SagaError{
				Step:      step.Name,
				Error:     err.Error(),
				Timestamp: time.Now(),
			})
			exec.State = SagaStateFailed
			exec.mu.Unlock()
			o.store.Save(ctx, exec)

			// Start compensation
			o.compensate(ctx, def, exec, i-1)
			return
		}

		exec.mu.Lock()
		exec.CompletedSteps = append(exec.CompletedSteps, step.Name)
		exec.mu.Unlock()
		o.store.Save(ctx, exec)
	}

	// Saga completed successfully
	exec.mu.Lock()
	exec.State = SagaStateCompleted
	now := time.Now()
	exec.CompletedAt = &now
	exec.mu.Unlock()
	o.store.Save(ctx, exec)
}

func (o *Orchestrator) executeStep(ctx context.Context, step Step, exec *SagaExecution) error {
	stepCtx := ctx
	if step.Timeout > 0 {
		var cancel context.CancelFunc
		stepCtx, cancel = context.WithTimeout(ctx, step.Timeout)
		defer cancel()
	}

	retries := step.RetryCount
	if retries == 0 {
		retries = 1
	}

	var lastErr error
	for i := 0; i < retries; i++ {
		if i > 0 && step.RetryDelay > 0 {
			time.Sleep(step.RetryDelay)
		}

		err := step.Execute(stepCtx, exec.Data)
		if err == nil {
			return nil
		}
		lastErr = err

		// Check if context is done
		if stepCtx.Err() != nil {
			return stepCtx.Err()
		}
	}

	return lastErr
}

func (o *Orchestrator) compensate(ctx context.Context, def *SagaDefinition, exec *SagaExecution, fromStep int) {
	exec.mu.Lock()
	exec.State = SagaStateCompensating
	exec.mu.Unlock()
	o.store.Save(ctx, exec)

	// Compensate in reverse order
	for i := fromStep; i >= 0; i-- {
		step := def.Steps[i]
		if step.Compensate == nil {
			continue
		}

		err := step.Compensate(ctx, exec.Data)
		if err != nil {
			exec.mu.Lock()
			exec.Errors = append(exec.Errors, SagaError{
				Step:           step.Name,
				Error:          err.Error(),
				Timestamp:      time.Now(),
				IsCompensation: true,
			})
			exec.mu.Unlock()
			// Continue with other compensations even if one fails
		}
	}

	exec.mu.Lock()
	exec.State = SagaStateCompensated
	now := time.Now()
	exec.CompletedAt = &now
	exec.mu.Unlock()
	o.store.Save(ctx, exec)
}

// GetStatus returns the current status of a saga
func (o *Orchestrator) GetStatus(ctx context.Context, sagaID string) (*SagaExecution, error) {
	return o.store.Load(ctx, sagaID)
}

// RecoverPending recovers pending sagas after restart
func (o *Orchestrator) RecoverPending(ctx context.Context) error {
	pending, err := o.store.ListPending(ctx)
	if err != nil {
		return err
	}

	for _, exec := range pending {
		o.mu.RLock()
		def, exists := o.definitions[exec.DefinitionID]
		o.mu.RUnlock()

		if !exists {
			continue
		}

		// Resume or compensate based on state
		if exec.State == SagaStateRunning {
			go o.executeSaga(ctx, def, exec)
		} else if exec.State == SagaStateCompensating {
			go o.compensate(ctx, def, exec, exec.CurrentStep)
		}
	}

	return nil
}

// InMemorySagaStore is an in-memory implementation
type InMemorySagaStore struct {
	sagas map[string]*SagaExecution
	mu    sync.RWMutex
}

func NewInMemorySagaStore() *InMemorySagaStore {
	return &InMemorySagaStore{
		sagas: make(map[string]*SagaExecution),
	}
}

func (s *InMemorySagaStore) Save(ctx context.Context, saga *SagaExecution) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	
	// Deep copy to prevent mutations
	data, _ := json.Marshal(saga)
	var copy SagaExecution
	json.Unmarshal(data, &copy)
	s.sagas[saga.ID] = &copy
	return nil
}

func (s *InMemorySagaStore) Load(ctx context.Context, id string) (*SagaExecution, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	
	saga, exists := s.sagas[id]
	if !exists {
		return nil, errors.New("saga not found")
	}
	return saga, nil
}

func (s *InMemorySagaStore) ListPending(ctx context.Context) ([]*SagaExecution, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	
	pending := make([]*SagaExecution, 0)
	for _, saga := range s.sagas {
		if saga.State == SagaStateRunning || saga.State == SagaStateCompensating {
			pending = append(pending, saga)
		}
	}
	return pending, nil
}

func generateSagaID() string {
	return "saga_" + time.Now().Format("20060102150405") + "_" + randomString(8)
}

func randomString(n int) string {
	const letters = "abcdefghijklmnopqrstuvwxyz0123456789"
	b := make([]byte, n)
	for i := range b {
		b[i] = letters[time.Now().UnixNano()%int64(len(letters))]
	}
	return string(b)
}
