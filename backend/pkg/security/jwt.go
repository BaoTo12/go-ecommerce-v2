package security

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"errors"
	"strings"
	"time"
)

var (
	ErrInvalidToken     = errors.New("invalid token")
	ErrTokenExpired     = errors.New("token expired")
	ErrInvalidSignature = errors.New("invalid signature")
)

// JWTClaims represents JWT claims
type JWTClaims struct {
	Subject   string            `json:"sub"`
	Issuer    string            `json:"iss,omitempty"`
	Audience  string            `json:"aud,omitempty"`
	ExpiresAt int64             `json:"exp"`
	IssuedAt  int64             `json:"iat"`
	NotBefore int64             `json:"nbf,omitempty"`
	TokenID   string            `json:"jti,omitempty"`
	Role      string            `json:"role,omitempty"`
	Extra     map[string]string `json:"extra,omitempty"`
}

// JWTConfig configures the JWT service
type JWTConfig struct {
	Secret           []byte
	Issuer           string
	Audience         string
	AccessTokenTTL   time.Duration
	RefreshTokenTTL  time.Duration
}

// DefaultJWTConfig returns secure defaults
func DefaultJWTConfig(secret []byte) *JWTConfig {
	return &JWTConfig{
		Secret:           secret,
		Issuer:           "titan-commerce",
		Audience:         "titan-commerce-api",
		AccessTokenTTL:   15 * time.Minute,
		RefreshTokenTTL:  7 * 24 * time.Hour,
	}
}

// JWTService handles JWT operations
type JWTService struct {
	config *JWTConfig
}

// NewJWTService creates a JWT service
func NewJWTService(config *JWTConfig) *JWTService {
	return &JWTService{config: config}
}

// TokenPair represents access and refresh tokens
type TokenPair struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	ExpiresIn    int64  `json:"expires_in"`
	TokenType    string `json:"token_type"`
}

// GenerateTokenPair creates access and refresh tokens
func (j *JWTService) GenerateTokenPair(userID, role string, extra map[string]string) (*TokenPair, error) {
	now := time.Now()

	// Access token
	accessClaims := &JWTClaims{
		Subject:   userID,
		Issuer:    j.config.Issuer,
		Audience:  j.config.Audience,
		IssuedAt:  now.Unix(),
		ExpiresAt: now.Add(j.config.AccessTokenTTL).Unix(),
		Role:      role,
		Extra:     extra,
	}

	accessToken, err := j.sign(accessClaims)
	if err != nil {
		return nil, err
	}

	// Refresh token (longer lived, minimal claims)
	refreshClaims := &JWTClaims{
		Subject:   userID,
		IssuedAt:  now.Unix(),
		ExpiresAt: now.Add(j.config.RefreshTokenTTL).Unix(),
		TokenID:   generateTokenID(),
	}

	refreshToken, err := j.sign(refreshClaims)
	if err != nil {
		return nil, err
	}

	return &TokenPair{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		ExpiresIn:    int64(j.config.AccessTokenTTL.Seconds()),
		TokenType:    "Bearer",
	}, nil
}

// Verify validates a token and returns claims
func (j *JWTService) Verify(token string) (*JWTClaims, error) {
	parts := strings.Split(token, ".")
	if len(parts) != 3 {
		return nil, ErrInvalidToken
	}

	// Verify signature
	signingInput := parts[0] + "." + parts[1]
	signature, err := base64.RawURLEncoding.DecodeString(parts[2])
	if err != nil {
		return nil, ErrInvalidToken
	}

	expectedSig := j.computeSignature(signingInput)
	if !hmac.Equal(signature, expectedSig) {
		return nil, ErrInvalidSignature
	}

	// Decode claims
	claimsBytes, err := base64.RawURLEncoding.DecodeString(parts[1])
	if err != nil {
		return nil, ErrInvalidToken
	}

	var claims JWTClaims
	if err := json.Unmarshal(claimsBytes, &claims); err != nil {
		return nil, ErrInvalidToken
	}

	// Check expiration
	if claims.ExpiresAt < time.Now().Unix() {
		return nil, ErrTokenExpired
	}

	// Check not before
	if claims.NotBefore > 0 && claims.NotBefore > time.Now().Unix() {
		return nil, ErrInvalidToken
	}

	return &claims, nil
}

// RefreshToken generates new token pair from refresh token
func (j *JWTService) RefreshToken(refreshToken string) (*TokenPair, error) {
	claims, err := j.Verify(refreshToken)
	if err != nil {
		return nil, err
	}

	return j.GenerateTokenPair(claims.Subject, claims.Role, claims.Extra)
}

func (j *JWTService) sign(claims *JWTClaims) (string, error) {
	// Header
	header := map[string]string{
		"alg": "HS256",
		"typ": "JWT",
	}
	headerBytes, _ := json.Marshal(header)
	headerB64 := base64.RawURLEncoding.EncodeToString(headerBytes)

	// Claims
	claimsBytes, err := json.Marshal(claims)
	if err != nil {
		return "", err
	}
	claimsB64 := base64.RawURLEncoding.EncodeToString(claimsBytes)

	// Signature
	signingInput := headerB64 + "." + claimsB64
	signature := j.computeSignature(signingInput)
	signatureB64 := base64.RawURLEncoding.EncodeToString(signature)

	return signingInput + "." + signatureB64, nil
}

func (j *JWTService) computeSignature(input string) []byte {
	h := hmac.New(sha256.New, j.config.Secret)
	h.Write([]byte(input))
	return h.Sum(nil)
}

func generateTokenID() string {
	token, _ := GenerateSecureToken(16)
	return token
}

// RefreshTokenStore stores refresh tokens for rotation
type RefreshTokenStore interface {
	Store(tokenID, userID string, expiresAt time.Time) error
	Verify(tokenID string) (userID string, valid bool, err error)
	Revoke(tokenID string) error
	RevokeAllForUser(userID string) error
}
