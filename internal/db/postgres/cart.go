package postgres

import (
	cart_errors "cart-api/internal/errors"
	"cart-api/internal/models"
	"context"
	"database/sql"
	"errors"

	"github.com/jmoiron/sqlx"
)

// CartRepository provides methods to interact with the carts table in the database.
type CartRepository struct {
	db *sqlx.DB
}

// NewCartRepository creates a new instance of CartRepository.
func NewCartRepository(db *sqlx.DB) *CartRepository {
	return &CartRepository{db: db}
}

// Create inserts a new cart into the database and returns the created cart.
// The cart is initialized with an empty list of items.
func (r *CartRepository) Create(ctx context.Context) (*models.Cart, error) {
	var cart models.Cart
	query := `INSERT INTO carts DEFAULT VALUES RETURNING id`
	err := r.db.QueryRowxContext(ctx, query).Scan(&cart.ID)
	if err != nil {
		return nil, cart_errors.ErrFailedPostgresOpperation
	}
	cart.Items = []models.CartItem{}
	return &cart, nil
}

// Get retrieves a cart by its ID, including all associated cart items.
func (r *CartRepository) Get(ctx context.Context, id string) (*models.Cart, error) {
	var cart models.Cart
	query := `SELECT id FROM carts WHERE id = $1`
	err := r.db.GetContext(ctx, &cart, query, id)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, cart_errors.ErrCartDoesNotExist
	}
	if err != nil {
		return nil, cart_errors.ErrFailedToRetrieveCart
	}

	itemsQuery := `SELECT id, cart_id, product, quantity FROM cart_items WHERE cart_id = $1`
	err = r.db.SelectContext(ctx, &cart.Items, itemsQuery, id)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, cart_errors.ErrCartDoesNotExist
	}
	if err != nil {
		return nil, cart_errors.ErrFailedToRetrieveCartItems
	}
	return &cart, nil
}
