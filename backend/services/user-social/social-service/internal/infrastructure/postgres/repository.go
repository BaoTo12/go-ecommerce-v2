package postgres

import (
	"context"
	"database/sql"
	"time"

	"github.com/titan-commerce/backend/social-service/internal/domain"
	"github.com/titan-commerce/backend/pkg/errors"
	"github.com/titan-commerce/backend/pkg/logger"
	_ "github.com/lib/pq"
)

type SocialRepository struct {
	db     *sql.DB
	logger *logger.Logger
}

func NewSocialRepository(databaseURL string, logger *logger.Logger) (*SocialRepository, error) {
	db, err := sql.Open("postgres", databaseURL)
	if err != nil {
		return nil, errors.Wrap(errors.ErrInternal, "failed to connect to database", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	
	if err := db.PingContext(ctx); err != nil {
		return nil, errors.Wrap(errors.ErrInternal, "failed to ping database", err)
	}

	logger.Info("Social PostgreSQL repository initialized")
	return &SocialRepository{db: db, logger: logger}, nil
}

func (r *SocialRepository) SaveFollow(ctx context.Context, follow *domain.Follow) error {
	query := `
		INSERT INTO follows (follower_id, followee_id, created_at)
		VALUES ($1, $2, $3)
		ON CONFLICT (follower_id, followee_id) DO NOTHING
	`

	_, err := r.db.ExecContext(ctx, query, follow.FollowerID, follow.FolloweeID, follow.CreatedAt)
	if err != nil {
		return errors.Wrap(errors.ErrInternal, "failed to save follow", err)
	}

	return nil
}

func (r *SocialRepository) DeleteFollow(ctx context.Context, followerID, followeeID string) error {
	query := `DELETE FROM follows WHERE follower_id = $1 AND followee_id = $2`
	_, err := r.db.ExecContext(ctx, query, followerID, followeeID)
	if err != nil {
		return errors.Wrap(errors.ErrInternal, "failed to delete follow", err)
	}
	return nil
}

func (r *SocialRepository) GetFollowers(ctx context.Context, userID string, page, pageSize int) ([]*domain.Follow, int, error) {
	offset := (page - 1) * pageSize

	query := `
		SELECT follower_id, followee_id, created_at
		FROM follows
		WHERE followee_id = $1
		ORDER BY created_at DESC
		LIMIT $2 OFFSET $3
	`

	rows, err := r.db.QueryContext(ctx, query, userID, pageSize, offset)
	if err != nil {
		return nil, 0, errors.Wrap(errors.ErrInternal, "failed to get followers", err)
	}
	defer rows.Close()

	var follows []*domain.Follow
	for rows.Next() {
		var f domain.Follow
		if err := rows.Scan(&f.FollowerID, &f.FolloweeID, &f.CreatedAt); err != nil {
			return nil, 0, errors.Wrap(errors.ErrInternal, "failed to scan follow", err)
		}
		follows = append(follows, &f)
	}

	var total int
	countQuery := `SELECT COUNT(*) FROM follows WHERE followee_id = $1`
	err = r.db.QueryRowContext(ctx, countQuery, userID).Scan(&total)
	if err != nil {
		return nil, 0, errors.Wrap(errors.ErrInternal, "failed to count followers", err)
	}

	return follows, total, nil
}

func (r *SocialRepository) GetFollowing(ctx context.Context, userID string, page, pageSize int) ([]*domain.Follow, int, error) {
	offset := (page - 1) * pageSize

	query := `
		SELECT follower_id, followee_id, created_at
		FROM follows
		WHERE follower_id = $1
		ORDER BY created_at DESC
		LIMIT $2 OFFSET $3
	`

	rows, err := r.db.QueryContext(ctx, query, userID, pageSize, offset)
	if err != nil {
		return nil, 0, errors.Wrap(errors.ErrInternal, "failed to get following", err)
	}
	defer rows.Close()

	var follows []*domain.Follow
	for rows.Next() {
		var f domain.Follow
		if err := rows.Scan(&f.FollowerID, &f.FolloweeID, &f.CreatedAt); err != nil {
			return nil, 0, errors.Wrap(errors.ErrInternal, "failed to scan follow", err)
		}
		follows = append(follows, &f)
	}

	var total int
	countQuery := `SELECT COUNT(*) FROM follows WHERE follower_id = $1`
	err = r.db.QueryRowContext(ctx, countQuery, userID).Scan(&total)
	if err != nil {
		return nil, 0, errors.Wrap(errors.ErrInternal, "failed to count following", err)
	}

	return follows, total, nil
}

func (r *SocialRepository) GetStats(ctx context.Context, userID string) (*domain.SocialStats, error) {
	stats := &domain.SocialStats{UserID: userID}

	err := r.db.QueryRowContext(ctx, "SELECT COUNT(*) FROM follows WHERE followee_id = $1", userID).Scan(&stats.FollowersCount)
	if err != nil {
		return nil, errors.Wrap(errors.ErrInternal, "failed to count followers", err)
	}

	err = r.db.QueryRowContext(ctx, "SELECT COUNT(*) FROM follows WHERE follower_id = $1", userID).Scan(&stats.FollowingCount)
	if err != nil {
		return nil, errors.Wrap(errors.ErrInternal, "failed to count following", err)
	}

	return stats, nil
}
