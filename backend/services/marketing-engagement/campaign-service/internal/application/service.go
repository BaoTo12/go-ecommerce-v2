package application

import (
	"context"
	"time"

	"github.com/titan-commerce/backend/campaign-service/internal/domain"
	"github.com/titan-commerce/backend/pkg/errors"
	"github.com/titan-commerce/backend/pkg/logger"
)

type CampaignRepository interface {
	Save(ctx context.Context, campaign *domain.Campaign) error
	FindByID(ctx context.Context, campaignID string) (*domain.Campaign, error)
	Update(ctx context.Context, campaign *domain.Campaign) error
	FindActive(ctx context.Context) ([]*domain.Campaign, error)
	FindByType(ctx context.Context, campaignType domain.CampaignType) ([]*domain.Campaign, error)
	FindByProduct(ctx context.Context, productID string) ([]*domain.Campaign, error)
	FindScheduled(ctx context.Context) ([]*domain.Campaign, error)
	FindExpired(ctx context.Context) ([]*domain.Campaign, error)
}

type CampaignService struct {
	repo   CampaignRepository
	logger *logger.Logger
}

func NewCampaignService(repo CampaignRepository, logger *logger.Logger) *CampaignService {
	svc := &CampaignService{
		repo:   repo,
		logger: logger,
	}
	
	// Start background job to check campaign schedules
	go svc.runScheduler()
	
	return svc
}

func (s *CampaignService) runScheduler() {
	ticker := time.NewTicker(1 * time.Minute)
	for range ticker.C {
		ctx := context.Background()
		
		// Activate scheduled campaigns
		scheduled, _ := s.repo.FindScheduled(ctx)
		for _, c := range scheduled {
			if time.Now().After(c.StartTime) {
				c.Activate()
				s.repo.Update(ctx, c)
				s.logger.Infof("Campaign activated: %s", c.ID)
			}
		}

		// End expired campaigns
		active, _ := s.repo.FindActive(ctx)
		for _, c := range active {
			if time.Now().After(c.EndTime) {
				c.End()
				s.repo.Update(ctx, c)
				s.logger.Infof("Campaign ended: %s", c.ID)
			}
		}
	}
}

// CreateCampaign creates a new marketing campaign
func (s *CampaignService) CreateCampaign(ctx context.Context, name, description string, campaignType domain.CampaignType, start, end time.Time, budget float64) (*domain.Campaign, error) {
	if end.Before(start) {
		return nil, errors.New(errors.ErrInvalidInput, "end time must be after start time")
	}

	campaign := domain.NewCampaign(name, description, campaignType, start, end, budget)
	
	if err := s.repo.Save(ctx, campaign); err != nil {
		return nil, err
	}

	s.logger.Infof("Campaign created: %s (%s)", name, campaignType)
	return campaign, nil
}

// SetCampaignRules sets the rules for a campaign
func (s *CampaignService) SetCampaignRules(ctx context.Context, campaignID string, rules domain.CampaignRules) error {
	campaign, err := s.repo.FindByID(ctx, campaignID)
	if err != nil {
		return err
	}

	campaign.SetRules(rules)
	return s.repo.Update(ctx, campaign)
}

// AddProductsToCampaign adds products to a campaign
func (s *CampaignService) AddProductsToCampaign(ctx context.Context, campaignID string, productIDs []string) error {
	campaign, err := s.repo.FindByID(ctx, campaignID)
	if err != nil {
		return err
	}

	campaign.AddProducts(productIDs)
	return s.repo.Update(ctx, campaign)
}

// ScheduleCampaign schedules a campaign for activation
func (s *CampaignService) ScheduleCampaign(ctx context.Context, campaignID string) error {
	campaign, err := s.repo.FindByID(ctx, campaignID)
	if err != nil {
		return err
	}

	if campaign.Status != domain.CampaignStatusDraft {
		return errors.New(errors.ErrInvalidInput, "campaign must be in draft status")
	}

	campaign.Schedule()
	s.logger.Infof("Campaign scheduled: %s", campaignID)
	return s.repo.Update(ctx, campaign)
}

// ActivateCampaign immediately activates a campaign
func (s *CampaignService) ActivateCampaign(ctx context.Context, campaignID string) error {
	campaign, err := s.repo.FindByID(ctx, campaignID)
	if err != nil {
		return err
	}

	campaign.Activate()
	s.logger.Infof("Campaign activated: %s", campaignID)
	return s.repo.Update(ctx, campaign)
}

// PauseCampaign pauses an active campaign
func (s *CampaignService) PauseCampaign(ctx context.Context, campaignID string) error {
	campaign, err := s.repo.FindByID(ctx, campaignID)
	if err != nil {
		return err
	}

	campaign.Pause()
	s.logger.Infof("Campaign paused: %s", campaignID)
	return s.repo.Update(ctx, campaign)
}

// EndCampaign ends a campaign
func (s *CampaignService) EndCampaign(ctx context.Context, campaignID string) error {
	campaign, err := s.repo.FindByID(ctx, campaignID)
	if err != nil {
		return err
	}

	campaign.End()
	s.logger.Infof("Campaign ended: %s, stats: %+v", campaignID, campaign.Stats)
	return s.repo.Update(ctx, campaign)
}

// GetCampaign retrieves a campaign by ID
func (s *CampaignService) GetCampaign(ctx context.Context, campaignID string) (*domain.Campaign, error) {
	return s.repo.FindByID(ctx, campaignID)
}

// GetActiveCampaigns returns all active campaigns
func (s *CampaignService) GetActiveCampaigns(ctx context.Context) ([]*domain.Campaign, error) {
	return s.repo.FindActive(ctx)
}

// GetCampaignsForProduct returns campaigns applicable to a product
func (s *CampaignService) GetCampaignsForProduct(ctx context.Context, productID string) ([]*domain.Campaign, error) {
	return s.repo.FindByProduct(ctx, productID)
}

// RecordConversion records a conversion for a campaign
func (s *CampaignService) RecordConversion(ctx context.Context, campaignID string, orderValue float64) error {
	campaign, err := s.repo.FindByID(ctx, campaignID)
	if err != nil {
		return err
	}

	campaign.RecordConversion(orderValue)
	return s.repo.Update(ctx, campaign)
}

// RecordImpression records an impression for a campaign
func (s *CampaignService) RecordImpression(ctx context.Context, campaignID string) error {
	campaign, err := s.repo.FindByID(ctx, campaignID)
	if err != nil {
		return err
	}

	campaign.Stats.Impressions++
	return s.repo.Update(ctx, campaign)
}

// RecordClick records a click for a campaign
func (s *CampaignService) RecordClick(ctx context.Context, campaignID string) error {
	campaign, err := s.repo.FindByID(ctx, campaignID)
	if err != nil {
		return err
	}

	campaign.Stats.Clicks++
	return s.repo.Update(ctx, campaign)
}

// GetCampaignStats returns stats for a campaign
func (s *CampaignService) GetCampaignStats(ctx context.Context, campaignID string) (*domain.CampaignStats, error) {
	campaign, err := s.repo.FindByID(ctx, campaignID)
	if err != nil {
		return nil, err
	}
	return &campaign.Stats, nil
}
