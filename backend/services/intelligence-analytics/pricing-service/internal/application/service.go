package application

import (
	"context"
	"time"

	"github.com/titan-commerce/backend/pricing-service/internal/domain"
	"github.com/titan-commerce/backend/pkg/logger"
)

type PricingRepository interface {
	SavePrice(ctx context.Context, price *domain.ProductPrice) error
	GetPrice(ctx context.Context, productID string) (*domain.ProductPrice, error)
	UpdatePrice(ctx context.Context, price *domain.ProductPrice) error
	SaveHistory(ctx context.Context, history *domain.PriceHistory) error
	GetPriceHistory(ctx context.Context, productID string, limit int) ([]*domain.PriceHistory, error)
	GetActiveRules(ctx context.Context) ([]*domain.PricingRule, error)
	SaveRule(ctx context.Context, rule *domain.PricingRule) error
	UpdateDemandMetrics(ctx context.Context, productID string, metrics *domain.DemandMetrics) error
	UpdateCompetitorData(ctx context.Context, productID string, data *domain.CompetitorData) error
	GetProductsByCategory(ctx context.Context, categoryID string) ([]string, error)
}

type PricingService struct {
	repo   PricingRepository
	logger *logger.Logger
}

func NewPricingService(repo PricingRepository, logger *logger.Logger) *PricingService {
	svc := &PricingService{
		repo:   repo,
		logger: logger,
	}

	// Start background price optimization
	go svc.runOptimizer()

	return svc
}

func (s *PricingService) runOptimizer() {
	ticker := time.NewTicker(5 * time.Minute)
	for range ticker.C {
		ctx := context.Background()
		s.optimizeAllPrices(ctx)
	}
}

func (s *PricingService) optimizeAllPrices(ctx context.Context) {
	rules, err := s.repo.GetActiveRules(ctx)
	if err != nil {
		s.logger.Error(err, "failed to get pricing rules")
		return
	}

	for _, rule := range rules {
		for _, productID := range rule.ProductIDs {
			s.OptimizePrice(ctx, productID)
		}
	}
}

// GetPrice returns the current price for a product
func (s *PricingService) GetPrice(ctx context.Context, productID string) (*domain.ProductPrice, error) {
	return s.repo.GetPrice(ctx, productID)
}

// SetBasePrice sets the base price for a product
func (s *PricingService) SetBasePrice(ctx context.Context, productID string, basePrice, minPrice, maxPrice float64) (*domain.ProductPrice, error) {
	price, err := s.repo.GetPrice(ctx, productID)
	if err != nil {
		price = domain.NewProductPrice(productID, basePrice, minPrice, maxPrice)
		if err := s.repo.SavePrice(ctx, price); err != nil {
			return nil, err
		}
	} else {
		price.BasePrice = basePrice
		price.MinPrice = minPrice
		price.MaxPrice = maxPrice
		if err := s.repo.UpdatePrice(ctx, price); err != nil {
			return nil, err
		}
	}

	s.logger.Infof("Base price set: product=%s, price=%.2f", productID, basePrice)
	return price, nil
}

// SetStrategy sets the pricing strategy for a product
func (s *PricingService) SetStrategy(ctx context.Context, productID string, strategy domain.PricingStrategy) error {
	price, err := s.repo.GetPrice(ctx, productID)
	if err != nil {
		return err
	}

	price.Strategy = strategy
	return s.repo.UpdatePrice(ctx, price)
}

// UpdateDemandMetrics updates demand data for dynamic pricing
func (s *PricingService) UpdateDemandMetrics(ctx context.Context, productID string, metrics *domain.DemandMetrics) error {
	// Calculate demand score
	metrics.DemandScore = s.calculateDemandScore(metrics)
	return s.repo.UpdateDemandMetrics(ctx, productID, metrics)
}

