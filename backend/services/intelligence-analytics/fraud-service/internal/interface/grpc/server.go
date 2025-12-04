package grpc

import (
	"context"

	"github.com/titan-commerce/backend/fraud-service/internal/application"
	"github.com/titan-commerce/backend/fraud-service/internal/domain"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type FraudServer struct {
	UnimplementedFraudServiceServer
	service *application.FraudService
}

func NewFraudServer(service *application.FraudService) *FraudServer {
	return &FraudServer{service: service}
}

func (s *FraudServer) Register(server *grpc.Server) {
	RegisterFraudServiceServer(server, s)
}

func (s *FraudServer) CheckTransaction(ctx context.Context, req *CheckTransactionRequest) (*FraudCheckResult, error) {
	check, err := s.service.CheckTransaction(
		ctx,
		req.TransactionId,
		req.UserId,
		req.Amount,
		req.Currency,
		req.IpAddress,
		req.DeviceId,
		req.UserAgent,
	)
	if err != nil {
		return nil, status.Errorf(codes.Internal, err.Error())
	}

	return domainCheckToProto(check), nil
}

func (s *FraudServer) GetFraudCheck(ctx context.Context, req *GetFraudCheckRequest) (*FraudCheckResult, error) {
	check, err := s.service.GetFraudCheck(ctx, req.CheckId)
	if err != nil {
		return nil, status.Errorf(codes.NotFound, "fraud check not found")
	}
	return domainCheckToProto(check), nil
}

func (s *FraudServer) GetUserFraudHistory(ctx context.Context, req *GetUserFraudHistoryRequest) (*GetUserFraudHistoryResponse, error) {
	checks, err := s.service.GetUserFraudHistory(ctx, req.UserId, int(req.Limit))
	if err != nil {
		return nil, status.Errorf(codes.Internal, err.Error())
	}

	pbChecks := make([]*FraudCheckResult, len(checks))
	for i, c := range checks {
		pbChecks[i] = domainCheckToProto(c)
	}

	return &GetUserFraudHistoryResponse{Checks: pbChecks}, nil
}

func (s *FraudServer) OverrideDecision(ctx context.Context, req *OverrideDecisionRequest) (*OverrideDecisionResponse, error) {
	err := s.service.OverrideDecision(ctx, req.CheckId, domain.FraudDecision(req.NewDecision), req.Reason)
	if err != nil {
		return nil, status.Errorf(codes.Internal, err.Error())
	}
	return &OverrideDecisionResponse{Success: true}, nil
}

func (s *FraudServer) GetPendingAlerts(ctx context.Context, req *GetPendingAlertsRequest) (*GetPendingAlertsResponse, error) {
	alerts, err := s.service.GetPendingAlerts(ctx)
	if err != nil {
		return nil, status.Errorf(codes.Internal, err.Error())
	}

	pbAlerts := make([]*FraudAlert, len(alerts))
	for i, a := range alerts {
		pbAlerts[i] = &FraudAlert{
			Id:           a.ID,
			FraudCheckId: a.FraudCheckID,
			AlertType:    a.AlertType,
			Severity:     string(a.Severity),
			Message:      a.Message,
			Acknowledged: a.Acknowledged,
			CreatedAt:    timestamppb.New(a.CreatedAt),
		}
	}

	return &GetPendingAlertsResponse{Alerts: pbAlerts}, nil
}

func domainCheckToProto(check *domain.FraudCheck) *FraudCheckResult {
	return &FraudCheckResult{
		CheckId:          check.ID,
		TransactionId:    check.TransactionID,
		Score:            check.Score,
		RiskLevel:        string(check.RiskLevel),
		Decision:         string(check.Decision),
		Reasons:          check.Reasons,
		ProcessingTimeMs: check.ProcessingTime,
		Features: &FraudFeatures{
			AccountAgeDays:       int32(check.Features.AccountAgeDays),
			TotalOrders:          int32(check.Features.TotalOrders),
			AvgOrderValue:        check.Features.AvgOrderValue,
			OrdersLast_24H:       int32(check.Features.OrdersLast24h),
			OrdersLast_7D:        int32(check.Features.OrdersLast7d),
			FailedPaymentsLast_7D: int32(check.Features.FailedPaymentsLast7d),
			UniqueDevices:        int32(check.Features.UniqueDevices),
			UniqueIps:            int32(check.Features.UniqueIPs),
			NewDevice:            check.Features.NewDevice,
			NewIp:                check.Features.NewIP,
			VelocityScore:        check.Features.VelocityScore,
			AmountDeviation:      check.Features.AmountDeviation,
		},
		CreatedAt: timestamppb.New(check.CreatedAt),
	}
}

// Placeholder types and registration
type UnimplementedFraudServiceServer struct{}

func (UnimplementedFraudServiceServer) CheckTransaction(context.Context, *CheckTransactionRequest) (*FraudCheckResult, error) {
	return nil, status.Errorf(codes.Unimplemented, "unimplemented")
}
func (UnimplementedFraudServiceServer) GetFraudCheck(context.Context, *GetFraudCheckRequest) (*FraudCheckResult, error) {
	return nil, status.Errorf(codes.Unimplemented, "unimplemented")
}
func (UnimplementedFraudServiceServer) GetUserFraudHistory(context.Context, *GetUserFraudHistoryRequest) (*GetUserFraudHistoryResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "unimplemented")
}
func (UnimplementedFraudServiceServer) OverrideDecision(context.Context, *OverrideDecisionRequest) (*OverrideDecisionResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "unimplemented")
}
func (UnimplementedFraudServiceServer) GetPendingAlerts(context.Context, *GetPendingAlertsRequest) (*GetPendingAlertsResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "unimplemented")
}
func (UnimplementedFraudServiceServer) mustEmbedUnimplementedFraudServiceServer() {}

