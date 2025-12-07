package domain_test

import (
	"testing"
	"time"

	"github.com/titan-commerce/backend/recommendation-service/internal/domain"
)

func TestCalculateSimilarity(t *testing.T) {
	tests := []struct {
		name     string
		p1       domain.Product
		p2       domain.Product
		minScore float64
		maxScore float64
	}{
		{
			name: "Same category and similar price",
			p1: domain.Product{
				ID:       "p1",
				Category: "Electronics",
				Price:    100000,
				ShopID:   "shop1",
				Rating:   4.5,
			},
			p2: domain.Product{
				ID:       "p2",
				Category: "Electronics",
				Price:    110000,
				ShopID:   "shop2",
				Rating:   4.6,
			},
			minScore: 0.6,
			maxScore: 0.8,
		},
		{
			name: "Different category but same shop",
			p1: domain.Product{
				ID:       "p1",
				Category: "Electronics",
				Price:    100000,
				ShopID:   "shop1",
				Rating:   4.5,
			},
			p2: domain.Product{
				ID:       "p2",
				Category: "Fashion",
				Price:    50000,
				ShopID:   "shop1",
				Rating:   4.5,
			},
			minScore: 0.2,
			maxScore: 0.4,
		},
		{
			name: "Completely different products",
			p1: domain.Product{
				ID:       "p1",
				Category: "Electronics",
				Price:    100000,
				ShopID:   "shop1",
				Rating:   4.5,
			},
			p2: domain.Product{
				ID:       "p2",
				Category: "Fashion",
				Price:    10000,
				ShopID:   "shop2",
				Rating:   2.0,
			},
			minScore: 0.0,
			maxScore: 0.2,
		},
		{
			name: "Identical products",
			p1: domain.Product{
				ID:       "p1",
				Category: "Electronics",
				Price:    100000,
				ShopID:   "shop1",
				Rating:   4.5,
			},
			p2: domain.Product{
				ID:       "p2",
				Category: "Electronics",
				Price:    100000,
				ShopID:   "shop1",
				Rating:   4.5,
			},
			minScore: 0.9,
			maxScore: 1.0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			score := domain.CalculateSimilarity(tt.p1, tt.p2)
			if score < tt.minScore || score > tt.maxScore {
				t.Errorf("CalculateSimilarity() = %v, want between %v and %v", score, tt.minScore, tt.maxScore)
			}
		})
	}
}

func TestScoreByPreferences(t *testing.T) {
	products := []domain.Product{
		{ID: "p1", Category: "Electronics", Price: 100000, ShopName: "TechStore", Rating: 4.8},
		{ID: "p2", Category: "Fashion", Price: 50000, ShopName: "FashionStore", Rating: 4.2},
		{ID: "p3", Category: "Electronics", Price: 200000, ShopName: "TechStore", Rating: 4.9},
	}

	prefs := domain.UserPreference{
		UserID: "u1",
		Categories: map[string]float64{
			"Electronics": 1.0,
			"Fashion":     0.3,
		},
		PriceRangeMin:   50000,
		PriceRangeMax:   150000,
		PreferredBrands: []string{"TechStore"},
	}

	scored := domain.ScoreByPreferences(products, prefs)

	// Electronics + in price range + preferred brand should be first
	if len(scored) == 0 {
		t.Fatal("Expected scored products")
	}
	if scored[0].Product.ID != "p1" {
		t.Errorf("Expected p1 to be first, got %s", scored[0].Product.ID)
	}
}

func TestFindSimilarProducts(t *testing.T) {
	target := domain.Product{
		ID:       "target",
		Category: "Electronics",
		Price:    100000,
		ShopID:   "shop1",
		Rating:   4.5,
	}

	candidates := []domain.Product{
		{ID: "p1", Category: "Electronics", Price: 110000, ShopID: "shop1", Rating: 4.6},
		{ID: "p2", Category: "Fashion", Price: 50000, ShopID: "shop2", Rating: 4.0},
		{ID: "p3", Category: "Electronics", Price: 95000, ShopID: "shop2", Rating: 4.4},
		{ID: "target", Category: "Electronics", Price: 100000, ShopID: "shop1", Rating: 4.5}, // Should be excluded
	}

	similar := domain.FindSimilarProducts(target, candidates, 2)

	if len(similar) > 2 {
		t.Errorf("Expected at most 2 similar products, got %d", len(similar))
	}

	for _, sp := range similar {
		if sp.Product.ID == "target" {
			t.Error("Target product should be excluded from similar products")
		}
	}

	// Electronics should rank higher
	if len(similar) > 0 && similar[0].Product.Category != "Electronics" {
		t.Error("Expected electronics product to rank first")
	}
}

func BenchmarkCalculateSimilarity(b *testing.B) {
	p1 := domain.Product{ID: "p1", Category: "Electronics", Price: 100000, ShopID: "shop1", Rating: 4.5}
	p2 := domain.Product{ID: "p2", Category: "Electronics", Price: 110000, ShopID: "shop2", Rating: 4.6}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		domain.CalculateSimilarity(p1, p2)
	}
}

func BenchmarkScoreByPreferences(b *testing.B) {
	products := make([]domain.Product, 100)
	for i := 0; i < 100; i++ {
		products[i] = domain.Product{
			ID:       string(rune('0' + i)),
			Category: "Electronics",
			Price:    float64(100000 + i*1000),
			Rating:   4.5,
		}
	}

	prefs := domain.UserPreference{
		UserID:        "u1",
		Categories:    map[string]float64{"Electronics": 1.0},
		PriceRangeMin: 50000,
		PriceRangeMax: 150000,
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		domain.ScoreByPreferences(products, prefs)
	}
}

func TestBrowsingEvent(t *testing.T) {
	event := &domain.BrowsingEvent{
		UserID:    "u1",
		ProductID: "p1",
		EventType: "view",
		Timestamp: time.Now(),
	}

	if event.UserID != "u1" {
		t.Errorf("Expected UserID u1, got %s", event.UserID)
	}
	if event.EventType != "view" {
		t.Errorf("Expected EventType view, got %s", event.EventType)
	}
}
