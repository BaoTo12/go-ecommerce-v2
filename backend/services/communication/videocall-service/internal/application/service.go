package application

import (
	"context"

	"github.com/titan-commerce/backend/videocall-service/internal/domain"
	"github.com/titan-commerce/backend/pkg/logger"
)

type VideocallService struct {
	callRepo      domain.VideoCallRepository
	signalingRepo domain.SignalingRepository
	logger        *logger.Logger
}

func NewVideocallService(
	callRepo domain.VideoCallRepository,
	signalingRepo domain.SignalingRepository,
	logger *logger.Logger,
) *VideocallService {
	return &VideocallService{
		callRepo:      callRepo,
		signalingRepo: signalingRepo,
		logger:        logger,
	}
}

// InitiateCall initiates a new video call (Command)
func (s *VideocallService) InitiateCall(
	ctx context.Context,
	callType domain.CallType,
	callerID, callerName, calleeID, calleeName string,
) (*domain.VideoCall, error) {
	// Check if user already in a call
	activeCall, _ := s.callRepo.GetActiveCall(ctx, callerID)
	if activeCall != nil {
		s.logger.Warnf("User already in active call: %s", callerID)
		return nil, nil
	}

	call, err := domain.NewVideoCall(callType, callerID, callerName, calleeID, calleeName)
	if err != nil {
		return nil, err
	}

	if err := s.callRepo.CreateCall(ctx, call); err != nil {
		s.logger.Error(err, "failed to create call")
		return nil, err
	}

	s.logger.Infof("Call initiated: id=%s, caller=%s, callee=%s", call.CallID, callerID, calleeID)
	return call, nil
}

// RingCall marks call as ringing (Command)
func (s *VideocallService) RingCall(ctx context.Context, callID string) error {
	call, err := s.callRepo.GetCall(ctx, callID)
	if err != nil {
		return err
	}

	if err := call.Ring(); err != nil {
		return err
	}

	if err := s.callRepo.UpdateCall(ctx, call); err != nil {
		s.logger.Error(err, "failed to update call status")
		return err
	}

	s.logger.Infof("Call ringing: id=%s", callID)
	return nil
}

// AnswerCall answers an incoming call (Command)
func (s *VideocallService) AnswerCall(ctx context.Context, callID string) error {
	call, err := s.callRepo.GetCall(ctx, callID)
	if err != nil {
		return err
	}

	if err := call.Answer(); err != nil {
		return err
	}

	if err := s.callRepo.UpdateCall(ctx, call); err != nil {
		s.logger.Error(err, "failed to answer call")
		return err
	}

	s.logger.Infof("Call answered: id=%s, duration starting", callID)
	return nil
}

// EndCall ends a video call (Command)
func (s *VideocallService) EndCall(ctx context.Context, callID string) error {
	call, err := s.callRepo.GetCall(ctx, callID)
	if err != nil {
		return err
	}

	if err := call.End(); err != nil {
		return err
	}

	if err := s.callRepo.UpdateCall(ctx, call); err != nil {
		s.logger.Error(err, "failed to end call")
		return err
	}

	s.logger.Infof("Call ended: id=%s, duration=%ds, status=%s",
		callID, call.Duration, call.Status)

	return nil
}

// RejectCall rejects an incoming call (Command)
func (s *VideocallService) RejectCall(ctx context.Context, callID string) error {
	call, err := s.callRepo.GetCall(ctx, callID)
	if err != nil {
		return err
	}

	if err := call.Reject(); err != nil {
		return err
	}

	if err := s.callRepo.UpdateCall(ctx, call); err != nil {
		s.logger.Error(err, "failed to reject call")
		return err
	}

	s.logger.Infof("Call rejected: id=%s", callID)
	return nil
}

// SendSignal sends WebRTC signaling data (Command)
func (s *VideocallService) SendSignal(
	ctx context.Context,
	callID, userID, eventType, payload string,
) error {
	event := domain.NewCallEvent(callID, userID, eventType, payload)

	if err := s.signalingRepo.PublishEvent(ctx, event); err != nil {
		s.logger.Error(err, "failed to publish signal")
		return err
	}

	s.logger.Infof("Signal sent: call=%s, user=%s, type=%s", callID, userID, eventType)
	return nil
}

// UpdateCallQuality updates call quality metrics (Command)
func (s *VideocallService) UpdateCallQuality(
	ctx context.Context,
	callID string,
	quality *domain.CallQuality,
) error {
	call, err := s.callRepo.GetCall(ctx, callID)
	if err != nil {
		return err
	}

	call.UpdateQuality(quality)

	if err := s.callRepo.UpdateCall(ctx, call); err != nil {
		s.logger.Error(err, "failed to update call quality")
		return err
	}

	return nil
}

// GetCall retrieves call details (Query)
func (s *VideocallService) GetCall(ctx context.Context, callID string) (*domain.VideoCall, error) {
	return s.callRepo.GetCall(ctx, callID)
}

// GetUserCalls retrieves user's call history (Query)
func (s *VideocallService) GetUserCalls(ctx context.Context, userID string, limit int) ([]*domain.VideoCall, error) {
	return s.callRepo.GetUserCalls(ctx, userID, limit)
}

// SubscribeToSignals subscribes to call signaling events (Query)
func (s *VideocallService) SubscribeToSignals(ctx context.Context, callID string) (<-chan *domain.CallEvent, error) {
	return s.signalingRepo.SubscribeEvents(ctx, callID)
}

