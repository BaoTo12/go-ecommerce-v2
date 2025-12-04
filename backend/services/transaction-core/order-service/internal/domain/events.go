package domain

import (
	"time"
)

// EventType represents the type of domain event
type EventType string

const (
	EventTypeOrderCreated       EventType = "order.created"
	EventTypeOrderConfirmed     EventType = "order.confirmed"
	EventTypeOrderCancelled     EventType = "order.cancelled"
	EventTypeOrderStatusUpdated EventType = "order.status.updated"
)

// DomainEvent interface for all domain events
type DomainEvent interface {
	EventType() string
	OccurredAt() time.Time
}

// OrderEvent represents a domain event
type OrderEvent struct {
	ID          string
	AggregateID string // Order ID
	Type        EventType
	Data        map[string]interface{}
	UserID      string
	Timestamp   time.Time
	Version     int
}

// NewOrderCreatedEvent creates an order created event
func NewOrderCreatedEvent(order *Order) *OrderEvent {
	return &OrderEvent{
		AggregateID: order.ID,
		Type:        EventTypeOrderCreated,
		Data: map[string]interface{}{
			"user_id":          order.UserID,
			"total_amount":     order.TotalAmount,
			"shipping_address": order.ShippingAddress,
			"items_count":      len(order.Items),
		},
		UserID:    order.UserID,
		Timestamp: order.CreatedAt,
		Version:   order.Version,
	}
}

// OrderCreatedEvent for eventstore deserialization
type OrderCreatedEvent struct {
	OrderID   string    `json:"order_id"`
	UserID    string    `json:"user_id"`
	Total     float64   `json:"total"`
	CreatedAt time.Time `json:"created_at"`
}

func (e *OrderCreatedEvent) EventType() string    { return "OrderCreated" }
func (e *OrderCreatedEvent) OccurredAt() time.Time { return e.CreatedAt }

// OrderPaidEvent
type OrderPaidEvent struct {
	OrderID   string    `json:"order_id"`
	PaymentID string    `json:"payment_id"`
	Amount    float64   `json:"amount"`
	PaidAt    time.Time `json:"paid_at"`
}

func (e *OrderPaidEvent) EventType() string    { return "OrderPaid" }
func (e *OrderPaidEvent) OccurredAt() time.Time { return e.PaidAt }

// OrderShippedEvent
type OrderShippedEvent struct {
	OrderID    string    `json:"order_id"`
	TrackingNo string    `json:"tracking_no"`
	ShippedAt  time.Time `json:"shipped_at"`
}

func (e *OrderShippedEvent) EventType() string    { return "OrderShipped" }
func (e *OrderShippedEvent) OccurredAt() time.Time { return e.ShippedAt }

// OrderCompletedEvent
type OrderCompletedEvent struct {
	OrderID     string    `json:"order_id"`
	CompletedAt time.Time `json:"completed_at"`
}

func (e *OrderCompletedEvent) EventType() string    { return "OrderCompleted" }
func (e *OrderCompletedEvent) OccurredAt() time.Time { return e.CompletedAt }

// OrderCancelledEvent
type OrderCancelledEvent struct {
	OrderID     string    `json:"order_id"`
	Reason      string    `json:"reason"`
	CancelledAt time.Time `json:"cancelled_at"`
}

func (e *OrderCancelledEvent) EventType() string    { return "OrderCancelled" }
func (e *OrderCancelledEvent) OccurredAt() time.Time { return e.CancelledAt }
