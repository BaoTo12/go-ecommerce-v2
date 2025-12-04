package application

import (
	"context"
	"sync"
	"time"

	"github.com/titan-commerce/backend/videocall-service/internal/domain"
	"github.com/titan-commerce/backend/pkg/logger"
)

type VideoCallRepository interface {
	Save(ctx context.Context, call *domain.VideoCall) error
	FindByID(ctx context.Context, callID string) (*domain.VideoCall, error)
	FindByRoomID(ctx context.Context, roomID string) (*domain.VideoCall, error)
	Update(ctx context.Context, call *domain.VideoCall) error
	FindCallHistory(ctx context.Context, userID string, limit int) ([]*domain.VideoCall, error)
}

// SignalingHub manages WebRTC signaling between peers
type SignalingHub struct {
	channels map[string]chan *domain.SignalingMessage // userID -> channel
	mu       sync.RWMutex
}

func NewSignalingHub() *SignalingHub {
	return &SignalingHub{
		channels: make(map[string]chan *domain.SignalingMessage),
	}
}

func (h *SignalingHub) Register(userID string) chan *domain.SignalingMessage {
	h.mu.Lock()
	defer h.mu.Unlock()
	ch := make(chan *domain.SignalingMessage, 10)
	h.channels[userID] = ch
	return ch
}

func (h *SignalingHub) Unregister(userID string) {
	h.mu.Lock()
	defer h.mu.Unlock()
	if ch, ok := h.channels[userID]; ok {
		close(ch)
		delete(h.channels, userID)
	}
}

func (h *SignalingHub) SendToUser(userID string, msg *domain.SignalingMessage) bool {
	h.mu.RLock()
	defer h.mu.RUnlock()
	if ch, ok := h.channels[userID]; ok {
		select {
		case ch <- msg:
			return true
		default:
			return false
		}
	}
	return false
}

type VideoCallService struct {
	repo       VideoCallRepository
	sigHub     *SignalingHub
	logger     *logger.Logger
	ringTTL    time.Duration
}

func NewVideoCallService(repo VideoCallRepository, logger *logger.Logger) *VideoCallService {
	return &VideoCallService{
		repo:    repo,
		sigHub:  NewSignalingHub(),
		logger:  logger,
		ringTTL: 60 * time.Second, // 60 seconds to answer
	}
}

func (s *VideoCallService) GetSignalingHub() *SignalingHub {
	return s.sigHub
}

// InitiateCall starts a call to another user
func (s *VideoCallService) InitiateCall(ctx context.Context, callerID, calleeID string, callType domain.CallType) (*domain.VideoCall, error) {
	call := domain.NewVideoCall(callerID, calleeID, callType)
	
	if err := s.repo.Save(ctx, call); err != nil {
		s.logger.Error(err, "failed to save call")
		return nil, err
	}

	// Send incoming call signal
	s.sigHub.SendToUser(calleeID, &domain.SignalingMessage{
		Type:      "incoming_call",
		CallID:    call.ID,
		FromUser:  callerID,
		ToUser:    calleeID,
		Payload:   map[string]interface{}{"room_id": call.RoomID, "call_type": string(callType)},
		CreatedAt: time.Now(),
	})

	// Set timeout for missed call
	go s.handleCallTimeout(call.ID)

	s.logger.Infof("Call initiated: %s from %s to %s", call.ID, callerID, calleeID)
	return call, nil
}

func (s *VideoCallService) handleCallTimeout(callID string) {
	time.Sleep(s.ringTTL)
	
	call, err := s.repo.FindByID(context.Background(), callID)
	if err != nil {
		return
	}
	
	if call.Status == domain.CallStatusRinging {
		call.Miss()
		s.repo.Update(context.Background(), call)
		s.logger.Infof("Call missed: %s", callID)
	}
}

// AcceptCall accepts an incoming call
func (s *VideoCallService) AcceptCall(ctx context.Context, callID, userID string) (*domain.VideoCall, error) {
	call, err := s.repo.FindByID(ctx, callID)
	if err != nil {
		return nil, err
	}

	call.Accept()
	if err := s.repo.Update(ctx, call); err != nil {
		return nil, err
	}

	// Notify caller
	s.sigHub.SendToUser(call.CallerID, &domain.SignalingMessage{
		Type:      "call_accepted",
		CallID:    call.ID,
		FromUser:  userID,
		ToUser:    call.CallerID,
		Payload:   map[string]interface{}{"room_id": call.RoomID},
		CreatedAt: time.Now(),
	})

	s.logger.Infof("Call accepted: %s", callID)
	return call, nil
}

// RejectCall rejects an incoming call
func (s *VideoCallService) RejectCall(ctx context.Context, callID, userID string) error {
	call, err := s.repo.FindByID(ctx, callID)
	if err != nil {
		return err
	}

	call.Reject()
	if err := s.repo.Update(ctx, call); err != nil {
		return err
	}

	// Notify caller
	s.sigHub.SendToUser(call.CallerID, &domain.SignalingMessage{
		Type:      "call_rejected",
		CallID:    call.ID,
		FromUser:  userID,
		ToUser:    call.CallerID,
		CreatedAt: time.Now(),
	})

	s.logger.Infof("Call rejected: %s", callID)
	return nil
}

// EndCall ends an active call
func (s *VideoCallService) EndCall(ctx context.Context, callID, userID string) error {
	call, err := s.repo.FindByID(ctx, callID)
	if err != nil {
		return err
	}

	call.End()
	if err := s.repo.Update(ctx, call); err != nil {
		return err
	}

	// Notify other party
	otherUser := call.CallerID
	if userID == call.CallerID {
		otherUser = call.CalleeID
	}

	s.sigHub.SendToUser(otherUser, &domain.SignalingMessage{
		Type:      "call_ended",
		CallID:    call.ID,
		FromUser:  userID,
		ToUser:    otherUser,
		CreatedAt: time.Now(),
	})

	s.logger.Infof("Call ended: %s, duration: %ds", callID, call.Duration)
	return nil
}

// RelaySignaling relays WebRTC signaling messages (SDP, ICE candidates)
func (s *VideoCallService) RelaySignaling(ctx context.Context, msg *domain.SignalingMessage) error {
	msg.CreatedAt = time.Now()
	
	if !s.sigHub.SendToUser(msg.ToUser, msg) {
		s.logger.Warnf("Failed to relay signaling to %s", msg.ToUser)
	}

	return nil
}

// GetCallHistory gets call history for a user
func (s *VideoCallService) GetCallHistory(ctx context.Context, userID string, limit int) ([]*domain.VideoCall, error) {
	return s.repo.FindCallHistory(ctx, userID, limit)
}

// GetCall gets a call by ID
func (s *VideoCallService) GetCall(ctx context.Context, callID string) (*domain.VideoCall, error) {
	return s.repo.FindByID(ctx, callID)
}
