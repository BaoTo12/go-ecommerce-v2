package domain

import (
	"time"
)

// TrackingStatus represents the status of an order in tracking
type TrackingStatus string

const (
	StatusProcessing     TrackingStatus = "processing"
	StatusShipped        TrackingStatus = "shipped"
	StatusInTransit      TrackingStatus = "in_transit"
	StatusOutForDelivery TrackingStatus = "out_for_delivery"
	StatusDelivered      TrackingStatus = "delivered"
)

// OrderTracking represents complete tracking information for an order
type OrderTracking struct {
	OrderID           string         `json:"order_id"`
	Status            TrackingStatus `json:"status"`
	EstimatedDelivery string         `json:"estimated_delivery"`
	Carrier           string         `json:"carrier"`
	TrackingNumber    string         `json:"tracking_number"`
	CurrentLocation   *Location      `json:"current_location"`
	Steps             []TrackingStep `json:"steps"`
	Driver            *Driver        `json:"driver,omitempty"`
}

// Location represents a geographic location
type Location struct {
	Lat     float64 `json:"lat"`
	Lng     float64 `json:"lng"`
	Address string  `json:"address"`
}

// Driver represents the delivery driver
type Driver struct {
	ID        string `json:"id"`
	Name      string `json:"name"`
	Phone     string `json:"phone"`
	Vehicle   string `json:"vehicle"`
	PlateNo   string `json:"plate_no"`
	AvatarURL string `json:"avatar_url"`
}

// TrackingStep represents a step in the tracking timeline
type TrackingStep struct {
	Status      string `json:"status"`
	Description string `json:"description"`
	Time        string `json:"time"`
	Location    string `json:"location,omitempty"`
	Completed   bool   `json:"completed"`
}

// NewOrderTracking creates a new tracking record
func NewOrderTracking(orderID, carrier string) *OrderTracking {
	return &OrderTracking{
		OrderID:        orderID,
		Status:         StatusProcessing,
		Carrier:        carrier,
		TrackingNumber: "SPXVN" + orderID[3:],
		Steps:          make([]TrackingStep, 0),
	}
}

// UpdateDriverLocation updates the driver's current location
func (t *OrderTracking) UpdateDriverLocation(lat, lng float64, address string) {
	t.CurrentLocation = &Location{
		Lat:     lat,
		Lng:     lng,
		Address: address,
	}
}

// AddStep adds a new tracking step
func (t *OrderTracking) AddStep(status, description, location string) {
	step := TrackingStep{
		Status:      status,
		Description: description,
		Time:        time.Now().Format("15:04"),
		Location:    location,
		Completed:   true,
	}
	// Prepend to keep newest first
	t.Steps = append([]TrackingStep{step}, t.Steps...)
}

// Repository interface for tracking data
type Repository interface {
	GetTrackingByOrderID(ctx interface{}, orderID string) (*OrderTracking, error)
	SaveTracking(ctx interface{}, tracking *OrderTracking) error
	UpdateDriverLocation(ctx interface{}, orderID string, lat, lng float64) error
}
