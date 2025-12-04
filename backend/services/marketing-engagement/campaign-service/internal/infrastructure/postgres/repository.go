package postgres

import (
	"context"

	"github.com/titan-commerce/backend/campaign-service/internal/domain"
)

type CampaignRepository struct{}

func NewCampaignRepository() *CampaignRepository {
	return &CampaignRepository{}
}

func (r *CampaignRepository) Save(ctx context.Context, campaign *domain.Campaign) error {
	return nil
}

func (r *CampaignRepository) FindByID(ctx context.Context, campaignID string) (*domain.Campaign, error) {
	return &domain.Campaign{ID: campaignID}, nil
}

func (r *CampaignRepository) Update(ctx context.Context, campaign *domain.Campaign) error {
	return nil
}

func (r *CampaignRepository) FindActive(ctx context.Context) ([]*domain.Campaign, error) {
	return nil, nil
}

func (r *CampaignRepository) FindByType(ctx context.Context, campaignType domain.CampaignType) ([]*domain.Campaign, error) {
	return nil, nil
}

func (r *CampaignRepository) FindByProduct(ctx context.Context, productID string) ([]*domain.Campaign, error) {
	return nil, nil
}

func (r *CampaignRepository) FindScheduled(ctx context.Context) ([]*domain.Campaign, error) {
	return nil, nil
}

func (r *CampaignRepository) FindExpired(ctx context.Context) ([]*domain.Campaign, error) {
	return nil, nil
}
