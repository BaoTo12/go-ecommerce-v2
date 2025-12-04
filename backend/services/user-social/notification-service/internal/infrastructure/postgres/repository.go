package postgres

import (
	"context"
	"database/sql"
	"time"

	"github.com/titan-commerce/backend/notification-service/internal/domain"
	"github.com/titan-commerce/backend/pkg/errors"
	"github.com/titan-commerce/backend/pkg/logger"
	_ "github.com/lib/pq"
)

type NotificationRepository struct {
	db     *sql.DB
	logger *logger.Logger
}

func NewNotificationRepository(databaseURL string, logger *logger.Logger) (*NotificationRepository, error) {
	db, err := sql.Open("postgres", databaseURL)
	if err != nil {
		return nil, errors.Wrap(errors.ErrInternal, "failed to connect to database", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	
	if err := db.PingContext(ctx); err != nil {
		return nil, errors.Wrap(errors.ErrInternal, "failed to ping database", err)
	}

	logger.Info("Notification PostgreSQL repository initialized")
	return &NotificationRepository{db: db, logger: logger}, nil
}

func (r *NotificationRepository) Save(ctx context.Context, notification *domain.Notification) error {
	query := `
		INSERT INTO notifications (
			id, user_id, type, channel, title, content, 
			read, created_at, sent_at
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
	`

	_, err := r.db.ExecContext(ctx, query,
		notification.ID, notification.UserID, notification.Type, notification.Channel,
		notification.Title, notification.Content, notification.Read,
		notification.CreatedAt, notification.SentAt,
	)

	if err != nil {
		return errors.Wrap(errors.ErrInternal, "failed to save notification", err)
	}

	return nil
}

func (r *NotificationRepository) FindByUserID(ctx context.Context, userID string, page, pageSize int) ([]*domain.Notification, int, error) {
	offset := (page - 1) * pageSize

	query := `
		SELECT id, user_id, type, channel, title, content, read, created_at, sent_at
		FROM notifications
		WHERE user_id = $1
		ORDER BY created_at DESC
		LIMIT $2 OFFSET $3
	`

	rows, err := r.db.QueryContext(ctx, query, userID, pageSize, offset)
	if err != nil {
		return nil, 0, errors.Wrap(errors.ErrInternal, "failed to find notifications", err)
	}
	defer rows.Close()

	var notifications []*domain.Notification
	for rows.Next() {
		var n domain.Notification
		err := rows.Scan(
			&n.ID, &n.UserID, &n.Type, &n.Channel, &n.Title, &n.Content,
			&n.Read, &n.CreatedAt, &n.SentAt,
		)
		if err != nil {
			return nil, 0, errors.Wrap(errors.ErrInternal, "failed to scan notification", err)
		}
		notifications = append(notifications, &n)
	}

	// Get total count
	var total int
	countQuery := `SELECT COUNT(*) FROM notifications WHERE user_id = $1`
	err = r.db.QueryRowContext(ctx, countQuery, userID).Scan(&total)
	if err != nil {
		return nil, 0, errors.Wrap(errors.ErrInternal, "failed to count notifications", err)
	}

	return notifications, total, nil
}

func (r *NotificationRepository) MarkAsRead(ctx context.Context, notificationID string) error {
	query := `UPDATE notifications SET read = true WHERE id = $1`
	_, err := r.db.ExecContext(ctx, query, notificationID)
	if err != nil {
		return errors.Wrap(errors.ErrInternal, "failed to mark notification as read", err)
	}
	return nil
}
