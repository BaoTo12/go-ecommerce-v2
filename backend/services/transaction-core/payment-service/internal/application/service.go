package application

import (
	"context"

	"github.com/titan-commerce/backend/payment-service/internal/domain"
	"github.com/titan-commerce/backend/pkg/errors"
	"github.com/titan-commerce/backend/pkg/logger"
)

type PaymentService struct {
	repo     domain.Repository
	gateways map[domain.PaymentGateway]domain.PaymentGatewayProvider
	logger   *logger.Logger
}

func NewPaymentService(repo domain.Repository, gateways map[domain.PaymentGateway]domain.PaymentGatewayProvider, logger *logger.Logger) *PaymentService {
	return &PaymentService{
		repo:     repo,
		gateways: gateways,
		logger:   logger,
	}
}

// ProcessPayment processes a payment (Command)
func (s *PaymentService) ProcessPayment(ctx context.Context, orderID, userID string, amount float64, currency string, gatewayType domain.PaymentGateway, paymentMethodID, idempotencyKey string) (*domain.Payment, string, error) {
	// Check idempotency - prevent duplicate payments
	existing, err := s.repo.FindByIdempotencyKey(ctx, idempotencyKey)
	if err == nil && existing != nil {
		s.logger.Infof("Idempotent request detected: %s", idempotencyKey)
		return existing, "", nil
	}

	// Create payment aggregate
	payment, err := domain.NewPayment(orderID, userID, amount, currency, gatewayType, idempotencyKey)
	if err != nil {
		s.logger.Error(err, "failed to create payment")
		return nil, "", err
	}

	// Save payment (pending state)
	if err := s.repo.Save(ctx, payment); err != nil {
		s.logger.Error(err, "failed to save payment")
		return nil, "", err
	}

	// Get payment gateway
	gateway, ok := s.gateways[gatewayType]
	if !ok {
		return nil, "", errors.New(errors.ErrInvalidInput, "unsupported payment gateway")
	}

	// Process payment through gateway
	gatewayTxnID, clientSecret, err := gateway.ProcessPayment(ctx, payment, paymentMethodID)
	if err != nil {
		payment.MarkFailed(err.Error())
		s.repo.Update(ctx, payment)
		s.logger.Error(err, "gateway payment failed")
		return nil, "", errors.Wrap(errors.ErrPaymentFailed, "payment processing failed", err)
	}

	// Mark as processing
	if err := payment.MarkProcessing(gatewayTxnID); err != nil {
		return nil, "", err
	}

	// Update payment
	if err := s.repo.Update(ctx, payment); err != nil {
		return nil, "", err
	}

	s.logger.Infof("Payment processed: %s for order: %s, gateway txn: %s", payment.ID, orderID, gatewayTxnID)
	return payment, clientSecret, nil
}

// GetPayment retrieves a payment (Query)
func (s *PaymentService) GetPayment(ctx context.Context, paymentID string) (*domain.Payment, error) {
	payment, err := s.repo.FindByID(ctx, paymentID)
	if err != nil {
		s.logger.Error(err, "failed to get payment")
		return nil, err
	}
	return payment, nil
}

// RefundPayment issues a refund (Command)
func (s *PaymentService) RefundPayment(ctx context.Context, paymentID string, amount float64, reason string) (string, error) {
	payment, err := s.repo.FindByID(ctx, paymentID)
	if err != nil {
		return "", err
	}

	// Get gateway
	gateway, ok := s.gateways[payment.Gateway]
	if !ok {
		return "", errors.New(errors.ErrInternal, "gateway not found")
	}

	// Process refund through gateway
	refundID, err := gateway.RefundPayment(ctx, payment.GatewayTransactionID, amount)
	if err != nil {
		s.logger.Error(err, "refund failed")
		return "", err
	}

	// Mark payment as refunded
	if err := payment.Refund(); err != nil {
		return "", err
	}

	if err := s.repo.Update(ctx, payment); err != nil {
		return "", err
	}

	s.logger.Infof("Payment refunded: %s, refund ID: %s", paymentID, refundID)
	return refundID, nil
}
