-- name: SelectUserHashedPasswordByLogin :one
SELECT
  hashed_password
FROM
  users
WHERE
  login = $1;
