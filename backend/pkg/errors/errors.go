package errors

import (
	"fmt"
	"net/http"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type ErrorCode string

const (
	ErrInternal            ErrorCode = "INTERNAL_ERROR"
	ErrNotFound            ErrorCode = "NOT_FOUND"
	ErrInvalidInput        ErrorCode = "INVALID_INPUT"
	ErrUnauthorized        ErrorCode = "UNAUTHORIZED"
	ErrForbidden           ErrorCode = "FORBIDDEN"
	ErrConflict            ErrorCode = "CONFLICT"
	ErrInsufficientStock   ErrorCode = "INSUFFICIENT_STOCK"
	ErrInsufficientBalance ErrorCode = "INSUFFICIENT_BALANCE"
	ErrPaymentFailed       ErrorCode = "PAYMENT_FAILED"
	ErrOrderNotCancellable ErrorCode = "ORDER_NOT_CANCELLABLE"
)

type AppError struct {
	Code       ErrorCode              `json:"code"`
	Message    string                 `json:"message"`
	HTTPStatus int                    `json:"-"`
	GRPCCode   codes.Code             `json:"-"`
	Details    map[string]interface{} `json:"details,omitempty"`
	Err        error                  `json:"-"`
}

func (e *AppError) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("%s: %s (caused by: %v)", e.Code, e.Message, e.Err)
	}
	return fmt.Sprintf("%s: %s", e.Code, e.Message)
}

func New(code ErrorCode, message string) *AppError {
	return &AppError{
		Code:       code,
		Message:    message,
		HTTPStatus: codeToHTTP(code),
		GRPCCode:   codeToGRPC(code),
	}
}

func Wrap(code ErrorCode, message string, err error) *AppError {
	return &AppError{
		Code:       code,
		Message:    message,
		HTTPStatus: codeToHTTP(code),
		GRPCCode:   codeToGRPC(code),
		Err:        err,
	}
}

func (e *AppError) ToGRPCError() error {
	return status.Error(e.GRPCCode, e.Message)
}

func codeToHTTP(code ErrorCode) int {
	switch code {
	case ErrNotFound:
		return http.StatusNotFound
	case ErrInvalidInput:
		return http.StatusBadRequest
	case ErrUnauthorized:
		return http.StatusUnauthorized
	case ErrForbidden:
		return http.StatusForbidden
	case ErrConflict:
		return http.StatusConflict
	case ErrInsufficientStock, ErrInsufficientBalance, ErrOrderNotCancellable:
		return http.StatusUnprocessableEntity
	case ErrPaymentFailed:
		return http.StatusPaymentRequired
	default:
		return http.StatusInternalServerError
	}
}

func codeToGRPC(code ErrorCode) codes.Code {
	switch code {
	case ErrNotFound:
		return codes.NotFound
	case ErrInvalidInput:
		return codes.InvalidArgument
	case ErrUnauthorized:
		return codes.Unauthenticated
	case ErrForbidden:
		return codes.PermissionDenied
	case ErrConflict:
		return codes.AlreadyExists
	case ErrInsufficientStock, ErrInsufficientBalance, ErrOrderNotCancellable:
		return codes.FailedPrecondition
	case ErrPaymentFailed:
		return codes.Aborted
	default:
		return codes.Internal
	}
}
