-- +goose Up
CREATE TABLE IF NOT EXISTS orders (
  id UUID PRIMARY KEY,
  number TEXT NOT NULL,
  user_id UUID NOT NULL,
  status TEXT NOT NULL,
  created_at TIMESTAMPTZ NOT NULL,
  CONSTRAINT user_id_order_number UNIQUE (user_id, number)
);


-- +goose Down
DROP TABLE IF EXISTS orders;
