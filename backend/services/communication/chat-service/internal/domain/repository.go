package domain

import "context"

// MessageRepository defines the interface for message persistence
type MessageRepository interface {
	// Message operations
	CreateMessage(ctx context.Context, message *Message) error
	GetMessage(ctx context.Context, messageID string) (*Message, error)
	UpdateMessage(ctx context.Context, message *Message) error
	GetMessagesByConversation(ctx context.Context, conversationID string, limit int, beforeTimestamp *int64) ([]*Message, error)
	MarkMessagesAsRead(ctx context.Context, conversationID, userID string) error
	DeleteMessage(ctx context.Context, messageID string) error
}

// ConversationRepository defines the interface for conversation persistence
type ConversationRepository interface {
	// Conversation operations
	CreateConversation(ctx context.Context, conversation *Conversation) error
	GetConversation(ctx context.Context, conversationID string) (*Conversation, error)
	UpdateConversation(ctx context.Context, conversation *Conversation) error
	GetUserConversations(ctx context.Context, userID string) ([]*Conversation, error)
	FindConversationByParticipants(ctx context.Context, participantIDs []string) (*Conversation, error)
	DeleteConversation(ctx context.Context, conversationID string) error
}

// PresenceRepository defines the interface for user presence tracking
type PresenceRepository interface {
	// Presence operations
	SetOnlineStatus(ctx context.Context, userID string, online bool) error
	GetOnlineStatus(ctx context.Context, userID string) (*OnlineStatus, error)
	GetOnlineUsers(ctx context.Context, userIDs []string) (map[string]bool, error)

	// Typing indicators
	SetTypingIndicator(ctx context.Context, indicator *TypingIndicator) error
	GetTypingIndicators(ctx context.Context, conversationID string) ([]*TypingIndicator, error)
}

