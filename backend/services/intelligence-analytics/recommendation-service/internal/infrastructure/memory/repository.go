package memory

import (
	"context"
	"sort"
	"sync"
	"time"

	"github.com/titan-commerce/backend/recommendation-service/internal/domain"
)

// RecommendationRepository is an in-memory implementation for development
type RecommendationRepository struct {
	products    []domain.Product
	preferences map[string]*domain.UserPreference
	events      []domain.BrowsingEvent
	mu          sync.RWMutex
}

// NewRecommendationRepository creates a new in-memory repository with sample data
func NewRecommendationRepository() *RecommendationRepository {
	repo := &RecommendationRepository{
		preferences: make(map[string]*domain.UserPreference),
		events:      make([]domain.BrowsingEvent, 0),
	}
	repo.seedProducts()
	return repo
}

func (r *RecommendationRepository) seedProducts() {
	r.products = []domain.Product{
		{ID: "p1", Name: "iPhone 15 Pro Max 256GB", Price: 29990000, OriginalPrice: 34990000, Discount: 14, Category: "Điện thoại", CategoryID: "phones", Thumbnail: "https://images.unsplash.com/photo-1695048133142-1a20484d2569?w=300", Rating: 4.9, Reviews: 8560, Sold: 12300, SoldDisplay: "12.3k", Location: "TP. Hồ Chí Minh", ShopID: "shop1", ShopName: "Apple Store"},
		{ID: "p2", Name: "Samsung Galaxy S24 Ultra 512GB", Price: 25990000, OriginalPrice: 29990000, Discount: 13, Category: "Điện thoại", CategoryID: "phones", Thumbnail: "https://images.unsplash.com/photo-1610945415295-d9bbf067e59c?w=300", Rating: 4.8, Reviews: 5430, Sold: 8700, SoldDisplay: "8.7k", Location: "Hà Nội", ShopID: "shop2", ShopName: "Samsung Store"},
		{ID: "p3", Name: "MacBook Air M3 13 inch", Price: 27990000, OriginalPrice: 31990000, Discount: 12, Category: "Laptop", CategoryID: "laptops", Thumbnail: "https://images.unsplash.com/photo-1517336714731-489689fd1ca8?w=300", Rating: 4.9, Reviews: 2340, Sold: 3200, SoldDisplay: "3.2k", Location: "TP. Hồ Chí Minh", ShopID: "shop1", ShopName: "Apple Store"},
		{ID: "p4", Name: "Áo Hoodie Unisex Form Rộng", Price: 199000, OriginalPrice: 350000, Discount: 43, Category: "Thời trang", CategoryID: "fashion", Thumbnail: "https://images.unsplash.com/photo-1556821840-3a63f95609a7?w=300", Rating: 4.7, Reviews: 12340, Sold: 45200, SoldDisplay: "45.2k", Location: "Hà Nội", ShopID: "shop3", ShopName: "Fashion Store"},
		{ID: "p5", Name: "Nike Air Force 1 07 Low White", Price: 2590000, OriginalPrice: 3200000, Discount: 19, Category: "Giày dép", CategoryID: "shoes", Thumbnail: "https://images.unsplash.com/photo-1600269452121-4f2416e55c28?w=300", Rating: 4.8, Reviews: 3456, Sold: 5200, SoldDisplay: "5.2k", Location: "TP. Hồ Chí Minh", ShopID: "shop4", ShopName: "Nike Store"},
		{ID: "p6", Name: "Son Dior Addict Lip Glow", Price: 950000, OriginalPrice: 1200000, Discount: 21, Category: "Làm đẹp", CategoryID: "beauty", Thumbnail: "https://images.unsplash.com/photo-1586495777744-4413f21062fa?w=300", Rating: 4.9, Reviews: 8765, Sold: 18700, SoldDisplay: "18.7k", Location: "TP. Hồ Chí Minh", ShopID: "shop5", ShopName: "Dior Beauty"},
		{ID: "p7", Name: "Nồi Chiên Không Dầu Lock&Lock 5.2L", Price: 1290000, OriginalPrice: 2490000, Discount: 48, Category: "Nhà cửa", CategoryID: "home", Thumbnail: "https://images.unsplash.com/photo-1585515320310-259814833e62?w=300", Rating: 4.8, Reviews: 4567, Sold: 23400, SoldDisplay: "23.4k", Location: "Hà Nội", ShopID: "shop6", ShopName: "Lock&Lock"},
		{ID: "p8", Name: "Laptop Dell XPS 13 Plus", Price: 32990000, OriginalPrice: 38990000, Discount: 15, Category: "Laptop", CategoryID: "laptops", Thumbnail: "https://images.unsplash.com/photo-1593642702821-c8da6771f0c6?w=300", Rating: 4.7, Reviews: 1234, Sold: 1200, SoldDisplay: "1.2k", Location: "TP. Hồ Chí Minh", ShopID: "shop7", ShopName: "Dell Store"},
		{ID: "p9", Name: "Quần Jean Nam Slim Fit", Price: 299000, OriginalPrice: 450000, Discount: 34, Category: "Thời trang", CategoryID: "fashion", Thumbnail: "https://images.unsplash.com/photo-1542272604-787c3835535d?w=300", Rating: 4.6, Reviews: 6780, Sold: 67800, SoldDisplay: "67.8k", Location: "Hà Nội", ShopID: "shop3", ShopName: "Fashion Store"},
		{ID: "p10", Name: "Serum Vitamin C The Ordinary", Price: 350000, OriginalPrice: 500000, Discount: 30, Category: "Làm đẹp", CategoryID: "beauty", Thumbnail: "https://images.unsplash.com/photo-1620916566398-39f1143ab7be?w=300", Rating: 4.8, Reviews: 9876, Sold: 34500, SoldDisplay: "34.5k", Location: "TP. Hồ Chí Minh", ShopID: "shop8", ShopName: "Beauty Zone"},
	}
}

