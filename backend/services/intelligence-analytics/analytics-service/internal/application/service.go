package application

import (
	"context"
	"time"

	"github.com/titan-commerce/backend/analytics-service/internal/domain"
	"github.com/titan-commerce/backend/pkg/logger"
)

type AnalyticsRepository interface {
	SaveEvent(ctx context.Context, event *domain.AnalyticsEvent) error
	GetEventsByUser(ctx context.Context, userID string, limit int) ([]*domain.AnalyticsEvent, error)
	GetEventsBySession(ctx context.Context, sessionID string) ([]*domain.AnalyticsEvent, error)
	GetEventCount(ctx context.Context, eventType domain.EventType, start, end time.Time) (int, error)
	GetUniqueUsers(ctx context.Context, start, end time.Time) (int, error)
	GetDashboardMetrics(ctx context.Context) (*domain.DashboardMetrics, error)
	GetProductAnalytics(ctx context.Context, productID, period string, date time.Time) (*domain.ProductAnalytics, error)
	GenerateSalesReport(ctx context.Context, period string, start, end time.Time) (*domain.SalesReport, error)
	GetConversionFunnel(ctx context.Context, name, period string, date time.Time) (*domain.ConversionFunnel, error)
	GetUserCohorts(ctx context.Context, months int) ([]*domain.UserCohort, error)
	GetTopProducts(ctx context.Context, limit int, start, end time.Time) ([]domain.ProductSales, error)
	IncrementRealTimeMetric(ctx context.Context, metric string, value float64) error
	GetRealTimeActiveUsers(ctx context.Context) (int, error)
}

type AnalyticsService struct {
	repo   AnalyticsRepository
	logger *logger.Logger
}

func NewAnalyticsService(repo AnalyticsRepository, logger *logger.Logger) *AnalyticsService {
	return &AnalyticsService{
		repo:   repo,
		logger: logger,
	}
}

// TrackEvent records an analytics event
func (s *AnalyticsService) TrackEvent(ctx context.Context, userID, sessionID string, eventType domain.EventType, properties map[string]interface{}, deviceType, platform, country, city string) error {
	event := domain.NewAnalyticsEvent(userID, sessionID, eventType, properties)
	event.DeviceType = deviceType
	event.Platform = platform
	event.Country = country
	event.City = city

	if err := s.repo.SaveEvent(ctx, event); err != nil {
		s.logger.Error(err, "failed to save analytics event")
		return err
	}

	// Update real-time metrics
	switch eventType {
	case domain.EventPurchase:
		if revenue, ok := properties["revenue"].(float64); ok {
			s.repo.IncrementRealTimeMetric(ctx, "revenue", revenue)
		}
		s.repo.IncrementRealTimeMetric(ctx, "orders", 1)
	case domain.EventPageView:
		s.repo.IncrementRealTimeMetric(ctx, "page_views", 1)
	}

	return nil
}

// GetDashboard returns real-time dashboard metrics
func (s *AnalyticsService) GetDashboard(ctx context.Context) (*domain.DashboardMetrics, error) {
	metrics, err := s.repo.GetDashboardMetrics(ctx)
	if err != nil {
		return nil, err
	}

	// Get real-time active users
	activeUsers, _ := s.repo.GetRealTimeActiveUsers(ctx)
	metrics.ActiveUsers = activeUsers
	metrics.UpdatedAt = time.Now()

	return metrics, nil
}

// GetProductAnalytics returns analytics for a specific product
func (s *AnalyticsService) GetProductAnalytics(ctx context.Context, productID string, period string) (*domain.ProductAnalytics, error) {
	return s.repo.GetProductAnalytics(ctx, productID, period, time.Now())
}

// GetSalesReport generates a sales report
func (s *AnalyticsService) GetSalesReport(ctx context.Context, period string) (*domain.SalesReport, error) {
	var start, end time.Time
	now := time.Now()

	switch period {
	case "daily":
		start = time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
		end = start.AddDate(0, 0, 1)
	case "weekly":
		start = now.AddDate(0, 0, -7)
		end = now
	case "monthly":
		start = time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, now.Location())
		end = start.AddDate(0, 1, 0)
	default:
		start = now.AddDate(0, 0, -1)
		end = now
	}

	return s.repo.GenerateSalesReport(ctx, period, start, end)
}

