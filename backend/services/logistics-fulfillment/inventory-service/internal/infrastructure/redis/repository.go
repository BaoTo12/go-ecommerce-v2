package redis

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/titan-commerce/backend/inventory-service/internal/domain"
	"github.com/titan-commerce/backend/pkg/errors"
)

const (
	stockAvailablePrefix  = "stock:available:"
	stockReservedPrefix   = "stock:reserved:"
	reservationPrefix     = "reservation:"
	alertPrefix           = "alert:"
	defaultReservationTTL = 15 * time.Minute
)

type StockRepository struct {
	client *redis.Client
}

func NewStockRepository(client *redis.Client) *StockRepository {
	return &StockRepository{client: client}
}

// ReserveStock atomically reserves stock using Lua script
func (r *StockRepository) ReserveStock(ctx context.Context, productID string, quantity int, reservationID string, ttlMinutes int) (bool, error) {
	availableKey := stockAvailablePrefix + productID
	reservedKey := stockReservedPrefix + productID
	reservationKey := reservationPrefix + reservationID

	// Lua script for atomic reservation - NO OVERSELLING
	script := redis.NewScript(`
		local available_key = KEYS[1]
		local reserved_key = KEYS[2]
		local reservation_key = KEYS[3]
		local quantity = tonumber(ARGV[1])
		local reservation_data = ARGV[2]
		local ttl = tonumber(ARGV[3])

		local available = tonumber(redis.call('GET', available_key) or 0)

		if available >= quantity then
			-- Decrement available stock
			redis.call('DECRBY', available_key, quantity)
			-- Increment reserved stock
			redis.call('INCRBY', reserved_key, quantity)
			-- Store reservation details with TTL
			redis.call('SETEX', reservation_key, ttl, reservation_data)
			return 1
		else
			return 0
		end
	`)

	reservationData, _ := json.Marshal(map[string]interface{}{
		"reservation_id": reservationID,
		"product_id":     productID,
		"quantity":       quantity,
		"created_at":     time.Now().Unix(),
	})

	result, err := script.Run(ctx, r.client,
		[]string{availableKey, reservedKey, reservationKey},
		quantity, string(reservationData), ttlMinutes*60).Int64()

	if err != nil {
		return false, err
	}

	return result == 1, nil
}

// CommitReservation commits a reservation (remove from reserved, delete reservation record)
func (r *StockRepository) CommitReservation(ctx context.Context, reservationID string) error {
	reservationKey := reservationPrefix + reservationID

	// Get reservation data
	data, err := r.client.Get(ctx, reservationKey).Result()
	if err != nil {
		return err
	}

	var reservation map[string]interface{}
	if err := json.Unmarshal([]byte(data), &reservation); err != nil {
		return err
	}

	productID := reservation["product_id"].(string)
	quantity := int(reservation["quantity"].(float64))

	reservedKey := stockReservedPrefix + productID

	// Lua script for atomic commit
	script := redis.NewScript(`
		local reserved_key = KEYS[1]
		local reservation_key = KEYS[2]
		local quantity = tonumber(ARGV[1])

		-- Decrement reserved (stock already removed from available during reserve)
		redis.call('DECRBY', reserved_key, quantity)
		-- Delete reservation
		redis.call('DEL', reservation_key)
		return 1
	`)

	return script.Run(ctx, r.client, []string{reservedKey, reservationKey}, quantity).Err()
}

// RollbackReservation returns reserved stock to available pool
func (r *StockRepository) RollbackReservation(ctx context.Context, reservationID string) error {
	reservationKey := reservationPrefix + reservationID

	// Get reservation data
	data, err := r.client.Get(ctx, reservationKey).Result()
	if err != nil {
		return err
	}

	var reservation map[string]interface{}
	if err := json.Unmarshal([]byte(data), &reservation); err != nil {
		return err
	}

	productID := reservation["product_id"].(string)
	quantity := int(reservation["quantity"].(float64))

	availableKey := stockAvailablePrefix + productID
	reservedKey := stockReservedPrefix + productID

	// Lua script for atomic rollback
	script := redis.NewScript(`
		local available_key = KEYS[1]
		local reserved_key = KEYS[2]
		local reservation_key = KEYS[3]
		local quantity = tonumber(ARGV[1])

		-- Return to available stock
		redis.call('INCRBY', available_key, quantity)
		-- Decrement reserved
		redis.call('DECRBY', reserved_key, quantity)
		-- Delete reservation
		redis.call('DEL', reservation_key)
		return 1
	`)

	return script.Run(ctx, r.client, []string{availableKey, reservedKey, reservationKey}, quantity).Err()
}

