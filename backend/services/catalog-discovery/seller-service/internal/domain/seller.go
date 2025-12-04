package domain

import (
	"time"

	"github.com/google/uuid"
	"github.com/titan-commerce/backend/pkg/errors"
)

type SellerStatus string

const (
	SellerStatusPendingVerification SellerStatus = "PENDING_VERIFICATION"
	SellerStatusVerified            SellerStatus = "VERIFIED"
	SellerStatusActive              SellerStatus = "ACTIVE"
	SellerStatusSuspended           SellerStatus = "SUSPENDED"
)

type Seller struct {
	ID            string
	UserID        string
	BusinessName  string
	BusinessType  string // Individual, Company
	TaxID         string
	Address       string
	Phone         string
	BankAccountID string
	Status        SellerStatus
	Rating        float64
	TotalProducts int
	TotalSales    int
	JoinedAt      time.Time
	UpdatedAt     time.Time
}

func NewSeller(userID, businessName, businessType, taxID, address, phone, bankAccountID string) (*Seller, error) {
	if userID == "" || businessName == "" {
		return nil, errors.New(errors.ErrInvalidInput, "user ID and business name are required")
	}
	if taxID == "" {
		return nil, errors.New(errors.ErrInvalidInput, "tax ID is required for KYC")
	}

	now := time.Now()
	return &Seller{
		ID:            uuid.New().String(),
		UserID:        userID,
		BusinessName:  businessName,
		BusinessType:  businessType,
		TaxID:         taxID,
		Address:       address,
		Phone:         phone,
		BankAccountID: bankAccountID,
		Status:        SellerStatusPendingVerification,
		Rating:        0.0,
		TotalProducts: 0,
		TotalSales:    0,
		JoinedAt:      now,
		UpdatedAt:     now,
	}, nil
}

func (s *Seller) Verify() {
	s.Status = SellerStatusVerified
	s.UpdatedAt = time.Now()
}

func (s *Seller) Activate() {
	s.Status = SellerStatusActive
	s.UpdatedAt = time.Now()
}

func (s *Seller) Suspend(reason string) {
	s.Status = SellerStatusSuspended
	s.UpdatedAt = time.Now()
}

func (s *Seller) UpdateRating(newRating float64) {
	s.Rating = newRating
	s.UpdatedAt = time.Now()
}

func (s *Seller) IncrementSales() {
	s.TotalSales++
	s.UpdatedAt = time.Now()
}
