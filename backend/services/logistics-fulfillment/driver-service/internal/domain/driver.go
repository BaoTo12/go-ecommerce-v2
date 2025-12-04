package domain

import (
	"time"

	"github.com/google/uuid"
	"github.com/titan-commerce/backend/pkg/errors"
)

// DriverStatus represents the current status of a driver
type DriverStatus string

const (
	DriverAvailable  DriverStatus = "AVAILABLE"
	DriverOnDelivery DriverStatus = "ON_DELIVERY"
	DriverOnBreak    DriverStatus = "ON_BREAK"
	DriverOffDuty    DriverStatus = "OFF_DUTY"
)

// VehicleType represents the type of vehicle a driver uses
type VehicleType string

const (
	VehicleMotorcycle VehicleType = "MOTORCYCLE"
	VehicleCar        VehicleType = "CAR"
	VehicleVan        VehicleType = "VAN"
	VehicleTruck      VehicleType = "TRUCK"
)

// Driver represents a delivery driver in the system
type Driver struct {
	DriverID      string
	Name          string
	Phone         string
	Email         string
	VehicleType   VehicleType
	LicensePlate  string
	Status        DriverStatus
	CurrentLat    float64
	CurrentLng    float64
	Rating        float64
	TotalDeliveries int
	SuccessfulDeliveries int
	CreatedAt     time.Time
	UpdatedAt     time.Time
}

// NewDriver creates a new driver
func NewDriver(name, phone, email string, vehicleType VehicleType, licensePlate string) (*Driver, error) {
	if name == "" {
		return nil, errors.New(errors.ErrInvalidInput, "driver name is required")
	}
	if phone == "" {
		return nil, errors.New(errors.ErrInvalidInput, "driver phone is required")
	}
	if vehicleType == "" {
		return nil, errors.New(errors.ErrInvalidInput, "vehicle type is required")
	}

	now := time.Now()
	return &Driver{
		DriverID:             uuid.New().String(),
		Name:                 name,
		Phone:                phone,
		Email:                email,
		VehicleType:          vehicleType,
		LicensePlate:         licensePlate,
		Status:               DriverOffDuty,
		Rating:               5.0,
		TotalDeliveries:      0,
		SuccessfulDeliveries: 0,
		CreatedAt:            now,
		UpdatedAt:            now,
	}, nil
}

// UpdateStatus updates the driver's current status
func (d *Driver) UpdateStatus(status DriverStatus) error {
	if status == "" {
		return errors.New(errors.ErrInvalidInput, "status is required")
	}
	d.Status = status
	d.UpdatedAt = time.Now()
	return nil
}

// UpdateLocation updates the driver's current GPS location
func (d *Driver) UpdateLocation(lat, lng float64) error {
	if lat < -90 || lat > 90 {
		return errors.New(errors.ErrInvalidInput, "invalid latitude")
	}
	if lng < -180 || lng > 180 {
		return errors.New(errors.ErrInvalidInput, "invalid longitude")
	}
	d.CurrentLat = lat
	d.CurrentLng = lng
	d.UpdatedAt = time.Now()
	return nil
}

// RecordDeliveryCompletion records a successful delivery
func (d *Driver) RecordDeliveryCompletion(success bool) {
	d.TotalDeliveries++
	if success {
		d.SuccessfulDeliveries++
	}
	d.UpdatedAt = time.Now()
	d.updateRating()
}

// updateRating recalculates the driver's rating based on performance
func (d *Driver) updateRating() {
	if d.TotalDeliveries == 0 {
		d.Rating = 5.0
		return
	}
	successRate := float64(d.SuccessfulDeliveries) / float64(d.TotalDeliveries)
	d.Rating = 5.0 * successRate
}

// IsAvailableForAssignment checks if driver can be assigned new deliveries
func (d *Driver) IsAvailableForAssignment() bool {
	return d.Status == DriverAvailable
}

// DeliveryStatus represents the status of a single delivery
type DeliveryStatus string

const (
	DeliveryPending    DeliveryStatus = "PENDING"
	DeliveryAssigned   DeliveryStatus = "ASSIGNED"
	DeliveryInTransit  DeliveryStatus = "IN_TRANSIT"
	DeliveryAtLocation DeliveryStatus = "AT_LOCATION"
	DeliveryCompleted  DeliveryStatus = "COMPLETED"
	DeliveryFailed     DeliveryStatus = "FAILED"
	DeliveryRescheduled DeliveryStatus = "RESCHEDULED"
)

// Delivery represents a single delivery task
type Delivery struct {
	DeliveryID       string
	ShipmentID       string
	RouteID          string
	DriverID         string
	CustomerName     string
	CustomerPhone    string
	DeliveryAddress  string
	DeliveryLat      float64
	DeliveryLng      float64
	Status           DeliveryStatus
	SequenceNumber   int // Order in the route
	EstimatedArrival time.Time
	ActualArrival    *time.Time
	ProofOfDelivery  *ProofOfDelivery
	Notes            string
	CreatedAt        time.Time
	UpdatedAt        time.Time
}

// ProofOfDelivery represents evidence of delivery completion
type ProofOfDelivery struct {
	PhotoURL      string
	SignatureURL  string
	RecipientName string
	CompletedAt   time.Time
}

