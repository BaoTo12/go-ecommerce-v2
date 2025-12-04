package domain
}
	}
		Timestamp:  time.Now(),
		Dimensions: dimensions,
		Value:      value,
		Type:       metricType,
		Name:       name,
		MetricID:   uuid.New().String(),
	return &Metric{
func NewMetric(name, metricType string, value float64, dimensions map[string]string) *Metric {
// NewMetric creates a new metric

}
	}
		Timestamp:  time.Now(),
		Properties: properties,
		SessionID:  sessionID,
		UserID:     userID,
		EventType:  eventType,
		EventID:    uuid.New().String(),
	return &Event{
func NewEvent(eventType, userID, sessionID string, properties map[string]interface{}) *Event {
// NewEvent creates a new analytics event

}
	Data        interface{}
	GeneratedAt time.Time
	EndDate     time.Time
	StartDate   time.Time
	Filters     map[string]string
	Metrics     []string
	Type        string // sales, user_behavior, product_performance
	Description string
	Name        string
	ReportID    string
type Report struct {
// Report represents an analytics report

}
	Timestamp  time.Time
	Properties map[string]interface{}
	SessionID  string
	UserID     string
	EventType  string // page_view, add_to_cart, purchase, etc.
	EventID    string
type Event struct {
// Event represents a business event

}
	Timestamp  time.Time
	Dimensions map[string]string
	Value      float64
	Type       string // counter, gauge, histogram
	Name       string
	MetricID   string
type Metric struct {
// Metric represents an analytics metric

)
	"github.com/google/uuid"

	"time"
import (


