package postgres
}
	return err

	)
		shipment.ID,
		shipment.UpdatedAt,
		shipment.EstimatedDelivery,
		shipment.ShippingCost,
		shipment.Status,
	_, err := r.db.ExecContext(ctx, query,

	`
		WHERE shipment_id = $5
		SET status = $1, shipping_cost = $2, estimated_delivery = $3, updated_at = $4
		UPDATE shipments
	query := `
func (r *ShipmentRepository) Update(ctx context.Context, shipment *domain.Shipment) error {

}
	return &shipment, nil

	}
		return nil, err
	if err != nil {
	}
		return nil, errors.New(errors.ErrNotFound, "shipment not found")
	if err == sql.ErrNoRows {

	)
		&shipment.UpdatedAt,
		&shipment.CreatedAt,
		&shipment.EstimatedDelivery,
		&shipment.ShippingCost,
		&shipment.Weight,
		&shipment.DestinationAddress,
		&shipment.OriginAddress,
		&shipment.Status,
		&shipment.TrackingNumber,
		&shipment.Carrier,
		&shipment.OrderID,
		&shipment.ID,
	err := r.db.QueryRowContext(ctx, query, orderID).Scan(
	var shipment domain.Shipment

	`
		WHERE order_id = $1
		FROM shipments
			estimated_delivery, created_at, updated_at
			origin_address, destination_address, weight, shipping_cost,
		SELECT shipment_id, order_id, carrier, tracking_number, status,
	query := `
func (r *ShipmentRepository) FindByOrderID(ctx context.Context, orderID string) (*domain.Shipment, error) {

}
	return &shipment, nil

	}
		return nil, err
	if err != nil {
	}
		return nil, errors.New(errors.ErrNotFound, "shipment not found")
	if err == sql.ErrNoRows {

	)
		&shipment.UpdatedAt,
		&shipment.CreatedAt,
		&shipment.EstimatedDelivery,
		&shipment.ShippingCost,
		&shipment.Weight,
		&shipment.DestinationAddress,
		&shipment.OriginAddress,
		&shipment.Status,
		&shipment.TrackingNumber,
		&shipment.Carrier,
		&shipment.OrderID,
		&shipment.ID,
	err := r.db.QueryRowContext(ctx, query, shipmentID).Scan(
	var shipment domain.Shipment

	`
		WHERE shipment_id = $1
		FROM shipments
			estimated_delivery, created_at, updated_at
			origin_address, destination_address, weight, shipping_cost,
		SELECT shipment_id, order_id, carrier, tracking_number, status,
	query := `
func (r *ShipmentRepository) FindByID(ctx context.Context, shipmentID string) (*domain.Shipment, error) {

}
	return err

	)
		shipment.UpdatedAt,
		shipment.CreatedAt,
		shipment.EstimatedDelivery,
		shipment.ShippingCost,
		shipment.Weight,
		shipment.DestinationAddress,
		shipment.OriginAddress,
		shipment.Status,
		shipment.TrackingNumber,
		shipment.Carrier,
		shipment.OrderID,
		shipment.ID,
	_, err := r.db.ExecContext(ctx, query,

	`
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)
			estimated_delivery, created_at, updated_at
			origin_address, destination_address, weight, shipping_cost,
			shipment_id, order_id, carrier, tracking_number, status,
		INSERT INTO shipments (
	query := `
func (r *ShipmentRepository) Save(ctx context.Context, shipment *domain.Shipment) error {

}
	return &ShipmentRepository{db: db}
func NewShipmentRepository(db *sql.DB) *ShipmentRepository {

}
	db *sql.DB
type ShipmentRepository struct {

)
	"github.com/titan-commerce/backend/pkg/errors"
	"github.com/titan-commerce/backend/shipping-service/internal/domain"

	"database/sql"
	"context"
import (


