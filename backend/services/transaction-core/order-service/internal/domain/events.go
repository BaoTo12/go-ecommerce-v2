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
