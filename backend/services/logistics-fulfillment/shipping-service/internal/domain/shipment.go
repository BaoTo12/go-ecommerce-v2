package domain

import (
	"time"

	"github.com/google/uuid"
	"github.com/titan-commerce/backend/pkg/errors"
)

type ShipmentStatus string

const (
	ShipmentStatusPending         ShipmentStatus = "PENDING"
	ShipmentStatusPickedUp        ShipmentStatus = "PICKED_UP"
	ShipmentStatusInTransit       ShipmentStatus = "IN_TRANSIT"
	ShipmentStatusOutForDelivery  ShipmentStatus = "OUT_FOR_DELIVERY"
	ShipmentStatusDelivered       ShipmentStatus = "DELIVERED"
	ShipmentStatusReturned        ShipmentStatus = "RETURNED"
)

type Shipment struct {
	ID                 string
	OrderID            string
	Carrier            string // DHL, FedEx, UPS, etc.
	TrackingNumber     string
	Status             ShipmentStatus
	OriginAddress      string
	DestinationAddress string
	Weight             float64
	ShippingCost       float64
	EstimatedDelivery  time.Time
	CreatedAt          time.Time
	UpdatedAt          time.Time
}

func NewShipment(orderID, carrier, originAddr, destAddr string, weight float64) (*Shipment, error) {
	if orderID == "" {
		return nil, errors.New(errors.ErrInvalidInput, "order ID is required")
	}
	if weight <= 0 {
		return nil, errors.New(errors.ErrInvalidInput, "weight must be positive")
	}

	now := time.Now()
	return &Shipment{
		ID:                 uuid.New().String(),
		OrderID:            orderID,
		Carrier:            carrier,
		TrackingNumber:     generateTrackingNumber(),
		Status:             ShipmentStatusPending,
		OriginAddress:      originAddr,
		DestinationAddress: destAddr,
		Weight:             weight,
		EstimatedDelivery:  now.Add(72 * time.Hour), // 3 days default
		CreatedAt:          now,
		UpdatedAt:          now,
	}, nil
}

func (s *Shipment) UpdateStatus(status ShipmentStatus) {
	s.Status = status
	s.UpdatedAt = time.Now()
}

func (s *Shipment) SetShippingCost(cost float64) {
	s.ShippingCost = cost
}

func generateTrackingNumber() string {
	return "TRK" + uuid.New().String()[:8]
}
