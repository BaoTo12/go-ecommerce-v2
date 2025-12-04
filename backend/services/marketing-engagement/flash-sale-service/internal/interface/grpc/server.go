package grpc

import (
	"context"
	"time"

	"github.com/titan-commerce/backend/flash-sale-service/internal/application"
	"github.com/titan-commerce/backend/flash-sale-service/internal/domain"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type FlashSaleServer struct {
	UnimplementedFlashSaleServiceServer
	service *application.FlashSaleService
}

func NewFlashSaleServer(service *application.FlashSaleService) *FlashSaleServer {
	return &FlashSaleServer{service: service}
}

func (s *FlashSaleServer) Register(server *grpc.Server) {
	RegisterFlashSaleServiceServer(server, s)
}

func (s *FlashSaleServer) GetChallenge(ctx context.Context, req *GetChallengeRequest) (*GetChallengeResponse, error) {
	challenge := s.service.GetChallenge(req.SaleId, req.UserId)
	return &GetChallengeResponse{
		Challenge:  challenge,
		Difficulty: 4,
		ExpiresAt:  time.Now().Add(5 * time.Minute).Unix(),
	}, nil
}

func (s *FlashSaleServer) AttemptPurchase(ctx context.Context, req *AttemptPurchaseRequest) (*AttemptPurchaseResponse, error) {
	reservation, err := s.service.AttemptPurchase(ctx, req.SaleId, req.UserId, int(req.Quantity), req.Challenge, req.Nonce)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, err.Error())
	}

	sale, _ := s.service.GetFlashSale(ctx, req.SaleId)
	totalPrice := sale.SalePrice * float64(req.Quantity)

	return &AttemptPurchaseResponse{
		ReservationId: reservation.ID,
		ExpiresAt:     timestamppb.New(reservation.ExpiresAt),
		TotalPrice:    totalPrice,
	}, nil
}

func (s *FlashSaleServer) ConfirmPurchase(ctx context.Context, req *ConfirmPurchaseRequest) (*ConfirmPurchaseResponse, error) {
	err := s.service.ConfirmPurchase(ctx, req.ReservationId, req.UserId)
	if err != nil {
		return nil, status.Errorf(codes.Internal, err.Error())
	}

	return &ConfirmPurchaseResponse{
		Success: true,
		OrderId: req.ReservationId,
	}, nil
}

func (s *FlashSaleServer) GetActiveFlashSales(ctx context.Context, req *GetActiveFlashSalesRequest) (*GetActiveFlashSalesResponse, error) {
	sales, err := s.service.GetActiveFlashSales(ctx)
	if err != nil {
		return nil, status.Errorf(codes.Internal, err.Error())
	}

	pbSales := make([]*FlashSale, len(sales))
	for i, sale := range sales {
		pbSales[i] = domainToProto(sale)
	}

	return &GetActiveFlashSalesResponse{
		Sales: pbSales,
		Total: int32(len(sales)),
	}, nil
}

func (s *FlashSaleServer) GetFlashSale(ctx context.Context, req *GetFlashSaleRequest) (*FlashSale, error) {
	sale, err := s.service.GetFlashSale(ctx, req.SaleId)
	if err != nil {
		return nil, status.Errorf(codes.NotFound, "flash sale not found")
	}
	return domainToProto(sale), nil
}

func (s *FlashSaleServer) CreateFlashSale(ctx context.Context, req *CreateFlashSaleRequest) (*FlashSale, error) {
	sale, err := s.service.CreateFlashSale(
		ctx,
		req.ProductId,
		req.OriginalPrice,
		req.SalePrice,
		int(req.TotalQuantity),
		int(req.MaxPerUser),
		req.StartTime.AsTime(),
		req.EndTime.AsTime(),
	)
	if err != nil {
		return nil, status.Errorf(codes.Internal, err.Error())
	}
	return domainToProto(sale), nil
}

func domainToProto(sale *domain.FlashSale) *FlashSale {
	return &FlashSale{
		Id:                sale.ID,
		ProductId:         sale.ProductID,
		OriginalPrice:     sale.OriginalPrice,
		SalePrice:         sale.SalePrice,
		DiscountPercent:   int32(sale.DiscountPercent),
		TotalQuantity:     int32(sale.TotalQuantity),
		SoldQuantity:      int32(sale.SoldQuantity),
		RemainingQuantity: int32(sale.RemainingQuantity()),
		MaxPerUser:        int32(sale.MaxPerUser),
		Status:            string(sale.Status),
		StartTime:         timestamppb.New(sale.StartTime),
		EndTime:           timestamppb.New(sale.EndTime),
	}
}

