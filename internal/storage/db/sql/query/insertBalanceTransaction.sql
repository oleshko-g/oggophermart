-- name: InsertBalanceTransaction :exec
INSERT INTO
  transactions (id, user_id, order_id, kind, amount)
VALUES
  ($1, $2, $3, $4, $5);
