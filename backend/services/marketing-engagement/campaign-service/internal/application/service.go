package application

import (
	"context"

	"github.com/titan-commerce/backend/campaign-service/internal/domain"
	"github.com/titan-commerce/backend/pkg/logger"
)

type CampaignService struct {
	repo   domain.CampaignRepository
	logger *logger.Logger
}

func NewCampaignService(repo domain.CampaignRepository, logger *logger.Logger) *CampaignService {
	return &CampaignService{repo: repo, logger: logger}
}

// CreateCampaign creates a new marketing campaign
func (s *CampaignService) CreateCampaign(
	ctx context.Context,
	name, description string,
	campaignType domain.CampaignType,
) (*domain.Campaign, error) {
	campaign := domain.NewCampaign(name, description, campaignType)

	if err := s.repo.SaveCampaign(ctx, campaign); err != nil {
		s.logger.Error(err, "failed to create campaign")
		return nil, err
	}

	s.logger.Infof("Campaign created: %s", name)
	return campaign, nil
}

// ActivateCampaign starts a campaign
func (s *CampaignService) ActivateCampaign(ctx context.Context, campaignID string) error {
	campaign, err := s.repo.GetCampaign(ctx, campaignID)
	if err != nil {
		return err
	}

	campaign.Activate()

	if err := s.repo.UpdateCampaign(ctx, campaign); err != nil {
		s.logger.Error(err, "failed to activate campaign")
		return err
	}

	s.logger.Infof("Campaign activated: %s", campaignID)
	return nil
}

// RecordConversion records a campaign conversion
func (s *CampaignService) RecordConversion(ctx context.Context, campaignID string, revenue int64) error {
	campaign, err := s.repo.GetCampaign(ctx, campaignID)
	if err != nil {
		return err
	}

	campaign.RecordConversion(revenue)

	if err := s.repo.UpdateCampaign(ctx, campaign); err != nil {
		return err
	}

	return nil
}