func RegisterFraudServiceServer(s *grpc.Server, srv FraudServiceServer) {}

type FraudServiceServer interface {
	CheckTransaction(context.Context, *CheckTransactionRequest) (*FraudCheckResult, error)
	GetFraudCheck(context.Context, *GetFraudCheckRequest) (*FraudCheckResult, error)
	GetUserFraudHistory(context.Context, *GetUserFraudHistoryRequest) (*GetUserFraudHistoryResponse, error)
	OverrideDecision(context.Context, *OverrideDecisionRequest) (*OverrideDecisionResponse, error)
	GetPendingAlerts(context.Context, *GetPendingAlertsRequest) (*GetPendingAlertsResponse, error)
	mustEmbedUnimplementedFraudServiceServer()
}

// Message types
type CheckTransactionRequest struct {
	TransactionId string
	UserId        string
	Amount        float64
	Currency      string
	IpAddress     string
	DeviceId      string
	UserAgent     string
	Metadata      map[string]string
}
type FraudCheckResult struct {
	CheckId          string
	TransactionId    string
	Score            float64
	RiskLevel        string
	Decision         string
	Reasons          []string
	ProcessingTimeMs int64
	Features         *FraudFeatures
	CreatedAt        *timestamppb.Timestamp
}
type FraudFeatures struct {
	AccountAgeDays        int32
	TotalOrders           int32
	AvgOrderValue         float64
	OrdersLast_24H        int32
	OrdersLast_7D         int32
	FailedPaymentsLast_7D int32
	UniqueDevices         int32
	UniqueIps             int32
	NewDevice             bool
	NewIp                 bool
	VelocityScore         float64
	AmountDeviation       float64
}
type GetFraudCheckRequest struct{ CheckId string }
type GetUserFraudHistoryRequest struct {
	UserId string
	Limit  int32
}
type GetUserFraudHistoryResponse struct{ Checks []*FraudCheckResult }
type OverrideDecisionRequest struct {
	CheckId     string
	NewDecision string
	Reason      string
	AdminId     string
}
type OverrideDecisionResponse struct{ Success bool }
type GetPendingAlertsRequest struct{ Limit int32 }
type GetPendingAlertsResponse struct{ Alerts []*FraudAlert }
type FraudAlert struct {
	Id           string
	FraudCheckId string
	AlertType    string
	Severity     string
	Message      string
	Acknowledged bool
	CreatedAt    *timestamppb.Timestamp
}
