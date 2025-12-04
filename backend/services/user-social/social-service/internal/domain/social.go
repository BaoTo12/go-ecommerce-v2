package domain

import (
	"time"

	"github.com/titan-commerce/backend/pkg/errors"
)

type Follow struct {
	FollowerID string
	FolloweeID string
	CreatedAt  time.Time
}

type SocialStats struct {
	UserID         string
	FollowersCount int
	FollowingCount int
}

func NewFollow(followerID, followeeID string) (*Follow, error) {
	if followerID == "" || followeeID == "" {
		return nil, errors.New(errors.ErrInvalidInput, "follower and followee IDs are required")
	}
	if followerID == followeeID {
		return nil, errors.New(errors.ErrInvalidInput, "cannot follow yourself")
	}

	return &Follow{
		FollowerID: followerID,
		FolloweeID: followeeID,
		CreatedAt:  time.Now(),
	}, nil
}
