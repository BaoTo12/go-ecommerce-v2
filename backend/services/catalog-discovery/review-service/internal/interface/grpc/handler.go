package grpc

import (
	"context"

	"github.com/titan-commerce/backend/review-service/internal/application"
	"github.com/titan-commerce/backend/review-service/internal/domain"
	pb "github.com/titan-commerce/backend/review-service/proto/review/v1"
	"github.com/titan-commerce/backend/pkg/logger"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type ReviewServiceServer struct {
	pb.UnimplementedReviewServiceServer
	service *application.ReviewService
	logger  *logger.Logger
}

func NewReviewServiceServer(service *application.ReviewService, logger *logger.Logger) *ReviewServiceServer {
	return &ReviewServiceServer{
		service: service,
		logger:  logger,
	}
}

func (s *ReviewServiceServer) CreateReview(ctx context.Context, req *pb.CreateReviewRequest) (*pb.CreateReviewResponse, error) {
	review, err := s.service.CreateReview(ctx, req.UserId, req.ProductId, int(req.Rating), req.Comment, req.Images)
	if err != nil {
		s.logger.Error(err, "failed to create review")
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &pb.CreateReviewResponse{
		ReviewId: review.ID,
		Review:   domainToProto(review),
	}, nil
}

func (s *ReviewServiceServer) GetProductReviews(ctx context.Context, req *pb.GetProductReviewsRequest) (*pb.GetProductReviewsResponse, error) {
	reviews, total, err := s.service.GetProductReviews(ctx, req.ProductId, int(req.Page), int(req.PageSize))
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	var protoReviews []*pb.Review
	for _, r := range reviews {
		protoReviews = append(protoReviews, domainToProto(r))
	}

	return &pb.GetProductReviewsResponse{
		Reviews: protoReviews,
		Total:   int32(total),
	}, nil
}

func (s *ReviewServiceServer) GetReviewStats(ctx context.Context, req *pb.GetReviewStatsRequest) (*pb.GetReviewStatsResponse, error) {
	stats, err := s.service.GetStats(ctx, req.ProductId)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	dist := make(map[int32]int32)
	for k, v := range stats.RatingDistribution {
		dist[int32(k)] = int32(v)
	}

	return &pb.GetReviewStatsResponse{
		AverageRating:      stats.AverageRating,
		TotalReviews:       int32(stats.TotalReviews),
		RatingDistribution: dist,
	}, nil
}

func domainToProto(r *domain.Review) *pb.Review {
	return &pb.Review{
		Id:        r.ID,
		UserId:    r.UserID,
		ProductId: r.ProductID,
		Rating:    int32(r.Rating),
		Comment:   r.Comment,
		Images:    r.Images,
		CreatedAt: timestamppb.New(r.CreatedAt),
	}
}
