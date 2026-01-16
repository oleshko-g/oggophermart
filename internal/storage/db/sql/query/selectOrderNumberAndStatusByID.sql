-- name: SelectOrderNumberAndStatusByID :one
SELECT
  number,
  status
FROM
  orders
WHERE
  id = $1
FOR UPDATE;
