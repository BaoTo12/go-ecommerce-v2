package domain

import (
	"time"

	"github.com/google/uuid"
)

type EventType string

const (
	EventPageView     EventType = "page_view"
	EventProductView  EventType = "product_view"
	EventAddToCart    EventType = "add_to_cart"
	EventCheckout     EventType = "checkout"
	EventPurchase     EventType = "purchase"
	EventSearch       EventType = "search"
	EventClick        EventType = "click"
)

// AnalyticsEvent represents a tracked user event
type AnalyticsEvent struct {
	ID          string
	EventID     string
	UserID      string
	SessionID   string
	EventType   EventType
	Properties  map[string]interface{}
	Timestamp   time.Time
	DeviceType  string
	Platform    string
	Country     string
	City        string
}

// DashboardMetrics represents real-time dashboard data
type DashboardMetrics struct {
	// Real-time
	ActiveUsers       int
	OrdersInProgress  int
	CurrentRevenue    float64
	
	// Today's summary
	TodayOrders       int
	TodayRevenue      float64
	TodayPageViews    int
	TodayUniqueUsers  int
	
	// Comparison
	OrdersChange      float64 // % change vs yesterday
	RevenueChange     float64
	UsersChange       float64
	
	UpdatedAt         time.Time
}

type ProductAnalytics struct {
	ProductID       string
	Views           int
	UniqueViews     int
	AddToCartRate   float64
	PurchaseRate    float64
	AvgTimeOnPage   float64
	Revenue         float64
	Quantity        int
	Period          string // "daily", "weekly", "monthly"
	Date            time.Time
}

type SalesReport struct {
	ID              string
	Period          string // "daily", "weekly", "monthly"
	StartDate       time.Time
	EndDate         time.Time
	TotalOrders     int
	TotalRevenue    float64
	TotalItems      int
	AvgOrderValue   float64
	TopProducts     []ProductSales
	CategoryBreakdown []CategorySales
	HourlyDistribution []HourlySales
	GeneratedAt     time.Time
}

type ProductSales struct {
	ProductID   string
	ProductName string
	Quantity    int
	Revenue     float64
}

type CategorySales struct {
	CategoryID   string
	CategoryName string
	Orders       int
	Revenue      float64
	Percentage   float64
}

type HourlySales struct {
	Hour     int
	Orders   int
	Revenue  float64
}

type UserCohort struct {
	CohortDate      time.Time
	UsersAcquired   int
	RetentionDay1   float64
	RetentionDay7   float64
	RetentionDay30  float64
	AvgLTV          float64
}

type FunnelStep struct {
	Name        string
	Users       int
	Conversion  float64 // % from previous step
}

type ConversionFunnel struct {
	ID          string
	Name        string
	Steps       []FunnelStep
	OverallRate float64
	Period      string
	Date        time.Time
}

func NewAnalyticsEvent(userID, sessionID string, eventType EventType, properties map[string]interface{}) *AnalyticsEvent {
	return &AnalyticsEvent{
		ID:         uuid.New().String(),
		EventID:    uuid.New().String(),
		UserID:     userID,
		SessionID:  sessionID,
		EventType:  eventType,
		Properties: properties,
		Timestamp:  time.Now(),
	}
}

type Repository interface {
	SaveEvent(ctx interface{}, event *AnalyticsEvent) error
	GetEventsByUser(ctx interface{}, userID string, limit int) ([]*AnalyticsEvent, error)
	GetEventsBySession(ctx interface{}, sessionID string) ([]*AnalyticsEvent, error)
	GetEventCount(ctx interface{}, eventType EventType, start, end time.Time) (int, error)
	GetUniqueUsers(ctx interface{}, start, end time.Time) (int, error)
	GetDashboardMetrics(ctx interface{}) (*DashboardMetrics, error)
	GetProductAnalytics(ctx interface{}, productID, period string, date time.Time) (*ProductAnalytics, error)
	GenerateSalesReport(ctx interface{}, period string, start, end time.Time) (*SalesReport, error)
	GetConversionFunnel(ctx interface{}, name, period string, date time.Time) (*ConversionFunnel, error)
	GetUserCohorts(ctx interface{}, months int) ([]*UserCohort, error)
}
