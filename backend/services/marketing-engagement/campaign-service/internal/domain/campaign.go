package domain

import (
	"time"

	"github.com/google/uuid"
	"github.com/titan-commerce/backend/pkg/errors"
)

type CampaignType string

const (
	CampaignTypeFlashSale  CampaignType = "FLASH_SALE"
	CampaignTypeSeasonal   CampaignType = "SEASONAL"
	CampaignTypeCategory   CampaignType = "CATEGORY_SALE"
	CampaignTypeNewUser    CampaignType = "NEW_USER"
)

type Campaign struct {
	ID             string
	Name           string
	Description    string
	Type           CampaignType
	StartTime      time.Time
	EndTime        time.Time
	TargetUsers    int
	ReachedUsers   int
	Conversions    int
	Active         bool
	CreatedAt      time.Time
}

func NewCampaign(name, description string, campaignType CampaignType, startTime, endTime time.Time, targetUsers int) (*Campaign, error) {
	if name == "" {
		return nil, errors.New(errors.ErrInvalidInput, "campaign name is required")
	}
	if endTime.Before(startTime) {
		return nil, errors.New(errors.ErrInvalidInput, "end time must be after start time")
	}

	return &Campaign{
		ID:           uuid.New().String(),
		Name:         name,
		Description:  description,
		Type:         campaignType,
		StartTime:    startTime,
		EndTime:      endTime,
		TargetUsers:  targetUsers,
		ReachedUsers: 0,
		Conversions:  0,
		Active:       time.Now().After(startTime) && time.Now().Before(endTime),
		CreatedAt:    time.Now(),
	}, nil
}

func (c *Campaign) TrackConversion() {
	c.Conversions++
}

func (c *Campaign) IncrementReach() {
	c.ReachedUsers++
}

func (c *Campaign) GetConversionRate() float64 {
	if c.ReachedUsers == 0 {
		return 0
	}
	return float64(c.Conversions) / float64(c.ReachedUsers) * 100
}

func (c *Campaign) IsActive() bool {
	now := time.Now()
	return now.After(c.StartTime) && now.Before(c.EndTime)
}
