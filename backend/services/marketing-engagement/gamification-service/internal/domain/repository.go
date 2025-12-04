package domain

import "context"

type GamificationRepository interface {
	GetUserPoints(ctx context.Context, userID string) (*UserPoints, error)
	SaveUserPoints(ctx context.Context, points *UserPoints) error
	SaveTransaction(ctx context.Context, trans *PointsTransaction) error
	GetUserBadges(ctx context.Context, userID string) ([]*UserBadge, error)
	AwardBadge(ctx context.Context, userBadge *UserBadge) error
	GetRewards(ctx context.Context) ([]*Reward, error)
	GetReward(ctx context.Context, rewardID string) (*Reward, error)
}

