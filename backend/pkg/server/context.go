package server

import (
	"context"

	"github.com/titan-commerce/backend/pkg/security"
)

// Context keys
type claimsKey struct{}
type requestBodyKey struct{}
type userIDKey struct{}

// SetClaims adds JWT claims to context
func SetClaims(ctx context.Context, claims *security.JWTClaims) context.Context {
	return context.WithValue(ctx, claimsKey{}, claims)
}

// GetClaims retrieves JWT claims from context
func GetClaims(ctx context.Context) *security.JWTClaims {
	if claims, ok := ctx.Value(claimsKey{}).(*security.JWTClaims); ok {
		return claims
	}
	return nil
}

// GetUserID retrieves the current user ID from context
func GetUserID(ctx context.Context) string {
	claims := GetClaims(ctx)
	if claims != nil {
		return claims.Subject
	}
	return ""
}

// GetUserRole retrieves the current user role from context
func GetUserRole(ctx context.Context) string {
	claims := GetClaims(ctx)
	if claims != nil {
		return claims.Role
	}
	return ""
}

// SetRequestBody stores the parsed request body in context
func SetRequestBody[T any](ctx context.Context, body *T) context.Context {
	return context.WithValue(ctx, requestBodyKey{}, body)
}

// GetRequestBody retrieves the parsed request body from context
func GetRequestBody[T any](ctx context.Context) *T {
	if body, ok := ctx.Value(requestBodyKey{}).(*T); ok {
		return body
	}
	return nil
}

// GetRequestID retrieves the request ID from context
func GetRequestID(ctx context.Context) string {
	if id, ok := ctx.Value(requestIDKey{}).(string); ok {
		return id
	}
	return ""
}
