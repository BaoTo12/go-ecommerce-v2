package domain

import (
	"time"

	"github.com/google/uuid"
	"github.com/titan-commerce/backend/pkg/errors"
)

type VoucherType string

const (
	VoucherTypeDiscount     VoucherType = "DISCOUNT"
	VoucherTypeFreeShipping VoucherType = "FREE_SHIPPING"
	VoucherTypeGiftCard     VoucherType = "GIFT_CARD"
)

type Voucher struct {
	ID          string
	Code        string
	Type        VoucherType
	Value       float64
	UserID      string // Assigned to specific user
	Used        bool
	UsedAt      *time.Time
	ExpiresAt   time.Time
	CreatedAt   time.Time
}

func NewVoucher(code string, voucherType VoucherType, value float64, userID string, expiresAt time.Time) (*Voucher, error) {
	if code == "" {
		return nil, errors.New(errors.ErrInvalidInput, "voucher code is required")
	}
	if value <= 0 {
		return nil, errors.New(errors.ErrInvalidInput, "voucher value must be positive")
	}

	return &Voucher{
		ID:        uuid.New().String(),
		Code:      code,
		Type:      voucherType,
		Value:     value,
		UserID:    userID,
		Used:      false,
		ExpiresAt: expiresAt,
		CreatedAt: time.Now(),
	}, nil
}

func (v *Voucher) CanUse(userID string) (bool, string) {
	if v.UserID != "" && v.UserID != userID {
		return false, "voucher is assigned to another user"
	}
	if v.Used {
		return false, "voucher already used"
	}
	if time.Now().After(v.ExpiresAt) {
		return false, "voucher has expired"
	}
	return true, ""
}

func (v *Voucher) Redeem() error {
	if v.Used {
		return errors.New(errors.ErrInvalidInput, "voucher already used")
	}
	if time.Now().After(v.ExpiresAt) {
		return errors.New(errors.ErrInvalidInput, "voucher has expired")
	}

	v.Used = true
	now := time.Now()
	v.UsedAt = &now
	return nil
}
