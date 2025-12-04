package infrastructure

import (
	"context"
	"encoding/json"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/titan-commerce/backend/cart-service/internal/domain"
	"github.com/titan-commerce/backend/pkg/errors"
)

const (
	cartKeyPrefix = "cart:"
	cartTTL       = 7 * 24 * time.Hour // 7 days
)

type RedisCartRepository struct {
	client *redis.Client
}

func NewRedisCartRepository(addr, password string) (*RedisCartRepository, error) {
	client := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: password,
		DB:       0,
	})

	// Test connection
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	
	if err := client.Ping(ctx).Err(); err != nil {
		return nil, errors.Wrap(errors.ErrInternal, "failed to connect to Redis", err)
	}

	return &RedisCartRepository{client: client}, nil
}

func (r *RedisCartRepository) Save(ctx context.Context, cart *domain.Cart) error {
	key := cartKeyPrefix + cart.UserID
	
	data, err := json.Marshal(cart)
	if err != nil {
		return errors.Wrap(errors.ErrInternal, "failed to marshal cart", err)
	}

	if err := r.client.Set(ctx, key, data, cartTTL).Err(); err != nil {
		return errors.Wrap(errors.ErrInternal, "failed to save cart to Redis", err)
	}

	return nil
}

func (r *RedisCartRepository) FindByUserID(ctx context.Context, userID string) (*domain.Cart, error) {
	key := cartKeyPrefix + userID
	
	data, err := r.client.Get(ctx, key).Result()
	if err == redis.Nil {
		// Cart doesn't exist, return new empty cart
		return domain.NewCart(userID), nil
	}
	if err != nil {
		return nil, errors.Wrap(errors.ErrInternal, "failed to get cart from Redis", err)
	}

	var cart domain.Cart
	if err := json.Unmarshal([]byte(data), &cart); err != nil {
		return nil, errors.Wrap(errors.ErrInternal, "failed to unmarshal cart", err)
	}

	return &cart, nil
}

func (r *RedisCartRepository) Delete(ctx context.Context, userID string) error {
	key := cartKeyPrefix + userID
	
	if err := r.client.Del(ctx, key).Err(); err != nil {
		return errors.Wrap(errors.ErrInternal, "failed to delete cart from Redis", err)
	}

	return nil
}
