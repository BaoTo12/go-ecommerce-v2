package grpc

import (
	"context"

	"github.com/titan-commerce/backend/auth-service/internal/application"
	pb "github.com/titan-commerce/backend/auth-service/proto/auth/v1"
	"github.com/titan-commerce/backend/pkg/logger"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type AuthServiceServer struct {
	pb.UnimplementedAuthServiceServer
	service *application.AuthService
	logger  *logger.Logger
}

func NewAuthServiceServer(service *application.AuthService, logger *logger.Logger) *AuthServiceServer {
	return &AuthServiceServer{
		service: service,
		logger:  logger,
	}
}

func (s *AuthServiceServer) Register(ctx context.Context, req *pb.RegisterRequest) (*pb.RegisterResponse, error) {
	userID, accessToken, refreshToken, err := s.service.Register(ctx, req.Email, req.Password, req.FullName)
	if err != nil {
		s.logger.Error(err, "failed to register")
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &pb.RegisterResponse{
		UserId:       userID,
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}

func (s *AuthServiceServer) Login(ctx context.Context, req *pb.LoginRequest) (*pb.LoginResponse, error) {
	userID, accessToken, refreshToken, mfaRequired, err := s.service.Login(ctx, req.Email, req.Password, req.MfaCode)
	if err != nil {
		s.logger.Error(err, "failed to login")
		return nil, status.Error(codes.Unauthenticated, err.Error())
	}

	return &pb.LoginResponse{
		UserId:       userID,
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		MfaRequired:  mfaRequired,
	}, nil
}

func (s *AuthServiceServer) ValidateToken(ctx context.Context, req *pb.ValidateTokenRequest) (*pb.ValidateTokenResponse, error) {
	valid, userID, email, err := s.service.ValidateToken(ctx, req.AccessToken)
	if err != nil {
		return &pb.ValidateTokenResponse{Valid: false}, nil
	}

	return &pb.ValidateTokenResponse{
		Valid:  valid,
		UserId: userID,
		Email:  email,
	}, nil
}

func (s *AuthServiceServer) RefreshToken(ctx context.Context, req *pb.RefreshTokenRequest) (*pb.RefreshTokenResponse, error) {
	accessToken, refreshToken, err := s.service.RefreshToken(ctx, req.RefreshToken)
	if err != nil {
		return nil, status.Error(codes.Unauthenticated, err.Error())
	}

	return &pb.RefreshTokenResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}

func (s *AuthServiceServer) Logout(ctx context.Context, req *pb.LogoutRequest) (*pb.LogoutResponse, error) {
	if err := s.service.Logout(ctx, req.AccessToken); err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	return &pb.LogoutResponse{Success: true}, nil
}

func (s *AuthServiceServer) EnableMFA(ctx context.Context, req *pb.EnableMFARequest) (*pb.EnableMFAResponse, error) {
	secret, qrCodeURL, err := s.service.EnableMFA(ctx, req.UserId)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	return &pb.EnableMFAResponse{
		Secret:    secret,
		QrCodeUrl: qrCodeURL,
	}, nil
}

func (s *AuthServiceServer) VerifyMFA(ctx context.Context, req *pb.VerifyMFARequest) (*pb.VerifyMFAResponse, error) {
	success, err := s.service.VerifyMFA(ctx, req.UserId, req.Code)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	return &pb.VerifyMFAResponse{Success: success}, nil
}
