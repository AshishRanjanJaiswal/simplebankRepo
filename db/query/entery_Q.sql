-- name: CreateEnteries :one
INSERT INTO enteries (
  account_id,
  amount
) VALUES (
  $1, $2
) RETURNING *;

-- name: GetEnteries :one
SELECT * FROM enteries
WHERE id = $1 LIMIT 1;

-- name: ListEnteries :many
SELECT * FROM enteries
WHERE account_id = $1
ORDER BY id
LIMIT $2
OFFSET $3;