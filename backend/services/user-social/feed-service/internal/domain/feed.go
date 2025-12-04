package domain

import (
	"time"

	"github.com/titan-commerce/backend/pkg/errors"
	"github.com/google/uuid"
)

type Post struct {
	ID            string
	UserID        string
	Content       string
	MediaURL      string
	Tags          []string
	LikesCount    int
	CommentsCount int
	CreatedAt     time.Time
	UpdatedAt     time.Time
}

func NewPost(userID, content, mediaURL string, tags []string) (*Post, error) {
	if userID == "" {
		return nil, errors.New(errors.ErrInvalidInput, "user ID is required")
	}
	if content == "" && mediaURL == "" {
		return nil, errors.New(errors.ErrInvalidInput, "content or media URL is required")
	}

	return &Post{
		ID:        uuid.New().String(),
		UserID:    userID,
		Content:   content,
		MediaURL:  mediaURL,
		Tags:      tags,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}, nil
}
