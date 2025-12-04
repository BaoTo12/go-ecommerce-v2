package domain
}
	return w.Status == StatusActive
func (w *Warehouse) IsActive() bool {

}
	w.UpdatedAt = time.Now()
	}
		w.Priority = priority
	if priority > 0 {
	}
		w.Capacity = capacity
	if capacity > 0 {
	}
		w.Name = name
	if name != "" {
func (w *Warehouse) Update(name string, capacity, priority int) {

}
	w.UpdatedAt = time.Now()
	w.Status = status
func (w *Warehouse) UpdateStatus(status WarehouseStatus) {

}
	}, nil
		UpdatedAt: now,
		CreatedAt: now,
		Priority:  priority,
		Capacity:  capacity,
		Status:    StatusActive,
		Address:   address,
		Code:      code,
		Name:      name,
		ID:        uuid.New().String(),
	return &Warehouse{
	now := time.Now()

	}
		return nil, errors.New(errors.ErrInvalidInput, "capacity must be positive")
	if capacity <= 0 {
	}
		return nil, errors.New(errors.ErrInvalidInput, "warehouse code is required")
	if code == "" {
	}
		return nil, errors.New(errors.ErrInvalidInput, "warehouse name is required")
	if name == "" {
func NewWarehouse(name, code string, address Address, capacity, priority int) (*Warehouse, error) {

}
	UpdatedAt time.Time
	CreatedAt time.Time
	Priority  int // Lower number = higher priority for allocation
	Capacity  int
	Status    WarehouseStatus
	Address   Address
	Code      string
	Name      string
	ID        string
type Warehouse struct {

}
	Longitude  float64
	Latitude   float64
	Country    string
	PostalCode string
	State      string
	City       string
	Street     string
type Address struct {

)
	StatusMaintenance WarehouseStatus = "MAINTENANCE"
	StatusInactive    WarehouseStatus = "INACTIVE"
	StatusActive      WarehouseStatus = "ACTIVE"
const (

type WarehouseStatus string

)
	"github.com/titan-commerce/backend/pkg/errors"
	"github.com/google/uuid"

	"time"
import (


