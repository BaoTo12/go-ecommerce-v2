package application

import (
	"context"

	"github.com/titan-commerce/backend/notification-service/internal/domain"
	"github.com/titan-commerce/backend/pkg/logger"
)

type NotificationRepository interface {
	Save(ctx context.Context, notification *domain.Notification) error
	FindByID(ctx context.Context, notificationID string) (*domain.Notification, error)
	FindByUserID(ctx context.Context, userID string, pageSize int) ([]*domain.Notification, error)
	Update(ctx context.Context, notification *domain.Notification) error
}

// NotificationSender handles sending notifications via different channels
type NotificationSender interface {
	SendEmail(ctx context.Context, userID, title, message string) error
	SendSMS(ctx context.Context, userID, message string) error
	SendPush(ctx context.Context, userID, title, message string) error
}

type NotificationService struct {
	repo   NotificationRepository
	sender NotificationSender
	logger *logger.Logger
}

func NewNotificationService(repo NotificationRepository, sender NotificationSender, logger *logger.Logger) *NotificationService {
	return &NotificationService{
		repo:   repo,
		sender: sender,
		logger: logger,
	}
}

// SendNotification sends notification via specified channels (Command)
func (s *NotificationService) SendNotification(ctx context.Context, userID string, notifType domain.NotificationType, title, message string, channels []domain.NotificationChannel) (string, error) {
	// Create notification record
	notification := domain.NewNotification(userID, notifType, title, message)

	if err := s.repo.Save(ctx, notification); err != nil {
		s.logger.Error(err, "failed to save notification")
		return "", err
	}

	// Send via requested channels
	for _, channel := range channels {
		switch channel {
		case domain.NotificationChannelEmail:
			if err := s.sender.SendEmail(ctx, userID, title, message); err != nil {
				s.logger.Error(err, "failed to send email")
			}
		case domain.NotificationChannelSMS:
			if err := s.sender.SendSMS(ctx, userID, message); err != nil {
				s.logger.Error(err, "failed to send SMS")
			}
		case domain.NotificationChannelPush:
			if err := s.sender.SendPush(ctx, userID, title, message); err != nil {
				s.logger.Error(err, "failed to send push notification")
			}
		case domain.NotificationChannelInApp:
			// In-app notification already saved to database
		}
	}

	s.logger.Infof("Notification sent: user=%s, type=%s, channels=%d", userID, notifType, len(channels))
	return notification.ID, nil
}

// GetNotifications retrieves user's notifications (Query)
func (s *NotificationService) GetNotifications(ctx context.Context, userID string, pageSize int) ([]*domain.Notification, error) {
	return s.repo.FindByUserID(ctx, userID, pageSize)
}

// MarkAsRead marks notification as read (Command)
func (s *NotificationService) MarkAsRead(ctx context.Context, notificationID string) error {
	notification, err := s.repo.FindByID(ctx, notificationID)
	if err != nil {
		return err
	}

	notification.MarkAsRead()

	if err := s.repo.Update(ctx, notification); err != nil {
		return err
	}

	return nil
}