// Placeholder for generated code
type UnimplementedFlashSaleServiceServer struct{}

func (UnimplementedFlashSaleServiceServer) GetChallenge(context.Context, *GetChallengeRequest) (*GetChallengeResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "unimplemented")
}
func (UnimplementedFlashSaleServiceServer) AttemptPurchase(context.Context, *AttemptPurchaseRequest) (*AttemptPurchaseResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "unimplemented")
}
func (UnimplementedFlashSaleServiceServer) ConfirmPurchase(context.Context, *ConfirmPurchaseRequest) (*ConfirmPurchaseResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "unimplemented")
}
func (UnimplementedFlashSaleServiceServer) GetActiveFlashSales(context.Context, *GetActiveFlashSalesRequest) (*GetActiveFlashSalesResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "unimplemented")
}
func (UnimplementedFlashSaleServiceServer) GetFlashSale(context.Context, *GetFlashSaleRequest) (*FlashSale, error) {
	return nil, status.Errorf(codes.Unimplemented, "unimplemented")
}
func (UnimplementedFlashSaleServiceServer) CreateFlashSale(context.Context, *CreateFlashSaleRequest) (*FlashSale, error) {
	return nil, status.Errorf(codes.Unimplemented, "unimplemented")
}
func (UnimplementedFlashSaleServiceServer) mustEmbedUnimplementedFlashSaleServiceServer() {}

func RegisterFlashSaleServiceServer(s *grpc.Server, srv FlashSaleServiceServer) {
	// In production, this would be generated by protoc
}

type FlashSaleServiceServer interface {
	GetChallenge(context.Context, *GetChallengeRequest) (*GetChallengeResponse, error)
	AttemptPurchase(context.Context, *AttemptPurchaseRequest) (*AttemptPurchaseResponse, error)
	ConfirmPurchase(context.Context, *ConfirmPurchaseRequest) (*ConfirmPurchaseResponse, error)
	GetActiveFlashSales(context.Context, *GetActiveFlashSalesRequest) (*GetActiveFlashSalesResponse, error)
	GetFlashSale(context.Context, *GetFlashSaleRequest) (*FlashSale, error)
	CreateFlashSale(context.Context, *CreateFlashSaleRequest) (*FlashSale, error)
	mustEmbedUnimplementedFlashSaleServiceServer()
}

// Message types (would be generated by protoc)
type GetChallengeRequest struct {
	SaleId string
	UserId string
}
type GetChallengeResponse struct {
	Challenge  string
	Difficulty int32
	ExpiresAt  int64
}
type AttemptPurchaseRequest struct {
	SaleId    string
	UserId    string
	Quantity  int32
	Challenge string
	Nonce     string
}
type AttemptPurchaseResponse struct {
	ReservationId string
	ExpiresAt     *timestamppb.Timestamp
	TotalPrice    float64
}
type ConfirmPurchaseRequest struct {
	ReservationId string
	UserId        string
	PaymentId     string
}
type ConfirmPurchaseResponse struct {
	Success bool
	OrderId string
}
type GetActiveFlashSalesRequest struct {
	Limit  int32
	Offset int32
}
type GetActiveFlashSalesResponse struct {
	Sales []*FlashSale
	Total int32
}
type GetFlashSaleRequest struct {
	SaleId string
}
type CreateFlashSaleRequest struct {
	ProductId     string
	OriginalPrice float64
	SalePrice     float64
	TotalQuantity int32
	MaxPerUser    int32
	StartTime     *timestamppb.Timestamp
	EndTime       *timestamppb.Timestamp
}
type FlashSale struct {
	Id                string
	ProductId         string
	OriginalPrice     float64
	SalePrice         float64
	DiscountPercent   int32
	TotalQuantity     int32
	SoldQuantity      int32
	RemainingQuantity int32
	MaxPerUser        int32
	Status            string
	StartTime         *timestamppb.Timestamp
	EndTime           *timestamppb.Timestamp
}
