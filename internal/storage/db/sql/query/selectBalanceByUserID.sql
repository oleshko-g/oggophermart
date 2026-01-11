-- name: SelectBalanceByUserID :one
SELECT
  current,
  withdrawn_sum
FROM
  user_balances
WHERE
  user_id = $1;
