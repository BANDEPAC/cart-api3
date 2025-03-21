package postgres

import (
	cart_errors "cart-api/internal/errors"
	"cart-api/internal/models"
	"context"
	"log"

	"github.com/jmoiron/sqlx"
)

// CartItemRepository provides methods to interact with the cart_items table in the database.
type CartItemRepository struct {
	db *sqlx.DB
}

// NewCartItemRepository creates a new instance of CartItemRepository.
func NewCartItemRepository(db *sqlx.DB) *CartItemRepository {
	return &CartItemRepository{db: db}
}

// Create inserts a new cart item into the database.
// It returns an error if the operation fails.
func (r *CartItemRepository) Create(ctx context.Context, item *models.CartItem) error {
	query := `INSERT INTO cart_items (cart_id, product, quantity) VALUES ($1, $2, $3) RETURNING id`
	var id string
	err := r.db.QueryRowContext(ctx, query, item.CartID, item.Product, item.Quantity).Scan(&id)
	if err != nil {
		return cart_errors.ErrFailedPostgresOpperation
	}

	item.ID = id
	return nil
}

// CartExists checks if a cart with the given ID exists in the database.
// It returns a boolean indicating existence and an error if the query fails.
func (r *CartItemRepository) CartExists(ctx context.Context, cartID string) (bool, error) {
	var exists bool
	query := `SELECT EXISTS(SELECT 1 FROM carts WHERE id = $1)`
	err := r.db.QueryRowxContext(ctx, query, cartID).Scan(&exists)
	if err != nil {
		return false, cart_errors.ErrFailedPostgresOpperation
	}
	return exists, nil
}

// Delete removes a cart item from the database by its ID and cart ID.
// It returns an error if the cart does not exist or the operation fails.
func (r *CartItemRepository) Delete(ctx context.Context, cartID, cartItemID string) error {

	exists, err := r.CartExists(ctx, cartID)
	if err != nil {
		return err
	}
	if !exists {
		return cart_errors.ErrNotFound
	}
	log.Println("id : ", cartItemID, "cart_id : ", cartID)
	query := `DELETE FROM cart_items WHERE id = $1 AND cart_id = $2`
	_, err = r.db.ExecContext(ctx, query, cartItemID, cartID)
	return err
}
