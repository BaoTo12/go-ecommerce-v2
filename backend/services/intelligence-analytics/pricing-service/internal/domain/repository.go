package domain

import "context"

type PricingRepository interface {
	SaveRecommendation(ctx context.Context, rec *PriceRecommendation) error
	GetRecommendation(ctx context.Context, productID string) (*PriceRecommendation, error)
	SavePriceHistory(ctx context.Context, history *PriceHistory) error
	GetPriceHistory(ctx context.Context, productID string, limit int) ([]*PriceHistory, error)
	SaveElasticity(ctx context.Context, elasticity *PriceElasticity) error
	GetElasticity(ctx context.Context, productID string) (*PriceElasticity, error)
	SaveCompetitorPrice(ctx context.Context, comp *CompetitorPrice) error
	GetCompetitorPrices(ctx context.Context, productID string) ([]*CompetitorPrice, error)
	SaveRule(ctx context.Context, rule *DynamicPricingRule) error
	GetActiveRules(ctx context.Context, productID string) ([]*DynamicPricingRule, error)
}

