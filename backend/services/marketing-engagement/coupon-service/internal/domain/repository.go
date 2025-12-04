package domain

import "context"

type CouponRepository interface {
	SaveCoupon(ctx context.Context, coupon *Coupon) error
	GetCoupon(ctx context.Context, couponID string) (*Coupon, error)
	GetCouponByCode(ctx context.Context, code string) (*Coupon, error)
	UpdateCoupon(ctx context.Context, coupon *Coupon) error
	SaveUsage(ctx context.Context, usage *CouponUsage) error
	GetUserUsageCount(ctx context.Context, couponID, userID string) (int, error)
}