func (s *PricingService) calculateDemandScore(metrics *domain.DemandMetrics) float64 {
	score := 0.0

	// View velocity (weight: 30%)
	viewVelocity := float64(metrics.ViewsLast24h) / (float64(metrics.ViewsLast7d)/7 + 1)
	score += min(viewVelocity*10, 30)

	// Purchase velocity (weight: 40%)
	purchaseVelocity := float64(metrics.PurchasesLast24h) / (float64(metrics.PurchasesLast7d)/7 + 1)
	score += min(purchaseVelocity*20, 40)

	// Conversion rate (weight: 20%)
	score += metrics.ConversionRate * 20

	// Scarcity (weight: 10%)
	scarcity := 1 - metrics.InventoryLevel
	score += scarcity * 10

	return min(score, 100)
}

// UpdateCompetitorPrices updates competitor pricing data
func (s *PricingService) UpdateCompetitorPrices(ctx context.Context, productID string, competitors []float64) error {
	if len(competitors) == 0 {
		return nil
	}

	lowest := competitors[0]
	highest := competitors[0]
	sum := 0.0

	for _, price := range competitors {
		if price < lowest {
			lowest = price
		}
		if price > highest {
			highest = price
		}
		sum += price
	}

	data := &domain.CompetitorData{
		LowestPrice:     lowest,
		AveragePrice:    sum / float64(len(competitors)),
		HighestPrice:    highest,
		CompetitorCount: len(competitors),
		LastScraped:     time.Now(),
	}

	return s.repo.UpdateCompetitorData(ctx, productID, data)
}

// OptimizePrice recalculates the optimal price for a product
func (s *PricingService) OptimizePrice(ctx context.Context, productID string) (*domain.ProductPrice, error) {
	price, err := s.repo.GetPrice(ctx, productID)
	if err != nil {
		return nil, err
	}

	if price.Strategy == domain.PricingStrategyFixed {
		return price, nil
	}

	rules, err := s.repo.GetActiveRules(ctx)
	if err != nil {
		return nil, err
	}

	// Find applicable rule
	var applicableRule *domain.PricingRule
	for _, rule := range rules {
		for _, pid := range rule.ProductIDs {
			if pid == productID {
				applicableRule = rule
				break
			}
		}
	}

	if applicableRule == nil {
		return price, nil
	}

	// Calculate new price based on strategy
	var newPrice float64
	var reason string

	switch price.Strategy {
	case domain.PricingStrategyDynamic:
		newPrice = price.CalculateDynamicPrice(applicableRule)
		reason = "Dynamic pricing based on demand"
	case domain.PricingStrategyCompetitive:
		newPrice = price.CalculateCompetitivePrice(applicableRule)
		reason = "Competitive pricing adjustment"
	case domain.PricingStrategySurge:
		newPrice = price.CalculateSurgePrice(applicableRule)
		reason = "Surge pricing"
	default:
		return price, nil
	}

	// Apply price change
	history := price.SetPrice(newPrice, reason)

	if err := s.repo.UpdatePrice(ctx, price); err != nil {
		return nil, err
	}

	if err := s.repo.SaveHistory(ctx, history); err != nil {
		s.logger.Error(err, "failed to save price history")
	}

	s.logger.Infof("Price optimized: product=%s, old=%.2f, new=%.2f, strategy=%s",
		productID, history.OldPrice, history.NewPrice, price.Strategy)

	return price, nil
}

// GetPriceHistory returns price change history
func (s *PricingService) GetPriceHistory(ctx context.Context, productID string, limit int) ([]*domain.PriceHistory, error) {
	return s.repo.GetPriceHistory(ctx, productID, limit)
}

// CreateRule creates a new pricing rule
func (s *PricingService) CreateRule(ctx context.Context, rule *domain.PricingRule) error {
	return s.repo.SaveRule(ctx, rule)
}

// GetActiveRules returns all active pricing rules
func (s *PricingService) GetActiveRules(ctx context.Context) ([]*domain.PricingRule, error) {
	return s.repo.GetActiveRules(ctx)
}

func min(a, b float64) float64 {
	if a < b {
		return a
	}
	return b
}
