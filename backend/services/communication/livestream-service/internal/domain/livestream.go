package domain

import (
	"time"

	"github.com/google/uuid"
	"github.com/titan-commerce/backend/pkg/errors"
)

// StreamStatus represents the status of a livestream
type StreamStatus string

const (
	StreamScheduled StreamStatus = "SCHEDULED"
	StreamLive      StreamStatus = "LIVE"
	StreamPaused    StreamStatus = "PAUSED"
	StreamEnded     StreamStatus = "ENDED"
	StreamCancelled StreamStatus = "CANCELLED"
)

// Livestream represents a live shopping stream
type Livestream struct {
	StreamID          string
	SellerID          string
	SellerName        string
	Title             string
	Description       string
	ThumbnailURL      string
	Status            StreamStatus
	StreamKey         string // RTMP stream key
	PlaybackURL       string // HLS/DASH playback URL
	ViewerCount       int
	PeakViewerCount   int
	LikeCount         int
	ShareCount        int
	TotalRevenue      int64 // in cents
	FeaturedProducts  []string // Product IDs
	ScheduledStartAt  *time.Time
	ActualStartAt     *time.Time
	EndedAt           *time.Time
	Duration          int // in seconds
	IsRecorded        bool
	RecordingURL      string
	CreatedAt         time.Time
	UpdatedAt         time.Time
}

// NewLivestream creates a new livestream
func NewLivestream(sellerID, sellerName, title, description string, scheduledAt *time.Time) (*Livestream, error) {
	if sellerID == "" {
		return nil, errors.New(errors.ErrInvalidInput, "seller ID is required")
	}
	if title == "" {
		return nil, errors.New(errors.ErrInvalidInput, "title is required")
	}

	now := time.Now()
	streamID := uuid.New().String()

	return &Livestream{
		StreamID:         streamID,
		SellerID:         sellerID,
		SellerName:       sellerName,
		Title:            title,
		Description:      description,
		Status:           StreamScheduled,
		StreamKey:        generateStreamKey(streamID),
		ViewerCount:      0,
		PeakViewerCount:  0,
		LikeCount:        0,
		ShareCount:       0,
		TotalRevenue:     0,
		FeaturedProducts: []string{},
		ScheduledStartAt: scheduledAt,
		IsRecorded:       true,
		CreatedAt:        now,
		UpdatedAt:        now,
	}, nil
}

// Start starts the livestream
func (l *Livestream) Start(playbackURL string) error {
	if l.Status == StreamLive {
		return errors.New(errors.ErrInvalidInput, "stream is already live")
	}
	if l.Status == StreamEnded || l.Status == StreamCancelled {
		return errors.New(errors.ErrInvalidInput, "cannot start ended or cancelled stream")
	}

	now := time.Now()
	l.Status = StreamLive
	l.ActualStartAt = &now
	l.PlaybackURL = playbackURL
	l.UpdatedAt = now
	return nil
}

// End ends the livestream
func (l *Livestream) End() error {
	if l.Status != StreamLive && l.Status != StreamPaused {
		return errors.New(errors.ErrInvalidInput, "can only end live or paused streams")
	}

	now := time.Now()
	l.Status = StreamEnded
	l.EndedAt = &now
	l.UpdatedAt = now

	if l.ActualStartAt != nil {
		l.Duration = int(now.Sub(*l.ActualStartAt).Seconds())
	}

	return nil
}

// Pause pauses the livestream
func (l *Livestream) Pause() error {
	if l.Status != StreamLive {
		return errors.New(errors.ErrInvalidInput, "can only pause live streams")
	}

	l.Status = StreamPaused
	l.UpdatedAt = time.Now()
	return nil
}

// Resume resumes the paused livestream
func (l *Livestream) Resume() error {
	if l.Status != StreamPaused {
		return errors.New(errors.ErrInvalidInput, "can only resume paused streams")
	}

	l.Status = StreamLive
	l.UpdatedAt = time.Now()
	return nil
}

// UpdateViewerCount updates the current viewer count
func (l *Livestream) UpdateViewerCount(count int) {
	l.ViewerCount = count
	if count > l.PeakViewerCount {
		l.PeakViewerCount = count
	}
	l.UpdatedAt = time.Now()
}

