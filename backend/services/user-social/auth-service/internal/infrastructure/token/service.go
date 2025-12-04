package infrastructure

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/titan-commerce/backend/pkg/errors"
)

type TokenService struct {
	accessSecret  []byte
	refreshSecret []byte
	accessTTL     time.Duration
	refreshTTL    time.Duration
}

func NewTokenService(accessSecret, refreshSecret string) *TokenService {
	return &TokenService{
		accessSecret:  []byte(accessSecret),
		refreshSecret: []byte(refreshSecret),
		accessTTL:     15 * time.Minute,
		refreshTTL:    30 * 24 * time.Hour, // 30 days
	}
}

func (s *TokenService) GenerateTokens(userID, email string) (string, string, error) {
	// Access Token
	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub":   userID,
		"email": email,
		"exp":   time.Now().Add(s.accessTTL).Unix(),
		"type":  "access",
	})
	accessTokenString, err := accessToken.SignedString(s.accessSecret)
	if err != nil {
		return "", "", err
	}

	// Refresh Token
	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub":  userID,
		"exp":  time.Now().Add(s.refreshTTL).Unix(),
		"type": "refresh",
	})
	refreshTokenString, err := refreshToken.SignedString(s.refreshSecret)
	if err != nil {
		return "", "", err
	}

	return accessTokenString, refreshTokenString, nil
}

func (s *TokenService) ValidateAccessToken(tokenString string) (string, string, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New(errors.ErrInvalidInput, "unexpected signing method")
		}
		return s.accessSecret, nil
	})

	if err != nil {
		return "", "", err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		if claims["type"] != "access" {
			return "", "", errors.New(errors.ErrInvalidInput, "invalid token type")
		}
		return claims["sub"].(string), claims["email"].(string), nil
	}

	return "", "", errors.New(errors.ErrUnauthorized, "invalid token")
}

func (s *TokenService) ValidateRefreshToken(tokenString string) (string, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New(errors.ErrInvalidInput, "unexpected signing method")
		}
		return s.refreshSecret, nil
	})

	if err != nil {
		return "", err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		if claims["type"] != "refresh" {
			return "", errors.New(errors.ErrInvalidInput, "invalid token type")
		}
		return claims["sub"].(string), nil
	}

	return "", errors.New(errors.ErrUnauthorized, "invalid token")
}
