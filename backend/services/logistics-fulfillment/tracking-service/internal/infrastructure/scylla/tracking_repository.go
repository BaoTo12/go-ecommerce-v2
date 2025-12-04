package scylla

import (
	"context"
	"time"

	"github.com/gocql/gocql"
	"github.com/titan-commerce/backend/tracking-service/internal/domain"
	"github.com/titan-commerce/backend/pkg/errors"
)

type TrackingRepository struct {
	session *gocql.Session
}

func NewTrackingRepository(session *gocql.Session) *TrackingRepository {
	return &TrackingRepository{
		session: session,
	}
}

// Save creates a new tracking record
func (r *TrackingRepository) Save(ctx context.Context, tracking *domain.TrackingInfo) error {
	query := `
		INSERT INTO tracking_info (
			tracking_number, shipment_id, carrier, origin, destination,
			current_status, current_location_city, current_location_state,
			current_location_country, current_location_lat, current_location_lon,
			estimated_delivery, created_at, updated_at
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	`

	return r.session.Query(query,
		tracking.TrackingNumber,
		tracking.ShipmentID,
		tracking.Carrier,
		tracking.Origin,
		tracking.Destination,
		tracking.CurrentStatus,
		tracking.CurrentLocation.City,
		tracking.CurrentLocation.State,
		tracking.CurrentLocation.Country,
		tracking.CurrentLocation.Latitude,
		tracking.CurrentLocation.Longitude,
		tracking.EstimatedDelivery,
		tracking.CreatedAt,
		tracking.UpdatedAt,
	).WithContext(ctx).Exec()
}

// FindByTrackingNumber retrieves tracking info by tracking number
func (r *TrackingRepository) FindByTrackingNumber(ctx context.Context, trackingNumber string) (*domain.TrackingInfo, error) {
	query := `
		SELECT tracking_number, shipment_id, carrier, origin, destination,
			current_status, current_location_city, current_location_state,
			current_location_country, current_location_lat, current_location_lon,
			estimated_delivery, created_at, updated_at
		FROM tracking_info
		WHERE tracking_number = ?
	`

	var tracking domain.TrackingInfo
	var location domain.Location

	err := r.session.Query(query, trackingNumber).
		WithContext(ctx).
		Scan(
			&tracking.TrackingNumber,
			&tracking.ShipmentID,
			&tracking.Carrier,
			&tracking.Origin,
			&tracking.Destination,
			&tracking.CurrentStatus,
			&location.City,
			&location.State,
			&location.Country,
			&location.Latitude,
			&location.Longitude,
			&tracking.EstimatedDelivery,
			&tracking.CreatedAt,
			&tracking.UpdatedAt,
		)

	if err != nil {
		if err == gocql.ErrNotFound {
			return nil, errors.New(errors.ErrNotFound, "tracking info not found")
		}
		return nil, err
	}

	tracking.CurrentLocation = location
	return &tracking, nil
}

// FindByShipmentID retrieves tracking info by shipment ID
func (r *TrackingRepository) FindByShipmentID(ctx context.Context, shipmentID string) (*domain.TrackingInfo, error) {
	query := `
		SELECT tracking_number, shipment_id, carrier, origin, destination,
			current_status, current_location_city, current_location_state,
			current_location_country, current_location_lat, current_location_lon,
			estimated_delivery, created_at, updated_at
		FROM tracking_info
		WHERE shipment_id = ?
		ALLOW FILTERING
	`

	var tracking domain.TrackingInfo
	var location domain.Location

	err := r.session.Query(query, shipmentID).
		WithContext(ctx).
		Scan(
			&tracking.TrackingNumber,
			&tracking.ShipmentID,
			&tracking.Carrier,
			&tracking.Origin,
			&tracking.Destination,
			&tracking.CurrentStatus,
			&location.City,
			&location.State,
			&location.Country,
			&location.Latitude,
			&location.Longitude,
			&tracking.EstimatedDelivery,
			&tracking.CreatedAt,
			&tracking.UpdatedAt,
		)

	if err != nil {
		if err == gocql.ErrNotFound {
			return nil, errors.New(errors.ErrNotFound, "tracking info not found")
		}
		return nil, err
	}

	tracking.CurrentLocation = location
	return &tracking, nil
}

// AddEvent adds a tracking event (time-series data)
func (r *TrackingRepository) AddEvent(ctx context.Context, event *domain.TrackingEvent) error {
	query := `
		INSERT INTO tracking_events (
			tracking_number, event_id, event_type, timestamp,
			location_city, location_state, location_country,
			location_lat, location_lon, description, carrier, facility_name
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	`

	return r.session.Query(query,
		event.TrackingNumber,
		event.EventID,
		event.EventType,
		event.Timestamp,
		event.Location.City,
		event.Location.State,
		event.Location.Country,
		event.Location.Latitude,
		event.Location.Longitude,
		event.Description,
		event.Carrier,
		event.FacilityName,
	).WithContext(ctx).Exec()
}

// GetEventHistory retrieves all events for a tracking number (ordered by timestamp)
func (r *TrackingRepository) GetEventHistory(ctx context.Context, trackingNumber string) ([]domain.TrackingEvent, error) {
	query := `
		SELECT event_id, tracking_number, event_type, timestamp,
			location_city, location_state, location_country,
			location_lat, location_lon, description, carrier, facility_name
		FROM tracking_events
		WHERE tracking_number = ?
		ORDER BY timestamp DESC
	`

	iter := r.session.Query(query, trackingNumber).WithContext(ctx).Iter()
	defer iter.Close()

	var events []domain.TrackingEvent
	for {
		var event domain.TrackingEvent
		var location domain.Location

		if !iter.Scan(
			&event.EventID,
			&event.TrackingNumber,
			&event.EventType,
			&event.Timestamp,
			&location.City,
			&location.State,
			&location.Country,
			&location.Latitude,
			&location.Longitude,
			&event.Description,
			&event.Carrier,
			&event.FacilityName,
		) {
			break
		}

		event.Location = location
		events = append(events, event)
	}

	if err := iter.Close(); err != nil {
		return nil, err
	}

	return events, nil
}

// Update updates tracking info
func (r *TrackingRepository) Update(ctx context.Context, tracking *domain.TrackingInfo) error {
	query := `
		UPDATE tracking_info
		SET current_status = ?,
			current_location_city = ?,
			current_location_state = ?,
			current_location_country = ?,
			current_location_lat = ?,
			current_location_lon = ?,
			estimated_delivery = ?,
			updated_at = ?
		WHERE tracking_number = ?
	`

	return r.session.Query(query,
		tracking.CurrentStatus,
		tracking.CurrentLocation.City,
		tracking.CurrentLocation.State,
		tracking.CurrentLocation.Country,
		tracking.CurrentLocation.Latitude,
		tracking.CurrentLocation.Longitude,
		tracking.EstimatedDelivery,
		time.Now(),
		tracking.TrackingNumber,
	).WithContext(ctx).Exec()
}

