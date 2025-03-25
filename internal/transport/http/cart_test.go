package handler_test

import (
	"bytes"
	"cart-api/internal/model"
	handler "cart-api/internal/transport/http"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockCartService struct {
	mock.Mock
}

func (m *MockCartService) CreateCart(ctx context.Context) (*model.Cart, error) {
	args := m.Called(ctx)
	return args.Get(0).(*model.Cart), args.Error(1)
}

func (m *MockCartService) ViewCart(ctx context.Context, id string) (*model.Cart, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(*model.Cart), args.Error(1)
}

type MockCartItemService struct {
	mock.Mock
}

func (m *MockCartItemService) AddToCart(ctx context.Context, item *model.CartItem) error {
	args := m.Called(ctx, item)
	return args.Error(0)
}

func (m *MockCartItemService) RemoveFromCart(ctx context.Context, CartID, CartItemID string) error {
	args := m.Called(ctx, CartID, CartItemID)
	return args.Error(0)
}

func TestCreateCart(t *testing.T) {
	mockCartService := new(MockCartService)
	mockCartItemService := new(MockCartItemService)
	h := handler.NewCartHandler(mockCartService, mockCartItemService)

	cart := &model.Cart{ID: "123", Items: []model.CartItem{}}
	mockCartService.On("CreateCart", mock.Anything).Return(cart, nil)

	r := httptest.NewRequest(http.MethodPost, "/carts", nil)
	w := httptest.NewRecorder()

	h.CreateCart(w, r)
	res := w.Result()
	defer res.Body.Close()

	assert.Equal(t, http.StatusOK, res.StatusCode)
	var got model.Cart
	if err := json.NewDecoder(res.Body).Decode(&got); err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}
	assert.Equal(t, "123", got.ID)
}

func TestViewCart_NotFound(t *testing.T) {
	mockCartService := new(MockCartService)
	mockCartItemService := new(MockCartItemService)
	h := handler.NewCartHandler(mockCartService, mockCartItemService)

	mockCartService.On("ViewCart", mock.Anything, "123").Return((*model.Cart)(nil), errors.New("not found"))

	r := httptest.NewRequest(http.MethodGet, "/carts/123", nil)
	w := httptest.NewRecorder()

	h.ViewCart(w, r)
	res := w.Result()
	defer res.Body.Close()

	assert.Equal(t, http.StatusNotFound, res.StatusCode)
}

func TestAddToCart(t *testing.T) {
	mockCartService := new(MockCartService)
	mockCartItemService := new(MockCartItemService)
	h := handler.NewCartHandler(mockCartService, mockCartItemService)

	item := model.CartItem{CartID: "123", Product: "Apple", Quantity: 2}
	mockCartItemService.On("AddToCart", mock.Anything, &item).Return(nil)

	body, _ := json.Marshal(item)
	r := httptest.NewRequest(http.MethodPost, "/carts/123/items", bytes.NewReader(body))
	w := httptest.NewRecorder()

	h.AddToCart(w, r)
	res := w.Result()
	defer res.Body.Close()

	assert.Equal(t, http.StatusOK, res.StatusCode)
}

func TestRemoveFromCart(t *testing.T) {
	mockCartService := new(MockCartService)
	mockCartItemService := new(MockCartItemService)
	h := handler.NewCartHandler(mockCartService, mockCartItemService)

	mockCartItemService.On("RemoveFromCart", mock.Anything, "123", "456").Return(nil)

	r := httptest.NewRequest(http.MethodDelete, "/carts/123/items/456", nil)
	w := httptest.NewRecorder()

	h.RemoveFromCart(w, r)
	res := w.Result()
	defer res.Body.Close()

	assert.Equal(t, http.StatusNoContent, res.StatusCode)
}