// GetConversionFunnel returns conversion funnel analysis
func (s *AnalyticsService) GetConversionFunnel(ctx context.Context, funnelName, period string) (*domain.ConversionFunnel, error) {
	return s.repo.GetConversionFunnel(ctx, funnelName, period, time.Now())
}

// GetUserCohorts returns user cohort retention analysis
func (s *AnalyticsService) GetUserCohorts(ctx context.Context, months int) ([]*domain.UserCohort, error) {
	return s.repo.GetUserCohorts(ctx, months)
}

// GetTopProducts returns top selling products
func (s *AnalyticsService) GetTopProducts(ctx context.Context, limit int, daysBack int) ([]domain.ProductSales, error) {
	end := time.Now()
	start := end.AddDate(0, 0, -daysBack)
	return s.repo.GetTopProducts(ctx, limit, start, end)
}

// GetUserJourney returns the event sequence for a user
func (s *AnalyticsService) GetUserJourney(ctx context.Context, userID string, limit int) ([]*domain.AnalyticsEvent, error) {
	return s.repo.GetEventsByUser(ctx, userID, limit)
}

// GetSessionReplay returns events for a session (for session replay)
func (s *AnalyticsService) GetSessionReplay(ctx context.Context, sessionID string) ([]*domain.AnalyticsEvent, error) {
	return s.repo.GetEventsBySession(ctx, sessionID)
}

// CalculateMetrics calculates key metrics for a time range
func (s *AnalyticsService) CalculateMetrics(ctx context.Context, start, end time.Time) (map[string]interface{}, error) {
	metrics := make(map[string]interface{})

	// Page views
	pageViews, _ := s.repo.GetEventCount(ctx, domain.EventPageView, start, end)
	metrics["page_views"] = pageViews

	// Product views
	productViews, _ := s.repo.GetEventCount(ctx, domain.EventProductView, start, end)
	metrics["product_views"] = productViews

	// Add to cart
	addToCart, _ := s.repo.GetEventCount(ctx, domain.EventAddToCart, start, end)
	metrics["add_to_cart"] = addToCart

	// Purchases
	purchases, _ := s.repo.GetEventCount(ctx, domain.EventPurchase, start, end)
	metrics["purchases"] = purchases

	// Unique users
	uniqueUsers, _ := s.repo.GetUniqueUsers(ctx, start, end)
	metrics["unique_users"] = uniqueUsers

	// Conversion rates
	if productViews > 0 {
		metrics["view_to_cart_rate"] = float64(addToCart) / float64(productViews) * 100
	}
	if addToCart > 0 {
		metrics["cart_to_purchase_rate"] = float64(purchases) / float64(addToCart) * 100
	}

	return metrics, nil
}

// TrackPageView convenience method for tracking page views
func (s *AnalyticsService) TrackPageView(ctx context.Context, userID, sessionID, pageURL, pageTitle string) error {
	return s.TrackEvent(ctx, userID, sessionID, domain.EventPageView, map[string]interface{}{
		"url":   pageURL,
		"title": pageTitle,
	}, "", "", "", "")
}

// TrackProductView convenience method for tracking product views
func (s *AnalyticsService) TrackProductView(ctx context.Context, userID, sessionID, productID, productName string, price float64) error {
	return s.TrackEvent(ctx, userID, sessionID, domain.EventProductView, map[string]interface{}{
		"product_id":   productID,
		"product_name": productName,
		"price":        price,
	}, "", "", "", "")
}

// TrackPurchase convenience method for tracking purchases
func (s *AnalyticsService) TrackPurchase(ctx context.Context, userID, sessionID, orderID string, revenue float64, items int) error {
	return s.TrackEvent(ctx, userID, sessionID, domain.EventPurchase, map[string]interface{}{
		"order_id": orderID,
		"revenue":  revenue,
		"items":    items,
	}, "", "", "", "")
}
