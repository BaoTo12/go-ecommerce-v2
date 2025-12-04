package domain

import "context"

// Repository defines the interface for order persistence
type Repository interface {
	Save(ctx context.Context, order *Order) error
	FindByID(ctx context.Context, orderID string) (*Order, error)
	FindByUserID(ctx context.Context, userID string, limit, offset int) ([]*Order, error)
	Update(ctx context.Context, order *Order) error
}

// EventStore defines the interface for event sourcing
type EventStore interface {
	SaveEvent(ctx context.Context, event *OrderEvent) error
	LoadEvents(ctx context.Context, aggregateID string) ([]*OrderEvent, error)
}
