package infrastructure

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/titan-commerce/backend/pkg/errors"
	"github.com/titan-commerce/backend/pkg/logger"
)

// This implements "The 11.11 Problem" solution
// Goal: Handle 1M concurrent users hitting "Buy" at exactly 00:00:00

type FlashSaleRedisRepository struct {
	client *redis.Client
	logger *logger.Logger
}

func NewFlashSaleRedisRepository(addr, password string, logger *logger.Logger) (*FlashSaleRedisRepository, error) {
	client := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: password,
		DB:       2, // Separate DB for flash sales
	})

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := client.Ping(ctx).Err(); err != nil {
		return nil, errors.Wrap(errors.ErrInternal, "failed to connect to Redis", err)
	}

	logger.Info("Flash Sale Redis connected")
	return &FlashSaleRedisRepository{client: client, logger: logger}, nil
}

// AtomicPurchase performs atomic stock decrement using Lua script
// This is the CRITICAL function handling 1M concurrent requests
func (r *FlashSaleRedisRepository) AtomicPurchase(ctx context.Context, flashSaleID string, userID string, quantity int) (bool, error) {
	stockKey := "flashsale:stock:" + flashSaleID
	purchaseKey := "flashsale:purchase:" + flashSaleID + ":" + userID

	// Lua script for atomic operation
	// Ensures:
	// 1. Stock check
	// 2. Stock decrement  
	// 3. User purchase limit (prevent one user buying all)
	// 4. Record purchase
	script := redis.NewScript(`
		local stock_key = KEYS[1]
		local purchase_key = KEYS[2]
		local quantity = tonumber(ARGV[1])
		local max_per_user = tonumber(ARGV[2])
		
		-- Check if user already purchased
		local user_purchased = tonumber(redis.call('GET', purchase_key) or 0)
		if user_purchased + quantity > max_per_user then
			return -1  -- Exceeded purchase limit
		end
		
		-- Check stock
		local current_stock = tonumber(redis.call('GET', stock_key) or 0)
		if current_stock >= quantity then
			redis.call('DECRBY', stock_key, quantity)
			redis.call('INCRBY', purchase_key, quantity)
			redis.call('EXPIRE', purchase_key, 3600)  -- 1 hour expiry
			return 1  -- Success
		else
			return 0  -- Out of stock
		end
	`)

	maxPerUser := 5 // Limit each user to 5 items
	result, err := script.Run(ctx, r.client, []string{stockKey, purchaseKey}, quantity, maxPerUser).Result()
	if err != nil {
		return false, errors.Wrap(errors.ErrInternal, "failed to execute atomic purchase", err)
	}

	resultCode := result.(int64)
	if resultCode == -1 {
		return false, errors.New(errors.ErrInvalidInput, "purchase limit exceeded (max 5 per user)")
	}
	if resultCode == 0 {
		return false, errors.New(errors.ErrInsufficientStock, "flash sale sold out")
	}

	r.logger.Infof("Flash sale purchase: user=%s, flashSale=%s, qty=%d", userID, flashSaleID, quantity)
	return true, nil
}

// InitializeStock sets initial stock for flash sale
func (r *FlashSaleRedisRepository) InitializeStock(ctx context.Context, flashSaleID string, stock int) error {
	stockKey := "flashsale:stock:" + flashSaleID
	return r.client.Set(ctx, stockKey, stock, 0).Err()
}

// GetRemainingStock gets current stock level
func (r *FlashSaleRedisRepository) GetRemainingStock(ctx context.Context, flashSaleID string) (int, error) {
	stockKey := "flashsale:stock:" + flashSaleID
	
	stock, err := r.client.Get(ctx, stockKey).Int()
	if err == redis.Nil {
		return 0, nil
	}
	return stock, err
}

// VerifyProofOfWork verifies the PoW challenge solution
// This prevents bots from overwhelming the system
func (r *FlashSaleRedisRepository) VerifyProofOfWork(challenge, solution string, difficulty int) bool {
	// Verify that hash(challenge + solution) has 'difficulty' leading zeros
	data := challenge + solution
	hash := sha256.Sum256([]byte(data))
	hashHex := hex.EncodeToString(hash[:])

	// Check leading zeros
	prefix := ""
	for i := 0; i < difficulty; i++ {
		prefix += "0"
	}

	return len(hashHex) >= difficulty && hashHex[:difficulty] == prefix
}

// RateLimitCheck implements token bucket algorithm
// Allows burst of requests but limits sustained rate
func (r *FlashSaleRedisRepository) RateLimitCheck(ctx context.Context, userID string) (bool, error) {
	key := "ratelimit:" + userID
	maxRequests := 10   // 10 requests
	window := 1         // per 1 second

	// Use Redis sliding window
	script := redis.NewScript(`
		local key = KEYS[1]
		local max_requests = tonumber(ARGV[1])
		local window = tonumber(ARGV[2])
		local now = tonumber(ARGV[3])
		
		-- Remove old entries outside window
		redis.call('ZREMRANGEBYSCORE', key, 0, now - window)
		
		-- Count requests in current window
		local count = redis.call('ZCARD', key)
		
		if count < max_requests then
			redis.call('ZADD', key, now, now)
			redis.call('EXPIRE', key, window)
			return 1  -- Allowed
		else
			return 0  -- Rate limited
		end
	`)

	now := time.Now().Unix()
	result, err := script.Run(ctx, r.client, []string{key}, maxRequests, window, now).Result()
	if err != nil {
		return false, err
	}

	return result.(int64) == 1, nil
}
