package application

import (
	"context"

	"github.com/titan-commerce/backend/user-service/internal/domain"
	"github.com/titan-commerce/backend/pkg/logger"
)

type UserRepository interface {
	FindByID(ctx context.Context, userID string) (*domain.User, error)
	Save(ctx context.Context, user *domain.User) error
	Update(ctx context.Context, user *domain.User) error
}

type AddressRepository interface {
	Save(ctx context.Context, address *domain.Address) error
	FindByUserID(ctx context.Context, userID string) ([]*domain.Address, error)
	Update(ctx context.Context, address *domain.Address) error
}

type UserService struct {
	userRepo    UserRepository
	addressRepo AddressRepository
	logger      *logger.Logger
}

// NewUserService creates a new user service (with address repository)
func NewUserService(userRepo UserRepository, addressRepo AddressRepository, logger *logger.Logger) *UserService {
	return &UserService{
		userRepo:    userRepo,
		addressRepo: addressRepo,
		logger:      logger,
	}
}

// NewUserServiceSimple creates a user service without address repository
func NewUserServiceSimple(userRepo UserRepository, logger *logger.Logger) *UserService {
	return &UserService{
		userRepo: userRepo,
		logger:   logger,
	}
}

// GetUser retrieves user profile (Query)
func (s *UserService) GetUser(ctx context.Context, userID string) (*domain.User, error) {
	return s.userRepo.FindByID(ctx, userID)
}

// UpdateProfile updates user profile information (Command)
func (s *UserService) UpdateProfile(ctx context.Context, userID, name, phone, avatarURL string) (*domain.User, error) {
	user, err := s.userRepo.FindByID(ctx, userID)
	if err != nil {
		return nil, err
	}

	user.UpdateProfile(name, phone, avatarURL)

	if err := s.userRepo.Update(ctx, user); err != nil {
		s.logger.Error(err, "failed to update user profile")
		return nil, err
	}

	s.logger.Infof("Profile updated: user=%s", userID)
	return user, nil
}

// AddAddress adds a new shipping address (Command)
func (s *UserService) AddAddress(ctx context.Context, userID, fullName, phone, street, city, postalCode, country string, isDefault bool) (*domain.Address, error) {
	// If address repo not available, return error
	if s.addressRepo == nil {
		address, err := domain.NewAddress(userID, fullName, phone, street, city, postalCode, country, isDefault)
		if err != nil {
			return nil, err
		}
		return address, nil
	}

	// If this is default, unset other defaults
	if isDefault {
		addresses, _ := s.addressRepo.FindByUserID(ctx, userID)
		for _, addr := range addresses {
			if addr.IsDefault {
				addr.IsDefault = false
				s.addressRepo.Update(ctx, addr)
			}
		}
	}

	address, err := domain.NewAddress(userID, fullName, phone, street, city, postalCode, country, isDefault)
	if err != nil {
		return nil, err
	}

	if err := s.addressRepo.Save(ctx, address); err != nil {
		s.logger.Error(err, "failed to save address")
		return nil, err
	}

	s.logger.Infof("Address added: user=%s", userID)
	return address, nil
}

// GetAddresses retrieves all addresses for a user (Query)
func (s *UserService) GetAddresses(ctx context.Context, userID string) ([]*domain.Address, error) {
	if s.addressRepo == nil {
		return nil, nil
	}
	return s.addressRepo.FindByUserID(ctx, userID)
}

// UpdatePreferences updates user notification preferences (Command)
func (s *UserService) UpdatePreferences(ctx context.Context, userID string, emailNotifs, smsNotifs bool, language string) error {
	user, err := s.userRepo.FindByID(ctx, userID)
	if err != nil {
		return err
	}

	user.UpdatePreferences(emailNotifs, smsNotifs, language)

	if err := s.userRepo.Update(ctx, user); err != nil {
		return err
	}

	s.logger.Infof("Preferences updated: user=%s", userID)
	return nil
}
