package domain

import "context"

type FraudRepository interface {
	SaveFraudCheck(ctx context.Context, check *FraudCheck) error
	GetFraudCheck(ctx context.Context, orderID string) (*FraudCheck, error)
	SaveUserProfile(ctx context.Context, profile *UserRiskProfile) error
	GetUserProfile(ctx context.Context, userID string) (*UserRiskProfile, error)
	GetRecentChecks(ctx context.Context, userID string, limit int) ([]*FraudCheck, error)
}

