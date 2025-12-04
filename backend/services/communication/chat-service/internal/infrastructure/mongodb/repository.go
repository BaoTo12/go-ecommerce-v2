package mongodb

import (
	"context"
	"time"

	"github.com/titan-commerce/backend/chat-service/internal/domain"
	"github.com/titan-commerce/backend/pkg/errors"
	"github.com/titan-commerce/backend/pkg/logger"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type ChatRepository struct {
	db     *mongo.Database
	logger *logger.Logger
}

func NewChatRepository(mongoURI, database string, logger *logger.Logger) (*ChatRepository, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(mongoURI))
	if err != nil {
		return nil, errors.Wrap(errors.ErrInternal, "failed to connect to MongoDB", err)
	}

	if err := client.Ping(ctx, nil); err != nil {
		return nil, errors.Wrap(errors.ErrInternal, "failed to ping MongoDB", err)
	}

	db := client.Database(database)

	// Create indexes
	msgColl := db.Collection("messages")
	msgColl.Indexes().CreateOne(ctx, mongo.IndexModel{
		Keys: bson.D{{Key: "conversation_id", Value: 1}, {Key: "created_at", Value: -1}},
	})

	convColl := db.Collection("conversations")
	convColl.Indexes().CreateOne(ctx, mongo.IndexModel{
		Keys: bson.D{{Key: "participants", Value: 1}},
	})

	logger.Info("Chat MongoDB repository initialized")
	return &ChatRepository{db: db, logger: logger}, nil
}

func (r *ChatRepository) SaveMessage(ctx context.Context, message *domain.Message) error {
	_, err := r.db.Collection("messages").InsertOne(ctx, message)
	if err != nil {
		return errors.Wrap(errors.ErrInternal, "failed to save message", err)
	}
	return nil
}

func (r *ChatRepository) FindMessagesByConversation(ctx context.Context, conversationID string, limit, offset int) ([]*domain.Message, error) {
	opts := options.Find().
		SetSort(bson.D{{Key: "created_at", Value: -1}}).
		SetLimit(int64(limit)).
		SetSkip(int64(offset))

	cursor, err := r.db.Collection("messages").Find(ctx, bson.M{"conversation_id": conversationID}, opts)
	if err != nil {
		return nil, errors.Wrap(errors.ErrInternal, "failed to find messages", err)
	}
	defer cursor.Close(ctx)

	var messages []*domain.Message
	if err := cursor.All(ctx, &messages); err != nil {
		return nil, errors.Wrap(errors.ErrInternal, "failed to decode messages", err)
	}

	return messages, nil
}

func (r *ChatRepository) SaveConversation(ctx context.Context, conversation *domain.Conversation) error {
	_, err := r.db.Collection("conversations").InsertOne(ctx, conversation)
	if err != nil {
		return errors.Wrap(errors.ErrInternal, "failed to save conversation", err)
	}
	return nil
}

func (r *ChatRepository) FindConversationByID(ctx context.Context, conversationID string) (*domain.Conversation, error) {
	var conv domain.Conversation
	err := r.db.Collection("conversations").FindOne(ctx, bson.M{"id": conversationID}).Decode(&conv)
	if err == mongo.ErrNoDocuments {
		return nil, errors.New(errors.ErrNotFound, "conversation not found")
	}
	if err != nil {
		return nil, errors.Wrap(errors.ErrInternal, "failed to find conversation", err)
	}
	return &conv, nil
}

func (r *ChatRepository) FindConversationsByUser(ctx context.Context, userID string) ([]*domain.Conversation, error) {
	cursor, err := r.db.Collection("conversations").Find(ctx, bson.M{"participants": userID})
	if err != nil {
		return nil, errors.Wrap(errors.ErrInternal, "failed to find conversations", err)
	}
	defer cursor.Close(ctx)

	var conversations []*domain.Conversation
	if err := cursor.All(ctx, &conversations); err != nil {
		return nil, errors.Wrap(errors.ErrInternal, "failed to decode conversations", err)
	}
	return conversations, nil
}

func (r *ChatRepository) FindDirectConversation(ctx context.Context, user1ID, user2ID string) (*domain.Conversation, error) {
	filter := bson.M{
		"type": "direct",
		"participants": bson.M{
			"$all": []string{user1ID, user2ID},
		},
	}

	var conv domain.Conversation
	err := r.db.Collection("conversations").FindOne(ctx, filter).Decode(&conv)
	if err == mongo.ErrNoDocuments {
		return nil, errors.New(errors.ErrNotFound, "conversation not found")
	}
	if err != nil {
		return nil, errors.Wrap(errors.ErrInternal, "failed to find conversation", err)
	}
	return &conv, nil
}

func (r *ChatRepository) UpdateConversation(ctx context.Context, conversation *domain.Conversation) error {
	conversation.UpdatedAt = time.Now()
	_, err := r.db.Collection("conversations").ReplaceOne(ctx, bson.M{"id": conversation.ID}, conversation)
	if err != nil {
		return errors.Wrap(errors.ErrInternal, "failed to update conversation", err)
	}
	return nil
}
