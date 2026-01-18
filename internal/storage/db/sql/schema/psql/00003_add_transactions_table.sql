-- +goose Up
CREATE TABLE IF NOT EXISTS transactions (
  id UUID PRIMARY KEY,
  user_id UUID NOT NULL REFERENCES users (id) ON DELETE CASCADE,
  order_id UUID NOT NULL REFERENCES orders (id) ON DELETE CASCADE,
  kind TEXT NOT NULL,
  amount INTEGER NOT NULL
);


-- +goose Down
DROP TABLE IF EXISTS transactions;
