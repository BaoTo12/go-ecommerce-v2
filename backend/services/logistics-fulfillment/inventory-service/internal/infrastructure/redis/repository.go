package infrastructure

import (
	"context"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/titan-commerce/backend/pkg/errors"
)

const (
	stockKeyPrefix       = "stock:"
	reservationKeyPrefix = "reservation:"
	reservationTTL       = 15 * time.Minute
)

type RedisInventoryRepository struct {
	client *redis.Client
}

func NewRedisInventoryRepository(addr, password string) (*RedisInventoryRepository, error) {
	client := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: password,
		DB:       1, // Use different DB than cart
	})

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := client.Ping(ctx).Err(); err != nil {
		return nil, errors.Wrap(errors.ErrInternal, "failed to connect to Redis", err)
	}

	return &RedisInventoryRepository{client: client}, nil
}

// ReserveStock atomically reserves stock using Lua script
func (r *RedisInventoryRepository) ReserveStock(ctx context.Context, productID string, quantity int, reservationID string) error {
	stockKey := stockKeyPrefix + productID
	reservationKey := reservationKeyPrefix + reservationID + ":" + productID

	// Lua script for atomic stock reservation
	script := redis.NewScript(`
		local stock_key = KEYS[1]
		local reservation_key = KEYS[2]
		local quantity = tonumber(ARGV[1])
		local ttl = tonumber(ARGV[2])
		
		local current_stock = tonumber(redis.call('GET', stock_key) or 0)
		
		if current_stock >= quantity then
			redis.call('DECRBY', stock_key, quantity)
			redis.call('SET', reservation_key, quantity, 'EX', ttl)
			return 1
		else
			return 0
		end
	`)

	result, err := script.Run(ctx, r.client, []string{stockKey, reservationKey}, quantity, int(reservationTTL.Seconds())).Result()
	if err != nil {
		return errors.Wrap(errors.ErrInternal, "failed to execute reservation script", err)
	}

	if result.(int64) == 0 {
		return errors.New(errors.ErrInsufficientStock, fmt.Sprintf("insufficient stock for product %s", productID))
	}

	return nil
}

// CommitReservation commits a reservation (no-op, just delete reservation key)
func (r *RedisInventoryRepository) CommitReservation(ctx context.Context, reservationID string, productID string) error {
	reservationKey := reservationKeyPrefix + reservationID + ":" + productID
	return r.client.Del(ctx, reservationKey).Err()
}

// RollbackReservation releases reserved stock back
func (r *RedisInventoryRepository) RollbackReservation(ctx context.Context, reservationID string, productID string) error {
	stockKey := stockKeyPrefix + productID
	reservationKey := reservationKeyPrefix + reservationID + ":" + productID

	// Lua script for atomic rollback
	script := redis.NewScript(`
		local stock_key = KEYS[1]
		local reservation_key = KEYS[2]
		
		local reserved_quantity = tonumber(redis.call('GET', reservation_key) or 0)
		
		if reserved_quantity > 0 then
			redis.call('INCRBY', stock_key, reserved_quantity)
			redis.call('DEL', reservation_key)
			return 1
		end
		return 0
	`)

	_, err := script.Run(ctx, r.client, []string{stockKey, reservationKey}).Result()
	return err
}

// GetStock returns current stock level
func (r *RedisInventoryRepository) GetStock(ctx context.Context, productID string) (int, error) {
	stockKey := stockKeyPrefix + productID
	
	quantity, err := r.client.Get(ctx, stockKey).Int()
	if err == redis.Nil {
		return 0, nil
	}
	if err != nil {
		return 0, errors.Wrap(errors.ErrInternal, "failed to get stock", err)
	}

	return quantity, nil
}

// SetStock sets initial stock level (for admin/seller)
func (r *RedisInventoryRepository) SetStock(ctx context.Context, productID string, quantity int) error {
	stockKey := stockKeyPrefix + productID
	return r.client.Set(ctx, stockKey, quantity, 0).Err()
}
