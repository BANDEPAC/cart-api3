package postgres

import (
	"cart-api/internal/carterror"
	"cart-api/internal/model"
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
func (r *CartRepository) Create(ctx context.Context) (*model.Cart, error) {
	var cart model.Cart
	query := `INSERT INTO carts DEFAULT VALUES RETURNING id`
	err := r.db.QueryRowxContext(ctx, query).Scan(&cart.ID)
	if err != nil {
		return nil, carterror.ErrFailedPostgresOpperation
	}
	cart.Items = []model.CartItem{}
	return &cart, nil
}

// Get retrieves a cart by its ID, including all associated cart items.
func (r *CartRepository) Get(ctx context.Context, id string) (*model.Cart, error) {
	var cart model.Cart
	query := `SELECT id FROM carts WHERE id = $1`
	err := r.db.GetContext(ctx, &cart, query, id)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, carterror.ErrCartDoesNotExist
	}
	if err != nil {
		return nil, carterror.ErrFailedToRetrieveCart
	}

	itemsQuery := `SELECT id, cart_id, product, quantity FROM cart_items WHERE cart_id = $1`
	err = r.db.SelectContext(ctx, &cart.Items, itemsQuery, id)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, carterror.ErrCartDoesNotExist
	}
	if err != nil {
		return nil, carterror.ErrFailedToRetrieveCartItems
	}
	return &cart, nil
}
