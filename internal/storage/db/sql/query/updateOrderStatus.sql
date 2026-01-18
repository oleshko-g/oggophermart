-- name: UpdateOrderStatus :exec
UPDATE orders
SET
  status = $2
WHERE
  id = $1;
