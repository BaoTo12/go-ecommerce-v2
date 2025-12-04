package domain
}
	}
		CreatedAt: time.Now(),
		Payload:   payload,
		EventType: eventType,
		UserID:    userID,
		CallID:    callID,
		EventID:   uuid.New().String(),
	return &CallEvent{
func NewCallEvent(callID, userID, eventType, payload string) *CallEvent {
// NewCallEvent creates a new call event

}
	CreatedAt  time.Time
	Payload    string // SDP or ICE candidate JSON
	EventType  string // offer, answer, ice-candidate, hangup
	UserID     string
	CallID     string
	EventID    string
type CallEvent struct {
// CallEvent represents a signaling event during the call

}
	}
		},
			URLs: []string{"stun:stun1.l.google.com:19302"},
		{
		},
			URLs: []string{"stun:stun.l.google.com:19302"},
		{
	return []IceServer{
func getDefaultIceServers() []IceServer {
// getDefaultIceServers returns default STUN/TURN servers

}
	c.UpdatedAt = time.Now()
	c.IsRecorded = true
	c.RecordingURL = url
func (c *VideoCall) SetRecordingURL(url string) {
// SetRecordingURL sets the recording URL

}
	c.UpdatedAt = time.Now()
	c.Quality = quality
func (c *VideoCall) UpdateQuality(quality *CallQuality) {
// UpdateQuality updates the call quality metrics

}
	return nil
	c.UpdatedAt = time.Now()
	*c.EndedAt = time.Now()
	c.EndedAt = new(time.Time)
	c.Status = CallRejected

	}
		return errors.New(errors.ErrInvalidInput, "can only reject incoming calls")
	if c.Status != CallRinging && c.Status != CallInitiated {
func (c *VideoCall) Reject() error {
// Reject rejects the incoming call

}
	return nil

	}
		c.Status = CallEnded
	} else {
		c.Status = CallFailed
	} else if c.Status == CallInitiated {
		c.Status = CallMissed
	} else if c.Status == CallRinging {
		c.Status = CallEnded
		c.Duration = int(now.Sub(*c.StartedAt).Seconds())
	if c.Status == CallActive && c.StartedAt != nil {

	c.UpdatedAt = now
	c.EndedAt = &now
	now := time.Now()

	}
		return errors.New(errors.ErrInvalidInput, "call already ended")
	if c.Status == CallEnded {
func (c *VideoCall) End() error {
// End ends the call

}
	return nil
	c.UpdatedAt = now
	c.StartedAt = &now
	c.Status = CallActive
	now := time.Now()

	}
		return errors.New(errors.ErrInvalidInput, "can only answer ringing calls")
	if c.Status != CallRinging {
func (c *VideoCall) Answer() error {
// Answer answers the call and starts the session

}
	return nil
	c.UpdatedAt = time.Now()
	c.Status = CallRinging

	}
		return errors.New(errors.ErrInvalidInput, "call must be initiated to ring")
	if c.Status != CallInitiated {
func (c *VideoCall) Ring() error {
// Ring marks the call as ringing

}
	}, nil
		UpdatedAt:  now,
		CreatedAt:  now,
		IsRecorded: false,
		IceServers: getDefaultIceServers(),
		RoomID:     "room_" + callID,
		Status:     CallInitiated,
		CalleeName: calleeName,
		CalleeID:   calleeID,
		CallerName: callerName,
		CallerID:   callerID,
		Type:       callType,
		CallID:     callID,
	return &VideoCall{

	callID := uuid.New().String()
	now := time.Now()

	}
		return nil, errors.New(errors.ErrInvalidInput, "cannot call yourself")
	if callerID == calleeID {
	}
		return nil, errors.New(errors.ErrInvalidInput, "callee ID is required")
	if calleeID == "" {
	}
		return nil, errors.New(errors.ErrInvalidInput, "caller ID is required")
	if callerID == "" {
func NewVideoCall(callType CallType, callerID, callerName, calleeID, calleeName string) (*VideoCall, error) {
// NewVideoCall creates a new video call

}
	ConnectionDrops    int
	AudioQuality       string
	VideoResolution    string
	AverageBitrate     int     // in kbps
	PacketLoss         float64 // percentage
	AverageLatency     int     // in ms
type CallQuality struct {
// CallQuality represents quality metrics for the call

}
	Credential string
	Username   string
	URLs       []string
type IceServer struct {
// IceServer represents a STUN/TURN server config

}
	UpdatedAt    time.Time
	CreatedAt    time.Time
	IsRecorded   bool
	RecordingURL string
	Quality      *CallQuality
	Duration     int // in seconds
	EndedAt      *time.Time
	StartedAt    *time.Time
	IceServers   []IceServer
	RoomID       string // WebRTC room identifier
	Status       CallStatus
	CalleeName   string
	CalleeID     string
	CallerName   string
	CallerID     string
	Type         CallType
	CallID       string
type VideoCall struct {
// VideoCall represents a 1-on-1 video call session

)
	CallTypeConsult CallType = "CONSULT"
	CallTypeSales   CallType = "SALES"
	CallTypeSupport CallType = "SUPPORT"
const (

type CallType string
// CallType represents the type of call

)
	CallFailed    CallStatus = "FAILED"
	CallRejected  CallStatus = "REJECTED"
	CallMissed    CallStatus = "MISSED"
	CallEnded     CallStatus = "ENDED"
	CallActive    CallStatus = "ACTIVE"
	CallRinging   CallStatus = "RINGING"
	CallInitiated CallStatus = "INITIATED"
const (

type CallStatus string
// CallStatus represents the status of a video call

)
	"github.com/titan-commerce/backend/pkg/errors"
	"github.com/google/uuid"

	"time"
import (


