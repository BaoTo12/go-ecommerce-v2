package application

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/titan-commerce/backend/flash-sale-service/internal/domain"
	"github.com/titan-commerce/backend/pkg/errors"
	"github.com/titan-commerce/backend/pkg/logger"
)

type FlashSaleRepository interface {
	Save(ctx context.Context, sale *domain.FlashSale) error
	FindByID(ctx context.Context, saleID string) (*domain.FlashSale, error)
	FindActive(ctx context.Context) ([]*domain.FlashSale, error)
	FindUpcoming(ctx context.Context) ([]*domain.FlashSale, error)
	Update(ctx context.Context, sale *domain.FlashSale) error
	SaveReservation(ctx context.Context, res *domain.FlashSaleReservation) error
	SavePurchase(ctx context.Context, purchase *domain.FlashSalePurchase) error
	GetUserPurchases(ctx context.Context, saleID, userID string) (int, error)
	DecrementStock(ctx context.Context, saleID string, quantity int) (bool, error)
}

// RateLimiter for token bucket rate limiting
type RateLimiter struct {
	buckets map[string]*tokenBucket
	mu      sync.RWMutex
}

type tokenBucket struct {
	tokens     float64
	maxTokens  float64
	refillRate float64 // tokens per second
	lastRefill time.Time
}

func NewRateLimiter() *RateLimiter {
	return &RateLimiter{
		buckets: make(map[string]*tokenBucket),
	}
}

func (rl *RateLimiter) Allow(userID string, tokensNeeded float64) bool {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	bucket, ok := rl.buckets[userID]
	if !ok {
		bucket = &tokenBucket{
			tokens:     10,
			maxTokens:  10,
			refillRate: 1, // 1 token per second
			lastRefill: time.Now(),
		}
		rl.buckets[userID] = bucket
	}

	// Refill tokens
	now := time.Now()
	elapsed := now.Sub(bucket.lastRefill).Seconds()
	bucket.tokens += elapsed * bucket.refillRate
	if bucket.tokens > bucket.maxTokens {
		bucket.tokens = bucket.maxTokens
	}
	bucket.lastRefill = now

	if bucket.tokens >= tokensNeeded {
		bucket.tokens -= tokensNeeded
		return true
	}
	return false
}

// ProofOfWork validator
type PoWValidator struct {
	difficulty int // Number of leading zeros required
}

func NewPoWValidator(difficulty int) *PoWValidator {
	return &PoWValidator{difficulty: difficulty}
}

func (pow *PoWValidator) GenerateChallenge(saleID, userID string) string {
	return fmt.Sprintf("%s:%s:%d", saleID, userID, time.Now().UnixNano())
}

func (pow *PoWValidator) ValidateProof(challenge, nonce string) bool {
	data := challenge + nonce
	hash := sha256.Sum256([]byte(data))
	hashHex := hex.EncodeToString(hash[:])
	
	// Check leading zeros
	prefix := strings.Repeat("0", pow.difficulty)
	return strings.HasPrefix(hashHex, prefix)
}

type FlashSaleService struct {
	repo        FlashSaleRepository
	rateLimiter *RateLimiter
	powValidator *PoWValidator
	logger      *logger.Logger
	reserveTTL  time.Duration
}

func NewFlashSaleService(repo FlashSaleRepository, logger *logger.Logger) *FlashSaleService {
	return &FlashSaleService{
		repo:         repo,
		rateLimiter:  NewRateLimiter(),
		powValidator: NewPoWValidator(4), // 4 leading zeros
		logger:       logger,
		reserveTTL:   5 * time.Minute,
	}
}

// CreateFlashSale creates a new flash sale
func (s *FlashSaleService) CreateFlashSale(ctx context.Context, productID string, originalPrice, salePrice float64, totalQty, maxPerUser int, start, end time.Time) (*domain.FlashSale, error) {
	sale := domain.NewFlashSale(productID, originalPrice, salePrice, totalQty, maxPerUser, start, end)
	
	if err := s.repo.Save(ctx, sale); err != nil {
		return nil, err
	}

	s.logger.Infof("Flash sale created: %s for product %s, %d%% off", 
		sale.ID, productID, sale.DiscountPercent)
	return sale, nil
}

