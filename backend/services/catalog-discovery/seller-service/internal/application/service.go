package application

import (
	"context"

	"github.com/titan-commerce/backend/seller-service/internal/domain"
	"github.com/titan-commerce/backend/pkg/logger"
)

type SellerRepository interface {
	Save(ctx context.Context, seller *domain.Seller) error
	FindByID(ctx context.Context, sellerID string) (*domain.Seller, error)
	FindByUserID(ctx context.Context, userID string) (*domain.Seller, error)
	Update(ctx context.Context, seller *domain.Seller) error
}

type SellerService struct {
	repo   SellerRepository
	logger *logger.Logger
}

func NewSellerService(repo SellerRepository, logger *logger.Logger) *SellerService {
	return &SellerService{
		repo:   repo,
		logger: logger,
	}
}

// RegisterSeller creates a new seller account (Command)
func (s *SellerService) RegisterSeller(ctx context.Context, userID, businessName, businessType, taxID, address, phone, bankAccountID string, kycDocuments []string) (*domain.Seller, error) {
	seller, err := domain.NewSeller(userID, businessName, businessType, taxID, address, phone, bankAccountID)
	if err != nil {
		return nil, err
	}

	if err := s.repo.Save(ctx, seller); err != nil {
		s.logger.Error(err, "failed to save seller")
		return nil, err
	}

	s.logger.Infof("Seller registered: %s (%s), status=%s", businessName, userID, seller.Status)
	return seller, nil
}

// GetSeller retrieves seller information (Query)
func (s *SellerService) GetSeller(ctx context.Context, sellerID string) (*domain.Seller, error) {
	return s.repo.FindByID(ctx, sellerID)
}

// UpdateSellerStatus updates seller verification status (Command)
func (s *SellerService) UpdateSellerStatus(ctx context.Context, sellerID string, status domain.SellerStatus, reason string) (*domain.Seller, error) {
	seller, err := s.repo.FindByID(ctx, sellerID)
	if err != nil {
		return nil, err
	}

	switch status {
	case domain.SellerStatusVerified:
		seller.Verify()
	case domain.SellerStatusActive:
		seller.Activate()
	case domain.SellerStatusSuspended:
		seller.Suspend(reason)
	}

	if err := s.repo.Update(ctx, seller); err != nil {
		return nil, err
	}

	s.logger.Infof("Seller status updated: seller=%s, status=%s", sellerID, status)
	return seller, nil
}
