package postgres

import (
	"context"
	"database/sql"
	"time"

	"github.com/titan-commerce/backend/payment-service/internal/domain"
	"github.com/titan-commerce/backend/pkg/errors"
	"github.com/titan-commerce/backend/pkg/logger"
	_ "github.com/lib/pq"
)

type PaymentRepository struct {
	db     *sql.DB
	logger *logger.Logger
}

func NewPaymentRepository(databaseURL string, logger *logger.Logger) (*PaymentRepository, error) {
	db, err := sql.Open("postgres", databaseURL)
	if err != nil {
		return nil, errors.Wrap(errors.ErrInternal, "failed to connect to database", err)
	}

	// Test connection
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	
	if err := db.PingContext(ctx); err != nil {
		return nil, errors.Wrap(errors.ErrInternal, "failed to ping database", err)
	}

	logger.Info("Payment PostgreSQL repository initialized")
	return &PaymentRepository{db: db, logger: logger}, nil
}

func (r *PaymentRepository) Save(ctx context.Context, payment *domain.Payment) error {
	query := `
		INSERT INTO payments (
			id, order_id, user_id, amount, currency, gateway, status,
			gateway_transaction_id, idempotency_key, created_at, updated_at, version
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)
	`

	_, err := r.db.ExecContext(ctx, query,
		payment.ID, payment.OrderID, payment.UserID, payment.Amount, payment.Currency,
		payment.Gateway, payment.Status, payment.GatewayTransactionID, payment.IdempotencyKey,
		payment.CreatedAt, payment.UpdatedAt, payment.Version,
	)

	if err != nil {
		return errors.Wrap(errors.ErrInternal, "failed to save payment", err)
	}

	return nil
}

func (r *PaymentRepository) FindByID(ctx context.Context, paymentID string) (*domain.Payment, error) {
	query := `
		SELECT id, order_id, user_id, amount, currency, gateway, status,
			   gateway_transaction_id, idempotency_key, created_at, updated_at, version
		FROM payments WHERE id = $1
	`

	var payment domain.Payment
	err := r.db.QueryRowContext(ctx, query, paymentID).Scan(
		&payment.ID, &payment.OrderID, &payment.UserID, &payment.Amount, &payment.Currency,
		&payment.Gateway, &payment.Status, &payment.GatewayTransactionID, &payment.IdempotencyKey,
		&payment.CreatedAt, &payment.UpdatedAt, &payment.Version,
	)

	if err == sql.ErrNoRows {
		return nil, errors.New(errors.ErrNotFound, "payment not found")
	}
	if err != nil {
		return nil, errors.Wrap(errors.ErrInternal, "failed to find payment", err)
	}

	return &payment, nil
}

func (r *PaymentRepository) FindByIdempotencyKey(ctx context.Context, idempotencyKey string) (*domain.Payment, error) {
	query := `
		SELECT id, order_id, user_id, amount, currency, gateway, status,
			   gateway_transaction_id, idempotency_key, created_at, updated_at, version
		FROM payments WHERE idempotency_key = $1
	`

	var payment domain.Payment
	err := r.db.QueryRowContext(ctx, query, idempotencyKey).Scan(
		&payment.ID, &payment.OrderID, &payment.UserID, &payment.Amount, &payment.Currency,
		&payment.Gateway, &payment.Status, &payment.GatewayTransactionID, &payment.IdempotencyKey,
		&payment.CreatedAt, &payment.UpdatedAt, &payment.Version,
	)

	if err == sql.ErrNoRows {
		return nil, errors.New(errors.ErrNotFound, "payment not found")
	}
	if err != nil {
		return nil, errors.Wrap(errors.ErrInternal, "failed to find payment", err)
	}

	return &payment, nil
}

func (r *PaymentRepository) FindByOrderID(ctx context.Context, orderID string) (*domain.Payment, error) {
	query := `
		SELECT id, order_id, user_id, amount, currency, gateway, status,
			   gateway_transaction_id, idempotency_key, created_at, updated_at, version
		FROM payments WHERE order_id = $1
	`

	var payment domain.Payment
	err := r.db.QueryRowContext(ctx, query, orderID).Scan(
		&payment.ID, &payment.OrderID, &payment.UserID, &payment.Amount, &payment.Currency,
		&payment.Gateway, &payment.Status, &payment.GatewayTransactionID, &payment.IdempotencyKey,
		&payment.CreatedAt, &payment.UpdatedAt, &payment.Version,
	)

	if err == sql.ErrNoRows {
		return nil, errors.New(errors.ErrNotFound, "payment not found")
	}
	if err != nil {
		return nil, errors.Wrap(errors.ErrInternal, "failed to find payment", err)
	}

	return &payment, nil
}

func (r *PaymentRepository) Update(ctx context.Context, payment *domain.Payment) error {
	query := `
		UPDATE payments
		SET status = $1, gateway_transaction_id = $2, updated_at = $3, version = $4
		WHERE id = $5 AND version = $6
	`

	result, err := r.db.ExecContext(ctx, query,
		payment.Status, payment.GatewayTransactionID, payment.UpdatedAt,
		payment.Version, payment.ID, payment.Version-1,
	)

	if err != nil {
		return errors.Wrap(errors.ErrInternal, "failed to update payment", err)
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return errors.Wrap(errors.ErrInternal, "failed to get rows affected", err)
	}

	if rows == 0 {
		return errors.New(errors.ErrConflict, "payment was modified by another transaction (optimistic lock)")
	}

	return nil
}

func (r *PaymentRepository) Close() error {
	return r.db.Close()
}
