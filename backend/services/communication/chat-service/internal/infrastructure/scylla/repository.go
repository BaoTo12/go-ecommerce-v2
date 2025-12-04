package scylla

import (
	"context"
	"encoding/json"
	"time"

	"github.com/gocql/gocql"
	"github.com/titan-commerce/backend/chat-service/internal/domain"
	"github.com/titan-commerce/backend/pkg/errors"
	"github.com/titan-commerce/backend/pkg/logger"
)

type ChatRepository struct {
	session *gocql.Session
	logger  *logger.Logger
}

func NewChatRepository(hosts []string, keyspace string, logger *logger.Logger) (*ChatRepository, error) {
	cluster := gocql.NewCluster(hosts...)
	cluster.Keyspace = keyspace
	cluster.Consistency = gocql.Quorum
	cluster.Timeout = 5 * time.Second

	session, err := cluster.CreateSession()
	if err != nil {
		return nil, errors.Wrap(errors.ErrInternal, "failed to connect to ScyllaDB", err)
	}

	logger.Info("Chat ScyllaDB repository initialized")
	return &ChatRepository{session: session, logger: logger}, nil
}

func (r *ChatRepository) SaveMessage(ctx context.Context, msg *domain.Message) error {
	metadataJSON, _ := json.Marshal(msg.Metadata)

	query := `INSERT INTO messages (message_id, conversation_id, sender_id, sender_name, message_type, content, metadata, status, is_edited, is_deleted, created_at, updated_at, read_at)
			  VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`

	return r.session.Query(query,
		msg.MessageID, msg.ConversationID, msg.SenderID, msg.SenderName, msg.MessageType,
		msg.Content, metadataJSON, msg.Status, msg.IsEdited, msg.IsDeleted,
		msg.CreatedAt, msg.UpdatedAt, msg.ReadAt,
	).WithContext(ctx).Exec()
}

func (r *ChatRepository) GetMessage(ctx context.Context, messageID string) (*domain.Message, error) {
	query := `SELECT message_id, conversation_id, sender_id, sender_name, message_type, content, metadata, status, is_edited, is_deleted, created_at, updated_at, read_at
			  FROM messages WHERE message_id = ? LIMIT 1`

	var msg domain.Message
	var metadataJSON []byte

	if err := r.session.Query(query, messageID).WithContext(ctx).Scan(
		&msg.MessageID, &msg.ConversationID, &msg.SenderID, &msg.SenderName, &msg.MessageType,
		&msg.Content, &metadataJSON, &msg.Status, &msg.IsEdited, &msg.IsDeleted,
		&msg.CreatedAt, &msg.UpdatedAt, &msg.ReadAt,
	); err != nil {
		if err == gocql.ErrNotFound {
			return nil, errors.New(errors.ErrNotFound, "message not found")
		}
		return nil, err
	}

	if metadataJSON != nil {
		json.Unmarshal(metadataJSON, &msg.Metadata)
	}
	return &msg, nil
}

func (r *ChatRepository) UpdateMessage(ctx context.Context, msg *domain.Message) error {
	metadataJSON, _ := json.Marshal(msg.Metadata)

	query := `UPDATE messages SET content = ?, metadata = ?, status = ?, is_edited = ?, is_deleted = ?, updated_at = ?, read_at = ?
			  WHERE message_id = ?`

	return r.session.Query(query,
		msg.Content, metadataJSON, msg.Status, msg.IsEdited, msg.IsDeleted, msg.UpdatedAt, msg.ReadAt, msg.MessageID,
	).WithContext(ctx).Exec()
}

