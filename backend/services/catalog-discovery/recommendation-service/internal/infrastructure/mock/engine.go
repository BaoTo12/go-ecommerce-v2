package mock

import (
	"context"
	"math/rand"
	"time"

	"github.com/titan-commerce/backend/recommendation-service/internal/domain"
	"github.com/titan-commerce/backend/pkg/logger"
	"github.com/google/uuid"
)

type MockEngine struct {
	logger *logger.Logger
}

func NewMockEngine(logger *logger.Logger) *MockEngine {
	return &MockEngine{logger: logger}
}

func (e *MockEngine) GetRecommendations(ctx context.Context, userID string, limit int, context string) ([]*domain.RecommendedItem, error) {
	// Simulate AI computation latency
	time.Sleep(50 * time.Millisecond)

	var items []*domain.RecommendedItem
	reasons := []string{"Trending now", "Because you viewed electronics", "Top rated", "New arrival"}

	for i := 0; i < limit; i++ {
		items = append(items, &domain.RecommendedItem{
			ProductID: uuid.New().String(),
			Score:     rand.Float64(),
			Reason:    reasons[rand.Intn(len(reasons))],
		})
	}

	return items, nil
}

func (e *MockEngine) TrackInteraction(ctx context.Context, interaction *domain.Interaction) error {
	// In a real system, this would push to Kafka/Redpanda
	e.logger.Infof("[MOCK AI] Learning from interaction: %v", interaction)
	return nil
}
