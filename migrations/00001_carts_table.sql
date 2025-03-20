-- +goose Up
CREATE TABLE carts (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid()
);

-- +goose Down
DROP TABLE carts;