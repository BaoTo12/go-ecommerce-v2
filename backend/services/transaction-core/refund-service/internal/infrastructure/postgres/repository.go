package postgres

import (
	"context"
	"database/sql"
	"time"

	"github.com/titan-commerce/backend/refund-service/internal/domain"
	"github.com/titan-commerce/backend/pkg/errors"
	"github.com/titan-commerce/backend/pkg/logger"
	_ "github.com/lib/pq"
)

type RefundRepository struct {
	db     *sql.DB
	logger *logger.Logger
}

func NewRefundRepository(databaseURL string, logger *logger.Logger) (*RefundRepository, error) {
	db, err := sql.Open("postgres", databaseURL)
	if err != nil {
		return nil, errors.Wrap(errors.ErrInternal, "failed to connect to database", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	
	if err := db.PingContext(ctx); err != nil {
		return nil, errors.Wrap(errors.ErrInternal, "failed to ping database", err)
	}

	logger.Info("Refund PostgreSQL repository initialized")
	return &RefundRepository{db: db, logger: logger}, nil
}

func (r *RefundRepository) Save(ctx context.Context, refund *domain.Refund) error {
	query := `
		INSERT INTO refunds (id, payment_id, order_id, amount, reason, status, gateway_refund_id, created_at, processed_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
	`

	_, err := r.db.ExecContext(ctx, query,
		refund.ID, refund.PaymentID, refund.OrderID, refund.Amount, refund.Reason,
		refund.Status, refund.GatewayRefundID, refund.CreatedAt, refund.ProcessedAt)

	if err != nil {
		return errors.Wrap(errors.ErrInternal, "failed to save refund", err)
	}

	return nil
}

func (r *RefundRepository) FindByID(ctx context.Context, refundID string) (*domain.Refund, error) {
	query := `
		SELECT id, payment_id, order_id, amount, reason, status, gateway_refund_id, created_at, processed_at
		FROM refunds WHERE id = $1
	`

	var refund domain.Refund
	err := r.db.QueryRowContext(ctx, query, refundID).Scan(
		&refund.ID, &refund.PaymentID, &refund.OrderID, &refund.Amount, &refund.Reason,
		&refund.Status, &refund.GatewayRefundID, &refund.CreatedAt, &refund.ProcessedAt)

	if err == sql.ErrNoRows {
		return nil, errors.New(errors.ErrNotFound, "refund not found")
	}
	if err != nil {
		return nil, errors.Wrap(errors.ErrInternal, "failed to find refund", err)
	}

	return &refund, nil
}

func (r *RefundRepository) Update(ctx context.Context, refund *domain.Refund) error {
	query := `
		UPDATE refunds
		SET status = $1, gateway_refund_id = $2, processed_at = $3
		WHERE id = $4
	`

	_, err := r.db.ExecContext(ctx, query,
		refund.Status, refund.GatewayRefundID, refund.ProcessedAt, refund.ID)

	if err != nil {
		return errors.Wrap(errors.ErrInternal, "failed to update refund", err)
	}

	return nil
}
