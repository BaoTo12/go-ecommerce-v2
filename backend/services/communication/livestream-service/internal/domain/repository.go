package domain
}
	GetAnalytics(ctx context.Context, streamID string) (*StreamAnalytics, error)
	RecordSale(ctx context.Context, streamID string, amount int64) error
	RecordShare(ctx context.Context, streamID, userID string) error
	RecordLike(ctx context.Context, streamID, userID string) error
	RecordView(ctx context.Context, streamID, userID string) error
	// Analytics operations
type AnalyticsRepository interface {
// AnalyticsRepository defines the interface for stream analytics

}
	UpdateViewerWatchTime(ctx context.Context, streamID, userID string, watchTime int) error
	GetViewerCount(ctx context.Context, streamID string) (int, error)
	GetActiveViewers(ctx context.Context, streamID string) ([]*StreamViewer, error)
	RemoveViewer(ctx context.Context, streamID, userID string) error
	AddViewer(ctx context.Context, viewer *StreamViewer) error
	// Viewer operations
type ViewerRepository interface {
// ViewerRepository defines the interface for viewer tracking

}
	PinComment(ctx context.Context, commentID string) error
	DeleteComment(ctx context.Context, commentID string) error
	GetComments(ctx context.Context, streamID string, limit int) ([]*StreamComment, error)
	CreateComment(ctx context.Context, comment *StreamComment) error
	// Comment operations
type CommentRepository interface {
// CommentRepository defines the interface for stream comments

}
	DeleteStream(ctx context.Context, streamID string) error
	GetScheduledStreams(ctx context.Context, limit int) ([]*Livestream, error)
	GetStreamsBySeller(ctx context.Context, sellerID string) ([]*Livestream, error)
	GetLiveStreams(ctx context.Context, limit int) ([]*Livestream, error)
	UpdateStream(ctx context.Context, stream *Livestream) error
	GetStream(ctx context.Context, streamID string) (*Livestream, error)
	CreateStream(ctx context.Context, stream *Livestream) error
	// Stream operations
type LivestreamRepository interface {
// LivestreamRepository defines the interface for livestream persistence

import "context"


