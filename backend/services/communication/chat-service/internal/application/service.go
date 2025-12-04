package application

import (
	"context"
	"sync"

	"github.com/titan-commerce/backend/chat-service/internal/domain"
	"github.com/titan-commerce/backend/pkg/logger"
)

type ChatRepository interface {
	SaveMessage(ctx context.Context, message *domain.Message) error
	FindMessagesByConversation(ctx context.Context, conversationID string, limit, offset int) ([]*domain.Message, error)
	SaveConversation(ctx context.Context, conversation *domain.Conversation) error
	FindConversationByID(ctx context.Context, conversationID string) (*domain.Conversation, error)
	FindConversationsByUser(ctx context.Context, userID string) ([]*domain.Conversation, error)
	FindDirectConversation(ctx context.Context, user1ID, user2ID string) (*domain.Conversation, error)
	UpdateConversation(ctx context.Context, conversation *domain.Conversation) error
}

// WebSocket connection manager
type ConnectionManager struct {
	connections map[string]map[string]chan *domain.Message // userID -> connectionID -> channel
	mu          sync.RWMutex
}

func NewConnectionManager() *ConnectionManager {
	return &ConnectionManager{
		connections: make(map[string]map[string]chan *domain.Message),
	}
}

func (cm *ConnectionManager) AddConnection(userID, connID string, ch chan *domain.Message) {
	cm.mu.Lock()
	defer cm.mu.Unlock()
	if cm.connections[userID] == nil {
		cm.connections[userID] = make(map[string]chan *domain.Message)
	}
	cm.connections[userID][connID] = ch
}

func (cm *ConnectionManager) RemoveConnection(userID, connID string) {
	cm.mu.Lock()
	defer cm.mu.Unlock()
	if conns, ok := cm.connections[userID]; ok {
		if ch, ok := conns[connID]; ok {
			close(ch)
			delete(conns, connID)
		}
		if len(conns) == 0 {
			delete(cm.connections, userID)
		}
	}
}

func (cm *ConnectionManager) SendToUser(userID string, msg *domain.Message) {
	cm.mu.RLock()
	defer cm.mu.RUnlock()
	if conns, ok := cm.connections[userID]; ok {
		for _, ch := range conns {
			select {
			case ch <- msg:
			default:
				// Channel full, skip
			}
		}
	}
}

type ChatService struct {
	repo       ChatRepository
	connMgr    *ConnectionManager
	logger     *logger.Logger
}

func NewChatService(repo ChatRepository, logger *logger.Logger) *ChatService {
	return &ChatService{
		repo:    repo,
		connMgr: NewConnectionManager(),
		logger:  logger,
	}
}

func (s *ChatService) GetConnectionManager() *ConnectionManager {
	return s.connMgr
}

// SendMessage sends a message to a conversation
func (s *ChatService) SendMessage(ctx context.Context, conversationID, senderID, content string, msgType domain.MessageType) (*domain.Message, error) {
	// Verify conversation exists
	conv, err := s.repo.FindConversationByID(ctx, conversationID)
	if err != nil {
		return nil, err
	}

	// Create message
	message := domain.NewMessage(conversationID, senderID, content, msgType)
	
	if err := s.repo.SaveMessage(ctx, message); err != nil {
		s.logger.Error(err, "failed to save message")
		return nil, err
	}

	// Update last message
	conv.LastMessage = message
	s.repo.UpdateConversation(ctx, conv)

	// Broadcast to all participants
	for _, userID := range conv.Participants {
		if userID != senderID {
			s.connMgr.SendToUser(userID, message)
		}
	}

	s.logger.Infof("Message sent: %s in conversation %s", message.ID, conversationID)
	return message, nil
}

// GetMessages retrieves messages in a conversation
func (s *ChatService) GetMessages(ctx context.Context, conversationID string, limit, offset int) ([]*domain.Message, error) {
	return s.repo.FindMessagesByConversation(ctx, conversationID, limit, offset)
}

// GetOrCreateDirectConversation gets or creates a direct conversation between two users
func (s *ChatService) GetOrCreateDirectConversation(ctx context.Context, user1ID, user2ID string) (*domain.Conversation, error) {
	conv, err := s.repo.FindDirectConversation(ctx, user1ID, user2ID)
	if err == nil {
		return conv, nil
	}

	// Create new conversation
	conv = domain.NewDirectConversation(user1ID, user2ID)
	if err := s.repo.SaveConversation(ctx, conv); err != nil {
		return nil, err
	}

	s.logger.Infof("Created direct conversation: %s", conv.ID)
	return conv, nil
}

// CreateGroupConversation creates a new group conversation
func (s *ChatService) CreateGroupConversation(ctx context.Context, name string, participants []string) (*domain.Conversation, error) {
	conv := domain.NewGroupConversation(name, participants)
	if err := s.repo.SaveConversation(ctx, conv); err != nil {
		return nil, err
	}

	s.logger.Infof("Created group conversation: %s with %d participants", conv.ID, len(participants))
	return conv, nil
}

// GetConversations retrieves all conversations for a user
func (s *ChatService) GetConversations(ctx context.Context, userID string) ([]*domain.Conversation, error) {
	return s.repo.FindConversationsByUser(ctx, userID)
}

// MarkAsRead marks a message as read
func (s *ChatService) MarkAsRead(ctx context.Context, conversationID, userID string) error {
	// In production, update read status in database
	s.logger.Infof("Marked conversation %s as read by %s", conversationID, userID)
	return nil
}
