package handler

import (
	"cart-api/internal/carterror"
	"cart-api/internal/model"
	"context"
	"encoding/json"
	"log"
	"net/http"
	"strings"
)

// CartService defines the interface for cart-related operations.
type CartService interface {
	CreateCart(ctx context.Context) (*model.Cart, error)
	ViewCart(ctx context.Context, cartID string) (*model.Cart, error)
}

// CartItemService defines the interface for cart item-related operations.
type CartItemService interface {
	AddToCart(ctx context.Context, item *model.CartItem) error
	RemoveFromCart(ctx context.Context, cartID, itemID string) error
}

// CartHandler provides HTTP handlers for cart and cart item operations.
type CartHandler struct {
	cartService     CartService
	cartItemService CartItemService
}

// NewCartHandler creates a new instance of CartHandler.
func NewCartHandler(c CartService, ci CartItemService) *CartHandler {
	return &CartHandler{cartService: c, cartItemService: ci}
}

// CreateCart handles the creation of a new cart.
func (h *CartHandler) CreateCart(w http.ResponseWriter, r *http.Request) {
	log.Println("CreateCart is called")
	if r.Method != http.MethodPost {
		http.Error(w, carterror.ErrInvalidRequestMethod.Error(), http.StatusMethodNotAllowed)
		return
	}
	cart, err := h.cartService.CreateCart(r.Context())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(cart); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

// ViewCart handles the retrieval of a cart by its ID.
func (h *CartHandler) ViewCart(w http.ResponseWriter, r *http.Request) {
	log.Println("ViewCart is called")
	if r.Method != http.MethodGet {
		http.Error(w, carterror.ErrInvalidRequestMethod.Error(), http.StatusMethodNotAllowed)
		return
	}
	cartID := r.URL.Path[len("/carts/"):]
	log.Println("Extracted id: ", cartID)

	cart, err := h.cartService.ViewCart(r.Context(), cartID)
	if err != nil {
		http.Error(w, carterror.ErrCartDoesNotExist.Error(), http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(cart); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

// AddToCart handles the addition of an item to the cart.
func (h *CartHandler) AddToCart(w http.ResponseWriter, r *http.Request) {
	log.Println("AddToCart is called")
	if r.Method != http.MethodPost {
		http.Error(w, carterror.ErrInvalidRequestMethod.Error(), http.StatusMethodNotAllowed)
		return
	}
	cartID := strings.Split(r.URL.Path, "/")[2]
	log.Println("Extracted id: ", cartID)

	if cartID == "" {
		http.Error(w, carterror.ErrInvalidQuery.Error(), http.StatusBadRequest)
		return
	}

	var request struct {
		Product  string `json:"product"`
		Quantity int    `json:"quantity"`
	}

	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, carterror.ErrInvalidRequestBody.Error(), http.StatusBadRequest)
		return
	}

	item := model.CartItem{ID: "", CartID: cartID, Product: request.Product, Quantity: request.Quantity}
	err := h.cartItemService.AddToCart(r.Context(), &item)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(item); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

// RemoveFromCart handles the removal of an item from the cart.
func (h *CartHandler) RemoveFromCart(w http.ResponseWriter, r *http.Request) {
	log.Println("RemoveFromCart is called")
	if r.Method != http.MethodDelete {
		http.Error(w, carterror.ErrInvalidRequestMethod.Error(), http.StatusMethodNotAllowed)
		return
	}
	paths := strings.Split(r.URL.Path, "/")
	cartID := paths[2]
	itemID := paths[4]
	log.Println("Extracted cartID: ", cartID)
	log.Println("Extracted itemID: ", itemID)

	if cartID == "" {
		http.Error(w, carterror.ErrCartIDRequired.Error(), http.StatusBadRequest)
		return
	}
	if itemID == "" {
		http.Error(w, carterror.ErrItemIDRequired.Error(), http.StatusBadRequest)
		return
	}

	err := h.cartItemService.RemoveFromCart(r.Context(), cartID, itemID)
	if err != nil {
		if err.Error() == "cart does not exist" {
			http.Error(w, carterror.ErrCartDoesNotExist.Error(), http.StatusNotFound)
			return
		}
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
