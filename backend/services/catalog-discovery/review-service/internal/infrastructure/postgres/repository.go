package postgres

import (
	"context"
	"database/sql"
	"encoding/json"
	"time"

	"github.com/titan-commerce/backend/review-service/internal/domain"
	"github.com/titan-commerce/backend/pkg/errors"
	"github.com/titan-commerce/backend/pkg/logger"
	_ "github.com/lib/pq"
)

type ReviewRepository struct {
	db     *sql.DB
	logger *logger.Logger
}

func NewReviewRepository(databaseURL string, logger *logger.Logger) (*ReviewRepository, error) {
	db, err := sql.Open("postgres", databaseURL)
	if err != nil {
		return nil, errors.Wrap(errors.ErrInternal, "failed to connect to database", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	
	if err := db.PingContext(ctx); err != nil {
		return nil, errors.Wrap(errors.ErrInternal, "failed to ping database", err)
	}

	logger.Info("Review PostgreSQL repository initialized")
	return &ReviewRepository{db: db, logger: logger}, nil
}

func (r *ReviewRepository) Save(ctx context.Context, review *domain.Review) error {
	imagesJSON, err := json.Marshal(review.Images)
	if err != nil {
		return errors.Wrap(errors.ErrInternal, "failed to marshal images", err)
	}

	query := `
		INSERT INTO reviews (
			id, user_id, product_id, rating, comment, images, created_at
		) VALUES ($1, $2, $3, $4, $5, $6, $7)
	`

	_, err = r.db.ExecContext(ctx, query,
		review.ID, review.UserID, review.ProductID, review.Rating,
		review.Comment, imagesJSON, review.CreatedAt,
	)

	if err != nil {
		return errors.Wrap(errors.ErrInternal, "failed to save review", err)
	}

	return nil
}

func (r *ReviewRepository) GetByProduct(ctx context.Context, productID string, page, pageSize int) ([]*domain.Review, int, error) {
	offset := (page - 1) * pageSize

	query := `
		SELECT id, user_id, product_id, rating, comment, images, created_at
		FROM reviews
		WHERE product_id = $1
		ORDER BY created_at DESC
		LIMIT $2 OFFSET $3
	`

	rows, err := r.db.QueryContext(ctx, query, productID, pageSize, offset)
	if err != nil {
		return nil, 0, errors.Wrap(errors.ErrInternal, "failed to get reviews", err)
	}
	defer rows.Close()

	var reviews []*domain.Review
	for rows.Next() {
		var rev domain.Review
		var imagesJSON []byte
		err := rows.Scan(
			&rev.ID, &rev.UserID, &rev.ProductID, &rev.Rating,
			&rev.Comment, &imagesJSON, &rev.CreatedAt,
		)
		if err != nil {
			return nil, 0, errors.Wrap(errors.ErrInternal, "failed to scan review", err)
		}
		json.Unmarshal(imagesJSON, &rev.Images)
		reviews = append(reviews, &rev)
	}

	var total int
	r.db.QueryRowContext(ctx, "SELECT COUNT(*) FROM reviews WHERE product_id = $1", productID).Scan(&total)

	return reviews, total, nil
}

func (r *ReviewRepository) GetStats(ctx context.Context, productID string) (*domain.ReviewStats, error) {
	stats := &domain.ReviewStats{
		ProductID:          productID,
		RatingDistribution: make(map[int]int),
	}

	// Get total and average
	query := `
		SELECT COUNT(*), COALESCE(AVG(rating), 0)
		FROM reviews
		WHERE product_id = $1
	`
	err := r.db.QueryRowContext(ctx, query, productID).Scan(&stats.TotalReviews, &stats.AverageRating)
	if err != nil {
		return nil, errors.Wrap(errors.ErrInternal, "failed to get stats", err)
	}

	// Get distribution
	distQuery := `
		SELECT rating, COUNT(*)
		FROM reviews
		WHERE product_id = $1
		GROUP BY rating
	`
	rows, err := r.db.QueryContext(ctx, distQuery, productID)
	if err != nil {
		return nil, errors.Wrap(errors.ErrInternal, "failed to get distribution", err)
	}
	defer rows.Close()

	for rows.Next() {
		var rating, count int
		if err := rows.Scan(&rating, &count); err != nil {
			continue
		}
		stats.RatingDistribution[rating] = count
	}

	return stats, nil
}
