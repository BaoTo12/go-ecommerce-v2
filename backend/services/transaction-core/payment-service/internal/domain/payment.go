package domain

import (
	"time"

	"github.com/google/uuid"
	"github.com/titan-commerce/backend/pkg/errors"
)

type PaymentStatus string

const (
	PaymentStatusPending    PaymentStatus = "PENDING"
	PaymentStatusProcessing PaymentStatus = "PROCESSING"
	PaymentStatusCompleted  PaymentStatus = "COMPLETED"
	PaymentStatusFailed     PaymentStatus = "FAILED"
	PaymentStatusRefunded   PaymentStatus = "REFUNDED"
)

type PaymentGateway string

const (
	PaymentGatewayStripe PaymentGateway = "STRIPE"
	PaymentGatewayPayPal PaymentGateway = "PAYPAL"
	PaymentGatewayAdyen  PaymentGateway = "ADYEN"
)

// Payment is the aggregate root
type Payment struct {
	ID                   string
	OrderID              string
	UserID               string
	Amount               float64
	Currency             string
	Gateway              PaymentGateway
	Status               PaymentStatus
	GatewayTransactionID string
	IdempotencyKey       string
	CreatedAt            time.Time
	UpdatedAt            time.Time
	Version              int
}

// NewPayment creates a new payment (Factory method)
func NewPayment(orderID, userID string, amount float64, currency string, gateway PaymentGateway, idempotencyKey string) (*Payment, error) {
	if orderID == "" {
		return nil, errors.New(errors.ErrInvalidInput, "order ID is required")
	}
	if userID == "" {
		return nil, errors.New(errors.ErrInvalidInput, "user ID is required")
	}
	if amount <= 0 {
		return nil, errors.New(errors.ErrInvalidInput, "amount must be positive")
	}
	if currency == "" {
		currency = "USD"
	}
	if idempotencyKey == "" {
		return nil, errors.New(errors.ErrInvalidInput, "idempotency key is required")
	}

	now := time.Now()
	return &Payment{
		ID:             uuid.New().String(),
		OrderID:        orderID,
		UserID:         userID,
		Amount:         amount,
		Currency:       currency,
		Gateway:        gateway,
		Status:         PaymentStatusPending,
		IdempotencyKey: idempotencyKey,
		CreatedAt:      now,
		UpdatedAt:      now,
		Version:        1,
	}, nil
}

// MarkProcessing transitions payment to processing state
func (p *Payment) MarkProcessing(gatewayTransactionID string) error {
	if p.Status != PaymentStatusPending {
		return errors.New(errors.ErrInvalidInput, "only pending payments can be marked as processing")
	}
	p.Status = PaymentStatusProcessing
	p.GatewayTransactionID = gatewayTransactionID
	p.UpdatedAt = time.Now()
	p.Version++
	return nil
}

// MarkCompleted marks payment as successfully completed
func (p *Payment) MarkCompleted() error {
	if p.Status != PaymentStatusProcessing {
		return errors.New(errors.ErrInvalidInput, "only processing payments can be completed")
	}
	p.Status = PaymentStatusCompleted
	p.UpdatedAt = time.Now()
	p.Version++
	return nil
}

// MarkFailed marks payment as failed
func (p *Payment) MarkFailed(reason string) error {
	if p.Status == PaymentStatusCompleted || p.Status == PaymentStatusRefunded {
		return errors.New(errors.ErrInvalidInput, "cannot fail completed or refunded payment")
	}
	p.Status = PaymentStatusFailed
	p.UpdatedAt = time.Now()
	p.Version++
	return nil
}

// Refund issues a refund for this payment
func (p *Payment) Refund() error {
	if p.Status != PaymentStatusCompleted {
		return errors.New(errors.ErrInvalidInput, "only completed payments can be refunded")
	}
	p.Status = PaymentStatusRefunded
	p.UpdatedAt = time.Now()
	p.Version++
	return nil
}
