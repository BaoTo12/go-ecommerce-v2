package domain

import (
	"time"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"github.com/titan-commerce/backend/pkg/errors"
)

type User struct {
	ID           string
	Email        string
	PasswordHash string
	Name         string
	Roles        []string
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

func NewUser(email, password, name string) (*User, error) {
	if email == "" {
		return nil, errors.New(errors.ErrInvalidInput, "email is required")
	}
	if password == "" || len(password) < 8 {
		return nil, errors.New(errors.ErrInvalidInput, "password must be at least 8 characters")
	}
	if name == "" {
		return nil, errors.New(errors.ErrInvalidInput, "name is required")
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, errors.Wrap(errors.ErrInternal, "failed to hash password", err)
	}

	now := time.Now()
	return &User{
		ID:           uuid.New().String(),
		Email:        email,
		PasswordHash: string(hash),
		Name:         name,
		Roles:        []string{"customer"},
		CreatedAt:    now,
		UpdatedAt:    now,
	}, nil
}

func (u *User) VerifyPassword(password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(u.PasswordHash), []byte(password))
	return err == nil
}

func (u *User) AddRole(role string) {
	for _, r := range u.Roles {
		if r == role {
			return
		}
	}
	u.Roles = append(u.Roles, role)
	u.UpdatedAt = time.Now()
}

type RefreshToken struct {
	TokenID   string
	UserID    string
	Token     string
	ExpiresAt time.Time
	CreatedAt time.Time
}

func NewRefreshToken(userID string, expiryDays int) *RefreshToken {
	return &RefreshToken{
		TokenID:   uuid.New().String(),
		UserID:    userID,
		Token:     uuid.New().String(),
		ExpiresAt: time.Now().Add(time.Duration(expiryDays) * 24 * time.Hour),
		CreatedAt: time.Now(),
	}
}

func (t *RefreshToken) IsExpired() bool {
	return time.Now().After(t.ExpiresAt)
}
