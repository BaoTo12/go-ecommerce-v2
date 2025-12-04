package grpc

import (
	"context"

	"github.com/titan-commerce/backend/product-service/internal/application"
	"github.com/titan-commerce/backend/product-service/internal/domain"
	pb "github.com/titan-commerce/backend/product-service/proto/product/v1"
	"github.com/titan-commerce/backend/pkg/logger"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type ProductServiceServer struct {
	pb.UnimplementedProductServiceServer
	service *application.ProductService
	logger  *logger.Logger
}

func NewProductServiceServer(service *application.ProductService, logger *logger.Logger) *ProductServiceServer {
	return &ProductServiceServer{
		service: service,
		logger:  logger,
	}
}

func (s *ProductServiceServer) CreateProduct(ctx context.Context, req *pb.CreateProductRequest) (*pb.CreateProductResponse, error) {
	product := &domain.Product{
		Name:        req.Name,
		Description: req.Description,
		Price:       req.Price,
		Currency:    req.Currency,
		CategoryID:  req.CategoryId,
		Images:      req.Images,
		Attributes:  req.Attributes,
		Stock:       int(req.Stock),
	}

	created, err := s.service.CreateProduct(ctx, product)
	if err != nil {
		s.logger.Error(err, "failed to create product")
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &pb.CreateProductResponse{
		Product: domainToProto(created),
	}, nil
}

func (s *ProductServiceServer) GetProduct(ctx context.Context, req *pb.GetProductRequest) (*pb.GetProductResponse, error) {
	product, err := s.service.GetProduct(ctx, req.Id)
	if err != nil {
		return nil, status.Error(codes.NotFound, err.Error())
	}

	return &pb.GetProductResponse{
		Product: domainToProto(product),
	}, nil
}

func (s *ProductServiceServer) ListProducts(ctx context.Context, req *pb.ListProductsRequest) (*pb.ListProductsResponse, error) {
	products, total, err := s.service.ListProducts(ctx, int(req.Page), int(req.PageSize))
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	var protoProducts []*pb.Product
	for _, p := range products {
		protoProducts = append(protoProducts, domainToProto(p))
	}

	return &pb.ListProductsResponse{
		Products: protoProducts,
		Total:    int32(total),
	}, nil
}

func domainToProto(p *domain.Product) *pb.Product {
	return &pb.Product{
		Id:          p.ID,
		Name:        p.Name,
		Description: p.Description,
		Price:       p.Price,
		Currency:    p.Currency,
		CategoryId:  p.CategoryID,
		Images:      p.Images,
		Attributes:  p.Attributes,
		Stock:       int32(p.Stock),
	}
}
