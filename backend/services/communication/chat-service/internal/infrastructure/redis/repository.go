package redis

import (
	"context"
	"encoding/json"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/titan-commerce/backend/chat-service/internal/domain"
	"github.com/titan-commerce/backend/pkg/logger"
)

type PresenceRepository struct {
	client *redis.Client
	logger *logger.Logger
}

func NewPresenceRepository(redisURL string, logger *logger.Logger) (*PresenceRepository, error) {
	opt, err := redis.ParseURL(redisURL)
	if err != nil {
		return nil, err
	}

	client := redis.NewClient(opt)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := client.Ping(ctx).Err(); err != nil {
		return nil, err
	}

	logger.Info("Chat Redis repository initialized")
	return &PresenceRepository{client: client, logger: logger}, nil
}

func (r *PresenceRepository) SetUserOnline(ctx context.Context, userID string, socketID string) error {
	key := "presence:" + userID
	return r.client.Set(ctx, key, socketID, 5*time.Minute).Err()
}

func (r *PresenceRepository) SetUserOffline(ctx context.Context, userID string) error {
	key := "presence:" + userID
	return r.client.Del(ctx, key).Err()
}

func (r *PresenceRepository) IsUserOnline(ctx context.Context, userID string) (bool, error) {
	key := "presence:" + userID
	exists, err := r.client.Exists(ctx, key).Result()
	return exists > 0, err
}

func (r *PresenceRepository) SetTypingIndicator(ctx context.Context, conversationID, userID string) error {
	key := "typing:" + conversationID + ":" + userID
	return r.client.Set(ctx, key, "1", 5*time.Second).Err()
}

func (r *PresenceRepository) GetTypingUsers(ctx context.Context, conversationID string) ([]string, error) {
	pattern := "typing:" + conversationID + ":*"
	keys, err := r.client.Keys(ctx, pattern).Result()
	if err != nil {
		return nil, err
	}

	var users []string
	for _, key := range keys {
		// Extract userID from key
		parts := []rune(key)
		if len(parts) > len("typing:"+conversationID+":") {
			userID := string(parts[len("typing:"+conversationID+":"):])
			users = append(users, userID)
		}
	}
	return users, nil
}

func (r *PresenceRepository) PublishMessage(ctx context.Context, conversationID string, msg *domain.Message) error {
	channel := "conversation:" + conversationID
	msgJSON, _ := json.Marshal(msg)
	return r.client.Publish(ctx, channel, msgJSON).Err()
}

func (r *PresenceRepository) SubscribeToConversation(ctx context.Context, conversationID string) *redis.PubSub {
	channel := "conversation:" + conversationID
	return r.client.Subscribe(ctx, channel)
}

func (r *PresenceRepository) Close() error {
	return r.client.Close()
}
