
-- name: CreateVendor :one
INSERT INTO vendors (
  name,
  email,
  phone,
  payment_info,
  social_links
) VALUES ($1, $2, $3, $4, $5)
RETURNING *;

-- name: UpdateVendor :one
UPDATE vendors SET
  name = $2,
  email = $3,
  phone = $4,
  payment_info = $5,
  social_links = $6
WHERE id = $1
RETURNING *;

-- name: DeleteVendor :exec
DELETE FROM vendors WHERE id = $1;

-- name: GetVendor :one
SELECT * FROM vendors WHERE id = $1;

-- name: ListVendors :many
SELECT * FROM vendors ORDER BY id LIMIT $1 OFFSET $2;

-- name: CountVendors :one
SELECT COUNT(*) FROM vendors;

-- name: SearchVendors :many
SELECT * FROM vendors
WHERE name ILIKE $1 OR email ILIKE $1 OR phone ILIKE $1 ORDER BY id LIMIT $2 OFFSET $3;





