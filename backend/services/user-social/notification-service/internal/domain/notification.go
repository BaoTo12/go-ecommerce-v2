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
	NotificationChannelEmail  NotificationChannel = "EMAIL"
	NotificationChannelSMS    NotificationChannel = "SMS"
	NotificationChannelPush   NotificationChannel = "PUSH"
	NotificationChannelInApp  NotificationChannel = "IN_APP"
)

type Notification struct {
	ID        string
	UserID    string
	Type      NotificationType
	Title     string
	Message   string
	Read      bool
	CreatedAt time.Time
}

func NewNotification(userID string, notifType NotificationType, title, message string) *Notification {
	return &Notification{
		ID:        uuid.New().String(),
		UserID:    userID,
		Type:      notifType,
		Title:     title,
		Message:   message,
		Read:      false,
		CreatedAt: time.Now(),
	}
}

func (n *Notification) MarkAsRead() {
	n.Read = true
}
