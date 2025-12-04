package domain

import (
	"time"

	"github.com/google/uuid"
)

type StreamStatus string

const (
	StreamStatusPending  StreamStatus = "PENDING"
	StreamStatusLive     StreamStatus = "LIVE"
	StreamStatusEnded    StreamStatus = "ENDED"
	StreamStatusBanned   StreamStatus = "BANNED"
)

type Livestream struct {
	ID              string
	SellerID        string
	Title           string
	Description     string
	ThumbnailURL    string
	StreamKey       string       // RTMP stream key
	RTMPUrl         string       // RTMP ingest URL
	HLSUrl          string       // HLS playback URL
	Status          StreamStatus
	ViewerCount     int
	PeakViewerCount int
	PinnedProducts  []string     // Product IDs
	StartedAt       *time.Time
	EndedAt         *time.Time
	CreatedAt       time.Time
	UpdatedAt       time.Time
}

type StreamChat struct {
	ID          string
	StreamID    string
	UserID      string
	Username    string
	Message     string
	IsPinned    bool
	CreatedAt   time.Time
}

type StreamEvent struct {
	ID        string
	StreamID  string
	Type      string // "flash_sale", "product_pin", "giveaway"
	Payload   map[string]interface{}
	CreatedAt time.Time
}

func NewLivestream(sellerID, title, description string) *Livestream {
	streamKey := uuid.New().String()
	return &Livestream{
		ID:             uuid.New().String(),
		SellerID:       sellerID,
		Title:          title,
		Description:    description,
		StreamKey:      streamKey,
		RTMPUrl:        "rtmp://live.titancommerce.io/live/" + streamKey,
		Status:         StreamStatusPending,
		ViewerCount:    0,
		PinnedProducts: []string{},
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
	}
}

func (l *Livestream) GoLive(hlsUrl string) {
	now := time.Now()
	l.Status = StreamStatusLive
	l.HLSUrl = hlsUrl
	l.StartedAt = &now
	l.UpdatedAt = now
}

func (l *Livestream) End() {
	now := time.Now()
	l.Status = StreamStatusEnded
	l.EndedAt = &now
	l.UpdatedAt = now
}

func (l *Livestream) UpdateViewerCount(count int) {
	l.ViewerCount = count
	if count > l.PeakViewerCount {
		l.PeakViewerCount = count
	}
	l.UpdatedAt = time.Now()
}

func (l *Livestream) PinProduct(productID string) {
	// Add to front
	l.PinnedProducts = append([]string{productID}, l.PinnedProducts...)
	l.UpdatedAt = time.Now()
}

func (l *Livestream) UnpinProduct(productID string) {
	var newPinned []string
	for _, p := range l.PinnedProducts {
		if p != productID {
			newPinned = append(newPinned, p)
		}
	}
	l.PinnedProducts = newPinned
	l.UpdatedAt = time.Now()
}

func NewStreamChat(streamID, userID, username, message string) *StreamChat {
	return &StreamChat{
		ID:        uuid.New().String(),
		StreamID:  streamID,
		UserID:    userID,
		Username:  username,
		Message:   message,
		CreatedAt: time.Now(),
	}
}

type Repository interface {
	Save(ctx interface{}, stream *Livestream) error
	FindByID(ctx interface{}, streamID string) (*Livestream, error)
	FindByStreamKey(ctx interface{}, streamKey string) (*Livestream, error)
	FindLiveStreams(ctx interface{}, limit, offset int) ([]*Livestream, error)
	Update(ctx interface{}, stream *Livestream) error
	SaveChat(ctx interface{}, chat *StreamChat) error
	GetRecentChats(ctx interface{}, streamID string, limit int) ([]*StreamChat, error)
}
