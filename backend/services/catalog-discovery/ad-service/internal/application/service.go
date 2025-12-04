package application

import (
	"context"
	"time"

	"github.com/titan-commerce/backend/ad-service/internal/domain"
	"github.com/titan-commerce/backend/pkg/logger"
)

type AdRepository interface {
	SaveCampaign(ctx context.Context, campaign *domain.Campaign) error
	GetActiveCampaigns(ctx context.Context) ([]*domain.Campaign, error)
	TrackEvent(ctx context.Context, event *domain.AdEvent) error
	DeductBudget(ctx context.Context, campaignID string, amount float64) error
}

type AdService struct {
	repo   AdRepository
	logger *logger.Logger
}

func NewAdService(repo AdRepository, logger *logger.Logger) *AdService {
	return &AdService{
		repo:   repo,
		logger: logger,
	}
}

func (s *AdService) CreateCampaign(ctx context.Context, sellerID, productID string, budget, bid float64, start, end time.Time) (*domain.Campaign, error) {
	campaign, err := domain.NewCampaign(sellerID, productID, budget, bid, start, end)
	if err != nil {
		return nil, err
	}

	if err := s.repo.SaveCampaign(ctx, campaign); err != nil {
		s.logger.Error(err, "failed to save campaign")
		return nil, err
	}

	s.logger.Infof("Campaign created for product %s", productID)
	return campaign, nil
}

func (s *AdService) GetAds(ctx context.Context, context string, limit int) ([]*domain.Campaign, error) {
	// Simplified ad serving logic: just return active campaigns
	// In reality, this would involve complex bidding logic
	return s.repo.GetActiveCampaigns(ctx)
}

func (s *AdService) TrackEvent(ctx context.Context, adID, userID, eventType string) error {
	// Deduct budget for clicks
	if eventType == "click" {
		// Assuming adID maps to campaignID for simplicity
		// In reality, we'd look up the campaign and bid amount
		// Hardcoding deduction for MVP
		if err := s.repo.DeductBudget(ctx, adID, 0.50); err != nil {
			s.logger.Error(err, "failed to deduct budget")
		}
	}

	return s.repo.TrackEvent(ctx, &domain.AdEvent{
		AdID:      adID,
		UserID:    userID,
		EventType: eventType,
		Timestamp: time.Now(),
	})
}
