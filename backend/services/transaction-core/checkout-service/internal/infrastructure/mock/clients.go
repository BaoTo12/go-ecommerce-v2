package mock

import (
	"context"
	"github.com/google/uuid"
)

// Mock Clients for Checkout Service

type MockInventoryClient struct{}
func (m *MockInventoryClient) ReserveStock(ctx context.Context, productIDs []string) (string, error) {
	return "res_" + uuid.New().String(), nil
}
func (m *MockInventoryClient) CommitReservation(ctx context.Context, reservationID string) error { return nil }
func (m *MockInventoryClient) RollbackReservation(ctx context.Context, reservationID string) error { return nil }

type MockPaymentClient struct{}
func (m *MockPaymentClient) ProcessPayment(ctx context.Context, userID string, amount float64, paymentMethodID string) (string, error) {
	return "pay_" + uuid.New().String(), nil
}
func (m *MockPaymentClient) RefundPayment(ctx context.Context, paymentID string) error { return nil }

type MockOrderClient struct{}
func (m *MockOrderClient) CreateOrder(ctx context.Context, userID string, productIDs []string, shippingAddress string) (string, error) {
	return "ord_" + uuid.New().String(), nil
}
func (m *MockOrderClient) CancelOrder(ctx context.Context, orderID string) error { return nil }

type MockCartClient struct{}
func (m *MockCartClient) GetCart(ctx context.Context, userID string) (float64, []string, error) {
	return 100.0, []string{"prod_1", "prod_2"}, nil
}
func (m *MockCartClient) ClearCart(ctx context.Context, userID string) error { return nil }
