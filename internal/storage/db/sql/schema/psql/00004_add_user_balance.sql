-- +goose Up
CREATE VIEW user_balances (user_id, accrued_sum, withdrawn_sum, current) AS
WITH
  cte_users AS (
    SELECT
      id AS user_id
    FROM
      users
  ),
  accruals AS (
    SELECT
      user_id,
      sum(amount) AS amount_sum
    FROM
      transactions
    WHERE
      kind = 'ACCRUAL'
    GROUP BY
      user_id
  ),
  withdrawals AS (
    SELECT
      user_id,
      sum(amount) AS amount_sum
    FROM
      transactions
    WHERE
      kind = 'WITHDRAWAL'
    GROUP BY
      user_id
  ),
  balances AS (
    SELECT
      cte_users.user_id,
      COALESCE(accruals.amount_sum, 0) AS accrued_sum,
      COALESCE(withdrawals.amount_sum, 0) AS withdrawn_sum,
      COALESCE(accruals.amount_sum, 0) - COALESCE(withdrawals.amount_sum, 0) AS current
    FROM
      cte_users
      LEFT JOIN accruals USING (user_id)
      LEFT JOIN withdrawals USING (user_id)
  )
SELECT
  user_id,
  accrued_sum,
  withdrawn_sum,
  accrued_sum - withdrawn_sum AS current
FROM
  balances;


-- +goose Down
DROP VIEW user_balances;
