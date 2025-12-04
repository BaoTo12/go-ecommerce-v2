package grpc

import (
	"context"

	"github.com/titan-commerce/backend/category-service/internal/application"
	"github.com/titan-commerce/backend/category-service/internal/domain"
	pb "github.com/titan-commerce/backend/category-service/proto/category/v1"
	"github.com/titan-commerce/backend/pkg/logger"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type CategoryServiceServer struct {
	pb.UnimplementedCategoryServiceServer
	service *application.CategoryService
	logger  *logger.Logger
}

func NewCategoryServiceServer(service *application.CategoryService, logger *logger.Logger) *CategoryServiceServer {
	return &CategoryServiceServer{
		service: service,
		logger:  logger,
	}
}

func (s *CategoryServiceServer) CreateCategory(ctx context.Context, req *pb.CreateCategoryRequest) (*pb.CreateCategoryResponse, error) {
	category, err := s.service.CreateCategory(ctx, req.Name, req.Description, req.ParentId, req.ImageUrl)
	if err != nil {
		s.logger.Error(err, "failed to create category")
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &pb.CreateCategoryResponse{
		Category: domainToProto(category),
	}, nil
}

func (s *CategoryServiceServer) GetCategory(ctx context.Context, req *pb.GetCategoryRequest) (*pb.GetCategoryResponse, error) {
	category, err := s.service.GetCategory(ctx, req.Id)
	if err != nil {
		return nil, status.Error(codes.NotFound, err.Error())
	}

	return &pb.GetCategoryResponse{
		Category: domainToProto(category),
	}, nil
}

func (s *CategoryServiceServer) ListCategories(ctx context.Context, req *pb.ListCategoriesRequest) (*pb.ListCategoriesResponse, error) {
	categories, total, err := s.service.ListCategories(ctx, int(req.Page), int(req.PageSize))
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	var protoCategories []*pb.Category
	for _, c := range categories {
		protoCategories = append(protoCategories, domainToProto(c))
	}

	return &pb.ListCategoriesResponse{
		Categories: protoCategories,
		Total:      int32(total),
	}, nil
}

func (s *CategoryServiceServer) GetCategoryTree(ctx context.Context, req *pb.GetCategoryTreeRequest) (*pb.GetCategoryTreeResponse, error) {
	roots, err := s.service.GetCategoryTree(ctx)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	var protoRoots []*pb.CategoryNode
	for _, r := range roots {
		protoRoots = append(protoRoots, nodeToProto(r))
	}

	return &pb.GetCategoryTreeResponse{
		Roots: protoRoots,
	}, nil
}

func domainToProto(c *domain.Category) *pb.Category {
	return &pb.Category{
		Id:          c.ID,
		Name:        c.Name,
		Description: c.Description,
		ParentId:    c.ParentID,
		ImageUrl:    c.ImageURL,
	}
}

func nodeToProto(n *domain.CategoryNode) *pb.CategoryNode {
	node := &pb.CategoryNode{
		Category: domainToProto(n.Category),
	}
	for _, child := range n.Children {
		node.Children = append(node.Children, nodeToProto(child))
	}
	return node
}
