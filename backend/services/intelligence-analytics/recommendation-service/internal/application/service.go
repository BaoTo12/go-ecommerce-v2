package application

import (
	"context"
	"time"

	"github.com/titan-commerce/backend/recommendation-service/internal/domain"
	"github.com/titan-commerce/backend/pkg/logger"
)

// RecommendationService handles product recommendations
type RecommendationService struct {
	repo   domain.Repository
	logger *logger.Logger
}

// NewRecommendationService creates a new recommendation service
func NewRecommendationService(repo domain.Repository, log *logger.Logger) *RecommendationService {
	return &RecommendationService{
		repo:   repo,
		logger: log,
	}
}

// GetAlsoBought returns "customers also bought" recommendations
func (s *RecommendationService) GetAlsoBought(ctx context.Context, productID string) (*domain.RecommendationResult, error) {
	product, err := s.repo.GetProductByID(ctx, productID)
	if err != nil || product == nil {
		return &domain.RecommendationResult{
			Products:   []domain.Product{},
			Reason:     "",
			Confidence: 0,
		}, nil
	}

	// Get all products and find similar ones
	allProducts, err := s.repo.GetProducts(ctx, 100)
	if err != nil {
		return nil, err
	}

	scored := domain.FindSimilarProducts(*product, allProducts, 6)
	products := make([]domain.Product, len(scored))
	for i, sp := range scored {
		products[i] = sp.Product
	}

	confidence := 0.0
	if len(scored) > 0 {
		confidence = scored[0].Score
	}

	return &domain.RecommendationResult{
		Products:   products,
		Reason:     "Khách hàng cũng mua",
		Confidence: confidence,
	}, nil
}

// GetPersonalized returns personalized recommendations based on user history
func (s *RecommendationService) GetPersonalized(ctx context.Context, userID string) (*domain.RecommendationResult, error) {
	prefs, err := s.repo.GetUserPreferences(ctx, userID)
	if err != nil || prefs == nil {
		// Return popular products if no preferences
		products, _ := s.repo.GetProducts(ctx, 12)
		return &domain.RecommendationResult{
			Products:   products,
			Reason:     "Gợi ý cho bạn",
			Confidence: 0.5,
		}, nil
	}

	allProducts, err := s.repo.GetProducts(ctx, 50)
	if err != nil {
		return nil, err
	}

	scored := domain.ScoreByPreferences(allProducts, *prefs)
	products := make([]domain.Product, 0, 12)
	for i, sp := range scored {
		if i >= 12 {
			break
		}
		products = append(products, sp.Product)
	}

	confidence := 0.0
	if len(scored) > 0 {
		confidence = scored[0].Score
	}

	return &domain.RecommendationResult{
		Products:   products,
		Reason:     "Dành riêng cho bạn",
		Confidence: confidence,
	}, nil
}

// GetFrequentlyBoughtTogether returns complementary products
func (s *RecommendationService) GetFrequentlyBoughtTogether(ctx context.Context, productID string) (*domain.RecommendationResult, error) {
	product, err := s.repo.GetProductByID(ctx, productID)
	if err != nil || product == nil {
		return &domain.RecommendationResult{
			Products:   []domain.Product{},
			Reason:     "",
			Confidence: 0,
		}, nil
	}

	// Get products from different categories
	allProducts, _ := s.repo.GetProducts(ctx, 30)
	complementary := make([]domain.Product, 0, 3)
	for _, p := range allProducts {
		if p.ID != productID && p.Category != product.Category {
			complementary = append(complementary, p)
			if len(complementary) >= 3 {
				break
			}
		}
	}

	return &domain.RecommendationResult{
		Products:   complementary,
		Reason:     "Thường được mua cùng",
		Confidence: 0.7,
	}, nil
}

// GetBecauseYouViewed returns similar products to a viewed product
func (s *RecommendationService) GetBecauseYouViewed(ctx context.Context, productID string) (*domain.RecommendationResult, error) {
	product, err := s.repo.GetProductByID(ctx, productID)
	if err != nil || product == nil {
		return &domain.RecommendationResult{
			Products:   []domain.Product{},
			Reason:     "",
			Confidence: 0,
		}, nil
	}

	similar, _ := s.repo.GetProductsByCategory(ctx, product.Category, 8)
	products := make([]domain.Product, 0, 6)
	for _, p := range similar {
		if p.ID != productID {
			products = append(products, p)
			if len(products) >= 6 {
				break
			}
		}
	}

	reason := "Vì bạn đã xem"
	if len(product.Name) > 30 {
		reason = reason + " \"" + product.Name[:30] + "...\""
	} else {
		reason = reason + " \"" + product.Name + "\""
	}

	return &domain.RecommendationResult{
		Products:   products,
		Reason:     reason,
		Confidence: 0.8,
	}, nil
}

// GetTrendingNearYou returns trending products in a location
func (s *RecommendationService) GetTrendingNearYou(ctx context.Context, location string) (*domain.RecommendationResult, error) {
	if location == "" {
		location = "TP. Hồ Chí Minh"
	}

	products, err := s.repo.GetTrendingByLocation(ctx, location, 6)
	if err != nil {
		return nil, err
	}

	return &domain.RecommendationResult{
		Products:   products,
		Reason:     "Đang hot tại " + location,
		Confidence: 0.75,
	}, nil
}

// RecordView records a product view for preference learning
func (s *RecommendationService) RecordView(ctx context.Context, userID, productID string) error {
	event := &domain.BrowsingEvent{
		UserID:    userID,
		ProductID: productID,
		EventType: "view",
		Timestamp: time.Now(),
	}
	return s.repo.SaveBrowsingEvent(ctx, event)
}
