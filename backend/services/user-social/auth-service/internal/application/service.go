package application

import (
	"context"
	"crypto/sha256"
	"encoding/hex"

	"github.com/titan-commerce/backend/auth-service/internal/domain"
	"github.com/titan-commerce/backend/pkg/auth"
	"github.com/titan-commerce/backend/pkg/errors"
	"github.com/titan-commerce/backend/pkg/logger"
)

type UserRepository interface {
	Save(ctx context.Context, user *domain.User) error
	FindByID(ctx context.Context, userID string) (*domain.User, error)
	FindByEmail(ctx context.Context, email string) (*domain.User, error)
}

type RefreshTokenRepository interface {
	Save(ctx context.Context, token *domain.RefreshToken) error
	FindByToken(ctx context.Context, token string) (*domain.RefreshToken, error)
	DeleteByToken(ctx context.Context, token string) error
}

type AuthService struct {
	userRepo         UserRepository
	refreshTokenRepo RefreshTokenRepository
	jwtService       *auth.JWTService
	logger           *logger.Logger
}

func NewAuthService(userRepo UserRepository, tokenRepo RefreshTokenRepository, jwtService *auth.JWTService, logger *logger.Logger) *AuthService {
	return &AuthService{
		userRepo:         userRepo,
		refreshTokenRepo: tokenRepo,
		jwtService:       jwtService,
		logger:           logger,
	}
}

// Register creates a new user account (Command)
func (s *AuthService) Register(ctx context.Context, email, password, name string) (*domain.User, string, string, error) {
	// Check if user already exists
	existing, err := s.userRepo.FindByEmail(ctx, email)
	if err == nil && existing != nil {
		return nil, "", "", errors.New(errors.ErrAlreadyExists, "user with this email already exists")
	}

	// Create user
	user, err := domain.NewUser(email, password, name)
	if err != nil {
		return nil, "", "", err
	}

	// Save user
	if err := s.userRepo.Save(ctx, user); err != nil {
		s.logger.Error(err, "failed to save user")
		return nil, "", "", err
	}

	// Generate tokens
	cellID := s.computeCellID(user.ID)
	accessToken, err := s.jwtService.GenerateAccessToken(user.ID, user.Email, cellID, user.Roles)
	if err != nil {
		return nil, "", "", err
	}

	refreshToken := domain.NewRefreshToken(user.ID, 30) // 30 days
	if err := s.refreshTokenRepo.Save(ctx, refreshToken); err != nil {
		return nil, "", "", err
	}

	s.logger.Infof("User registered: %s (%s)", user.ID, user.Email)
	return user, accessToken, refreshToken.Token, nil
}

// Login authenticates user and returns tokens (Command)
func (s *AuthService) Login(ctx context.Context, email, password string) (*domain.User, string, string, error) {
	user, err := s.userRepo.FindByEmail(ctx, email)
	if err != nil {
		return nil, "", "", errors.New(errors.ErrUnauthorized, "invalid credentials")
	}

	if !user.VerifyPassword(password) {
		return nil, "", "", errors.New(errors.ErrUnauthorized, "invalid credentials")
	}

	// Generate tokens
	cellID := s.computeCellID(user.ID)
	accessToken, err := s.jwtService.GenerateAccessToken(user.ID, user.Email, cellID, user.Roles)
	if err != nil {
		return nil, "", "", err
	}

	refreshToken := domain.NewRefreshToken(user.ID, 30)
	if err := s.refreshTokenRepo.Save(ctx, refreshToken); err != nil {
		return nil, "", "", err
	}

	s.logger.Infof("User logged in: %s (%s)", user.ID, user.Email)
	return user, accessToken, refreshToken.Token, nil
}

// RefreshAccessToken generates new access token from refresh token (Command)
func (s *AuthService) RefreshAccessToken(ctx context.Context, refreshTokenStr string) (string, string, error) {
	refreshToken, err := s.refreshTokenRepo.FindByToken(ctx, refreshTokenStr)
	if err != nil {
		return "", "", errors.New(errors.ErrUnauthorized, "invalid refresh token")
	}

	if refreshToken.IsExpired() {
		s.refreshTokenRepo.DeleteByToken(ctx, refreshTokenStr)
		return "", "", errors.New(errors.ErrUnauthorized, "refresh token expired")
	}

	user, err := s.userRepo.FindByID(ctx, refreshToken.UserID)
	if err != nil {
		return "", "", err
	}

	// Generate new access token
	cellID := s.computeCellID(user.ID)
	accessToken, err := s.jwtService.GenerateAccessToken(user.ID, user.Email, cellID, user.Roles)
	if err != nil {
		return "", "", err
	}

	// Generate new refresh token
	newRefreshToken := domain.NewRefreshToken(user.ID, 30)
	if err := s.refreshTokenRepo.Save(ctx, newRefreshToken); err != nil {
		return "", "", err
	}

	// Delete old refresh token
	s.refreshTokenRepo.DeleteByToken(ctx, refreshTokenStr)

	return accessToken, newRefreshToken.Token, nil
}

// Logout invalidates refresh token (Command)
func (s *AuthService) Logout(ctx context.Context, refreshTokenStr string) error {
	return s.refreshTokenRepo.DeleteByToken(ctx, refreshTokenStr)
}

// VerifyToken verifies access token and returns user (Query)
func (s *AuthService) VerifyToken(ctx context.Context, accessToken string) (*domain.User, error) {
	claims, err := s.jwtService.VerifyAccessToken(accessToken)
	if err != nil {
		return nil, err
	}

	user, err := s.userRepo.FindByID(ctx, claims.UserID)
	if err != nil {
		return nil, errors.New(errors.ErrUnauthorized, "user not found")
	}

	return user, nil
}

// computeCellID computes cell ID for user using consistent hashing
func (s *AuthService) computeCellID(userID string) string {
	hash := sha256.Sum256([]byte(userID))
	hashInt := uint64(hash[0]) | uint64(hash[1])<<8 | uint64(hash[2])<<16 | uint64(hash[3])<<24
	cellNum := (hashInt % 500) + 1
	return "cell-" + hex.EncodeToString([]byte{byte(cellNum / 100), byte((cellNum % 100) / 10), byte(cellNum % 10)})
}
