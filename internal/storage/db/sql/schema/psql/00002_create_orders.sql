-- +goose Up
CREATE TABLE if NOT exists orders (
  id UUID PRIMARY KEY,
  number TEXT NOT NULL UNIQUE,
  user_id UUID NOT NULL,
  status TEXT NOT NULL,
  created_at TIMESTAMPTZ NOT NULL
);


-- +goose Down
DROP TABLE IF EXISTS ordrers;
