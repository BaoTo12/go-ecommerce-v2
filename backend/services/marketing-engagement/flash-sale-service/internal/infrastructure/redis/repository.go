package redis

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"strings"
	"time"

	goredis "github.com/redis/go-redis/v9"
	"github.com/titan-commerce/backend/flash-sale-service/internal/domain"
	"github.com/titan-commerce/backend/pkg/errors"
	"github.com/titan-commerce/backend/pkg/logger"
)

type FlashSaleRepository struct {
	client *goredis.Client
	logger *logger.Logger
}

func NewFlashSaleRepository(addr, password string, logger *logger.Logger) (*FlashSaleRepository, error) {
	client := goredis.NewClient(&goredis.Options{
		Addr:     addr,
		Password: password,
		DB:       2,
	})

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := client.Ping(ctx).Err(); err != nil {
		return nil, errors.Wrap(errors.ErrInternal, "failed to ping Redis", err)
	}

	logger.Info("Flash Sale Redis repository initialized")
	return &FlashSaleRepository{client: client, logger: logger}, nil
}

func (r *FlashSaleRepository) Save(ctx context.Context, sale *domain.FlashSale) error {
	return nil
}

func (r *FlashSaleRepository) FindByID(ctx context.Context, saleID string) (*domain.FlashSale, error) {
	return &domain.FlashSale{
		ID:            saleID,
		Status:        domain.FlashSaleStatusActive,
		TotalQuantity: 1000,
		SoldQuantity:  0,
		MaxPerUser:    5,
		StartTime:     time.Now().Add(-1 * time.Hour),
		EndTime:       time.Now().Add(1 * time.Hour),
	}, nil
}

func (r *FlashSaleRepository) FindActive(ctx context.Context) ([]*domain.FlashSale, error) {
	return nil, nil
}

func (r *FlashSaleRepository) FindUpcoming(ctx context.Context) ([]*domain.FlashSale, error) {
	return nil, nil
}

func (r *FlashSaleRepository) Update(ctx context.Context, sale *domain.FlashSale) error {
	return nil
}

func (r *FlashSaleRepository) SaveReservation(ctx context.Context, res *domain.FlashSaleReservation) error {
	return nil
}

func (r *FlashSaleRepository) SavePurchase(ctx context.Context, purchase *domain.FlashSalePurchase) error {
	return nil
}

func (r *FlashSaleRepository) GetUserPurchases(ctx context.Context, saleID, userID string) (int, error) {
	key := "flashsale:purchase:" + saleID + ":" + userID
	count, err := r.client.Get(ctx, key).Int()
	if err == goredis.Nil {
		return 0, nil
	}
	return count, err
}

// DecrementStock atomically decrements stock using Lua script
func (r *FlashSaleRepository) DecrementStock(ctx context.Context, saleID string, quantity int) (bool, error) {
	stockKey := "flashsale:stock:" + saleID

	script := goredis.NewScript(`
		local stock_key = KEYS[1]
		local quantity = tonumber(ARGV[1])
		
		local current_stock = tonumber(redis.call('GET', stock_key) or 0)
		if current_stock >= quantity then
			redis.call('DECRBY', stock_key, quantity)
			return 1
		else
			return 0
		end
	`)

	result, err := script.Run(ctx, r.client, []string{stockKey}, quantity).Result()
	if err != nil {
		return false, errors.Wrap(errors.ErrInternal, "failed to decrement stock", err)
	}

	return result.(int64) == 1, nil
}

// VerifyProofOfWork verifies the PoW challenge solution
func (r *FlashSaleRepository) VerifyProofOfWork(challenge, nonce string, difficulty int) bool {
	data := challenge + nonce
	hash := sha256.Sum256([]byte(data))
	hashHex := hex.EncodeToString(hash[:])
	prefix := strings.Repeat("0", difficulty)
	return strings.HasPrefix(hashHex, prefix)
}
