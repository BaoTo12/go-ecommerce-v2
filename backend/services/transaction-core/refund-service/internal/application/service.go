package application

import (
	"context"

	"github.com/titan-commerce/backend/refund-service/internal/domain"
	"github.com/titan-commerce/backend/pkg/logger"
)

type RefundRepository interface {
	Save(ctx context.Context, refund *domain.Refund) error
	FindByID(ctx context.Context, refundID string) (*domain.Refund, error)
	Update(ctx context.Context, refund *domain.Refund) error
}

type PaymentGateway interface {
	ProcessRefund(ctx context.Context, paymentID string, amount float64) (string, error)
}

type RefundService struct {
	repo    RefundRepository
	gateway PaymentGateway
	logger  *logger.Logger
}

func NewRefundService(repo RefundRepository, gateway PaymentGateway, logger *logger.Logger) *RefundService {
	return &RefundService{
		repo:    repo,
		gateway: gateway,
		logger:  logger,
	}
}

// ProcessRefund initiates a refund (Command)
func (s *RefundService) ProcessRefund(ctx context.Context, paymentID, orderID string, amount float64, reason string) (*domain.Refund, error) {
	refund, err := domain.NewRefund(paymentID, orderID, amount, reason)
	if err != nil {
		return nil, err
	}

	if err := s.repo.Save(ctx, refund); err != nil {
		s.logger.Error(err, "failed to save refund")
		return nil, err
	}

	// Process refund via payment gateway
	gatewayRefundID, err := s.gateway.ProcessRefund(ctx, paymentID, amount)
	if err != nil {
		refund.Fail()
		s.repo.Update(ctx, refund)
		s.logger.Error(err, "gateway refund failed")
		return refund, err
	}

	refund.Process(gatewayRefundID)
	refund.Complete()

	if err := s.repo.Update(ctx, refund); err != nil {
		return refund, err
	}

	s.logger.Infof("Refund processed: refund=%s, payment=%s, amount=%.2f", refund.ID, paymentID, amount)
	return refund, nil
}

// GetRefund retrieves refund status (Query)
func (s *RefundService) GetRefund(ctx context.Context, refundID string) (*domain.Refund, error) {
	return s.repo.FindByID(ctx, refundID)
}
