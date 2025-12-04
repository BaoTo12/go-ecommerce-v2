package domain
}
	ChangedAt time.Time
	Reason    string
	Strategy  PricingStrategy
	Price     int64
	ProductID string
type PriceHistory struct {
// PriceHistory represents historical pricing data

}
	}
		UpdatedAt:  now,
		CreatedAt:  now,
		Priority:   100,
		IsActive:   true,
		Conditions: []PricingCondition{},
		MaxPrice:   maxPrice,
		MinPrice:   minPrice,
		ProductID:  productID,
		Name:       name,
		RuleID:     uuid.New().String(),
	return &DynamicPricingRule{
	now := time.Now()
func NewDynamicPricingRule(name, productID string, minPrice, maxPrice int64) *DynamicPricingRule {
// NewDynamicPricingRule creates a new pricing rule

}
	Operator  string // gt, lt, eq
	Threshold float64
	Type      string  // inventory_low, demand_high, competitor_lower, time_of_day
type PricingCondition struct {
// PricingCondition represents a condition for price adjustment

}
	UpdatedAt    time.Time
	CreatedAt    time.Time
	Priority     int
	IsActive     bool
	Conditions   []PricingCondition
	AdjustmentPct float64
	MaxPrice     int64
	MinPrice     int64
	CategoryID   string
	ProductID    string
	Name         string
	RuleID       string
type DynamicPricingRule struct {
// DynamicPricingRule represents business rules for dynamic pricing

}
	FetchedAt    time.Time
	URL          string
	Price        int64
	CompetitorID string
	ProductID    string
type CompetitorPrice struct {
// CompetitorPrice represents competitor pricing data

}
	Demand int
	Price  int64
type DemandPoint struct {
// DemandPoint represents a point on the demand curve

}
	UpdatedAt         time.Time
	DemandCurve       []DemandPoint
	OptimalPricePoint int64
	ElasticityScore   float64 // How demand changes with price
	ProductID         string
type PriceElasticity struct {
// PriceElasticity represents price elasticity for a product

}
	}
		CreatedAt:        now,
		ValidUntil:       now.Add(24 * time.Hour),
		ValidFrom:        now,
		Factors:          make(map[string]float64),
		Confidence:       0.8,
		Strategy:         strategy,
		RecommendedPrice: recommendedPrice,
		CurrentPrice:     currentPrice,
		ProductID:        productID,
		RecommendationID: uuid.New().String(),
	return &PriceRecommendation{
	now := time.Now()
func NewPriceRecommendation(productID string, currentPrice, recommendedPrice int64, strategy PricingStrategy) *PriceRecommendation {
// NewPriceRecommendation creates a new price recommendation

}
	CreatedAt        time.Time
	ValidUntil       time.Time
	ValidFrom        time.Time
	EstimatedRevenue int64
	EstimatedDemand  int
	Factors          map[string]float64 // demand, competition, inventory, etc.
	Confidence       float64 // 0.0 - 1.0
	Strategy         PricingStrategy
	RecommendedPrice int64
	CurrentPrice     int64 // in cents
	ProductID        string
	RecommendationID string
type PriceRecommendation struct {
// PriceRecommendation represents an ML-powered price recommendation

)
	StrategyPenetration  PricingStrategy = "PENETRATION"  // Market entry
	StrategySegmented    PricingStrategy = "SEGMENTED"    // User-segment based
	StrategyTimeBased    PricingStrategy = "TIME_BASED"   // Flash sales, time-of-day
	StrategyCompetitive  PricingStrategy = "COMPETITIVE"  // Match competitors
	StrategyDynamic      PricingStrategy = "DYNAMIC"      // Based on demand/supply
const (

type PricingStrategy string
// PricingStrategy represents different pricing algorithms

)
	"github.com/google/uuid"

	"time"
import (


