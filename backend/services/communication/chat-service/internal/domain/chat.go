package domain

import (
	"time"

	"github.com/google/uuid"
	"github.com/titan-commerce/backend/pkg/errors"
)

// MessageType represents the type of chat message
type MessageType string

const (
	MessageTypeText     MessageType = "TEXT"
	MessageTypeImage    MessageType = "IMAGE"
	MessageTypeFile     MessageType = "FILE"
	MessageTypeProduct  MessageType = "PRODUCT"
	MessageTypeOrder    MessageType = "ORDER"
	MessageTypeSystem   MessageType = "SYSTEM"
)

// MessageStatus represents the delivery status of a message
type MessageStatus string

const (
	MessageStatusSent      MessageStatus = "SENT"
	MessageStatusDelivered MessageStatus = "DELIVERED"
	MessageStatusRead      MessageStatus = "READ"
	MessageStatusFailed    MessageStatus = "FAILED"
)

// ConversationType represents the type of conversation
type ConversationType string

const (
	ConversationTypeBuyerSeller ConversationType = "BUYER_SELLER"
	ConversationTypeSupport     ConversationType = "SUPPORT"
	ConversationTypeGroup       ConversationType = "GROUP"
)

// Message represents a single chat message
type Message struct {
	MessageID      string
	ConversationID string
	SenderID       string
	SenderName     string
	MessageType    MessageType
	Content        string        // Text content or JSON for rich messages
	Metadata       *MessageMetadata
	Status         MessageStatus
	IsEdited       bool
	IsDeleted      bool
	CreatedAt      time.Time
	UpdatedAt      time.Time
	ReadAt         *time.Time
}

// MessageMetadata contains additional message information
type MessageMetadata struct {
	ImageURL      string            `json:"image_url,omitempty"`
	FileURL       string            `json:"file_url,omitempty"`
	FileName      string            `json:"file_name,omitempty"`
	FileSize      int64             `json:"file_size,omitempty"`
	ProductID     string            `json:"product_id,omitempty"`
	OrderID       string            `json:"order_id,omitempty"`
	ThumbnailURL  string            `json:"thumbnail_url,omitempty"`
	CustomData    map[string]string `json:"custom_data,omitempty"`
}

// NewMessage creates a new chat message
func NewMessage(conversationID, senderID, senderName string, messageType MessageType, content string) (*Message, error) {
	if conversationID == "" {
		return nil, errors.New(errors.ErrInvalidInput, "conversation ID is required")
	}
	if senderID == "" {
		return nil, errors.New(errors.ErrInvalidInput, "sender ID is required")
	}
	if content == "" && messageType == MessageTypeText {
		return nil, errors.New(errors.ErrInvalidInput, "message content is required")
	}

	now := time.Now()
	return &Message{
		MessageID:      uuid.New().String(),
		ConversationID: conversationID,
		SenderID:       senderID,
		SenderName:     senderName,
		MessageType:    messageType,
		Content:        content,
		Status:         MessageStatusSent,
		IsEdited:       false,
		IsDeleted:      false,
		CreatedAt:      now,
		UpdatedAt:      now,
	}, nil
}

// MarkAsDelivered marks the message as delivered
func (m *Message) MarkAsDelivered() {
	m.Status = MessageStatusDelivered
	m.UpdatedAt = time.Now()
}

// MarkAsRead marks the message as read
func (m *Message) MarkAsRead() {
	if m.Status != MessageStatusRead {
		now := time.Now()
		m.Status = MessageStatusRead
		m.ReadAt = &now
		m.UpdatedAt = now
	}
}

// Edit edits the message content
func (m *Message) Edit(newContent string) error {
	if m.IsDeleted {
		return errors.New(errors.ErrInvalidInput, "cannot edit deleted message")
	}
	if newContent == "" {
		return errors.New(errors.ErrInvalidInput, "message content cannot be empty")
	}

	m.Content = newContent
	m.IsEdited = true
	m.UpdatedAt = time.Now()
	return nil
}

// Delete soft deletes the message
func (m *Message) Delete() error {
	if m.IsDeleted {
		return errors.New(errors.ErrInvalidInput, "message already deleted")
	}

	m.IsDeleted = true
	m.Content = "[Message deleted]"
	m.UpdatedAt = time.Now()
	return nil
}

// Conversation represents a chat conversation between users
type Conversation struct {
	ConversationID   string
	Type             ConversationType
	ParticipantIDs   []string
	ParticipantNames map[string]string // userID -> userName
	LastMessage      *Message
	UnreadCounts     map[string]int // userID -> unread count
	IsMuted          map[string]bool // userID -> mute status
	CreatedAt        time.Time
	UpdatedAt        time.Time
}

