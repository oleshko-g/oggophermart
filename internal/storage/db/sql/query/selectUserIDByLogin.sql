-- name: SelectUserIDByLogin :one
SELECT
  id
FROM
  users
WHERE
  login = $1;
