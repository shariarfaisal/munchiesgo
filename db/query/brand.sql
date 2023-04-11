
-- name: CreateBrand :one
INSERT INTO brands (
  name,
  meta_tags,
  slug,
  type,
  phone,
  email,
  email_verified,
  logo,
  banner,
  rating,
  vendor_id,
  prefix,
  status,
  availability,
  location,
  address
) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16)
RETURNING *;

-- name: UpdateBrand :one
UPDATE brands SET
  name = $2,
  meta_tags = $3,
  slug = $4,
  type = $5,
  phone = $6,
  email = $7,
  email_verified = $8,
  logo = $9,
  banner = $10,
  rating = $11,
  vendor_id = $12,
  prefix = $13,
  status = $14,
  availability = $15,
  location = $16,
  address = $17
WHERE id = $1
RETURNING *;

-- name: DeleteBrand :exec
DELETE FROM brands WHERE id = $1;

-- name: GetBrand :one
SELECT * FROM brands WHERE id = $1;

-- name: ListBrands :many
SELECT * FROM brands ORDER BY id LIMIT $1 OFFSET $2;

-- name: ListBrandsByVendorID :many
SELECT * FROM brands WHERE vendor_id = $1 ORDER BY id LIMIT $2 OFFSET $3;


-- name: CreateOperationTime :one
INSERT INTO operation_times (
  brand_id,
  day_of_week,
  start_time,
  end_time
) VALUES ($1, $2, $3, $4)
RETURNING *;

-- name: UpdateOperationTime :one
UPDATE operation_times SET
  brand_id = $2,
  day_of_week = $3,
  start_time = $4,
  end_time = $5
WHERE id = $1
RETURNING *;

-- name: DeleteOperationTime :exec
DELETE FROM operation_times WHERE id = $1;

-- name: GetOperationTime :one
SELECT * FROM operation_times WHERE id = $1;

-- name: ListOperationTimesByBrandId :many
SELECT * FROM operation_times WHERE brand_id = $1 ORDER BY id LIMIT $2 OFFSET $3;

