INSERT INTO
  users (
    id,
    login,
    hashed_password,
    created_at,
    updated_at
  )
VALUES
  ($1, $2, $3, $4, $5);
