package application

import (
	"context"
	"time"

	"github.com/titan-commerce/backend/livestream-service/internal/domain"
	"github.com/titan-commerce/backend/pkg/logger"
)

type LivestreamService struct {
	streamRepo    domain.LivestreamRepository
	commentRepo   domain.CommentRepository
	viewerRepo    domain.ViewerRepository
	analyticsRepo domain.AnalyticsRepository
	logger        *logger.Logger
}

func NewLivestreamService(
	streamRepo domain.LivestreamRepository,
	commentRepo domain.CommentRepository,
	viewerRepo domain.ViewerRepository,
	analyticsRepo domain.AnalyticsRepository,
	logger *logger.Logger,
) *LivestreamService {
	return &LivestreamService{
		streamRepo:    streamRepo,
		commentRepo:   commentRepo,
		viewerRepo:    viewerRepo,
		analyticsRepo: analyticsRepo,
		logger:        logger,
	}
}

// CreateStream creates a new livestream (Command)
func (s *LivestreamService) CreateStream(
	ctx context.Context,
	sellerID, sellerName, title, description string,
	scheduledAt *time.Time,
) (*domain.Livestream, error) {
	stream, err := domain.NewLivestream(sellerID, sellerName, title, description, scheduledAt)
	if err != nil {
		return nil, err
	}

	if err := s.streamRepo.CreateStream(ctx, stream); err != nil {
		s.logger.Error(err, "failed to create stream")
		return nil, err
	}

	s.logger.Infof("Livestream created: id=%s, seller=%s, title=%s",
		stream.StreamID, sellerID, title)

	return stream, nil
}

// StartStream starts a livestream (Command)
func (s *LivestreamService) StartStream(ctx context.Context, streamID, playbackURL string) error {
	stream, err := s.streamRepo.GetStream(ctx, streamID)
	if err != nil {
		return err
	}

	if err := stream.Start(playbackURL); err != nil {
		return err
	}

	if err := s.streamRepo.UpdateStream(ctx, stream); err != nil {
		s.logger.Error(err, "failed to start stream")
		return err
	}

	s.logger.Infof("Livestream started: id=%s, playback=%s", streamID, playbackURL)
	return nil
}

// EndStream ends a livestream (Command)
func (s *LivestreamService) EndStream(ctx context.Context, streamID string) error {
	stream, err := s.streamRepo.GetStream(ctx, streamID)
	if err != nil {
		return err
	}

	if err := stream.End(); err != nil {
		return err
	}

	if err := s.streamRepo.UpdateStream(ctx, stream); err != nil {
		s.logger.Error(err, "failed to end stream")
		return err
	}

	s.logger.Infof("Livestream ended: id=%s, duration=%ds, peak_viewers=%d",
		streamID, stream.Duration, stream.PeakViewerCount)

	return nil
}

// JoinStream adds a viewer to the stream (Command)
func (s *LivestreamService) JoinStream(ctx context.Context, streamID, userID, userName string) error {
	viewer := domain.NewStreamViewer(streamID, userID, userName)

	if err := s.viewerRepo.AddViewer(ctx, viewer); err != nil {
		s.logger.Error(err, "failed to add viewer")
		return err
	}

	// Update viewer count
	count, _ := s.viewerRepo.GetViewerCount(ctx, streamID)
	stream, err := s.streamRepo.GetStream(ctx, streamID)
	if err == nil {
		stream.UpdateViewerCount(count)
		s.streamRepo.UpdateStream(ctx, stream)
	}

	// Record analytics
	s.analyticsRepo.RecordView(ctx, streamID, userID)

	s.logger.Infof("Viewer joined stream: stream=%s, user=%s", streamID, userID)
	return nil
}