func (r *RecommendationRepository) GetProducts(ctx interface{}, limit int) ([]domain.Product, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	if limit <= 0 || limit > len(r.products) {
		return r.products, nil
	}
	return r.products[:limit], nil
}

func (r *RecommendationRepository) GetProductByID(ctx interface{}, productID string) (*domain.Product, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	for _, p := range r.products {
		if p.ID == productID {
			return &p, nil
		}
	}
	return nil, nil
}

func (r *RecommendationRepository) GetProductsByCategory(ctx interface{}, category string, limit int) ([]domain.Product, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	result := make([]domain.Product, 0)
	for _, p := range r.products {
		if p.Category == category {
			result = append(result, p)
			if len(result) >= limit {
				break
			}
		}
	}
	return result, nil
}

func (r *RecommendationRepository) GetUserPreferences(ctx interface{}, userID string) (*domain.UserPreference, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	if pref, ok := r.preferences[userID]; ok {
		return pref, nil
	}
	return nil, nil
}

func (r *RecommendationRepository) SaveUserPreferences(ctx interface{}, prefs *domain.UserPreference) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	prefs.UpdatedAt = time.Now()
	r.preferences[prefs.UserID] = prefs
	return nil
}

func (r *RecommendationRepository) GetBrowsingHistory(ctx interface{}, userID string, limit int) ([]domain.BrowsingEvent, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	result := make([]domain.BrowsingEvent, 0)
	for i := len(r.events) - 1; i >= 0; i-- {
		if r.events[i].UserID == userID {
			result = append(result, r.events[i])
			if len(result) >= limit {
				break
			}
		}
	}
	return result, nil
}

func (r *RecommendationRepository) SaveBrowsingEvent(ctx interface{}, event *domain.BrowsingEvent) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.events = append(r.events, *event)
	return nil
}

func (r *RecommendationRepository) GetTrendingByLocation(ctx interface{}, location string, limit int) ([]domain.Product, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	local := make([]domain.Product, 0)
	for _, p := range r.products {
		if p.Location == location {
			local = append(local, p)
		}
	}

	// Sort by sold count descending
	sort.Slice(local, func(i, j int) bool {
		return local[i].Sold > local[j].Sold
	})

	if len(local) > limit {
		local = local[:limit]
	}
	return local, nil
}
