package postgres

import (
	"context"

	"github.com/titan-commerce/backend/videocall-service/internal/domain"
)

type VideoCallRepository struct {
	// In production: use *sql.DB
}

func NewVideoCallRepository() *VideoCallRepository {
	return &VideoCallRepository{}
}

func (r *VideoCallRepository) Save(ctx context.Context, call *domain.VideoCall) error {
	return nil
}

func (r *VideoCallRepository) FindByID(ctx context.Context, callID string) (*domain.VideoCall, error) {
	return &domain.VideoCall{ID: callID}, nil
}

func (r *VideoCallRepository) FindByRoomID(ctx context.Context, roomID string) (*domain.VideoCall, error) {
	return nil, nil
}

func (r *VideoCallRepository) Update(ctx context.Context, call *domain.VideoCall) error {
	return nil
}

func (r *VideoCallRepository) FindCallHistory(ctx context.Context, userID string, limit int) ([]*domain.VideoCall, error) {
	return nil, nil
}
