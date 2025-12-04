package domain

import "context"

type AnalyticsRepository interface {
	SaveEvent(ctx context.Context, event *Event) error
	SaveMetric(ctx context.Context, metric *Metric) error
	QueryEvents(ctx context.Context, eventType string, startTime, endTime time.Time) ([]*Event, error)
	QueryMetrics(ctx context.Context, name string, startTime, endTime time.Time) ([]*Metric, error)
	GenerateReport(ctx context.Context, report *Report) error
}

