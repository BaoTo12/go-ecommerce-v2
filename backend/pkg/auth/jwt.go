package pkg

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/titan-commerce/backend/pkg/errors"
)

type JWTClaims struct {
	UserID string   `json:"user_id"`
	Email  string   `json:"email"`
	Roles  []string `json:"roles"`
	CellID string   `json:"cell_id"`
	jwt.RegisteredClaims
}

type JWTService struct {
	secretKey       string
	accessTokenTTL  time.Duration
	refreshTokenTTL time.Duration
}

func NewJWTService(secretKey string, accessMinutes, refreshDays int) *JWTService {
	return &JWTService{
		secretKey:       secretKey,
		accessTokenTTL:  time.Duration(accessMinutes) * time.Minute,
		refreshTokenTTL: time.Duration(refreshDays) * 24 * time.Hour,
	}
}

func (s *JWTService) GenerateAccessToken(userID, email, cellID string, roles []string) (string, error) {
	claims := &JWTClaims{
		UserID: userID,
		Email:  email,
		Roles:  roles,
		CellID: cellID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(s.accessTokenTTL)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(s.secretKey))
	if err != nil {
		return "", errors.Wrap(errors.ErrInternal, "failed to sign token", err)
	}

	return tokenString, nil
}

func (s *JWTService) VerifyAccessToken(tokenString string) (*JWTClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &JWTClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New(errors.ErrUnauthorized, "invalid signing method")
		}
		return []byte(s.secretKey), nil
	})

	if err != nil {
		return nil, errors.Wrap(errors.ErrUnauthorized, "invalid token", err)
	}

	if claims, ok := token.Claims.(*JWTClaims); ok && token.Valid {
		return claims, nil
	}

	return nil, errors.New(errors.ErrUnauthorized, "invalid token claims")
}
