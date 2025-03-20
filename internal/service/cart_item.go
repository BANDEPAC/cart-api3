package service

import (
	"cart-api/internal/models"
	"context"
)

// CartItemStorage defines the interface for interacting with cart item storage.
type CartItemStorage interface {
	Create(ctx context.Context, item *models.CartItem) error
	Delete(ctx context.Context, CartID, CartItemID string) error
}

// CartItemService provides business logic for managing cart items.
type CartItemService struct {
	repo CartItemStorage
}

// NewCartItemRepository creates a new instance of CartItemService.
// It accepts a CartItemStorage implementation as a dependency.
func NewCartItemRepository(repo CartItemStorage) *CartItemService {
	return &CartItemService{repo: repo}
}

// AddToCart adds a new item to the cart.
// It delegates the operation to the underlying storage.
func (s CartItemService) AddToCart(ctx context.Context, item *models.CartItem) error {
	err := s.repo.Create(ctx, item)
	if err != nil {
		return err
	}
	return nil
}

// RemoveFromCart removes an item from the cart by its ID and cart ID.
// It delegates the operation to the underlying storage.
func (s CartItemService) RemoveFromCart(ctx context.Context, CartID, CartItemID string) error {
	err := s.repo.Delete(ctx, CartID, CartItemID)
	if err != nil {
		return err
	}
	return err
}
