package application

import (
	"context"

	"github.com/titan-commerce/backend/tracking-service/internal/domain"
	"github.com/titan-commerce/backend/pkg/logger"
)

type TrackingService struct {
	repo   domain.TrackingRepository
	logger *logger.Logger
}

func NewTrackingService(repo domain.TrackingRepository, logger *logger.Logger) *TrackingService {
	return &TrackingService{
		repo:   repo,
		logger: logger,
	}
}

// CreateTracking creates a new tracking record (Command)
func (s *TrackingService) CreateTracking(ctx context.Context, trackingNumber, shipmentID, carrier, origin, destination string) error {
	tracking, err := domain.NewTrackingInfo(trackingNumber, shipmentID, carrier, origin, destination)
	if err != nil {
		return err
	}

	if err := s.repo.Save(ctx, tracking); err != nil {
		s.logger.Error(err, "failed to save tracking info")
		return err
	}

	s.logger.Infof("Tracking created: tracking_number=%s, shipment_id=%s", trackingNumber, shipmentID)
	return nil
}

// UpdateLocation adds a new tracking event (Command)
func (s *TrackingService) UpdateLocation(ctx context.Context, trackingNumber string, eventType domain.TrackingEventType, location domain.Location, description, facilityName string) (*domain.TrackingEvent, error) {
	tracking, err := s.repo.FindByTrackingNumber(ctx, trackingNumber)
	if err != nil {
		return nil, err
	}

	event := tracking.AddEvent(eventType, location, description, facilityName)

	// Save event to ScyllaDB
	if err := s.repo.AddEvent(ctx, event); err != nil {
		s.logger.Error(err, "failed to save tracking event")
		return nil, err
	}

	// Update tracking info
	if err := s.repo.Update(ctx, tracking); err != nil {
		s.logger.Error(err, "failed to update tracking info")
		return nil, err
	}

	s.logger.Infof("Tracking updated: tracking_number=%s, event_type=%s, location=%s",
		trackingNumber, eventType, location.City)

	return event, nil
}

// GetTrackingHistory retrieves complete tracking history (Query)
func (s *TrackingService) GetTrackingHistory(ctx context.Context, trackingNumber string) (*domain.TrackingInfo, error) {
	tracking, err := s.repo.FindByTrackingNumber(ctx, trackingNumber)
	if err != nil {
		return nil, err
	}

	// Get event history from ScyllaDB
	events, err := s.repo.GetEventHistory(ctx, trackingNumber)
	if err != nil {
		s.logger.Error(err, "failed to get event history")
		return nil, err
	}

	tracking.Events = events
	return tracking, nil
}

// GetCurrentStatus retrieves current tracking status (Query)
func (s *TrackingService) GetCurrentStatus(ctx context.Context, trackingNumber string) (*domain.TrackingInfo, error) {
	return s.repo.FindByTrackingNumber(ctx, trackingNumber)
}

