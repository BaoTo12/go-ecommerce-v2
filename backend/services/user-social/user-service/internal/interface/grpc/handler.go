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
		User: domainToProto(user),
	}, nil
}

func (s *UserServiceServer) UpdateProfile(ctx context.Context, req *pb.UpdateProfileRequest) (*pb.UpdateProfileResponse, error) {
	user, err := s.service.UpdateProfile(ctx, req.UserId, req.FullName, req.PhoneNumber, req.AvatarUrl)
	if err != nil {
		s.logger.Error(err, "failed to update profile")
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &pb.UpdateProfileResponse{
		User: domainToProto(user),
	}, nil
}

func (s *UserServiceServer) AddAddress(ctx context.Context, req *pb.AddAddressRequest) (*pb.AddAddressResponse, error) {
	address := domain.Address{
		Street:  req.Street,
		City:    req.City,
		State:   req.State,
		ZipCode: req.ZipCode,
		Country: req.Country,
	}

	user, err := s.service.AddAddress(ctx, req.UserId, address)
	if err != nil {
		s.logger.Error(err, "failed to add address")
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &pb.AddAddressResponse{
		User: domainToProto(user),
	}, nil
}

func domainToProto(user *domain.User) *pb.User {
	addresses := make([]*pb.Address, len(user.Addresses))
	for i, addr := range user.Addresses {
		addresses[i] = &pb.Address{
			Street:  addr.Street,
			City:    addr.City,
			State:   addr.State,
			ZipCode: addr.ZipCode,
			Country: addr.Country,
		}
	}

	return &pb.User{
		UserId:      user.ID,
		Email:       user.Email,
		FullName:    user.FullName,
		PhoneNumber: user.PhoneNumber,
		AvatarUrl:   user.AvatarURL,
		Addresses:   addresses,
	}
}
