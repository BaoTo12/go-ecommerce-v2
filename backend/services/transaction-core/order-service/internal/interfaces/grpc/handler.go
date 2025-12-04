package grpc

import (
	"context"

	"github.com/titan-commerce/backend/order-service/internal/application"
	"github.com/titan-commerce/backend/order-service/internal/domain"
	pb "github.com/titan-commerce/backend/order-service/proto/order/v1"
	"github.com/titan-commerce/backend/pkg/logger"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type OrderServiceServer struct {
	pb.UnimplementedOrderServiceServer
	service *application.OrderService
	logger  *logger.Logger
}

func NewOrderServiceServer(server *grpc.Server, service *application.OrderService, logger *logger.Logger) *OrderServiceServer {
	handler := &OrderServiceServer{
		service: service,
		logger:  logger,
	}
	pb.RegisterOrderServiceServer(server, handler)
	return handler
}

func (s *OrderServiceServer) CreateOrder(ctx context.Context, req *pb.CreateOrderRequest) (*pb.CreateOrderResponse, error) {
	items := make([]domain.OrderItem, len(req.Items))
	for i, item := range req.Items {
		items[i] = domain.OrderItem{
			ProductID:   item.ProductId,
			ProductName: item.ProductName,
			Quantity:    int(item.Quantity),
			UnitPrice:   item.Price,
		}
	}

	order, err := s.service.CreateOrder(ctx, req.UserId, items, req.ShippingAddress)
	if err != nil {
		s.logger.Error(err, "failed to create order")
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &pb.CreateOrderResponse{
		Order: domainToProto(order),
	}, nil
}

func (s *OrderServiceServer) GetOrder(ctx context.Context, req *pb.GetOrderRequest) (*pb.GetOrderResponse, error) {
	order, err := s.service.GetOrder(ctx, req.OrderId)
	if err != nil {
		return nil, status.Error(codes.NotFound, err.Error())
	}

	return &pb.GetOrderResponse{
		Order: domainToProto(order),
	}, nil
}

func (s *OrderServiceServer) ListOrders(ctx context.Context, req *pb.ListOrdersRequest) (*pb.ListOrdersResponse, error) {
	return &pb.ListOrdersResponse{}, nil
}

func (s *OrderServiceServer) CancelOrder(ctx context.Context, req *pb.CancelOrderRequest) (*pb.CancelOrderResponse, error) {
	_, err := s.service.CancelOrder(ctx, req.OrderId, req.Reason)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	return &pb.CancelOrderResponse{Success: true}, nil
}

func (s *OrderServiceServer) UpdateOrderStatus(ctx context.Context, req *pb.UpdateOrderStatusRequest) (*pb.UpdateOrderStatusResponse, error) {
	return &pb.UpdateOrderStatusResponse{}, nil
}

func domainToProto(order *domain.Order) *pb.Order {
	items := make([]*pb.OrderItem, len(order.Items))
	for i, item := range order.Items {
		items[i] = &pb.OrderItem{
			ProductId:   item.ProductID,
			ProductName: item.ProductName,
			Quantity:    int32(item.Quantity),
			Price:       item.UnitPrice,
		}
	}

	return &pb.Order{
		OrderId:         order.ID,
		UserId:          order.UserID,
		Items:           items,
		TotalAmount:     order.TotalAmount,
		Status:          string(order.Status),
		ShippingAddress: order.ShippingAddress,
		CreatedAt:       nil,
	}
}
