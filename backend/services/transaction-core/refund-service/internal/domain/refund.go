package domain

import (
	"time"

	"github.com/google/uuid"
	"github.com/titan-commerce/backend/pkg/errors"
)

type Refund struct {
	ID              string
	PaymentID       string
	OrderID         string
	Amount          float64
	Reason          string
	Status          RefundStatus
	GatewayRefundID string
	CreatedAt       time.Time
	ProcessedAt     *time.Time
}

type RefundStatus string

const (
	RefundStatusPending   RefundStatus = "PENDING"
	RefundStatusProcessing RefundStatus = "PROCESSING"
	RefundStatusCompleted RefundStatus = "COMPLETED"
	RefundStatusFailed    RefundStatus = "FAILED"
)

func NewRefund(paymentID, orderID string, amount float64, reason string) (*Refund, error) {
	if paymentID == "" {
		return nil, errors.New(errors.ErrInvalidInput, "payment ID is required")
	}
	if amount <= 0 {
		return nil, errors.New(errors.ErrInvalidInput, "refund amount must be positive")
	}

	return &Refund{
		ID:        uuid.New().String(),
		PaymentID: paymentID,
		OrderID:   orderID,
		Amount:    amount,
		Reason:    reason,
		Status:    RefundStatusPending,
		CreatedAt: time.Now(),
	}, nil
}

func (r *Refund) Process(gatewayRefundID string) {
	r.Status = RefundStatusProcessing
	r.GatewayRefundID = gatewayRefundID
}

func (r *Refund) Complete() {
	r.Status = RefundStatusCompleted
	now := time.Now()
	r.ProcessedAt = &now
}

func (r *Refund) Fail() {
	r.Status = RefundStatusFailed
	now := time.Now()
	r.ProcessedAt = &now
}
