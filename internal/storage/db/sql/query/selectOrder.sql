-- name: SelectOrder :one
SELECT
  *
FROM
  orders
WHERE
  id = $1
FOR UPDATE;
