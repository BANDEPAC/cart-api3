package service_test

import (
	"cart-api/internal/models"
	"cart-api/internal/service"
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

type mockCartItemStorage struct {
	createErr error
	deleteErr error
}

func (m *mockCartItemStorage) Create(ctx context.Context, item *models.CartItem) error {
	return m.createErr
}

func (m *mockCartItemStorage) Delete(ctx context.Context, cartID, cartItemID string) error {
	return m.deleteErr
}

func TestAddToCart_Success(t *testing.T) {
	mockRepo := &mockCartItemStorage{
		createErr: nil,
	}

	service := service.NewCartItemRepository(mockRepo)

	item := &models.CartItem{
		CartID:   "cart-id",
		Product:  "product1",
		Quantity: 2,
	}

	err := service.AddToCart(context.Background(), item)
	assert.NoError(t, err)
}

func TestAddToCart_Error(t *testing.T) {
	mockRepo := &mockCartItemStorage{
		createErr: errors.New("failed to add item to cart"),
	}

	service := service.NewCartItemRepository(mockRepo)

	item := &models.CartItem{
		CartID:   "cart-id",
		Product:  "product1",
		Quantity: 2,
	}

	err := service.AddToCart(context.Background(), item)
	assert.Error(t, err)
	assert.Equal(t, "failed to add item to cart", err.Error())
}

func TestRemoveFromCart_Success(t *testing.T) {
	mockRepo := &mockCartItemStorage{
		deleteErr: nil,
	}

	service := service.NewCartItemRepository(mockRepo)

	err := service.RemoveFromCart(context.Background(), "cart-id", "item-id")
	assert.NoError(t, err)
}

func TestRemoveFromCart_Error(t *testing.T) {
	mockRepo := &mockCartItemStorage{
		deleteErr: errors.New("failed to remove item from cart"),
	}

	service := service.NewCartItemRepository(mockRepo)

	err := service.RemoveFromCart(context.Background(), "cart-id", "item-id")
	assert.Error(t, err)
	assert.Equal(t, "failed to remove item from cart", err.Error())
}
