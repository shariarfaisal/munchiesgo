

-- name: CreateRider :one
INSERT INTO riders (name, phone, hashed_password)
VALUES ($1, $2, $3) RETURNING *;

-- name: GetRider :one
SELECT * FROM riders WHERE id = $1;

-- name: ListRider :many
SELECT * FROM riders ORDER BY id Limit $1 OFFSET $2;