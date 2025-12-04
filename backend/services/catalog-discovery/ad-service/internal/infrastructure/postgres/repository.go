package postgres

import (
	"context"
	"database/sql"
	"time"

	"github.com/titan-commerce/backend/ad-service/internal/domain"
	"github.com/titan-commerce/backend/pkg/errors"
	"github.com/titan-commerce/backend/pkg/logger"
	_ "github.com/lib/pq"
)

type AdRepository struct {
	db     *sql.DB
	logger *logger.Logger
}

func NewAdRepository(databaseURL string, logger *logger.Logger) (*AdRepository, error) {
	db, err := sql.Open("postgres", databaseURL)
	if err != nil {
		return nil, errors.Wrap(errors.ErrInternal, "failed to connect to database", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	
	if err := db.PingContext(ctx); err != nil {
		return nil, errors.Wrap(errors.ErrInternal, "failed to ping database", err)
	}

	logger.Info("Ad PostgreSQL repository initialized")
	return &AdRepository{db: db, logger: logger}, nil
}

func (r *AdRepository) SaveCampaign(ctx context.Context, c *domain.Campaign) error {
	query := `
		INSERT INTO campaigns (
			id, seller_id, product_id, budget, remaining_budget, 
			bid_amount, status, start_time, end_time, created_at
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
	`
	_, err := r.db.ExecContext(ctx, query,
		c.ID, c.SellerID, c.ProductID, c.Budget, c.RemainingBudget,
		c.BidAmount, c.Status, c.StartTime, c.EndTime, c.CreatedAt,
	)
	if err != nil {
		return errors.Wrap(errors.ErrInternal, "failed to save campaign", err)
	}
	return nil
}

func (r *AdRepository) GetActiveCampaigns(ctx context.Context) ([]*domain.Campaign, error) {
	query := `
		SELECT id, seller_id, product_id, budget, remaining_budget, 
			   bid_amount, status, start_time, end_time, created_at
		FROM campaigns
		WHERE status = 'active' AND remaining_budget > 0
		LIMIT 10
	`
	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, errors.Wrap(errors.ErrInternal, "failed to get campaigns", err)
	}
	defer rows.Close()

	var campaigns []*domain.Campaign
	for rows.Next() {
		var c domain.Campaign
		if err := rows.Scan(
			&c.ID, &c.SellerID, &c.ProductID, &c.Budget, &c.RemainingBudget,
			&c.BidAmount, &c.Status, &c.StartTime, &c.EndTime, &c.CreatedAt,
		); err != nil {
			return nil, errors.Wrap(errors.ErrInternal, "failed to scan campaign", err)
		}
		campaigns = append(campaigns, &c)
	}
	return campaigns, nil
}

func (r *AdRepository) TrackEvent(ctx context.Context, e *domain.AdEvent) error {
	query := `
		INSERT INTO ad_events (ad_id, user_id, event_type, timestamp)
		VALUES ($1, $2, $3, $4)
	`
	_, err := r.db.ExecContext(ctx, query, e.AdID, e.UserID, e.EventType, e.Timestamp)
	if err != nil {
		return errors.Wrap(errors.ErrInternal, "failed to track event", err)
	}
	return nil
}

func (r *AdRepository) DeductBudget(ctx context.Context, campaignID string, amount float64) error {
	query := `
		UPDATE campaigns 
		SET remaining_budget = remaining_budget - $1 
		WHERE id = $2 AND remaining_budget >= $1
	`
	res, err := r.db.ExecContext(ctx, query, amount, campaignID)
	if err != nil {
		return errors.Wrap(errors.ErrInternal, "failed to deduct budget", err)
	}
	rows, _ := res.RowsAffected()
	if rows == 0 {
		return errors.New(errors.ErrInternal, "insufficient budget or campaign not found")
	}
	return nil
}
