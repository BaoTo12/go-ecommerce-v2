package application

import (
	"context"
	"fmt"
	"time"

	"github.com/titan-commerce/backend/checkout-service/internal/domain"
	"github.com/titan-commerce/backend/pkg/logger"
)

// Service Interfaces for external dependencies
type InventoryClient interface {
	ReserveStock(ctx context.Context, productIDs []string) (string, error)
	CommitReservation(ctx context.Context, reservationID string) error
	RollbackReservation(ctx context.Context, reservationID string) error
}

type PaymentClient interface {
	ProcessPayment(ctx context.Context, userID string, amount float64, paymentMethodID string) (string, error)
	RefundPayment(ctx context.Context, paymentID string) error
}

type OrderClient interface {
	CreateOrder(ctx context.Context, userID string, productIDs []string, shippingAddress string) (string, error)
	CancelOrder(ctx context.Context, orderID string) error
}

type CartClient interface {
	GetCart(ctx context.Context, userID string) (float64, []string, error)
	ClearCart(ctx context.Context, userID string) error
}

type CheckoutService struct {
	inventory InventoryClient
	payment   PaymentClient
	order     OrderClient
	cart      CartClient
	sessions  map[string]*domain.CheckoutSession // In-memory storage for simplicity, replace with Redis/DB
	logger    *logger.Logger
}

func NewCheckoutService(inv InventoryClient, pay PaymentClient, ord OrderClient, crt CartClient, logger *logger.Logger) *CheckoutService {
	return &CheckoutService{
		inventory: inv,
		payment:   pay,
		order:     ord,
		cart:      crt,
		sessions:  make(map[string]*domain.CheckoutSession),
		logger:    logger,
	}
}

func (s *CheckoutService) InitiateCheckout(ctx context.Context, userID, shippingAddress, paymentMethodID string) (*domain.CheckoutSession, error) {
	// 1. Get Cart
	totalAmount, productIDs, err := s.cart.GetCart(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get cart: %w", err)
	}

	if len(productIDs) == 0 {
		return nil, fmt.Errorf("cart is empty")
	}

	// 2. Create Session
	session := domain.NewCheckoutSession(userID, shippingAddress, paymentMethodID, productIDs, totalAmount)
	s.sessions[session.SessionID] = session

	// 3. Start Saga (Async)
	go s.runSaga(session)

	return session, nil
}

func (s *CheckoutService) GetCheckoutStatus(ctx context.Context, sessionID string) (*domain.CheckoutSession, error) {
	session, ok := s.sessions[sessionID]
	if !ok {
		return nil, fmt.Errorf("session not found")
	}
	return session, nil
}

func (s *CheckoutService) CancelCheckout(ctx context.Context, sessionID string) error {
	session, ok := s.sessions[sessionID]
	if !ok {
		return fmt.Errorf("session not found")
	}

	// Trigger compensation if in progress
	if session.Status != domain.CheckoutStatusCompleted && session.Status != domain.CheckoutStatusFailed {
		session.MarkCompensating()
		// In a real system, this would signal the running saga or trigger a compensation workflow
		// For this simple implementation, we'll just mark it. The running saga would need to check this status.
	}

	return nil
}

func (s *CheckoutService) runSaga(session *domain.CheckoutSession) {
	ctx := context.Background()
	s.logger.Infof("Starting Saga for session %s", session.SessionID)

	// Step 1: Reserve Inventory
	session.MarkReservingInventory()
	reservationID, err := s.inventory.ReserveStock(ctx, session.ProductIDs)
	if err != nil {
		s.failSaga(session, "Inventory reservation failed: "+err.Error())
		return
	}
	session.ReservationID = reservationID

	// Step 2: Process Payment
	session.MarkProcessingPayment(reservationID)
	paymentID, err := s.payment.ProcessPayment(ctx, session.UserID, session.TotalAmount, session.PaymentMethodID)
	if err != nil {
		s.compensateInventory(ctx, session)
		s.failSaga(session, "Payment failed: "+err.Error())
		return
	}
	session.PaymentID = paymentID

	// Step 3: Create Order
	session.MarkCreatingOrder(paymentID)
	orderID, err := s.order.CreateOrder(ctx, session.UserID, session.ProductIDs, session.ShippingAddress)
	if err != nil {
		s.compensatePayment(ctx, session)
		s.compensateInventory(ctx, session)
		s.failSaga(session, "Order creation failed: "+err.Error())
		return
	}

	// Step 4: Commit Inventory & Clear Cart
	if err := s.inventory.CommitReservation(ctx, reservationID); err != nil {
		s.logger.Error(err, "Failed to commit reservation (critical)")
		// Manual intervention might be needed here
	}
	
	if err := s.cart.ClearCart(ctx, session.UserID); err != nil {
		s.logger.Error(err, "Failed to clear cart")
	}

	session.MarkCompleted(orderID)
	s.logger.Infof("Saga completed successfully for session %s", session.SessionID)
}

func (s *CheckoutService) failSaga(session *domain.CheckoutSession, reason string) {
	session.MarkFailed(reason)
	s.logger.Errorf("Saga failed for session %s: %s", session.SessionID, reason)
}

func (s *CheckoutService) compensateInventory(ctx context.Context, session *domain.CheckoutSession) {
	if session.ReservationID != "" {
		s.logger.Infof("Compensating: Rolling back inventory for session %s", session.SessionID)
		if err := s.inventory.RollbackReservation(ctx, session.ReservationID); err != nil {
			s.logger.Error(err, "Failed to rollback inventory")
		}
	}
}

func (s *CheckoutService) compensatePayment(ctx context.Context, session *domain.CheckoutSession) {
	if session.PaymentID != "" {
		s.logger.Infof("Compensating: Refunding payment for session %s", session.SessionID)
		if err := s.payment.RefundPayment(ctx, session.PaymentID); err != nil {
			s.logger.Error(err, "Failed to refund payment")
		}
	}
}
