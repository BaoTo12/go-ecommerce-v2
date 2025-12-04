package domain

import (
	"time"

	"github.com/google/uuid"
)

type MessageType string

const (
	MessageTypeText  MessageType = "TEXT"
	MessageTypeImage MessageType = "IMAGE"
	MessageTypeFile  MessageType = "FILE"
)

type Message struct {
	ID           string
	ConversationID string
	SenderID     string
	Content      string
	Type         MessageType
	Metadata     map[string]string
	ReadBy       []string
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

type Conversation struct {
	ID           string
	Type         string // "direct" or "group"
	Participants []string
	Name         string
	LastMessage  *Message
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

func NewMessage(conversationID, senderID, content string, msgType MessageType) *Message {
	return &Message{
		ID:             uuid.New().String(),
		ConversationID: conversationID,
		SenderID:       senderID,
		Content:        content,
		Type:           msgType,
		Metadata:       make(map[string]string),
		ReadBy:         []string{senderID},
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
	}
}

func NewDirectConversation(user1ID, user2ID string) *Conversation {
	return &Conversation{
		ID:           uuid.New().String(),
		Type:         "direct",
		Participants: []string{user1ID, user2ID},
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}
}

func NewGroupConversation(name string, participants []string) *Conversation {
	return &Conversation{
		ID:           uuid.New().String(),
		Type:         "group",
		Name:         name,
		Participants: participants,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}
}

func (m *Message) MarkRead(userID string) {
	for _, u := range m.ReadBy {
		if u == userID {
			return
		}
	}
	m.ReadBy = append(m.ReadBy, userID)
	m.UpdatedAt = time.Now()
}

type Repository interface {
	SaveMessage(ctx interface{}, message *Message) error
	FindMessagesByConversation(ctx interface{}, conversationID string, limit, offset int) ([]*Message, error)
	SaveConversation(ctx interface{}, conversation *Conversation) error
	FindConversationByID(ctx interface{}, conversationID string) (*Conversation, error)
	FindConversationsByUser(ctx interface{}, userID string) ([]*Conversation, error)
	FindDirectConversation(ctx interface{}, user1ID, user2ID string) (*Conversation, error)
}
