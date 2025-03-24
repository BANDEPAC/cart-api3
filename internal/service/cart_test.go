package service_test

import (
	"cart-api/internal/model"
	"cart-api/internal/service"
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

type mockCartStorage struct {
	createCartResult *model.Cart
	createCartErr    error
	getCartResult    *model.Cart
	getCartErr       error
}

func (m *mockCartStorage) Create(ctx context.Context) (*model.Cart, error) {
	return m.createCartResult, m.createCartErr
}

func (m *mockCartStorage) Get(ctx context.Context, id string) (*model.Cart, error) {
	return m.getCartResult, m.getCartErr
}

func TestCreateCart_Success(t *testing.T) {
	mockRepo := &mockCartStorage{
		createCartResult: &model.Cart{ID: "cart-id"},
		createCartErr:    nil,
	}

	service := service.NewCartService(mockRepo)

	cart, err := service.CreateCart(context.Background())
	assert.NoError(t, err)
	assert.NotNil(t, cart)
	assert.Equal(t, "cart-id", cart.ID)
}

func TestCreateCart_Error(t *testing.T) {
	mockRepo := &mockCartStorage{
		createCartResult: nil,
		createCartErr:    errors.New("failed to create cart"),
	}

	service := service.NewCartService(mockRepo)

	cart, err := service.CreateCart(context.Background())
	assert.Error(t, err)
	assert.Nil(t, cart)
	assert.Equal(t, "failed to create cart", err.Error())
}

func TestViewCart_Success(t *testing.T) {
	mockRepo := &mockCartStorage{
		getCartResult: &model.Cart{ID: "cart-id"},
		getCartErr:    nil,
	}

	service := service.NewCartService(mockRepo)

	cart, err := service.ViewCart(context.Background(), "cart-id")
	assert.NoError(t, err)
	assert.NotNil(t, cart)
	assert.Equal(t, "cart-id", cart.ID)
}

func TestViewCart_Error(t *testing.T) {
	mockRepo := &mockCartStorage{
		getCartResult: nil,
		getCartErr:    errors.New("cart not found"),
	}

	service := service.NewCartService(mockRepo)

	cart, err := service.ViewCart(context.Background(), "cart-id")
	assert.Error(t, err)
	assert.Nil(t, cart)
	assert.Equal(t, "cart not found", err.Error())
}
