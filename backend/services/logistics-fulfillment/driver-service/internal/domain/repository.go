package domain
}
	GetRouteHistory(ctx context.Context, driverID string, limit int) ([]*DeliveryRoute, error)
	GetActiveRouteByDriver(ctx context.Context, driverID string) (*DeliveryRoute, error)
	UpdateRoute(ctx context.Context, route *DeliveryRoute) error
	GetRoute(ctx context.Context, routeID string) (*DeliveryRoute, error)
	CreateRoute(ctx context.Context, route *DeliveryRoute) error
	// Route operations
type RouteRepository interface {
// RouteRepository defines the interface for route persistence

}
	GetDeliveriesByShipment(ctx context.Context, shipmentID string) (*Delivery, error)
	GetDeliveriesByRoute(ctx context.Context, routeID string) ([]*Delivery, error)
	UpdateDelivery(ctx context.Context, delivery *Delivery) error
	GetDelivery(ctx context.Context, deliveryID string) (*Delivery, error)
	CreateDelivery(ctx context.Context, delivery *Delivery) error
	// Delivery operations
type DeliveryRepository interface {
// DeliveryRepository defines the interface for delivery persistence

}
	UpdateDriverLocation(ctx context.Context, driverID string, lat, lng float64) error
	UpdateDriverStatus(ctx context.Context, driverID string, status DriverStatus) error
	GetAvailableDrivers(ctx context.Context, lat, lng float64, radiusKm int) ([]*Driver, error)
	UpdateDriver(ctx context.Context, driver *Driver) error
	GetDriver(ctx context.Context, driverID string) (*Driver, error)
	CreateDriver(ctx context.Context, driver *Driver) error
	// Driver operations
type DriverRepository interface {
// DriverRepository defines the interface for driver persistence

import "context"


