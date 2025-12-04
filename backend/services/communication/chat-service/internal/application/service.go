package application

import (
	"context"

	"github.com/titan-commerce/backend/chat-service/internal/domain"
	"github.com/titan-commerce/backend/pkg/logger"
)

type MessageRepository interface {
	Save(ctx context.Context, message *domain.Message) error
	FindByConversationID(ctx context.Context, conversationID string, page, pageSize int) ([]*domain.Message, int, error)
	Update(ctx context.Context, message *domain.Message) error
}

type ConversationRepository interface {
	Save(ctx context.Context, conversation *domain.Conversation) error
	FindByID(ctx context.Context, conversationID string) (*domain.Conversation, error)
	FindByUserID(ctx context.Context, userID string, pageSize int) ([]*domain.Conversation, error)
	Update(ctx context.Context, conversation *domain.Conversation) error
}

type ChatService struct {
	msgRepo  MessageRepository
	convRepo ConversationRepository
	logger   *logger.Logger
}

func NewChatService(msgRepo MessageRepository, convRepo ConversationRepository, logger *logger.Logger) *ChatService {
	return &ChatService{
		msgRepo:  msgRepo,
		convRepo: convRepo,
		logger:   logger,
	}
}

// SendMessage sends a message (Command)
func (s *ChatService) SendMessage(ctx context.Context, senderID, receiverID string, msgType domain.MessageType, content, mediaURL string) (*domain.Message, error) {
	// Find or create conversation
	conversationID := generateConversationID(senderID, receiverID)

	message := domain.NewMessage(conversationID, senderID, receiverID, msgType, content, mediaURL)

	if err := s.msgRepo.Save(ctx, message); err != nil {
		s.logger.Error(err, "failed to save message")
		return nil, err
	}

	// Update conversation
	conversation, err := s.convRepo.FindByID(ctx, conversationID)
	if err != nil {
		conversation = domain.NewConversation(senderID, receiverID)
		s.convRepo.Save(ctx, conversation)
	}

	conversation.UpdateLastMessage(message.ID)
	s.convRepo.Update(ctx, conversation)

	s.logger.Infof("Message sent: from=%s, to=%s, type=%s", senderID, receiverID, msgType)
	return message, nil
}

// GetMessages retrieves messages for a conversation (Query)
func (s *ChatService) GetMessages(ctx context.Context, conversationID string, page, pageSize int) ([]*domain.Message, int, error) {
	return s.msgRepo.FindByConversationID(ctx, conversationID, page, pageSize)
}

// GetConversations retrieves user's conversations (Query)
func (s *ChatService) GetConversations(ctx context.Context, userID string, pageSize int) ([]*domain.Conversation, error) {
	return s.convRepo.FindByUserID(ctx, userID, pageSize)
}

func generateConversationID(user1, user2 string) string {
	// Simple deterministic conversation ID
	if user1 < user2 {
		return user1 + "_" + user2
	}
	return user2 + "_" + user1
}
