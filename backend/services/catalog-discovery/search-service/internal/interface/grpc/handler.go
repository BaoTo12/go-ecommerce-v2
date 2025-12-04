package grpc

import (
	"context"

	"github.com/titan-commerce/backend/search-service/internal/application"
	"github.com/titan-commerce/backend/search-service/internal/domain"
	pb "github.com/titan-commerce/backend/search-service/proto/search/v1"
	"github.com/titan-commerce/backend/pkg/logger"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type SearchServiceServer struct {
	pb.UnimplementedSearchServiceServer
	service *application.SearchService
	logger  *logger.Logger
}

func NewSearchServiceServer(service *application.SearchService, logger *logger.Logger) *SearchServiceServer {
	return &SearchServiceServer{
		service: service,
		logger:  logger,
	}
}

func (s *SearchServiceServer) IndexProduct(ctx context.Context, req *pb.IndexProductRequest) (*pb.IndexProductResponse, error) {
	doc := &domain.ProductDocument{
		ID:          req.Id,
		Name:        req.Name,
		Description: req.Description,
		Category:    req.Category,
		Price:       req.Price,
		ImageURL:    req.ImageUrl,
	}

	if err := s.service.IndexProduct(ctx, doc); err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &pb.IndexProductResponse{Success: true}, nil
}

func (s *SearchServiceServer) SearchProducts(ctx context.Context, req *pb.SearchProductsRequest) (*pb.SearchProductsResponse, error) {
	results, total, err := s.service.SearchProducts(ctx, req.Query, int(req.Page), int(req.PageSize))
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	var protoResults []*pb.ProductResult
	for _, r := range results {
		protoResults = append(protoResults, &pb.ProductResult{
			Id:          r.ID,
			Name:        r.Name,
			Description: r.Description,
			Price:       r.Price,
			ImageUrl:    r.ImageURL,
			Score:       1.0, // Placeholder score
		})
	}

	return &pb.SearchProductsResponse{
		Results: protoResults,
		Total:   int32(total),
	}, nil
}
