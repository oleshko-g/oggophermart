-- name: SelectBalanceByUserID :one
SELECT
  *
FROM
  user_balances
WHERE
  user_id = $1;
