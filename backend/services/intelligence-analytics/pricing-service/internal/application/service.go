package application

import (
	"context"
	"math"

	"github.com/titan-commerce/backend/pricing-service/internal/domain"
	"github.com/titan-commerce/backend/pkg/logger"
)

type PricingService struct {
	repo   domain.PricingRepository
	logger *logger.Logger
}

func NewPricingService(repo domain.PricingRepository, logger *logger.Logger) *PricingService {
	return &PricingService{repo: repo, logger: logger}
}

// CalculateDynamicPrice calculates optimal price using ML/demand data
func (s *PricingService) CalculateDynamicPrice(
	ctx context.Context,
	productID string,
	currentPrice int64,
	demand int,
	inventory int,
	competitorPrices []int64,
) (*domain.PriceRecommendation, error) {

	// Simple dynamic pricing algorithm (in real-world, use ML model)
	avgCompetitorPrice := s.calculateAverage(competitorPrices)
	demandFactor := s.normalizeDemand(demand)
	inventoryFactor := s.normalizeInventory(inventory)

	// Calculate recommended price
	baseAdjustment := (demandFactor * 0.4) + (inventoryFactor * 0.3) + ((float64(avgCompetitorPrice)/float64(currentPrice) - 1) * 0.3)
	recommendedPrice := int64(float64(currentPrice) * (1 + baseAdjustment))

	// Apply min/max constraints
	rules, _ := s.repo.GetActiveRules(ctx, productID)
	if len(rules) > 0 {
		rule := rules[0]
		if recommendedPrice < rule.MinPrice {
			recommendedPrice = rule.MinPrice
		}
		if recommendedPrice > rule.MaxPrice {
			recommendedPrice = rule.MaxPrice
		}
	}

	rec := domain.NewPriceRecommendation(productID, currentPrice, recommendedPrice, domain.StrategyDynamic)
	rec.Confidence = 0.85
	rec.Factors = map[string]float64{
		"demand":     demandFactor,
		"inventory":  inventoryFactor,
		"competitor": float64(avgCompetitorPrice) / float64(currentPrice),
	}
	rec.EstimatedDemand = int(float64(demand) * (1 + (baseAdjustment * -0.5))) // Price elasticity

	s.repo.SaveRecommendation(ctx, rec)
	s.logger.Infof("Price recommendation: product=%s, current=%d, recommended=%d, confidence=%.2f",
		productID, currentPrice, recommendedPrice, rec.Confidence)

	return rec, nil
}

// CreatePricingRule creates a new dynamic pricing rule
func (s *PricingService) CreatePricingRule(
	ctx context.Context,
	name, productID string,
	minPrice, maxPrice int64,
) (*domain.DynamicPricingRule, error) {
	rule := domain.NewDynamicPricingRule(name, productID, minPrice, maxPrice)

	if err := s.repo.SaveRule(ctx, rule); err != nil {
		s.logger.Error(err, "failed to save pricing rule")
		return nil, err
	}

	s.logger.Infof("Pricing rule created: %s for product %s", name, productID)
	return rule, nil
}

// GetPriceRecommendation retrieves price recommendation
func (s *PricingService) GetPriceRecommendation(ctx context.Context, productID string) (*domain.PriceRecommendation, error) {
	return s.repo.GetRecommendation(ctx, productID)
}

func (s *PricingService) calculateAverage(prices []int64) int64 {
	if len(prices) == 0 {
		return 0
	}
	var sum int64
	for _, p := range prices {
		sum += p
	}
	return sum / int64(len(prices))
}

func (s *PricingService) normalizeDemand(demand int) float64 {
	// Normalize to -0.2 to +0.2 range
	return math.Tanh(float64(demand)/1000.0) * 0.2
}

func (s *PricingService) normalizeInventory(inventory int) float64 {
	// Low inventory = higher price
	if inventory < 10 {
		return 0.15
	} else if inventory < 50 {
		return 0.05
	} else if inventory > 500 {
		return -0.10
	}
	return 0.0
}

