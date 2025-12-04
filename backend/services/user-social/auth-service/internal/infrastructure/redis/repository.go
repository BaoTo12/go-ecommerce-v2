package redis

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/titan-commerce/backend/pkg/errors"
)

type RedisRepository struct {
	client *redis.Client
}

func NewRedisRepository(addr, password string) (*RedisRepository, error) {
	client := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: password,
		DB:       0,
	})

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := client.Ping(ctx).Err(); err != nil {
		return nil, errors.Wrap(errors.ErrInternal, "failed to connect to redis", err)
	}

	return &RedisRepository{client: client}, nil
}

func (r *RedisRepository) BlacklistToken(ctx context.Context, token string, expiration time.Duration) error {
	return r.client.Set(ctx, "blacklist:"+token, "true", expiration).Err()
}

func (r *RedisRepository) IsBlacklisted(ctx context.Context, token string) (bool, error) {
	exists, err := r.client.Exists(ctx, "blacklist:"+token).Result()
	if err != nil {
		return false, err
	}
	return exists > 0, nil
}

func (r *RedisRepository) StoreRefreshToken(ctx context.Context, userID, token string, expiration time.Duration) error {
	return r.client.Set(ctx, "refresh:"+userID, token, expiration).Err()
}

func (r *RedisRepository) GetRefreshToken(ctx context.Context, userID string) (string, error) {
	return r.client.Get(ctx, "refresh:"+userID).Result()
}

func (r *RedisRepository) RevokeRefreshToken(ctx context.Context, userID string) error {
	return r.client.Del(ctx, "refresh:"+userID).Err()
}
