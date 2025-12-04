package mongodb

import (
	"context"
	"time"

	"github.com/titan-commerce/backend/livestream-service/internal/domain"
	"github.com/titan-commerce/backend/pkg/errors"
	"github.com/titan-commerce/backend/pkg/logger"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type LivestreamRepository struct {
	db     *mongo.Database
	logger *logger.Logger
}

func NewLivestreamRepository(mongoURI, database string, logger *logger.Logger) (*LivestreamRepository, error) {
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
	streamColl := db.Collection("streams")
	streamColl.Indexes().CreateMany(ctx, []mongo.IndexModel{
		{Keys: bson.D{{Key: "stream_key", Value: 1}}},
		{Keys: bson.D{{Key: "status", Value: 1}}},
		{Keys: bson.D{{Key: "seller_id", Value: 1}}},
	})

	logger.Info("Livestream MongoDB repository initialized")
	return &LivestreamRepository{db: db, logger: logger}, nil
}

func (r *LivestreamRepository) Save(ctx context.Context, stream *domain.Livestream) error {
	_, err := r.db.Collection("streams").InsertOne(ctx, stream)
	if err != nil {
		return errors.Wrap(errors.ErrInternal, "failed to save stream", err)
	}
	return nil
}

func (r *LivestreamRepository) FindByID(ctx context.Context, streamID string) (*domain.Livestream, error) {
	var stream domain.Livestream
	err := r.db.Collection("streams").FindOne(ctx, bson.M{"id": streamID}).Decode(&stream)
	if err == mongo.ErrNoDocuments {
		return nil, errors.New(errors.ErrNotFound, "stream not found")
	}
	if err != nil {
		return nil, errors.Wrap(errors.ErrInternal, "failed to find stream", err)
	}
	return &stream, nil
}

func (r *LivestreamRepository) FindByStreamKey(ctx context.Context, streamKey string) (*domain.Livestream, error) {
	var stream domain.Livestream
	err := r.db.Collection("streams").FindOne(ctx, bson.M{"streamkey": streamKey}).Decode(&stream)
	if err == mongo.ErrNoDocuments {
		return nil, errors.New(errors.ErrNotFound, "stream not found")
	}
	if err != nil {
		return nil, errors.Wrap(errors.ErrInternal, "failed to find stream", err)
	}
	return &stream, nil
}

func (r *LivestreamRepository) FindLiveStreams(ctx context.Context, limit, offset int) ([]*domain.Livestream, error) {
	opts := options.Find().
		SetSort(bson.D{{Key: "viewer_count", Value: -1}}).
		SetLimit(int64(limit)).
		SetSkip(int64(offset))

	cursor, err := r.db.Collection("streams").Find(ctx, bson.M{"status": "LIVE"}, opts)
	if err != nil {
		return nil, errors.Wrap(errors.ErrInternal, "failed to find streams", err)
	}
	defer cursor.Close(ctx)

	var streams []*domain.Livestream
	if err := cursor.All(ctx, &streams); err != nil {
		return nil, errors.Wrap(errors.ErrInternal, "failed to decode streams", err)
	}

	return streams, nil
}

func (r *LivestreamRepository) Update(ctx context.Context, stream *domain.Livestream) error {
	_, err := r.db.Collection("streams").ReplaceOne(ctx, bson.M{"id": stream.ID}, stream)
	if err != nil {
		return errors.Wrap(errors.ErrInternal, "failed to update stream", err)
	}
	return nil
}

func (r *LivestreamRepository) SaveChat(ctx context.Context, chat *domain.StreamChat) error {
	_, err := r.db.Collection("stream_chats").InsertOne(ctx, chat)
	return err
}

func (r *LivestreamRepository) GetRecentChats(ctx context.Context, streamID string, limit int) ([]*domain.StreamChat, error) {
	opts := options.Find().
		SetSort(bson.D{{Key: "created_at", Value: -1}}).
		SetLimit(int64(limit))

	cursor, err := r.db.Collection("stream_chats").Find(ctx, bson.M{"stream_id": streamID}, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var chats []*domain.StreamChat
	cursor.All(ctx, &chats)
	return chats, nil
}
