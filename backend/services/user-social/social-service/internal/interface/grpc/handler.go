package grpc

import (
	"context"

	"github.com/titan-commerce/backend/social-service/internal/application"
	"github.com/titan-commerce/backend/social-service/internal/domain"
	pb "github.com/titan-commerce/backend/social-service/proto/social/v1"
	"github.com/titan-commerce/backend/pkg/logger"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type SocialServiceServer struct {
	pb.UnimplementedSocialServiceServer
	service *application.SocialService
	logger  *logger.Logger
}

func NewSocialServiceServer(service *application.SocialService, logger *logger.Logger) *SocialServiceServer {
	return &SocialServiceServer{
		service: service,
		logger:  logger,
	}
}

func (s *SocialServiceServer) FollowUser(ctx context.Context, req *pb.FollowUserRequest) (*pb.FollowUserResponse, error) {
	if err := s.service.FollowUser(ctx, req.FollowerId, req.FolloweeId); err != nil {
		s.logger.Error(err, "failed to follow user")
		return nil, status.Error(codes.Internal, err.Error())
	}
	return &pb.FollowUserResponse{Success: true}, nil
}

func (s *SocialServiceServer) UnfollowUser(ctx context.Context, req *pb.UnfollowUserRequest) (*pb.UnfollowUserResponse, error) {
	if err := s.service.UnfollowUser(ctx, req.FollowerId, req.FolloweeId); err != nil {
		s.logger.Error(err, "failed to unfollow user")
		return nil, status.Error(codes.Internal, err.Error())
	}
	return &pb.UnfollowUserResponse{Success: true}, nil
}

func (s *SocialServiceServer) GetFollowers(ctx context.Context, req *pb.GetFollowersRequest) (*pb.GetFollowersResponse, error) {
	followers, total, err := s.service.GetFollowers(ctx, req.UserId, int(req.Page), int(req.PageSize))
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	var protoFollowers []*pb.SocialUser
	for _, f := range followers {
		protoFollowers = append(protoFollowers, &pb.SocialUser{
			UserId:     f.FollowerID,
			FollowedAt: timestamppb.New(f.CreatedAt),
		})
	}

	return &pb.GetFollowersResponse{
		Followers: protoFollowers,
		Total:     int32(total),
	}, nil
}

func (s *SocialServiceServer) GetFollowing(ctx context.Context, req *pb.GetFollowingRequest) (*pb.GetFollowingResponse, error) {
	following, total, err := s.service.GetFollowing(ctx, req.UserId, int(req.Page), int(req.PageSize))
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	var protoFollowing []*pb.SocialUser
	for _, f := range following {
		protoFollowing = append(protoFollowing, &pb.SocialUser{
			UserId:     f.FolloweeID,
			FollowedAt: timestamppb.New(f.CreatedAt),
		})
	}

	return &pb.GetFollowingResponse{
		Following: protoFollowing,
		Total:     int32(total),
	}, nil
}

func (s *SocialServiceServer) GetSocialStats(ctx context.Context, req *pb.GetSocialStatsRequest) (*pb.GetSocialStatsResponse, error) {
	stats, err := s.service.GetStats(ctx, req.UserId)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &pb.GetSocialStatsResponse{
		FollowersCount: int32(stats.FollowersCount),
		FollowingCount: int32(stats.FollowingCount),
	}, nil
}
