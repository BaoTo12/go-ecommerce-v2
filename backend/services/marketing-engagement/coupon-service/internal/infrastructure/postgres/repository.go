package postgres

import (
	"context"

	"github.com/titan-commerce/backend/coupon-service/internal/domain"
)

type CouponRepository struct{}

func NewCouponRepository() *CouponRepository {
	return &CouponRepository{}
}

func (r *CouponRepository) Save(ctx context.Context, coupon *domain.Coupon) error {
	return nil
}

func (r *CouponRepository) FindByID(ctx context.Context, couponID string) (*domain.Coupon, error) {
	return nil, nil
}

func (r *CouponRepository) FindByCode(ctx context.Context, code string) (*domain.Coupon, error) {
	return nil, nil
}

func (r *CouponRepository) Update(ctx context.Context, coupon *domain.Coupon) error {
	return nil
}

func (r *CouponRepository) FindActive(ctx context.Context) ([]*domain.Coupon, error) {
	return nil, nil
}

func (r *CouponRepository) SaveUsage(ctx context.Context, usage *domain.CouponUsage) error {
	return nil
}

func (r *CouponRepository) GetUserUsage(ctx context.Context, couponID, userID string) (int, error) {
	return 0, nil
}
