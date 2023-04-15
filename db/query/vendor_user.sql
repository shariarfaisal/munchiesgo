
-- name: CreateVendorUser :one
INSERT INTO vendor_users (
  username,
  full_name,
  email,
  hashed_password,
  password_changed_at,
  role,
  vendor_id
) VALUES ($1, $2, $3, $4, $5, $6, $7)
RETURNING *;

-- name: UpdateVendorUser :one
UPDATE vendor_users SET
  username = $2,
  full_name = $3,
  email = $4,
  hashed_password = $5,
  password_changed_at = $6,
  role = $7,
  vendor_id = $8
WHERE id = $1 
RETURNING *;

-- name: DeleteVendorUser :exec
DELETE FROM vendor_users WHERE id = $1;

-- name: GetVendorUser :one
SELECT * FROM vendor_users WHERE id = $1;

-- name: GetVendorUserByUsername :one
SELECT * FROM vendor_users WHERE username = $1;

-- name: ListVendorUsers :many
SELECT * FROM vendor_users WHERE vendor_id = $1 ORDER BY id LIMIT $2 OFFSET $3;


--TODO: 
-- CREATE TABLE "vendor_permissions" (
--   "id" int PRIMARY KEY,
--   "name" varchar,
--   "created_at" timestamptz NOT NULL DEFAULT (now())
-- );

-- CREATE TABLE "vendor_user_permissions" (
--   "id" bigserial PRIMARY KEY,
--   "user_id" int NOT NULL,
--   "permission_id" int NOT NULL,
--   "created_at" timestamptz NOT NULL DEFAULT (now())
-- );

