package postgres

import (
	"context"
	"database/sql"
	"encoding/json"
	"time"

	"github.com/titan-commerce/backend/feed-service/internal/domain"
	"github.com/titan-commerce/backend/pkg/errors"
	"github.com/titan-commerce/backend/pkg/logger"
	_ "github.com/lib/pq"
)

type FeedRepository struct {
	db     *sql.DB
	logger *logger.Logger
}

func NewFeedRepository(databaseURL string, logger *logger.Logger) (*FeedRepository, error) {
	db, err := sql.Open("postgres", databaseURL)
	if err != nil {
		return nil, errors.Wrap(errors.ErrInternal, "failed to connect to database", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	
	if err := db.PingContext(ctx); err != nil {
		return nil, errors.Wrap(errors.ErrInternal, "failed to ping database", err)
	}

	logger.Info("Feed PostgreSQL repository initialized")
	return &FeedRepository{db: db, logger: logger}, nil
}

func (r *FeedRepository) SavePost(ctx context.Context, post *domain.Post) error {
	tagsJSON, err := json.Marshal(post.Tags)
	if err != nil {
		return errors.Wrap(errors.ErrInternal, "failed to marshal tags", err)
	}

	query := `
		INSERT INTO posts (
			id, user_id, content, media_url, tags, 
			likes_count, comments_count, created_at, updated_at
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
	`

	_, err = r.db.ExecContext(ctx, query,
		post.ID, post.UserID, post.Content, post.MediaURL, tagsJSON,
		post.LikesCount, post.CommentsCount, post.CreatedAt, post.UpdatedAt,
	)

	if err != nil {
		return errors.Wrap(errors.ErrInternal, "failed to save post", err)
	}

	return nil
}

func (r *FeedRepository) DeletePost(ctx context.Context, postID string) error {
	query := `DELETE FROM posts WHERE id = $1`
	_, err := r.db.ExecContext(ctx, query, postID)
	if err != nil {
		return errors.Wrap(errors.ErrInternal, "failed to delete post", err)
	}
	return nil
}

func (r *FeedRepository) GetGlobalFeed(ctx context.Context, page, pageSize int) ([]*domain.Post, error) {
	offset := (page - 1) * pageSize

	query := `
		SELECT id, user_id, content, media_url, tags, 
			   likes_count, comments_count, created_at, updated_at
		FROM posts
		ORDER BY created_at DESC
		LIMIT $1 OFFSET $2
	`

	rows, err := r.db.QueryContext(ctx, query, pageSize, offset)
	if err != nil {
		return nil, errors.Wrap(errors.ErrInternal, "failed to get feed", err)
	}
	defer rows.Close()

	var posts []*domain.Post
	for rows.Next() {
		var p domain.Post
		var tagsJSON []byte
		err := rows.Scan(
			&p.ID, &p.UserID, &p.Content, &p.MediaURL, &tagsJSON,
			&p.LikesCount, &p.CommentsCount, &p.CreatedAt, &p.UpdatedAt,
		)
		if err != nil {
			return nil, errors.Wrap(errors.ErrInternal, "failed to scan post", err)
		}
		if err := json.Unmarshal(tagsJSON, &p.Tags); err != nil {
			return nil, errors.Wrap(errors.ErrInternal, "failed to unmarshal tags", err)
		}
		posts = append(posts, &p)
	}

	return posts, nil
}

func (r *FeedRepository) GetUserFeed(ctx context.Context, userID string, page, pageSize int) ([]*domain.Post, error) {
	// Placeholder for personalized feed logic
	return r.GetGlobalFeed(ctx, page, pageSize)
}
