package domain
}
	return t.CurrentStatus == EventTypeDelivered
func (t *TrackingInfo) IsDelivered() bool {

}
	return &t.Events[len(t.Events)-1]
	}
		return nil
	if len(t.Events) == 0 {
func (t *TrackingInfo) GetLatestEvent() *TrackingEvent {

}
	return event

	t.UpdatedAt = time.Now()
	t.CurrentLocation = location
	t.CurrentStatus = eventType
	t.Events = append(t.Events, *event)

	}
		FacilityName:   facilityName,
		Carrier:        t.Carrier,
		Timestamp:      time.Now(),
		Description:    description,
		Location:       location,
		EventType:      eventType,
		TrackingNumber: t.TrackingNumber,
		EventID:        uuid.New().String(),
	event := &TrackingEvent{
func (t *TrackingInfo) AddEvent(eventType TrackingEventType, location Location, description, facilityName string) *TrackingEvent {

}
	}, nil
		Events:            []TrackingEvent{},
		UpdatedAt:         now,
		CreatedAt:         now,
		EstimatedDelivery: now.Add(72 * time.Hour), // 3 days default
		CurrentStatus:     EventTypePickedUp,
		Destination:       destination,
		Origin:            origin,
		Carrier:           carrier,
		ShipmentID:        shipmentID,
		TrackingNumber:    trackingNumber,
	return &TrackingInfo{
	now := time.Now()

	}
		return nil, errors.New(errors.ErrInvalidInput, "shipment ID is required")
	if shipmentID == "" {
	}
		return nil, errors.New(errors.ErrInvalidInput, "tracking number is required")
	if trackingNumber == "" {
func NewTrackingInfo(trackingNumber, shipmentID, carrier, origin, destination string) (*TrackingInfo, error) {

}
	Events            []TrackingEvent
	UpdatedAt         time.Time
	CreatedAt         time.Time
	EstimatedDelivery time.Time
	CurrentLocation   Location
	CurrentStatus     TrackingEventType
	Destination       string
	Origin            string
	Carrier           string
	ShipmentID        string
	TrackingNumber    string
type TrackingInfo struct {

}
	FacilityName   string
	Carrier        string
	Timestamp      time.Time
	Description    string
	Location       Location
	EventType      TrackingEventType
	TrackingNumber string
	EventID        string
type TrackingEvent struct {

}
	Longitude float64
	Latitude  float64
	Country   string
	State     string
	City      string
type Location struct {

)
	EventTypeReturned        TrackingEventType = "RETURNED"
	EventTypeException       TrackingEventType = "EXCEPTION"
	EventTypeDelivered       TrackingEventType = "DELIVERED"
	EventTypeOutForDelivery  TrackingEventType = "OUT_FOR_DELIVERY"
	EventTypeAtFacility      TrackingEventType = "AT_FACILITY"
	EventTypeInTransit       TrackingEventType = "IN_TRANSIT"
	EventTypePickedUp        TrackingEventType = "PICKED_UP"
const (

type TrackingEventType string

)
	"github.com/titan-commerce/backend/pkg/errors"
	"github.com/google/uuid"

	"time"
import (


