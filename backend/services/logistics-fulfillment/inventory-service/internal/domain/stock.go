package domain
}
	return nil
	r.Status = ReservationRolledBack
	}
		return errors.New(errors.ErrInvalidInput, "cannot rollback committed reservation")
	if r.Status == ReservationCommitted {
func (r *Reservation) Rollback() error {

}
	return nil
	r.Status = ReservationCommitted
	}
		return errors.New(errors.ErrInvalidInput, "reservation already processed")
	if r.Status != ReservationPending {
	}
		return errors.New(errors.ErrInvalidInput, "reservation has expired")
	if r.IsExpired() {
func (r *Reservation) Commit() error {

}
	return time.Now().After(r.ExpiresAt)
func (r *Reservation) IsExpired() bool {

}
	}, nil
		Status:        ReservationPending,
		CreatedAt:     now,
		ExpiresAt:     now.Add(time.Duration(ttlMinutes) * time.Minute),
		Quantity:      quantity,
		ProductID:     productID,
		ReservationID: uuid.New().String(),
	return &Reservation{
	now := time.Now()

	}
		return nil, errors.New(errors.ErrInvalidInput, "quantity must be positive")
	if quantity <= 0 {
	}
		return nil, errors.New(errors.ErrInvalidInput, "product ID is required")
	if productID == "" {
func NewReservation(productID string, quantity int, ttlMinutes int) (*Reservation, error) {

}
	CreatedAt     time.Time
	ThresholdType string // LOW_STOCK, OUT_OF_STOCK
	CurrentStock  int
	ProductID     string
type StockAlert struct {

)
	ReservationExpired   ReservationStatus = "EXPIRED"
	ReservationRolledBack ReservationStatus = "ROLLED_BACK"
	ReservationCommitted ReservationStatus = "COMMITTED"
	ReservationPending   ReservationStatus = "PENDING"
const (

type ReservationStatus string

}
	Status        ReservationStatus
	CreatedAt     time.Time
	ExpiresAt     time.Time
	Quantity      int
	ProductID     string
	ReservationID string
type Reservation struct {

}
	UpdatedAt         time.Time
	WarehouseID       string
	TotalQuantity     int
	ReservedQuantity  int
	AvailableQuantity int
	ProductID         string
type Stock struct {

)
	"github.com/titan-commerce/backend/pkg/errors"
	"github.com/google/uuid"

	"time"
import (


