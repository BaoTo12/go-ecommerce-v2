package domain

import "context"

type CampaignRepository interface {
	SaveCampaign(ctx context.Context, campaign *Campaign) error
	GetCampaign(ctx context.Context, campaignID string) (*Campaign, error)
	UpdateCampaign(ctx context.Context, campaign *Campaign) error
	GetActiveCampaigns(ctx context.Context) ([]*Campaign, error)
}

