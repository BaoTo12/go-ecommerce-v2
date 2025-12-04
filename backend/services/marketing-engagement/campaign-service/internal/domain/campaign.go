package domain

import (
	"time"

	"github.com/google/uuid"
)

type CampaignStatus string
type CampaignType string

const (
	CampaignStatusDraft     CampaignStatus = "DRAFT"
	CampaignStatusScheduled CampaignStatus = "SCHEDULED"
	CampaignStatusActive    CampaignStatus = "ACTIVE"
	CampaignStatusPaused    CampaignStatus = "PAUSED"
	CampaignStatusEnded     CampaignStatus = "ENDED"

	CampaignTypeDiscount    CampaignType = "DISCOUNT"
	CampaignTypeBundleDeal  CampaignType = "BUNDLE_DEAL"
	CampaignTypeBuyXGetY    CampaignType = "BUY_X_GET_Y"
	CampaignTypeFlashSale   CampaignType = "FLASH_SALE"
	CampaignTypeFreeGift    CampaignType = "FREE_GIFT"
)

type Campaign struct {
	ID          string
	Name        string
	Description string
	Type        CampaignType
	Status      CampaignStatus
	Banner      string
	StartTime   time.Time
	EndTime     time.Time
	Budget      float64
	SpentBudget float64
	Rules       CampaignRules
	Products    []string // Featured product IDs
	Categories  []string // Featured category IDs
	Stats       CampaignStats
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

type CampaignRules struct {
	DiscountPercent   int
	MinPurchase       float64
	MaxDiscount       float64
	BundleProducts    []string
	BundlePrice       float64
	BuyQuantity       int
	GetQuantity       int
	GetProductID      string
	FreeGiftProductID string
	FreeGiftMinOrder  float64
}

type CampaignStats struct {
	Impressions   int
	Clicks        int
	Conversions   int
	Revenue       float64
	AvgOrderValue float64
}

func NewCampaign(name, description string, campaignType CampaignType, start, end time.Time, budget float64) *Campaign {
	return &Campaign{
		ID:          uuid.New().String(),
		Name:        name,
		Description: description,
		Type:        campaignType,
		Status:      CampaignStatusDraft,
		StartTime:   start,
		EndTime:     end,
		Budget:      budget,
		Products:    []string{},
		Categories:  []string{},
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}
}

func (c *Campaign) Schedule() {
	c.Status = CampaignStatusScheduled
	c.UpdatedAt = time.Now()
}

func (c *Campaign) Activate() {
	c.Status = CampaignStatusActive
	c.UpdatedAt = time.Now()
}

func (c *Campaign) Pause() {
	c.Status = CampaignStatusPaused
	c.UpdatedAt = time.Now()
}

func (c *Campaign) End() {
	c.Status = CampaignStatusEnded
	c.UpdatedAt = time.Now()
}

func (c *Campaign) IsActive() bool {
	now := time.Now()
	return c.Status == CampaignStatusActive && 
		now.After(c.StartTime) && 
		now.Before(c.EndTime) &&
		c.SpentBudget < c.Budget
}

func (c *Campaign) SetRules(rules CampaignRules) {
	c.Rules = rules
	c.UpdatedAt = time.Now()
}

func (c *Campaign) AddProducts(productIDs []string) {
	c.Products = append(c.Products, productIDs...)
	c.UpdatedAt = time.Now()
}

func (c *Campaign) RecordConversion(orderValue float64) {
	c.Stats.Conversions++
	c.Stats.Revenue += orderValue
	c.SpentBudget += orderValue * float64(c.Rules.DiscountPercent) / 100
	if c.Stats.Conversions > 0 {
		c.Stats.AvgOrderValue = c.Stats.Revenue / float64(c.Stats.Conversions)
	}
	c.UpdatedAt = time.Now()
}

type Repository interface {
	Save(ctx interface{}, campaign *Campaign) error
	FindByID(ctx interface{}, campaignID string) (*Campaign, error)
	Update(ctx interface{}, campaign *Campaign) error
	FindActive(ctx interface{}) ([]*Campaign, error)
	FindByType(ctx interface{}, campaignType CampaignType) ([]*Campaign, error)
	FindByProduct(ctx interface{}, productID string) ([]*Campaign, error)
}
