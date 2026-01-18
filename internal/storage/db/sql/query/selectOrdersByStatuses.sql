-- name: SelectOrdersIDsByStatuses :many
SELECT
  orders.id
FROM
  orders
WHERE
  orders.status = ANY(sqlc.arg(statuses)::TEXT[]);
