package domain

import (
	"time"

	"github.com/google/uuid"
)

type MessageType string

const (
	MessageTypeText        MessageType = "TEXT"
	MessageTypeImage       MessageType = "IMAGE"
	MessageTypeVideo       MessageType = "VIDEO"
	MessageTypeFile        MessageType = "FILE"
	MessageTypeProductLink MessageType = "PRODUCT_LINK"
)

type Message struct {
	ID             string
	ConversationID string
	SenderID       string
	ReceiverID     string
	Type           MessageType
	Content        string
	MediaURL       string
	Read           bool
	CreatedAt      time.Time
}

type Conversation struct {
	ID            string
	ParticipantIDs []string
	LastMessageID string
	UnreadCount   int
	UpdatedAt     time.Time
}

func NewMessage(conversationID, senderID, receiverID string, msgType MessageType, content, mediaURL string) *Message {
	return &Message{
		ID:             uuid.New().String(),
		ConversationID: conversationID,
		SenderID:       senderID,
		ReceiverID:     receiverID,
		Type:           msgType,
		Content:        content,
		MediaURL:       mediaURL,
		Read:           false,
		CreatedAt:      time.Now(),
	}
}

func (m *Message) MarkAsRead() {
	m.Read = true
}

func NewConversation(participant1ID, participant2ID string) *Conversation {
	return &Conversation{
		ID:             uuid.New().String(),
		ParticipantIDs: []string{participant1ID, participant2ID},
		UnreadCount:    0,
		UpdatedAt:      time.Now(),
	}
}

func (c *Conversation) UpdateLastMessage(messageID string) {
	c.LastMessageID = messageID
	c.UnreadCount++
	c.UpdatedAt = time.Now()
}
