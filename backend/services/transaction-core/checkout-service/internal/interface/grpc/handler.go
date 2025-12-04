package grpc

import (
	"context"

	"github.com/titan-commerce/backend/checkout-service/internal/application"
	"github.com/titan-commerce/backend/checkout-service/internal/domain"
	pb "github.com/titan-commerce/backend/checkout-service/proto/checkout/v1"
	"github.com/titan-commerce/backend/pkg/logger"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type CheckoutServiceServer struct {
	pb.UnimplementedCheckoutServiceServer
	service *application.CheckoutService
	logger  *logger.Logger
}

func NewCheckoutServiceServer(service *application.CheckoutService, logger *logger.Logger) *CheckoutServiceServer {
	return &CheckoutServiceServer{
		service: service,
		logger:  logger,
	}
}

func (s *CheckoutServiceServer) InitiateCheckout(ctx context.Context, req *pb.InitiateCheckoutRequest) (*pb.InitiateCheckoutResponse, error) {
	session, err := s.service.InitiateCheckout(ctx, req.UserId, req.ShippingAddress, req.PaymentMethodId)
	if err != nil {
		s.logger.Error(err, "failed to initiate checkout")
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &pb.InitiateCheckoutResponse{
		Session: domainToProto(session),
	}, nil
}

func (s *CheckoutServiceServer) GetCheckoutStatus(ctx context.Context, req *pb.GetCheckoutStatusRequest) (*pb.GetCheckoutStatusResponse, error) {
	session, err := s.service.GetCheckoutStatus(ctx, req.SessionId)
	if err != nil {
		return nil, status.Error(codes.NotFound, err.Error())
	}

	return &pb.GetCheckoutStatusResponse{
		Session: domainToProto(session),
	}, nil
}

func (s *CheckoutServiceServer) CancelCheckout(ctx context.Context, req *pb.CancelCheckoutRequest) (*pb.CancelCheckoutResponse, error) {
	if err := s.service.CancelCheckout(ctx, req.SessionId); err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	return &pb.CancelCheckoutResponse{Success: true}, nil
}

func domainToProto(session *domain.CheckoutSession) *pb.CheckoutSession {
	// Map domain status to proto status
	var status pb.CheckoutStatus
	switch session.Status {
	case domain.CheckoutStatusInitiated:
		status = pb.CheckoutStatus_CHECKOUT_STATUS_INITIATED
	case domain.CheckoutStatusReservingInventory:
		status = pb.CheckoutStatus_CHECKOUT_STATUS_RESERVING_INVENTORY
	case domain.CheckoutStatusProcessingPayment:
		status = pb.CheckoutStatus_CHECKOUT_STATUS_PROCESSING_PAYMENT
	case domain.CheckoutStatusCreatingOrder:
		status = pb.CheckoutStatus_CHECKOUT_STATUS_CREATING_ORDER
	case domain.CheckoutStatusCompleted:
		status = pb.CheckoutStatus_CHECKOUT_STATUS_COMPLETED
	case domain.CheckoutStatusFailed:
		status = pb.CheckoutStatus_CHECKOUT_STATUS_FAILED
	case domain.CheckoutStatusCompensating:
		status = pb.CheckoutStatus_CHECKOUT_STATUS_COMPENSATING
	default:
		status = pb.CheckoutStatus_CHECKOUT_STATUS_UNSPECIFIED
	}

	return &pb.CheckoutSession{
		SessionId:    session.SessionID,
		UserId:       session.UserID,
		ProductIds:   session.ProductIDs,
		TotalAmount:  session.TotalAmount,
		Status:       status,
		ErrorMessage: session.ErrorMessage,
		OrderId:      session.OrderID,
		PaymentId:    session.PaymentID,
		// Timestamps omitted for brevity
	}
}
