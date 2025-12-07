package domain

import (
	"math"
	"sort"
	"time"
)

// Product represents a product for recommendations
type Product struct {
	ID            string    `json:"id"`
	Name          string    `json:"name"`
	Description   string    `json:"description"`
	Price         float64   `json:"price"`
	OriginalPrice float64   `json:"original_price"`
	Discount      int       `json:"discount"`
	Category      string    `json:"category"`
	CategoryID    string    `json:"category_id"`
	Thumbnail     string    `json:"thumbnail"`
	Rating        float64   `json:"rating"`
	Reviews       int       `json:"reviews"`
	Sold          int       `json:"sold"`
	SoldDisplay   string    `json:"sold_display"`
	Stock         int       `json:"stock"`
	Location      string    `json:"location"`
	ShopID        string    `json:"shop_id"`
	ShopName      string    `json:"shop_name"`
	CreatedAt     time.Time `json:"created_at"`
}

// RecommendationResult holds recommendation data
type RecommendationResult struct {
	Products   []Product `json:"products"`
	Reason     string    `json:"reason"`
	Confidence float64   `json:"confidence"`
}

// UserPreference represents learned user preferences
type UserPreference struct {
	UserID          string             `json:"user_id"`
	Categories      map[string]float64 `json:"categories"`
	PriceRangeMin   float64            `json:"price_range_min"`
	PriceRangeMax   float64            `json:"price_range_max"`
	PreferredBrands []string           `json:"preferred_brands"`
	ViewedProducts  []string           `json:"viewed_products"`
	UpdatedAt       time.Time          `json:"updated_at"`
}

// BrowsingEvent represents a user browsing action
type BrowsingEvent struct {
	UserID    string    `json:"user_id"`
	ProductID string    `json:"product_id"`
	EventType string    `json:"event_type"` // "view", "cart", "purchase", "wishlist"
	Timestamp time.Time `json:"timestamp"`
}

// ScoredProduct holds a product with its recommendation score
type ScoredProduct struct {
	Product Product
	Score   float64
}

// CalculateSimilarity computes similarity between two products
func CalculateSimilarity(p1, p2 Product) float64 {
	score := 0.0

	// Same category (40% weight)
	if p1.Category == p2.Category {
		score += 0.4
	}

	// Similar price range within 30% (30% weight)
	maxPrice := math.Max(p1.Price, p2.Price)
	if maxPrice > 0 {
		priceDiff := math.Abs(p1.Price-p2.Price) / maxPrice
		if priceDiff < 0.3 {
			score += 0.3 * (1 - priceDiff)
		}
	}

	// Same shop (20% weight)
	if p1.ShopID == p2.ShopID {
		score += 0.2
	}

	// Similar rating (10% weight)
	ratingDiff := math.Abs(p1.Rating - p2.Rating)
	if ratingDiff < 1 {
		score += 0.1 * (1 - ratingDiff)
	}

	return score
}

// ScoreByPreferences scores products based on user preferences
func ScoreByPreferences(products []Product, prefs UserPreference) []ScoredProduct {
	scored := make([]ScoredProduct, 0, len(products))

	for _, p := range products {
		score := 0.0

		// Category preference (40% weight)
		if catScore, ok := prefs.Categories[p.Category]; ok {
			score += catScore * 0.4
		}

		// Price range fit (30% weight)
		if p.Price >= prefs.PriceRangeMin && p.Price <= prefs.PriceRangeMax {
			score += 0.3
		}

		// Brand preference (20% weight)
		for _, brand := range prefs.PreferredBrands {
			if p.ShopName == brand {
				score += 0.2
				break
			}
		}

		// High rating bonus (10% weight)
		if p.Rating >= 4.5 {
			score += 0.1
		}

		scored = append(scored, ScoredProduct{Product: p, Score: score})
	}

	// Sort by score descending
	sort.Slice(scored, func(i, j int) bool {
		return scored[i].Score > scored[j].Score
	})

	return scored
}

// FindSimilarProducts finds products similar to the given product
func FindSimilarProducts(target Product, candidates []Product, limit int) []ScoredProduct {
	scored := make([]ScoredProduct, 0, len(candidates))

	for _, p := range candidates {
		if p.ID == target.ID {
			continue
		}
		score := CalculateSimilarity(target, p)
		if score > 0 {
			scored = append(scored, ScoredProduct{Product: p, Score: score})
		}
	}

	sort.Slice(scored, func(i, j int) bool {
		return scored[i].Score > scored[j].Score
	})

	if len(scored) > limit {
		scored = scored[:limit]
	}

	return scored
}

// Repository interface for recommendation data
type Repository interface {
	GetProducts(ctx interface{}, limit int) ([]Product, error)
	GetProductByID(ctx interface{}, productID string) (*Product, error)
	GetProductsByCategory(ctx interface{}, category string, limit int) ([]Product, error)
	GetUserPreferences(ctx interface{}, userID string) (*UserPreference, error)
	SaveUserPreferences(ctx interface{}, prefs *UserPreference) error
	GetBrowsingHistory(ctx interface{}, userID string, limit int) ([]BrowsingEvent, error)
	SaveBrowsingEvent(ctx interface{}, event *BrowsingEvent) error
	GetTrendingByLocation(ctx interface{}, location string, limit int) ([]Product, error)
}
