package application

import (
	"context"
	"math"

	"github.com/titan-commerce/backend/driver-service/internal/domain"
	"github.com/titan-commerce/backend/pkg/logger"
)

type DriverService struct {
	driverRepo   domain.DriverRepository
	deliveryRepo domain.DeliveryRepository
	routeRepo    domain.RouteRepository
	logger       *logger.Logger
}

func NewDriverService(
	driverRepo domain.DriverRepository,
	deliveryRepo domain.DeliveryRepository,
	routeRepo domain.RouteRepository,
	logger *logger.Logger,
) *DriverService {
	return &DriverService{
		driverRepo:   driverRepo,
		deliveryRepo: deliveryRepo,
		routeRepo:    routeRepo,
		logger:       logger,
	}
}

// RegisterDriver registers a new driver (Command)
func (s *DriverService) RegisterDriver(
	ctx context.Context,
	name, phone, email string,
	vehicleType domain.VehicleType,
	licensePlate string,
) (*domain.Driver, error) {
	driver, err := domain.NewDriver(name, phone, email, vehicleType, licensePlate)
	if err != nil {
		return nil, err
	}

	if err := s.driverRepo.CreateDriver(ctx, driver); err != nil {
		s.logger.Error(err, "failed to create driver")
		return nil, err
	}

	s.logger.Infof("Driver registered: id=%s, name=%s", driver.DriverID, driver.Name)
	return driver, nil
}

// UpdateDriverStatus updates driver availability status (Command)
func (s *DriverService) UpdateDriverStatus(ctx context.Context, driverID string, status domain.DriverStatus) error {
	driver, err := s.driverRepo.GetDriver(ctx, driverID)
	if err != nil {
		return err
	}

	if err := driver.UpdateStatus(status); err != nil {
		return err
	}

	if err := s.driverRepo.UpdateDriver(ctx, driver); err != nil {
		s.logger.Error(err, "failed to update driver status")
		return err
	}

	s.logger.Infof("Driver status updated: id=%s, status=%s", driverID, status)
	return nil
}

// UpdateDriverLocation updates driver's GPS location (Command)
func (s *DriverService) UpdateDriverLocation(ctx context.Context, driverID string, lat, lng float64) error {
	driver, err := s.driverRepo.GetDriver(ctx, driverID)
	if err != nil {
		return err
	}

	if err := driver.UpdateLocation(lat, lng); err != nil {
		return err
	}

	if err := s.driverRepo.UpdateDriverLocation(ctx, driverID, lat, lng); err != nil {
		s.logger.Error(err, "failed to update driver location")
		return err
	}

	return nil
}

// CreateDelivery creates a new delivery task (Command)
func (s *DriverService) CreateDelivery(
	ctx context.Context,
	shipmentID, customerName, customerPhone, address string,
	lat, lng float64,
) (*domain.Delivery, error) {
	delivery, err := domain.NewDelivery(shipmentID, customerName, customerPhone, address, lat, lng)
	if err != nil {
		return nil, err
	}

	if err := s.deliveryRepo.CreateDelivery(ctx, delivery); err != nil {
		s.logger.Error(err, "failed to create delivery")
		return nil, err
	}

	s.logger.Infof("Delivery created: id=%s, shipment=%s", delivery.DeliveryID, shipmentID)
	return delivery, nil
}

// AssignRoute assigns deliveries to a driver and creates optimized route (Command)
func (s *DriverService) AssignRoute(
	ctx context.Context,
	driverID string,
	deliveryIDs []string,
) (*domain.DeliveryRoute, error) {
	// Verify driver is available
	driver, err := s.driverRepo.GetDriver(ctx, driverID)
	if err != nil {
		return nil, err
	}

	if !driver.IsAvailableForAssignment() {
		s.logger.Warnf("Driver not available: id=%s, status=%s", driverID, driver.Status)
		return nil, nil
	}

	// Create route
	route, err := domain.NewDeliveryRoute(driverID, len(deliveryIDs))
	if err != nil {
		return nil, err
	}

	if err := s.routeRepo.CreateRoute(ctx, route); err != nil {
		s.logger.Error(err, "failed to create route")
		return nil, err
	}

	// Assign deliveries to route with optimized sequence
	for i, deliveryID := range deliveryIDs {
		delivery, err := s.deliveryRepo.GetDelivery(ctx, deliveryID)
		if err != nil {
			continue
		}

		if err := delivery.AssignToRoute(route.RouteID, driverID, i+1); err != nil {
			s.logger.Error(err, "failed to assign delivery to route")
			continue
		}

		if err := s.deliveryRepo.UpdateDelivery(ctx, delivery); err != nil {
			s.logger.Error(err, "failed to update delivery")
		}
	}

	// Update driver status
	if err := driver.UpdateStatus(domain.DriverOnDelivery); err != nil {
		s.logger.Error(err, "failed to update driver status")
	}
	s.driverRepo.UpdateDriver(ctx, driver)

	s.logger.Infof("Route assigned: route=%s, driver=%s, deliveries=%d",
		route.RouteID, driverID, len(deliveryIDs))

	return route, nil
}

