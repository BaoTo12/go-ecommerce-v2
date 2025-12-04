package mock

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/titan-commerce/backend/payment-service/internal/domain"
	"github.com/titan-commerce/backend/pkg/logger"
)

// MockPaymentGateway simulates payment gateway for testing
type MockPaymentGateway struct {
	logger *logger.Logger
}

func NewMockPaymentGateway(logger *logger.Logger) *MockPaymentGateway {
	return &MockPaymentGateway{logger: logger}
}

func (g *MockPaymentGateway) ProcessPayment(ctx context.Context, payment *domain.Payment, paymentMethodID string) (string, string, error) {
	// Simulate payment processing
	gatewayTxnID := "mock_txn_" + uuid.New().String()[:8]
	clientSecret := "mock_secret_" + uuid.New().String()[:12]

	g.logger.Infof("MOCK: Processing payment %s for $%.2f via mock gateway", payment.ID, payment.Amount)
	
	// In production, this would call Stripe/PayPal/Adyen API
	// For now, we simulate success
	
	return gatewayTxnID, clientSecret, nil
}

func (g *MockPaymentGateway) RefundPayment(ctx context.Context, gatewayTransactionID string, amount float64) (string, error) {
	refundID := "mock_refund_" + uuid.New().String()[:8]
	
	g.logger.Infof("MOCK: Refunding transaction %s for $%.2f", gatewayTransactionID, amount)
	
	return refundID, nil
}

func (g *MockPaymentGateway) GetPaymentStatus(ctx context.Context, gatewayTransactionID string) (string, error) {
	return "completed", nil
}
