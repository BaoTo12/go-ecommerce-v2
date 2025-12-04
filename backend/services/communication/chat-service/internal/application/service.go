package application
}
	return s.presenceRepo.GetOnlineUsers(ctx, userIDs)
func (s *ChatService) GetMultipleOnlineStatuses(ctx context.Context, userIDs []string) (map[string]bool, error) {
// GetMultipleOnlineStatuses retrieves online status for multiple users (Query)

}
	return s.presenceRepo.GetTypingIndicators(ctx, conversationID)
func (s *ChatService) GetTypingIndicators(ctx context.Context, conversationID string) ([]*domain.TypingIndicator, error) {
// GetTypingIndicators retrieves typing indicators for a conversation (Query)

}
	return s.presenceRepo.GetOnlineStatus(ctx, userID)
func (s *ChatService) GetOnlineStatus(ctx context.Context, userID string) (*domain.OnlineStatus, error) {
// GetOnlineStatus retrieves user's online status (Query)

}
	return s.messageRepo.GetMessagesByConversation(ctx, conversationID, limit, beforeTimestamp)

	}
		limit = 50 // Default page size
	if limit <= 0 || limit > 100 {
) ([]*domain.Message, error) {
	beforeTimestamp *int64,
	limit int,
	conversationID string,
	ctx context.Context,
func (s *ChatService) GetMessages(
// GetMessages retrieves messages from a conversation (Query)

}
	return conversations, nil

	})
		return conversations[i].LastMessage.CreatedAt.After(conversations[j].LastMessage.CreatedAt)
		}
			return true
		if conversations[j].LastMessage == nil {
		}
			return false
		if conversations[i].LastMessage == nil {
	sort.Slice(conversations, func(i, j int) bool {
	// Sort by last message timestamp (most recent first)

	}
		return nil, err
	if err != nil {
	conversations, err := s.conversationRepo.GetUserConversations(ctx, userID)
func (s *ChatService) GetUserConversations(ctx context.Context, userID string) ([]*domain.Conversation, error) {
// GetUserConversations retrieves all conversations for a user (Query)

}
	return s.conversationRepo.GetConversation(ctx, conversationID)
func (s *ChatService) GetConversation(ctx context.Context, conversationID string) (*domain.Conversation, error) {
// GetConversation retrieves a conversation (Query)

}
	return nil
	s.logger.Infof("User status updated: user=%s, online=%v", userID, online)

	}
		return err
		s.logger.Error(err, "failed to set online status")
	if err := s.presenceRepo.SetOnlineStatus(ctx, userID, online); err != nil {
func (s *ChatService) SetUserOnlineStatus(ctx context.Context, userID string, online bool) error {
// SetUserOnlineStatus sets user's online/offline status (Command)

}
	return nil

	}
		return err
		s.logger.Error(err, "failed to set typing indicator")
	if err := s.presenceRepo.SetTypingIndicator(ctx, indicator); err != nil {

	indicator := domain.NewTypingIndicator(conversationID, userID, userName, isTyping)
) error {
	isTyping bool,
	conversationID, userID, userName string,
	ctx context.Context,
func (s *ChatService) SetTypingIndicator(
// SetTypingIndicator sets typing indicator for a user in a conversation (Command)

}
	return nil
	s.logger.Infof("Messages marked as read: conversation=%s, user=%s", conversationID, userID)

	}
		return err
		s.logger.Error(err, "failed to mark messages as read")
	if err := s.messageRepo.MarkMessagesAsRead(ctx, conversationID, userID); err != nil {

	}
		return err
		s.logger.Error(err, "failed to mark conversation as read")
	if err := s.conversationRepo.UpdateConversation(ctx, conversation); err != nil {
	conversation.MarkAsRead(userID)

	}
		return errors.New(errors.ErrUnauthorized, "user is not a participant")
	if !conversation.IsParticipant(userID) {

	}
		return err
	if err != nil {
	conversation, err := s.conversationRepo.GetConversation(ctx, conversationID)
func (s *ChatService) MarkMessagesAsRead(ctx context.Context, conversationID, userID string) error {
// MarkMessagesAsRead marks all messages in a conversation as read (Command)

}
	return nil
	s.logger.Infof("Message deleted: id=%s", messageID)

	}
		return err
		s.logger.Error(err, "failed to delete message")
	if err := s.messageRepo.UpdateMessage(ctx, message); err != nil {

	}
		return err
	if err := message.Delete(); err != nil {

	}
		return errors.New(errors.ErrUnauthorized, "only the sender can delete their message")
	if message.SenderID != userID {

	}
		return err
	if err != nil {
	message, err := s.messageRepo.GetMessage(ctx, messageID)
func (s *ChatService) DeleteMessage(ctx context.Context, messageID, userID string) error {
// DeleteMessage soft deletes a message (Command)

}
	return nil
	s.logger.Infof("Message edited: id=%s", messageID)

	}
		return err
		s.logger.Error(err, "failed to edit message")
	if err := s.messageRepo.UpdateMessage(ctx, message); err != nil {

	}
		return err
	if err := message.Edit(newContent); err != nil {

	}
		return errors.New(errors.ErrUnauthorized, "only the sender can edit their message")
	if message.SenderID != userID {

	}
		return err
	if err != nil {
	message, err := s.messageRepo.GetMessage(ctx, messageID)
func (s *ChatService) EditMessage(ctx context.Context, messageID, userID, newContent string) error {
// EditMessage edits an existing message (Command)

}
	return message, nil

		message.MessageID, conversationID, senderID)
	s.logger.Infof("Message sent: id=%s, conversation=%s, sender=%s",

	s.conversationRepo.UpdateConversation(ctx, conversation)
	conversation.UpdateLastMessage(message)
	// Update conversation's last message

	}
		return nil, err
		s.logger.Error(err, "failed to create message")
	if err := s.messageRepo.CreateMessage(ctx, message); err != nil {

	}
		message.Metadata = metadata
	if metadata != nil {

	}
		return nil, err
	if err != nil {
	message, err := domain.NewMessage(conversationID, senderID, senderName, messageType, content)

	}
		return nil, errors.New(errors.ErrUnauthorized, "user is not a participant in this conversation")
	if !conversation.IsParticipant(senderID) {

	}
		return nil, err
	if err != nil {
	conversation, err := s.conversationRepo.GetConversation(ctx, conversationID)
	// Verify conversation exists and sender is a participant
) (*domain.Message, error) {
	metadata *domain.MessageMetadata,
	content string,
	messageType domain.MessageType,
	conversationID, senderID, senderName string,
	ctx context.Context,
func (s *ChatService) SendMessage(
// SendMessage sends a message in a conversation (Command)

}
	return conversation, nil

		conversation.ConversationID, len(participantIDs))
	s.logger.Infof("Conversation created: id=%s, participants=%d",

	}
		return nil, err
		s.logger.Error(err, "failed to create conversation")
	if err := s.conversationRepo.CreateConversation(ctx, conversation); err != nil {

	}
		return nil, err
	if err != nil {
	conversation, err := domain.NewConversation(conversationType, participantIDs)

	}
		return existingConv, nil
		s.logger.Infof("Conversation already exists: %s", existingConv.ConversationID)
	if existingConv != nil {
	existingConv, _ := s.conversationRepo.FindConversationByParticipants(ctx, participantIDs)
	// Check if conversation already exists between these participants
) (*domain.Conversation, error) {
	participantIDs []string,
	conversationType domain.ConversationType,
	ctx context.Context,
func (s *ChatService) CreateConversation(
// CreateConversation creates a new conversation (Command)

}
	}
		logger:           logger,
		presenceRepo:     presenceRepo,
		conversationRepo: conversationRepo,
		messageRepo:      messageRepo,
	return &ChatService{
) *ChatService {
	logger *logger.Logger,
	presenceRepo domain.PresenceRepository,
	conversationRepo domain.ConversationRepository,
	messageRepo domain.MessageRepository,
func NewChatService(

}
	logger           *logger.Logger
	presenceRepo     domain.PresenceRepository
	conversationRepo domain.ConversationRepository
	messageRepo      domain.MessageRepository
type ChatService struct {

)
	"github.com/titan-commerce/backend/pkg/logger"
	"github.com/titan-commerce/backend/pkg/errors"
	"github.com/titan-commerce/backend/chat-service/internal/domain"

	"sort"
	"context"
import (