// StartRoute marks route as started (Command)
func (s *DriverService) StartRoute(ctx context.Context, routeID string) error {
	route, err := s.routeRepo.GetRoute(ctx, routeID)
	if err != nil {
		return err
	}

	if err := route.Start(); err != nil {
		return err
	}

	if err := s.routeRepo.UpdateRoute(ctx, route); err != nil {
		s.logger.Error(err, "failed to start route")
		return err
	}

	s.logger.Infof("Route started: route=%s, driver=%s", routeID, route.DriverID)
	return nil
}

// UpdateDeliveryStatus updates the status of a delivery (Command)
func (s *DriverService) UpdateDeliveryStatus(
	ctx context.Context,
	deliveryID string,
	status domain.DeliveryStatus,
) error {
	delivery, err := s.deliveryRepo.GetDelivery(ctx, deliveryID)
	if err != nil {
		return err
	}

	if err := delivery.UpdateStatus(status); err != nil {
		return err
	}

	if err := s.deliveryRepo.UpdateDelivery(ctx, delivery); err != nil {
		s.logger.Error(err, "failed to update delivery status")
		return err
	}

	s.logger.Infof("Delivery status updated: id=%s, status=%s", deliveryID, status)
	return nil
}

// RecordProofOfDelivery records proof of delivery (Command)
func (s *DriverService) RecordProofOfDelivery(
	ctx context.Context,
	deliveryID, photoURL, signatureURL, recipientName string,
) error {
	delivery, err := s.deliveryRepo.GetDelivery(ctx, deliveryID)
	if err != nil {
		return err
	}

	if err := delivery.RecordProofOfDelivery(photoURL, signatureURL, recipientName); err != nil {
		return err
	}

	if err := s.deliveryRepo.UpdateDelivery(ctx, delivery); err != nil {
		s.logger.Error(err, "failed to record POD")
		return err
	}

	// Update route progress
	route, err := s.routeRepo.GetRoute(ctx, delivery.RouteID)
	if err == nil {
		route.RecordDeliveryCompletion()
		s.routeRepo.UpdateRoute(ctx, route)

		// Update driver stats
		driver, err := s.driverRepo.GetDriver(ctx, delivery.DriverID)
		if err == nil {
			driver.RecordDeliveryCompletion(true)
			s.driverRepo.UpdateDriver(ctx, driver)
		}

		// If route is complete, mark driver as available
		if route.Status == domain.RouteCompleted {
			driver.UpdateStatus(domain.DriverAvailable)
			s.driverRepo.UpdateDriver(ctx, driver)
		}
	}

	s.logger.Infof("POD recorded: delivery=%s, recipient=%s", deliveryID, recipientName)
	return nil
}

// GetAvailableDrivers finds available drivers near a location (Query)
func (s *DriverService) GetAvailableDrivers(
	ctx context.Context,
	lat, lng float64,
	radiusKm int,
) ([]*domain.Driver, error) {
	return s.driverRepo.GetAvailableDrivers(ctx, lat, lng, radiusKm)
}

// GetDriver retrieves driver details (Query)
func (s *DriverService) GetDriver(ctx context.Context, driverID string) (*domain.Driver, error) {
	return s.driverRepo.GetDriver(ctx, driverID)
}

// GetDriverRoute gets driver's current route (Query)
func (s *DriverService) GetDriverRoute(ctx context.Context, driverID string) (*domain.DeliveryRoute, error) {
	return s.routeRepo.GetActiveRouteByDriver(ctx, driverID)
}

// GetDeliveryStatus gets delivery status (Query)
func (s *DriverService) GetDeliveryStatus(ctx context.Context, deliveryID string) (*domain.Delivery, error) {
	return s.deliveryRepo.GetDelivery(ctx, deliveryID)
}

// GetRouteDeliveries gets all deliveries in a route (Query)
func (s *DriverService) GetRouteDeliveries(ctx context.Context, routeID string) ([]*domain.Delivery, error) {
	return s.deliveryRepo.GetDeliveriesByRoute(ctx, routeID)
}

// Haversine formula to calculate distance between two GPS coordinates
func haversineDistance(lat1, lng1, lat2, lng2 float64) float64 {
	const earthRadiusKm = 6371

	dLat := (lat2 - lat1) * math.Pi / 180
	dLng := (lng2 - lng1) * math.Pi / 180

	lat1Rad := lat1 * math.Pi / 180
	lat2Rad := lat2 * math.Pi / 180

	a := math.Sin(dLat/2)*math.Sin(dLat/2) +
		math.Sin(dLng/2)*math.Sin(dLng/2)*math.Cos(lat1Rad)*math.Cos(lat2Rad)
	c := 2 * math.Atan2(math.Sqrt(a), math.Sqrt(1-a))

	return earthRadiusKm * c
}

