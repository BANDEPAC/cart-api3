-- +goose Up
CREATE TABLE cart_items (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    cart_id UUID REFERENCES carts(id) ON DELETE CASCADE,
    product TEXT NOT NULL,
    quantity INT NOT NULL CHECK (quantity > 0)
);

-- +goose Down
DROP TABLE cart_items;