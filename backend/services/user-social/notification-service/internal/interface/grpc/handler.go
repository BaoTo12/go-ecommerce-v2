package grpc

import (
	"context"

	"github.com/titan-commerce/backend/notification-service/internal/application"
	"github.com/titan-commerce/backend/notification-service/internal/domain"
	pb "github.com/titan-commerce/backend/notification-service/proto/notification/v1"
	"github.com/titan-commerce/backend/pkg/logger"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type NotificationServiceServer struct {
	pb.UnimplementedNotificationServiceServer
	service *application.NotificationService
	logger  *logger.Logger
}

func NewNotificationServiceServer(service *application.NotificationService, logger *logger.Logger) *NotificationServiceServer {
	return &NotificationServiceServer{
		service: service,
		logger:  logger,
	}
}

func (s *NotificationServiceServer) SendNotification(ctx context.Context, req *pb.SendNotificationRequest) (*pb.SendNotificationResponse, error) {
	// Convert single channel to slice for service
	channels := []domain.NotificationChannel{domain.NotificationChannel(req.Channel)}

	notificationID, err := s.service.SendNotification(
		ctx,
		req.UserId,
		domain.NotificationType(req.Type),
		req.Title,
		req.Content,
		channels,
	)
	if err != nil {
		s.logger.Error(err, "failed to send notification")
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &pb.SendNotificationResponse{
		NotificationId: notificationID,
		Success:        true,
	}, nil
}

func (s *NotificationServiceServer) GetNotifications(ctx context.Context, req *pb.GetNotificationsRequest) (*pb.GetNotificationsResponse, error) {
	notifications, err := s.service.GetNotifications(ctx, req.UserId, int(req.PageSize))
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	var protoNotifications []*pb.Notification
	for _, n := range notifications {
		protoNotifications = append(protoNotifications, domainToProto(n))
	}

	return &pb.GetNotificationsResponse{
		Notifications: protoNotifications,
		Total:         int32(len(notifications)),
	}, nil
}

func (s *NotificationServiceServer) MarkAsRead(ctx context.Context, req *pb.MarkAsReadRequest) (*pb.MarkAsReadResponse, error) {
	if err := s.service.MarkAsRead(ctx, req.NotificationId); err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	return &pb.MarkAsReadResponse{Success: true}, nil
}

func domainToProto(n *domain.Notification) *pb.Notification {
	return &pb.Notification{
		Id:        n.ID,
		UserId:    n.UserID,
		Type:      string(n.Type),
		Channel:   string(n.Channel),
		Title:     n.Title,
		Content:   n.Content,
		Read:      n.Read,
		CreatedAt: n.CreatedAt.String(),
	}
}
