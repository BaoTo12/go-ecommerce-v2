package application

import (
	"context"

	"github.com/titan-commerce/backend/campaign-service/internal/domain"
	"github.com/titan-commerce/backend/pkg/logger"
)

type CampaignRepository interface {
	Save(ctx context.Context, campaign *domain.Campaign) error
	FindByID(ctx context.Context, campaignID string) (*domain.Campaign, error)
	FindActive(ctx context.Context, pageSize int) ([]*domain.Campaign, error)
	Update(ctx context.Context, campaign *domain.Campaign) error
}

type CampaignService struct {
	repo   CampaignRepository
	logger *logger.Logger
}

func NewCampaignService(repo CampaignRepository, logger *logger.Logger) *CampaignService {
	return &CampaignService{
		repo:   repo,
		logger: logger,
	}
}

// CreateCampaign creates a marketing campaign (Command)
func (s *CampaignService) CreateCampaign(ctx context.Context, name, description string, campaignType domain.CampaignType, startTime, endTime time.Time, targetUsers int) (*domain.Campaign, error) {
	campaign, err := domain.NewCampaign(name, description, campaignType, startTime, endTime, targetUsers)
	if err != nil {
		return nil, err
	}

	if err := s.repo.Save(ctx, campaign); err != nil {
		s.logger.Error(err, "failed to save campaign")
		return nil, err
	}

	s.logger.Infof("Campaign created: %s (%s), target=%d users", name, campaignType, targetUsers)
	return campaign, nil
}

// GetCampaign retrieves campaign details (Query)
func (s *CampaignService) GetCampaign(ctx context.Context, campaignID string) (*domain.Campaign, error) {
	return s.repo.FindByID(ctx, campaignID)
}

// ListActiveCampaigns returns currently active campaigns (Query)
func (s *CampaignService) ListActiveCampaigns(ctx context.Context, pageSize int) ([]*domain.Campaign, error) {
	return s.repo.FindActive(ctx, pageSize)
}

// TrackConversion records a conversion for a campaign (Command)
func (s *CampaignService) TrackConversion(ctx context.Context, campaignID, userID, orderID string) error {
	campaign, err := s.repo.FindByID(ctx, campaignID)
	if err != nil {
		return err
	}

	campaign.TrackConversion()

	if err := s.repo.Update(ctx, campaign); err != nil {
		return err
	}

	s.logger.Infof("Conversion tracked: campaign=%s, user=%s, order=%s, rate=%.2f%%",
		campaignID, userID, orderID, campaign.GetConversionRate())
	return nil
}
