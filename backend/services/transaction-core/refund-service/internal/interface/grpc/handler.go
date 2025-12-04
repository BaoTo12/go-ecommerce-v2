package handler

import (
	"context"

	"github.com/titan-commerce/backend/refund-service/internal/application"
	"github.com/titan-commerce/backend/refund-service/internal/domain"
	pb "github.com/titan-commerce/backend/refund-service/proto/refund/v1"
	"github.com/titan-commerce/backend/pkg/logger"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type RefundServiceServer struct {
	pb.UnimplementedRefundServiceServer
	service *application.RefundService
	logger  *logger.Logger
}

func NewRefundServiceServer(service *application.RefundService, logger *logger.Logger) *RefundServiceServer {
	return &RefundServiceServer{
		service: service,
		logger:  logger,
	}
}

func (s *RefundServiceServer) ProcessRefund(ctx context.Context, req *pb.ProcessRefundRequest) (*pb.ProcessRefundResponse, error) {
	refund, err := s.service.ProcessRefund(ctx, req.PaymentId, req.OrderId, req.Amount, req.Reason)
	if err != nil {
		s.logger.Error(err, "failed to process refund")
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &pb.ProcessRefundResponse{
		Refund: domainToProto(refund),
	}, nil
}

func (s *RefundServiceServer) GetRefund(ctx context.Context, req *pb.GetRefundRequest) (*pb.GetRefundResponse, error) {
	refund, err := s.service.GetRefund(ctx, req.RefundId)
	if err != nil {
		return nil, status.Error(codes.NotFound, err.Error())
	}

	return &pb.GetRefundResponse{
		Refund: domainToProto(refund),
	}, nil
}

func domainToProto(refund *domain.Refund) *pb.Refund {
	return &pb.Refund{
		RefundId:        refund.ID,
		PaymentId:       refund.PaymentID,
		OrderId:         refund.OrderID,
		Amount:          refund.Amount,
		Reason:          refund.Reason,
		Status:          string(refund.Status),
		GatewayRefundId: refund.GatewayRefundID,
	}
}
