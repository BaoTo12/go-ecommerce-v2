package saga

import (
	"context"
	"errors"
	"testing"
	"time"
)

func TestSagaOrchestrator_Success(t *testing.T) {
	store := NewInMemorySagaStore()
	orchestrator := NewOrchestrator(store)

	// Track execution order
	executed := make([]string, 0)

	// Register saga with 3 steps
	orchestrator.RegisterSaga(&SagaDefinition{
		Name:    "test-saga",
		Timeout: 30 * time.Second,
		Steps: []Step{
			{
				Name: "step1",
				Execute: func(ctx context.Context, data map[string]interface{}) error {
					executed = append(executed, "step1")
					return nil
				},
				Compensate: func(ctx context.Context, data map[string]interface{}) error {
					executed = append(executed, "compensate1")
					return nil
				},
			},
			{
				Name: "step2",
				Execute: func(ctx context.Context, data map[string]interface{}) error {
					executed = append(executed, "step2")
					return nil
				},
				Compensate: func(ctx context.Context, data map[string]interface{}) error {
					executed = append(executed, "compensate2")
					return nil
				},
			},
			{
				Name: "step3",
				Execute: func(ctx context.Context, data map[string]interface{}) error {
					executed = append(executed, "step3")
					return nil
				},
			},
		},
	})

	// Execute
	exec, err := orchestrator.Execute(context.Background(), "test-saga", map[string]interface{}{
		"order_id": "123",
	})
	if err != nil {
		t.Fatalf("Failed to start saga: %v", err)
	}

	// Wait for completion
	time.Sleep(100 * time.Millisecond)

	// Verify
	status, _ := orchestrator.GetStatus(context.Background(), exec.ID)
	if status.State != SagaStateCompleted {
		t.Errorf("Expected COMPLETED, got %s", status.State)
	}

	if len(executed) != 3 {
		t.Errorf("Expected 3 steps executed, got %d", len(executed))
	}
}

func TestSagaOrchestrator_Compensation(t *testing.T) {
	store := NewInMemorySagaStore()
	orchestrator := NewOrchestrator(store)

	compensated := make([]string, 0)

	orchestrator.RegisterSaga(&SagaDefinition{
		Name: "fail-saga",
		Steps: []Step{
			{
				Name: "step1",
				Execute: func(ctx context.Context, data map[string]interface{}) error {
					return nil
				},
				Compensate: func(ctx context.Context, data map[string]interface{}) error {
					compensated = append(compensated, "step1")
					return nil
				},
			},
			{
				Name: "step2",
				Execute: func(ctx context.Context, data map[string]interface{}) error {
					return errors.New("step2 failed")
				},
				Compensate: func(ctx context.Context, data map[string]interface{}) error {
					compensated = append(compensated, "step2")
					return nil
				},
			},
		},
	})

	exec, _ := orchestrator.Execute(context.Background(), "fail-saga", nil)
	time.Sleep(100 * time.Millisecond)

	status, _ := orchestrator.GetStatus(context.Background(), exec.ID)
	if status.State != SagaStateCompensated {
		t.Errorf("Expected COMPENSATED, got %s", status.State)
	}

	if len(compensated) != 1 {
		t.Errorf("Expected 1 compensation, got %d", len(compensated))
	}
}

func TestInMemorySagaStore(t *testing.T) {
	store := NewInMemorySagaStore()
	ctx := context.Background()

	saga := &SagaExecution{
		ID:           "test-saga-1",
		DefinitionID: "test",
		State:        SagaStateRunning,
	}

	// Save
	err := store.Save(ctx, saga)
	if err != nil {
		t.Fatalf("Save failed: %v", err)
	}

	// Load
	loaded, err := store.Load(ctx, "test-saga-1")
	if err != nil {
		t.Fatalf("Load failed: %v", err)
	}

	if loaded.State != SagaStateRunning {
		t.Errorf("Expected RUNNING, got %s", loaded.State)
	}

	// List pending
	pending, _ := store.ListPending(ctx)
	if len(pending) != 1 {
		t.Errorf("Expected 1 pending, got %d", len(pending))
	}
}
