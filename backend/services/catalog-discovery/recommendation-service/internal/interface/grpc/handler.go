package grpc

import (
	"context"

	"github.com/titan-commerce/backend/recommendation-service/internal/application"
	pb "github.com/titan-commerce/backend/recommendation-service/proto/recommendation/v1"
	"github.com/titan-commerce/backend/pkg/logger"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type RecommendationServiceServer struct {
	pb.UnimplementedRecommendationServiceServer
	service *application.RecommendationService
	logger  *logger.Logger
}

func NewRecommendationServiceServer(service *application.RecommendationService, logger *logger.Logger) *RecommendationServiceServer {
	return &RecommendationServiceServer{
		service: service,
		logger:  logger,
	}
}

func (s *RecommendationServiceServer) GetRecommendations(ctx context.Context, req *pb.GetRecommendationsRequest) (*pb.GetRecommendationsResponse, error) {
	items, err := s.service.GetRecommendations(ctx, req.UserId, int(req.Limit), req.Context)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	var protoItems []*pb.RecommendedItem
	for _, item := range items {
		protoItems = append(protoItems, &pb.RecommendedItem{
			ProductId: item.ProductID,
			Score:     float32(item.Score),
			Reason:    item.Reason,
		})
	}

	return &pb.GetRecommendationsResponse{
		Items: protoItems,
	}, nil
}

func (s *RecommendationServiceServer) TrackInteraction(ctx context.Context, req *pb.TrackInteractionRequest) (*pb.TrackInteractionResponse, error) {
	if err := s.service.TrackInteraction(ctx, req.UserId, req.ProductId, req.InteractionType); err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	return &pb.TrackInteractionResponse{Success: true}, nil
}