// NewDelivery creates a new delivery task
func NewDelivery(shipmentID, customerName, customerPhone, address string, lat, lng float64) (*Delivery, error) {
	if shipmentID == "" {
		return nil, errors.New(errors.ErrInvalidInput, "shipment ID is required")
	}
	if address == "" {
		return nil, errors.New(errors.ErrInvalidInput, "delivery address is required")
	}

	now := time.Now()
	return &Delivery{
		DeliveryID:      uuid.New().String(),
		ShipmentID:      shipmentID,
		CustomerName:    customerName,
		CustomerPhone:   customerPhone,
		DeliveryAddress: address,
		DeliveryLat:     lat,
		DeliveryLng:     lng,
		Status:          DeliveryPending,
		CreatedAt:       now,
		UpdatedAt:       now,
	}, nil
}

// AssignToRoute assigns delivery to a route
func (d *Delivery) AssignToRoute(routeID, driverID string, sequence int) error {
	if routeID == "" {
		return errors.New(errors.ErrInvalidInput, "route ID is required")
	}
	if driverID == "" {
		return errors.New(errors.ErrInvalidInput, "driver ID is required")
	}

	d.RouteID = routeID
	d.DriverID = driverID
	d.SequenceNumber = sequence
	d.Status = DeliveryAssigned
	d.UpdatedAt = time.Now()
	return nil
}

// UpdateStatus updates delivery status
func (d *Delivery) UpdateStatus(status DeliveryStatus) error {
	d.Status = status
	d.UpdatedAt = time.Now()

	if status == DeliveryAtLocation {
		now := time.Now()
		d.ActualArrival = &now
	}

	return nil
}

// RecordProofOfDelivery records POD when delivery is completed
func (d *Delivery) RecordProofOfDelivery(photoURL, signatureURL, recipientName string) error {
	if d.Status != DeliveryAtLocation {
		return errors.New(errors.ErrInvalidInput, "delivery must be at location to record POD")
	}

	now := time.Now()
	d.ProofOfDelivery = &ProofOfDelivery{
		PhotoURL:      photoURL,
		SignatureURL:  signatureURL,
		RecipientName: recipientName,
		CompletedAt:   now,
	}
	d.Status = DeliveryCompleted
	d.UpdatedAt = now
	return nil
}

// RouteStatus represents the status of a delivery route
type RouteStatus string

const (
	RouteCreated    RouteStatus = "CREATED"
	RouteAssigned   RouteStatus = "ASSIGNED"
	RouteInProgress RouteStatus = "IN_PROGRESS"
	RouteCompleted  RouteStatus = "COMPLETED"
)

// DeliveryRoute represents a collection of deliveries assigned to a driver
type DeliveryRoute struct {
	RouteID           string
	DriverID          string
	Status            RouteStatus
	TotalDeliveries   int
	CompletedDeliveries int
	EstimatedDuration int // minutes
	ActualDuration    *int
	StartedAt         *time.Time
	CompletedAt       *time.Time
	CreatedAt         time.Time
	UpdatedAt         time.Time
}

// NewDeliveryRoute creates a new delivery route
func NewDeliveryRoute(driverID string, totalDeliveries int) (*DeliveryRoute, error) {
	if driverID == "" {
		return nil, errors.New(errors.ErrInvalidInput, "driver ID is required")
	}
	if totalDeliveries <= 0 {
		return nil, errors.New(errors.ErrInvalidInput, "route must have at least one delivery")
	}

	now := time.Now()
	return &DeliveryRoute{
		RouteID:             uuid.New().String(),
		DriverID:            driverID,
		Status:              RouteCreated,
		TotalDeliveries:     totalDeliveries,
		CompletedDeliveries: 0,
		CreatedAt:           now,
		UpdatedAt:           now,
	}, nil
}

// Start marks the route as started
func (r *DeliveryRoute) Start() error {
	if r.Status != RouteCreated && r.Status != RouteAssigned {
		return errors.New(errors.ErrInvalidInput, "route already started or completed")
	}

	now := time.Now()
	r.StartedAt = &now
	r.Status = RouteInProgress
	r.UpdatedAt = now
	return nil
}

// RecordDeliveryCompletion records completion of a delivery in the route
func (r *DeliveryRoute) RecordDeliveryCompletion() error {
	r.CompletedDeliveries++
	r.UpdatedAt = time.Now()

	if r.CompletedDeliveries >= r.TotalDeliveries {
		return r.Complete()
	}

	return nil
}

// Complete marks the entire route as completed
func (r *DeliveryRoute) Complete() error {
	now := time.Now()
	r.CompletedAt = &now
	r.Status = RouteCompleted
	r.UpdatedAt = now

	if r.StartedAt != nil {
		duration := int(now.Sub(*r.StartedAt).Minutes())
		r.ActualDuration = &duration
	}

	return nil
}

// GetProgress returns the completion percentage
func (r *DeliveryRoute) GetProgress() float64 {
	if r.TotalDeliveries == 0 {
		return 0
	}
	return float64(r.CompletedDeliveries) / float64(r.TotalDeliveries) * 100
}

