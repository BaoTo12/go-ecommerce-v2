package application

import (
	"context"
	"sync"
	"time"

	"github.com/titan-commerce/backend/livestream-service/internal/domain"
	"github.com/titan-commerce/backend/pkg/logger"
)

type LivestreamRepository interface {
	Save(ctx context.Context, stream *domain.Livestream) error
	FindByID(ctx context.Context, streamID string) (*domain.Livestream, error)
	FindByStreamKey(ctx context.Context, streamKey string) (*domain.Livestream, error)
	FindLiveStreams(ctx context.Context, limit, offset int) ([]*domain.Livestream, error)
	Update(ctx context.Context, stream *domain.Livestream) error
	SaveChat(ctx context.Context, chat *domain.StreamChat) error
	GetRecentChats(ctx context.Context, streamID string, limit int) ([]*domain.StreamChat, error)
}

// ViewerManager tracks real-time viewers per stream
type ViewerManager struct {
	viewers map[string]map[string]bool // streamID -> userID -> connected
	mu      sync.RWMutex
}

func NewViewerManager() *ViewerManager {
	return &ViewerManager{
		viewers: make(map[string]map[string]bool),
	}
}

func (vm *ViewerManager) AddViewer(streamID, userID string) int {
	vm.mu.Lock()
	defer vm.mu.Unlock()
	if vm.viewers[streamID] == nil {
		vm.viewers[streamID] = make(map[string]bool)
	}
	vm.viewers[streamID][userID] = true
	return len(vm.viewers[streamID])
}

func (vm *ViewerManager) RemoveViewer(streamID, userID string) int {
	vm.mu.Lock()
	defer vm.mu.Unlock()
	if vm.viewers[streamID] != nil {
		delete(vm.viewers[streamID], userID)
	}
	return len(vm.viewers[streamID])
}

func (vm *ViewerManager) GetViewerCount(streamID string) int {
	vm.mu.RLock()
	defer vm.mu.RUnlock()
	return len(vm.viewers[streamID])
}

type LivestreamService struct {
	repo      LivestreamRepository
	viewerMgr *ViewerManager
	logger    *logger.Logger
}

func NewLivestreamService(repo LivestreamRepository, logger *logger.Logger) *LivestreamService {
	return &LivestreamService{
		repo:      repo,
		viewerMgr: NewViewerManager(),
		logger:    logger,
	}
}

// CreateStream creates a new livestream session
func (s *LivestreamService) CreateStream(ctx context.Context, sellerID, title, description string) (*domain.Livestream, error) {
	stream := domain.NewLivestream(sellerID, title, description)
	
	if err := s.repo.Save(ctx, stream); err != nil {
		s.logger.Error(err, "failed to create stream")
		return nil, err
	}

	s.logger.Infof("Created stream: %s for seller: %s", stream.ID, sellerID)
	return stream, nil
}

// StartStream marks stream as live (called by RTMP ingest callback)
func (s *LivestreamService) StartStream(ctx context.Context, streamKey string) (*domain.Livestream, error) {
	stream, err := s.repo.FindByStreamKey(ctx, streamKey)
	if err != nil {
		return nil, err
	}

	// Generate HLS URL (in production, this comes from transcoding service)
	hlsUrl := "https://cdn.titancommerce.io/live/" + stream.ID + "/index.m3u8"
	stream.GoLive(hlsUrl)

	if err := s.repo.Update(ctx, stream); err != nil {
		return nil, err
	}

	s.logger.Infof("Stream started: %s", stream.ID)
	return stream, nil
}

// EndStream ends a livestream
func (s *LivestreamService) EndStream(ctx context.Context, streamID string) error {
	stream, err := s.repo.FindByID(ctx, streamID)
	if err != nil {
		return err
	}

	stream.End()
	if err := s.repo.Update(ctx, stream); err != nil {
		return err
	}

	s.logger.Infof("Stream ended: %s, peak viewers: %d", streamID, stream.PeakViewerCount)
	return nil
}

// JoinStream handles a user joining a stream
func (s *LivestreamService) JoinStream(ctx context.Context, streamID, userID string) (*domain.Livestream, error) {
	stream, err := s.repo.FindByID(ctx, streamID)
	if err != nil {
		return nil, err
	}

	count := s.viewerMgr.AddViewer(streamID, userID)
	stream.UpdateViewerCount(count)
	s.repo.Update(ctx, stream)

	return stream, nil
}

// LeaveStream handles a user leaving a stream
func (s *LivestreamService) LeaveStream(ctx context.Context, streamID, userID string) {
	count := s.viewerMgr.RemoveViewer(streamID, userID)
	
	stream, err := s.repo.FindByID(ctx, streamID)
	if err == nil {
		stream.UpdateViewerCount(count)
		s.repo.Update(ctx, stream)
	}
}

// SendChat sends a chat message in the stream
func (s *LivestreamService) SendChat(ctx context.Context, streamID, userID, username, message string) (*domain.StreamChat, error) {
	chat := domain.NewStreamChat(streamID, userID, username, message)
	if err := s.repo.SaveChat(ctx, chat); err != nil {
		return nil, err
	}
	return chat, nil
}

// GetRecentChats gets recent chat messages
func (s *LivestreamService) GetRecentChats(ctx context.Context, streamID string, limit int) ([]*domain.StreamChat, error) {
	return s.repo.GetRecentChats(ctx, streamID, limit)
}

// PinProduct pins a product during livestream
func (s *LivestreamService) PinProduct(ctx context.Context, streamID, productID string) error {
	stream, err := s.repo.FindByID(ctx, streamID)
	if err != nil {
		return err
	}

	stream.PinProduct(productID)
	return s.repo.Update(ctx, stream)
}

// GetLiveStreams returns currently live streams
func (s *LivestreamService) GetLiveStreams(ctx context.Context, limit, offset int) ([]*domain.Livestream, error) {
	return s.repo.FindLiveStreams(ctx, limit, offset)
}

// GetStream returns a stream by ID
func (s *LivestreamService) GetStream(ctx context.Context, streamID string) (*domain.Livestream, error) {
	return s.repo.FindByID(ctx, streamID)
}

// TriggerFlashSale triggers a flash sale during livestream
func (s *LivestreamService) TriggerFlashSale(ctx context.Context, streamID, productID string, discountPercent int, duration time.Duration) error {
	stream, err := s.repo.FindByID(ctx, streamID)
	if err != nil {
		return err
	}

	// In production, this would publish an event to Kafka for the flash-sale-service
	s.logger.Infof("Flash sale triggered: stream=%s, product=%s, discount=%d%%, duration=%v",
		stream.ID, productID, discountPercent, duration)
	
	return nil
}
