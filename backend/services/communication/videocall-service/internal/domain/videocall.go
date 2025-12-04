package domain

import (
	"time"

	"github.com/google/uuid"
)

type CallStatus string

const (
	CallStatusRinging    CallStatus = "RINGING"
	CallStatusInProgress CallStatus = "IN_PROGRESS"
	CallStatusEnded      CallStatus = "ENDED"
	CallStatusMissed     CallStatus = "MISSED"
	CallStatusRejected   CallStatus = "REJECTED"
)

type CallType string

const (
	CallTypeVideo CallType = "VIDEO"
	CallTypeAudio CallType = "AUDIO"
)

type VideoCall struct {
	ID         string
	CallerID   string
	CalleeID   string
	Type       CallType
	Status     CallStatus
	RoomID     string
	StartedAt  *time.Time
	EndedAt    *time.Time
	Duration   int // seconds
	CreatedAt  time.Time
}

type SignalingMessage struct {
	Type      string                 `json:"type"` // offer, answer, ice-candidate
	CallID    string                 `json:"call_id"`
	FromUser  string                 `json:"from_user"`
	ToUser    string                 `json:"to_user"`
	Payload   map[string]interface{} `json:"payload"`
	CreatedAt time.Time              `json:"created_at"`
}

func NewVideoCall(callerID, calleeID string, callType CallType) *VideoCall {
	return &VideoCall{
		ID:        uuid.New().String(),
		CallerID:  callerID,
		CalleeID:  calleeID,
		Type:      callType,
		Status:    CallStatusRinging,
		RoomID:    uuid.New().String(),
		CreatedAt: time.Now(),
	}
}

func (c *VideoCall) Accept() {
	now := time.Now()
	c.Status = CallStatusInProgress
	c.StartedAt = &now
}

func (c *VideoCall) End() {
	now := time.Now()
	c.Status = CallStatusEnded
	c.EndedAt = &now
	if c.StartedAt != nil {
		c.Duration = int(now.Sub(*c.StartedAt).Seconds())
	}
}

func (c *VideoCall) Reject() {
	c.Status = CallStatusRejected
	now := time.Now()
	c.EndedAt = &now
}

func (c *VideoCall) Miss() {
	c.Status = CallStatusMissed
	now := time.Now()
	c.EndedAt = &now
}

type Repository interface {
	Save(ctx interface{}, call *VideoCall) error
	FindByID(ctx interface{}, callID string) (*VideoCall, error)
	FindByRoomID(ctx interface{}, roomID string) (*VideoCall, error)
	Update(ctx interface{}, call *VideoCall) error
	FindCallHistory(ctx interface{}, userID string, limit int) ([]*VideoCall, error)
}
