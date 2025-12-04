package domain

import "context"

// VideoCallRepository defines the interface for video call persistence
type VideoCallRepository interface {
	CreateCall(ctx context.Context, call *VideoCall) error
	GetCall(ctx context.Context, callID string) (*VideoCall, error)
	UpdateCall(ctx context.Context, call *VideoCall) error
	GetUserCalls(ctx context.Context, userID string, limit int) ([]*VideoCall, error)
	GetActiveCall(ctx context.Context, userID string) (*VideoCall, error)
}

// SignalingRepository defines the interface for WebRTC signaling
type SignalingRepository interface {
	PublishEvent(ctx context.Context, event *CallEvent) error
	SubscribeEvents(ctx context.Context, callID string) (<-chan *CallEvent, error)
	GetCallEvents(ctx context.Context, callID string) ([]*CallEvent, error)
}