// LeaveStream removes a viewer from the stream (Command)
func (s *LivestreamService) LeaveStream(ctx context.Context, streamID, userID string) error {
	if err := s.viewerRepo.RemoveViewer(ctx, streamID, userID); err != nil {
		s.logger.Error(err, "failed to remove viewer")
		return err
	}

	// Update viewer count
	count, _ := s.viewerRepo.GetViewerCount(ctx, streamID)
	stream, err := s.streamRepo.GetStream(ctx, streamID)
	if err == nil {
		stream.UpdateViewerCount(count)
		s.streamRepo.UpdateStream(ctx, stream)
	}

	s.logger.Infof("Viewer left stream: stream=%s, user=%s", streamID, userID)
	return nil
}

// PostComment posts a comment in the stream (Command)
func (s *LivestreamService) PostComment(
	ctx context.Context,
	streamID, userID, userName, content string,
) (*domain.StreamComment, error) {
	comment, err := domain.NewStreamComment(streamID, userID, userName, content)
	if err != nil {
		return nil, err
	}

	if err := s.commentRepo.CreateComment(ctx, comment); err != nil {
		s.logger.Error(err, "failed to create comment")
		return nil, err
	}

	s.logger.Infof("Comment posted: stream=%s, user=%s", streamID, userID)
	return comment, nil
}

// LikeStream increments the like count (Command)
func (s *LivestreamService) LikeStream(ctx context.Context, streamID, userID string) error {
	stream, err := s.streamRepo.GetStream(ctx, streamID)
	if err != nil {
		return err
	}

	stream.IncrementLikes()

	if err := s.streamRepo.UpdateStream(ctx, stream); err != nil {
		s.logger.Error(err, "failed to update stream likes")
		return err
	}

	s.analyticsRepo.RecordLike(ctx, streamID, userID)
	s.logger.Infof("Stream liked: stream=%s, user=%s, total=%d", streamID, userID, stream.LikeCount)
	return nil
}

// AddFeaturedProduct adds a product to the stream (Command)
func (s *LivestreamService) AddFeaturedProduct(ctx context.Context, streamID, productID string) error {
	stream, err := s.streamRepo.GetStream(ctx, streamID)
	if err != nil {
		return err
	}

	if err := stream.AddFeaturedProduct(productID); err != nil {
		return err
	}

	if err := s.streamRepo.UpdateStream(ctx, stream); err != nil {
		s.logger.Error(err, "failed to add featured product")
		return err
	}

	s.logger.Infof("Product featured: stream=%s, product=%s", streamID, productID)
	return nil
}

// RecordSale records a sale made during the stream (Command)
func (s *LivestreamService) RecordSale(ctx context.Context, streamID string, amount int64) error {
	stream, err := s.streamRepo.GetStream(ctx, streamID)
	if err != nil {
		return err
	}

	stream.RecordSale(amount)

	if err := s.streamRepo.UpdateStream(ctx, stream); err != nil {
		s.logger.Error(err, "failed to record sale")
		return err
	}

	s.analyticsRepo.RecordSale(ctx, streamID, amount)
	s.logger.Infof("Sale recorded: stream=%s, amount=%d", streamID, amount)
	return nil
}

// GetStream retrieves stream details (Query)
func (s *LivestreamService) GetStream(ctx context.Context, streamID string) (*domain.Livestream, error) {
	return s.streamRepo.GetStream(ctx, streamID)
}

// GetLiveStreams retrieves all live streams (Query)
func (s *LivestreamService) GetLiveStreams(ctx context.Context, limit int) ([]*domain.Livestream, error) {
	return s.streamRepo.GetLiveStreams(ctx, limit)
}

// GetStreamComments retrieves stream comments (Query)
func (s *LivestreamService) GetStreamComments(ctx context.Context, streamID string, limit int) ([]*domain.StreamComment, error) {
	return s.commentRepo.GetComments(ctx, streamID, limit)
}

// GetStreamAnalytics retrieves stream analytics (Query)
func (s *LivestreamService) GetStreamAnalytics(ctx context.Context, streamID string) (*domain.StreamAnalytics, error) {
	return s.analyticsRepo.GetAnalytics(ctx, streamID)
}

