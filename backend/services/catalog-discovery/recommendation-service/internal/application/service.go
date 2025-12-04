package application

import (
	"context"

	"github.com/titan-commerce/backend/recommendation-service/internal/domain"
	"github.com/titan-commerce/backend/pkg/logger"
)

type RecommendationEngine interface {
	GetRecommendations(ctx context.Context, userID string, limit int, context string) ([]*domain.RecommendedItem, error)
	TrackInteraction(ctx context.Context, interaction *domain.Interaction) error
}

type RecommendationService struct {
	engine RecommendationEngine
	logger *logger.Logger
}

func NewRecommendationService(engine RecommendationEngine, logger *logger.Logger) *RecommendationService {
	return &RecommendationService{
		engine: engine,
		logger: logger,
	}
}

func (s *RecommendationService) GetRecommendations(ctx context.Context, userID string, limit int, context string) ([]*domain.RecommendedItem, error) {
	return s.engine.GetRecommendations(ctx, userID, limit, context)
}

func (s *RecommendationService) TrackInteraction(ctx context.Context, userID, productID, interactionType string) error {
	interaction := &domain.Interaction{
		UserID:          userID,
		ProductID:       productID,
		InteractionType: interactionType,
	}
	if err := s.engine.TrackInteraction(ctx, interaction); err != nil {
		s.logger.Error(err, "failed to track interaction")
		return err
	}
	s.logger.Infof("Tracked interaction: User %s %s Product %s", userID, interactionType, productID)
	return nil
}