// NewConversation creates a new conversation
func NewConversation(conversationType ConversationType, participantIDs []string) (*Conversation, error) {
	if len(participantIDs) < 2 {
		return nil, errors.New(errors.ErrInvalidInput, "conversation must have at least 2 participants")
	}

	now := time.Now()
	return &Conversation{
		ConversationID:   uuid.New().String(),
		Type:             conversationType,
		ParticipantIDs:   participantIDs,
		ParticipantNames: make(map[string]string),
		UnreadCounts:     make(map[string]int),
		IsMuted:          make(map[string]bool),
		CreatedAt:        now,
		UpdatedAt:        now,
	}, nil
}

// AddParticipant adds a new participant to the conversation
func (c *Conversation) AddParticipant(userID, userName string) error {
	if userID == "" {
		return errors.New(errors.ErrInvalidInput, "user ID is required")
	}

	for _, id := range c.ParticipantIDs {
		if id == userID {
			return errors.New(errors.ErrInvalidInput, "user already in conversation")
		}
	}

	c.ParticipantIDs = append(c.ParticipantIDs, userID)
	c.ParticipantNames[userID] = userName
	c.UnreadCounts[userID] = 0
	c.IsMuted[userID] = false
	c.UpdatedAt = time.Now()
	return nil
}

// RemoveParticipant removes a participant from the conversation
func (c *Conversation) RemoveParticipant(userID string) error {
	newParticipants := []string{}
	found := false

	for _, id := range c.ParticipantIDs {
		if id != userID {
			newParticipants = append(newParticipants, id)
		} else {
			found = true
		}
	}

	if !found {
		return errors.New(errors.ErrNotFound, "user not in conversation")
	}

	c.ParticipantIDs = newParticipants
	delete(c.ParticipantNames, userID)
	delete(c.UnreadCounts, userID)
	delete(c.IsMuted, userID)
	c.UpdatedAt = time.Now()
	return nil
}

// UpdateLastMessage updates the last message in the conversation
func (c *Conversation) UpdateLastMessage(message *Message) {
	c.LastMessage = message
	c.UpdatedAt = time.Now()

	// Increment unread count for all participants except sender
	for _, participantID := range c.ParticipantIDs {
		if participantID != message.SenderID {
			c.UnreadCounts[participantID]++
		}
	}
}

// MarkAsRead marks all messages as read for a user
func (c *Conversation) MarkAsRead(userID string) {
	c.UnreadCounts[userID] = 0
	c.UpdatedAt = time.Now()
}

// Mute mutes the conversation for a user
func (c *Conversation) Mute(userID string) {
	c.IsMuted[userID] = true
	c.UpdatedAt = time.Now()
}

// Unmute unmutes the conversation for a user
func (c *Conversation) Unmute(userID string) {
	c.IsMuted[userID] = false
	c.UpdatedAt = time.Now()
}

// IsParticipant checks if a user is a participant in the conversation
func (c *Conversation) IsParticipant(userID string) bool {
	for _, id := range c.ParticipantIDs {
		if id == userID {
			return true
		}
	}
	return false
}

// TypingIndicator represents a user typing in a conversation
type TypingIndicator struct {
	ConversationID string
	UserID         string
	UserName       string
	IsTyping       bool
	Timestamp      time.Time
}

// NewTypingIndicator creates a new typing indicator
func NewTypingIndicator(conversationID, userID, userName string, isTyping bool) *TypingIndicator {
	return &TypingIndicator{
		ConversationID: conversationID,
		UserID:         userID,
		UserName:       userName,
		IsTyping:       isTyping,
		Timestamp:      time.Now(),
	}
}

// OnlineStatus represents a user's online status
type OnlineStatus struct {
	UserID       string
	IsOnline     bool
	LastSeenAt   time.Time
	DeviceInfo   string
}

// NewOnlineStatus creates a new online status
func NewOnlineStatus(userID string) *OnlineStatus {
	return &OnlineStatus{
		UserID:     userID,
		IsOnline:   true,
		LastSeenAt: time.Now(),
	}
}

// SetOffline marks the user as offline
func (o *OnlineStatus) SetOffline() {
	o.IsOnline = false
	o.LastSeenAt = time.Now()
}

// SetOnline marks the user as online
func (o *OnlineStatus) SetOnline() {
	o.IsOnline = true
	o.LastSeenAt = time.Now()
}