func (r *ChatRepository) GetConversationMessages(ctx context.Context, conversationID string, limit int) ([]*domain.Message, error) {
	query := `SELECT message_id, conversation_id, sender_id, sender_name, message_type, content, metadata, status, is_edited, is_deleted, created_at, updated_at, read_at
			  FROM messages_by_conversation WHERE conversation_id = ? ORDER BY created_at DESC LIMIT ?`

	iter := r.session.Query(query, conversationID, limit).WithContext(ctx).Iter()
	defer iter.Close()

	var messages []*domain.Message
	for {
		var msg domain.Message
		var metadataJSON []byte

		if !iter.Scan(
			&msg.MessageID, &msg.ConversationID, &msg.SenderID, &msg.SenderName, &msg.MessageType,
			&msg.Content, &metadataJSON, &msg.Status, &msg.IsEdited, &msg.IsDeleted,
			&msg.CreatedAt, &msg.UpdatedAt, &msg.ReadAt,
		) {
			break
		}

		if metadataJSON != nil {
			json.Unmarshal(metadataJSON, &msg.Metadata)
		}
		messages = append(messages, &msg)
	}

	if err := iter.Close(); err != nil {
		return nil, err
	}
	return messages, nil
}

func (r *ChatRepository) SaveConversation(ctx context.Context, conv *domain.Conversation) error {
	participantsJSON, _ := json.Marshal(conv.Participants)

	query := `INSERT INTO conversations (conversation_id, type, participants, created_at, updated_at, last_message_at)
			  VALUES (?, ?, ?, ?, ?, ?)`

	return r.session.Query(query,
		conv.ConversationID, conv.Type, participantsJSON, conv.CreatedAt, conv.UpdatedAt, conv.LastMessageAt,
	).WithContext(ctx).Exec()
}

func (r *ChatRepository) GetConversation(ctx context.Context, conversationID string) (*domain.Conversation, error) {
	query := `SELECT conversation_id, type, participants, created_at, updated_at, last_message_at
			  FROM conversations WHERE conversation_id = ? LIMIT 1`

	var conv domain.Conversation
	var participantsJSON []byte

	if err := r.session.Query(query, conversationID).WithContext(ctx).Scan(
		&conv.ConversationID, &conv.Type, &participantsJSON, &conv.CreatedAt, &conv.UpdatedAt, &conv.LastMessageAt,
	); err != nil {
		if err == gocql.ErrNotFound {
			return nil, errors.New(errors.ErrNotFound, "conversation not found")
		}
		return nil, err
	}

	json.Unmarshal(participantsJSON, &conv.Participants)
	return &conv, nil
}

func (r *ChatRepository) UpdateConversation(ctx context.Context, conv *domain.Conversation) error {
	participantsJSON, _ := json.Marshal(conv.Participants)

	query := `UPDATE conversations SET participants = ?, updated_at = ?, last_message_at = ?
			  WHERE conversation_id = ?`

	return r.session.Query(query,
		participantsJSON, conv.UpdatedAt, conv.LastMessageAt, conv.ConversationID,
	).WithContext(ctx).Exec()
}

func (r *ChatRepository) GetUserConversations(ctx context.Context, userID string) ([]*domain.Conversation, error) {
	query := `SELECT conversation_id, type, participants, created_at, updated_at, last_message_at
			  FROM conversations_by_user WHERE user_id = ? ORDER BY last_message_at DESC`

	iter := r.session.Query(query, userID).WithContext(ctx).Iter()
	defer iter.Close()

	var conversations []*domain.Conversation
	for {
		var conv domain.Conversation
		var participantsJSON []byte

		if !iter.Scan(
			&conv.ConversationID, &conv.Type, &participantsJSON, &conv.CreatedAt, &conv.UpdatedAt, &conv.LastMessageAt,
		) {
			break
		}

		json.Unmarshal(participantsJSON, &conv.Participants)
		conversations = append(conversations, &conv)
	}

	if err := iter.Close(); err != nil {
		return nil, err
	}
	return conversations, nil
}

func (r *ChatRepository) GetUnreadCount(ctx context.Context, conversationID, userID string) (int, error) {
	query := `SELECT COUNT(*) FROM messages_by_conversation 
			  WHERE conversation_id = ? AND sender_id != ? AND status != 'READ'`

	var count int
	if err := r.session.Query(query, conversationID, userID).WithContext(ctx).Scan(&count); err != nil {
		return 0, err
	}
	return count, nil
}

func (r *ChatRepository) Close() error {
	r.session.Close()
	return nil
}
