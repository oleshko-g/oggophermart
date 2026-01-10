-- name: InsertOrder :execresult
INSERT INTO
  orders (id, number, user_id, status, created_at)
VALUES
  ($1, $2, $3, $4, $5);
