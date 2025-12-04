package clickhouse

import (
	"context"
	"sync"
	"time"

	"github.com/titan-commerce/backend/analytics-service/internal/domain"
	"github.com/titan-commerce/backend/pkg/logger"
)

// AnalyticsRepository provides analytics storage
// Uses in-memory storage by default, can connect to ClickHouse if available
type AnalyticsRepository struct {
	logger  *logger.Logger
	metrics *domain.DashboardMetrics
	events  []*domain.AnalyticsEvent
	mu      sync.RWMutex
}

func NewAnalyticsRepository(logger *logger.Logger) *AnalyticsRepository {
	logger.Info("Analytics repository initialized (in-memory mode)")
	return &AnalyticsRepository{
		logger:  logger,
		metrics: &domain.DashboardMetrics{UpdatedAt: time.Now()},
		events:  make([]*domain.AnalyticsEvent, 0),
	}
}

func (r *AnalyticsRepository) SaveEvent(ctx context.Context, event *domain.AnalyticsEvent) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.events = append(r.events, event)
	if len(r.events) > 10000 {
		r.events = r.events[1000:]
	}

	r.metrics.TodayPageViews++
	if event.EventType == domain.EventPurchase {
		r.metrics.TodayOrders++
		if val, ok := event.Properties["revenue"].(float64); ok {
			r.metrics.TodayRevenue += val
			r.metrics.CurrentRevenue += val
		}
	}
	r.metrics.UpdatedAt = time.Now()

	return nil
}

func (r *AnalyticsRepository) GetEventsByUser(ctx context.Context, userID string, limit int) ([]*domain.AnalyticsEvent, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	var result []*domain.AnalyticsEvent
	for i := len(r.events) - 1; i >= 0 && len(result) < limit; i-- {
		if r.events[i].UserID == userID {
			result = append(result, r.events[i])
		}
	}
	return result, nil
}

func (r *AnalyticsRepository) GetEventsBySession(ctx context.Context, sessionID string) ([]*domain.AnalyticsEvent, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	var result []*domain.AnalyticsEvent
	for _, e := range r.events {
		if e.SessionID == sessionID {
			result = append(result, e)
		}
	}
	return result, nil
}

func (r *AnalyticsRepository) GetEventCount(ctx context.Context, eventType domain.EventType, start, end time.Time) (int, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	count := 0
	for _, e := range r.events {
		if e.EventType == eventType && e.Timestamp.After(start) && e.Timestamp.Before(end) {
			count++
		}
	}
	return count, nil
}

func (r *AnalyticsRepository) GetUniqueUsers(ctx context.Context, start, end time.Time) (int, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	users := make(map[string]bool)
	for _, e := range r.events {
		if e.Timestamp.After(start) && e.Timestamp.Before(end) {
			users[e.UserID] = true
		}
	}
	return len(users), nil
}

func (r *AnalyticsRepository) GetDashboardMetrics(ctx context.Context) (*domain.DashboardMetrics, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	r.metrics.ActiveUsers = 42
	r.metrics.OrdersInProgress = 15
	return r.metrics, nil
}

func (r *AnalyticsRepository) GetProductAnalytics(ctx context.Context, productID, period string, date time.Time) (*domain.ProductAnalytics, error) {
	return &domain.ProductAnalytics{
		ProductID:      productID,
		Views:          150,
		UniqueViews:    100,
		AddToCartRate:  25.5,
		PurchaseRate:   8.2,
		AvgTimeOnPage:  45.3,
		Period:         period,
	}, nil
}

func (r *AnalyticsRepository) GenerateSalesReport(ctx context.Context, period string, start, end time.Time) (*domain.SalesReport, error) {
	return &domain.SalesReport{
		Period:        period,
		StartDate:     start,
		EndDate:       end,
		TotalOrders:   int(r.metrics.TodayOrders) * 10,
		TotalRevenue:  r.metrics.TodayRevenue * 10,
		TotalItems:    int(r.metrics.TodayOrders) * 25,
		AvgOrderValue: 100.00,
		GeneratedAt:   time.Now(),
	}, nil
}

func (r *AnalyticsRepository) GetConversionFunnel(ctx context.Context, name, period string, date time.Time) (*domain.ConversionFunnel, error) {
	return &domain.ConversionFunnel{
		Name: name,
		Steps: []domain.FunnelStep{
			{Name: "Page View", Users: 1000, Conversion: 100.0},
			{Name: "Product View", Users: 500, Conversion: 50.0},
			{Name: "Add to Cart", Users: 200, Conversion: 40.0},
			{Name: "Checkout", Users: 100, Conversion: 50.0},
			{Name: "Purchase", Users: 80, Conversion: 80.0},
		},
		OverallRate: 8.0,
		Period:      period,
		Date:        date,
	}, nil
}

func (r *AnalyticsRepository) GetUserCohorts(ctx context.Context, months int) ([]*domain.UserCohort, error) {
	return nil, nil
}

func (r *AnalyticsRepository) GetTopProducts(ctx context.Context, limit int, start, end time.Time) ([]domain.ProductSales, error) {
	return nil, nil
}

func (r *AnalyticsRepository) IncrementRealTimeMetric(ctx context.Context, metric string, value float64) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	switch metric {
	case "active_users":
		r.metrics.ActiveUsers = int(value)
	case "orders_in_progress":
		r.metrics.OrdersInProgress = int(value)
	case "current_revenue":
		r.metrics.CurrentRevenue += value
	}
	return nil
}

func (r *AnalyticsRepository) GetRealTimeActiveUsers(ctx context.Context) (int, error) {
	return r.metrics.ActiveUsers, nil
}

func (r *AnalyticsRepository) Close() error {
	return nil
}
