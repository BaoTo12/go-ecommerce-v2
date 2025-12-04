package application_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/titan-commerce/backend/order-service/internal/application"
	"github.com/titan-commerce/backend/order-service/internal/domain"
	"github.com/titan-commerce/backend/pkg/logger"
)

// MockRepository is a mock implementation of the order repository
type MockRepository struct {
	mock.Mock
}

func (m *MockRepository) Save(ctx context.Context, order *domain.Order) error {
	args := m.Called(ctx, order)
	return args.Error(0)
}

func (m *MockRepository) FindByID(ctx context.Context, orderID string) (*domain.Order, error) {
	args := m.Called(ctx, orderID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Order), args.Error(1)
}

func (m *MockRepository) FindByUserID(ctx context.Context, userID string, limit, offset int) ([]*domain.Order, error) {
	args := m.Called(ctx, userID, limit, offset)
	return args.Get(0).([]*domain.Order), args.Error(1)
}

func (m *MockRepository) Update(ctx context.Context, order *domain.Order) error {
	args := m.Called(ctx, order)
	return args.Error(0)
}

// MockEventStore is a mock implementation of the event store
type MockEventStore struct {
	mock.Mock
}

func (m *MockEventStore) SaveEvent(ctx context.Context, event interface{}) error {
	args := m.Called(ctx, event)
	return args.Error(0)
}

func TestOrderService_CreateOrder(t *testing.T) {
	// Setup
	mockRepo := new(MockRepository)
	mockEvents := new(MockEventStore)
	log := logger.New(logger.Config{Level: "debug", ServiceName: "test"})
	
	service := application.NewOrderService(mockRepo, mockEvents, log)

	ctx := context.Background()
	userID := "user-123"
	items := []domain.OrderItem{
		{ProductID: "prod-1", ProductName: "Test Product", Quantity: 2, UnitPrice: 29.99},
	}
	shippingAddress := "123 Test Street"

	// Expectations
	mockRepo.On("Save", ctx, mock.AnythingOfType("*domain.Order")).Return(nil)
	mockEvents.On("SaveEvent", ctx, mock.Anything).Return(nil)

	// Execute
	order, err := service.CreateOrder(ctx, userID, items, shippingAddress)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, order)
	assert.Equal(t, userID, order.UserID)
	assert.Equal(t, shippingAddress, order.ShippingAddress)
	assert.Len(t, order.Items, 1)
	assert.Equal(t, domain.OrderStatusPending, order.Status)
	
	mockRepo.AssertExpectations(t)
	mockEvents.AssertExpectations(t)
}

func TestOrderService_GetOrder(t *testing.T) {
	// Setup
	mockRepo := new(MockRepository)
	mockEvents := new(MockEventStore)
	log := logger.New(logger.Config{Level: "debug", ServiceName: "test"})
	
	service := application.NewOrderService(mockRepo, mockEvents, log)

	ctx := context.Background()
	orderID := "order-123"
	expectedOrder := &domain.Order{
		ID:     orderID,
		UserID: "user-123",
		Status: domain.OrderStatusPending,
	}

	// Expectations
	mockRepo.On("FindByID", ctx, orderID).Return(expectedOrder, nil)

	// Execute
	order, err := service.GetOrder(ctx, orderID)

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, expectedOrder, order)
	
	mockRepo.AssertExpectations(t)
}

func TestOrderService_CancelOrder(t *testing.T) {
	// Setup
	mockRepo := new(MockRepository)
	mockEvents := new(MockEventStore)
	log := logger.New(logger.Config{Level: "debug", ServiceName: "test"})
	
	service := application.NewOrderService(mockRepo, mockEvents, log)

	ctx := context.Background()
	orderID := "order-123"
	reason := "Customer requested cancellation"
	
	existingOrder := &domain.Order{
		ID:     orderID,
		UserID: "user-123",
		Status: domain.OrderStatusPending,
	}

	// Expectations
	mockRepo.On("FindByID", ctx, orderID).Return(existingOrder, nil)
	mockRepo.On("Update", ctx, mock.AnythingOfType("*domain.Order")).Return(nil)

	// Execute
	order, err := service.CancelOrder(ctx, orderID, reason)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, order)
	assert.Equal(t, domain.OrderStatusCancelled, order.Status)
	
	mockRepo.AssertExpectations(t)
}

func TestOrderService_CreateOrder_EmptyItems(t *testing.T) {
	// Setup
	mockRepo := new(MockRepository)
	mockEvents := new(MockEventStore)
	log := logger.New(logger.Config{Level: "debug", ServiceName: "test"})
	
	service := application.NewOrderService(mockRepo, mockEvents, log)

	ctx := context.Background()
	userID := "user-123"
	items := []domain.OrderItem{} // Empty
	shippingAddress := "123 Test Street"

	// Execute
	order, err := service.CreateOrder(ctx, userID, items, shippingAddress)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, order)
}
