package application_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/titan-commerce/backend/cart-service/internal/application"
	"github.com/titan-commerce/backend/cart-service/internal/domain"
	"github.com/titan-commerce/backend/pkg/logger"
)

type MockCartRepository struct {
	mock.Mock
}

func (m *MockCartRepository) Save(ctx context.Context, cart *domain.Cart) error {
	args := m.Called(ctx, cart)
	return args.Error(0)
}

func (m *MockCartRepository) FindByUserID(ctx context.Context, userID string) (*domain.Cart, error) {
	args := m.Called(ctx, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Cart), args.Error(1)
}

func (m *MockCartRepository) Delete(ctx context.Context, userID string) error {
	args := m.Called(ctx, userID)
	return args.Error(0)
}

func TestCartService_AddToCart(t *testing.T) {
	// Setup
	mockRepo := new(MockCartRepository)
	log := logger.New(logger.Config{Level: "debug", ServiceName: "test"})
	
	service := application.NewCartService(mockRepo, log)

	ctx := context.Background()
	userID := "user-123"
	productID := "prod-123"
	productName := "Test Product"
	quantity := 2
	unitPrice := 29.99

	existingCart := domain.NewCart(userID)

	// Expectations
	mockRepo.On("FindByUserID", ctx, userID).Return(existingCart, nil)
	mockRepo.On("Save", ctx, mock.AnythingOfType("*domain.Cart")).Return(nil)

	// Execute
	cart, err := service.AddToCart(ctx, userID, productID, productName, quantity, unitPrice)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, cart)
	assert.Len(t, cart.Items, 1)
	assert.Equal(t, productID, cart.Items[0].ProductID)
	assert.Equal(t, quantity, cart.Items[0].Quantity)
	
	mockRepo.AssertExpectations(t)
}

func TestCartService_AddToCart_ExistingItem(t *testing.T) {
	// Setup
	mockRepo := new(MockCartRepository)
	log := logger.New(logger.Config{Level: "debug", ServiceName: "test"})
	
	service := application.NewCartService(mockRepo, log)

	ctx := context.Background()
	userID := "user-123"
	productID := "prod-123"
	productName := "Test Product"
	unitPrice := 29.99

	existingCart := domain.NewCart(userID)
	existingCart.AddItem(productID, productName, 1, unitPrice)

	// Expectations
	mockRepo.On("FindByUserID", ctx, userID).Return(existingCart, nil)
	mockRepo.On("Save", ctx, mock.AnythingOfType("*domain.Cart")).Return(nil)

	// Execute - add 2 more of the same product
	cart, err := service.AddToCart(ctx, userID, productID, productName, 2, unitPrice)

	// Assert
	assert.NoError(t, err)
	assert.Len(t, cart.Items, 1) // Still one item, quantity increased
	assert.Equal(t, 3, cart.Items[0].Quantity) // 1 + 2 = 3
	
	mockRepo.AssertExpectations(t)
}

func TestCartService_RemoveFromCart(t *testing.T) {
	// Setup
	mockRepo := new(MockCartRepository)
	log := logger.New(logger.Config{Level: "debug", ServiceName: "test"})
	
	service := application.NewCartService(mockRepo, log)

	ctx := context.Background()
	userID := "user-123"
	productID := "prod-123"

	existingCart := domain.NewCart(userID)
	existingCart.AddItem(productID, "Test Product", 2, 29.99)
	existingCart.AddItem("prod-456", "Other Product", 1, 19.99)

	// Expectations
	mockRepo.On("FindByUserID", ctx, userID).Return(existingCart, nil)
	mockRepo.On("Save", ctx, mock.AnythingOfType("*domain.Cart")).Return(nil)

	// Execute
	cart, err := service.RemoveFromCart(ctx, userID, productID)

	// Assert
	assert.NoError(t, err)
	assert.Len(t, cart.Items, 1) // Only one item left
	assert.Equal(t, "prod-456", cart.Items[0].ProductID)
	
	mockRepo.AssertExpectations(t)
}

func TestCartService_GetCart(t *testing.T) {
	// Setup
	mockRepo := new(MockCartRepository)
	log := logger.New(logger.Config{Level: "debug", ServiceName: "test"})
	
	service := application.NewCartService(mockRepo, log)

	ctx := context.Background()
	userID := "user-123"

	expectedCart := domain.NewCart(userID)
	expectedCart.AddItem("prod-123", "Test Product", 2, 29.99)

	// Expectations
	mockRepo.On("FindByUserID", ctx, userID).Return(expectedCart, nil)

	// Execute
	cart, err := service.GetCart(ctx, userID)

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, expectedCart, cart)
	
	mockRepo.AssertExpectations(t)
}

func TestCartService_ClearCart(t *testing.T) {
	// Setup
	mockRepo := new(MockCartRepository)
	log := logger.New(logger.Config{Level: "debug", ServiceName: "test"})
	
	service := application.NewCartService(mockRepo, log)

	ctx := context.Background()
	userID := "user-123"

	// Expectations
	mockRepo.On("Delete", ctx, userID).Return(nil)

	// Execute
	err := service.ClearCart(ctx, userID)

	// Assert
	assert.NoError(t, err)
	
	mockRepo.AssertExpectations(t)
}

func TestCartService_UpdateQuantity(t *testing.T) {
	// Setup
	mockRepo := new(MockCartRepository)
	log := logger.New(logger.Config{Level: "debug", ServiceName: "test"})
	
	service := application.NewCartService(mockRepo, log)

	ctx := context.Background()
	userID := "user-123"
	productID := "prod-123"

	existingCart := domain.NewCart(userID)
	existingCart.AddItem(productID, "Test Product", 2, 29.99)

	// Expectations
	mockRepo.On("FindByUserID", ctx, userID).Return(existingCart, nil)
	mockRepo.On("Save", ctx, mock.AnythingOfType("*domain.Cart")).Return(nil)

	// Execute - update quantity to 5
	cart, err := service.UpdateQuantity(ctx, userID, productID, 5)

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, 5, cart.Items[0].Quantity)
	
	mockRepo.AssertExpectations(t)
}
