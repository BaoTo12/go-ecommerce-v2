package domain

import "context"

// Repository defines the payment persistence interface
type Repository interface {
	Save(ctx context.Context, payment *Payment) error
	FindByID(ctx context.Context, paymentID string) (*Payment, error)
	FindByOrderID(ctx context.Context, orderID string) (*Payment, error)
	FindByIdempotencyKey(ctx context.Context, key string) (*Payment, error)
	Update(ctx context.Context, payment *Payment) error
}

// PaymentGatewayProvider defines the interface for payment gateway integrations
type PaymentGatewayProvider interface {
	ProcessPayment(ctx context.Context, payment *Payment, paymentMethodID string) (gatewayTransactionID string, clientSecret string, err error)
	RefundPayment(ctx context.Context, gatewayTransactionID string, amount float64) (refundID string, err error)
	VerifyPayment(ctx context.Context, gatewayTransactionID string) (verified bool, err error)
}
