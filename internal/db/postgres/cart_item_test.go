package postgres_test

import (
	"cart-api/internal/db/postgres"
	"cart-api/internal/model"
	"context"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/assert"
)

func TestCreateCartItem(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to open mock database: %s", err)
	}
	defer db.Close()

	sqlxDB := sqlx.NewDb(db, "sqlmock")
	repo := postgres.NewCartItemRepository(sqlxDB)

	item := &model.CartItem{
		CartID:   "cart-id",
		Product:  "product1",
		Quantity: 2,
	}

	mock.ExpectQuery(`INSERT INTO cart_items \(cart_id, product, quantity\) VALUES \(\$1, \$2, \$3\) RETURNING id`).
		WithArgs(item.CartID, item.Product, item.Quantity).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow("item-id"))

	err = repo.Create(context.Background(), item)
	assert.NoError(t, err)
	assert.Equal(t, "item-id", item.ID)

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestCartExists(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to open mock database: %s", err)
	}
	defer db.Close()

	sqlxDB := sqlx.NewDb(db, "sqlmock")
	repo := postgres.NewCartItemRepository(sqlxDB)

	mock.ExpectQuery(`SELECT EXISTS\(SELECT 1 FROM carts WHERE id = \$1\)`).
		WithArgs("cart-id").
		WillReturnRows(sqlmock.NewRows([]string{"exists"}).AddRow(true))

	exists, err := repo.CartExists(context.Background(), "cart-id")
	assert.NoError(t, err)
	assert.True(t, exists)

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestDeleteCartItem_CartNotFound(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to open mock database: %s", err)
	}
	defer db.Close()

	sqlxDB := sqlx.NewDb(db, "sqlmock")
	repo := postgres.NewCartItemRepository(sqlxDB)

	mock.ExpectQuery(`SELECT EXISTS\(SELECT 1 FROM carts WHERE id = \$1\)`).
		WithArgs("cart-id").
		WillReturnRows(sqlmock.NewRows([]string{"exists"}).AddRow(false))

	err = repo.Delete(context.Background(), "cart-id", "item-id")
	assert.Error(t, err)
	assert.Equal(t, "cart does not exist", err.Error())

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestDeleteCartItem_Success(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to open mock database: %s", err)
	}
	defer db.Close()

	sqlxDB := sqlx.NewDb(db, "sqlmock")
	repo := postgres.NewCartItemRepository(sqlxDB)

	mock.ExpectQuery(`SELECT EXISTS\(SELECT 1 FROM carts WHERE id = \$1\)`).
		WithArgs("cart-id").
		WillReturnRows(sqlmock.NewRows([]string{"exists"}).AddRow(true))

	mock.ExpectExec(`DELETE FROM cart_items WHERE id = \$1 AND cart_id = \$2`).
		WithArgs("item-id", "cart-id").
		WillReturnResult(sqlmock.NewResult(0, 1))

	err = repo.Delete(context.Background(), "cart-id", "item-id")
	assert.NoError(t, err)

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}
