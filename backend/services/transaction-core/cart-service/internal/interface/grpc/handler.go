package handler

import (
	"context"

	"github.com/titan-commerce/backend/cart-service/internal/application"
	pb "github.com/titan-commerce/backend/cart-service/proto/cart/v1"
	"github.com/titan-commerce/backend/pkg/logger"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type CartServiceServer struct {
	pb.UnimplementedCartServiceServer
	service *application.CartService
	logger  *logger.Logger
}

func NewCartServiceServer(service *application.CartService, logger *logger.Logger) *CartServiceServer {
	return &CartServiceServer{
		service: service,
		logger:  logger,
	}
}

func (s *CartServiceServer) AddItem(ctx context.Context, req *pb.AddItemRequest) (*pb.AddItemResponse, error) {
	cart, err := s.service.AddItem(ctx, req.UserId, req.ProductId, req.ProductName, req.Price, int(req.Quantity))
	if err != nil {
		s.logger.Error(err, "failed to add item to cart")
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &pb.AddItemResponse{
		Cart: domainToProto(cart),
	}, nil
}

func (s *CartServiceServer) RemoveItem(ctx context.Context, req *pb.RemoveItemRequest) (*pb.RemoveItemResponse, error) {
	cart, err := s.service.RemoveItem(ctx, req.UserId, req.ProductId)
	if err != nil {
		s.logger.Error(err, "failed to remove item from cart")
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &pb.RemoveItemResponse{
		Cart: domainToProto(cart),
	}, nil
}

func (s *CartServiceServer) GetCart(ctx context.Context, req *pb.GetCartRequest) (*pb.GetCartResponse, error) {
	cart, err := s.service.GetCart(ctx, req.UserId)
	if err != nil {
		s.logger.Error(err, "failed to get cart")
		return nil, status.Error(codes.NotFound, err.Error())
	}

	return &pb.GetCartResponse{
		Cart: domainToProto(cart),
	}, nil
}

func (s *CartServiceServer) ClearCart(ctx context.Context, req *pb.ClearCartRequest) (*pb.ClearCartResponse, error) {
	if err := s.service.ClearCart(ctx, req.UserId); err != nil {
		s.logger.Error(err, "failed to clear cart")
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &pb.ClearCartResponse{Success: true}, nil
}

func domainToProto(cart *domain.Cart) *pb.Cart {
	items := make([]*pb.CartItem, len(cart.Items))
	for i, item := range cart.Items {
		items[i] = &pb.CartItem{
			ProductId:   item.ProductID,
			ProductName: item.ProductName,
			Price:       item.Price,
			Quantity:    int32(item.Quantity),
		}
	}

	return &pb.Cart{
		UserId:      cart.UserID,
		Items:       items,
		TotalAmount: cart.TotalAmount,
	}
}
