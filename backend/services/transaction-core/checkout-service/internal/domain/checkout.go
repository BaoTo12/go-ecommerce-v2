package domain

import (
	"time"

	"github.com/google/uuid"
)

type CheckoutStatus string

const (
	CheckoutStatusInitiated          CheckoutStatus = "INITIATED"
	CheckoutStatusReservingInventory CheckoutStatus = "RESERVING_INVENTORY"
	CheckoutStatusProcessingPayment  CheckoutStatus = "PROCESSING_PAYMENT"
	CheckoutStatusCreatingOrder      CheckoutStatus = "CREATING_ORDER"
	CheckoutStatusCompleted          CheckoutStatus = "COMPLETED"
	CheckoutStatusFailed             CheckoutStatus = "FAILED"
	CheckoutStatusCompensating       CheckoutStatus = "COMPENSATING"
)

// CheckoutSession represents a Saga instance
type CheckoutSession struct {
	SessionID       string
	UserID          string
	ProductIDs      []string
	TotalAmount     float64
	ShippingAddress string
	PaymentMethodID string
	Status          CheckoutStatus
	ErrorMessage    string
	OrderID         string
	PaymentID       string
	ReservationID   string
	CreatedAt       time.Time
	UpdatedAt       time.Time
}

func NewCheckoutSession(userID, shippingAddress, paymentMethodID string, productIDs []string, totalAmount float64) *CheckoutSession {
	now := time.Now()
	return &CheckoutSession{
		SessionID:       uuid.New().String(),
		UserID:          userID,
		ProductIDs:      productIDs,
		TotalAmount:     totalAmount,
		ShippingAddress: shippingAddress,
		PaymentMethodID: paymentMethodID,
		Status:          CheckoutStatusInitiated,
		CreatedAt:       now,
		UpdatedAt:       now,
	}
}

func (s *CheckoutSession) MarkReservingInventory() {
	s.Status = CheckoutStatusReservingInventory
	s.UpdatedAt = time.Now()
}

func (s *CheckoutSession) MarkProcessingPayment(reservationID string) {
	s.Status = CheckoutStatusProcessingPayment
	s.ReservationID = reservationID
	s.UpdatedAt = time.Now()
}

func (s *CheckoutSession) MarkCreatingOrder(paymentID string) {
	s.Status = CheckoutStatusCreatingOrder
	s.PaymentID = paymentID
	s.UpdatedAt = time.Now()
}

func (s *CheckoutSession) MarkCompleted(orderID string) {
	s.Status = CheckoutStatusCompleted
	s.OrderID = orderID
	s.UpdatedAt = time.Now()
}

func (s *CheckoutSession) MarkFailed(errorMessage string) {
	s.Status = CheckoutStatusFailed
	s.ErrorMessage = errorMessage
	s.UpdatedAt = time.Now()
}

func (s *CheckoutSession) MarkCompensating() {
	s.Status = CheckoutStatusCompensating
	s.UpdatedAt = time.Now()
}
