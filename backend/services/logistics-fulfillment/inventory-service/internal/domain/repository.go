package domain

import "context"

type StockRepository interface {
	// Atomic operations using Redis Lua scripts
	ReserveStock(ctx context.Context, productID string, quantity int, reservationID string, ttlMinutes int) (bool, error)
	CommitReservation(ctx context.Context, reservationID string) error
	RollbackReservation(ctx context.Context, reservationID string) error

	// Query operations
	GetAvailableStock(ctx context.Context, productID string) (int, error)
	GetReservedStock(ctx context.Context, productID string) (int, error)
	CheckAvailability(ctx context.Context, productID string, quantity int) (bool, error)

	// Stock management
	AddStock(ctx context.Context, productID string, quantity int) error
	RemoveStock(ctx context.Context, productID string, quantity int) error
	SetStock(ctx context.Context, productID string, quantity int) error

	// Reservations
	GetReservation(ctx context.Context, reservationID string) (*Reservation, error)
	ListReservations(ctx context.Context, productID string) ([]*Reservation, error)
	CleanExpiredReservations(ctx context.Context) error

	// Alerts
	SaveAlert(ctx context.Context, alert *StockAlert) error
	GetAlerts(ctx context.Context, productID string) ([]*StockAlert, error)
}

