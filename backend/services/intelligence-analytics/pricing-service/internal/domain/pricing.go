package domain

import (
	"time"

	"github.com/google/uuid"
)

type PricingStrategy string

const (
	PricingStrategyFixed       PricingStrategy = "FIXED"
	PricingStrategyDynamic     PricingStrategy = "DYNAMIC"
	PricingStrategyCompetitive PricingStrategy = "COMPETITIVE"
	PricingStrategySurge       PricingStrategy = "SURGE"
)

type ProductPrice struct {
	ProductID      string
	BasePrice      float64
	CurrentPrice   float64
	MinPrice       float64
	MaxPrice       float64
	Strategy       PricingStrategy
	Demand         DemandMetrics
	Competition    CompetitorData
	LastUpdated    time.Time
}

type DemandMetrics struct {
	ViewsLast24h        int
	ViewsLast7d         int
	PurchasesLast24h    int
	PurchasesLast7d     int
	CartAddsLast24h     int
	ConversionRate      float64
	InventoryLevel      float64 // 0-1, lower = less stock
	DemandScore         float64 // 0-100
}

type CompetitorData struct {
	LowestPrice    float64
	AveragePrice   float64
	HighestPrice   float64
	CompetitorCount int
	LastScraped    time.Time
}

type PriceHistory struct {
	ID         string
	ProductID  string
	OldPrice   float64
	NewPrice   float64
	Reason     string
	Strategy   PricingStrategy
	CreatedAt  time.Time
}

type PricingRule struct {
	ID              string
	Name            string
	ProductIDs      []string
	CategoryIDs     []string
	Strategy        PricingStrategy
	Parameters      PricingParameters
	IsActive        bool
	Priority        int
	CreatedAt       time.Time
}

type PricingParameters struct {
	// Dynamic pricing
	DemandMultiplierMin float64 // e.g., 0.9 (10% discount at low demand)
	DemandMultiplierMax float64 // e.g., 1.3 (30% surge at high demand)
	
	// Competitive pricing
	TargetPosition     string  // "lowest", "average", "premium"
	PriceMargin        float64 // % below/above competition
	
	// Surge pricing
	SurgeThreshold     float64 // demand score to trigger surge
	SurgeMultiplier    float64 // max surge multiplier
	
	// Time-based
	TimeOfDayFactor    bool
	DayOfWeekFactor    bool
}

func NewProductPrice(productID string, basePrice, minPrice, maxPrice float64) *ProductPrice {
	return &ProductPrice{
		ProductID:    productID,
		BasePrice:    basePrice,
		CurrentPrice: basePrice,
		MinPrice:     minPrice,
		MaxPrice:     maxPrice,
		Strategy:     PricingStrategyFixed,
		LastUpdated:  time.Now(),
	}
}

func (p *ProductPrice) SetPrice(newPrice float64, reason string) *PriceHistory {
	// Enforce bounds
	if newPrice < p.MinPrice {
		newPrice = p.MinPrice
	}
	if newPrice > p.MaxPrice {
		newPrice = p.MaxPrice
	}

	history := &PriceHistory{
		ID:        uuid.New().String(),
		ProductID: p.ProductID,
		OldPrice:  p.CurrentPrice,
		NewPrice:  newPrice,
		Reason:    reason,
		Strategy:  p.Strategy,
		CreatedAt: time.Now(),
	}

	p.CurrentPrice = newPrice
	p.LastUpdated = time.Now()

	return history
}

func (p *ProductPrice) CalculateDynamicPrice(rule *PricingRule) float64 {
	params := rule.Parameters
	demand := p.Demand.DemandScore / 100.0 // normalize to 0-1

	// Linear interpolation between min and max multiplier based on demand
	multiplier := params.DemandMultiplierMin + demand*(params.DemandMultiplierMax-params.DemandMultiplierMin)

	return p.BasePrice * multiplier
}

func (p *ProductPrice) CalculateCompetitivePrice(rule *PricingRule) float64 {
	if p.Competition.CompetitorCount == 0 {
		return p.BasePrice
	}

	params := rule.Parameters
	var targetPrice float64

	switch params.TargetPosition {
	case "lowest":
		targetPrice = p.Competition.LowestPrice * (1 - params.PriceMargin/100)
	case "average":
		targetPrice = p.Competition.AveragePrice * (1 - params.PriceMargin/100)
	case "premium":
		targetPrice = p.Competition.HighestPrice * (1 + params.PriceMargin/100)
	default:
		targetPrice = p.Competition.AveragePrice
	}

	return targetPrice
}

func (p *ProductPrice) CalculateSurgePrice(rule *PricingRule) float64 {
	params := rule.Parameters
	
	if p.Demand.DemandScore < params.SurgeThreshold {
		return p.CurrentPrice
	}

	// Calculate surge multiplier based on demand above threshold
	excessDemand := (p.Demand.DemandScore - params.SurgeThreshold) / (100 - params.SurgeThreshold)
	surgeMultiplier := 1 + excessDemand*(params.SurgeMultiplier-1)

	return p.BasePrice * surgeMultiplier
}

type Repository interface {
	SavePrice(ctx interface{}, price *ProductPrice) error
	GetPrice(ctx interface{}, productID string) (*ProductPrice, error)
	UpdatePrice(ctx interface{}, price *ProductPrice) error
	SaveHistory(ctx interface{}, history *PriceHistory) error
	GetPriceHistory(ctx interface{}, productID string, limit int) ([]*PriceHistory, error)
	GetActiveRules(ctx interface{}) ([]*PricingRule, error)
	SaveRule(ctx interface{}, rule *PricingRule) error
	UpdateDemandMetrics(ctx interface{}, productID string, metrics *DemandMetrics) error
	UpdateCompetitorData(ctx interface{}, productID string, data *CompetitorData) error
}
