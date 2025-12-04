package handler

import (
	"context"

	"github.com/titan-commerce/backend/wallet-service/internal/application"
	pb "github.com/titan-commerce/backend/wallet-service/proto/wallet/v1"
	"github.com/titan-commerce/backend/pkg/logger"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type WalletServiceServer struct {
	pb.UnimplementedWalletServiceServer
	service *application.WalletService
	logger  *logger.Logger
}

func NewWalletServiceServer(service *application.WalletService, logger *logger.Logger) *WalletServiceServer {
	return &WalletServiceServer{
		service: service,
		logger:  logger,
	}
}

func (s *WalletServiceServer) GetBalance(ctx context.Context, req *pb.GetBalanceRequest) (*pb.GetBalanceResponse, error) {
	wallet, err := s.service.GetBalance(ctx, req.UserId)
	if err != nil {
		return nil, status.Error(codes.NotFound, err.Error())
	}

	return &pb.GetBalanceResponse{
		Wallet: &pb.Wallet{
			WalletId:         wallet.WalletID,
			UserId:           wallet.UserID,
			AvailableBalance: wallet.AvailableBalance,
			HeldBalance:      wallet.HeldBalance,
			TotalBalance:     wallet.GetTotalBalance(),
			Currency:         wallet.Currency,
		},
	}, nil
}

func (s *WalletServiceServer) Deposit(ctx context.Context, req *pb.DepositRequest) (*pb.DepositResponse, error) {
	wallet, err := s.service.Deposit(ctx, req.UserId, req.Amount)
	if err != nil {
		s.logger.Error(err, "failed to deposit")
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &pb.DepositResponse{
		Wallet: &pb.Wallet{
			WalletId:         wallet.WalletID,
			UserId:           wallet.UserID,
			AvailableBalance: wallet.AvailableBalance,
			HeldBalance:      wallet.HeldBalance,
			TotalBalance:     wallet.GetTotalBalance(),
			Currency:         wallet.Currency,
		},
	}, nil
}

func (s *WalletServiceServer) HoldFunds(ctx context.Context, req *pb.HoldFundsRequest) (*pb.HoldFundsResponse, error) {
	holdID, err := s.service.HoldFunds(ctx, req.UserId, req.OrderId, req.Amount)
	if err != nil {
		s.logger.Error(err, "failed to hold funds")
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &pb.HoldFundsResponse{
		Success: true,
		HoldId:  holdID,
	}, nil
}

func (s *WalletServiceServer) ReleaseFunds(ctx context.Context, req *pb.ReleaseFundsRequest) (*pb.ReleaseFundsResponse, error) {
	// Extract userID from holdID or lookup
	// For simplicity, assuming holdID format includes userID
	err := s.service.ReleaseFunds(ctx, "", req.HoldId, 0, req.ReleaseToUser)
	if err != nil {
		s.logger.Error(err, "failed to release funds")
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &pb.ReleaseFundsResponse{
		Success: true,
	}, nil
}