// GetChallenge returns a PoW challenge for the user
func (s *FlashSaleService) GetChallenge(saleID, userID string) string {
	return s.powValidator.GenerateChallenge(saleID, userID)
}

// AttemptPurchase attempts to purchase from flash sale with PoW verification
func (s *FlashSaleService) AttemptPurchase(ctx context.Context, saleID, userID string, quantity int, challenge, nonce string) (*domain.FlashSaleReservation, error) {
	// Step 1: Rate limiting
	if !s.rateLimiter.Allow(userID, 1) {
		return nil, errors.New(errors.ErrInvalidInput, "rate limit exceeded, please wait")
	}

	// Step 2: Validate Proof of Work
	if !s.powValidator.ValidateProof(challenge, nonce) {
		return nil, errors.New(errors.ErrInvalidInput, "invalid proof of work")
	}

	// Step 3: Get flash sale
	sale, err := s.repo.FindByID(ctx, saleID)
	if err != nil {
		return nil, err
	}

	// Step 4: Check if active
	if !sale.IsActive() {
		return nil, errors.New(errors.ErrInvalidInput, "flash sale not active")
	}

	// Step 5: Check user limit
	purchased, _ := s.repo.GetUserPurchases(ctx, saleID, userID)
	if purchased+quantity > sale.MaxPerUser {
		return nil, errors.New(errors.ErrInvalidInput, "max per user limit exceeded")
	}

	// Step 6: Atomic stock decrement (Redis Lua script in production)
	success, err := s.repo.DecrementStock(ctx, saleID, quantity)
	if err != nil {
		return nil, err
	}
	if !success {
		return nil, errors.New(errors.ErrInvalidInput, "sold out")
	}

	// Step 7: Create reservation (user has X minutes to complete payment)
	reservation := domain.NewReservation(saleID, userID, quantity, s.reserveTTL)
	if err := s.repo.SaveReservation(ctx, reservation); err != nil {
		// Rollback stock
		return nil, err
	}

	s.logger.Infof("Flash sale reservation: %s by user %s for %d items", 
		reservation.ID, userID, quantity)
	return reservation, nil
}

// ConfirmPurchase confirms a reservation after payment
func (s *FlashSaleService) ConfirmPurchase(ctx context.Context, reservationID, userID string) error {
	// In production: verify reservation, create purchase record
	purchase := &domain.FlashSalePurchase{
		ID:          reservationID,
		UserID:      userID,
		PurchasedAt: time.Now(),
	}
	return s.repo.SavePurchase(ctx, purchase)
}

// GetActiveFlashSales returns currently active flash sales
func (s *FlashSaleService) GetActiveFlashSales(ctx context.Context) ([]*domain.FlashSale, error) {
	return s.repo.FindActive(ctx)
}

// GetUpcomingFlashSales returns scheduled flash sales
func (s *FlashSaleService) GetUpcomingFlashSales(ctx context.Context) ([]*domain.FlashSale, error) {
	return s.repo.FindUpcoming(ctx)
}

// GetFlashSale returns a flash sale by ID
func (s *FlashSaleService) GetFlashSale(ctx context.Context, saleID string) (*domain.FlashSale, error) {
	return s.repo.FindByID(ctx, saleID)
}

// ActivateFlashSale activates a scheduled flash sale
func (s *FlashSaleService) ActivateFlashSale(ctx context.Context, saleID string) error {
	sale, err := s.repo.FindByID(ctx, saleID)
	if err != nil {
		return err
	}

	sale.Activate()
	return s.repo.Update(ctx, sale)
}

// EndFlashSale ends a flash sale
func (s *FlashSaleService) EndFlashSale(ctx context.Context, saleID string) error {
	sale, err := s.repo.FindByID(ctx, saleID)
	if err != nil {
		return err
	}

	sale.End()
	
	s.logger.Infof("Flash sale ended: %s, sold %d/%d", 
		saleID, sale.SoldQuantity, sale.TotalQuantity)
	return s.repo.Update(ctx, sale)
}
