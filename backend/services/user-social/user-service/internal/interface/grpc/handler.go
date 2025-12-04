package grpc

import (
	"context"

	"github.com/titan-commerce/backend/user-service/internal/application"
	"github.com/titan-commerce/backend/user-service/internal/domain"
	pb "github.com/titan-commerce/backend/user-service/proto/user/v1"
	"github.com/titan-commerce/backend/pkg/logger"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type UserServiceServer struct {
	pb.UnimplementedUserServiceServer
	service *application.UserService
	logger  *logger.Logger
}

func NewUserServiceServer(service *application.UserService, logger *logger.Logger) *UserServiceServer {
	return &UserServiceServer{
		service: service,
		logger:  logger,
	}
}

func (s *UserServiceServer) GetUser(ctx context.Context, req *pb.GetUserRequest) (*pb.GetUserResponse, error) {
	user, err := s.service.GetUser(ctx, req.UserId)
	if err != nil {
		return nil, status.Error(codes.NotFound, err.Error())
	}

	return &pb.GetUserResponse{
		User: domainUserToProto(user),
	}, nil
}

func (s *UserServiceServer) UpdateProfile(ctx context.Context, req *pb.UpdateProfileRequest) (*pb.UpdateProfileResponse, error) {
	user, err := s.service.UpdateProfile(ctx, req.UserId, req.FullName, req.PhoneNumber, req.AvatarUrl)
	if err != nil {
		s.logger.Error(err, "failed to update profile")
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &pb.UpdateProfileResponse{
		User: domainUserToProto(user),
	}, nil
}

func (s *UserServiceServer) AddAddress(ctx context.Context, req *pb.AddAddressRequest) (*pb.AddAddressResponse, error) {
	address, err := s.service.AddAddress(ctx, req.UserId, "", "", req.Street, req.City, req.ZipCode, req.Country, false)
	if err != nil {
		s.logger.Error(err, "failed to add address")
		return nil, status.Error(codes.Internal, err.Error())
	}

	// Get user to return full response
	user, err := s.service.GetUser(ctx, req.UserId)
	if err != nil {
		s.logger.Error(err, "failed to get user")
		return nil, status.Error(codes.Internal, err.Error())
	}
	_ = address // address created successfully

	return &pb.AddAddressResponse{
		User: domainUserToProto(user),
	}, nil
}

func domainUserToProto(user *domain.User) *pb.User {
	return &pb.User{
		UserId:      user.ID,
		Email:       user.Email,
		FullName:    user.FullName,
		PhoneNumber: user.PhoneNumber,
		AvatarUrl:   user.AvatarURL,
		Addresses:   nil,
	}
}
