package grpc

import (
	"context"

	"github.com/titan-commerce/backend/feed-service/internal/application"
	"github.com/titan-commerce/backend/feed-service/internal/domain"
	pb "github.com/titan-commerce/backend/feed-service/proto/feed/v1"
	"github.com/titan-commerce/backend/pkg/logger"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type FeedServiceServer struct {
	pb.UnimplementedFeedServiceServer
	service *application.FeedService
	logger  *logger.Logger
}

func NewFeedServiceServer(service *application.FeedService, logger *logger.Logger) *FeedServiceServer {
	return &FeedServiceServer{
		service: service,
		logger:  logger,
	}
}

func (s *FeedServiceServer) PublishPost(ctx context.Context, req *pb.PublishPostRequest) (*pb.PublishPostResponse, error) {
	post, err := s.service.PublishPost(ctx, req.UserId, req.Content, req.MediaUrl, req.Tags)
	if err != nil {
		s.logger.Error(err, "failed to publish post")
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &pb.PublishPostResponse{
		PostId: post.ID,
		Item:   domainToProto(post),
	}, nil
}

func (s *FeedServiceServer) DeletePost(ctx context.Context, req *pb.DeletePostRequest) (*pb.DeletePostResponse, error) {
	if err := s.service.DeletePost(ctx, req.PostId, req.UserId); err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	return &pb.DeletePostResponse{Success: true}, nil
}

func (s *FeedServiceServer) GetFeed(ctx context.Context, req *pb.GetFeedRequest) (*pb.GetFeedResponse, error) {
	posts, err := s.service.GetFeed(ctx, req.UserId, int(req.Page), int(req.PageSize))
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	var items []*pb.FeedItem
	for _, p := range posts {
		items = append(items, domainToProto(p))
	}

	return &pb.GetFeedResponse{
		Items:      items,
		NextCursor: "", // Pagination cursor logic omitted for MVP
	}, nil
}

func domainToProto(p *domain.Post) *pb.FeedItem {
	return &pb.FeedItem{
		PostId:        p.ID,
		UserId:        p.UserID,
		Content:       p.Content,
		MediaUrl:      p.MediaURL,
		LikesCount:    int32(p.LikesCount),
		CommentsCount: int32(p.CommentsCount),
		Tags:          p.Tags,
		CreatedAt:     timestamppb.New(p.CreatedAt),
	}
}
