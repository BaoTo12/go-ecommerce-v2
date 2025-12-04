package grpc

import (
	"context"

	"github.com/titan-commerce/backend/ad-service/internal/application"
	"github.com/titan-commerce/backend/ad-service/internal/domain"
	pb "github.com/titan-commerce/backend/ad-service/proto/ad/v1"
	"github.com/titan-commerce/backend/pkg/logger"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type AdServiceServer struct {
	pb.UnimplementedAdServiceServer
	service *application.AdService
	logger  *logger.Logger
}

func NewAdServiceServer(service *application.AdService, logger *logger.Logger) *AdServiceServer {
	return &AdServiceServer{
		service: service,
		logger:  logger,
	}
}

func (s *AdServiceServer) CreateCampaign(ctx context.Context, req *pb.CreateCampaignRequest) (*pb.CreateCampaignResponse, error) {
	campaign, err := s.service.CreateCampaign(
		ctx, req.SellerId, req.ProductId, req.Budget, req.BidAmount,
		req.StartTime.AsTime(), req.EndTime.AsTime(),
	)
	if err != nil {
		s.logger.Error(err, "failed to create campaign")
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &pb.CreateCampaignResponse{
		Campaign: domainToProto(campaign),
	}, nil
}

func (s *AdServiceServer) GetAds(ctx context.Context, req *pb.GetAdsRequest) (*pb.GetAdsResponse, error) {
	campaigns, err := s.service.GetAds(ctx, req.Context, int(req.Limit))
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	var ads []*pb.Ad
	for _, c := range campaigns {
		ads = append(ads, &pb.Ad{
			Id:            c.ID, // Using Campaign ID as Ad ID for simplicity
			ProductId:     c.ProductID,
			Title:         "Sponsored Product", // Placeholder
			ImageUrl:      "",                  // Placeholder
			Price:         0,                   // Placeholder
			TrackingToken: c.ID,
		})
	}

	return &pb.GetAdsResponse{Ads: ads}, nil
}

func (s *AdServiceServer) TrackAdEvent(ctx context.Context, req *pb.TrackAdEventRequest) (*pb.TrackAdEventResponse, error) {
	if err := s.service.TrackEvent(ctx, req.AdId, req.UserId, req.EventType); err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	return &pb.TrackAdEventResponse{Success: true}, nil
}

func domainToProto(c *domain.Campaign) *pb.Campaign {
	return &pb.Campaign{
		Id:              c.ID,
		SellerId:        c.SellerID,
		ProductId:       c.ProductID,
		Budget:          c.Budget,
		RemainingBudget: c.RemainingBudget,
		Status:          c.Status,
	}
}
