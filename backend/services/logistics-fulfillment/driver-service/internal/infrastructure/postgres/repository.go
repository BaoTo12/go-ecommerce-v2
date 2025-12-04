package postgres

import (
	"context"
	"database/sql"
	"encoding/json"
	"time"

	"github.com/titan-commerce/backend/driver-service/internal/domain"
	"github.com/titan-commerce/backend/pkg/errors"
	"github.com/titan-commerce/backend/pkg/logger"
	_ "github.com/lib/pq"
)

type DriverRepository struct {
	db     *sql.DB
	logger *logger.Logger
}

type DeliveryRepository struct {
	db     *sql.DB
	logger *logger.Logger
}

type RouteRepository struct {
	db     *sql.DB
	logger *logger.Logger
}

func NewDriverRepository(databaseURL string, logger *logger.Logger) (*DriverRepository, error) {
	db, err := sql.Open("postgres", databaseURL)
	if err != nil {
		return nil, errors.Wrap(errors.ErrInternal, "failed to connect to database", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	
	if err := db.PingContext(ctx); err != nil {
		return nil, errors.Wrap(errors.ErrInternal, "failed to ping database", err)
	}

	logger.Info("Driver PostgreSQL repository initialized")
	return &DriverRepository{db: db, logger: logger}, nil
}

func NewDeliveryRepository(db *sql.DB, logger *logger.Logger) *DeliveryRepository {
	return &DeliveryRepository{db: db, logger: logger}
}

func NewRouteRepository(db *sql.DB, logger *logger.Logger) *RouteRepository {
	return &RouteRepository{db: db, logger: logger}
}

// DriverRepository implementations
func (r *DriverRepository) CreateDriver(ctx context.Context, driver *domain.Driver) error {
	query := `
		INSERT INTO drivers (driver_id, name, phone, email, vehicle_type, license_plate, status, current_lat, current_lng, rating, total_deliveries, successful_deliveries, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14)
	`
	_, err := r.db.ExecContext(ctx, query,
		driver.DriverID, driver.Name, driver.Phone, driver.Email, driver.VehicleType,
		driver.LicensePlate, driver.Status, driver.CurrentLat, driver.CurrentLng,
		driver.Rating, driver.TotalDeliveries, driver.SuccessfulDeliveries,
		driver.CreatedAt, driver.UpdatedAt,
	)
	return err
}

func (r *DriverRepository) GetDriver(ctx context.Context, driverID string) (*domain.Driver, error) {
	query := `SELECT driver_id, name, phone, email, vehicle_type, license_plate, status, current_lat, current_lng, rating, total_deliveries, successful_deliveries, created_at, updated_at FROM drivers WHERE driver_id = $1`
	
	var driver domain.Driver
	err := r.db.QueryRowContext(ctx, query, driverID).Scan(
		&driver.DriverID, &driver.Name, &driver.Phone, &driver.Email, &driver.VehicleType,
		&driver.LicensePlate, &driver.Status, &driver.CurrentLat, &driver.CurrentLng,
		&driver.Rating, &driver.TotalDeliveries, &driver.SuccessfulDeliveries,
		&driver.CreatedAt, &driver.UpdatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, errors.New(errors.ErrNotFound, "driver not found")
	}
	return &driver, err
}

func (r *DriverRepository) UpdateDriver(ctx context.Context, driver *domain.Driver) error {
	query := `
		UPDATE drivers 
		SET name = $2, phone = $3, email = $4, vehicle_type = $5, license_plate = $6, 
		    status = $7, current_lat = $8, current_lng = $9, rating = $10, 
		    total_deliveries = $11, successful_deliveries = $12, updated_at = $13
		WHERE driver_id = $1
	`
	_, err := r.db.ExecContext(ctx, query,
		driver.DriverID, driver.Name, driver.Phone, driver.Email, driver.VehicleType,
		driver.LicensePlate, driver.Status, driver.CurrentLat, driver.CurrentLng,
		driver.Rating, driver.TotalDeliveries, driver.SuccessfulDeliveries, driver.UpdatedAt,
	)
	return err
}

func (r *DriverRepository) GetAvailableDrivers(ctx context.Context, lat, lng float64, radiusKm int) ([]*domain.Driver, error) {
	// Simplified query - in production would use PostGIS for geo queries
	query := `SELECT driver_id, name, phone, email, vehicle_type, license_plate, status, current_lat, current_lng, rating, total_deliveries, successful_deliveries, created_at, updated_at 
			  FROM drivers WHERE status = $1 LIMIT 20`
	
	rows, err := r.db.QueryContext(ctx, query, domain.DriverAvailable)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var drivers []*domain.Driver
	for rows.Next() {
		var driver domain.Driver
		if err := rows.Scan(&driver.DriverID, &driver.Name, &driver.Phone, &driver.Email,
			&driver.VehicleType, &driver.LicensePlate, &driver.Status, &driver.CurrentLat,
			&driver.CurrentLng, &driver.Rating, &driver.TotalDeliveries, &driver.SuccessfulDeliveries,
			&driver.CreatedAt, &driver.UpdatedAt); err != nil {
			return nil, err
		}
		drivers = append(drivers, &driver)
	}
	return drivers, nil
}

// DeliveryRepository implementations
func (r *DeliveryRepository) CreateDelivery(ctx context.Context, delivery *domain.Delivery) error {
	addressJSON, _ := json.Marshal(delivery.DeliveryAddress)
	
	query := `
		INSERT INTO deliveries (delivery_id, order_id, driver_id, route_id, delivery_address, status, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
	`
	_, err := r.db.ExecContext(ctx, query,
		delivery.DeliveryID, delivery.OrderID, delivery.DriverID, delivery.RouteID,
		addressJSON, delivery.Status, delivery.CreatedAt, delivery.UpdatedAt,
	)
	return err
}

func (r *DeliveryRepository) GetDelivery(ctx context.Context, deliveryID string) (*domain.Delivery, error) {
	query := `SELECT delivery_id, order_id, driver_id, route_id, delivery_address, status, created_at, updated_at, picked_up_at, delivered_at FROM deliveries WHERE delivery_id = $1`
	
	var delivery domain.Delivery
	var addressJSON []byte
	
	err := r.db.QueryRowContext(ctx, query, deliveryID).Scan(
		&delivery.DeliveryID, &delivery.OrderID, &delivery.DriverID, &delivery.RouteID,
		&addressJSON, &delivery.Status, &delivery.CreatedAt, &delivery.UpdatedAt,
		&delivery.PickedUpAt, &delivery.DeliveredAt,
	)
	if err == sql.ErrNoRows {
		return nil, errors.New(errors.ErrNotFound, "delivery not found")
	}
	if err != nil {
		return nil, err
	}

	json.Unmarshal(addressJSON, &delivery.DeliveryAddress)
	return &delivery, nil
}

func (r *DeliveryRepository) UpdateDelivery(ctx context.Context, delivery *domain.Delivery) error {
	addressJSON, _ := json.Marshal(delivery.DeliveryAddress)
	
	query := `
		UPDATE deliveries 
		SET driver_id = $2, route_id = $3, delivery_address = $4, status = $5, 
		    updated_at = $6, picked_up_at = $7, delivered_at = $8
		WHERE delivery_id = $1
	`
	_, err := r.db.ExecContext(ctx, query,
		delivery.DeliveryID, delivery.DriverID, delivery.RouteID, addressJSON,
		delivery.Status, delivery.UpdatedAt, delivery.PickedUpAt, delivery.DeliveredAt,
	)
	return err
}

func (r *DeliveryRepository) GetDeliveriesByRoute(ctx context.Context, routeID string) ([]*domain.Delivery, error) {
	query := `SELECT delivery_id, order_id, driver_id, route_id, delivery_address, status, created_at, updated_at, picked_up_at, delivered_at 
			  FROM deliveries WHERE route_id = $1 ORDER BY created_at`
	
	rows, err := r.db.QueryContext(ctx, query, routeID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var deliveries []*domain.Delivery
	for rows.Next() {
		var delivery domain.Delivery
		var addressJSON []byte
		
		if err := rows.Scan(&delivery.DeliveryID, &delivery.OrderID, &delivery.DriverID,
			&delivery.RouteID, &addressJSON, &delivery.Status, &delivery.CreatedAt,
			&delivery.UpdatedAt, &delivery.PickedUpAt, &delivery.DeliveredAt); err != nil {
			return nil, err
		}
		
		json.Unmarshal(addressJSON, &delivery.DeliveryAddress)
		deliveries = append(deliveries, &delivery)
	}
	return deliveries, nil
}

// RouteRepository implementations
func (r *RouteRepository) CreateRoute(ctx context.Context, route *domain.DeliveryRoute) error {
	stopsJSON, _ := json.Marshal(route.Stops)
	
	query := `
		INSERT INTO delivery_routes (route_id, driver_id, stops, total_distance_km, estimated_duration_min, status, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
	`
	_, err := r.db.ExecContext(ctx, query,
		route.RouteID, route.DriverID, stopsJSON, route.TotalDistanceKm,
		route.EstimatedDurationMin, route.Status, route.CreatedAt, route.UpdatedAt,
	)
	return err
}

func (r *RouteRepository) GetRoute(ctx context.Context, routeID string) (*domain.DeliveryRoute, error) {
	query := `SELECT route_id, driver_id, stops, total_distance_km, estimated_duration_min, status, created_at, updated_at, started_at, completed_at FROM delivery_routes WHERE route_id = $1`
	
	var route domain.DeliveryRoute
	var stopsJSON []byte
	
	err := r.db.QueryRowContext(ctx, query, routeID).Scan(
		&route.RouteID, &route.DriverID, &stopsJSON, &route.TotalDistanceKm,
		&route.EstimatedDurationMin, &route.Status, &route.CreatedAt, &route.UpdatedAt,
		&route.StartedAt, &route.CompletedAt,
	)
	if err == sql.ErrNoRows {
		return nil, errors.New(errors.ErrNotFound, "route not found")
	}
	if err != nil {
		return nil, err
	}

	json.Unmarshal(stopsJSON, &route.Stops)
	return &route, nil
}

func (r *RouteRepository) UpdateRoute(ctx context.Context, route *domain.DeliveryRoute) error {
	stopsJSON, _ := json.Marshal(route.Stops)
	
	query := `
		UPDATE delivery_routes 
		SET stops = $2, total_distance_km = $3, estimated_duration_min = $4, status = $5, 
		    updated_at = $6, started_at = $7, completed_at = $8
		WHERE route_id = $1
	`
	_, err := r.db.ExecContext(ctx, query,
		route.RouteID, stopsJSON, route.TotalDistanceKm, route.EstimatedDurationMin,
		route.Status, route.UpdatedAt, route.StartedAt, route.CompletedAt,
	)
	return err
}

func (r *RouteRepository) GetActiveRouteByDriver(ctx context.Context, driverID string) (*domain.DeliveryRoute, error) {
	query := `SELECT route_id, driver_id, stops, total_distance_km, estimated_duration_min, status, created_at, updated_at, started_at, completed_at 
			  FROM delivery_routes WHERE driver_id = $1 AND status IN ($2, $3) ORDER BY created_at DESC LIMIT 1`
	
	var route domain.DeliveryRoute
	var stopsJSON []byte
	
	err := r.db.QueryRowContext(ctx, query, driverID, domain.RouteAssigned, domain.RouteInProgress).Scan(
		&route.RouteID, &route.DriverID, &stopsJSON, &route.TotalDistanceKm,
		&route.EstimatedDurationMin, &route.Status, &route.CreatedAt, &route.UpdatedAt,
		&route.StartedAt, &route.CompletedAt,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	json.Unmarshal(stopsJSON, &route.Stops)
	return &route, nil
}

func (r *DriverRepository) Close() error {
	return r.db.Close()
}
