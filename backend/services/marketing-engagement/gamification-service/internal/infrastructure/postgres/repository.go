package postgres

import (
	"context"
	"database/sql"
	"encoding/json"
	"time"

	_ "github.com/lib/pq"
	"github.com/titan-commerce/backend/gamification-service/internal/domain"
	"github.com/titan-commerce/backend/pkg/errors"
	"github.com/titan-commerce/backend/pkg/logger"
)

type GamificationRepository struct {
	db     *sql.DB
	logger *logger.Logger
}

func NewGamificationRepository(databaseURL string, logger *logger.Logger) (*GamificationRepository, error) {
	db, err := sql.Open("postgres", databaseURL)
	if err != nil {
		return nil, errors.Wrap(errors.ErrInternal, "failed to connect to database", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := db.PingContext(ctx); err != nil {
		return nil, errors.Wrap(errors.ErrInternal, "failed to ping database", err)
	}

	logger.Info("Gamification PostgreSQL repository initialized")
	return &GamificationRepository{db: db, logger: logger}, nil
}

func (r *GamificationRepository) GetUserPoints(ctx context.Context, userID string) (*domain.UserPoints, error) {
	query := `SELECT user_id, total_points, available_points, lifetime_earned, lifetime_spent, level, updated_at 
	          FROM user_points WHERE user_id = $1`

	var points domain.UserPoints
	err := r.db.QueryRowContext(ctx, query, userID).Scan(
		&points.UserID, &points.TotalPoints, &points.AvailablePoints, 
		&points.LifetimeEarned, &points.LifetimeSpent, &points.Level, &points.UpdatedAt,
	)
	if err == sql.ErrNoRows {
		// Create new user points entry
		newPoints := domain.NewUserPoints(userID)
		r.SaveUserPoints(ctx, newPoints)
		return newPoints, nil
	}
	return &points, err
}

func (r *GamificationRepository) SaveUserPoints(ctx context.Context, points *domain.UserPoints) error {
	query := `
		INSERT INTO user_points (user_id, total_points, available_points, lifetime_earned, lifetime_spent, level, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		ON CONFLICT (user_id) DO UPDATE 
		SET total_points = $2, available_points = $3, lifetime_earned = $4, lifetime_spent = $5, level = $6, updated_at = $7
	`
	_, err := r.db.ExecContext(ctx, query,
		points.UserID, points.TotalPoints, points.AvailablePoints, points.LifetimeEarned,
		points.LifetimeSpent, points.Level, points.UpdatedAt,
	)
	return err
}

func (r *GamificationRepository) SaveTransaction(ctx context.Context, tx *domain.PointsTransaction) error {
	query := `
		INSERT INTO points_transactions (transaction_id, user_id, points, transaction_type, reference, description, created_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
	`
	_, err := r.db.ExecContext(ctx, query,
		tx.TransactionID, tx.UserID, tx.Points, tx.Type, tx.Reference, tx.Description, tx.CreatedAt,
	)
	return err
}

func (r *GamificationRepository) GetTransactionHistory(ctx context.Context, userID string, limit int) ([]*domain.PointsTransaction, error) {
	query := `SELECT transaction_id, user_id, points, transaction_type, reference, description, created_at 
			  FROM points_transactions WHERE user_id = $1 ORDER BY created_at DESC LIMIT $2`

	rows, err := r.db.QueryContext(ctx, query, userID, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var transactions []*domain.PointsTransaction
	for rows.Next() {
		var tx domain.PointsTransaction
		if err := rows.Scan(&tx.TransactionID, &tx.UserID, &tx.Points, &tx.Type,
			&tx.Reference, &tx.Description, &tx.CreatedAt); err != nil {
			return nil, err
		}
		transactions = append(transactions, &tx)
	}
	return transactions, nil
}

func (r *GamificationRepository) SaveBadge(ctx context.Context, userBadge *domain.UserBadge) error {
	query := `
		INSERT INTO user_badges (user_id, badge_id, earned_at)
		VALUES ($1, $2, $3)
	`
	_, err := r.db.ExecContext(ctx, query,
		userBadge.UserID, userBadge.BadgeID, userBadge.EarnedAt,
	)
	return err
}

func (r *GamificationRepository) GetUserBadges(ctx context.Context, userID string) ([]*domain.UserBadge, error) {
	query := `SELECT user_id, badge_id, earned_at FROM user_badges WHERE user_id = $1`

	rows, err := r.db.QueryContext(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var badges []*domain.UserBadge
	for rows.Next() {
		var badge domain.UserBadge
		if err := rows.Scan(&badge.UserID, &badge.BadgeID, &badge.EarnedAt); err != nil {
			return nil, err
		}
		badges = append(badges, &badge)
	}
	return badges, nil
}

func (r *GamificationRepository) SaveReward(ctx context.Context, reward *domain.Reward) error {
	query := `
		INSERT INTO rewards (reward_id, name, description, points_cost, reward_type, value, stock, is_active)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
	`
	_, err := r.db.ExecContext(ctx, query,
		reward.RewardID, reward.Name, reward.Description, reward.PointsCost,
		reward.RewardType, reward.Value, reward.Stock, reward.IsActive,
	)
	return err
}

func (r *GamificationRepository) GetAvailableRewards(ctx context.Context) ([]*domain.Reward, error) {
	query := `SELECT reward_id, name, description, points_cost, reward_type, value, stock, is_active 
			  FROM rewards WHERE is_active = true AND stock > 0`

	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var rewards []*domain.Reward
	for rows.Next() {
		var reward domain.Reward

		if err := rows.Scan(&reward.RewardID, &reward.Name, &reward.Description, &reward.PointsCost,
			&reward.RewardType, &reward.Value, &reward.Stock, &reward.IsActive); err != nil {
			return nil, err
		}

		rewards = append(rewards, &reward)
	}
	return rewards, nil
}

func (r *GamificationRepository) SaveRedemption(ctx context.Context, redemption *domain.RewardRedemption) error {
	query := `
		INSERT INTO reward_redemptions (redemption_id, user_id, reward_id, points_spent, status, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
	`
	_, err := r.db.ExecContext(ctx, query,
		redemption.RedemptionID, redemption.UserID, redemption.RewardID, redemption.PointsSpent,
		redemption.Status, redemption.CreatedAt, redemption.UpdatedAt,
	)
	return err
}

func (r *GamificationRepository) GetLeaderboard(ctx context.Context, limit int) ([]*domain.UserPoints, error) {
	query := `SELECT user_id, total_points, available_points, lifetime_earned, lifetime_spent, level, updated_at 
			  FROM user_points ORDER BY lifetime_earned DESC LIMIT $1`

	rows, err := r.db.QueryContext(ctx, query, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []*domain.UserPoints
	for rows.Next() {
		var points domain.UserPoints
		if err := rows.Scan(&points.UserID, &points.TotalPoints, &points.AvailablePoints,
			&points.LifetimeEarned, &points.LifetimeSpent, &points.Level, &points.UpdatedAt); err != nil {
			return nil, err
		}
		users = append(users, &points)
	}
	return users, nil
}

func (r *GamificationRepository) Close() error {
	return r.db.Close()
}
