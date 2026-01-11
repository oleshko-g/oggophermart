-- name: SelectOrdersByUserID :many
SELECT
  number,
  status,
  created_at
FROM
  orders
WHERE
  user_id = $1
ORDER BY
  created_at ASC;
