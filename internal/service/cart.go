package service

import (
	"cart-api/internal/model"
	"context"
)

// CartStorage defines the interface for interacting with cart storage
type CartStorage interface {
	Create(ctx context.Context) (*model.Cart, error)
	Get(ctx context.Context, id string) (*model.Cart, error)
}

// CartService provides business logic for managing carts.
type CartService struct {
	repo CartStorage
}

// NewCartRepository creates a new instance of CartService.
// It accepts a CartStorage implementation as a dependency.
func NewCartService(repo CartStorage) *CartService {
	return &CartService{repo: repo}
}

// CreateCart creates a new cart
// It delegates the operation to the underlying storage.
func (s *CartService) CreateCart(ctx context.Context) (*model.Cart, error) {
	cart, err := s.repo.Create(ctx)
	if err != nil {
		return nil, err
	}
	return cart, nil

}

// CreateCart retrieves a cart by its ID,including all associated items.
// It delegates the operation to the underlying storage.
func (s *CartService) ViewCart(ctx context.Context, id string) (*model.Cart, error) {
	cart, err := s.repo.Get(ctx, id)
	if err != nil {
		return nil, err
	}
	return cart, nil
}
