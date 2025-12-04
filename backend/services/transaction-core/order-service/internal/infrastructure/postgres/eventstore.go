package postgres

import (
	"context"
	"database/sql"
	"encoding/json"
	"time"

	"github.com/titan-commerce/backend/order-service/internal/domain"
	"github.com/titan-commerce/backend/pkg/errors"
	"github.com/titan-commerce/backend/pkg/logger"
	_ "github.com/lib/pq"
)

type EventStoreRepository struct {
	db     *sql.DB
	logger *logger.Logger
}

func NewEventStoreRepository(databaseURL string, logger *logger.Logger) (*EventStoreRepository, error) {
	db, err := sql.Open("postgres", databaseURL)
	if err != nil {
		return nil, errors.Wrap(errors.ErrInternal, "failed to connect to database", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	
	if err := db.PingContext(ctx); err != nil {
		return nil, errors.Wrap(errors.ErrInternal, "failed to ping database", err)
	}

	logger.Info("Order EventStore repository initialized")
	return &EventStoreRepository{db: db, logger: logger}, nil
}

// SaveEvents appends events to the event stream
func (r *EventStoreRepository) SaveEvents(ctx context.Context, aggregateID string, events []domain.DomainEvent, expectedVersion int) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return errors.Wrap(errors.ErrInternal, "failed to begin transaction", err)
	}
	defer tx.Rollback()

	// Check current version (optimistic concurrency)
	var currentVersion int
	err = tx.QueryRowContext(ctx, 
		"SELECT COALESCE(MAX(version), 0) FROM events WHERE aggregate_id = $1", aggregateID).Scan(&currentVersion)
	if err != nil {
		return errors.Wrap(errors.ErrInternal, "failed to get current version", err)
	}

	if currentVersion != expectedVersion {
		return errors.New(errors.ErrConflict, "concurrent modification detected")
	}

	// Insert events
	query := `
		INSERT INTO events (aggregate_id, aggregate_type, event_type, event_data, version, occurred_at)
		VALUES ($1, $2, $3, $4, $5, $6)
	`

	for i, event := range events {
		eventData, err := json.Marshal(event)
		if err != nil {
			return errors.Wrap(errors.ErrInternal, "failed to marshal event", err)
		}

		version := expectedVersion + i + 1
		_, err = tx.ExecContext(ctx, query,
			aggregateID, "Order", event.EventType(), eventData, version, event.OccurredAt())
		if err != nil {
			return errors.Wrap(errors.ErrInternal, "failed to save event", err)
		}
	}

	if err := tx.Commit(); err != nil {
		return errors.Wrap(errors.ErrInternal, "failed to commit transaction", err)
	}

	return nil
}

// GetEvents retrieves all events for an aggregate
func (r *EventStoreRepository) GetEvents(ctx context.Context, aggregateID string) ([]domain.DomainEvent, error) {
	query := `
		SELECT event_type, event_data, occurred_at
		FROM events
		WHERE aggregate_id = $1
		ORDER BY version ASC
	`

	rows, err := r.db.QueryContext(ctx, query, aggregateID)
	if err != nil {
		return nil, errors.Wrap(errors.ErrInternal, "failed to query events", err)
	}
	defer rows.Close()

	var events []domain.DomainEvent
	for rows.Next() {
		var eventType string
		var eventData []byte
		var occurredAt time.Time

		if err := rows.Scan(&eventType, &eventData, &occurredAt); err != nil {
			return nil, errors.Wrap(errors.ErrInternal, "failed to scan event", err)
		}

		// Deserialize based on event type
		var event domain.DomainEvent
		switch eventType {
		case "OrderCreated":
			event = &domain.OrderCreatedEvent{}
		case "OrderPaid":
			event = &domain.OrderPaidEvent{}
		case "OrderShipped":
			event = &domain.OrderShippedEvent{}
		case "OrderCompleted":
			event = &domain.OrderCompletedEvent{}
		case "OrderCancelled":
			event = &domain.OrderCancelledEvent{}
		default:
			r.logger.Warnf("Unknown event type: %s", eventType)
			continue
		}

		if err := json.Unmarshal(eventData, event); err != nil {
			return nil, errors.Wrap(errors.ErrInternal, "failed to unmarshal event", err)
		}

		events = append(events, event)
	}

	return events, nil
}

// Read model repository
type OrderReadModelRepository struct {
	db     *sql.DB
	logger *logger.Logger
}

func NewOrderReadModelRepository(databaseURL string, logger *logger.Logger) (*OrderReadModelRepository, error) {
	db, err := sql.Open("postgres", databaseURL)
	if err != nil {
		return nil, errors.Wrap(errors.ErrInternal, "failed to connect to database", err)
	}

	logger.Info("Order read model repository initialized")
	return &OrderReadModelRepository{db: db, logger: logger}, nil
}

// SaveReadModel updates denormalized read model
func (r *OrderReadModelRepository) SaveReadModel(ctx context.Context, order *domain.Order) error {
	itemsJSON, err := json.Marshal(order.Items)
	if err != nil {
		return errors.Wrap(errors.ErrInternal, "failed to marshal items", err)
	}

	query := `
		INSERT INTO orders_read_model (order_id, user_id, status, total_amount, items, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		ON CONFLICT (order_id) DO UPDATE SET
			status = EXCLUDED.status,
			updated_at = EXCLUDED.updated_at
	`

	_, err = r.db.ExecContext(ctx, query,
		order.ID, order.UserID, order.Status, order.TotalAmount, itemsJSON,
		order.CreatedAt, order.UpdatedAt)

	if err != nil {
		return errors.Wrap(errors.ErrInternal, "failed to save read model", err)
	}

	return nil
}

// FindByID retrieves order from read model
func (r *OrderReadModelRepository) FindByID(ctx context.Context, orderID string) (*domain.Order, error) {
	query := `
		SELECT order_id, user_id, status, total_amount, items, created_at, updated_at
		FROM orders_read_model
		WHERE order_id = $1
	`

	var order domain.Order
	var itemsJSON []byte

	err := r.db.QueryRowContext(ctx, query, orderID).Scan(
		&order.ID, &order.UserID, &order.Status, &order.TotalAmount,
		&itemsJSON, &order.CreatedAt, &order.UpdatedAt)

	if err == sql.ErrNoRows {
		return nil, errors.New(errors.ErrNotFound, "order not found")
	}
	if err != nil {
		return nil, errors.Wrap(errors.ErrInternal, "failed to find order", err)
	}

	if err := json.Unmarshal(itemsJSON, &order.Items); err != nil {
		return nil, errors.Wrap(errors.ErrInternal, "failed to unmarshal items", err)
	}

	return &order, nil
}
