package application

import (
	"context"

	"github.com/titan-commerce/backend/analytics-service/internal/domain"
	"github.com/titan-commerce/backend/pkg/logger"
)

type AnalyticsService struct {
	repo   domain.AnalyticsRepository
	logger *logger.Logger
}

func NewAnalyticsService(repo domain.AnalyticsRepository, logger *logger.Logger) *AnalyticsService {
	return &AnalyticsService{repo: repo, logger: logger}
}

// TrackEvent records a business event
func (s *AnalyticsService) TrackEvent(
	ctx context.Context,
	eventType, userID, sessionID string,
	properties map[string]interface{},
) error {
	event := domain.NewEvent(eventType, userID, sessionID, properties)

	if err := s.repo.SaveEvent(ctx, event); err != nil {
		s.logger.Error(err, "failed to save event")
		return err
	}

	s.logger.Infof("Event tracked: type=%s, user=%s", eventType, userID)
	return nil
}

// RecordMetric records a metric value
func (s *AnalyticsService) RecordMetric(
	ctx context.Context,
	name, metricType string,
	value float64,
	dimensions map[string]string,
) error {
	metric := domain.NewMetric(name, metricType, value, dimensions)

	if err := s.repo.SaveMetric(ctx, metric); err != nil {
		s.logger.Error(err, "failed to save metric")
		return err
	}

	return nil
}

