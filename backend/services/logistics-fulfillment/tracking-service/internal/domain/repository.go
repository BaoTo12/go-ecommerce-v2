package domain

import "context"

type TrackingRepository interface {
	Save(ctx context.Context, tracking *TrackingInfo) error
	FindByTrackingNumber(ctx context.Context, trackingNumber string) (*TrackingInfo, error)
	FindByShipmentID(ctx context.Context, shipmentID string) (*TrackingInfo, error)
	AddEvent(ctx context.Context, event *TrackingEvent) error
	GetEventHistory(ctx context.Context, trackingNumber string) ([]TrackingEvent, error)
	Update(ctx context.Context, tracking *TrackingInfo) error
}

