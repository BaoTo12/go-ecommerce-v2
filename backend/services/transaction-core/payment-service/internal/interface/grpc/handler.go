package handler

import (
	"context"

	"github.com/titan-commerce/backend/payment-service/internal/application"
	"github.com/titan-commerce/backend/payment-service/internal/domain"
	pb "github.com/titan-commerce/backend/payment-service/proto/payment/v1"
	"github.com/titan-commerce/backend/pkg/logger"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type PaymentServiceServer struct {
	pb.UnimplementedPaymentServiceServer
	service *application.PaymentService
	logger  *logger.Logger
}

func NewPaymentServiceServer(service *application.PaymentService, logger *logger.Logger) *PaymentServiceServer {
	return &PaymentServiceServer{
		service: service,
		logger:  logger,
	}
}

func (s *PaymentServiceServer) ProcessPayment(ctx context.Context, req *pb.ProcessPaymentRequest) (*pb.ProcessPaymentResponse, error) {
	// Convert proto gateway to domain gateway
	gateway := domain.PaymentGateway(req.Gateway)

	payment, clientSecret, err := s.service.ProcessPayment(
		ctx,
		req.OrderId,
		req.UserId,
		req.Amount,
		req.Currency,
		gateway,
		req.PaymentMethodId,
		req.IdempotencyKey,
	)

	if err != nil {
		s.logger.Error(err, "failed to process payment")
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &pb.ProcessPaymentResponse{
		Payment:      domainToProto(payment),
		ClientSecret: clientSecret,
	}, nil
}

func (s *PaymentServiceServer) GetPayment(ctx context.Context, req *pb.GetPaymentRequest) (*pb.GetPaymentResponse, error) {
	payment, err := s.service.GetPayment(ctx, req.PaymentId)
	if err != nil {
		return nil, status.Error(codes.NotFound, err.Error())
	}

	return &pb.GetPaymentResponse{
		Payment: domainToProto(payment),
	}, nil
}

func (s *PaymentServiceServer) RefundPayment(ctx context.Context, req *pb.RefundPaymentRequest) (*pb.RefundPaymentResponse, error) {
	refundID, err := s.service.RefundPayment(ctx, req.PaymentId, req.Amount, req.Reason)
	if err != nil {
		s.logger.Error(err, "failed to refund payment")
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &pb.RefundPaymentResponse{
		RefundId: refundID,
		Success:  true,
	}, nil
}

func domainToProto(payment *domain.Payment) *pb.Payment {
	return &pb.Payment{
		PaymentId:            payment.ID,
		OrderId:              payment.OrderID,
		UserId:               payment.UserID,
		Amount:               payment.Amount,
		Currency:             payment.Currency,
		Gateway:              string(payment.Gateway),
		Status:               string(payment.Status),
		GatewayTransactionId: payment.GatewayTransactionID,
	}
}
