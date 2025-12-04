package domain

import (
	"time"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID           string
	Email        string
	PasswordHash string
	FullName     string
	MFAEnabled   bool
	MFASecret    string
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

func NewUser(email, password, fullName string) (*User, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	return &User{
		ID:           uuid.New().String(),
		Email:        email,
		PasswordHash: string(hash),
		FullName:     fullName,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}, nil
}

func (u *User) CheckPassword(password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(u.PasswordHash), []byte(password))
	return err == nil
}
