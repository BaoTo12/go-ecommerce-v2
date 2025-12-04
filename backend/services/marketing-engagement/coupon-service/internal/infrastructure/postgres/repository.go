package postgres

import (
	"context"
	"database/sql"
	"encoding/json"
	"time"

	_ "github.com/lib/pq"
	"github.com/titan-commerce/backend/coupon-service/internal/domain"
	"github.com/titan-commerce/backend/pkg/errors"
	"github.com/titan-commerce/backend/pkg/logger"
)

type CouponRepository struct {
	db     *sql.DB
	logger *logger.Logger
}

func NewCouponRepository(databaseURL string, logger *logger.Logger) (*CouponRepository, error) {
	db, err := sql.Open("postgres", databaseURL)
	if err != nil {
		return nil, errors.Wrap(errors.ErrInternal, "failed to connect to database", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := db.PingContext(ctx); err != nil {
		return nil, errors.Wrap(errors.ErrInternal, "failed to ping database", err)
	}

	logger.Info("Coupon PostgreSQL repository initialized")
	return &CouponRepository{db: db, logger: logger}, nil
}

func (r *CouponRepository) SaveCoupon(ctx context.Context, coupon *domain.Coupon) error {
	productsJSON, _ := json.Marshal(coupon.ApplicableProducts)
	categoriesJSON, _ := json.Marshal(coupon.ApplicableCategories)

	query := `
		INSERT INTO coupons (coupon_id, code, type, discount_value, min_order_value, max_discount, 
			usage_limit, usage_count, per_user_limit, valid_from, valid_until, is_active, 
			applicable_products, applicable_categories, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16)
	`
	_, err := r.db.ExecContext(ctx, query,
		coupon.CouponID, coupon.Code, coupon.Type, coupon.DiscountValue, coupon.MinOrderValue,
		coupon.MaxDiscount, coupon.UsageLimit, coupon.UsageCount, coupon.PerUserLimit,
		coupon.ValidFrom, coupon.ValidUntil, coupon.IsActive, productsJSON, categoriesJSON,
		coupon.CreatedAt, coupon.UpdatedAt,
	)
	return err
}

func (r *CouponRepository) GetCoupon(ctx context.Context, couponID string) (*domain.Coupon, error) {
	query := `SELECT coupon_id, code, type, discount_value, min_order_value, max_discount, 
		usage_limit, usage_count, per_user_limit, valid_from, valid_until, is_active, 
		applicable_products, applicable_categories, created_at, updated_at 
		FROM coupons WHERE coupon_id = $1`

	var coupon domain.Coupon
	var productsJSON, categoriesJSON []byte

	err := r.db.QueryRowContext(ctx, query, couponID).Scan(
		&coupon.CouponID, &coupon.Code, &coupon.Type, &coupon.DiscountValue, &coupon.MinOrderValue,
		&coupon.MaxDiscount, &coupon.UsageLimit, &coupon.UsageCount, &coupon.PerUserLimit,
		&coupon.ValidFrom, &coupon.ValidUntil, &coupon.IsActive, &productsJSON, &categoriesJSON,
		&coupon.CreatedAt, &coupon.UpdatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, errors.New(errors.ErrNotFound, "coupon not found")
	}
	if err != nil {
		return nil, err
	}

	json.Unmarshal(productsJSON, &coupon.ApplicableProducts)
	json.Unmarshal(categoriesJSON, &coupon.ApplicableCategories)
	return &coupon, nil
}

func (r *CouponRepository) GetCouponByCode(ctx context.Context, code string) (*domain.Coupon, error) {
	query := `SELECT coupon_id, code, type, discount_value, min_order_value, max_discount, 
		usage_limit, usage_count, per_user_limit, valid_from, valid_until, is_active, 
		applicable_products, applicable_categories, created_at, updated_at 
		FROM coupons WHERE code = $1`

	var coupon domain.Coupon
	var productsJSON, categoriesJSON []byte

	err := r.db.QueryRowContext(ctx, query, code).Scan(
		&coupon.CouponID, &coupon.Code, &coupon.Type, &coupon.DiscountValue, &coupon.MinOrderValue,
		&coupon.MaxDiscount, &coupon.UsageLimit, &coupon.UsageCount, &coupon.PerUserLimit,
		&coupon.ValidFrom, &coupon.ValidUntil, &coupon.IsActive, &productsJSON, &categoriesJSON,
		&coupon.CreatedAt, &coupon.UpdatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, errors.New(errors.ErrNotFound, "coupon not found")
	}
	if err != nil {
		return nil, err
	}

	json.Unmarshal(productsJSON, &coupon.ApplicableProducts)
	json.Unmarshal(categoriesJSON, &coupon.ApplicableCategories)
	return &coupon, nil
}

func (r *CouponRepository) UpdateCoupon(ctx context.Context, coupon *domain.Coupon) error {
	productsJSON, _ := json.Marshal(coupon.ApplicableProducts)
	categoriesJSON, _ := json.Marshal(coupon.ApplicableCategories)

	query := `
		UPDATE coupons 
		SET code = $2, type = $3, discount_value = $4, min_order_value = $5, max_discount = $6,
			usage_limit = $7, usage_count = $8, per_user_limit = $9, valid_from = $10, valid_until = $11,
			is_active = $12, applicable_products = $13, applicable_categories = $14, updated_at = $15
		WHERE coupon_id = $1
	`
	_, err := r.db.ExecContext(ctx, query,
		coupon.CouponID, coupon.Code, coupon.Type, coupon.DiscountValue, coupon.MinOrderValue,
		coupon.MaxDiscount, coupon.UsageLimit, coupon.UsageCount, coupon.PerUserLimit,
		coupon.ValidFrom, coupon.ValidUntil, coupon.IsActive, productsJSON, categoriesJSON,
		coupon.UpdatedAt,
	)
	return err
}

func (r *CouponRepository) SaveUsage(ctx context.Context, usage *domain.CouponUsage) error {
	query := `
		INSERT INTO coupon_usage (usage_id, coupon_id, user_id, order_id, discount_applied, used_at)
		VALUES ($1, $2, $3, $4, $5, $6)
	`
	_, err := r.db.ExecContext(ctx, query,
		usage.UsageID, usage.CouponID, usage.UserID, usage.OrderID, usage.DiscountApplied, usage.UsedAt,
	)
	return err
}

func (r *CouponRepository) GetUserUsageCount(ctx context.Context, couponID, userID string) (int, error) {
	query := `SELECT COUNT(*) FROM coupon_usage WHERE coupon_id = $1 AND user_id = $2`

	var count int
	err := r.db.QueryRowContext(ctx, query, couponID, userID).Scan(&count)
	return count, err
}

func (r *CouponRepository) Close() error {
	return r.db.Close()
}
