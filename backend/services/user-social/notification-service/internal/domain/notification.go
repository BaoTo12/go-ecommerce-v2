package domain

import (
	"time"

	"github.com/google/uuid"
)

type NotificationType string

const (
	NotificationTypeOrderPlaced    NotificationType = "ORDER_PLACED"
	NotificationTypePaymentSuccess NotificationType = "PAYMENT_SUCCESS"
	NotificationTypeShipmentUpdate NotificationType = "SHIPMENT_UPDATE"
	NotificationTypeFlashSaleAlert NotificationType = "FLASH_SALE_ALERT"
	NotificationTypeCoinReward     NotificationType = "COIN_REWARD"
)

type NotificationChannel string

const (
	NotificationChannelEmail NotificationChannel = "EMAIL"
	NotificationChannelSMS   NotificationChannel = "SMS"
	NotificationChannelPush  NotificationChannel = "PUSH"
	NotificationChannelInApp NotificationChannel = "IN_APP"
)

type Notification struct {
	ID        string
	UserID    string
	Type      NotificationType
	Channel   NotificationChannel
	Title     string
	Content   string
	Message   string // Alias for Content
	Read      bool
	CreatedAt time.Time
	SentAt    time.Time
}

func NewNotification(userID string, notifType NotificationType, title, message string) *Notification {
	now := time.Now()
	return &Notification{
		ID:        uuid.New().String(),
		UserID:    userID,
		Type:      notifType,
		Title:     title,
		Content:   message,
		Message:   message,
		Read:      false,
		CreatedAt: now,
		SentAt:    now,
	}
}

func NewNotificationWithChannel(userID string, notifType NotificationType, channel NotificationChannel, title, content string) *Notification {
	now := time.Now()
	return &Notification{
		ID:        uuid.New().String(),
		UserID:    userID,
		Type:      notifType,
		Channel:   channel,
		Title:     title,
		Content:   content,
		Message:   content,
		Read:      false,
		CreatedAt: now,
		SentAt:    now,
	}
}

func (n *Notification) MarkAsRead() {
	n.Read = true
}
