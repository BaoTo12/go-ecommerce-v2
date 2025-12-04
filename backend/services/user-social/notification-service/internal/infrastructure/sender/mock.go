package sender

import (
	"context"

	"github.com/titan-commerce/backend/pkg/logger"
)

type MockNotificationSender struct {
	logger *logger.Logger
}

func NewMockNotificationSender(logger *logger.Logger) *MockNotificationSender {
	return &MockNotificationSender{logger: logger}
}

func (s *MockNotificationSender) SendEmail(ctx context.Context, userID, title, message string) error {
	s.logger.Infof("MOCK EMAIL to User %s: %s - %s", userID, title, message)
	return nil
}

func (s *MockNotificationSender) SendSMS(ctx context.Context, userID, message string) error {
	s.logger.Infof("MOCK SMS to User %s: %s", userID, message)
	return nil
}

func (s *MockNotificationSender) SendPush(ctx context.Context, userID, title, message string) error {
	s.logger.Infof("MOCK PUSH to User %s: %s - %s", userID, title, message)
	return nil
}
