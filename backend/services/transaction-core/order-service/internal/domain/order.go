package domain

import (
	"time"

	"github.com/google/uuid"
	"github.com/titan-commerce/backend/pkg/errors"
)

// OrderStatus represents the status of an order
type OrderStatus string

const (
	OrderStatusPending    OrderStatus = "PENDING"
	OrderStatusConfirmed  OrderStatus = "CONFIRMED"
	OrderStatusProcessing OrderStatus = "PROCESSING"
	OrderStatusShipped    OrderStatus = "SHIPPED"
	OrderStatusDelivered  OrderStatus = "DELIVERED"
	OrderStatusCancelled  OrderStatus = "CANCELLED"
	OrderStatusRefunded   OrderStatus = "REFUNDED"
)

// OrderItem represents a line item in an order
type OrderItem struct {
	ProductID   string
	ProductName string
	Quantity    int
	UnitPrice   float64
	Subtotal    float64
}

// Order is the aggregate root for order domain
type Order struct {
	ID              string
	UserID          string
	Items           []OrderItem
	TotalAmount     float64
	Status          OrderStatus
	ShippingAddress string
	CreatedAt       time.Time
	UpdatedAt       time.Time
	Version         int // For optimistic locking
}

// NewOrder creates a new order (Factory method)
func NewOrder(userID string, items []OrderItem, shippingAddress string) (*Order, error) {
	if userID == "" {
		return nil, errors.New(errors.ErrInvalidInput, "user ID is required")
	}
	if len(items) == 0 {
		return nil, errors.New(errors.ErrInvalidInput, "order must have at least one item")
	}
	if shippingAddress == "" {
		return nil, errors.New(errors.ErrInvalidInput, "shipping address is required")
	}

	// Calculate total
	var total float64
	for i, item := range items {
		if item.Quantity <= 0 {
			return nil, errors.New(errors.ErrInvalidInput, "item quantity must be positive")
		}
		if item.UnitPrice <= 0 {
			return nil, errors.New(errors.ErrInvalidInput, "item unit price must be positive")
		}
		items[i].Subtotal = float64(item.Quantity) * item.UnitPrice
		total += items[i].Subtotal
	}

	now := time.Now()
	return &Order{
		ID:              uuid.New().String(),
		UserID:          userID,
		Items:           items,
		TotalAmount:     total,
		Status:          OrderStatusPending,
		ShippingAddress: shippingAddress,
		CreatedAt:       now,
		UpdatedAt:       now,
		Version:         1,
	}, nil
}

// Confirm transitions order to confirmed status
func (o *Order) Confirm() error {
	if o.Status != OrderStatusPending {
		return errors.New(errors.ErrInvalidInput, "only pending orders can be confirmed")
	}
	o.Status = OrderStatusConfirmed
	o.UpdatedAt = time.Now()
	o.Version++
	return nil
}

// Cancel cancels the order
func (o *Order) Cancel(reason string) error {
	if o.Status == OrderStatusDelivered || o.Status == OrderStatusCancelled {
		return errors.New(errors.ErrOrderNotCancellable, "order cannot be cancelled")
	}
	o.Status = OrderStatusCancelled
	o.UpdatedAt = time.Now()
	o.Version++
	return nil
}

// UpdateStatus updates the order status
func (o *Order) UpdateStatus(newStatus OrderStatus) error {
	// Validate status transitions
	validTransitions := map[OrderStatus][]OrderStatus{
		OrderStatusPending:    {OrderStatusConfirmed, OrderStatusCancelled},
		OrderStatusConfirmed:  {OrderStatusProcessing, OrderStatusCancelled},
		OrderStatusProcessing: {OrderStatusShipped, OrderStatusCancelled},
		OrderStatusShipped:    {OrderStatusDelivered},
		OrderStatusDelivered:  {OrderStatusRefunded},
	}

	allowedStatuses, ok := validTransitions[o.Status]
	if !ok {
		return errors.New(errors.ErrInvalidInput, "no valid transitions from current status")
	}

	isValid := false
	for _, allowed := range allowedStatuses {
		if newStatus == allowed {
			isValid = true
			break
		}
	}

	if !isValid {
		return errors.New(errors.ErrInvalidInput, "invalid status transition")
	}

	o.Status = newStatus
	o.UpdatedAt = time.Now()
	o.Version++
	return nil
}
