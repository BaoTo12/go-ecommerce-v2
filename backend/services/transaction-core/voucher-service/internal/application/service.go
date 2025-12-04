package application

import (
	"context"

	"github.com/titan-commerce/backend/voucher-service/internal/domain"
	"github.com/titan-commerce/backend/pkg/logger"
)

type VoucherRepository interface {
	Save(ctx context.Context, voucher *domain.Voucher) error
	FindByCode(ctx context.Context, code string) (*domain.Voucher, error)
	FindByUserID(ctx context.Context, userID string) ([]*domain.Voucher, error)
	Update(ctx context.Context, voucher *domain.Voucher) error
}

type VoucherService struct {
	repo   VoucherRepository
	logger *logger.Logger
}

func NewVoucherService(repo VoucherRepository, logger *logger.Logger) *VoucherService {
	return &VoucherService{
		repo:   repo,
		logger: logger,
	}
}

// CreateVoucher generates a new voucher (Command)
func (s *VoucherService) CreateVoucher(ctx context.Context, code string, voucherType domain.VoucherType, value float64, userID string, expiresAt time.Time) (*domain.Voucher, error) {
	voucher, err := domain.NewVoucher(code, voucherType, value, userID, expiresAt)
	if err != nil {
		return nil, err
	}

	if err := s.repo.Save(ctx, voucher); err != nil {
		s.logger.Error(err, "failed to save voucher")
		return nil, err
	}

	s.logger.Infof("Voucher created: code=%s, type=%s, user=%s", code, voucherType, userID)
	return voucher, nil
}

// ValidateVoucher checks if voucher can be used (Query)
func (s *VoucherService) ValidateVoucher(ctx context.Context, code, userID string) (*domain.Voucher, string, error) {
	voucher, err := s.repo.FindByCode(ctx, code)
	if err != nil {
		return nil, "voucher not found", err
	}

	canUse, reason := voucher.CanUse(userID)
	if !canUse {
		return voucher, reason, errors.New(errors.ErrInvalidInput, reason)
	}

	return voucher, "", nil
}

// RedeemVoucher marks voucher as used (Command)
func (s *VoucherService) RedeemVoucher(ctx context.Context, code, userID string) (*domain.Voucher, error) {
	voucher, _, err := s.ValidateVoucher(ctx, code, userID)
	if err != nil {
		return nil, err
	}

	if err := voucher.Redeem(); err != nil {
		return nil, err
	}

	if err := s.repo.Update(ctx, voucher); err != nil {
		return nil, err
	}

	s.logger.Infof("Voucher redeemed: code=%s, user=%s, value=%.2f", code, userID, voucher.Value)
	return voucher, nil
}

// GetUserVouchers retrieves all vouchers for a user (Query)
func (s *VoucherService) GetUserVouchers(ctx context.Context, userID string) ([]*domain.Voucher, error) {
	return s.repo.FindByUserID(ctx, userID)
}
