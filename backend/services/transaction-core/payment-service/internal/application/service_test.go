package application_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/titan-commerce/backend/payment-service/internal/application"
	"github.com/titan-commerce/backend/payment-service/internal/domain"
	"github.com/titan-commerce/backend/pkg/logger"
)

type MockPaymentRepository struct {
	mock.Mock
}

func (m *MockPaymentRepository) Save(ctx context.Context, payment *domain.Payment) error {
	args := m.Called(ctx, payment)
	return args.Error(0)
}

func (m *MockPaymentRepository) FindByID(ctx context.Context, paymentID string) (*domain.Payment, error) {
	args := m.Called(ctx, paymentID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Payment), args.Error(1)
}

func (m *MockPaymentRepository) FindByIdempotencyKey(ctx context.Context, key string) (*domain.Payment, error) {
	args := m.Called(ctx, key)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Payment), args.Error(1)
}

func (m *MockPaymentRepository) Update(ctx context.Context, payment *domain.Payment) error {
	args := m.Called(ctx, payment)
	return args.Error(0)
}

type MockPaymentGateway struct {
	mock.Mock
}

func (m *MockPaymentGateway) ProcessPayment(ctx context.Context, payment *domain.Payment, paymentMethodID string) (string, string, error) {
	args := m.Called(ctx, payment, paymentMethodID)
	return args.String(0), args.String(1), args.Error(2)
}

func (m *MockPaymentGateway) RefundPayment(ctx context.Context, transactionID string, amount float64) (string, error) {
	args := m.Called(ctx, transactionID, amount)
	return args.String(0), args.Error(1)
}

func TestPaymentService_ProcessPayment(t *testing.T) {
	// Setup
	mockRepo := new(MockPaymentRepository)
	mockGateway := new(MockPaymentGateway)
	log := logger.New(logger.Config{Level: "debug", ServiceName: "test"})

	gateways := map[domain.PaymentGateway]domain.PaymentGatewayProvider{
		domain.GatewayStripe: mockGateway,
	}
	
	service := application.NewPaymentService(mockRepo, gateways, log)

	ctx := context.Background()
	orderID := "order-123"
	userID := "user-123"
	amount := 99.99
	currency := "USD"
	gatewayType := domain.GatewayStripe
	paymentMethodID := "pm_123"
	idempotencyKey := "idem-123"

	// Expectations
	mockRepo.On("FindByIdempotencyKey", ctx, idempotencyKey).Return(nil, assert.AnError)
	mockRepo.On("Save", ctx, mock.AnythingOfType("*domain.Payment")).Return(nil)
	mockGateway.On("ProcessPayment", ctx, mock.AnythingOfType("*domain.Payment"), paymentMethodID).
		Return("txn-123", "cs_123", nil)
	mockRepo.On("Update", ctx, mock.AnythingOfType("*domain.Payment")).Return(nil)

	// Execute
	payment, clientSecret, err := service.ProcessPayment(ctx, orderID, userID, amount, currency, gatewayType, paymentMethodID, idempotencyKey)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, payment)
	assert.Equal(t, "cs_123", clientSecret)
	assert.Equal(t, orderID, payment.OrderID)
	assert.Equal(t, amount, payment.Amount)
	
	mockRepo.AssertExpectations(t)
	mockGateway.AssertExpectations(t)
}

func TestPaymentService_IdempotencyCheck(t *testing.T) {
	// Setup
	mockRepo := new(MockPaymentRepository)
	log := logger.New(logger.Config{Level: "debug", ServiceName: "test"})
	
	service := application.NewPaymentService(mockRepo, nil, log)

	ctx := context.Background()
	idempotencyKey := "idem-123"
	
	existingPayment := &domain.Payment{
		ID:             "pay-123",
		OrderID:        "order-123",
		Amount:         99.99,
		IdempotencyKey: idempotencyKey,
	}

	// Expectations - return existing payment for idempotency check
	mockRepo.On("FindByIdempotencyKey", ctx, idempotencyKey).Return(existingPayment, nil)

	// Execute
	payment, _, err := service.ProcessPayment(ctx, "order-123", "user-123", 99.99, "USD", domain.GatewayStripe, "pm_123", idempotencyKey)

	// Assert - should return existing payment
	assert.NoError(t, err)
	assert.Equal(t, existingPayment, payment)
	
	mockRepo.AssertExpectations(t)
}

func TestPaymentService_RefundPayment(t *testing.T) {
	// Setup
	mockRepo := new(MockPaymentRepository)
	mockGateway := new(MockPaymentGateway)
	log := logger.New(logger.Config{Level: "debug", ServiceName: "test"})

	gateways := map[domain.PaymentGateway]domain.PaymentGatewayProvider{
		domain.GatewayStripe: mockGateway,
	}
	
	service := application.NewPaymentService(mockRepo, gateways, log)

	ctx := context.Background()
	paymentID := "pay-123"
	refundAmount := 50.00
	reason := "Customer requested refund"

	existingPayment := &domain.Payment{
		ID:                   paymentID,
		Amount:               100.00,
		Status:               domain.PaymentStatusCompleted,
		Gateway:              domain.GatewayStripe,
		GatewayTransactionID: "txn-123",
	}

	// Expectations
	mockRepo.On("FindByID", ctx, paymentID).Return(existingPayment, nil)
	mockGateway.On("RefundPayment", ctx, "txn-123", refundAmount).Return("refund-123", nil)
	mockRepo.On("Update", ctx, mock.AnythingOfType("*domain.Payment")).Return(nil)

	// Execute
	refundID, err := service.RefundPayment(ctx, paymentID, refundAmount, reason)

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, "refund-123", refundID)
	
	mockRepo.AssertExpectations(t)
	mockGateway.AssertExpectations(t)
}