// GetAvailableStock returns available stock
func (r *StockRepository) GetAvailableStock(ctx context.Context, productID string) (int, error) {
	key := stockAvailablePrefix + productID
	val, err := r.client.Get(ctx, key).Int()
	if err == redis.Nil {
		return 0, nil
	}
	return val, err
}

// GetReservedStock returns reserved stock
func (r *StockRepository) GetReservedStock(ctx context.Context, productID string) (int, error) {
	key := stockReservedPrefix + productID
	val, err := r.client.Get(ctx, key).Int()
	if err == redis.Nil {
		return 0, nil
	}
	return val, err
}

// CheckAvailability checks if sufficient stock is available
func (r *StockRepository) CheckAvailability(ctx context.Context, productID string, quantity int) (bool, error) {
	available, err := r.GetAvailableStock(ctx, productID)
	if err != nil {
		return false, err
	}
	return available >= quantity, nil
}

// AddStock adds to available stock
func (r *StockRepository) AddStock(ctx context.Context, productID string, quantity int) error {
	key := stockAvailablePrefix + productID
	return r.client.IncrBy(ctx, key, int64(quantity)).Err()
}

// RemoveStock removes from available stock (for adjustments)
func (r *StockRepository) RemoveStock(ctx context.Context, productID string, quantity int) error {
	key := stockAvailablePrefix + productID
	return r.client.DecrBy(ctx, key, int64(quantity)).Err()
}

// SetStock sets stock level
func (r *StockRepository) SetStock(ctx context.Context, productID string, quantity int) error {
	key := stockAvailablePrefix + productID
	return r.client.Set(ctx, key, quantity, 0).Err()
}

// GetReservation retrieves reservation details
func (r *StockRepository) GetReservation(ctx context.Context, reservationID string) (*domain.Reservation, error) {
	key := reservationPrefix + reservationID
	data, err := r.client.Get(ctx, key).Result()
	if err != nil {
		return nil, err
	}

	var resData map[string]interface{}
	if err := json.Unmarshal([]byte(data), &resData); err != nil {
		return nil, err
	}

	reservation := &domain.Reservation{
		ReservationID: resData["reservation_id"].(string),
		ProductID:     resData["product_id"].(string),
		Quantity:      int(resData["quantity"].(float64)),
		CreatedAt:     time.Unix(int64(resData["created_at"].(float64)), 0),
		Status:        domain.ReservationPending,
	}

	return reservation, nil
}

// ListReservations lists all reservations for a product (not commonly used in production)
func (r *StockRepository) ListReservations(ctx context.Context, productID string) ([]*domain.Reservation, error) {
	pattern := reservationPrefix + "*"
	keys, err := r.client.Keys(ctx, pattern).Result()
	if err != nil {
		return nil, err
	}

	var reservations []*domain.Reservation
	for _, key := range keys {
		reservation, err := r.GetReservation(ctx, key[len(reservationPrefix):])
		if err != nil {
			continue
		}
		if reservation.ProductID == productID {
			reservations = append(reservations, reservation)
		}
	}

	return reservations, nil
}

// CleanExpiredReservations - Redis handles this automatically with TTL
func (r *StockRepository) CleanExpiredReservations(ctx context.Context) error {
	// Redis automatically deletes expired keys, so this is a no-op
	// Could be used to track and log expired reservations if needed
	return nil
}

// SaveAlert saves a stock alert
func (r *StockRepository) SaveAlert(ctx context.Context, alert *domain.StockAlert) error {
	key := fmt.Sprintf("%s%s:%d", alertPrefix, alert.ProductID, time.Now().Unix())
	data, _ := json.Marshal(alert)
	return r.client.Set(ctx, key, data, 24*time.Hour).Err()
}

// GetAlerts retrieves alerts for a product
func (r *StockRepository) GetAlerts(ctx context.Context, productID string) ([]*domain.StockAlert, error) {
	pattern := alertPrefix + productID + ":*"
	keys, err := r.client.Keys(ctx, pattern).Result()
	if err != nil {
		return nil, err
	}

	var alerts []*domain.StockAlert
	for _, key := range keys {
		data, err := r.client.Get(ctx, key).Result()
		if err != nil {
			continue
		}

		var alert domain.StockAlert
		if err := json.Unmarshal([]byte(data), &alert); err != nil {
			continue
		}
		alerts = append(alerts, &alert)
	}

	return alerts, nil
}
