-- name: SelectOrdersIDsByStatuses :many
SELECT
  orders.id
FROM
  orders
WHERE
  orders.status = ALL(sqlc.arg(statuses)::TEXT[]);
