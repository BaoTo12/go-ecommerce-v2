package mock

import (
	"context"

	"github.com/titan-commerce/backend/notification-service/internal/domain"
	"github.com/titan-commerce/backend/pkg/logger"
)

type MockNotificationSender struct {
	logger *logger.Logger
}

func NewMockNotificationSender(logger *logger.Logger) *MockNotificationSender {
	return &MockNotificationSender{logger: logger}
}

func (s *MockNotificationSender) Send(ctx context.Context, notification *domain.Notification) error {
	s.logger.Infof("MOCK SEND [%s] to User %s: %s - %s", 
		notification.Channel, notification.UserID, notification.Title, notification.Content)
	return nil
}
