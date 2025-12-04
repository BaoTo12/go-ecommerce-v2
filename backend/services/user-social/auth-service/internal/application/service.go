package application

import (
	"context"
	"time"

	"github.com/titan-commerce/backend/auth-service/internal/domain"
	"github.com/titan-commerce/backend/auth-service/internal/infrastructure/token"
	"github.com/titan-commerce/backend/pkg/errors"
	"github.com/titan-commerce/backend/pkg/logger"
	"github.com/pquerna/otp/totp"
)

type AuthRepository interface {
	SaveUser(ctx context.Context, user *domain.User) error
	FindByEmail(ctx context.Context, email string) (*domain.User, error)
	FindByID(ctx context.Context, userID string) (*domain.User, error)
	UpdateMFA(ctx context.Context, userID, secret string, enabled bool) error
}

type TokenRepository interface {
	BlacklistToken(ctx context.Context, token string, expiration time.Duration) error
	IsBlacklisted(ctx context.Context, token string) (bool, error)
	StoreRefreshToken(ctx context.Context, userID, token string, expiration time.Duration) error
	GetRefreshToken(ctx context.Context, userID string) (string, error)
	RevokeRefreshToken(ctx context.Context, userID string) error
}

type AuthService struct {
	repo         AuthRepository
	tokenRepo    TokenRepository
	tokenService *token.TokenService
	logger       *logger.Logger
}

func NewAuthService(repo AuthRepository, tokenRepo TokenRepository, tokenService *token.TokenService, logger *logger.Logger) *AuthService {
	return &AuthService{
		repo:         repo,
		tokenRepo:    tokenRepo,
		tokenService: tokenService,
		logger:       logger,
	}
}

func (s *AuthService) Register(ctx context.Context, email, password, fullName string) (string, string, string, error) {
	// Check if user exists
	if _, err := s.repo.FindByEmail(ctx, email); err == nil {
		return "", "", "", errors.New(errors.ErrConflict, "email already exists")
	}

	// Create user
	user, err := domain.NewUser(email, password, fullName)
	if err != nil {
		return "", "", "", err
	}

	if err := s.repo.SaveUser(ctx, user); err != nil {
		return "", "", "", err
	}

	// Generate tokens
	accessToken, refreshToken, err := s.tokenService.GenerateTokens(user.ID, user.Email)
	if err != nil {
		return "", "", "", err
	}

	// Store refresh token
	if err := s.tokenRepo.StoreRefreshToken(ctx, user.ID, refreshToken, 30*24*time.Hour); err != nil {
		return "", "", "", err
	}

	return user.ID, accessToken, refreshToken, nil
}

func (s *AuthService) Login(ctx context.Context, email, password, mfaCode string) (string, string, string, bool, error) {
	user, err := s.repo.FindByEmail(ctx, email)
	if err != nil {
		return "", "", "", false, errors.New(errors.ErrUnauthorized, "invalid credentials")
	}

	if !user.CheckPassword(password) {
		return "", "", "", false, errors.New(errors.ErrUnauthorized, "invalid credentials")
	}

	// MFA Check
	if user.MFAEnabled {
		if mfaCode == "" {
			return user.ID, "", "", true, nil // Signal MFA required
		}
		if !totp.Validate(mfaCode, user.MFASecret) {
			return "", "", "", true, errors.New(errors.ErrUnauthorized, "invalid MFA code")
		}
	}

	// Generate tokens
	accessToken, refreshToken, err := s.tokenService.GenerateTokens(user.ID, user.Email)
	if err != nil {
		return "", "", "", false, err
	}

	// Store refresh token
	if err := s.tokenRepo.StoreRefreshToken(ctx, user.ID, refreshToken, 30*24*time.Hour); err != nil {
		return "", "", "", false, err
	}

	return user.ID, accessToken, refreshToken, false, nil
}

func (s *AuthService) ValidateToken(ctx context.Context, accessToken string) (bool, string, string, error) {
	// Check blacklist
	blacklisted, err := s.tokenRepo.IsBlacklisted(ctx, accessToken)
	if err != nil {
		return false, "", "", err
	}
	if blacklisted {
		return false, "", "", errors.New(errors.ErrUnauthorized, "token is blacklisted")
	}

	// Validate JWT
	userID, email, err := s.tokenService.ValidateAccessToken(accessToken)
	if err != nil {
		return false, "", "", err
	}

	return true, userID, email, nil
}

func (s *AuthService) RefreshToken(ctx context.Context, refreshToken string) (string, string, error) {
	// Validate JWT
	userID, err := s.tokenService.ValidateRefreshToken(refreshToken)
	if err != nil {
		return "", "", err
	}

	// Check if stored
	storedToken, err := s.tokenRepo.GetRefreshToken(ctx, userID)
	if err != nil || storedToken != refreshToken {
		return "", "", errors.New(errors.ErrUnauthorized, "invalid refresh token")
	}

	// Get user email
	user, err := s.repo.FindByID(ctx, userID)
	if err != nil {
		return "", "", err
	}

	// Generate new tokens
	newAccessToken, newRefreshToken, err := s.tokenService.GenerateTokens(userID, user.Email)
	if err != nil {
		return "", "", err
	}

	// Rotate refresh token
	if err := s.tokenRepo.StoreRefreshToken(ctx, userID, newRefreshToken, 30*24*time.Hour); err != nil {
		return "", "", err
	}

	return newAccessToken, newRefreshToken, nil
}

func (s *AuthService) Logout(ctx context.Context, accessToken string) error {
	// Blacklist access token (15 mins)
	if err := s.tokenRepo.BlacklistToken(ctx, accessToken, 15*time.Minute); err != nil {
		return err
	}
	return nil
}

func (s *AuthService) EnableMFA(ctx context.Context, userID string) (string, string, error) {
	key, err := totp.Generate(totp.GenerateOpts{
		Issuer:      "TitanCommerce",
		AccountName: userID,
	})
	if err != nil {
		return "", "", err
	}

	// Store secret temporarily or pending verification? 
	// For simplicity, we'll update user but mark enabled=false until verified, 
	// or just return secret and expect VerifyMFA to enable it.
	// Let's do: Update secret, enabled=false.
	if err := s.repo.UpdateMFA(ctx, userID, key.Secret(), false); err != nil {
		return "", "", err
	}

	return key.Secret(), key.URL(), nil
}

func (s *AuthService) VerifyMFA(ctx context.Context, userID, code string) (bool, error) {
	user, err := s.repo.FindByID(ctx, userID)
	if err != nil {
		return false, err
	}

	if !totp.Validate(code, user.MFASecret) {
		return false, errors.New(errors.ErrUnauthorized, "invalid code")
	}

	if err := s.repo.UpdateMFA(ctx, userID, user.MFASecret, true); err != nil {
		return false, err
	}

	return true, nil
}
