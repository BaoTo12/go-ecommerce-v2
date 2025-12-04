package domain

import (
	"time"

	"github.com/google/uuid"
	"github.com/titan-commerce/backend/pkg/errors"
)

type User struct {
	ID          string
	Email       string
	Name        string
	FullName    string // Alias for repository compatibility
	Phone       string
	PhoneNumber string // Alias for repository compatibility
	AvatarURL   string
	Addresses   []Address
	Preferences UserPreferences
	CreatedAt   time.Time
	UpdatedAt   time.Time
	Version     int
}

type UserPreferences struct {
	EmailNotifications bool
	SMSNotifications   bool
	PreferredLanguage  string
}

type Address struct {
	ID            string
	UserID        string
	FullName      string
	Phone         string
	Street        string
	StreetAddress string // Alias for repository
	City          string
	State         string
	ZipCode       string
	PostalCode    string // Alias for repository
	Country       string
	IsDefault     bool
	CreatedAt     time.Time
}

func NewAddress(userID, fullName, phone, street, city, postalCode, country string, isDefault bool) (*Address, error) {
	if userID == "" {
		return nil, errors.New(errors.ErrInvalidInput, "user ID is required")
	}
	if fullName == "" {
		return nil, errors.New(errors.ErrInvalidInput, "full name is required")
	}
	if street == "" {
		return nil, errors.New(errors.ErrInvalidInput, "street address is required")
	}

	return &Address{
		ID:            uuid.New().String(),
		UserID:        userID,
		FullName:      fullName,
		Phone:         phone,
		Street:        street,
		StreetAddress: street,
		City:          city,
		PostalCode:    postalCode,
		ZipCode:       postalCode,
		Country:       country,
		IsDefault:     isDefault,
		CreatedAt:     time.Now(),
	}, nil
}

func (u *User) UpdateProfile(name, phone, avatarURL string) {
	if name != "" {
		u.Name = name
		u.FullName = name
	}
	if phone != "" {
		u.Phone = phone
		u.PhoneNumber = phone
	}
	if avatarURL != "" {
		u.AvatarURL = avatarURL
	}
	u.UpdatedAt = time.Now()
}

func (u *User) UpdatePreferences(emailNotifs, smsNotifs bool, language string) {
	u.Preferences.EmailNotifications = emailNotifs
	u.Preferences.SMSNotifications = smsNotifs
	if language != "" {
		u.Preferences.PreferredLanguage = language
	}
	u.UpdatedAt = time.Now()
}