// IncrementLikes increments the like count
func (l *Livestream) IncrementLikes() {
	l.LikeCount++
	l.UpdatedAt = time.Now()
}

// IncrementShares increments the share count
func (l *Livestream) IncrementShares() {
	l.ShareCount++
	l.UpdatedAt = time.Now()
}

// AddFeaturedProduct adds a product to featured list
func (l *Livestream) AddFeaturedProduct(productID string) error {
	if productID == "" {
		return errors.New(errors.ErrInvalidInput, "product ID is required")
	}

	// Check if already featured
	for _, id := range l.FeaturedProducts {
		if id == productID {
			return nil
		}
	}

	l.FeaturedProducts = append(l.FeaturedProducts, productID)
	l.UpdatedAt = time.Now()
	return nil
}

// RecordSale records a sale made during the stream
func (l *Livestream) RecordSale(amount int64) {
	l.TotalRevenue += amount
	l.UpdatedAt = time.Now()
}

// SetRecordingURL sets the recording URL after stream ends
func (l *Livestream) SetRecordingURL(url string) error {
	if l.Status != StreamEnded {
		return errors.New(errors.ErrInvalidInput, "can only set recording URL for ended streams")
	}

	l.RecordingURL = url
	l.UpdatedAt = time.Now()
	return nil
}

// StreamComment represents a comment in the livestream chat
type StreamComment struct {
	CommentID string
	StreamID  string
	UserID    string
	UserName  string
	Content   string
	IsPinned  bool
	IsDeleted bool
	CreatedAt time.Time
}

// NewStreamComment creates a new stream comment
func NewStreamComment(streamID, userID, userName, content string) (*StreamComment, error) {
	if streamID == "" {
		return nil, errors.New(errors.ErrInvalidInput, "stream ID is required")
	}
	if userID == "" {
		return nil, errors.New(errors.ErrInvalidInput, "user ID is required")
	}
	if content == "" {
		return nil, errors.New(errors.ErrInvalidInput, "comment content is required")
	}

	return &StreamComment{
		CommentID: uuid.New().String(),
		StreamID:  streamID,
		UserID:    userID,
		UserName:  userName,
		Content:   content,
		IsPinned:  false,
		IsDeleted: false,
		CreatedAt: time.Now(),
	}, nil
}

// Pin pins the comment
func (c *StreamComment) Pin() {
	c.IsPinned = true
}

// Unpin unpins the comment
func (c *StreamComment) Unpin() {
	c.IsPinned = false
}

// Delete soft deletes the comment
func (c *StreamComment) Delete() {
	c.IsDeleted = true
	c.Content = "[Comment deleted]"
}

// StreamViewer represents a viewer watching the stream
type StreamViewer struct {
	StreamID   string
	UserID     string
	UserName   string
	JoinedAt   time.Time
	LeftAt     *time.Time
	WatchTime  int // in seconds
	IsActive   bool
}

// NewStreamViewer creates a new stream viewer
func NewStreamViewer(streamID, userID, userName string) *StreamViewer {
	return &StreamViewer{
		StreamID: streamID,
		UserID:   userID,
		UserName: userName,
		JoinedAt: time.Now(),
		WatchTime: 0,
		IsActive: true,
	}
}

// Leave marks the viewer as left
func (v *StreamViewer) Leave() {
	if !v.IsActive {
		return
	}

	now := time.Now()
	v.LeftAt = &now
	v.IsActive = false
	v.WatchTime = int(now.Sub(v.JoinedAt).Seconds())
}

// UpdateWatchTime updates the current watch time
func (v *StreamViewer) UpdateWatchTime() {
	if v.IsActive {
		v.WatchTime = int(time.Since(v.JoinedAt).Seconds())
	}
}

// StreamAnalytics represents analytics for a stream
type StreamAnalytics struct {
	StreamID              string
	TotalViews            int
	UniqueViewers         int
	AverageWatchTime      int
	PeakViewers           int
	TotalLikes            int
	TotalShares           int
	TotalComments         int
	TotalRevenue          int64
	ProductClickThrough   int
	ConversionRate        float64
	ViewerRetentionRate   float64
	CreatedAt             time.Time
	UpdatedAt             time.Time
}

// generateStreamKey generates a unique stream key for RTMP
func generateStreamKey(streamID string) string {
	return "stream_" + streamID + "_" + uuid.New().String()[:8]
}

