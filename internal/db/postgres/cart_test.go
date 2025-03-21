package postgres_test

import (
	"cart-api/internal/db/postgres"
	"context"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/assert"
)

func TestCreateCart(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' occurred when opening a stub database connection", err)
	}
	defer db.Close()

	sqlxDB := sqlx.NewDb(db, "sqlmock")
	repo := postgres.NewCartRepository(sqlxDB)

	mock.ExpectQuery(`INSERT INTO carts DEFAULT VALUES RETURNING id`).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow("cart-id"))

	cart, err := repo.Create(context.Background())
	assert.NoError(t, err)
	assert.NotNil(t, cart)
	assert.Equal(t, "cart-id", cart.ID)
	assert.Empty(t, cart.Items)

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestGetCart(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' occurred when opening a stub database connection", err)
	}
	defer db.Close()

	sqlxDB := sqlx.NewDb(db, "sqlmock")
	repo := postgres.NewCartRepository(sqlxDB)

	mock.ExpectQuery(`SELECT id FROM carts WHERE id = \$1`).
		WithArgs("cart-id").
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow("cart-id"))

	mock.ExpectQuery(`SELECT id, cart_id, product, quantity FROM cart_items WHERE cart_id = \$1`).
		WithArgs("cart-id").
		WillReturnRows(sqlmock.NewRows([]string{"id", "cart_id", "product", "quantity"}).
			AddRow("item-id", "cart-id", "product1", 2))

	cart, err := repo.Get(context.Background(), "cart-id")
	assert.NoError(t, err)
	assert.NotNil(t, cart)
	assert.Equal(t, "cart-id", cart.ID)
	assert.Equal(t, "product1", cart.Items[0].Product)

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}
